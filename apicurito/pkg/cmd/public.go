/*
 * Copyright (C) 2019 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
	Context context.Context
	Command *cobra.Command
}

// Creates new Apicurito command
func NewApicuritoCommand(ctx context.Context) (*cobra.Command, error) {
	options := Options{
		Context: ctx,
	}

	var cmd = cobra.Command{
		Use:   `apicurito`,
		Short: `the apicurito operator`,
		Long:  `the apicurito operator takes care of installing and running apicurito on a cluster.`,
	}

	// Lets rexport the flags installed by the controller runtime, and make them a little less kube specific
	f := *flag.CommandLine.Lookup("kubeconfig")
	f.Name = "config"
	f.Usage = "path to the config file to connect to the cluster"
	cmd.PersistentFlags().AddGoFlag(&f)

	f = *flag.CommandLine.Lookup("master")
	f.Usage = "the address of the cluster API server."
	cmd.PersistentFlags().AddGoFlag(&f)

	cmd.AddCommand(newRunCommand(&options))
	return &cmd, nil
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
