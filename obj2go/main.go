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
type Mesh struct {
    Vertices []float32
    Normals []float32
    TextureCoords []float32
    Objects []*MeshObject
}

func (m *Mesh) VTN() ([]float32, []float32, []float32) {
    return m.Vertices, m.TextureCoords, m.Normals
}

type MeshObject struct {
    Name string
    Offset int
    Count int
    MaterialRef string
    Smooth bool
}

func (mo *MeshObject) VertexOffset() int32 {
    return int32(mo.Offset / 3)
}

func (mo *MeshObject) VertexCount() int32 {
    return int32(mo.Count / 3)
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
    var pkg, name string
    flag.StringVar(&pkg, "package", "model", "Target package")
    flag.StringVar(&name, "name", "", "Go-friendly model name (Should start with upper case letter)")
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
    if name == "" {
        flag.Usage()
        fmt.Fprintln(os.Stderr,
            "Please specify a go-friendly name for the model.\n")
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
    fmt.Fprintf(f, "var %sMesh *Mesh = &Mesh {\n", name)
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
    fmt.Fprintf(f, "\n    // Groups / Objects")
    fmt.Fprintf(f, "\n    []*MeshObject{")
    for _, obj := range mesh.Objects {
        fmt.Fprintf(f, "\n        &MeshObject{\"%s\", %d, %d, \"%s\", %t},",
            obj.Name, obj.Offset, obj.Count,
            obj.MaterialRef, obj.Smooth)
    }
    fmt.Fprintf(f, "\n    },")
    fmt.Fprintf(f, "\n}")
    // Process materials
    if len(materials) == 0 {
        return
    }
    fmt.Fprintf(f,
        "\n\nvar %sMaterials map[string]*Material = map[string]*Material {",
        name)
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
