# cn_go_patterns
This repository contains cloud native patterns written in Go


## Stability Patterns
### Circuit breaker
Circuit Breaker automatically degrades service functions in response to a likely fault,
preventing larger or cascading failures by eliminating recurring errors and providing reasonable error responses.

#### Participants
This pattern includes the following participants:

- Circuit
  The function that interacts with the service.

- Breaker
  A closure with the same function signature as Circuit.


### Debounce
Debounce limits the frequency of a function invocation so that only the first or last in a cluster of calls is actually performed.

#### Participants
This pattern includes the following participants:

- Circuit
  The function to regulate.

- Debounce
  A closure with the same function signature as Circuit.


### Retry
Retry accounts for a possible transient fault in a distributed system by transparently retrying a failed operation.

#### Participants
This pattern includes the following participants:

- Effector
  The function that interacts with the service.

- Retry
  A function that accepts Effector and returns a closure with the same function signature as effector.


### Throttle
Throttle limits the frequency of a function call to some maximum number of invocations per unit of time.

#### Participants
This pattern includes the following participants:

- Effector
  The function to regulate

- Throttle
  A function that accepts Effector and returns a closure with the same function signature as Effector.


### Timout
Timeout allows a process to stop waiting for an answer once it's clear that an answer may not be coming.

#### Participants
This pattern includes the following participants:

- Client
  The client who wants to execute SlowFunction.

- SlowFunction
  The long-running function that implements the functionality desired by Client.

- Timeout
  A wrapper function around SlowFunction that implements the timeout logic.

## Concurrency Patterns
### Fan-In
Fan-in multiplexes multiple input channels onto one output channel. Services that have  some number of workers that all generate output may find
it useful to combine all of the workers outputs to be prociessed as a single unified stream.

#### Participants
This pattern includes the following participants:

- Sources
  A set of one or more input channels with the same type. Accepted by Funnel

- Destination
  An output channel of the same type as Sources. Created and provided by Funnel.

- Funnel
  Accepts Sources and immediately returns Destination. Any input from any Sources will be output by Destination.


### Fan-Out
Fan-out evenly distributes messages from an input channel to multiple output channels.
Fan-out receives messages from an input channel, distributing them evenly among output channels, and is a useful pattern for parallelizing CPU and I/O utilization.

#### Participants
This pattern includes the following participants:

- Source
  An input channel. Accepted by Split.

- Destinations
  An output channel of the same type as Source. Created and provided by Split.

- Split
  A function that accepts Source and immediately returns Destinations. Any input from Source will be output to a Destination.


## Other notes

### Throttle and the Token Bucket algorithm
The Token Bucket algorithm uses the analogy of a bucket that can hold some maximum number of tokens. When a function is called, a token is taken 
from the bucket, which then refills at some fixed rate.

Common strategies for handling requests using Throttle pattern:
- Return an error
    - Useful to restrict unreasonable or potentially abusive numbers of client requests. A RESTful service example might return a status code of 429 (Too Many Requests).
  
- Replay the response of the last successful function call
    - Useful when a service or expensive fucntion call is likely to provide an identical result if called to soon.

- Enqueue the request for execution when sufficient tokens are available
    - Useful when you eventualy want to handle all requests. More complex and may require care to be taken to ensure that memory isn't exhausted.

### Difference between Throttle and Debounce
Throttle limits the event rate; Debounce allows only one event in a cluster.
example for 20 input requests within a given time frame say once per second, throttle might limit the event rate to 1 request every other second
effectively handling 10 out of 20 requests. Debounce on the other hand will handle the first out of 20 requests if debounce rate is set to 20sec and might cache the response to be returned immediately.

### Multiplexer (mux)
In electronics a mutliplexer or mux, also known as a data selector, is a device that selects between multiple digital or analog input signals
and forwards the selected input to a single output line. In CS the same aspects can be applied to provide a single channel for multiple user or process requests.

### Variadic function
A variadic function is a function of indefinite arity i.e., one which accepts a variable number of arguments.

### Fan-Out and 2 forwarding techniques
- Using a single goroutine that reads the values from Source and forwards them to the Destinations in a round-robin fashion. 
This has the virtue of requiring only one master goroutine, but if the next channel isn't ready to read yet, it'll slow the entire process.

- Using separate goroutines for each Destination that competes to read the next value from Source and forward it to their respective Destination.
This requires slightly more resources, but is less likely to get bogged down by a single slow-running worker.

