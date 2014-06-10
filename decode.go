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
	unfair.AddEmition("6", 2.0)
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
			v[s.Name] = make([]float64, len(observation))
			e, _ := s.GetEmitionProbability(string(observation[0]))
			t := model.GetInitialProb(s.Name)
			v_j := e * t
			v[s.Name][0] = v_j
			fmt.Println(s.Name, e, t, v[s.Name][0])
			if v_j > max_val {
				max_val = v_j
				max_state = s.Name
			}
		}
		pi[0] = max_state

		for j := 1; j < len(observation); j++ {
			obs_j := string(observation[j])
			var max_state string
			max_val := math.Inf(-1)
			for _, state_j := range model.States {
				e_j, _ := state_j.GetEmitionProbability(obs_j)
				vi_sum := 0.0
				for _, state_i := range model.States {
					v_prev := v[state_i.Name][j-1]
					t_prob, _ := state_i.GetTransitionProbability(obs_j)
					vi_sum += v_prev * t_prob
				}
				v_j := e_j * vi_sum
				v[state_j.Name][j] = v_j
				if v_j > max_val {
					max_val = v_j
					max_state = state_j.Name
				}
			}
			pi[j] = max_state
		}

		fmt.Print(pi)

	}
}
