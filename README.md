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