package markov

import (
	"math"
)

type Posterior struct {
	Model
}

func (p *Posterior) Decode(observations []string) []string {

	// setup return decoding
	decoding := make([]string, len(observations))
	// mapping of state names to lists of previous state probabilities
	vProb := map[string][]float64{}

	// iterate over all states in the model using their initial
	// probabilities and their emition probabilites to determine
	// the first decoded state
	var bestState string
	maxProb := math.Inf(-1)

	end := len(observations) - 1
	for _, s := range p.States {
		vProb[s.Name] = make([]float64, len(observations))
		e, _ := s.GetEmitionProbability(observations[end])
		t := p.GetInitialProb(s.Name)
		v_j := math.Exp(math.Log(e) + math.Log(t))
		vProb[s.Name][end] = v_j
		if v_j > maxProb {
			maxProb = v_j
			bestState = s.Name
		}
	}
	decoding[end] = bestState

	// iterate over the rest of the observed states using the
	// previous state's probabilities
	for j := end - 1; j > 0; j-- {
		obs_j := observations[j]

		// iterate over all states to determine the
		// state with the max probability that observations[j]
		// was generated by it
		var bestJState string
		maxJProb := math.Inf(-1)
		for _, state_j := range p.States {
			e_j, _ := state_j.GetEmitionProbability(obs_j)

			// second iterate over all states to determine the best
			// pervious state
			maxIProb := math.Inf(-1)
			for _, state_i := range p.States {
				v_i := vProb[state_i.Name][j+1]
				trans_j2i_prob, _ := state_j.GetTransitionProbability(state_i.Name)

				mul := math.Exp(math.Log(v_i) + math.Log(trans_j2i_prob))

				if mul > maxIProb {
					maxIProb = mul
				}
			}

			v_j := math.Exp(math.Log(e_j) + maxIProb)

			vProb[state_j.Name][j] = v_j
			if v_j > maxJProb {
				maxJProb = v_j
				bestJState = state_j.Name
			}
		}
		decoding[j] = bestJState
	}

	return decoding
}
