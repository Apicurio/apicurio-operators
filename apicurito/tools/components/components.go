package components

import (
	"strings"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"

	monv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	consolev1 "github.com/openshift/api/console/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetDeployment(operatorName, repository, context, imageName, tag, imagePullPolicy string, cfg *configuration.Config) *appsv1.Deployment {
	registryName := strings.Join([]string{repository, context, imageName}, "/")
	image := strings.Join([]string{registryName, tag}, ":")
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appsv1.SchemeGroupVersion.String(),
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: operatorName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": operatorName,
					"app":  "apicurito",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": operatorName,
						"app":  "apicurito",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "apicurito",
					Containers: []corev1.Container{
						{
							Name:            operatorName,
							Image:           image,
							ImagePullPolicy: corev1.PullPolicy(imagePullPolicy),
							Env: []corev1.EnvVar{

								{
									Name: "WATCH_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name:  "RELATED_IMAGE_APICURITO_OPERATOR",
									Value: image,
								},
								{
									Name:  "RELATED_IMAGE_APICURITO",
									Value: cfg.UiImage,
								},
								{
									Name:  "RELATED_IMAGE_GENERATOR",
									Value: cfg.GeneratorImage,
								},
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},

								{
									Name:  "OPERATOR_NAME",
									Value: operatorName,
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func GetRole(operatorName string) *rbacv1.Role {
	role := &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: rbacv1.SchemeGroupVersion.String(),
			Kind:       "Role",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "apicurito",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"pods",
					"services",
					"endpoints",
					"persistentvolumeclaims",
					"events",
					"configmaps",
					"secrets",
					"serviceaccounts",
				},
				Verbs: []string{"*"},
			},
			{
				APIGroups: []string{
					"apps",
				},
				Resources: []string{
					"deployments",
					"daemonsets",
					"replicasets",
					"statefulsets",
				},
				Verbs: []string{"*"},
			},

			{
				APIGroups: []string{
					monv1.SchemeGroupVersion.Group,
				},
				Resources: []string{"servicemonitors"},
				Verbs:     []string{"get", "create"},
			},
			{
				APIGroups: []string{
					"apps",
				}, ResourceNames: []string{

					"apicurito",
				},
				Resources: []string{
					"deployments/finalizers",
				},
				Verbs: []string{"update"},
			},
			{
				APIGroups: []string{
					"apicur.io",
				},

				Resources: []string{
					"*",
				},
				Verbs: []string{"*"},
			},

			{
				APIGroups: []string{
					"route.openshift.io",
				},
				Resources: []string{
					"routes",
				},
				Verbs: []string{"get", "list", "create", "update", "watch"},
			},
		},
	}
	return role
}

func int32Ptr(i int32) *int32 {
	return &i
}

func GetClusterRole(operatorName string) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: rbacv1.SchemeGroupVersion.String(),
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: operatorName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{consolev1.GroupVersion.Group},
				Resources: []string{"consolelinks", "consoleyamlsamples"},
				Verbs: []string{
					"get",
					"create",
					"list",
					"update",
					"delete",
				},
			},
		},
	}
}
