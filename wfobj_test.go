package go3dm

import (
    "fmt"
    "testing"
    "strings"
)

const squareMeshOBJ string = `
mtllib square.mtl
o square
v -1.000000 0.000000 1.000000
v 1.000000 0.000000 1.000000
v -1.000000 0.000000 -1.000000
v 1.000000 0.000000 -1.000000
vn 0.000000 1.000000 0.000000
usemtl square
s off
f 2//1 4//1 3//1
f 1//1 2//1 3//1
`

const square2MeshOBJ string = `
mtllib square2.mtl
o square
v -1.000000 0.000000 1.000000
v 1.000000 0.000000 1.000000
v -1.000000 0.000000 -1.000000
v 1.000000 0.000000 -1.000000
usemtl square
s off
f 2 4 3 
f 1 2 3
`

const squareMTL string = `
newmtl square
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.000000 0.000000 0.640000
Ks 0.500000 0.500000 0.500000
Ni 1.000000
d 1.000000
illum 2
`

var squareMeshVertices = []float32{
    1.000000, 0.000000, 1.000000,
    1.000000, 0.000000, -1.000000,
    -1.000000, 0.000000, -1.000000,
    -1.000000, 0.000000, 1.000000,
    1.000000, 0.000000, 1.000000,
    -1.000000, 0.000000, -1.000000,
}

var squareMeshNormals = []float32{
    0.000000, 1.000000, 0.000000,
    0.000000, 1.000000, 0.000000,
    0.000000, 1.000000, 0.000000,
    0.000000, 1.000000, 0.000000,
    0.000000, 1.000000, 0.000000,
    0.000000, 1.000000, 0.000000,
}

const cubesMeshOBJ string = `
mtllib cubes.mtl
o redCube
v 0.414876 -1.594997 -0.480466
v 0.414876 -0.387669 -0.480466
v 0.414876 -1.594997 -1.687793
v 0.414876 -0.387669 -1.687793
v 1.622203 -1.594997 -0.480466
v 1.622203 -0.387669 -0.480466
v 1.622203 -1.594997 -1.687793
v 1.622203 -0.387669 -1.687793
vn -0.577300 0.577300 -0.577300
vn -0.577300 -0.577300 -0.577300
vn -0.577300 -0.577300 0.577300
vn 0.577300 0.577300 -0.577300
vn 0.577300 -0.577300 -0.577300
vn 0.577300 0.577300 0.577300
vn 0.577300 -0.577300 0.577300
vn -0.577300 0.577300 0.577300
usemtl redCube
s 1
f 4//1 3//2 1//3
f 8//4 7//5 3//2
f 6//6 5//7 7//5
f 2//8 1//3 5//7
f 3//2 7//5 5//7
f 8//4 4//1 2//8
f 2//8 4//1 1//3
f 4//1 8//4 3//2
f 8//4 6//6 7//5
f 6//6 2//8 5//7
f 1//3 3//2 5//7
f 6//6 8//4 2//8
o blueCube
v -0.714148 -1.000000 -1.000000
v -0.714148 -1.000000 1.000000
v -2.714148 -1.000000 1.000000
v -2.714148 -1.000000 -1.000000
v -0.714147 1.000000 -0.999999
v -0.714148 1.000000 1.000001
v -2.714148 1.000000 1.000000
v -2.714148 1.000000 -1.000000
vn 0.000000 -1.000000 0.000000
vn 0.000000 1.000000 0.000000
vn 1.000000 0.000000 0.000000
vn 0.000000 0.000000 1.000000
vn -1.000000 0.000000 0.000000
vn 0.000000 0.000000 -1.000000
usemtl blueCube
s off
f 10//9 11//9 12//9
f 16//10 15//10 14//10
f 13//11 14//11 10//11
f 14//12 15//12 11//12
f 15//13 16//13 12//13
f 9//14 12//14 16//14
f 9//9 10//9 12//9
f 13//10 16//10 14//10
f 9//11 13//11 10//11
f 10//12 14//12 11//12
f 11//13 15//13 12//13
f 13//14 9//14 16//14
`

const cubesMTL string = `
newmtl blueCube
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.000000 0.000000 0.640000
Ks 0.500000 0.500000 0.500000
Ni 1.000000
d 1.000000
illum 2

newmtl redCube
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.640000 0.000000 0.000000
Ks 0.500000 0.500000 0.500000
Ni 1.000000
d 1.000000
illum 2
`

var cubesMeshVertices = []float32 {
    // red cube
    0.414876, -0.387669, -1.687793, // v04
    0.414876, -1.594997, -1.687793, // v03
    0.414876, -1.594997, -0.480466, // v01
    1.622203, -0.387669, -1.687793, // v08
    1.622203, -1.594997, -1.687793, // v07
    0.414876, -1.594997, -1.687793, // v03
    1.622203, -0.387669, -0.480466, // v06
    1.622203, -1.594997, -0.480466, // v05
    1.622203, -1.594997, -1.687793, // v07
    0.414876, -0.387669, -0.480466, // v02
    0.414876, -1.594997, -0.480466, // v01
    1.622203, -1.594997, -0.480466, // v05
    0.414876, -1.594997, -1.687793, // v03
    1.622203, -1.594997, -1.687793, // v07
    1.622203, -1.594997, -0.480466, // v05
    1.622203, -0.387669, -1.687793, // v08
    0.414876, -0.387669, -1.687793, // v04
    0.414876, -0.387669, -0.480466, // v02
    0.414876, -0.387669, -0.480466, // v02
    0.414876, -0.387669, -1.687793, // v04
    0.414876, -1.594997, -0.480466, // v01
    0.414876, -0.387669, -1.687793, // v04
    1.622203, -0.387669, -1.687793, // v08
    0.414876, -1.594997, -1.687793, // v03
    1.622203, -0.387669, -1.687793, // v08
    1.622203, -0.387669, -0.480466, // v06
    1.622203, -1.594997, -1.687793, // v07
    1.622203, -0.387669, -0.480466, // v06
    0.414876, -0.387669, -0.480466, // v02
    1.622203, -1.594997, -0.480466, // v05
    0.414876, -1.594997, -0.480466, // v01
    0.414876, -1.594997, -1.687793, // v03
    1.622203, -1.594997, -0.480466, // v05
    1.622203, -0.387669, -0.480466, // v06
    1.622203, -0.387669, -1.687793, // v08
    0.414876, -0.387669, -0.480466, // v02
    //blue cube
    -0.714148, -1.000000, 1.000000, // v10
    -2.714148, -1.000000, 1.000000, // v11
    -2.714148, -1.000000, -1.000000, // v12
    -2.714148, 1.000000, -1.000000, // v16
    -2.714148, 1.000000, 1.000000, // v15
    -0.714148, 1.000000, 1.000001, // v14
    -0.714147, 1.000000, -0.999999, // v13
    -0.714148, 1.000000, 1.000001, // v14
    -0.714148, -1.000000, 1.000000, // v10
    -0.714148, 1.000000, 1.000001, // v14
    -2.714148, 1.000000, 1.000000, // v15
    -2.714148, -1.000000, 1.000000, // v11
    -2.714148, 1.000000, 1.000000, // v15
    -2.714148, 1.000000, -1.000000, // v16
    -2.714148, -1.000000, -1.000000, // v12
    -0.714148, -1.000000, -1.000000, // v09
    -2.714148, -1.000000, -1.000000, // v12
    -2.714148, 1.000000, -1.000000, // v16
    -0.714148, -1.000000, -1.000000, // v09
    -0.714148, -1.000000, 1.000000, // v10
    -2.714148, -1.000000, -1.000000, // v12
    -0.714147, 1.000000, -0.999999, // v13
    -2.714148, 1.000000, -1.000000, // v16
    -0.714148, 1.000000, 1.000001, // v14
    -0.714148, -1.000000, -1.000000, // v09
    -0.714147, 1.000000, -0.999999, // v13
    -0.714148, -1.000000, 1.000000, // v10
    -0.714148, -1.000000, 1.000000, // v10
    -0.714148, 1.000000, 1.000001, // v14
    -2.714148, -1.000000, 1.000000, // v11
    -2.714148, -1.000000, 1.000000, // v11
    -2.714148, 1.000000, 1.000000, // v15
    -2.714148, -1.000000, -1.000000, // v12
    -0.714147, 1.000000, -0.999999, // v13
    -0.714148, -1.000000, -1.000000, // v09
    -2.714148, 1.000000, -1.000000, // v16
}

var cubesMeshNormals = []float32 {
    // red cube
    -0.577300, 0.577300, -0.577300, // n01
    -0.577300, -0.577300, -0.577300, // n02
    -0.577300, -0.577300, 0.577300, // n03
    0.577300, 0.577300, -0.577300, // n04
    0.577300, -0.577300, -0.577300, // n05
    -0.577300, -0.577300, -0.577300, // n02
    0.577300, 0.577300, 0.577300, // n06
    0.577300, -0.577300, 0.577300, // n07
    0.577300, -0.577300, -0.577300, // n05
    -0.577300, 0.577300, 0.577300, // n08
    -0.577300, -0.577300, 0.577300, // n03
    0.577300, -0.577300, 0.577300, // n07
    -0.577300, -0.577300, -0.577300, // n02
    0.577300, -0.577300, -0.577300, // n05
    0.577300, -0.577300, 0.577300, // n07
    0.577300, 0.577300, -0.577300, // n04
    -0.577300, 0.577300, -0.577300, // n01
    -0.577300, 0.577300, 0.577300, // n08
    -0.577300, 0.577300, 0.577300, // n08
    -0.577300, 0.577300, -0.577300, // n01
    -0.577300, -0.577300, 0.577300, // n03
    -0.577300, 0.577300, -0.577300, // n01
    0.577300, 0.577300, -0.577300, // n04
    -0.577300, -0.577300, -0.577300, // n02
    0.577300, 0.577300, -0.577300, // n04
    0.577300, 0.577300, 0.577300, // n06
    0.577300, -0.577300, -0.577300, // n05
    0.577300, 0.577300, 0.577300, // n06
    -0.577300, 0.577300, 0.577300, // n08
    0.577300, -0.577300, 0.577300, // n07
    -0.577300, -0.577300, 0.577300, // n03
    -0.577300, -0.577300, -0.577300, // n02
    0.577300, -0.577300, 0.577300, // n07
    0.577300, 0.577300, 0.577300, // n06
    0.577300, 0.577300, -0.577300, // n04
    -0.577300, 0.577300, 0.577300, // n08
    // blue cube
    0.000000, -1.000000, 0.000000, // n09
    0.000000, -1.000000, 0.000000, // n09
    0.000000, -1.000000, 0.000000, // n09
    0.000000, 1.000000, 0.000000, // n10
    0.000000, 1.000000, 0.000000, // n10
    0.000000, 1.000000, 0.000000, // n10
    1.000000, 0.000000, 0.000000, // n11
    1.000000, 0.000000, 0.000000, // n11
    1.000000, 0.000000, 0.000000, // n11
    0.000000, 0.000000, 1.000000, // n12
    0.000000, 0.000000, 1.000000, // n12
    0.000000, 0.000000, 1.000000, // n12
    -1.000000, 0.000000, 0.000000, // n13
    -1.000000, 0.000000, 0.000000, // n13
    -1.000000, 0.000000, 0.000000, // n13
    0.000000, 0.000000, -1.000000, // n14
    0.000000, 0.000000, -1.000000, // n14
    0.000000, 0.000000, -1.000000, // n14
    0.000000, -1.000000, 0.000000, // n09
    0.000000, -1.000000, 0.000000, // n09
    0.000000, -1.000000, 0.000000, // n09
    0.000000, 1.000000, 0.000000, // n10
    0.000000, 1.000000, 0.000000, // n10
    0.000000, 1.000000, 0.000000, // n10
    1.000000, 0.000000, 0.000000, // n11
    1.000000, 0.000000, 0.000000, // n11
    1.000000, 0.000000, 0.000000, // n11
    0.000000, 0.000000, 1.000000, // n12
    0.000000, 0.000000, 1.000000, // n12
    0.000000, 0.000000, 1.000000, // n12
    -1.000000, 0.000000, 0.000000, // n13
    -1.000000, 0.000000, 0.000000, // n13
    -1.000000, 0.000000, 0.000000, // n13
    0.000000, 0.000000, -1.000000, // n14
    0.000000, 0.000000, -1.000000, // n14
    0.000000, 0.000000, -1.000000, // n14
}

func TestChanges(t *testing.T) {
    t.Log("Testing recent changes")
    r := strings.NewReader(squareMeshOBJ)
    mesh, err := LoadOBJFrom(r)
    if err != nil {
        t.Error(err)
    }
    printMesh(mesh)
}

/*
func TestLoadSquareMesh(t *testing.T) {
    t.Log("Testing LoadOBJFrom with square mesh")
    testLoadOBJFrom(t, squareMeshOBJ,
                squareMeshVertices,
                nil,
                squareMeshNormals)
}

func TestLoadSquare2Mesh(t *testing.T) {
    t.Log("Testing LoadOBJFrom with square 2 mesh")
    testLoadOBJFrom(t, square2MeshOBJ,
                squareMeshVertices,
                nil,
                nil)
}

func TestLoadCubesMesh(t *testing.T) {
    t.Log("Testing LoadOBJFrom with cubes mesh")
    testLoadOBJFrom(t, cubesMeshOBJ,
                cubesMeshVertices,
                nil,
                cubesMeshNormals)
}

func TestLoadCubesMTL(t *testing.T) {
    t.Log("Testing LoadMTLFrom with cubes material")
    r := strings.NewReader(cubesMTL)
    materials, err := LoadMTLFrom(r)
    if err != nil { t.Error(err) }
    if len(materials) != 2 {
        t.Error("Expected 2 materials")
    }
    if materials[0].Name() != "blueCube" {
        t.Error("Expected blueCube material")
    }
    if materials[0].Ns() != 96.078431 {
        t.Error("Wrong Ns value")
    }
    if materials[0].Ka()[2] != 0.0{
        t.Error("Wrong Ka value")
    }
    if materials[0].Kd()[2] != 0.64{
        t.Error("Wrong Kd value")
    }
    if materials[0].Ks()[1] != 0.5{
        t.Error("Wrong Ks value")
    }

    if materials[1].Name() != "redCube" {
        t.Error("Expected redCube material")
    }
    if materials[1].Ns() != 96.078431 {
        t.Error("Wrong Ns value")
    }
    if materials[1].Ka()[2] != 0.0{
        t.Error("Wrong Ka value")
    }
    if materials[1].Kd()[0] != 0.64{
        t.Error("Wrong Kd value")
    }
    if materials[1].Ks()[2] != 0.5{
        t.Error("Wrong Ks value")
    }

}

func TestLoadOBJ(t *testing.T) {
    t.Log("Testing LoadOBJ with cubes mesh")
    mesh, materials, err := LoadOBJ("test-meshes/cubes.obj")
    if err != nil { t.Error(err) }
    fmt.Println(mesh)
    for _, mat := range materials {
        fmt.Println(mat)
    }
}

func testLoadOBJFrom(t *testing.T, meshStr string,
                  expectedVertices []float32,
                  expectedTexCoords []float32,
                  expectedNormals []float32) {
    var i int
    r := strings.NewReader(meshStr)
    mesh, err := LoadOBJFrom(r)
    if err != nil { t.Error(err) }
    vertices, texcoords, normals := mesh.VTN()
    if vertices == nil {
        t.Error("Didn't load any vertices")
    }
    if expectedVertices != nil {
        if len(vertices) != len(expectedVertices) {
            t.Error("Unexpected number of vertices")
        }
        for i = 0; i < len(vertices); i++ {
            if vertices[i] != expectedVertices[i] {
                t.Errorf("Unexpected vertex data at index %d", i)
            }
        }
    }
    if normals != nil {
        if len(vertices) != len(normals) {
            t.Error("Number of normals should equal number of vertices")
        }
        if expectedNormals != nil {
            if len(normals) != len(expectedNormals) {
                t.Error("Unexpected number of normals")
            }
            for i = 0; i < len(normals); i++ {
                if normals[i] != expectedNormals[i] {
                    t.Errorf("Unexpected normal data at index %d", i)
                }
            }
        }
    }
    if texcoords != nil {
        if len(vertices)/3 != len(texcoords)/2 {
            t.Error("Number of texture coords should equal number of vertices")
        }
        if expectedTexCoords != nil {
            if len(texcoords) != len(expectedTexCoords) {
                t.Error("Unexpected number of texture coordinates")
            }
            for i = 0; i < len(normals); i++ {
                if texcoords[i] != expectedTexCoords[i] {
                    t.Errorf("Unexpected textrue data at index %d", i)
                }
            }
        }
    }
    //printMesh(mesh)
}
*/

func printMesh(mesh Mesh) {
    vertices, texcoords, normals := mesh.VTN()
    fmt.Println("\nVertices:")
    for idx, v := range vertices {
        if idx % 3 == 0 { fmt.Println("") } else { fmt.Print(", ") }
        fmt.Printf("%f", v)
    }
    if texcoords != nil {
        fmt.Println("\n\nTexture Coordinates:\n", normals)
    }
    if normals != nil {
        fmt.Println("\n\nNormals:")
        for idx, n := range normals {
            if idx % 3 == 0 { fmt.Println("") } else { fmt.Print(", ") }
            fmt.Printf("%f", n)
        }
    }
    fmt.Println("\n\nElements")
    for idx, e := range mesh.Elements() {
        if idx % 3 == 0 { fmt.Println("") } else { fmt.Print(", ") }
        fmt.Printf("%d", e)
    }
    fmt.Println("\n\nObjects")
    for _, g := range mesh.Objects() {
        fmt.Printf("------------\n%s\n------------\n", g.Name)
        fmt.Printf("Material: %s\n", g.MaterialRef)
        fmt.Printf("Smooth: %t\n", g.Smooth)
        fmt.Printf("Offset: %d\n", g.Offset)
        fmt.Printf("Count: %d\n", g.Count)
    }
}
