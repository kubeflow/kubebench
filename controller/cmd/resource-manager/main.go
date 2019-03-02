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
	"io"
	"io/ioutil"
	"os"

	"github.com/kubeflow/kubebench/controller/cmd/resource-manager/app"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func run(opt *app.AppOption) error {

	// create input reader
	inputFile := opt.InputFile
	inputReader, err := os.Open(inputFile)
	if err != nil {
		log.Errorf("Failed to open input file: %s", inputFile)
		return err
	}
	defer inputReader.Close()

	// create output writer
	outputFile := opt.OutputFile
	var outputWriter io.Writer
	if outputFile == "" {
		outputWriter = ioutil.Discard
	} else {
		fileWriter, err := os.Create(outputFile)
		if err != nil {
			log.Errorf("Failed to create output file: %s", outputFile)
			return err
		}
		defer fileWriter.Close()
		outputWriter = fileWriter
	}

	// get k8s client config
	var config *rest.Config
	if opt.Kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", opt.Kubeconfig)
		if err != nil {
			log.Errorf("Failed to get config: %s", err)
			return err
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Errorf("Failed to get config: %s", err)
			return err
		}
	}

	// run resource manager
	resourceManager := app.NewResourceManager(config, inputReader, outputWriter)
	err = resourceManager.Run(opt)
	if err != nil {
		log.Errorf("Failed to run resource manager: %s", err)
		return err
	}

	return nil
}

func main() {
	opt := app.AppOption{}
	opt.AddFlags(flag.CommandLine)
	flag.Parse()

	if err := run(&opt); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
}
