package main

import (
	"fmt"
	"time"

	"chunky"
)

var count int
var gCount int
var parrellelItemsCount int
var cap = 100
var pjobs int
var generalItemsCount int

type cType chan int

func main() {
	fmt.Println("Chunky, scheduled channel buffer jobs, split the que!")

	cBuf := make(cType, cap)
	factor := 6
	chunker := chunky.NewChunk(cBuf, Sump, factor)
	done := make(chan int)
	tChunk := time.Now()
	go func() {
		for parrellelItemsCount < 1000 {
			//if len(cBuf) == 0 {
			Pump(cBuf)
			//	}
		}
		done <- 1
	}()

	go chunker.Schedule()
	<-done
	for parrellelItemsCount < 1000 {
	}
	totalPComplete := parrellelItemsCount
	tdone := time.Now().Sub(tChunk)

	chunker.Close()

	fmt.Println("------------------------")

	cBuf = make(cType, cap)
	tRegular := time.Now()

	go func() {
		for generalItemsCount < 1000 {
			Pump(cBuf)
		}
		done <- 1
	}()
	go GeneralSump(cBuf)
	<-done
	for generalItemsCount < 1000 {
	}
	totalGeneralComplete := generalItemsCount
	tRegDone := time.Now().Sub(tRegular)
	fmt.Printf("Chunker Done in %v seconds, processed %d jobs with total of %d values\n", tdone.Seconds(), pjobs, totalPComplete)

	fmt.Printf("Normal Buffered Channel Done in %v seconds, processed 1 job with total of %d values\n", tRegDone.Seconds(), totalGeneralComplete)
}

func Pump(b cType) {

	for i := 0; i < cap; i++ {
		b <- i
	}

}

func GeneralSump(b interface{}) {
	gCount++
	jobs := gCount
	switch val := b.(type) {
	case cType:
		for v := range val {
			generalItemsCount++
			fmt.Printf("general job %d : %d\n", jobs, v)
		}

	}
}

func Sump(b interface{}) {
	count++
	pjobs = count
	switch val := b.(type) {
	case cType:
		length := len(val)
		for i := 0; i < length; i++ {
			select {
			case v := <-val:
				parrellelItemsCount++
				//if jobs <= 10 {
				fmt.Printf("parallel job %d : %d\n", pjobs, v)
				//}
			default:
			}

		}

	}
}
