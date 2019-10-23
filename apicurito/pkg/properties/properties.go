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

package properties

import (
    "encoding/json"
    "github.com/apicurio/apicurio-operators/apicurito/pkg/configuration"
    "io/ioutil"
    "k8s.io/apimachinery/pkg/util/yaml"
    "strings"
)

type Properties struct {
    Image string
}

// Returns all processed properties for Apicurito
func GetProperties() (*Properties, error) {
    c := configuration.GetConfiguration()
    p, err := loadFromFile(c.ConfigFile)
    if err != nil {
        return nil, err
    }

    return p, nil
}

// Load configuration from config file. Config file is expected to be a yaml
// The returned configuration is parsed to JSON and returned as Properties
func loadFromFile(config string ) (*Properties, error) {
    data, err := ioutil.ReadFile(config)
    if err != nil {
        return nil, err
    }

    if strings.HasSuffix(config, ".yaml") || strings.HasSuffix(config, ".yml") {
        data, err = yaml.ToJSON(data)
        if err != nil {
            return nil, err
        }
    }

    p := &Properties{}
    if err := json.Unmarshal(data, p); err != nil {
        return nil, err
    }

    return p, nil
}