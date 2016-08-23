# go-money

[![godoc reference](https://cdn.rawgit.com/mvader/2faf5060e6cb109617ef5548836532aa/raw/2f5e2f2e934f6dde4ec4652ff0ae6d5c83cbfd6a/godoc.svg)](https://godoc.org/github.com/mvader/go-money) [![Build Status](https://travis-ci.org/mvader/go-money.svg?branch=master)](https://travis-ci.org/mvader/go-money) [![codecov](https://codecov.io/gh/mvader/go-money/branch/master/graph/badge.svg)](https://codecov.io/gh/mvader/go-money)  [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

Small Golang library to parse and print money amounts.

### Install

```
go get gopkg.in/mvader/go-money.v1
```

### Usage

```go
a, err := money.Parse("3500000 eur")
if err != nil {
  log.Fatal(err)
}

fmt.Println(a) // prints 3.5M €

fmt.Println(money.NewAmount(3500, "€")) // prints 3.5K €
```
