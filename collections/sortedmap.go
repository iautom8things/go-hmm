package collections

import (
	"errors"
	"math/Rand"
	"sort"
)

type SortedMap struct {
	m map[string]float64
	s []string
}

func CreateSortedMap() SortedMap {
	sm := SortedMap{}
	sm.m = make(map[string]float64)
	sm.s = make([]string, 0)
	return sm
}

func (sm *SortedMap) Has(key string) bool {
	_, ok := sm.m[key]
	return ok
}

func (sm *SortedMap) Add(key string, value float64) error {
	found := sm.Has(key)
	if found {
		return errors.New("key already seen")
	}

	sm.m[key] = value
	sm.s = append(sm.s, key)
	sort.Strings(sm.s)

	return nil
}

func (sm *SortedMap) Edit(key string, value float64) error {
	found := sm.Has(key)
	if !found {
		return errors.New("key not found")
	}

	sm.m[key] = value
	sort.Strings(sm.s)
	return nil
}

func (sm *SortedMap) Remove(key string) error {
	found := sm.Has(key)
	if !found {
		return errors.New("key not found")
	}

	delete(sm.m, key)
	i := sort.SearchStrings(sm.s, key)
	if i == len(sm.s) {
		sm.s = sm.s[:i]
	} else {
		copy(sm.s[:i], sm.s[i+1:])
	}

	return nil
}

func (sm *SortedMap) GetProbabilityOf(key string) (float64, error) {

	found := sm.Has(key)
	if !found {
		return 0.0, errors.New("key not found")
	}

	return sm.m[key] / sm.Total(), nil
}

func (sm *SortedMap) GetRandom() (string, error) {
	if len(sm.m) == 0 {
		return "", errors.New("empty map")
	}

	totalProb := sm.Total()
	r := rand.Float64()
	acc := 0.0
	for _, s := range sm.s {
		acc += sm.m[s] / totalProb
		if acc > r {
			return s, nil
		}
	}

	return "", errors.New("unexpected fall-through")
}

func (sm *SortedMap) Len() int {
	return len(sm.m)
}

func (sm *SortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *SortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func (sm *SortedMap) Total() float64 {
	totalProb := 0.0
	for _, relProb := range sm.m {
		totalProb += relProb
	}
	return totalProb
}
