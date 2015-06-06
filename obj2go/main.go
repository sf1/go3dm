package main

import (
    "fmt"
    "flag"
    "os"
    "path/filepath"
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

const types string = `
type TriangleMesh struct {
    Vertices []float32
    Normals []float32
    TextureCoords []float32
    Indicies []uint32
    Objects []*MeshObject
}

func (m *TriangleMesh) VTN() ([]float32, []float32, []float32) {
    return m.Vertices, m.TextureCoords, m.Normals
}

type MeshObject struct {
    Name string
    IndexOffset int
    IndexCount int
    MaterialRef string
    Smooth bool
}

type Material struct {
    Name string
    Ka []float32
    Kd []float32
    Ks []float32
    Ns float32
    Tr float32
    KaMapName string
    KdMapName string
    KsMapName string
}
`

func main() {
    var pkg string
    flag.StringVar(&pkg, "package", "", "Package name")
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
    if pkg == "" {
        flag.Usage()
        fmt.Fprintln(os.Stderr,
            "Please specify a package name.\n")
        return
    }
    if _, err := os.Stat(pkg); os.IsNotExist(err) {
        os.MkdirAll(pkg, 0755)
    }
    outFile := filepath.Join(pkg, filepath.Base(args[0]) + ".go")
    mesh, materials, err := go3dm.LoadOBJ(args[0])
    if err != nil { panic(err) }
    // Write types.go
    f, err := os.Create(filepath.Join(pkg, "types.go"))
    if err != nil { panic(err) }
    defer f.Close()
    fmt.Fprintf(f, "package %s\n\n", pkg)
    fmt.Fprintf(f, types)
    // Create <modelname>.go
    f, err = os.Create(outFile)
    fmt.Fprintf(f, "package %s\n\n", pkg)
    // Process mesh data
    fmt.Fprintf(f, "var Mesh *TriangleMesh = &TriangleMesh {\n")
    fmt.Fprintf(f, "    // Vertices\n")
    fmt.Fprintf(f, "    []float32{")
    for idx, val := range mesh.Vertices {
        if (idx % 6) == 0 {
            fmt.Fprintf(f, "\n       ")
        }
        fmt.Fprintf(f," %f,", val)
    }
    fmt.Fprintf(f, "\n    },")
    fmt.Fprintf(f, "\n    // Normals")
    if mesh.Normals != nil {
        fmt.Fprintf(f, "\n    []float32{")
        for idx, val := range mesh.Normals {
            if (idx % 6) == 0 {
                fmt.Fprintf(f, "\n       ")
            }
            fmt.Fprintf(f," %f,", val)
        }
        fmt.Fprintf(f, "\n    },")
    } else {
        fmt.Fprintf(f, "\n    nil,")
    }
    fmt.Fprintf(f, "\n    // Texture Coordinates")
    if mesh.TextureCoords != nil {
        fmt.Fprintf(f, "\n    []float32{")
        for idx, val := range mesh.TextureCoords {
            if (idx % 6) == 0 {
                fmt.Fprintf(f, "\n       ")
            }
            fmt.Fprintf(f," %f,", val)
        }
        fmt.Fprintf(f, "\n    },")
    } else {
        fmt.Fprintf(f, "\n    nil,")
    }

    fmt.Fprintf(f, "\n    // Indicies")
    fmt.Fprintf(f, "\n    []uint32{")
    for idx, val := range mesh.Indicies {
        if (idx % 10) == 0 { fmt.Fprintf(f, "\n       ") }
        fmt.Fprintf(f," %d,", val)
    }
    fmt.Fprintf(f, "\n    },")

    fmt.Fprintf(f, "\n    // Groups / Objects")
    fmt.Fprintf(f, "\n    []*MeshObject{")
    for _, obj := range mesh.Objects {
        fmt.Fprintf(f,
            "\n        &MeshObject{\"%s\", %d, %d, \"%s\", %t},",
            obj.Name, obj.IndexOffset, obj.IndexCount,
            obj.MaterialRef, obj.Smooth)
    }
    fmt.Fprintf(f, "\n    },")
    fmt.Fprintf(f, "\n}")
    // Process materials
    if len(materials) == 0 {
        return
    }
    fmt.Fprintf(f,
        "\n\nvar Materials map[string]*Material = map[string]*Material {")
    for key, mat := range materials {
        fmt.Fprintf(f,"\n    \"%s\": &Material{\n        \"%s\",",
            key, mat.Name)
        fmt.Fprintf(f,"\n        []float32{%f, %f, %f},",
                    mat.Ka[0], mat.Ka[1], mat.Ka[2])
        fmt.Fprintf(f,"\n        []float32{%f, %f, %f},",
                    mat.Kd[0], mat.Kd[1], mat.Kd[2])
        fmt.Fprintf(f,"\n        []float32{%f, %f, %f},",
                    mat.Ks[0], mat.Ks[1], mat.Ks[2])
        fmt.Fprintf(f,"\n        %f, %f,", mat.Ns, mat.Tr)
        fmt.Fprintf(f,"\n        \"%s\", \"%s\", \"%s\",",
            mat.KaMapName,
            mat.KdMapName,
            mat.KsMapName)
        fmt.Fprintf(f,"\n    },")
    }
    fmt.Fprintf(f, "\n}")
}
