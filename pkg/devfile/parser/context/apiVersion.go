//
// Copyright 2022 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/devfile/library/pkg/devfile/parser/data"
	"github.com/pkg/errors"
	"k8s.io/klog"
)

// SetDevfileAPIVersion returns the devfile APIVersion
func (d *DevfileCtx) SetDevfileAPIVersion() error {

	// Unmarshal JSON into map
	var r map[string]interface{}
	err := json.Unmarshal(d.rawContent, &r)
	if err != nil {
		return errors.Wrapf(err, "failed to decode devfile json")
	}

	// Get "schemaVersion" value from map for devfile V2
	schemaVersion, okSchema := r["schemaVersion"]
	var devfilePath string
	if d.GetAbsPath() != "" {
		devfilePath = d.GetAbsPath()
	} else if d.GetURL() != "" {
		devfilePath = d.GetURL()
	}

	if okSchema {
		// SchemaVersion cannot be empty
		if schemaVersion.(string) == "" {
			return fmt.Errorf("schemaVersion in devfile: %s cannot be empty", devfilePath)
		}
	} else {
		return fmt.Errorf("schemaVersion not present in devfile: %s", devfilePath)
	}

	// Successful
	// split by `-` and get the first substring as schema version, schemaVersion without `-` won't get affected
	// e.g. 2.2.0-latest => 2.2.0, 2.2.0 => 2.2.0
	d.apiVersion = strings.Split(schemaVersion.(string), "-")[0]
	klog.V(4).Infof("devfile schemaVersion: '%s'", d.apiVersion)
	return nil
}

// GetApiVersion returns apiVersion stored in devfile context
func (d *DevfileCtx) GetApiVersion() string {
	return d.apiVersion
}

// IsApiVersionSupported return true if the apiVersion in DevfileCtx is supported
func (d *DevfileCtx) IsApiVersionSupported() bool {
	return data.IsApiVersionSupported(d.apiVersion)
}
