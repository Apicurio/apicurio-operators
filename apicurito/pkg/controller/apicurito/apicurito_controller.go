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

package apicurito

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"

	"github.com/RHsyseng/operator-utils/pkg/resource/compare"
	"github.com/RHsyseng/operator-utils/pkg/resource/read"
	"github.com/RHsyseng/operator-utils/pkg/resource/write"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/resources"

	pkg "github.com/apicurio/apicurio-operators/apicurito/pkg"
	api "github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_apicurito")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Apicurito Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	v := &ReconcileApicurito{client: mgr.GetClient(), scheme: mgr.GetScheme()}
	if err := ConsoleYAMLSampleExists(); err == nil {
		createConsoleYAMLSamples(v.client)
	}
	return v
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("apicurito-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Apicurito
	err = c.Watch(&source.Kind{Type: &api.Apicurito{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Apicurito
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Apicurito{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileApicurito{}

// ReconcileApicurito reconciles a Apicurito object
type ReconcileApicurito struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Apicurito object and makes changes based on the state read
// and what is in the Apicurito.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileApicurito) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Apicurito.")

	// Fetch the Apicurito instance
	apicurito := &api.Apicurito{}
	err := r.client.Get(ctx, request.NamespacedName, apicurito)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not fd, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Apicurito resource not found. Ignoring since object must be deleted.")

			if err := consoleLinkExists(); err == nil {
				apicurito.ObjectMeta = metav1.ObjectMeta{
					Name:      request.Name,
					Namespace: request.Namespace,
				}
				removeConsoleLink(r.client, apicurito)
			}
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Apicurito.")
		return reconcile.Result{}, err
	}

	if apicurito.Status.Phase == api.ApicuritoPhaseMissing {
		r.updateStatus(ctx, apicurito, api.ApicuritoPhaseStarting, reqLogger)
	}

	c := &configuration.Config{}
	if err = c.Config(apicurito); err != nil {
		reqLogger.Error(err, "failed to generate configuration")
		r.updateStatus(ctx, apicurito, api.ApicuritoPhaseInstallError, reqLogger)
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: 10 * time.Second,
		}, err
	}

	var rs resources.Generator = resources.Resource{
		Client:    r.client,
		Apicurito: apicurito,
		Cfg:       c,
		Logger:    reqLogger,
	}

	if apicurito.Status.Phase == api.ApicuritoPhaseStarting || apicurito.Status.Phase == api.ApicuritoPhaseInstallError {
		r.updateStatus(ctx, apicurito, api.ApicuritoPhaseInstalling, reqLogger)
	}

	// Fetch routes resources and apply them before the rest
	// This is needed because ConfigMaps require the routes to be present and should run only once
	// at startup
	route := &routev1.Route{}
	err = r.client.Get(ctx, types.NamespacedName{Name: fmt.Sprintf("%s-%s", apicurito.Name, "generator"), Namespace: apicurito.Namespace}, route)
	if err != nil && errors.IsNotFound(err) {
		routes := rs.Routes()
		err = r.applyResources(apicurito, routes, reqLogger)
		if err != nil {
			reqLogger.Info("Apicurito CR resource changed in the meantime, requeue and rerun in 10 seconds", "err", err)
			r.updateStatus(ctx, apicurito, api.ApicuritoPhaseInstallError, reqLogger)
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: 10 * time.Second,
			}, err
		}

		time.Sleep(5 * time.Second)
	}
	if err := consoleLinkExists(); err == nil {
		createConsoleLink(r.client, apicurito)
	}
	// generate all resources and apply them
	res, err := rs.Generate()
	if err != nil {
		reqLogger.Error(err, "failed to generate resources")
		r.updateStatus(ctx, apicurito, api.ApicuritoPhaseInstallError, reqLogger)
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: 10 * time.Second,
		}, err
	}
	err = r.applyResources(apicurito, res, reqLogger)
	if err != nil {
		reqLogger.Info("Apicurito CR changed in the meantime, requeue and rerun in 10 seconds", "err", err)
		r.updateStatus(ctx, apicurito, api.ApicuritoPhaseInstallError, reqLogger)
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: 10 * time.Second,
		}, err
	}

	r.updateStatus(ctx, apicurito, api.ApicuritoPhaseInstalled, reqLogger)

	return reconcile.Result{
		Requeue:      true,
		RequeueAfter: 20 * time.Second,
	}, nil
}

func (r *ReconcileApicurito) updateStatus(ctx context.Context, apicurito *api.Apicurito, phase api.ApicuritoPhase, logger logr.Logger) {
	target := apicurito.DeepCopy()
	target.Status.Phase = phase
	target.Status.Version = pkg.Version
	err := r.client.Status().Update(ctx, target)
	time.Sleep(3 * time.Second)
	if err != nil {
		logger.Info("Failed to update apicurito status", "err", err)
	}
}

func (r *ReconcileApicurito) applyResources(apicurito *api.Apicurito, res []client.Object, logger logr.Logger) (err error) {
	deployed, err := getDeployedResources(apicurito, r.client)

	requested := compare.NewMapBuilder().Add(res...).ResourceMap()
	comparator := getComparator()
	deltas := comparator.Compare(deployed, requested)
	writer := write.New(r.client).WithOwnerController(apicurito, r.scheme)

	for resourceType, delta := range deltas {
		if !delta.HasChanges() {
			continue
		}

		logger.Info("", "instances of ", resourceType, "Will create ", len(delta.Added), "update ", len(delta.Updated), "and delete", len(delta.Removed))

		_, err := writer.AddResources(delta.Added)
		if err != nil {
			return fmt.Errorf("Apicurito CR changed: AddResources: %s", err)
		}

		_, err = writer.UpdateResources(deployed[resourceType], delta.Updated)
		if err != nil {
			return fmt.Errorf("Apicurito CR changed: UpdateResources : %s", err)
		}

		_, err = writer.RemoveResources(delta.Removed)
		if err != nil {
			return fmt.Errorf("Apicurito CR changed: RemoveResources: %s", err)
		}

	}

	return
}

func getDeployedResources(cr *api.Apicurito, client client.Client) (map[reflect.Type][]client.Object, error) {
	var log = logf.Log.WithName("getDeployedResources")

	reader := read.New(client).WithNamespace(cr.Namespace).WithOwnerObject(cr)
	resourceMap, err := reader.ListAll(
		&corev1.ConfigMapList{},
		&corev1.ServiceList{},
		&appsv1.DeploymentList{},
		&routev1.RouteList{},
		&corev1.ServiceAccountList{},
	)
	if err != nil {
		log.Error(err, "Failed to list deployed objects. ", err)
		return nil, err
	}

	return resourceMap, nil

}

func getComparator() compare.MapComparator {
	resourceComparator := compare.DefaultComparator()

	configMapType := reflect.TypeOf(corev1.ConfigMap{})
	resourceComparator.SetComparator(configMapType, func(deployed client.Object, requested client.Object) bool {
		configMap1 := deployed.(*corev1.ConfigMap)
		configMap2 := requested.(*corev1.ConfigMap)
		var pairs [][2]interface{}
		pairs = append(pairs, [2]interface{}{configMap1.Name, configMap2.Name})
		pairs = append(pairs, [2]interface{}{configMap1.Namespace, configMap2.Namespace})
		pairs = append(pairs, [2]interface{}{configMap1.Labels, configMap2.Labels})
		pairs = append(pairs, [2]interface{}{configMap1.Annotations, configMap2.Annotations})
		pairs = append(pairs, [2]interface{}{configMap1.Data, configMap2.Data})
		pairs = append(pairs, [2]interface{}{configMap1.BinaryData, configMap2.BinaryData})
		equal := compare.EqualPairs(pairs)
		if !equal {
			log.Info("Resources are not equal", "deployed", deployed, "requested", requested)
		}
		return equal
	})

	return compare.MapComparator{Comparator: resourceComparator}
}
