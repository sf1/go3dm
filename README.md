# go3dm
Go packages for importing and converting 3D models. **Highly experimental code**.

## Status

Rudimentary support for loading and converting Wavefront OBJ files (including materials)

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
...
```
