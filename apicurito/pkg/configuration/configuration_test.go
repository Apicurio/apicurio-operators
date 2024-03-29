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

package configuration

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadFromFile(t *testing.T) {
	type args struct {
		config string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			"When loading a config file, all parameter should match",
			args{config: "../../build/conf/config_test.yaml"},
			&Config{UiImage: "apicurio/apicurito-ui:1.2.4", GeneratorImage: "apicurio/fuse-apicurito-generator:latest"},
			false,
		},
	}
	for _, tt := range tests {
		p := &Config{}
		t.Run(tt.name, func(t *testing.T) {
			err := p.loadFromFile(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("loadFromFile() got = %v, want %v", p, tt.want)
			}
		})
	}
}

func TestConfig_setPropertiesFromEnv(t *testing.T) {
	type fields struct {
		UiImage        string
		GeneratorImage string
	}
	tests := []struct {
		name       string
		fields     fields
		env        map[string]string
		wantErr    bool
		wantConfig *Config
	}{
		{
			name:   "When env is provided it should replace the existing image",
			fields: fields{UiImage: "someimage", GeneratorImage: "someotherimage"},
			env: map[string]string{
				"RELATED_IMAGE_APICURITO": "image_from_env",
				"RELATED_IMAGE_GENERATOR": "image_from_env_g",
			},
			wantConfig: &Config{UiImage: "image_from_env", GeneratorImage: "image_from_env_g"},
			wantErr:    false,
		},
		{
			name:   "When env is provided and no images is set, env should prevail",
			fields: fields{},
			env: map[string]string{
				"RELATED_IMAGE_APICURITO": "image_from_env",
				"RELATED_IMAGE_GENERATOR": "image_from_env_g",
			},
			wantConfig: &Config{UiImage: "image_from_env", GeneratorImage: "image_from_env_g"},
			wantErr:    false,
		},
		{
			name:       "When no env is provided, the initial value should prevail",
			fields:     fields{UiImage: "someimage", GeneratorImage: "someotherimage"},
			env:        map[string]string{},
			wantConfig: &Config{UiImage: "someimage", GeneratorImage: "someotherimage"},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			c := &Config{
				UiImage:        tt.fields.UiImage,
				GeneratorImage: tt.fields.GeneratorImage,
			}
			if err := c.setPropertiesFromEnv(); (err != nil) != tt.wantErr {
				t.Errorf("setPropertiesFromEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantConfig, c)

			for k := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}
