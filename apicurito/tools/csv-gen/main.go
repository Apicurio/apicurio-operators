package main

import (
	"fmt"
	"os"

	"github.com/apicurio/apicurio-operators/apicurito/tools/run"
)

func main() {

	if err := run.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
