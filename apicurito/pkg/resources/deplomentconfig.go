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

	"github.com/RHsyseng/operator-utils/pkg/resource"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	"github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"
	"github.com/apicurio/apicurio-operators/apicurito/version"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Creates and returns a apicurito Deployment object
func apicuritoDeployment(c *configuration.Config, a *v1alpha1.Apicurito) (dep resource.KubernetesResource) {
	// Define a new deployment
	var dm int32 = 420
	deployLabels := map[string]string{
		"app":           "apicurito",
		"component":     fmt.Sprintf("%s-%s", a.Name, "ui"),
		"com.company":   "Red_Hat",
		"rht.prod_name": "Red_Hat_Integration",
		"rht.prod_ver":  version.ShortVersion(),
		"rht.comp":      "Fuse",
		"rht.comp_ver":  version.ShortVersion(),
		"rht.subcomp":   fmt.Sprintf("%s-%s", a.Name, "ui"),
		"rht.subcomp_t": "infrastructure",
	}
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
						Name:            fmt.Sprintf("%s-%s", a.Name, "ui"),
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
									DefaultMode: &dm,
								},
							},
						},
					},
				},
			},
		},
	}

	return
}

// Creates and returns a generator Deployment object
func generatorDeployment(c *configuration.Config, a *v1alpha1.Apicurito) (dep resource.KubernetesResource) {
	// Define a new deployment
	deployLabels := map[string]string{
		"app":           "apicurito",
		"component":     fmt.Sprintf("%s-%s", a.Name, "ui"),
		"com.company":   "Red_Hat",
		"rht.prod_name": "Red_Hat_Integration",
		"rht.prod_ver":  version.ShortVersion(),
		"rht.comp":      "Fuse",
		"rht.comp_ver":  version.ShortVersion(),
		"rht.subcomp":   fmt.Sprintf("%s-%s", a.Name, "generator"),
		"rht.subcomp_t": "infrastructure",
	}
	dep = &appsv1.Deployment{
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
						Name:            fmt.Sprintf("%s-%s", a.Name, "generator"),
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
								TCPSocket: &corev1.TCPSocketAction{
									Port: intstr.FromString("http"),
								},
							},
						},
						ReadinessProbe: &corev1.Probe{
							FailureThreshold:    3,
							InitialDelaySeconds: 10,
							PeriodSeconds:       10,
							SuccessThreshold:    1,
							TimeoutSeconds:      1,
							Handler: corev1.Handler{
								TCPSocket: &corev1.TCPSocketAction{
									Port: intstr.FromString("http"),
								},
							},
						},
					}},
				},
			},
		},
	}

	return
}
