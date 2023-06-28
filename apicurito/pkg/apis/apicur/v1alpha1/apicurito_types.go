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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApicuritoSpec defines the desired state of Apicurito
// +k8s:openapi-gen=true
type ApicuritoSpec struct {
	// The number of pods to scale
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Size"
	Size int32 `json:"size"`

	// Resources for the UI pod
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="ResourcesUI"
	ResourcesUI *corev1.ResourceRequirements `json:"resourcesUI,omitempty"`

	// Resources for the generator pod
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="ResourcesGenerator"
	ResourcesGenerator *corev1.ResourceRequirements `json:"resourcesGenerator,omitempty"`
}

// ApicuritoPhase defines the phase of the installation
type ApicuritoPhase string

const (
	// ApicuritoPhaseMissing defines if Apicurito CR not available
	ApicuritoPhaseMissing ApicuritoPhase = ""
	// ApicuritoPhaseStarting defines if Apicurito install starting
	ApicuritoPhaseStarting ApicuritoPhase = "Starting"
	// ApicuritoPhaseInstalling defines if Apicurito install underway
	ApicuritoPhaseInstalling ApicuritoPhase = "Installing"
	// ApicuritoPhaseInstalled defines if Apicurito install completed
	ApicuritoPhaseInstalled ApicuritoPhase = "Installed"
	// ApicuritoPhaseInstallError defines if Apicurito install produced an error
	ApicuritoPhaseInstallError ApicuritoPhase = "InstallError"
)

// ApicuritoStatus defines the observed state of Apicurito
// +k8s:openapi-gen=true
type ApicuritoStatus struct {
	// The phase the operator has reached, eg. INSTALLED, STARTING
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +operator-sdk:csv:customresourcedefinitions:type=status,displayName="Phase"
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:io.kubernetes.phase"
	Phase ApicuritoPhase `json:"phase,omitempty"`

	// The currently installed version of Apicurito
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +operator-sdk:csv:customresourcedefinitions:type=status,displayName="Version"
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:text"
	Version string `json:"version,omitempty"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Apicurito is the Schema for the apicuritoes API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=apicuritoes,scope=Namespaced
// +kubebuilder:printcolumn:name="Phase",description="The apicurito phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Version",description="The apicurito version",type=string,JSONPath=`.status.version`
// +operator-sdk:csv:customresourcedefinitions:displayName="Apicurito"
// +operator-sdk:csv:customresourcedefinitions:resources={{ServiceAccount,v1},{ClusterRole,rbac.authorization.k8s.io/v1},{Role,rbac.authorization.k8s.io/v1},{Deployment,apps/v1}}
type Apicurito struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApicuritoSpec   `json:"spec,omitempty"`
	Status ApicuritoStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApicuritoList contains a list of Apicurito
type ApicuritoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Apicurito `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Apicurito{}, &ApicuritoList{})
}
