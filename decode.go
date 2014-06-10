package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"git.bigodev.com/mazubieta/go-hmm/markov"
	"math/Rand"
	"os"
	"strconv"
	"time"
)

func main() {
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
	viterbi := markov.Viterbi{
		markov.Model{States: []*markov.State{&fair, &unfair}},
	}

	viterbi.Initialize()

	// read csv file of all tests
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
	observations := make([][]string, nRecords)
	hiddenStates := make([][]string, nRecords)

	var numObservations int64
	for i, line := range lines {
		seeds[i], _ = strconv.ParseInt(line[0], 10, 64)
		tprobs[i], _ = strconv.ParseFloat(line[1], 64)
		numObservations, _ = strconv.ParseInt(line[2], 10, 0)
		observations[i] = line[3 : 3+numObservations]
		hiddenStates[i] = line[3+numObservations : 3+2*numObservations]
	}

	n := int(numObservations)

	for i := 0; i < 1; i++ {
		observation := observations[i]
		states := hiddenStates[i]
		fair.EditNeighbor("f", 1.0-tprobs[i])
		fair.EditNeighbor("u", tprobs[i])
		unfair.EditNeighbor("f", tprobs[i])
		unfair.EditNeighbor("u", 1.0-tprobs[i])

		if len(observation) != len(states) {
			fmt.Println("Mismatch record [", i, "] ...skipping...")
			continue
		}

		// decode via Viterbi
		pi := viterbi.Decode(observation)

		var stateBuffer, piBuffer bytes.Buffer
		for j := 0; j < n; j++ {
			piBuffer.WriteString(pi[j])
			stateBuffer.WriteString(states[j])
		}
		fmt.Println(stateBuffer.String())
		fmt.Println(piBuffer.String())

	}
}
