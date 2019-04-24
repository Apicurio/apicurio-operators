// Copyright 2018 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	apis "github.com/apicurio/apicurio-operators/apicurito/pkg/apis"
	apicuritosv1alpha1 "github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestApicurito(t *testing.T) {
	apicuritoList := &apicuritosv1alpha1.ApicuritoList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Apicurito",
			APIVersion: "apicur.io/v1alpha1",
		},
	}
	if err := framework.AddToFrameworkScheme(apis.AddToScheme, apicuritoList); err != nil {
		t.Fatalf("Failed to add custom resource scheme to framework: %v", err)
	}

	// run subtests
	t.Run("apicurito-group", func(t *testing.T) {
		t.Run("Cluster1", ApicuritoCluster)
	})
}

func apicuritoScaleTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	n, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("Could not get namespace: %v", err)
	}

	// apicurito custom resource
	tn := "test-scale-apicurito"
	cr := apicuritosv1alpha1.Apicurito{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Apicurito",
			APIVersion: "apicur.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tn,
			Namespace: n,
		},
		Spec: apicuritosv1alpha1.ApicuritoSpec{
			Size:  3,
			Image: "apicurio/apicurito-ui:latest",
		},
	}

	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	t.Logf("Creating Apicurito CR with replicas Apicurito.Size (%d)", cr.Spec.Size)
	co := framework.CleanupOptions{
		TestContext:   ctx,
		Timeout:       cleanupTimeout,
		RetryInterval: cleanupRetryInterval,
	}
	if err := f.Client.Create(goctx.TODO(), &cr, &co); err != nil {
		return err
	}

	// wait for apicurito to reach 3 replicas
	t.Logf("Waiting for Apicurito Deployment to reach (%d) replicas", cr.Spec.Size)
	err = e2eutil.WaitForDeployment(t, f.KubeClient, n, tn, 3, retryInterval, timeout)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: tn, Namespace: n}, &cr)
	if err != nil {
		t.Errorf("Unable to get Apicurito CR, %v", err)
		return err
	}

	cr.Spec.Size = 4
	t.Logf("Modifying Apicurito CR with replicas Apicurito.Size (%d)", cr.Spec.Size)
	if err := f.Client.Update(goctx.TODO(), &cr); err != nil {
		t.Errorf("Unable to update Apicurito CR, %v", err)
		return err
	}

	// wait for example-apicurito to reach 4 replicas
	if err := e2eutil.WaitForDeployment(t, f.KubeClient, n, tn, 4, retryInterval, timeout); err != nil {
		return err
	}

	return nil
}

func apicuritoUpdateTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	n, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("Could not get namespace: %v", err)
	}

	// apicurito custom resource
	tn := "test-image-apicurito"
	cr := apicuritosv1alpha1.Apicurito{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Apicurito",
			APIVersion: "apicur.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tn,
			Namespace: n,
		},
		Spec: apicuritosv1alpha1.ApicuritoSpec{
			Size:  1,
			Image: "apicurio/apicurito-ui:latest",
		},
	}

	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	t.Logf("Creating Apicurito CR with image Apicurito.Spec.Image (%s)", cr.Spec.Image)
	co := framework.CleanupOptions{
		TestContext:   ctx,
		Timeout:       cleanupTimeout,
		RetryInterval: cleanupRetryInterval,
	}
	if err := f.Client.Create(goctx.TODO(), &cr, &co); err != nil {
		return err
	}

	// wait for apicurito to reach 1 replicas
	t.Logf("Waiting for Apicurito Deployment to reach (%d) replicas", cr.Spec.Size)
	err = e2eutil.WaitForDeployment(t, f.KubeClient, n, tn, 1, retryInterval, timeout)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: tn, Namespace: n}, &cr)
	if err != nil {
		t.Errorf("Unable to get Apicurito CR, %v", err)
		return err
	}

	ni := "apicurio/apicurito-ui:1.0.1"
	cr.Spec.Image = ni
	t.Logf("Modifying Apicurito CR with image Apicurito.Spec.Image (%s)", cr.Spec.Image)
	if err := f.Client.Update(goctx.TODO(), &cr); err != nil {
		t.Errorf("Unable to update Apicurito CR, %v", err)
		return err
	}

	// wait for example-apicurito to reach 1 replicas
	if err := e2eutil.WaitForDeployment(t, f.KubeClient, n, tn, 1, retryInterval, timeout); err != nil {
		return err
	}

	dep := appsv1.Deployment{}
	if err := f.Client.Get(goctx.TODO(), types.NamespacedName{Name: tn, Namespace: n}, &dep); err != nil {
		t.Errorf("Unable to get Apicurito Deployment, %v", err)
		return err
	}

	di := dep.Spec.Template.Spec.Containers[0].Image
	if di != ni {
		t.Errorf("Deployment should have image (%s) but got (%s).", ni, di)
	}
	t.Logf("Deployment should have image (%s) and got (%s).", ni, di)

	return nil
}

// Create resources from file in namespace and run
// integration tests
func ApicuritoCluster(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()

	co := framework.CleanupOptions{
		TestContext:   ctx,
		Timeout:       cleanupTimeout,
		RetryInterval: cleanupRetryInterval,
	}

	// Initialize cluster resources from yaml, including operator deployment
	// and needed roles / rbac
	if err := ctx.InitializeClusterResources(&co); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources from yaml")

	n, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}

	// get global framework variables
	f := framework.Global

	// wait for apicurito-operator to be ready
	if err := e2eutil.WaitForDeployment(t, f.KubeClient, n, "apicurito", 1, retryInterval, timeout); err != nil {
		t.Fatal(err)
	}

	if err = apicuritoScaleTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}

	if err = apicuritoUpdateTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
