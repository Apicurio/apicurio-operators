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

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/RHsyseng/operator-utils/pkg/resource"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var labels = map[string]string{"app": "apicurito"}

func apicuritoService(a *v1alpha1.Apicurito) (s resource.KubernetesResource) {

	// Define new service
	name := fmt.Sprintf("%s-%s", a.Name, "ui")
	s = &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: a.Namespace,
			Labels:    labelComponent(name),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(a, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    a.Kind,
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labelComponent(name),
			Ports: []corev1.ServicePort{
				{
					Name:       "api-port",
					Port:       8080,
					TargetPort: intstr.FromString("api-port"),
				},
			},
		},
	}

	return
}

func generatorService(a *v1alpha1.Apicurito) (s resource.KubernetesResource) {
	s = &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", a.Name, "generator"),
			Namespace: a.Namespace,
			Labels:    labelComponent("apicurito-service-generator"),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(a, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    a.Kind,
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labelComponent("apicurito-service-generator"),
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       8080,
					TargetPort: intstr.FromString("http"),
				},
			},
		},
	}

	return
}

func labelComponent(name string) map[string]string {
	_labels := make(map[string]string)
	for index, element := range labels {
		_labels[index] = element
	}
	_labels["component"] = name
	return _labels
}
