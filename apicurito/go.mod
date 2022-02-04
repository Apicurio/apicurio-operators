module github.com/apicurio/apicurio-operators/apicurito

go 1.13

require (
	github.com/RHsyseng/operator-utils v0.0.0-00010101000000-000000000000
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/coreos/prometheus-operator v0.38.1-0.20200424145508-7e176fda06cc // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-logr/logr v0.4.0
	github.com/go-openapi/spec v0.19.9 // indirect
	github.com/gobuffalo/packr/v2 v2.7.1
	github.com/heroku/docker-registry-client v0.0.0-20190909225348-afc9e1acc3d5 // indirect
	github.com/imdario/mergo v0.3.12
	github.com/openshift/api v3.9.1-0.20190924102528-32369d4db2ad+incompatible
	github.com/operator-framework/operator-lifecycle-manager v0.0.0-20191115003340-16619cd27fa5 // indirect
	github.com/operator-framework/operator-sdk v0.19.4
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/kube-openapi v0.0.0-20210421082810-95288971da7e
	sigs.k8s.io/controller-runtime v0.10.0

)

replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.10.0

replace (
	// OpenShift release-4.11
	github.com/openshift/api => github.com/openshift/api v0.0.0-20200930075302-db52bc4ef99f
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200929181438-91d71ef2122c
	k8s.io/api => k8s.io/api v0.22.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.1
	k8s.io/client-go => k8s.io/client-go v0.22.1
)

replace (
	github.com/RHsyseng/operator-utils => github.com/RHsyseng/operator-utils v1.4.7
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20191016112112-5190913f932d
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191016114015-74ad18325ed5
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20191016115326-20453efc2458
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20191016115129-c07a134afb42
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191004115455-8e001e5d1894
	k8s.io/component-base => k8s.io/component-base v0.0.0-20191016111319-039242c015a9
	k8s.io/cri-api => k8s.io/cri-api v0.0.0-20190828162817-608eb1dad4ac
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20191016115521-756ffa5af0bd
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191016112429-9587704a8ad4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20191016114939-2b2b218dc1df
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20191016114407-2e83b6f20229
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20191016114748-65049c67a58b
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191016120415-2ed914427d51
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20191016114556-7841ed97f1b2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20191016115753-cf0698c3a16b
	k8s.io/metrics => k8s.io/metrics v0.0.0-20191016113814-3b1a734dba6e
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20191016112829-06bb3c9d77c9
)

replace github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.19.4

replace golang.org/x/text => golang.org/x/text v0.3.3
