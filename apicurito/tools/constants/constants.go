package constants

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	RedHatImageRegistry = "registry.redhat.io"
	Apicurito16Image    = "fuse-apicurito"

	Apicurito16ImageTag = "1.8"
	Apicurito16ImageURL = RedHatImageRegistry + "/fuse7/" + Apicurito16Image + ":" + Apicurito16ImageTag

	Generator16Image    = "fuse-apicurito-generator"
	Generator16ImageURL = RedHatImageRegistry + "/fuse7/" + Generator16Image + ":" + Apicurito16ImageTag

	Apicurito16Component = "fuse-apicurito-openshift-container"
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
