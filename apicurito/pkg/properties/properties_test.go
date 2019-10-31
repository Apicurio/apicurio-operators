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
	"reflect"
	"testing"

	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
)

func TestGetProperties(t *testing.T) {
	type args struct {
		apicurito *v1alpha1.Apicurito
	}
	tests := []struct {
		name    string
		args    args
		want    *Properties
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProperties(tt.args.apicurito)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProperties() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadFromApicurito(t *testing.T) {
	type args struct {
		apicurito *v1alpha1.Apicurito
	}
	tests := []struct {
		name    string
		args    args
		want    *Properties
		wantErr bool
	}{
		{
			"When reading from apicurito, all fields should be loaded",
			args{&v1alpha1.Apicurito{Spec: v1alpha1.ApicuritoSpec{Image: "apicurio/apicurito-ui"}}},
			&Properties{Image: "apicurio/apicurito-ui"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadFromApicurito(tt.args.apicurito)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromApicurito() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFromApicurito() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadFromFile(t *testing.T) {
	type args struct {
		config string
	}
	tests := []struct {
		name    string
		args    args
		want    *Properties
		wantErr bool
	}{
		{
			"When loading a config file, all parameter should match",
			args{config: "../../build/conf/config_test.yaml"},
			&Properties{Image: "apicurio/apicurito-ui"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadFromFile(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeProperties(t *testing.T) {
	type args struct {
		config    Properties
		apicurito Properties
	}
	tests := []struct {
		name    string
		args    args
		want    *Properties
		wantErr bool
	}{
		{
			"When leaving Image in the CR empty, the default value from config should be taken",
			args{config: Properties{Image: "apicurito/from-config"}, apicurito: Properties{}},
			&Properties{"apicurito/from-config"},
			false,
		},
		{
			"When Image is present in the CR, the default value should be overwritten by the value from CR",
			args{config: Properties{Image: "apicurito/from-config"}, apicurito: Properties{Image: "apicurito/from-cr"}},
			&Properties{"apicurito/from-cr"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeProperties(tt.args.config, tt.args.apicurito)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeProperties() got = %v, want %v", got, tt.want)
			}
		})
	}
}
