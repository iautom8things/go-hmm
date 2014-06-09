package markov

import (
	"errors"
	"git.bigodev.com/mazubieta/go-hmm/collections"
)

type Model struct {
	States       []*State
	stateLookup  map[string]*State
	initialProbs collections.SortedMap
	currentState *State
}

func (m *Model) Initialize() error {
	m.initialProbs = collections.CreateSortedMap()
	m.stateLookup = make(map[string]*State)

	for _, s := range m.States {
		m.initialProbs.Add(s.Name, 1.0)
		m.stateLookup[s.Name] = s
	}

	if len(m.States) == 0 {
		return errors.New("empty states")
	}

	if m.initialProbs.Len() == 0 {
		return errors.New("empty probability")
	}

	rName, err := m.initialProbs.GetRandom()
	if err != nil {
		return err
	}

	rState, ok := m.stateLookup[rName]

	if !ok {
		return errors.New("state not found")
	}

	m.currentState = rState

	return nil
}

func (m *Model) TakeStep() (string, string, error) {
	currentStateName := m.currentState.Name
	emition, err1 := m.currentState.GetRandomEmition()

	if err1 != nil {
		return "", "", err1
	}

	nextState, err2 := m.currentState.GetRandomTransition()

	if err2 != nil {
		return "", "", err2
	}

	m.currentState = m.stateLookup[nextState]
	return currentStateName, emition, nil
}
