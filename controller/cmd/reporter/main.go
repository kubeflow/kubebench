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
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/kubeflow/kubebench/controller/cmd/reporter/app"
)

func main() {
	if len(os.Args) < 2 {
		log.Errorf("Subcommand is required.")
		os.Exit(1)
	}

	reporterFactory := app.ReporterFactory{Name: os.Args[1]}

	opt, err := reporterFactory.NewReporterOption()
	if err != nil {
		log.Errorf("Failed to create reporter option: %s", err)
		os.Exit(1)
	}
	err = opt.AddFlags(os.Args[2:])
	if err != nil {
		log.Errorf("Failed to parse command: %s", err)
		os.Exit(1)
	}

	reporter, err := reporterFactory.NewReporter()
	if err != nil {
		log.Errorf("Failed to create reporter: %s", err)
		os.Exit(1)
	}

	if err := reporter.Run(opt); err != nil {
		log.Errorf("Reporter failed: %s", err)
		os.Exit(1)
	}
}
