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
	"fmt"
	"testing"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"

	apicuritosv1alpha1 "github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/restmapper"

	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	// 	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const succeed = "\u2713"
const failed = "\u2717"

// TestApicuritoController runs ReconcileApicurito.Reconcile() against a
// fake client that tracks a apicurito object.
func TestApicuritoController(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	// logf.SetLogger(logf.ZapLogger(true))

	var (
		name              = "apicurito-operator"
		namespace         = "apicurito"
		replicas    int32 = 3
		metricsHost       = "0.0.0.0"
		metricsPort int32 = 8383
	)

	// An apicurito resource with metadata and spec.
	apicurito := &apicuritosv1alpha1.Apicurito{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: apicuritosv1alpha1.ApicuritoSpec{
			Size: replicas,
		},
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		t.Fatalf("\t%s\tFailed to get config: (%v)", failed, err)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		Namespace:          namespace,
		MapperProvider:     restmapper.NewDynamicRESTMapper,
		MetricsBindAddress: fmt.Sprintf("%s:%d", metricsHost, metricsPort),
	})
	if err != nil {
		t.Fatalf("\t%s\tFailed to create new cmd to provide shared dependencies: (%v)", failed, err)
	}

	if err := routev1.AddToScheme(mgr.GetScheme()); err != nil {
		t.Fatalf("\t%s\tFailed to add openshift route to scheme: (%v)", failed, err)
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		apicurito,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(apicuritosv1alpha1.SchemeGroupVersion, apicurito)

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClientWithScheme(s, objs...)

	// Create a ReconcileApicurito object with the scheme and fake client.
	r := &ReconcileApicurito{client: cl, scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	configuration.ConfigFile = "../../../build/conf/config_test.yaml"

	{
		t.Logf("\tTest 0\tWhen simulating reconcile the first time.")
		res, err := r.Reconcile(req)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to reconcile but got error: (%v)", failed, err)
		}
		t.Logf("\t%s\tShould be able to reconcile.", succeed)

		// Check the result of reconciliation to make sure it has the desired state.
		if !res.Requeue {
			t.Errorf("\t%s\tShould be able to requeue reconcile.", failed)
		} else {
			t.Logf("\t%s\tShould be able to requeue reconcile.", succeed)
		}
	}
}
