/*
 * Copyright (C) 2020 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"flag"
	"fmt"
	"runtime"

	kubemetrics "github.com/operator-framework/operator-sdk/pkg/kube-metrics"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/apicurio/apicurio-operators/apicurito/version"
	customMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis"
	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	"github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"
	"github.com/apicurio/apicurio-operators/apicurito/pkg/controller"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"github.com/operator-framework/operator-sdk/pkg/leader"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/operator-framework/operator-sdk/pkg/metrics"
	"github.com/operator-framework/operator-sdk/pkg/restmapper"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

type options struct {
	*Options
}

func newRunCommand(parent *Options) *cobra.Command {
	options := options{Options: parent}
	cmd := cobra.Command{
		Use:   "run",
		Short: "runs the operator",
		Run: func(_ *cobra.Command, _ []string) {
			exitOnError(options.run())
		},
	}

	cmd.PersistentFlags().StringVarP(&configuration.ConfigFile, "config", "", "/conf/config.yaml", "path to the operator configuration file.")
	cmd.PersistentFlags().AddFlagSet(zap.FlagSet())
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return &cmd
}

var (
	metricsHost               = "0.0.0.0"
	metricsPort         int32 = 8383
	operatorMetricsPort int32 = 8686
	operatorVersion           = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "apicurito_version_info",
			Help:        "Apicurito operator information",
			ConstLabels: prometheus.Labels{"operator_version": version.Version},
		},
	)
)

var log = logf.Log.WithName("cmd")

func init() {
	// Register custom metrics with the global prometheus registry
	customMetrics.Registry.MustRegister(operatorVersion)
}

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info(fmt.Sprintf("Version of operator-sdk: %v", sdkVersion.Version))
	log.Info(fmt.Sprintf("Version of apicurito operator: %v", version.Version))
}

func (o *options) run() error {

	// pflag.Parse()

	// Use a zap logr.Logger implementation. If none of the zap
	// flags are configured (or if the zap flag set is not being
	// used), this defaults to a production zap logger.
	//
	// The logger instantiated here can be changed to any logger
	// implementing the logr.Logger interface. This logger will
	// be propagated through the whole operator, generating
	// uniform and structured logs.
	logf.SetLogger(zap.Logger())

	printVersion()

	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		return err
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	ctx := o.Context

	// Become the leader before proceeding
	err = leader.Become(ctx, "apicurito-lock")
	if err != nil {
		return err
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		Namespace:          namespace,
		MapperProvider:     restmapper.NewDynamicRESTMapper,
		MetricsBindAddress: fmt.Sprintf("%s:%d", metricsHost, metricsPort),
	})
	if err != nil {
		return err
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		return err
	}

	if err := routev1.AddToScheme(mgr.GetScheme()); err != nil {
		exitOnError(err)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		exitOnError(err)
	}

	// Setup metrics. Serves Operator/CustomResource GVKs and generates metrics based on those types
	installationGVK := []schema.GroupVersionKind{v1alpha1.SchemaGroupVersionKind}

	// To generate metrics in other namespaces, add the values below.
	ns := []string{namespace}

	// Generate and serve custom resource specific metrics.
	err = kubemetrics.GenerateAndServeCRMetrics(cfg, ns, installationGVK, metricsHost, operatorMetricsPort)
	if err != nil {
		return err
	}

	// Create Service object to expose the metrics port.
	servicePorts := []v1.ServicePort{
		{Port: metricsPort, Name: metrics.OperatorPortName, Protocol: v1.ProtocolTCP, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: metricsPort}},
		{Port: operatorMetricsPort, Name: metrics.CRPortName, Protocol: v1.ProtocolTCP, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: operatorMetricsPort}},
	}

	service, err := metrics.CreateMetricsService(ctx, cfg, servicePorts)
	if err != nil {
		log.Info(err.Error())
		log.Info("Could not create metrics Service", "error", err.Error())
	}

	services := []*v1.Service{service}
	_, err = metrics.CreateServiceMonitors(cfg, namespace, services)
	if err != nil {
		log.Info("Could not create ServiceMonitor object", "error", err.Error())
		// If this operator is deployed to a cluster without the prometheus-operator running, it will return
		// ErrServiceMonitorNotPresent, which can be used to safely skip ServiceMonitor creation.
		if err == metrics.ErrServiceMonitorNotPresent {
			log.Info("Install prometheus-operator in your cluster to create ServiceMonitor objects", "error", err.Error())
		}
	}

	log.Info("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		return err
	}

	return nil
}
