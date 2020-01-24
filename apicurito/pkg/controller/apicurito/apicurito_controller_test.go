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
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"

	apicuritosv1alpha1 "github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/restmapper"

	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
		replicasCh  int32 = 1
		image             = "apicurio/apicurito-ui"
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

	// Check if deployment has been created and is correct.
	{
		t.Logf("\tTest 1\tWhen the deployment is created the first time.")

		dep := &appsv1.Deployment{}
		err := cl.Get(context.TODO(), req.NamespacedName, dep)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to get deployment but got error: (%v)", failed, err)
		}
		t.Logf("\t%s\tShould be able to get deployment.", succeed)

		dsize := *dep.Spec.Replicas
		if dsize != replicas {
			t.Errorf("\t%s\tDeployment should have replica size of (%d) but got (%d).", failed, replicas, dsize)
		} else {
			t.Logf("\t%s\tDeployment should have replica size of (%d) and got (%d).", succeed, replicas, dsize)
		}

		di := dep.Spec.Template.Spec.Containers[0].Image
		if di != image {
			t.Errorf("\t%s\tDeployment should have image (%s) but got (%s).", failed, image, di)
		} else {
			t.Logf("\t%s\tDeployment should have image (%s) and got (%s).", succeed, image, di)
		}

		ser := &corev1.Service{}
		if err = cl.Get(context.TODO(), req.NamespacedName, ser); err != nil {
			t.Errorf("\t%s\tShould be able to get service but got error: (%v)", failed, err)
		} else {
			t.Logf("\t%s\tShould be able to get service.", succeed)
		}
	}

	// Check that Reconcile update the Apicurito node list with some fake pod names
	{
		t.Logf("\tTest 2\tWhen pods are created, Apicurito CR should update with nodes names.")

		// Create the 3 expected pods in namespace and collect their names to check later
		podLabels := labelsForApicurito(name)
		pod := corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Labels:    podLabels,
			},
		}

		podNames := make([]string, 3)
		for i := 0; i < 3; i++ {
			pod.ObjectMeta.Name = name + ".pod." + strconv.Itoa(rand.Int())
			podNames[i] = pod.ObjectMeta.Name
			if err := cl.Create(context.TODO(), pod.DeepCopy()); err != nil {
				t.Fatalf("create pod %d: (%v)", i, err)
			}
		}

		// Reconcile again so Reconcile() checks pods and updates the apicurito
		// resources' Status.
		res, err := r.Reconcile(req)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to reconcile but got error: (%v)", failed, err)
		}
		if res != (reconcile.Result{}) {
			t.Errorf("\t%s\tReconcile should return an empty Result but got (%v)", failed, res)
		} else {
			t.Logf("\t%s\tReconcile should return an empty Result", succeed)
		}

		// Get the updated apicurito object.
		apicurito = &apicuritosv1alpha1.Apicurito{}
		err = r.client.Get(context.TODO(), req.NamespacedName, apicurito)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to get apicurito object but got an error instead: (%v)", failed, err)
		}

		// Ensure Reconcile() updated the apicurito's Status as expected.
		nodes := apicurito.Status.Nodes
		if !reflect.DeepEqual(podNames, nodes) {
			t.Errorf("\t%s\tPod names should be the same, but they dont match, got (%v), want (%v)", failed, podNames, nodes)
		} else {
			t.Logf("\t%s\tPod names should be the same.", succeed)
		}
	}

	// Change Apicurito after the Deployment exist and make sure deployment is
	// being
	{
		t.Logf("\tTest 3\tWhen changing the apicurito CR the deployment should reflect those changes.")

		apicurito.Spec.Size = replicasCh

		err := cl.Update(context.TODO(), apicurito)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to update Apicurito CR after changing CR: (%v)", failed, err)
		}

		res, err := r.Reconcile(req)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to reconcile after changing CR but got error: (%v)", failed, err)
		}
		if res != (reconcile.Result{Requeue: true}) {
			t.Errorf("\t%s\tAfter changing the replica size in CR, Reconcile should requeue, but got (%v)", failed, res)
		} else {
			t.Logf("\t%s\tAfter changing the replica size in CR, Reconcile should requeue.", succeed)
		}

		dep := &appsv1.Deployment{}
		err = cl.Get(context.TODO(), req.NamespacedName, dep)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to get deployment after changing CR but got error: (%v)", failed, err)
		}

		dsize := *dep.Spec.Replicas
		if dsize != replicasCh {
			t.Errorf("\t%s\tAfter changing the replica size in CR, Deployment should change, want (%d) but got (%d).", failed, replicasCh, dsize)
		} else {
			t.Logf("\t%s\tAfter changing the replica size in CR, Deployment should change, want (%d) and got (%d).", succeed, replicasCh, dsize)
		}
	}
}
