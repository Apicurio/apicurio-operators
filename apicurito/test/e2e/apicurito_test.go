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
	err := framework.AddToFrameworkScheme(apis.AddToScheme, apicuritoList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	// run subtests
	t.Run("apicurito-group", func(t *testing.T) {
		t.Run("Cluster1", ApicuritoCluster)
	})
}

func apicuritoScaleTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create apicurito custom resource
	exampleApicurito := &apicuritosv1alpha1.Apicurito{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Apicurito",
			APIVersion: "apicur.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service-apicurito",
			Namespace: namespace,
		},
		Spec: apicuritosv1alpha1.ApicuritoSpec{
			Size:  3,
			Image: "apicurio/apicurito-ui:latest",
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), exampleApicurito, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	// wait for example-apicurito to reach 3 replicas
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "test-service-apicurito", 3, retryInterval, timeout)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "test-service-apicurito", Namespace: namespace}, exampleApicurito)
	if err != nil {
		return err
	}

	exampleApicurito.Spec.Size = 4
	err = f.Client.Update(goctx.TODO(), exampleApicurito)
	if err != nil {
		return err
	}

	// wait for example-apicurito to reach 4 replicas
	return e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "test-service-apicurito", 4, retryInterval, timeout)
}

func ApicuritoCluster(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()

	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")

	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}

	// get global framework variables
	f := framework.Global

	// wait for apicurito-operator to be ready
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "apicurito", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = apicuritoScaleTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
