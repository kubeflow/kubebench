// Copyright 2018 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/kubeflow/kubebench/controller/cmd/configurator/app"
)

func run(opt *app.AppOption) error {

	configurator := app.Configurator{
		FileOperator:      &(app.FileOperator{}),
		ManifestGenerator: &(app.ManifestGenerator{}),
		ManifestModifier:  &(app.ManifestModifier{}),
	}

	if err := configurator.Run(opt); err != nil {
		log.Errorf("Configurator failed to run: %s", err)
		return err
	}

	return nil
}

func main() {
	opt := app.NewAppOption()
	opt.AddFlags(flag.CommandLine)

	flag.Parse()

	if err := run(opt); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
}
