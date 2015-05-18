# odTimeTracker libraries written in Go

Library is used in other [odTimeTracker](https://github.com/odTimeTracker) projects written using [Go](https://golang.org) language:

- [odtimetracker-go-cli](https://github.com/odTimeTracker/odtimetracker-go-cli)
- [odtimetracker-go-cgi](https://github.com/odTimeTracker/odtimetracker-go-cgi)

## Downloading & Building

Just write this:

	go get github.com/odTimeTracker/odtimetracker-go-lib
	go build
	go install

If you have correctly installed **Go** language this will download sources and build them immediatelly.

## Documentation

After `odtimetracker-go-lib` is successfully built documentation is available using [godoc](http://godoc.org/golang.org/x/tools/cmd/godoc). Just use this command:

	godoc -http=localhost:4040

And then navigate your browser to `http://localhost:4040/pkg/github.com/odTimeTracker/odtimetracker-go-lib/`.

