# Notifier

**Notifier** is a Go library for message delivery, according to the requirements of the refurbed challenge.
In order to handle requests efficiently, it implements a worker pool pattern, where messages are received and stored in a buffered processing queue, while a group of parallel workers pick jobs and process them accordingly.

## Packages

The library implements the following core packages:

- **Rest**: This package implements a basic rest client that allows to send POST requests to the specified URL. It also allows a few other small functionalities, like add custom headers to the request (in case they are needed). 
- **Processor**: It implements the logic for the worker pool pattern for message processing, through the use of channels and go routines. It also allows a few options that can help customize the pool for better performance.
- **Notifier**: This is just a wrapper that abstracts all the previous logic to make it simple for the consumer to implement.

## Usage

Use of the library is fair simple, as shown in the example below:

```golang
// url represents the url where messages will be sent
// interval represents the time interval at which messages will be delivered
// options represents processing queue custom options that will be configured at startup (maxQueueSize and maxSenders)
// headers represent additional headers that need to be included in each call
notifier := notifier.NewNotifier(*url, interval, &options, *headers)

notifier.Start() // Starts the worker pool

notifier.Notify("this is a message") // Processes a message

notifier.Stop() // Gracefully stops the worker pool

```

## Testing
The source code includes a makefile to generate test results and coverage reports.

- For test results: `make test`
- For coverage report: `make cover`

## Design considerations
Given the initial requirements, the following describes the considerations taken into account when designing the solution:

- Use of buffered channel over an unbuffered one: the purpose of this decision is to keep control over the amount of memory used for processing, as well as to avoid locks during message reception, allowing the caller to keep enqueuing messages (even during peaks) without delays. For example, considering an average message size of 4Kb, enqueuing 100,000 messages would take approximately 400Mb of allocated memory. Currently, the default queue size is 1,000 messages, with 100 default workers.
- A con of using a buffered channel is that delivery is not guaranteed for 100% of cases. In specific cases when the queue is full and there are no workers available, a message could be dropped. To partially fix this issue, a basic retry mechanism is implemented, that allows to re-process a message that was originally dropped. Also, given that jobs are not resource intensive (a basic API call), CPU cores can maintain lots of workers spawned to keep up with processing demand.
- The REST client expects an HTTP 201 response code without a body for a successful request. Also, it's configured by default to timeout after 1s without a response.
- Intervals for delivery have been restricted to a duration up to 10s. Having intervals bigger than 10s (which is a lot) didn't seemed to add value to the solution.

## Future Improvements
No code is perfect, so of course there's room for improvement. A few points I think could be improved in the future are the following:

- **Retry mechanism**: Current retry mechanism is fairly basic as it only waits for a few secs before retrying a message, hoping that by that time the queue have some room available. An improvement could be to implement another channel as retry queue, but that would require to handle new considerations regarding buffer size, and new sync to know when to send messages to the main queue, etc.
- **Rest client**: Current rest client is super basic, so another improvement could be to add more options to it, like retry policy, connection timeouts, authorization, etc.
- **Target API failure mechanism**: Another improvement could be to implement some sort of circuit breaker that stops delivery when the target API does not respond or respond with errors for a certain amount of time.