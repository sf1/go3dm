package go3dm

import (
    "fmt"
    "testing"
    "strings"
)

func TestLoadSquareMeshIndexed(t *testing.T) {
    t.Log("Testing: Square Mesh (Indexed)")
    testLoadOBJFrom(t, squareMeshOBJ,
                    squareMeshIndexedVertices,
                    nil,
                    squareMeshIndexedNormals,
                    squareMeshIndicies,
                    squareMeshObjects,
                    true)
}

func TestLoadCubesMesh(t *testing.T) {
    t.Log("Testing: Cubes Mesh")
    testLoadOBJFrom(t, cubesMeshOBJ,
                    cubesMeshVertices,
                    nil,
                    cubesMeshNormals,
                    nil,
                    cubesMeshObjects,
                    false)
}

func TestLoadCubesMeshIndexed(t *testing.T) {
    t.Log("Testing: Cubes Mesh (Indexed)")
    testLoadOBJFrom(t, cubesMeshOBJ,
                    cubesMeshIndexedVertices,
                    nil,
                    cubesMeshIndexedNormals,
                    cubesMeshIndicies,
                    cubesMeshObjects,
                    true)
}

func testLoadOBJFrom(t *testing.T, meshStr string,
                    expectedVertices []float32,
                    expectedTexCoords []float32,
                    expectedNormals []float32,
                    expectedIndicies []uint32,
                    expectedObjects []*MeshObject,
                    index bool) {
    var i int
    r := strings.NewReader(meshStr)
    mesh, err := LoadOBJFrom(r, index)
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
                    t.Errorf("Unexpected texture data at index %d", i)
                }
            }
        }
    }
    if expectedIndicies != nil {
        if mesh.Indicies == nil {
            t.Error("Expected indicies")
        }
        if len(expectedIndicies) != len(mesh.Indicies) {
            t.Error("Unexpected number of indicies")
        }
        for i = 0; i < len(expectedIndicies); i++ {
            if expectedIndicies[i] != mesh.Indicies[i] {
                t.Errorf("Unexpected index at %d", i)
            }
        }
    }
    if expectedObjects != nil {
        if mesh.Objects == nil {
            t.Error("Expected mesh objects")
        }
        if len(expectedObjects) != len(mesh.Objects) {
            t.Error("Unexpected number of mesh objects")
        }
        for i = 0; i < len(expectedObjects); i++ {
            if !expectedObjects[i].Equals(mesh.Objects[i]) {
                t.Errorf("Mesh object %d doesn't match expectation", i)
            }
        }
    }
    //printMesh(&mesh.TriangleMesh)
}

func TestLoadCubesMTL(t *testing.T) {
    t.Log("Testing: Cubes Mesh Material")
    r := strings.NewReader(cubesMTL)
    materials, err := LoadMTLFrom(r)
    if err != nil { t.Error(err) }
    if len(materials) != 2 {
        t.Error("Expected 2 materials")
    }
    if materials[0].Name != "blueCube" {
        t.Error("Expected blueCube material")
    }
    if materials[0].Ns != 96.078431 {
        t.Error("Wrong Ns value")
    }
    if materials[0].Ka[2] != 0.0{
        t.Error("Wrong Ka value")
    }
    if materials[0].Kd[2] != 0.64{
        t.Error("Wrong Kd value")
    }
    if materials[0].Ks[1] != 0.5{
        t.Error("Wrong Ks value")
    }

    if materials[1].Name != "redCube" {
        t.Error("Expected redCube material")
    }
    if materials[1].Ns != 96.078431 {
        t.Error("Wrong Ns value")
    }
    if materials[1].Ka[2] != 0.0{
        t.Error("Wrong Ka value")
    }
    if materials[1].Kd[0] != 0.64{
        t.Error("Wrong Kd value")
    }
    if materials[1].Ks[2] != 0.5{
        t.Error("Wrong Ks value")
    }
}

func printMesh(mesh *TriangleMesh) {
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
    if mesh.Indicies != nil {
        fmt.Println("\n\nTriangle Indicies:")
        for idx, e := range mesh.Indicies {
            if idx % 3 == 0 { fmt.Println("") } else { fmt.Print(", ") }
            fmt.Printf("%d", e)
        }
    }
    fmt.Println("\n\nObjects:")
    for _, mo := range mesh.Objects {
        fmt.Printf("\n  %s\n", mo.Name)
        fmt.Printf("  - Material: %s\n", mo.MaterialRef)
        fmt.Printf("  - Smooth: %t\n", mo.Smooth)
        fmt.Printf("  - Offset: %d\n", mo.IndexOffset)
        fmt.Printf("  - Count: %d\n", mo.IndexCount)
    }
    fmt.Println("")
}
