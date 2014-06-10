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

func CreateState(name string) State {
	state := State{Name: name}
	state.neighbors = collections.CreateSortedMap()
	state.emitions = collections.CreateSortedMap()
	return state
}

func (s *State) AddEmition(symbol string, relProb float64) error {
	found := s.emitions.Has(symbol)
	if found {
		return errors.New("Emition already seen")
	}

	s.emitions.Add(symbol, relProb)
	return nil
}

func (s *State) RemoveEmition(symbol string) error {
	// check if we have emition
	found := s.emitions.Has(symbol)
	if !found {
		return errors.New("Emition does not exist")
	}

	s.emitions.Remove(symbol)
	return nil
}

func (s *State) GetEmitionProbability(symbol string) (float64, error) {
	// check if we have emition
	found := s.emitions.Has(symbol)
	if !found {
		return 0.0, errors.New("State: Emition does not exist")
	}

	p, _ := s.emitions.GetProbabilityOf(symbol)
	return p, nil
}

func (s *State) GetRandomEmition() (string, error) {
	e, err := s.emitions.GetRandom()
	return e, err
}

func (s *State) AddNeighbor(state string, relProb float64) error {
	found := s.neighbors.Has(state)
	if found {
		return errors.New("neighbor already exists")
	}

	s.neighbors.Add(state, relProb)
	return nil
}

func (s *State) RemoveNeighbor(state string) error {
	// check if we have emition
	found := s.neighbors.Has(state)
	if !found {
		return errors.New("neihbor does not exist")
	}

	s.neighbors.Remove(state)
	return nil
}

func (s *State) GetTransitionProbability(state string) (float64, error) {
	// check if we have emition
	found := s.neighbors.Has(state)
	if !found {
		return 0.0, errors.New("transition does not exist")
	}

	p, _ := s.neighbors.GetProbabilityOf(state)
	return p, nil
}

func (s *State) GetRandomTransition() (string, error) {
	t, err := s.neighbors.GetRandom()
	return t, err
}
