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

import (
	"flag"

	corev1 "k8s.io/api/core/v1"
)

// ServerOption is the main options for the controller.
type ServerOption struct {
	Kubeconfig  string
	Config      string
	Namespace   string
	Threadiness int
}

// NewServerOption creates a new ServerOption
func NewServerOption() *ServerOption {
	return &ServerOption{}
}

// AddFlags adds flags to the specified FlagSet.
func (opt *ServerOption) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&opt.Kubeconfig, "kubeconfig", "", "The kubeconfig file to use (out of cluster only)")
	fs.StringVar(&opt.Config, "config", "", "The operator config file to use")

	fs.StringVar(&opt.Namespace, "namespace", corev1.NamespaceAll, "The namespace to "+
		"monitor kubebenchjobs. If not set, it monitors all namespaces cluster-wide, else it "+
		"only monitors given namespace.")

	fs.IntVar(&opt.Threadiness, "threadiness", 1, "Number of threads to process the main logic")
}
