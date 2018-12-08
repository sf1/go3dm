# go3dm
Go packages for importing and converting 3D models.

## Status

Rudimentary support for loading and converting Wavefront OBJ files (including materials). Only triangluar faces are supported.

## Installation

Package:

```
go get github.com/sf1/go3dm
```

OBJ to Go code converter:

```
go get github.com/sf1/go3dm/obj2go
```

## Usage

Load OBJ into an indexed triangle mesh:

```
mesh, materials, err := go3dm.LoadOBJ("al.obj", true)
if err != nil { panic(err) }

// Print vertex and normal slices
fmt.Println(mesh.Vertices)
fmt.Println(mesh.Normals}
fmt.Println(mesh.VertexIndex}

// Iterate through named objects / polygon groups
for _, obj := range mesh.Objects {
    fmt.Println(obj.Name)
    // ...
    // Get material, if any
    if obj.MaterialRef != "" {
        mat := materials[obj.MaterialRef]
        fmt.Println(mat)
    }
}
```
