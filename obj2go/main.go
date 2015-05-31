package main

import (
    "fmt"
    "flag"
    "os"
)

const usage string = `
Wavefront OBJ to Go Code Converter

Copyright © 2015 Sebastian Fleissner
Distributed under The MIT License. For details see:
https://github.com/sf1/go3dm/blob/master/LICENSE

Usage:
  %s [options] [model]

Options:
`

func main() {
    var pkg, dest string
    flag.StringVar(&pkg, "package", "main", "Target package")
    flag.StringVar(&dest, "o", ".", "Destination folder")
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, usage, os.Args[0])
        flag.PrintDefaults()
        fmt.Fprintln(os.Stderr, "")
    }
    flag.Parse()
    args := flag.Args()
    if len(args) == 0 {
        flag.Usage()
        return
    }
}
