package go3dm

import (
    "fmt"
    "testing"
    "strings"
)

const squareMesh string = `
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

const cubesMesh string = `
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

func TestLoadObj1(t *testing.T) {
    t.Log("Testing LoadObj with square mesh")
    r := strings.NewReader(squareMesh)
    mesh, err := LoadObj(r)
    if err != nil { t.Error(err) }
    printMesh(mesh)
}

func TestLoadObj2(t *testing.T) {
    t.Log("Testing LoadObj with cubes mesh")
    r := strings.NewReader(cubesMesh)
    mesh, err := LoadObj(r)
    if err != nil { t.Error(err) }
    printMesh(mesh)
}

func printMesh(mesh TriangleMesh) {
    vertices, _, normals := mesh.VTN()
    for i := 0; i < len(vertices)/3; i++  {
        vIdx := i*3
        fmt.Println(vertices[vIdx:vIdx+3],"//", normals[vIdx:vIdx+3])
    }
}
