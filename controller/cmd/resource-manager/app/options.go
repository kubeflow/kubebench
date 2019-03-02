package app

import "flag"

type AppOption struct {
	Action     string
	Kubeconfig string
	InputFile  string
	OutputFile string
	Timeout    string
}

func (opt *AppOption) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&opt.Action, "action", "", "Action to take against the k8s resources.")
	fs.StringVar(&opt.Kubeconfig, "kubeconfig", "", "Kubeconfig file (out of cluster only).")
	fs.StringVar(&opt.InputFile, "input-file", "", "Path to the input file.")
	fs.StringVar(&opt.OutputFile, "output-file", "", "Path to the output file.")
	fs.StringVar(&opt.Timeout, "timeout", "15m", "Timeout for auto-watch.")
}
