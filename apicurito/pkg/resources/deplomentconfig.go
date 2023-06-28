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
	"fmt"
	"gopkg.in/inf.v0"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"math/big"

	api "github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	"github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Creates and returns an Apicurito Deployment object
func apicuritoDeployment(c *configuration.Config, a *api.Apicurito) client.Object {
	// Define a new deployment
	var dm int32 = 420
	name := fmt.Sprintf("%s-%s", a.Name, "ui")
	deployLabels := labelComponent(name)
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: a.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(a, schema.GroupVersionKind{
					Group:   api.SchemeGroupVersion.Group,
					Version: api.SchemeGroupVersion.Version,
					Kind:    a.Kind,
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &a.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: deployLabels,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: deployLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           c.UiImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            name,
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "api-port",
							Protocol:      corev1.ProtocolTCP,
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
								Name:      name,
								MountPath: "/html/config",
							},
						},
						Resources: corev1.ResourceRequirements{
							Requests: Normalize(corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("500m"),
								corev1.ResourceMemory: resource.MustParse("64Mi"),
							}),
							Limits: Normalize(corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("1000m"),
								corev1.ResourceMemory: resource.MustParse("128Mi"),
							}),
						},
					}},
					Volumes: []corev1.Volume{
						{
							Name: name,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: name,
									},
									DefaultMode: &dm,
								},
							},
						},
					},
				},
			},
		},
	}

	if a.Spec.ResourcesUI != nil {
		if a.Spec.ResourcesUI.Requests != nil {
			deployment.Spec.Template.Spec.Containers[0].Resources.Requests = Normalize(a.Spec.ResourcesUI.Requests)
		}
		if a.Spec.ResourcesUI.Limits != nil {
			deployment.Spec.Template.Spec.Containers[0].Resources.Limits = Normalize(a.Spec.ResourcesUI.Limits)
		}
	}

	return deployment
}

// Creates and returns a generator Deployment object
func generatorDeployment(c *configuration.Config, a *api.Apicurito) client.Object {
	// Define a new deployment
	name := fmt.Sprintf("%s-%s", a.Name, "generator")
	deployLabels := labelComponent(name)
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: a.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(a, schema.GroupVersionKind{
					Group:   api.SchemeGroupVersion.Group,
					Version: api.SchemeGroupVersion.Version,
					Kind:    a.Kind,
				}),
			},
		},

		Spec: appsv1.DeploymentSpec{
			Replicas: &a.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: deployLabels,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: deployLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           c.GeneratorImage,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            name,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 8080,
								Name:          "http",
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 8181,
								Name:          "health",
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 9779,
								Name:          "prometheus",
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 8778,
								Name:          "jolokia",
								Protocol:      corev1.ProtocolTCP,
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
									Port:   intstr.FromString("health"),
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
									Port:   intstr.FromString("health"),
									Path:   "/health",
								}},
						},
						Resources: corev1.ResourceRequirements{
							Requests: Normalize(corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("500m"),
								corev1.ResourceMemory: resource.MustParse("256Mi"),
							}),
							Limits: Normalize(corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("1000m"),
								corev1.ResourceMemory: resource.MustParse("512Mi"),
							}),
						},
					}},
				},
			},
		},
	}

	if a.Spec.ResourcesGenerator != nil {
		if a.Spec.ResourcesGenerator.Requests != nil {
			deployment.Spec.Template.Spec.Containers[0].Resources.Requests = Normalize(a.Spec.ResourcesGenerator.Requests)
		}
		if a.Spec.ResourcesGenerator.Limits != nil {
			deployment.Spec.Template.Spec.Containers[0].Resources.Limits = Normalize(a.Spec.ResourcesGenerator.Limits)
		}
	}

	return deployment
}

// Normalize transforms each resource.Quantity item into a canonical form,
// so reflect.DeepEquals works correctly.
func Normalize(in corev1.ResourceList) (out corev1.ResourceList) {
	out = make(map[corev1.ResourceName]resource.Quantity)
	for k, v := range in {
		// Normalize scale
		dec := *(&v).AsDec()

		unscaled := dec.UnscaledBig()
		scale := dec.Scale()

		zero := big.NewInt(0)
		ten := big.NewInt(10)

		for scale > 0 {
			rem := big.NewInt(0)
			rem.Mod(unscaled, ten)
			if rem.Cmp(zero) == 0 {
				// We can rescale
				unscaled.Div(unscaled, big.NewInt(10))
				scale = scale - 1
			} else {
				break
			}
		}

		out[k] = *resource.NewDecimalQuantity(*inf.NewDecBig(unscaled, scale), v.Format)
	}
	return
}
