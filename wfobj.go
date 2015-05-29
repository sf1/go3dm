package go3dm

import (
    "io"
    "bufio"
    "strings"
    "strconv"
    "fmt"
)

func LoadObj(reader io.Reader) error {
    scanner := bufio.NewScanner(reader)
    var err error
    vertices := make([]float32, 0, 50)
    normalsTmp := make([]float32, 0, 50)
    texTmp := make([]float32, 0, 50)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        tokens := strings.Split(line, " ")
        switch tokens[0] {
        case "g","o":
            fmt.Println("Object/Group", tokens[1:])
        case "v":
            fmt.Println("Vertex", tokens[1:])
            vertices, err = parseAndAppendF32(tokens[1:], vertices)
            if err != nil {return err}
        case "vn":
            fmt.Println("Vertex Normal", tokens[1:])
            normalsTmp, err = parseAndAppendF32(tokens[1:], normalsTmp)
            if err != nil {return err}
        case "vt":
            fmt.Println("Texture Coordinate", tokens[1:])
            texTmp, err = parseAndAppendF32(tokens[1:], texTmp)
            if err != nil {return err}
        case "f":
            fmt.Println("Face", tokens[1:])
        case "s":
            fmt.Println("Smooth Shading Flag", tokens[1:])
        case "mtllib":
            fmt.Println("Material Library", tokens[1:])
        case "usemtl":
            fmt.Println("Material Reference", tokens[1:])

        }
    }
    fmt.Println(vertices)
    return nil
}

func parseAndAppendF32(tokens []string, floats []float32) ([]float32, error) {
    for _,t := range tokens {
        v, err := strconv.ParseFloat(t, 32)
        if err != nil {return nil, err}
        floats = append(floats, float32(v))
    }
    return floats, nil
}
