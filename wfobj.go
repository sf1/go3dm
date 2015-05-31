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

func LoadOBJ(objPath string) (*Mesh, map[string]*Material, error) {
    objPath, err := filepath.Abs(objPath)
    if err != nil { return nil, nil, err}
    absDir := filepath.Dir(objPath)
    matMap := make(map[string]*Material)
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
            return nil, nil, fmt.Errorf("Can't open mtllib: %s", objMesh.MTLLib)
        }
        defer mtlFile.Close()
        matList, err := LoadMTLFrom(mtlFile)
        if err != nil { return nil, nil, err}
        for _, mat := range matList {
            matMap[mat.Name] = mat
        }
    }
    return &objMesh.Mesh, matMap, nil
}

func LoadOBJFrom(reader io.Reader) (*OBJMesh, error) {
    // Arrays for holding original data
    verticesTmp := NewF32VA(3)
    normalsTmp := NewF32VA(3)
    texTmp := NewF32VA(2)
    // Arrays for holding processed data
    vertices := NewF32VA(3)
    normals := NewF32VA(3)
    texCoords := NewF32VA(3)

    groups := make([]*MeshObject, 0, 1)
    groups = append(groups, &MeshObject{"unkown", -1, -1, "", false})
    mtllib := ""

    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        tokens := strings.Split(line, " ")
        switch tokens[0] {
        case "g","o":
            groups = append(groups, &MeshObject{tokens[1], -1, -1, "",false})
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
            face := tokens[1:]
            if len(face) != 3 {
                return nil, fmt.Errorf(
                    "Loader currently only supports triangular faces")
            }
            group := groups[len(groups)-1]
            if group.Offset == -1 {
                group.Offset = len(vertices.Values)
                group.Count = 0
            }
            for _, fidx := range face {
                vIdx, tIdx, nIdx, err := parseFaceIndicies(fidx)
                if err != nil {return nil, err}
                vertices.AppendVector(verticesTmp.GetVector(vIdx-1))
                if nIdx > 0 {
                    normals.AppendVector(normalsTmp.GetVector(nIdx-1))
                }
                if tIdx > 0 {
                    texCoords.AppendVector(texTmp.GetVector(tIdx-1))
                }
                group.Count += 3
            }
        case "s":
            group := groups[len(groups)-1]
            group.Smooth = false
            if tokens[1] == "1" {
                group.Smooth = true
            }
        case "mtllib":
            mtllib = strings.Join(tokens[1:]," ")
        case "usemtl":
            groups[len(groups)-1].MaterialRef = strings.Join(tokens[1:]," ")
        }
    }

    if groups[0].Offset == -1 {
        groups = groups[1:]
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

    return &OBJMesh{
            Mesh{verticesFA, normalsFA, texCoordsFA, groups},
            mtllib}, nil
}

func LoadMTLFrom(reader io.Reader) ([]*Material, error) {
    scanner := bufio.NewScanner(reader)
    materials := make([]*Material, 0, 1)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        tokens := strings.Split(line, " ")
        if tokens[0] == "newmtl" {
            m := new(Material)
            m.Name = strings.Join(tokens[1:]," ")
            materials = append(materials, m)
        } else if len(materials) > 0 {
            var err error
            m := materials[len(materials)-1]
            switch tokens[0] {
            case "Ka":
                m.Ka, err = parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
            case "Kd":
                m.Kd, err = parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
            case "Ks":
                m.Ks, err = parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
            case "Ns":
                val, err := parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
                m.Ns = val[0]
            case "d","Tr":
                val, err := parseF32Tokens(tokens[1:])
                if err != nil {return nil, err}
                m.Tr = val[0]
            case "map_Ka":
                m.KaMapName = strings.Join(tokens[1:]," ")
            case "map_Kd":
                m.KdMapName = strings.Join(tokens[1:]," ")
            case "map_Ks":
                m.KsMapName = strings.Join(tokens[1:]," ")

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
    Mesh
    MTLLib string
}

