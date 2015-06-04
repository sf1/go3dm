# go3dm
Go packages for importing and converting 3D models. **Highly experimental code**.

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

```
mesh, materials, err := go3dm.LoadOBJ("al.obj")
if err != nil { panic(err) }

// Print vertex and normal slices
fmt.Println(mesh.Vertices())
fmt.Println(mesh.Normals()}

// Iterate through named objects / polygon groups
for _, obj := range mesh.Objects() {
    fmt.Println(obj.Name())
    // ...
    // Get material, if any
    if obj.MaterialRef() != "" {
        mat := materials[obj.MaterialRef()]
        fmt.Println(mat)
    }
}
```
