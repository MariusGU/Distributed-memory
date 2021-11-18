package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"text/tabwriter"
)

const max = 40
const divider = 1      //default - 10
const modifier = 2     //default - 2
const monitorSize = 10 //default - 10
const threadsCount = 4
const file = 1

//Thing is for storing string, int and float values in single structure
//Result value is added later to store filter parameter
type Preke struct {
	Pavadinimas string  `json:"pavadinimas"`
	Kiekis      int     `json:"kiekis"`
	Kaina       float32 `json:"kaina"`
	Result      int
}

//Array of things
type Array struct {
	Prekes []Preke
	Count  int
}

func (array *Array) add(preke Preke) {
	var i int
	for i = int(array.Count - 1); i >= 0 && array.Prekes[i].Result > preke.Result; i-- {
		(*array).Prekes[i+1] = (*array).Prekes[i]
	}
	(*array).Prekes[i+1] = preke
	(*array).Count++
}

func mainThread() {

	var data []Preke
	data = readJSON(fmt.Sprintf("IFK-9_GudzinskasM_L2_dat_%d.json", file))
	threadsCount := len(data) / 4

	workerch := make(chan Preke, monitorSize)
	resultsch := make(chan Preke)
	datach := make(chan Preke, 1)
	mainch := make(chan Array)

	var workersWaitGroup sync.WaitGroup
	for i := 0; i < threadsCount; i++ {
		workersWaitGroup.Add(1)
		go workerThread(&workersWaitGroup, i, workerch, resultsch)
	}
	go resultsThread(resultsch, mainch)
	go dataThread(datach, workerch)

	for _, preke := range data {
		datach <- preke
	}
	close(datach)
	workersWaitGroup.Wait()
	close(resultsch)

	var results Array
	for result := range mainch {
		results = result
	}

	fmt.Println("Buvo - ", len(data), "; Atfiltruota - ", results.Count)

	printToFile(data, results)
}

func workerThread(wg *sync.WaitGroup, id int, ch <-chan Preke, results chan<- Preke) {
	for value := range ch {
		primeCount := filterCondition(value)

		fmt.Println("gija - ", id, ", pavadinimas - ", value.Pavadinimas, " Kodas - ", primeCount)

		if (primeCount % modifier) == 0 {
			value.Result = primeCount
			results <- value
		}
	}
	wg.Done()
}

func dataThread(data <-chan Preke, worker chan<- Preke) {
	for value := range data {
		worker <- value
	}
	close(worker)
}

func resultsThread(resultsch <-chan Preke, mainch chan<- Array) {
	results := Array{Prekes: make([]Preke, max)}
	for value := range resultsch {
		results.add(value)
	}

	mainch <- results
	close(mainch)
}

func readJSON(path string) []Preke {
	file, _ := ioutil.ReadFile(path)
	var data []Preke
	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func primeCount(number int) int {
	var n = 3
	var is = true
	var count int

	for n < number {
		is = true
		if n%2 == 0 {
			is = false
		}
		for i := 3; i < (n / 3); i += 2 {
			if (n % i) == 0 {
				is = false
				break
			}
		}
		if is {
			count++
		}
		n++
	}
	return count
}

func filterCondition(preke Preke) int {
	sum := 0
	for _, value := range preke.Pavadinimas {
		sum += int(value)
	}
	sum += int(preke.Kaina)
	sum = sum * preke.Kiekis / divider
	primeCount := primeCount(sum)
	return primeCount
}

func printToFile(data []Preke, result Array) {
	f, _ := os.Create("./IFK-9_GudzinskasM_L2_rez.txt")
	w := tabwriter.NewWriter(f, 0, 0, 4, ' ', tabwriter.AlignRight|tabwriter.Debug)

	defer f.Close()

	fmt.Fprintln(w, "REZULTATAI")
	fmt.Fprintln(w, "-----------\t----------\t---------------\t------------------------------\t")
	fmt.Fprintln(w, "Pavadinimas\tKiekis\tkaina\tKodas\t")
	fmt.Fprintln(w, "-----------\t----------\t---------------\t------------------------------\t")

	for i := 0; i < int(result.Count); i++ {
		s := fmt.Sprintf("%s\t%d\t%f\t%d\t", result.Prekes[i].Pavadinimas, result.Prekes[i].Kiekis, result.Prekes[i].Kaina, result.Prekes[i].Result)
		fmt.Fprintln(w, s)
	}

	fmt.Fprintln(w, "DUOMENYS")
	fmt.Fprintln(w, "-----------\t----------\t---------------\t")
	fmt.Fprintln(w, "Pavadinimas\tKiekis\tKaina\t")
	fmt.Fprintln(w, "-----------\t----------\t---------------\t")
	for _, value := range data {
		s := fmt.Sprintf("%s\t%d\t%f\t", value.Pavadinimas, value.Kiekis, value.Kaina)
		fmt.Fprintln(w, s)
	}

	w.Flush()
}

func main() {
	mainThread()
}
