package markov

import (
	"errors"
	"git.bigodev.com/mazubieta/go-hmm/collections"
)

type Model struct {
	CurrentState *State
	States       []*State
	stateLookup  map[string]*State
	initialProbs collections.SortedMap
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

	m.CurrentState = rState

	return nil
}

func (m *Model) TakeStep() (string, string, error) {
	currentStateName := m.CurrentState.Name
	emition, err1 := m.CurrentState.GetRandomEmition()

	if err1 != nil {
		return "", "", err1
	}

	nextState, err2 := m.CurrentState.GetRandomTransition()

	if err2 != nil {
		return "", "", err2
	}

	m.CurrentState = m.stateLookup[nextState]
	return currentStateName, emition, nil
}

func (m *Model) GetInitialProb(s string) float64 {
	p, _ := m.initialProbs.GetProbabilityOf(s)
	return p
}

func (m *Model) GetState(s string) *State {
	return m.stateLookup[s]
}
