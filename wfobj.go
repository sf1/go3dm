package go3dm

import (
    "io"
    "bufio"
    "strings"
    "strconv"
    "fmt"
)

func LoadOBJFrom(reader io.Reader) (*OBJMesh, error) {
    // Arrays for holding original data
    verticesTmp := NewF32VA(3)
    normalsTmp := NewF32VA(3)
    texTmp := NewF32VA(2)
    // Arrays for holding processed data
    vertices := NewF32VA(3)
    normals := NewF32VA(3)
    texCoords := NewF32VA(3)

    groups := make([]*OBJGroup, 0, 1)
    groups = append(groups, newOBJGroup("unkown", -1, -1, false, ""))
    mtllib := ""

    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        tokens := strings.Split(line, " ")
        switch tokens[0] {
        case "g","o":
            groups = append(groups,
                newOBJGroup(tokens[1], -1, -1, false, ""))
        case "v":
            err := parseF32Tokens(tokens[1:], verticesTmp)
            if err != nil {return nil, err}
        case "vn":
            err := parseF32Tokens(tokens[1:], normalsTmp)
            if err != nil {return nil, err}
        case "vt":
            err := parseF32Tokens(tokens[1:], texTmp)
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

    return newOBJMesh(
        verticesFA, normalsFA, texCoordsFA,
        groups, mtllib), nil
}

func parseF32Tokens(tokens []string, floats *f32VA) error {
    for _,t := range tokens {
        v, err := strconv.ParseFloat(t, 32)
        if err != nil { return err }
        floats.Append(float32(v))
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
    return vIdx, tIdx, nIdx, nil
}

type OBJMesh struct {
    Mesh
    Groups []*OBJGroup
    MTLLib string
}

func newOBJMesh(verts, norms, texcoords []float32,
                groups []*OBJGroup, mtllib string) *OBJMesh {
    om := &OBJMesh{Mesh{verts, norms, texcoords, nil}, groups, mtllib}
    om.Objects = make([]*MeshObject, 0, len(groups))
    for _,g := range groups {
        om.Objects = append(om.Objects, &g.MeshObject)
    }
    return om
}

type OBJGroup struct {
    MeshObject
    Smooth bool
    MaterialRef string
}

func newOBJGroup(name string, offset, count int,
                      smooth bool, matref string) *OBJGroup {
    return &OBJGroup{
        MeshObject{name, offset, count, nil},
        smooth, matref}
}
