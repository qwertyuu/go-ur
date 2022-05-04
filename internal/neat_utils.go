package gour

import (
	"os"

	"github.com/pkg/errors"
	"github.com/yaricom/goNEAT/v2/neat"
	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

func LoadOptionsAndGenome(contextPath, genomePath string) (*neat.Options, *genetics.Genome, error) {
	// Load context configuration
	configFile, err := os.Open(contextPath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open context file")
	}
	context, err := neat.LoadNeatOptions(configFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load NEAT options")
	}

	// Load start Genome
	genomeFile, err := os.Open(genomePath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open genome file")
	}
	r, err := genetics.NewGenomeReader(genomeFile, genetics.YAMLGenomeEncoding)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load NEAT genome")
	}
	genome, err := r.Read()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load NEAT genome")
	}
	return context, genome, nil
}
