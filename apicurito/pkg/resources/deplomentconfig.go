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

package resources

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	"github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Creates and returns a apicurito Deployment object
func apicuritoDeployment(client client.Client, c *configuration.Config, a *v1alpha1.Apicurito) (dep *appsv1.Deployment, err error) {
	dep = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", a.Name, "ui"),
		},
	}
	err = client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-%s", a.Name, "ui"), Namespace: a.Namespace}, dep)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep = &appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", a.Name, "ui"),
				Namespace: a.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(a, schema.GroupVersionKind{
						Group:   v1alpha1.SchemeGroupVersion.Group,
						Version: v1alpha1.SchemeGroupVersion.Version,
						Kind:    a.Kind,
					}),
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &a.Spec.Size,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Image:           c.UiImage,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Name:            fmt.Sprintf("%s-%s", a.Name, "ui"),
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
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      fmt.Sprintf("%s-%s", a.Name, "ui"),
									MountPath: "/html/config",
								},
							},
						}},
						Volumes: []corev1.Volume{
							{
								Name: fmt.Sprintf("%s-%s", a.Name, "ui"),
								VolumeSource: corev1.VolumeSource{
									ConfigMap: &corev1.ConfigMapVolumeSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: fmt.Sprintf("%s-%s", a.Name, "ui"),
										},
									},
								},
							},
						},
					},
				},
			},
		}

		// Create the deployment
		if err = client.Create(context.TODO(), dep); err != nil {
			return
		}
	}

	return
}

// Creates and returns a generator Deployment object
func generatorDeployment(client client.Client, c *configuration.Config, a *v1alpha1.Apicurito) (dep *appsv1.Deployment, err error) {
	dep = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", a.Name, "ui"),
		},
	}
	err = client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-%s", a.Name, "generator"), Namespace: a.Namespace}, dep)

	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := &appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", a.Name, "generator"),
				Namespace: a.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(a, schema.GroupVersionKind{
						Group:   v1alpha1.SchemeGroupVersion.Group,
						Version: v1alpha1.SchemeGroupVersion.Version,
						Kind:    a.Kind,
					}),
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &a.Spec.Size,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Image:           c.GeneratorImage,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Name:            fmt.Sprintf("%s-%s", a.Name, "generator"),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									Name:          "http",
								},
								{
									ContainerPort: 9779,
									Name:          "prometheus",
								},
								{
									ContainerPort: 8778,
									Name:          "jolokia",
								},
							},
							LivenessProbe: &corev1.Probe{
								FailureThreshold:    3,
								InitialDelaySeconds: 180,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								TimeoutSeconds:      1,
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Scheme: corev1.URISchemeHTTP,
										Port:   intstr.FromInt(8181),
										Path:   "/health",
									}},
							},
							ReadinessProbe: &corev1.Probe{
								FailureThreshold:    3,
								InitialDelaySeconds: 10,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								TimeoutSeconds:      1,
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Scheme: corev1.URISchemeHTTP,
										Port:   intstr.FromInt(8181),
										Path:   "/health",
									}},
							},
						}},
					},
				},
			},
		}

		// Create the deployment
		if err = client.Create(context.TODO(), dep); err != nil {
			return dep, err
		}
	} else if err != nil {
		return
	}

	return
}
