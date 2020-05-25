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
	"time"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/resources"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"

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
	err = c.Watch(&source.Kind{Type: &v1alpha1.Apicurito{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Apicurito
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Apicurito{},
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
func (r *ReconcileApicurito) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Apicurito.")

	// Fetch the Apicurito instance
	apicurito := &v1alpha1.Apicurito{}
	err := r.client.Get(context.TODO(), request.NamespacedName, apicurito)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not fd, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Apicurito resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Apicurito.")
		return reconcile.Result{}, err
	}

	// Ensure the deployment image is the same as the one from configuration
	c := &configuration.Config{}
	if err = c.Config(apicurito); err != nil {
		reqLogger.Error(err, "failed to generate configuration")
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: 10 * time.Second,
		}, err
	}

	res := resources.Resource{
		Client:    r.client,
		Apicurito: apicurito,
		Cfg:       c,
		Logger:    reqLogger,
	}
	if err := res.Create(); err != nil {
		reqLogger.Error(err, "failed to create resources", "Namespace", res.Apicurito.Namespace)
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: 10 * time.Second,
		}, err
	}

	return reconcile.Result{
		Requeue:      true,
		RequeueAfter: 20 * time.Second,
	}, nil
}
