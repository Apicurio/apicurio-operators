package apicurito

//go:generate go run ./.packr/packr.go

import (
	"context"
	"fmt"
	"github.com/RHsyseng/operator-utils/pkg/logs"
	"github.com/RHsyseng/operator-utils/pkg/utils/kubernetes"
	"github.com/RHsyseng/operator-utils/pkg/utils/openshift"
	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	"github.com/apicurio/apicurio-operators/apicurito/pkg/resources"
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr/v2"
	consolev1 "github.com/openshift/api/console/v1"
	routev1 "github.com/openshift/api/route/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

var logu = logs.GetLogger("openshift-webconsole")

func consoleLinkExists() error {
	gvk := schema.GroupVersionKind{Group: "console.openshift.io", Version: "v1", Kind: "ConsoleLink"}
	return kubernetes.CustomResourceDefinitionExists(gvk)
}

func removeConsoleLink(c client.Client, api *v1alpha1.Apicurito) {
	doDeleteConsoleLink(getUIConsoleLinkName(api), c, api)
	doDeleteConsoleLink(getGeneratorConsoleLinkName(api), c, api)
}

func doDeleteConsoleLink(consoleLinkName string, c client.Client, api *v1alpha1.Apicurito) {
	consoleLink := &consolev1.ConsoleLink{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: consoleLinkName}, consoleLink)
	if err == nil && consoleLink != nil {
		err = c.Delete(context.TODO(), consoleLink)
		if err != nil {
			logu.Error(err, "Failed to delete the consolelink:", consoleLinkName)
		} else {
			logu.Info("deleted the consolelink:", consoleLinkName)
		}
	}
}

func createConsoleLink(c client.Client, api *v1alpha1.Apicurito) {
	doCreateConsoleLink(getUIConsoleLinkName(api), resources.GetUIRouteName(api), c, api)
	doCreateConsoleLink(getGeneratorConsoleLinkName(api), resources.GetGeneratorRouteName(api), c, api)
}

func doCreateConsoleLink(consoleLinkName string, routeName string, c client.Client, api *v1alpha1.Apicurito) {
	route := &routev1.Route{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: routeName, Namespace: api.Namespace}, route)
	if err == nil && route != nil {
		checkConsoleLink(route, consoleLinkName, api, c)
	}
}

func checkConsoleLink(route *routev1.Route, consoleLinkName string, api *v1alpha1.Apicurito, c client.Client) {
	consoleLink := &consolev1.ConsoleLink{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: consoleLinkName}, consoleLink)
	if err != nil && apierrors.IsNotFound(err) {
		consoleLink = createNamespaceDashboardLink(consoleLinkName, route, api)
		if err := c.Create(context.TODO(), consoleLink); err != nil {
			logu.Error(err, "Console link is not created.")
		} else {
			logu.Info("Console link has been created. ", consoleLinkName)
		}
	} else if err == nil && consoleLink != nil {
		reconcileConsoleLink(context.TODO(), route, consoleLink, c)
	}
}

func reconcileConsoleLink(ctx context.Context, route *routev1.Route, link *consolev1.ConsoleLink, client client.Client) {
	url := "https://" + route.Spec.Host
	linkTxt := ConsoleLinkText(route)
	if url != link.Spec.Href || linkTxt != link.Spec.Text {
		if err := client.Update(ctx, link); err != nil {
			logu.Error(err, "failed to reconcile Console Link", link)
		}
	}
}

func getUIConsoleLinkName(api *v1alpha1.Apicurito) string {
	return fmt.Sprintf("%s-%s", resources.GetUIRouteName(api), api.Namespace)
}

func getGeneratorConsoleLinkName(api *v1alpha1.Apicurito) string {
	return fmt.Sprintf("%s-%s", resources.GetGeneratorRouteName(api), api.Namespace)
}

func ConsoleLinkText(route *routev1.Route) string {
	name := route.Name
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "")
	name = strings.TrimPrefix(name, "apicurito")
	name = strings.TrimSuffix(name, "apicurito")
	name = strings.Title(name)
	return "Apicurito - " + name
}

func createNamespaceDashboardLink(consoleLinkname string, route *routev1.Route, api *v1alpha1.Apicurito) *consolev1.ConsoleLink {
	return &consolev1.ConsoleLink{
		ObjectMeta: metav1.ObjectMeta{
			Name: consoleLinkname,
			Labels: map[string]string{
				"apicurito.io/name": api.ObjectMeta.Name,
			},
		},
		Spec: consolev1.ConsoleLinkSpec{
			Link: consolev1.Link{
				Text: ConsoleLinkText(route),
				Href: "https://" + route.Spec.Host,
			},
			Location: consolev1.NamespaceDashboard,
			NamespaceDashboard: &consolev1.NamespaceDashboardSpec{
				Namespaces: []string{api.Namespace},
			},
		},
	}
}

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
