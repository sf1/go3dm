package main

import (
    "fmt"
    "flag"
    "os"
    "github.com/sf1/go3dm"
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
    if len(args) != 1 {
        flag.Usage()
        return
    }
    if _, err := os.Stat(dest); os.IsNotExist(err) {
        fmt.Fprintf(os.Stderr, "Error: folder %s does not exist\n", dest)
        os.Exit(1)
    }
    _, _, err := go3dm.LoadOBJ(args[0])
    if err != nil { complainAndExit(err) }
    f, err := os.Create("out.go")
    if err != nil { complainAndExit(err) }
    defer f.Close()
    fmt.Fprintf(f, "package %s", pkg)
}

func complainAndExit(err error) {
    fmt.Fprintf(os.Stderr, "Error: %s\n", err)
    os.Exit(1)
}
