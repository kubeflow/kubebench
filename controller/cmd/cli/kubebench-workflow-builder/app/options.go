// Copyright 2019 The Kubeflow Authors
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

import "flag"

type AppOption struct {
	Config          string
	Manifest        string
	Output          string
	NoDefaultConfig bool
}

func (opt *AppOption) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&opt.Config, "config", "", "Path to the Kubebench config file.")
	fs.StringVar(&opt.Manifest, "manifest", "", "Path to the KubebenchJob manifest file.")
	fs.StringVar(&opt.Output, "output", "", "Path to the output file.")
	fs.BoolVar(&opt.NoDefaultConfig, "no-default-config", false, "Do not apply default config.")
}
