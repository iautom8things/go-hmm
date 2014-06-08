package markov

import (
	"errors"
	"git.bigodev.com/mazubieta/go-hmm/collections"
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

func (s State) getRandomEmition() (string, error) {
	return s.emitions.GetRandom()
}

func (s *State) addNeighbor(state string, rel_prob float64) error {
	// check if we already have emition
	found := s.neighbors.Has(state)
	if !found {
		return errors.New("neighbor already exists")
	}

	s.neighbors.Add(state, rel_prob)
	return nil
}

func (s State) removeNeighbor(state string) error {
	// check if we have emition
	found := s.neighbors.Has(state)
	if !found {
		return errors.New("neihbor does not exist")
	}

	s.neighbors.Remove(state)
	return nil
}

func (s State) getTransitionProbability(state string) (float64, error) {
	// check if we have emition
	found := s.neighbors.Has(state)
	if !found {
		return 0.0, errors.New("transition does not exist")
	}

	totalProb := s.neighbors.Total()
	relProb, _ := s.neighbors.Get(state)
	return relProb / totalProb, nil
}

func (s State) getRandomTransition() (string, error) {
	return s.neighbors.GetRandom()
}
