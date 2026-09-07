package menu

import fzf "github.com/junegunn/fzf/src"

// RunOptions wraps the real fzf options internally so that MenuRunner — and any
// custom implementation of it, e.g. a test fake — never needs to import
// the fzf package.
type RunOptions struct {
	Input  chan string
	Output chan string

	fzfOptions *fzf.Options
}

// MenuRunner defines the interface for running FZF with options and parsing
// arguments.
type MenuRunner interface {
	// Run executes FZF with the given options and returns the exit code.
	Run(options *RunOptions) (int, error)

	// Parse converts command line arguments to FZF options, optionally applying
	// defaults.
	Parse(defaults bool, args Args) (*RunOptions, error)
}

type defaultRunner struct{}

func (d *defaultRunner) Run(options *RunOptions) (int, error) {
	options.fzfOptions.Input = options.Input
	options.fzfOptions.Output = options.Output
	return fzf.Run(options.fzfOptions)
}

func (d *defaultRunner) Parse(def bool, args Args) (*RunOptions, error) {
	fzfOpts, err := fzf.ParseOptions(def, args)
	if err != nil {
		return nil, err
	}
	return &RunOptions{fzfOptions: fzfOpts}, nil
}
