package cli

import (
	"flag"
	"fmt"
	"strings"
)

type ActionOpt string

const (
	PopulateIndex ActionOpt = "populate"
	CreateIndex   ActionOpt = "create"
)

func (a *ActionOpt) String() string {
	return string(*a)
}
func (a *ActionOpt) Set(value string) error {
	switch normalizedValue := strings.ToLower(value); normalizedValue {
	case string(PopulateIndex), string(CreateIndex):
		*a = ActionOpt(value)
		return nil
	default:
		return fmt.Errorf("must be %q or %q", PopulateIndex, CreateIndex)
	}
}

type Opts struct {
	IndexAction   ActionOpt
	ReadChunkSize int
}

func Parse() Opts {
	opts := Opts{}
	flag.IntVar(&opts.ReadChunkSize, "chunkSize", 1000, "The size of chunk to use for folder reads")
	flag.Var(&opts.IndexAction, "action", "Specific INDEX action that will be executed (create|populate)")
	flag.Parse()
	return opts
}
