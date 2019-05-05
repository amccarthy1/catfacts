catfacts
========

A cat facts API client library/CLI for golang

## Installation

`catfacts` supports go modules. To install using go > 1.11, just run `go get
github.com/amccarthy1/catfacts`

## Usage
```go
package main

import (
    "fmt"
    "github.com/amccarthy1/catfacts"
)

func main() {
    client := catfacts.NewClient()
    fact, err := client.GetRandomFact()
    if err != nil {
        panic(err)
    }
    fmt.Println(fact.Fact)
}
```

You can always look through the
[godoc](https://godoc.org/github.com/amccarthy1/catfacts) for more detailed
documentation.

If you'd prefer an example of an application using the client, check out
[amccarthy1/catfacts-cli](https://github.com/amccarthy1/catfacts-cli).

