# Flow

## A lightweight net/http middleware chainer written in Go.

This module aims to help developer to be able to chain middlewares in the easiest and simplest way possible.

## How to install this package

The easiest way to install this package is:

* Import to your project directly
  
  ```go
  import ("github.com/lutffmn/flow")
  ```

* Open terminal and navigate to your project directory and run

* ```shell
  go mod tidy
  ```

Or you can install it using 

* ```shell
  go get -u github.com/lutffmn/flow
  ```

## Usage

```go
package main

import (
    "fmt"
    "net/http"
    
    // import the module
    "github.com/lutffmn/flow"
)

// dummy http handler
func dummy(w http.ResponseWriter, r *http.Request){
    // your handler logic here
}

// dummy middlewares
func middleware1(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    next.ServeHTTP(w, r)
})    
}

func middleware2(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    next.ServeHTTP(w, r)
})    
}


func main() {
    r := http.NewServeMux()
    // registering your handler to Flow instance using Handle()
    // lists which middlewares you need your handler to flow through using Thru()
    http.Handle("GET /", flow.Handle(dummy).Thru(middleware1, middleware2))
    http.ListenAndServe(":9000", r)
}
```

# Find a bug?

If you found an issue or would like to submit an improvement to this project, please submit an issue using the issues tab above. If you would like to submit a PR with a fix, **please reference** the issue you created!
