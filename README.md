### Chunky : parallel process your buffered channels


Chunky will schedule out parallel runs of the same process for a single buffered channel.



### Create a new chunker

```golang
bufChannel := make(chan int,1000)
```
Decide how many chunks to make
```golang
factor := 2

```
Add your process function to process each chunk
```golang

func Sump(b interface{}) {
	switch val := b.(type) {
	case chan int:
		length := len(val)
		for i := 0; i < length; i++ {
		   // process each value on the buffer
		   <-val
		}

	}
}


```
Create a new chunker
```golang

chunker := chunky.NewChunk(bufChannel,Sump,factor)

```
Call Schedule and as data is added to the single buffer new chucks will be created
```golang
chunker.Schedule()
```

It is also possible to Close the scheduler and run again with an updated factor
```golang

chunker.Close()
chucker.Factor(3)
chucker.Schedule()

```
Or simply call
```golang
chucker.Factor(3)
```



craigo: ls
README.md     assets        chunk.go      chunk_test.go cpu.out       example
craigo: cat chunk.go
package chunky

import (
	"errors"
	"reflect"
)

type Builder interface {
	Read() (interface{}, bool)
	Write(interface{})
	Len() int
	Schedule()
	Factor(int) int
	Close() error
}

type Chunk struct {
	bufChan reflect.Value
	tt      reflect.Type
	cap     int
	factor  int
	fn      func(interface{})
	close   chan int
}

func NewChunk(input interface{}, fn func(interface{}), factor int) Builder {
	chunky := &Chunk{}
	chunky.bufChan = reflect.ValueOf(input)
	chunky.fn = fn
	chunky.tt = chunky.bufChan.Type()
	chunky.cap = chunky.bufChan.Cap()
	chunky.close = make(chan int)
	chunky.Factor(factor)
	return chunky
}

func (chunk *Chunk) Factor(factor int) int {
	if factor > 0 {
		chunk.factor = factor
	} else {
		chunk.factor = 1
	}
	return chunk.factor
}

func (chunk *Chunk) Len() int {
	return chunk.bufChan.Len()
}

func (chunk *Chunk) Read() (interface{}, bool) {
	v, ok := chunk.bufChan.Recv()
	return v.Interface(), ok
}

func (chunk *Chunk) Write(data interface{}) {
	d := reflect.ValueOf(data)
	chunk.bufChan.Send(d)
}

func (chunk *Chunk) Close() error {
	select {
	case chunk.close <- 1:
		return nil
	default:
		return errors.New("Failed to write")
	}
}

func (chunk *Chunk) Schedule() {
	for {
		select {
		case <-chunk.close:
			return
		default:
		}
		var remainder int
		length := chunk.Len()
		if length == 0 {
			continue
		}
		factor := chunk.factor
		reduce := length / factor
		remainder = length - (reduce * factor)

		if length < (chunk.cap / factor) {
			reduce = length
			factor = 1
		}
		for i := 1; i <= factor; i++ {
			if i == factor && remainder > 0 {
				reduce = reduce + remainder
			}
			newVal := reflect.MakeChan(chunk.tt, reduce)
			for r := 0; r < reduce; r++ {
				if v, ok := chunk.bufChan.Recv(); ok {
					newVal.Send(v)
				}
			}

			go chunk.fn(newVal.Interface())
		}
	}

}

