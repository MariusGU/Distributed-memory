package main

import (
	"fmt"
	"sync"
)

type Array struct {
	Nums  []int
	Count int
}

func (arr *Array) add(num int) {
	arr.Nums[arr.Count] = num
	arr.Count++
}

type Counter struct {
	Count int
}

func (cnt *Counter) increase() {
	cnt.Count++
}

func main() {
	gavch := make(chan int, 10)
	spausdch1 := make(chan int)
	spausdch2 := make(chan int)
	var s1arr *Array
	var s2arr *Array
	cntg := 0
	cnts := 0

	var gavWaitGoup sync.WaitGroup

	go siunThread1(spausdch1, s1arr)
	go siunThread1(spausdch2, s2arr)

	go gavThread(&gavWaitGoup, gavch, spausdch1, spausdch2, cnts)
	gavWaitGoup.Add(1)

	for i := 0; i < 11; i++ {
		gavch <- i
		cntg++
	}
	for i := 11; i < 20; i++ {
		gavch <- i
		cntg++
	}
	if cntg == 20 && cnts == 20 {
		close(gavch)
	}
	gavWaitGoup.Wait()
}

func gavThread(wg *sync.WaitGroup, gvch <-chan int, sch1 chan<- int, sch2 chan<- int, cnt int) {
	for value := range gvch {
		if value%2 == 0 {
			sch1 <- value
			cnt++
			fmt.Println("cnt1", cnt)
		} else {
			sch2 <- value
			cnt++
			fmt.Println("cnt2", cnt)
		}
	}
	wg.Done()
}
func siunThread1(sch <-chan int, arr *Array) {
	for value := range sch {
		fmt.Println(value)
	}
}
func siunThread2(sch <-chan int, arr *Array) {
	for value := range sch {
		fmt.Println(value)
	}
}
