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

package app

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Reporter interface {
	Run(options ReporterOption) error
}

type ReporterOption interface {
	AddFlags(args []string) error
}

type ReporterFactory struct {
	Name string
}

func (rf *ReporterFactory) NewReporter() (Reporter, error) {
	var r Reporter
	switch rf.Name {
	case "csv":
		r = &CsvReporter{}
	default:
		err := fmt.Errorf("Unknown reporter: %s", rf.Name)
		log.Errorf("%s", err)
		return nil, err
	}
	return r, nil
}

func (rf *ReporterFactory) NewReporterOption() (ReporterOption, error) {
	var ro ReporterOption
	switch rf.Name {
	case "csv":
		ro = &CsvReporterOption{}
	default:
		err := fmt.Errorf("Unknown subcommand: %s", rf.Name)
		log.Errorf("%s", err)
		return nil, err
	}
	return ro, nil
}
