package go3dm

import (
    "io"
    "bufio"
    "strings"
    "strconv"
    "os"
    "fmt"
    "path/filepath"
)

func LoadOBJ(objPath string) (TriangleMesh, map[string]Material, error) {
    objPath, err := filepath.Abs(objPath)
    if err != nil { return nil, nil, err}
    absDir := filepath.Dir(objPath)
    matMap := make(map[string]Material)
    objFile, err := os.Open(objPath)
    if err != nil { return nil, nil, err}
    defer objFile.Close()
    objMesh, err := LoadOBJFrom(objFile)
    if err != nil { return nil, nil, err}
    if objMesh.MTLLib != "" {
        mtlPath := objMesh.MTLLib
        if !filepath.IsAbs(mtlPath) {
            mtlPath = filepath.Join(absDir, objMesh.MTLLib)
        }
        mtlFile, err := os.Open(mtlPath)
        if err != nil {
            return nil, nil, fmt.Errorf(
                "Can't open mtllib: %s", objMesh.MTLLib)
        }
        defer mtlFile.Close()
        matList, err := LoadMTLFrom(mtlFile)
        if err != nil { return nil, nil, err}
        for _, mat := range matList {
            matMap[mat.Name()] = mat
        }
    }
    return objMesh, matMap, nil
}

func LoadOBJFrom(reader io.Reader) (*OBJMesh, error) {
    // Arrays for holding original data
    verticesTmp := NewF32VA(3)
    normalsTmp := NewF32VA(3)
    texTmp := NewF32VA(2)
    vtnMap := make(map[string]uint32)
    // Arrays for holding processed data
    vertices := NewF32VA(3)
    normals := NewF32VA(3)
    texCoords := NewF32VA(3)
    elements := make([]uint32, 0, 10)

    meshObjects := make([]*BasicMeshObject, 0, 1)
    meshObjects = append(meshObjects,
        &BasicMeshObject{"unkown", -1, -1, "", false})
    mtllib := ""

    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        tokens := strings.Split(line, " ")
        switch tokens[0] {
        case "g","o":
            meshObjects = append(meshObjects,
                &BasicMeshObject{tokens[1], -1, -1, "",false})
        case "v":
            err := parseAndAppendF32Tokens(tokens[1:], verticesTmp)
            if err != nil {return nil, err}
        case "vn":
            err := parseAndAppendF32Tokens(tokens[1:], normalsTmp)
            if err != nil {return nil, err}
        case "vt":
            err := parseAndAppendF32Tokens(tokens[1:], texTmp)
            if err != nil {return nil, err}
        case "f":
            faceIndicies := tokens[1:]
            if len(faceIndicies) != 3 {
                return nil, fmt.Errorf(
                    "Loader currently only supports triangle faces")
            }
            mo := meshObjects[len(meshObjects)-1]
            if mo.indexOffset == -1 {
                mo.indexOffset = len(elements)
                mo.indexCount = 0
            }
            for _, fidx := range faceIndicies {
                mo.indexCount++
                vtnIdx, ok := vtnMap[fidx]
                if ok {
                    elements = append(elements, vtnIdx)
                    continue
                }
                vtnIdx = uint32(vertices.VectorCount())
                vIdx, tIdx, nIdx, err := parseFaceIndicies(fidx)
                if err != nil {return nil, err}
                vertices.AppendVector(verticesTmp.GetVector(vIdx-1))
                if nIdx > 0 {
                    normals.AppendVector(normalsTmp.GetVector(nIdx-1))
                }
                if tIdx > 0 {
                    texCoords.AppendVector(texTmp.GetVector(tIdx-1))
                }
                vtnMap[fidx] = vtnIdx
                elements = append(elements, vtnIdx)
            }
        case "s":
            mo := meshObjects[len(meshObjects)-1]
            mo.smooth = false
            if tokens[1] == "1" {
                mo.smooth = true
            }
        case "mtllib":
            mtllib = strings.Join(tokens[1:]," ")
        case "usemtl":
            meshObjects[len(meshObjects)-1].materialRef = strings.Join(
                tokens[1:]," ")
        }
    }

    if meshObjects[0].indexOffset == -1 {
        meshObjects = meshObjects[1:]
    }

    var verticesFA, normalsFA, texCoordsFA []float32 = nil, nil, nil

    if len(vertices.Values) > 0 {
        verticesFA = vertices.Values
    }
    if len(normals.Values) > 0 {
        normalsFA = normals.Values
    }
    if len(texCoords.Values) > 0 {
        texCoordsFA = texCoords.Values
    }

    objs := make([]MeshObject, 0, len(meshObjects))
    for _, mo := range meshObjects {
        objs = append(objs, mo)
    }

    return &OBJMesh{
            BasicTriangleMesh{
                verticesFA,
                normalsFA,
                texCoordsFA,
                elements,
                objs},
            mtllib}, nil
}

func LoadMTLFrom(reader io.Reader) ([]Material, error) {
    scanner := bufio.NewScanner(reader)
    materials := make([]Material, 0, 1)
    var curMat *BasicMaterial = nil
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        tokens := strings.Split(line, " ")
        if tokens[0] == "newmtl" {
            curMat = new(BasicMaterial)
            curMat.name = strings.Join(tokens[1:]," ")
            materials = append(materials, curMat)
        } else if len(materials) > 0 {
            var err error
            switch tokens[0] {
            case "Ka":
                curMat.ka, err = parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
            case "Kd":
                curMat.kd, err = parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
            case "Ks":
                curMat.ks, err = parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
            case "Ns":
                val, err := parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
                curMat.ns = val[0]
            case "d","Tr":
                val, err := parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
                curMat.tr = val[0]
            case "map_Ka":
                curMat.kaMapName = strings.Join(tokens[1:]," ")
            case "map_Kd":
                curMat.kdMapName = strings.Join(tokens[1:]," ")
            case "map_Ks":
                curMat.ksMapName = strings.Join(tokens[1:]," ")

            }
        }
    }
    return materials, nil
}

func parseAndAppendF32Tokens(tokens []string, floats *f32VA) error {
    for _,t := range tokens {
        v, err := strconv.ParseFloat(t, 32)
        if err != nil { return err }
        floats.Append(float32(v))
    }
    return nil
}

func parseF32Tokens(tokens []string) ([]float32, error) {
    result := make([]float32, 0, 1)
    for _,t := range tokens {
        v, err := strconv.ParseFloat(t, 32)
        if err != nil { return nil, err }
        result = append(result, float32(v))
    }
    return result, nil
}

func parseFaceIndicies(fidx string) (int, int, int, error) {
   var vIdx, tIdx, nIdx int = 0,0,0
    parts := strings.Split(fidx,"/")
    if len(parts[0]) < 1 { return 0,0,0,fmt.Errorf("Parse error: %s", fidx) }
    val, err := strconv.ParseUint(parts[0], 10, 32)
    if err != nil {return 0,0,0,err}
    vIdx = int(val)
    if len(parts) == 1 { return vIdx,0,0,nil }
    if len(parts) != 3 { return 0,0,0,fmt.Errorf("Parse error: %s", fidx) }
    if len(parts[1]) < 1 {
        tIdx = 0
    }
    if len(parts[2]) < 1 { return 0,0,0,fmt.Errorf("Parse error: %s", fidx) }
    val, err = strconv.ParseUint(parts[2], 10, 32)
    if err != nil {return 0,0,0,err}
    nIdx = int(val)
    return vIdx, tIdx, nIdx, nil
}

type OBJMesh struct {
    BasicTriangleMesh
    MTLLib string
}

