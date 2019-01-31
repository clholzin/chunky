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




