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

func LoadOBJ(objPath string, index bool) (*TriangleMesh,
    map[string]*Material, error) {
    objPath, err := filepath.Abs(objPath)
    if err != nil { return nil, nil, err}
    absDir := filepath.Dir(objPath)
    matMap := make(map[string]*Material)
    objFile, err := os.Open(objPath)
    if err != nil { return nil, nil, err}
    defer objFile.Close()
    objMesh, err := LoadOBJFrom(objFile, index)
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
            matMap[mat.Name] = mat
        }
    }
    return &objMesh.TriangleMesh, matMap, nil
}

type OLState struct {
    verticesTmp *f32VA
    normalsTmp *f32VA
    texTmp *f32VA
    vtnMap map[string]uint32
    vertices *f32VA
    normals *f32VA
    texCoords *f32VA
    indicies []uint32
    meshObjects []*MeshObject
    index bool
}

func LoadOBJFrom(reader io.Reader, index bool) (*OBJMesh, error) {
    // Set up state struct
    state := &OLState {
        NewF32VA(3), NewF32VA(3), NewF32VA(2),
        make(map[string]uint32),
        NewF32VA(3), NewF32VA(3), NewF32VA(2),
        make([]uint32, 0, 10),
        make([]*MeshObject, 0, 1),
        index,
    }

    state.meshObjects = append(state.meshObjects,
        &MeshObject{"unkown", -1, -1, "", false})
    mtllib := ""

    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        tokens := strings.Split(line, " ")
        switch tokens[0] {
        case "g","o":
            state.meshObjects = append(state.meshObjects,
                &MeshObject{tokens[1], -1, -1, "",false})
        case "v":
            err := parseAndAppendF32Tokens(tokens[1:], state.verticesTmp)
            if err != nil {return nil, err}
        case "vn":
            err := parseAndAppendF32Tokens(tokens[1:], state.normalsTmp)
            if err != nil {return nil, err}
        case "vt":
            err := parseAndAppendF32Tokens(tokens[1:], state.texTmp)
            if err != nil {return nil, err}
        case "f":
            faceIndicies := tokens[1:]
            if len(faceIndicies) != 3 {
                return nil, fmt.Errorf(
                    "Loader currently only supports triangle faces")
            }
            err := processFace(faceIndicies, state)
            if err != nil { return nil, err }
        case "s":
            mo := state.meshObjects[len(state.meshObjects)-1]
            mo.Smooth = false
            if tokens[1] == "1" {
                mo.Smooth = true
            }
        case "mtllib":
            mtllib = strings.Join(tokens[1:]," ")
        case "usemtl":
            state.meshObjects[len(state.meshObjects)-1].MaterialRef =
                strings.Join(tokens[1:]," ")
        }
    }

    if state.meshObjects[0].VertexOffset == -1 {
        state.meshObjects = state.meshObjects[1:]
    }

    var verticesFA, normalsFA, texCoordsFA []float32 = nil, nil, nil

    if len(state.vertices.Values) > 0 {
        verticesFA = state.vertices.Values
    }
    if len(state.normals.Values) > 0 {
        normalsFA = state.normals.Values
    }
    if len(state.texCoords.Values) > 0 {
        texCoordsFA = state.texCoords.Values
    }
    if len(state.indicies) == 0 {
        state.indicies = nil
    }

    return &OBJMesh{
            TriangleMesh{
                verticesFA,
                normalsFA,
                texCoordsFA,
                state.indicies,
                state.meshObjects},
            mtllib}, nil
}

func processFace(faceIndicies []string, state *OLState) error {
    mo := state.meshObjects[len(state.meshObjects)-1]
    if mo.VertexOffset == -1 {
        if state.index {
            mo.VertexOffset = int32(len(state.indicies))
        } else {
            mo.VertexOffset = int32(len(state.vertices.Values) / 3)
        }
        mo.VertexCount = 0
    }
    for _, fidx := range faceIndicies {
        vtnIdx, ok := state.vtnMap[fidx]
        if state.index {
            if ok {
                state.indicies = append(state.indicies, vtnIdx)
                mo.VertexCount++
                continue
            }
            vtnIdx = uint32(state.vertices.VectorCount())
        }
        vIdx, tIdx, nIdx, err := parseFaceIndicies(fidx)
        if err != nil {return err}
        state.vertices.AppendVector(
            state.verticesTmp.GetVector(vIdx-1))
        if nIdx > 0 {
            state.normals.AppendVector(
                state.normalsTmp.GetVector(nIdx-1))
        }
        if tIdx > 0 {
            state.texCoords.AppendVector(
                state.texTmp.GetVector(tIdx-1))
        }
        if state.index {
            state.vtnMap[fidx] = vtnIdx
            state.indicies = append(state.indicies, vtnIdx)
        }
        mo.VertexCount++
    }
    return nil
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
    if parts[1] != "" {
        val, err = strconv.ParseUint(parts[1], 10, 32)
        if err != nil {return 0,0,0,err}
        tIdx = int(val)
    }
    return vIdx, tIdx, nIdx, nil
}

func LoadMTLFrom(reader io.Reader) ([]*Material, error) {
    scanner := bufio.NewScanner(reader)
    materials := make([]*Material, 0, 1)
    var curMat *Material= nil
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        tokens := strings.Split(line, " ")
        if tokens[0] == "newmtl" {
            curMat = new(Material)
            curMat.Name = strings.Join(tokens[1:]," ")
            materials = append(materials, curMat)
        } else if len(materials) > 0 {
            var err error
            switch tokens[0] {
            case "Ka":
                curMat.Ka, err = parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
            case "Kd":
                curMat.Kd, err = parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
            case "Ks":
                curMat.Ks, err = parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
            case "Ns":
                val, err := parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
                curMat.Ns = val[0]
            case "d","Tr":
                val, err := parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
                curMat.Tr = val[0]
            case "map_Ka":
                curMat.KaMapName = strings.Join(tokens[1:]," ")
            case "map_Kd":
                curMat.KdMapName = strings.Join(tokens[1:]," ")
            case "map_Ks":
                curMat.KsMapName = strings.Join(tokens[1:]," ")
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

type OBJMesh struct {
    TriangleMesh
    MTLLib string
}

