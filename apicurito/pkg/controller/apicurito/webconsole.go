package apicurito

//go:generate go run ./.packr/packr.go

import (
	"context"
	"github.com/RHsyseng/operator-utils/pkg/utils/kubernetes"
	"github.com/RHsyseng/operator-utils/pkg/utils/openshift"
	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ConsoleYAMLSampleExists() error {
	gvk := schema.GroupVersionKind{Group: "console.openshift.io", Version: "v1", Kind: "ConsoleYAMLSample"}
	return kubernetes.CustomResourceDefinitionExists(gvk)
}

func createConsoleYAMLSamples(c client.Client) {
	log.Info("Loading CR YAML samples.")
	box := packr.New("cryamlsamples", "../../../deploy/crs")
	if box.List() == nil {
		log.Error(nil, "CR YAML folder is empty. It is not loaded.")
		return
	}

	resMap := make(map[string]string)
	for _, filename := range box.List() {
		yamlStr, err := box.FindString(filename)
		if err != nil {
			resMap[filename] = err.Error()
			continue
		}
		apicurito := v1alpha1.Apicurito{}
		err = yaml.Unmarshal([]byte(yamlStr), &apicurito)
		if err != nil {
			resMap[filename] = err.Error()
			continue
		}
		yamlSample, err := openshift.GetConsoleYAMLSample(&apicurito)
		if err != nil {
			resMap[filename] = err.Error()
			continue
		}
		err = c.Create(context.TODO(), yamlSample)
		if err != nil {
			resMap[filename] = err.Error()
			continue
		}
		resMap[filename] = "Applied"
	}

	for k, v := range resMap {
		log.Info("yaml ", "name: ", k, " status: ", v)
	}
}
