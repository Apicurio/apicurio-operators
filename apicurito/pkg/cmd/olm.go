/*
 * Copyright (C) 2020 Red Hat, Inc.
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
	"flag"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"

	"github.com/apicurio/apicurio-operators/apicurito/tools/run"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/spf13/cobra"
)

type csv struct {
	*Options
}

func newOlmCommand(parent *Options) *cobra.Command {
	options := csv{Options: parent}
	cmd := cobra.Command{
		Use:   "olm",
		Short: "generates bundle files for OLM installation",
		Run: func(_ *cobra.Command, _ []string) {
			exitOnError(options.run())
		},
	}

	cmd.PersistentFlags().StringVarP(&configuration.ConfigFile, "config", "", "/conf/config.yaml", "path to the operator configuration file.")
	cmd.PersistentFlags().AddFlagSet(zap.FlagSet())
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return &cmd
}

func (c csv) run() error {
	return run.Run()
}
