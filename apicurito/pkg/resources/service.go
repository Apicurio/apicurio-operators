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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var labels = map[string]string{"app": "apicurito"}

func apicuritoService(client client.Client, a *v1alpha1.Apicurito) (s *corev1.Service, err error) {
	s = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", a.Name, "ui"),
		},
	}
	err = client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-%s", a.Name, "ui"), Namespace: a.Namespace}, s)
	if err != nil && errors.IsNotFound(err) {
		// Define new service
		labels["component"] = fmt.Sprintf("%s-%s", a.Name, "ui")
		s = &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Service",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", a.Name, "ui"),
				Namespace: a.Namespace,
				Labels:    labels,
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
				Selector: labels,
				Ports: []corev1.ServicePort{
					{
						Name: "api-port",
						Port: 8080,
					},
				},
			},
		}

		err = client.Create(context.TODO(), s)
		if err != nil {
			return
		}
	}

	return
}

func generatorService(client client.Client, a *v1alpha1.Apicurito) (s *corev1.Service, err error) {
	s = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", a.Name, "generator"),
		},
	}
	err = client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-%s", a.Name, "generator"), Namespace: a.Namespace}, s)
	if err != nil && errors.IsNotFound(err) {
		labels["component"] = "apicurito-generator"

		s = &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Service",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", a.Name, "generator"),
				Namespace: a.Namespace,
				Labels:    labels,
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
				Selector: labels,
				Ports: []corev1.ServicePort{
					{
						Name: "generator-port",
						Port: 8080,
					},
				},
			},
		}

		err = client.Create(context.TODO(), s)
		if err != nil {
			return
		}
	}

	return
}
