package markov

import (
	"errors"
	"git.bigodev.com/mazubieta/go-hmm/collections"
	"math/Rand"
)

type State struct {
	Name      string
	neighbors collections.SortedMap
	emitions  collections.SortedMap
}

func (s *State) addEmition(symbol string, rel_prob float64) error {
	// check if we already have emition
	found := s.emitions.Has(symbol)
	if !found {
		return errors.New("State: Emition already seen")
	}

	s.emitions.Add(symbol, rel_prob)
	return nil
}

func (s State) removeEmition(symbol string) error {
	// check if we have emition
	found := s.emitions.Has(symbol)
	if !found {
		return errors.New("State: Emition does not exist")
	}

	s.emitions.Remove(symbol)
	return nil
}

func (s State) getEmitionProbability(symbol string) (float64, error) {
	// check if we have emition
	found := s.emitions.Has(symbol)
	if !found {
		return 0.0, errors.New("State: Emition does not exist")
	}

	totalProb := s.emitions.Total()
	relProb, _ := s.emitions.Get(symbol)
	return relProb / totalProb, nil
}

func (s State) getRandomEmition() (float64, error) {
	r := rand.Float64()
	r += 1
	totalProb := s.emitions.Total()
	totalProb += 1

	return 0.0, nil
}
