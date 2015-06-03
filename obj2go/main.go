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
type Mesh interface {
    Vertices() []float32
    Normals() []float32
    TextureCoords() []float32
    Objects() []MeshObject
    VTN() ([]float32, []float32, []float32)
}

type BasicMesh struct {
    vertices []float32
    normals []float32
    textureCoords []float32
    objects []MeshObject
}

func (m *BasicMesh) Vertices() []float32 {
    return m.vertices
}

func (m *BasicMesh) Normals() []float32 {
    return m.normals
}

func (m *BasicMesh) TextureCoords() []float32 {
    return m.textureCoords
}

func (m *BasicMesh) Objects() []MeshObject {
    return m.objects
}

func (m *BasicMesh) VTN() ([]float32, []float32, []float32) {
    return m.vertices, m.textureCoords, m.normals
}

type MeshObject interface {
    Name() string
    Offset() int
    Count() int
    MaterialRef() string
    Smooth() bool
    VertexOffset() int32
    VertexCount() int32
}

type BasicMeshObject struct {
    name string
    offset int
    count int
    materialRef string
    smooth bool
}

func (mo *BasicMeshObject) Name() string {
    return mo.name
}

func (mo *BasicMeshObject) Offset() int {
    return mo.offset
}

func (mo *BasicMeshObject) Count() int {
    return mo.count
}

func (mo *BasicMeshObject) MaterialRef() string {
    return mo.materialRef
}

func (mo *BasicMeshObject) Smooth() bool {
    return mo.smooth
}

func (mo *BasicMeshObject) VertexOffset() int32 {
    return int32(mo.offset / 3)
}

func (mo *BasicMeshObject) VertexCount() int32 {
    return int32(mo.count / 3)
}

type Material interface {
    Name() string
    Ka() []float32
    Kd() []float32
    Ks() []float32
    Ns() float32
    Tr() float32
    KaMapName() string
    KdMapName() string
    KsMapName() string
}

type BasicMaterial struct {
    name string
    ka []float32
    kd []float32
    ks []float32
    ns float32
    tr float32
    kaMapName string
    kdMapName string
    ksMapName string
}

func NewBasicMaterial(name string) Material {
    mat := new(BasicMaterial)
    mat.name = name
    return mat
}

func (m *BasicMaterial) Name() string {
    return m.name
}

func (m *BasicMaterial) Ka() []float32 {
    return m.ka
}

func (m *BasicMaterial) Kd() []float32 {
    return m.kd
}

func (m *BasicMaterial) Ks() []float32 {
    return m.ks
}

func (m *BasicMaterial) Ns() float32 {
    return m.ns
}

func (m *BasicMaterial) Tr() float32 {
    return m.tr
}

func (m *BasicMaterial) KaMapName() string {
    return m.kaMapName
}

func (m *BasicMaterial) KdMapName() string {
    return m.kdMapName
}

func (m *BasicMaterial) KsMapName() string {
    return m.ksMapName
}
`

func main() {
    var pkg, name string
    flag.StringVar(&pkg, "package", "model", "Target package")
    flag.StringVar(&name, "name", "", "Go-friendly model name")
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
    fmt.Fprintf(f, "var %sMesh Mesh = &BasicMesh {\n", name)
    fmt.Fprintf(f, "    // Vertices\n")
    fmt.Fprintf(f, "    []float32{")
    for idx, val := range mesh.Vertices() {
        if (idx % 6) == 0 {
            fmt.Fprintf(f, "\n       ")
        }
        fmt.Fprintf(f," %f,", val)
    }
    fmt.Fprintf(f, "\n    },")
    fmt.Fprintf(f, "\n    // Normals")
    if mesh.Normals != nil {
        fmt.Fprintf(f, "\n    []float32{")
        for idx, val := range mesh.Normals() {
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
        for idx, val := range mesh.TextureCoords() {
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
    fmt.Fprintf(f, "\n    []MeshObject{")
    for _, obj := range mesh.Objects() {
        fmt.Fprintf(f,
            "\n        &BasicMeshObject{\"%s\", %d, %d, \"%s\", %t},",
            obj.Name(), obj.Offset(), obj.Count(),
            obj.MaterialRef(), obj.Smooth())
    }
    fmt.Fprintf(f, "\n    },")
    fmt.Fprintf(f, "\n}")
    // Process materials
    if len(materials) == 0 {
        return
    }
    fmt.Fprintf(f,
        "\n\nvar %sMaterials map[string]Material = map[string]Material {",
        name)
    for key, mat := range materials {
        fmt.Fprintf(f,"\n    \"%s\": &BasicMaterial{\n        \"%s\",",
            key, mat.Name)
        ka := mat.Ka()
        fmt.Fprintf(f,"\n        []float32{%f, %f, %f},",ka[0],ka[1],ka[2])
        kd := mat.Kd()
        fmt.Fprintf(f,"\n        []float32{%f, %f, %f},",kd[0],kd[1],kd[2])
        ks := mat.Ks()
        fmt.Fprintf(f,"\n        []float32{%f, %f, %f},",ks[0],ks[1],ks[2])
        fmt.Fprintf(f,"\n        %f, %f,", mat.Ns(), mat.Tr())
        fmt.Fprintf(f,"\n        \"%s\", \"%s\", \"%s\",",
            mat.KaMapName(),
            mat.KdMapName(),
            mat.KsMapName())
        fmt.Fprintf(f,"\n    },")
    }
    fmt.Fprintf(f, "\n}")
}
