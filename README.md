# fofi
Fan-out Fan-in Go example

We have a time/resource expensive job to do. In fact, we need to do these expensive jobs several times.<br>
Later, we want to process the output of these jobs and operate with the results.

We can run this within a loop, sequentially, something like:
```
for job in jobs:
  output = run(job)
  compute(output)
```

This would work, but `golang` provides very useful features to resolve this in a more performant way.

This example describes a `fan-out` `fan-in` pattern:

![fan-out fan-in](/docs/img/fan-out-fan-in.png)

Using `concurrency`, `parallelism` and `go channels` we can achieve a performant implementation easily.

Our main flow looks like this:

![main flow](/docs/img/fofi-main-flow.png)

Basically, we are spawning many producer processes (fan out), these producers are writing the output in a go channel. We are also spawning a single consumer process which is reading output go channel and processing the data (fan in).

Running `fofi.go` with `concurrentGoRoutines = 1`, to all effects, would be like running sequential solution:
```
$ time go run fofi.go
starting fofi processing
producer processing 5
consumer processing 5
producer wrote 5
producer processing 3
consumer processed 5
producer wrote 3
producer processing 2
consumer processing 3
consumer processed 3
producer wrote 2
producer processing 4
consumer processing 2
consumer processed 2
producer wrote 4
producer processing 1
consumer processing 4
consumer processed 4
producer wrote 1
consumer processing 1
consumer processed 1
consumer done
fofi processing finished
$ go run fofi.go  0.26s user 0.15s system 3% cpu 10.858 total
```

Running `fofi.go` with `concurrentGoRoutines = 5`:
```
time go run fofi.go
starting fofi processing
producer processing 5
producer processing 1
producer processing 3
producer processing 4
producer processing 2
producer wrote 4
consumer processing 4
producer wrote 3
producer wrote 5
producer wrote 2
producer wrote 1
consumer processed 4
consumer processing 1
consumer processed 1
consumer processing 3
consumer processed 3
consumer processing 2
consumer processed 2
consumer processing 5
consumer processed 5
consumer done
fofi processing finished
go run fofi.go  0.26s user 0.15s system 8% cpu 4.835 total
```

Caveats:
- Is important to close Consumer channel once all Producers have finished the work, otherwise we might produce a panic error
- We are using a channel to limit the number of Producer go routines running concurrently
- Producer channel size limits the number of concurrent Producer processes running
- Consumer channel size defines the # of elements that Producers can put there before become blocked. This would generate `backpressure` to Producers and this might be intended and beneficial or unintended and harmful
