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
	"flag"
)

type AppOption struct {
	Config             string
	ManifestOutput     string
	ExperimentIDOutput string
}

func NewAppOption() *AppOption {
	opt := AppOption{}
	return &opt
}

func (opt *AppOption) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&opt.Config, "config", "", "The configuration file.")
	fs.StringVar(&opt.ManifestOutput, "manifest-output", "kf-job-manifest.yaml",
		"The output manifest file.")
	fs.StringVar(&opt.ExperimentIDOutput, "experiment-id-output", "experiment-id",
		"The output experiment id file.")
}
