package apicurito

//go:generate go run ./.packr/packr.go

import (
	"context"
	"github.com/RHsyseng/operator-utils/pkg/logs"
	"github.com/RHsyseng/operator-utils/pkg/utils/kubernetes"
	"github.com/RHsyseng/operator-utils/pkg/utils/openshift"
	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var logu = logs.GetLogger("openshift-webconsole")

func ConsoleYAMLSampleExists() error {
	gvk := schema.GroupVersionKind{Group: "console.openshift.io", Version: "v1", Kind: "ConsoleYAMLSample"}
	return kubernetes.CustomResourceDefinitionExists(gvk)
}

func createConsoleYAMLSamples(c client.Client) {
	logu.Info("Loading CR YAML samples.")
	box := packr.New("cryamlsamples", "../../../deploy/crs")
	if box.List() == nil {
		logu.Error(nil, "CR YAML folder is empty. It is not loaded.")
		return
	}
	for _, filename := range box.List() {
		yamlStr, err := box.FindString(filename)
		if err != nil {
			logu.Info("yaml", " name: ", filename, " not created:  ", err.Error())
			continue
		}
		apicurito := v1alpha1.Apicurito{}
		err = yaml.Unmarshal([]byte(yamlStr), &apicurito)
		if err != nil {
			logu.Info("yaml", " name: ", filename, " not created:  ", err.Error())
			continue
		}
		yamlSample, err := openshift.GetConsoleYAMLSample(&apicurito)
		if err != nil {
			logu.Info("yaml", " name: ", filename, " not created:  ", err.Error())
			continue
		}
		err = c.Create(context.TODO(), yamlSample)
		if err != nil {
			if !apierrors.IsAlreadyExists(err) {
				logu.Info("yaml", " name: ", filename, " not created:+", err.Error())
			}
			continue
		}
		logu.Info("yaml", " name: ", filename, " Created.")
	}
}
