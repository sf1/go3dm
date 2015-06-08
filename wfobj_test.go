package go3dm

import (
    "fmt"
    "testing"
    "strings"
)

func TestLoadTexPlane(t *testing.T) {
    t.Log("Testing: Texplane Mesh")
    r := strings.NewReader(texplaneOBJ)
    mesh, err := LoadOBJFrom(r, false)
    if err != nil { t.Error(err); return }
    checkMesh(t, &mesh.TriangleMesh,
                texplaneVertices,
                texplaneTexCoords,
                texplaneNormals,
                nil,
                texplaneObjects)
}

func TestLoadTexPlaneAll(t *testing.T) {
    t.Log("Testing: Texplane Mesh V2, Materials, Textures")
    mesh, materials, err := LoadOBJ("test-meshes/texplane2.obj", false)
    if err != nil { t.Error(err); return }
    checkMesh(t, mesh,
                texplaneVertices,
                texplaneTexCoords,
                texplaneNormals,
                nil,
                texplaneObjects)
    checkMaterials(t, materials, texplaneV2Materials)
}

func TestLoadSquareIndexed(t *testing.T) {
    t.Log("Testing: Square Mesh (Indexed)")
    r := strings.NewReader(squareOBJ)
    mesh, err := LoadOBJFrom(r, true)
    if err != nil { t.Error(err); return }
    checkMesh(t, &mesh.TriangleMesh,
                squareIndexedVertices,
                nil,
                squareIndexedNormals,
                squareVertexIndex,
                squareObjects)
}

func TestLoadCubes(t *testing.T) {
    t.Log("Testing: Cubes Mesh")
    r := strings.NewReader(cubesOBJ)
    mesh, err := LoadOBJFrom(r, false)
    if err != nil { t.Error(err); return }
    checkMesh(t, &mesh.TriangleMesh,
                cubesVertices,
                nil,
                cubesNormals,
                nil,
                cubesObjects)
}

func TestLoadCubesIndexed(t *testing.T) {
    t.Log("Testing: Cubes Mesh (Indexed)")
    r := strings.NewReader(cubesOBJ)
    mesh, err := LoadOBJFrom(r, true)
    if err != nil { t.Error(err); return }
    checkMesh(t, &mesh.TriangleMesh,
                cubesIndexedVertices,
                nil,
                cubesIndexedNormals,
                cubesVertexIndex,
                cubesObjects)
}

func TestLoadCubesMTL(t *testing.T) {
    t.Log("Testing: Cubes Mesh Material")
    r := strings.NewReader(cubesMTL)
    materials, err := LoadMTLFrom(r)
    if err != nil { t.Error(err); return }
    matMap := make(map[string]*Material)
    for _, m := range materials { matMap[m.Name] = m }
    checkMaterials(t, matMap, cubesMaterials)
}

func checkMesh(t *testing.T, mesh *TriangleMesh,
                expectedVertices []float32,
                expectedTexCoords []float32,
                expectedNormals []float32,
                expectedIndicies []uint32,
                expectedObjects []*MeshObject) {
    var i int
    vertices, texcoords, normals := mesh.VTN()
    if vertices == nil {
        t.Error("Didn't load any vertices")
        return
    }
    if expectedVertices != nil {
        if len(vertices) != len(expectedVertices) {
            t.Error("Unexpected number of vertices")
            return
        }
        for i = 0; i < len(vertices); i++ {
            if vertices[i] != expectedVertices[i] {
                t.Errorf("Unexpected vertex data at index %d", i)
                return
            }
        }
    }
    if normals != nil {
        if len(vertices) != len(normals) {
            t.Error("Number of normals should equal number of vertices")
            return
        }
        if expectedNormals != nil {
            if len(normals) != len(expectedNormals) {
                t.Error("Unexpected number of normals")
                return
            }
            for i = 0; i < len(normals); i++ {
                if normals[i] != expectedNormals[i] {
                    t.Errorf("Unexpected normal data at index %d", i)
                    return
                }
            }
        }
    }
    if texcoords != nil {
        if len(vertices)/3 != len(texcoords)/2 {
            t.Error("Number of texture coords should equal number of vertices")
            return
        }
        if expectedTexCoords != nil {
            if len(texcoords) != len(expectedTexCoords) {
                t.Error("Unexpected number of texture coordinates")
                return
            }
            for i = 0; i < len(texcoords); i++ {
                if texcoords[i] != expectedTexCoords[i] {
                    t.Errorf("Unexpected texture data at index %d", i)
                    return
                }
            }
        }
    }
    if expectedIndicies != nil {
        if mesh.VertexIndex == nil {
            t.Error("Expected indicies")
            return
        }
        if len(expectedIndicies) != len(mesh.VertexIndex) {
            t.Error("Unexpected number of indicies")
            return
        }
        for i = 0; i < len(expectedIndicies); i++ {
            if expectedIndicies[i] != mesh.VertexIndex[i] {
                t.Errorf("Unexpected index at %d", i)
                return
            }
        }
    }
    if expectedObjects != nil {
        if mesh.Objects == nil {
            t.Error("Expected mesh objects")
            return
        }
        if len(expectedObjects) != len(mesh.Objects) {
            t.Error("Unexpected number of mesh objects")
            return
        }
        for i = 0; i < len(expectedObjects); i++ {
            if !expectedObjects[i].Equals(mesh.Objects[i]) {
                t.Errorf("Mesh object %d doesn't match expectation", i)
                return
            }
        }
    }
}

func checkMaterials(t *testing.T,
                    materials map[string]*Material,
                    expectedMaterials []*Material) {
    if len(materials) != len(expectedMaterials) {
        t.Errorf("Unexpected number of materials %d", len(materials))
        return
    }
    for _, em := range expectedMaterials {
        m, ok := materials[em.Name]
        if !ok {
            t.Errorf("Expected material not found: %s", em.Name)
            return
        }
        if !em.Equals(m) {
            t.Errorf("Material differs from expected material:\n%v\n%v", m, em)
            return
        }
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
        fmt.Println("\n\nTexture Coordinates:")
        for idx, t := range texcoords {
            if idx % 2 == 0 { fmt.Println("") } else { fmt.Print(", ") }
            fmt.Printf("%f", t)
        }
    }
    if normals != nil {
        fmt.Println("\n\nNormals:")
        for idx, n := range normals {
            if idx % 3 == 0 { fmt.Println("") } else { fmt.Print(", ") }
            fmt.Printf("%f", n)
        }
    }
    if mesh.VertexIndex != nil {
        fmt.Println("\n\nTriangle Indicies:")
        for idx, e := range mesh.VertexIndex {
            if idx % 3 == 0 { fmt.Println("") } else { fmt.Print(", ") }
            fmt.Printf("%d", e)
        }
    }
    fmt.Println("\n\nObjects:")
    for _, mo := range mesh.Objects {
        fmt.Printf("\n  %s\n", mo.Name)
        fmt.Printf("  - Material: %s\n", mo.MaterialRef)
        fmt.Printf("  - Smooth: %t\n", mo.Smooth)
        fmt.Printf("  - Offset: %d\n", mo.VertexOffset)
        fmt.Printf("  - Count: %d\n", mo.VertexCount)
    }
    fmt.Println("")
}
