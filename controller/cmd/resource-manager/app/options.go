package app

import "flag"

type AppOption struct {
	Action     string
	Kubeconfig string
	InputData  string
	InputFile  string
	OutputFile string
	NumCopies  int
	Timeout    string
}

func (opt *AppOption) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&opt.Action, "action", "", "Action to take against the k8s resources.")
	fs.StringVar(&opt.Kubeconfig, "kubeconfig", "", "Kubeconfig file (out of cluster only).")
	fs.StringVar(&opt.InputData, "input-data", "", "Input data.")
	fs.StringVar(&opt.InputFile, "input-file", "", "Path to the file containing input data.")
	fs.StringVar(&opt.OutputFile, "output-file", "", "Path to the output file.")
	fs.IntVar(&opt.NumCopies, "num-copies", 1, "Number of copies to create.")
	fs.StringVar(&opt.Timeout, "timeout", "15m", "Timeout for auto-watch.")
}
