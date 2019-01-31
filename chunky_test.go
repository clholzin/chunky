package chunky

import (
	"testing"
	"time"
)

var count int
var totalCountedItems int

type cType chan int

func Sump(b interface{}) {
	count++
	//cc := count
	switch val := b.(type) {
	case cType:
		length := len(val)
		for i := 0; i < length; i++ {
			select {
			case <-val:
				totalCountedItems++
				//fmt.Printf("Job-%d : %d\n", cc, v)
			default:
			}

		}

	}
}

func Pump(b cType, count int) {
	for i := 0; i < count; i++ {
		b <- i
	}
}

func TestNewChunk(t *testing.T) {
	count = 0
	jobs := 10
	cBuf := make(cType, jobs)
	chunker := NewChunk(cBuf, Sump, 2)
	done := make(chan int)
	go func() {

		for count < jobs {
			if len(cBuf) == 0 {
				Pump(cBuf, 10)

			}
		}
		done <- 1
	}()

	go chunker.Schedule()
	<-done

	for count >= jobs && len(cBuf) > 0 {

	}
	chunker.Close()
}

func TestChunkRead(t *testing.T) {
	cBuf := make(cType, 2)
	chunker := NewChunk(cBuf, Sump, 2)
	cBuf <- 1
	if _, ok := chunker.Read(); !ok {
		t.Fail()
	}
}

func TestChunkFactor(t *testing.T) {
	cBuf := make(cType, 2)
	chunker := NewChunk(cBuf, Sump, 2)
	if d := chunker.Factor(5); d != 5 {
		t.Fail()
	}
}

func TestChunkWrite(t *testing.T) {
	cBuf := make(cType, 2)
	chunker := NewChunk(cBuf, Sump, 2)
	chunker.Write(2)
	if len(cBuf) != 1 {
		t.Fail()
	}
}

func TestChunkLen(t *testing.T) {
	cBuf := make(cType, 2)
	chunker := NewChunk(cBuf, Sump, 2)
	cBuf <- 1
	if chunker.Len() != 1 {
		t.Fail()
	}

}

func BenchmarkChunk100(b *testing.B) {
	t := time.Now()
	count = 0
	totalCountedItems = 0
	parrallelProcesses := 2
	capacity := 100
	cBuf := make(cType, capacity)
	chunker := NewChunk(cBuf, Sump, parrallelProcesses)
	done := make(chan int)

	go func() {
		for totalCountedItems < 1000 { //hundred jobs
			Pump(cBuf, capacity)
		}
		done <- 1
	}()

	go chunker.Schedule()
	<-done
	for totalCountedItems < 1000 {
	}
	jobs := count
	doneTime := time.Now().Sub(t)
	chunker.Close()
	b.Logf("Created %d jobs of buffered channel capacity %d - with %d parrallel possible processes for each job and total items processed %d in %v sec", jobs, capacity, parrallelProcesses, totalCountedItems, doneTime.Seconds())

}

func BenchmarkChunk1000(b *testing.B) {
	t := time.Now()
	count = 0
	totalCountedItems = 0
	parrallelProcesses := 6
	capacity := 1000
	cBuf := make(cType, capacity)
	chunker := NewChunk(cBuf, Sump, parrallelProcesses)
	done := make(chan int)

	go func() {
		for totalCountedItems < 10000 {
			Pump(cBuf, capacity)
		}
		done <- 1
	}()

	go chunker.Schedule()
	<-done
	for totalCountedItems < 10000 {

	}
	jobs := count
	doneTime := time.Now().Sub(t)
	chunker.Close()
	b.Logf("Created %d jobs of buffered channel capacity %d - with %d parrallel possible processes for each job and total items processed %d in %v sec", jobs, capacity, parrallelProcesses, totalCountedItems, doneTime.Seconds())
}

func BenchmarkChunk10000(b *testing.B) {
	t := time.Now()
	count = 0
	totalCountedItems = 0
	parrallelProcesses := 6
	capacity := 10000
	cBuf := make(cType, capacity)
	chunker := NewChunk(cBuf, Sump, parrallelProcesses)
	done := make(chan int)

	go func() {
		for totalCountedItems < 100000 {
			Pump(cBuf, capacity)
		}
		done <- 1
	}()

	go chunker.Schedule()
	<-done
	for totalCountedItems < 100000 {

	}
	jobs := count
	doneTime := time.Now().Sub(t)
	chunker.Close()
	b.Logf("Created %d jobs of buffered channel capacity %d - with %d parrallel possible processes for each job and total items processed %d in %v sec", jobs, capacity, parrallelProcesses, totalCountedItems, doneTime.Seconds())
}
