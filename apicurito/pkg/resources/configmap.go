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

	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func apicuritoConfig(client client.Client, a *v1alpha1.Apicurito) (c *corev1.ConfigMap, err error) {
	c = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", a.Name, "ui"),
		},
	}

	err = client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-%s", a.Name, "ui"), Namespace: a.Namespace}, c)
	if err != nil && errors.IsNotFound(err) {
		c = &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
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
			Data: map[string]string{
				"config.js": "var ApicuritoConfig = { \"generators\": [ { \"name\":\"Fuse Camel Project\", \"url\":\"/api/v1/generate/camel-project.zip\" } ] }",
			},
		}

		// Create the deployment
		if err = client.Create(context.TODO(), c); err != nil {
			return
		}
	}
	return
}
