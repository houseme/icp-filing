# ICP-filing

[![Go Reference](https://pkg.go.dev/badge/github.com/houseme/icp-filing.svg)](https://pkg.go.dev/github.com/houseme/icp-filing)
[![ICP-Filing-CI](https://github.com/houseme/icp-filing/actions/workflows/go.yml/badge.svg)](https://github.com/houseme/icp-filing/actions/workflows/go.yml)
![GitHub](https://img.shields.io/github/license/houseme/icp-filing?style=flat-square)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/houseme/icp-filing/main?style=flat-square)

Domain name information filing

## Installation

```bash
go get -u -v github.com/houseme/icp-filing@main 
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    filing "github.com/houseme/icp-filing"
)

func main() {
    ctx := context.Background()
    f := filing.New(ctx, filing.WithLogPath(os.TempDir()))
    resp, err := f.DomainFilling(ctx, &filing.QueryRequest{
        UnitName: "baidu.com",
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Println("resp:", resp)
}

```

## License
FeiE is primarily distributed under the terms of both the [Apache License (Version 2.0)](LICENSE)