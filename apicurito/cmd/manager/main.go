package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/cmd"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	sdkVersion "github.com/operator-framework/operator-sdk/version"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// Change below variables to serve metrics on different host or port.
var (
	metricsHost       = "0.0.0.0"
	metricsPort int32 = 8383
)
var log = logf.Log.WithName("cmd")

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info(fmt.Sprintf("Version of operator-sdk: %v", sdkVersion.Version))
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel ctx as soon as main returns

	apicurito, err := cmd.NewApicuritoCommand(ctx)
	exeName := filepath.Base(os.Args[0])
	if !strings.Contains(exeName, "go_build_main_go") {
		apicurito.Use = exeName
	}
	exitOnError(err)

	err = apicurito.Execute()
	exitOnError(err)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
