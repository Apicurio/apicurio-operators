package run

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/blang/semver"

	api "github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	config "github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"
	"github.com/apicurio/apicurio-operators/apicurito/tools/components"
	"github.com/apicurio/apicurio-operators/apicurito/tools/constants"
	"github.com/apicurio/apicurio-operators/apicurito/tools/util"
	"github.com/apicurio/apicurio-operators/apicurito/version"
	"github.com/heroku/docker-registry-client/registry"
	oimagev1 "github.com/openshift/api/image/v1"
	csvv1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	olmversion "github.com/operator-framework/operator-lifecycle-manager/pkg/lib/version"
	"github.com/tidwall/sjson"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/json"

	"os"

	"strconv"
	"strings"
	"time"
)

var (
	rh         = "Red Hat"
	maturity   = "alpha"
	maintainer = "Apicurito Project"
	csv        = csvSetting{

		Name:         "fuse-apicurito",
		DisplayName:  "Red Hat Apicurito Operator",
		OperatorName: "fuse-apicurito",
		CsvDir:       "apicurito-operator",
		Registry:     "registry.redhat.io",
		Context:      "fuse7",
		ImageName:    "fuse-apicurito-rhel7-operator",
		Tag:          constants.Apicurito16ImageTag,
	}
)

func Run() error {
	c := &config.Config{}
	if err := c.Config(nil); err != nil {
		return err
	}

	imageShaMap := map[string]string{}

	operatorName := csv.Name

	templateStruct := &csvv1.ClusterServiceVersion{}
	templateStruct.SetGroupVersionKind(csvv1.SchemeGroupVersion.WithKind("ClusterServiceVersion"))
	csvStruct := &csvv1.ClusterServiceVersion{}
	strategySpec := &csvStrategySpec{}
	json.Unmarshal(csvStruct.Spec.InstallStrategy.StrategySpecRaw, strategySpec)

	templateStrategySpec := &csvStrategySpec{}

	deployment := components.GetDeployment(csv.OperatorName, csv.Registry, csv.Context, csv.ImageName, csv.Tag, "Always", c)
	templateStrategySpec.Deployments = append(templateStrategySpec.Deployments, []csvDeployments{{Name: csv.OperatorName, Spec: deployment.Spec}}...)
	role := components.GetRole(csv.OperatorName)
	templateStrategySpec.Permissions = append(templateStrategySpec.Permissions, []csvPermissions{{ServiceAccountName: deployment.Spec.Template.Spec.ServiceAccountName, Rules: role.Rules}}...)
	clusterRole := components.GetClusterRole(csv.OperatorName)
	templateStrategySpec.ClusterPermissions = append(templateStrategySpec.ClusterPermissions, []StrategyDeploymentPermissions{{ServiceAccountName: deployment.Spec.Template.Spec.ServiceAccountName, Rules: clusterRole.Rules}}...)

	updatedStrat, err := json.Marshal(templateStrategySpec)
	if err != nil {
		panic(err)
	}
	templateStruct.Spec.InstallStrategy.StrategySpecRaw = updatedStrat
	templateStruct.Spec.InstallStrategy.StrategyName = "deployment"
	csvVersionedName := operatorName + ".v" + version.Version
	templateStruct.Name = csvVersionedName
	templateStruct.Namespace = "placeholder"
	annotdescrip := "Manages the installation and upgrades of apicurito, a small/minimal version of Apicurio"
	var description = "Apicurito is a small/minimal version of Apicurio, a standalone API design studio that can be used to create new or edit existing API designs (using the OpenAPI specification).\n"
	description += "\n"
	description += "This operator supports the installation and upgrade of apicurito. Apicurito components are:\n"
	description += "   - apicurito-ui (Apicurito API designer application)\n"
	description += "   - apicurito-generator (Apicurito fuse project generator)\n"
	description += "   - apicurito route (to access apicurito from outside openshift)\n"
	description += "\n"
	description += "### How to install\n"
	description += "When the operator is installed (you have created a subscription and the operator is running in the selected namespace) create a new CR of Kind Apicurito (click the Create New button). The CR spec contains all defaults.\n"
	description += "\n"
	description += "At the moment, following fields are supported as part of the CR:\n"
	description += "   - size: how many pods your the apicurito operand will have.\n"
	description += "\n"

	repository := "https://github.com/Apicurio/apicurio-operators/tree/master/apicurito"
	examples := []string{"{\n        \"apiVersion\": \"apicur.io/v1alpha1\",\n        \"kind\": \"Apicurito\",\n        \"metadata\": {\n          \"name\": \"apicurito-service\"\n        },\n        \"spec\": {\n          \"size\": 3\n        }\n}"}

	templateStruct.SetAnnotations(
		map[string]string{
			"capabilities":   "Seamless Upgrades",
			"categories":     "Integration & Delivery",
			"certified":      "false",
			"createdAt":      time.Now().Format("2006-01-02 15:04:05"),
			"containerImage": deployment.Spec.Template.Spec.Containers[0].Image,
			"support":        "Apicurito Project",
			"description":    annotdescrip,
			"repository":     repository,
			"alm-examples":   "[" + strings.Join(examples, ",") + "]",
		},
	)
	/*templateStruct.SetLabels(
		map[string]string{
			"operator-" + csv.Name: "true",
		},
	)*/

	var opVersion olmversion.OperatorVersion
	templateStruct.Spec.DisplayName = csv.DisplayName
	templateStruct.Spec.Description = description
	templateStruct.Spec.Keywords = []string{"api", "apicurio", "apicurito"}
	opVersion.Version = semver.MustParse(version.Version)
	templateStruct.Spec.Version = opVersion
	templateStruct.Spec.Maturity = maturity
	templateStruct.Spec.Maintainers = []csvv1.Maintainer{{Name: maintainer, Email: "apicurio@lists.jboss.org"}}
	templateStruct.Spec.Provider = csvv1.AppLink{Name: rh}
	templateStruct.Spec.Links = []csvv1.AppLink{
		{Name: "Apicurito source code", URL: "https://github.com/Apicurio/apicurito"},
		{Name: "Apicurito operator source code", URL: "https://github.com/Apicurio/apicurio-operators/tree/master/apicurito"},
	}

	templateStruct.Spec.Icon = []csvv1.Icon{
		{
			Data:      "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",
			MediaType: "image/png",
		},
	}
	tLabels := map[string]string{
		//"alm-owner-" + csv.Name: operatorName,
		"name": operatorName,
	}
	templateStruct.Spec.Labels = tLabels
	templateStruct.Spec.Selector = &metav1.LabelSelector{MatchLabels: tLabels}
	templateStruct.Spec.InstallModes = []csvv1.InstallMode{
		{Type: csvv1.InstallModeTypeOwnNamespace, Supported: true},
		{Type: csvv1.InstallModeTypeSingleNamespace, Supported: true},
		{Type: csvv1.InstallModeTypeMultiNamespace, Supported: false},
		{Type: csvv1.InstallModeTypeAllNamespaces, Supported: false},
	}
	templateStruct.Spec.Replaces = operatorName + ".v" + version.PriorVersion
	templateStruct.Spec.CustomResourceDefinitions.Owned = []csvv1.CRDDescription{
		{
			Version:     api.SchemeGroupVersion.Version,
			Kind:        "Apicurito",
			DisplayName: "Apicurito CRD",
			Description: "CRD for Apicurito",
			Name:        "apicuritoes." + api.SchemeGroupVersion.Group,
			SpecDescriptors: []csvv1.SpecDescriptor{

				{
					Description:  "The number of Apicurito pods to deploy",
					DisplayName:  "Size",
					Path:         "size",
					XDescriptors: []string{"urn:alm:descriptor:com.tectonic.ui:fieldGroup:Deployment", "urn:alm:descriptor:com.tectonic.ui:podCount"},
				},
			},
		},
	}

	// opMajor, opMinor, opMicro := config.MajorMinorMicro(version.Version)
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path = path + "/deploy/"

	csvFile := path + "bundle/manifests/apicurito.clusterserviceversion.yaml"
	if err := ensureDir(path); err != nil {
		return err
	}

	imageName, _, _ := config.GetImage(deployment.Spec.Template.Spec.Containers[0].Image)
	relatedImages := []image{}

	templateStruct.Annotations["certified"] = "false"
	deployFile := path + "operator.yaml"
	createFile(deployFile, deployment)
	roleFile := path + "role.yaml"
	createFile(roleFile, role)

	relatedImages = append(relatedImages, image{Name: imageName, Image: deployment.Spec.Template.Spec.Containers[0].Image})

	imageRef := constants.ImageRef{
		TypeMeta: metav1.TypeMeta{
			APIVersion: oimagev1.SchemeGroupVersion.String(),
			Kind:       "ImageStream",
		},
		Spec: constants.ImageRefSpec{
			Tags: []constants.ImageRefTag{
				{
					// Needs to match the component name for upstream and downstream.
					Name: "fuse7/fuse-apicurito-rhel7-operator",
					From: &corev1.ObjectReference{
						// Needs to match the image that is in your CSV that you want to replace.
						Name: deployment.Spec.Template.Spec.Containers[0].Image,
						Kind: "DockerImage",
					},
				},
			},
		},
	}

	imageRef.Spec.Tags = append(imageRef.Spec.Tags, constants.ImageRefTag{
		Name: constants.Apicurito16Component,
		From: &corev1.ObjectReference{
			Name: constants.Apicurito16ImageURL,
			Kind: "DockerImage",
		},
	})

	relatedImages = append(relatedImages, getRelatedImage(constants.Apicurito16ImageURL))
	relatedImages = append(relatedImages, getRelatedImage(constants.Generator16ImageURL))

	if GetBoolEnv("DIGESTS") {

		for _, tagRef := range imageRef.Spec.Tags {

			if _, ok := imageShaMap[tagRef.From.Name]; !ok {
				imageShaMap[tagRef.From.Name] = ""
				imageName, imageTag, imageContext := config.GetImage(tagRef.From.Name)
				repo := imageContext + "/" + imageName

				digests, err := RetriveFromRedHatIO(repo, imageTag)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				if len(digests) > 1 {
					imageShaMap[tagRef.From.Name] = strings.ReplaceAll(tagRef.From.Name, ":"+imageTag, "@"+digests[len(digests)-1])
				}
			}
		}
	}

	//not sure if we required mage-references file in the future So comment out for now.

	//imageFile := "deploy/olm-catalog/" + csv.CsvDir + "/" + opMajor + "." + opMinor + "." + opMicro + "/" + "image-references"
	//createFile(imageFile, imageRef)

	var templateInterface interface{}
	if len(relatedImages) > 0 {
		templateJSON, err := json.Marshal(templateStruct)
		if err != nil {
			fmt.Println(err)
		}
		result, err := sjson.SetBytes(templateJSON, "spec.relatedImages", relatedImages)
		if err != nil {
			fmt.Println(err)

		}
		if err = json.Unmarshal(result, &templateInterface); err != nil {
			fmt.Println(err)
		}
	} else {
		templateInterface = templateStruct
	}

	// find and replace images with SHAs where necessary
	templateByte, err := json.Marshal(templateInterface)
	if err != nil {
		fmt.Println(err)
	}
	for from, to := range imageShaMap {
		if to != "" {
			templateByte = bytes.ReplaceAll(templateByte, []byte(from), []byte(to))
		}
	}
	if err = json.Unmarshal(templateByte, &templateInterface); err != nil {
		fmt.Println(err)
	}
	createFile(csvFile, &templateInterface)

	o, err := buildContainer()
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filepath.Join(path, "bundle", "container.yaml"), o, 0644); err != nil {
		return err
	}

	od, err := buildDocker(c)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filepath.Join(path, "bundle", "Dockerfile"), od, 0644); err != nil {
		return err
	}

	oa, err := buildAnnotation()
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filepath.Join(path, "bundle", "metadata", "annotations.yaml"), oa, 0644); err != nil {
		return err
	}

	_, err = copy(filepath.Join(path, "crds", "apicur.io_apicuritoes_crd.yaml"), filepath.Join(path, "bundle", "manifests", "apicur.io_apicuritoes_crd.yaml"))
	if err != nil {
		return err
	}

	return nil
}

func RetriveFromRedHatIO(image string, imageTag string) ([]string, error) {

	url := "https://" + constants.RedHatImageRegistry

	username := "" // anonymous
	password := "" // anonymous

	if userToken := strings.Split(os.Getenv("REDHATIO_TOKEN"), ":"); len(userToken) > 1 {
		username = userToken[0]
		password = userToken[1]
	}
	hub, err := registry.New(url, username, password)
	digests := []string{}
	if err != nil {
		fmt.Println(err)
	} else {
		tags, err := hub.Tags(image)
		if err != nil {
			fmt.Println(err)
		}
		// do not follow redirects - this is critical so we can get the registry digest from Location in redirect response
		hub.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		if _, exists := find(tags, imageTag); exists {
			req, err := http.NewRequest("GET", url+"/v2/"+image+"/manifests/"+imageTag, nil)
			if err != nil {
				fmt.Println(err)
			}
			req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
			resp, err := hub.Client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			if resp != nil {
				defer resp.Body.Close()
			}
			if resp.StatusCode == 302 || resp.StatusCode == 301 {
				digestURL, err := resp.Location()
				if err != nil {
					fmt.Println(err)

				}

				if digestURL != nil {
					if url := strings.Split(digestURL.EscapedPath(), "/"); len(url) > 1 {
						digests = url

						return url, err
					}
				}
			}
		}

	}
	return digests, err
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

type csvSetting struct {
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	OperatorName string `json:"operatorName"`
	CsvDir       string `json:"csvDir"`
	Registry     string `json:"repository"`
	Context      string `json:"context"`
	ImageName    string `json:"imageName"`
	Tag          string `json:"tag"`
}
type csvPermissions struct {
	ServiceAccountName string              `json:"serviceAccountName"`
	Rules              []rbacv1.PolicyRule `json:"rules"`
}
type csvDeployments struct {
	Name string                `json:"name"`
	Spec appsv1.DeploymentSpec `json:"spec,omitempty"`
}
type csvStrategySpec struct {
	Permissions        []csvPermissions                `json:"permissions"`
	Deployments        []csvDeployments                `json:"deployments"`
	ClusterPermissions []StrategyDeploymentPermissions `json:"clusterPermissions,omitempty"`
}

type StrategyDeploymentPermissions struct {
	ServiceAccountName string            `json:"serviceAccountName"`
	Rules              []rbac.PolicyRule `json:"rules"`
}
type channel struct {
	Name       string `json:"name"`
	CurrentCSV string `json:"currentCSV"`
}

type image struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func getRelatedImage(imageURL string) image {
	imageName, _, _ := config.GetImage(imageURL)
	return image{
		Name:  imageName,
		Image: imageURL,
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createFile(filepath string, obj interface{}) {
	f, err := os.Create(filepath)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	writer := bufio.NewWriter(f)
	util.MarshallObject(obj, writer)
	writer.Flush()
}

func GetBoolEnv(key string) bool {
	val := GetEnv(key, "false")
	ret, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return ret
}

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// Generate directory tree for OLM bundle. Path must already
// exist
func ensureDir(path string) (err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	err = os.Mkdir(filepath.Join(path, "bundle"), 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	err = os.Mkdir(filepath.Join(path, "bundle", "manifests"), 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	err = os.Mkdir(filepath.Join(path, "bundle", "metadata"), 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

func buildContainer() (out []byte, err error) {
	m := `operator_manifests:
  manifests_dir: manifests
  enable_digest_pinning: true
  enable_repo_replacements: true
  enable_registry_replacements: true
  repo_replacements:
    - registry: registry.redhat.io
      package_mappings:
        fuse-apicurito-operator-container: fuse7
        fuse-apicurito-openshift-container: fuse7
        fuse-apicurito-generator-openshift-container: fuse7`

	out = []byte(m)
	return
}

func buildDocker(c *config.Config) (out []byte, err error) {
	name := "fuse-apicurito"
	channel := fmt.Sprintf("%s-%s.x", name, version.ShortVersion())

	m := `FROM scratch

LABEL operators.operatorframework.io.bundle.mediatype.v1=registry+v1
LABEL operators.operatorframework.io.bundle.manifests.v1=manifests/
LABEL operators.operatorframework.io.bundle.metadata.v1=metadata/
LABEL operators.operatorframework.io.bundle.package.v1=-%s
LABEL operators.operatorframework.io.bundle.channels.v1=%s
LABEL operators.operatorframework.io.bundle.channel.default.v1=%s
LABEL com.redhat.delivery.operator.bundle=true
LABEL com.redhat.openshift.versions="%s"

COPY manifests /manifests/
COPY metadata/annotations.yaml /metadata/annotations.yaml

LABEL name="fuse7/fuse-online-operator-metadata" \
      version="%s" \
      maintainer="Otavio Piske <opiske@redhat.com>" \
      summary="Operator which manages the lifecycle of the Apicurito application." \
      description="Operator which manages the lifecycle of the Apicurito application." \
      com.redhat.component="fuse-apicurito-operator-metadata-container" \
      io.k8s.description="Operator which manages the lifecycle of the Apicurito application." \
      io.k8s.display-name="Red Hat Apicurito Operator" \
      io.openshift.tags="fuse,API"
`
	m = fmt.Sprintf(m, name, channel, channel, c.SupportedOpenShiftVersions, fmt.Sprintf("%s.x", version.ShortVersion()))
	out = []byte(m)
	return
}

func buildAnnotation() (out []byte, err error) {
	name := "fuse-apicurito"

	channel := fmt.Sprintf("apicurito-%s.x", version.ShortVersion())

	m := map[string]map[string]string{
		"annotations": {
			"operators.operatorframework.io.bundle.mediatype.v1":       "registry+v1",
			"operators.operatorframework.io.bundle.manifests.v1":       "manifests/",
			"operators.operatorframework.io.bundle.metadata.v1":        "metadata/",
			"operators.operatorframework.io.bundle.package.v1":         name,
			"operators.operatorframework.io.bundle.channels.v1":        channel,
			"operators.operatorframework.io.bundle.channel.default.v1": channel,
		},
	}
	out, err = yaml.Marshal(m)
	return
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
