package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"git.bigodev.com/mazubieta/go-hmm/markov"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/Rand"
	"os"
	"strconv"
	"time"
)

var (
	TRACE   *log.Logger
	INFO    *log.Logger
	WARNING *log.Logger
	ERROR   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	TRACE = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	INFO = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	WARNING = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ERROR = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	var seed int64
	var numIter int
	var transProb float64
	var inputFile, outputFile string

	// set flags for simulation
	flag.Int64Var(&seed, "s", time.Now().Unix(), "Seed for PRNG (Default: time.Time())")
	flag.IntVar(&numIter, "n", 1000, "Number of Time Steps (Default: 1000)")
	flag.Float64Var(&transProb, "t", 0.1, "Transition Porbability (Default: 0.1)")
	flag.StringVar(&inputFile, "i", "output_hmm", "Output file name(Default: \"output_hmm\")")
	flag.StringVar(&outputFile, "o", "output_hmm", "Output file name(Default: \"output_hmm\")")
	flag.Parse()

	// seed PRNG
	rand.Seed(seed)

	// create fair die
	fair := markov.CreateState("f")
	fair.AddEmition("1", 1.0)
	fair.AddEmition("2", 1.0)
	fair.AddEmition("3", 1.0)
	fair.AddEmition("4", 1.0)
	fair.AddEmition("5", 1.0)
	fair.AddEmition("6", 1.0)
	fair.AddNeighbor("f", 1.0-transProb)
	fair.AddNeighbor("u", transProb)

	// craete unfair die
	unfair := markov.CreateState("u")
	unfair.AddEmition("1", 1.0)
	unfair.AddEmition("2", 1.0)
	unfair.AddEmition("3", 1.0)
	unfair.AddEmition("4", 1.0)
	unfair.AddEmition("5", 1.0)
	unfair.AddEmition("6", 5.0)
	unfair.AddNeighbor("f", transProb)
	unfair.AddNeighbor("u", 1.0-transProb)

	// lookup map of states
	model := markov.Model{
		States: []*markov.State{&fair, &unfair},
	}

	// randomly start at a state (assumed equal probability of start state)
	err := model.Initialize()
	if err != nil {
		panic(err)
	}

	var emitionBuffer, stateBuffer bytes.Buffer

	for i := 0; i < numIter; i++ {
		s, e, err2 := model.TakeStep()
		if err2 != nil {
			panic(err2)
		}
		emitionBuffer.WriteString(e)
		stateBuffer.WriteString(s)
	}

	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	nRecords := len(lines)
	seeds := make([]int64, nRecords)
	tprobs := make([]float64, nRecords)
	observations := make([]string, nRecords)
	hiddenStates := make([]string, nRecords)

	for i, line := range lines {
		fmt.Println(i, line)
		seeds[i], err = strconv.ParseInt(line[0], 10, 64)
		tprobs[i], err = strconv.ParseFloat(line[1], 64)
		observations[i] = string(line[2])
		hiddenStates[i] = string(line[3])
	}

	for i := 0; i < 1; i++ {
		observation := observations[i]
		INFO.Println("workong on observation: ", observation)
		states := hiddenStates[i]

		if len(observation) != len(states) {
			fmt.Println("Mismatch record [", i, "] ...skipping...")
			continue
		}

		pi := make([]string, len(observation))
		v := map[string][]float64{}

		var max_state string
		max_val := math.Inf(-1)
		for _, s := range model.States {
			INFO.Println("working on initial state: ", s.Name)
			v[s.Name] = make([]float64, len(observation))
			e, err := s.GetEmitionProbability(string(observation[0]))
			if err != nil {
				panic(err)
			}
			t := model.GetInitialProb(s.Name)
			v_j := math.Exp(math.Log(e) + math.Log(t))
			v[s.Name][0] = v_j
			INFO.Println("state:", s.Name, "e_prob:", e, "t_prob:", t, "prev_prob:", v[s.Name][0])
			if v_j > max_val {
				max_val = v_j
				max_state = s.Name
			}
		}
		pi[0] = max_state

		INFO.Println("ITER OVER OBS")
		for j := 1; j < len(observation); j++ {
			obs_j := string(observation[j])
			INFO.Println("obs: <", obs_j, ">")

			var max_j_state string
			max_j_val := math.Inf(-1)
			for _, state_j := range model.States {
				INFO.Println("Trying to beat maxes: ", max_j_state, "v:", max_j_val)
				e_j, _ := state_j.GetEmitionProbability(obs_j)

				WARNING.Println("state: <", state_j.Name, "> emit_p: ", e_j)

				max_i_val := math.Inf(-1)
				for _, state_i := range model.States {
					v_prev := v[state_i.Name][j-1]
					t_prob, _ := state_i.GetTransitionProbability(state_j.Name)

					mul := math.Exp(math.Log(v_prev) + math.Log(t_prob))

					INFO.Println("State: ", state_i.Name, " v_prev: ", v_prev, "t_prob:", t_prob, "mul: ", mul)

					if mul > max_i_val {
						max_i_val = mul
					}
				}

				v_j := math.Exp(math.Log(e_j) + max_i_val)

				WARNING.Println("See if we beat maxes with: v_j: <", v_j)

				v[state_j.Name][j] = v_j
				if v_j > max_j_val {
					INFO.Println("OLD: ", max_j_state, "v:", max_j_val)
					max_j_val = v_j
					max_j_state = state_j.Name
				}
			}
			pi[j] = max_j_state
		}

		fmt.Println(states)
		fmt.Println(pi)

	}
}
