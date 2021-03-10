package constants

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	RedHatImageRegistry = "quay.io"
	Apicurito16Image    = "apicurito"

	Apicurito16ImageTag = "latest"
	Apicurito16ImageURL = RedHatImageRegistry + "/apicurio/" + Apicurito16Image + ":" + Apicurito16ImageTag

	Generator16Image    = "apicurito-generator"
	Generator16ImageURL = RedHatImageRegistry + "/apicurio/" + Generator16Image + ":" + Apicurito16ImageTag

	Apicurito16Component = "apicurito-openshift-container"
)

type ImageEnv struct {
	Var       string
	Component string
	Registry  string
}
type ImageRef struct {
	metav1.TypeMeta `json:",inline"`
	Spec            ImageRefSpec `json:"spec"`
}
type ImageRefSpec struct {
	Tags []ImageRefTag `json:"tags"`
}
type ImageRefTag struct {
	Name string                  `json:"name"`
	From *corev1.ObjectReference `json:"from"`
}
