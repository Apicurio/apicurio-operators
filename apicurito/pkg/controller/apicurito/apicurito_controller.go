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
	"reflect"
	"time"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"

	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
	} else {
		log.Info("Yaml Samples", "Console YAML sample is not added:", err)
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
			reqLogger.Info("Apicurito resource not fd. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Apicurito.")
		return reconcile.Result{}, err
	}

	// Check if the service exist and create it otherwise
	sf := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: apicurito.Name, Namespace: apicurito.Namespace}, sf)
	if err != nil && errors.IsNotFound(err) {
		// Define new service
		ser := r.serviceForApicurito(apicurito)
		reqLogger.Info("Creating a new Service.", "Service.Namespace", ser.Namespace, "Service.Name", ser.Name)
		err = r.client.Create(context.TODO(), ser)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service.", "Service.Namespace", ser.Namespace, "Service.Name", ser.Name)
			return reconcile.Result{}, err
		}
	}

	// Check if the deployment already exists, if not create a new one
	fd := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: apicurito.Name, Namespace: apicurito.Namespace}, fd)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep, err := r.deploymentForApicurito(apicurito)
		if err != nil {
			return reconcile.Result{}, err
		}
		reqLogger.Info("Creating a new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment.")
		return reconcile.Result{}, err
	}

	// Check if route already exists, if not create a new one
	rf := &routev1.Route{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: apicurito.Name, Namespace: apicurito.Namespace}, rf)
	if err != nil && errors.IsNotFound(err) {
		// Define new route
		route := r.routeForApicurito(apicurito)
		reqLogger.Info("Creating a new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
		err = r.client.Create(context.TODO(), route)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Route.", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
			return reconcile.Result{}, err
		}

		// Route takes some time to come up, let's give it 5s to come up
		reqLogger.Info("Route created, waiting 5s")
		time.Sleep(5 * time.Second)
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Route.")
		return reconcile.Result{}, err
	}

	// Ensure the deployment image is the same as the one from configuration
	c := &configuration.Config{}
	if err = c.Config(apicurito); err != nil {
		return reconcile.Result{}, err
	}
	image := c.Image
	if fd.Spec.Template.Spec.Containers[0].Image != image {
		fd.Spec.Template.Spec.Containers[0].Image = image
		err = r.client.Update(context.TODO(), fd)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment.", "Deployment.Namespace", fd.Namespace, "Deployment.Name", fd.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	// Ensure the deployment size are the same as the spec
	size := apicurito.Spec.Size
	if *fd.Spec.Replicas != size {
		fd.Spec.Replicas = &size
		err = r.client.Update(context.TODO(), fd)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment.", "Deployment.Namespace", fd.Namespace, "Deployment.Name", fd.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// Update the Apicurito status with the pod names
	// List the pods for this apicurito's deployment
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForApicurito(apicurito.Name))
	listOps := &client.ListOptions{
		Namespace:     apicurito.Namespace,
		LabelSelector: labelSelector,
	}

	err = r.client.List(context.TODO(), podList, listOps)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods.", "Apicurito.Namespace", apicurito.Namespace, "apicurito.Name", apicurito.Name)
		return reconcile.Result{}, err
	}

	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, apicurito.Status.Nodes) {
		apicurito.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), apicurito)
		if err != nil {
			reqLogger.Error(err, "Failed to update Apicurito status.")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// deploymentForApicurito returns a apicurito Deployment object
func (r *ReconcileApicurito) deploymentForApicurito(m *v1alpha1.Apicurito) (*appsv1.Deployment, error) {
	c := &configuration.Config{}
	if err := c.Config(m); err != nil {
		return nil, err
	}

	ls := labelsForApicurito(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           c.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            "apicurito",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "api-port",
						}},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Scheme: corev1.URISchemeHTTP,
									Port:   intstr.FromString("api-port"),
									Path:   "/",
								}},
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Scheme: corev1.URISchemeHTTP,
									Port:   intstr.FromString("api-port"),
									Path:   "/",
								}},
							PeriodSeconds:    5,
							FailureThreshold: 2,
						},
					}},
				},
			},
		},
	}

	// Set Apicurito instance as the owner and controller
	if err := controllerutil.SetControllerReference(m, dep, r.scheme); err != nil {
		return nil, err
	}

	return dep, nil
}

// labelsForApicurito returns the labels for selecting the resources
// belonging to the given apicurito CR name.
func labelsForApicurito(name string) map[string]string {
	return map[string]string{"app": "apicurito", "apicurito_cr": name}
}

// serviceForApicurito returns an apicurito Service
func (r *ReconcileApicurito) serviceForApicurito(a *v1alpha1.Apicurito) *corev1.Service {
	ls := labelsForApicurito(a.Name)

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: a.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Name: "api-port",
					Port: 8080,
				},
			},
		},
	}

	return service
}

func (r *ReconcileApicurito) routeForApicurito(a *v1alpha1.Apicurito) *routev1.Route {
	ls := labelsForApicurito(a.Name)
	route := routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: routev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: a.Namespace,
			Labels:    ls,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(a, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    a.Kind,
				}),
			},
		},
		Spec: routev1.RouteSpec{
			// Host: a.Spec.Route,
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: a.Name,
			},
		},
	}

	return &route
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
