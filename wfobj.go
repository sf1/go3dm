package go3dm

import (
    "io"
    "bufio"
    "strings"
    "fmt"
)

func LoadObj(reader io.Reader) {
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        tokens := strings.Split(line, " ")
        switch tokens[0] {
        case "g","o":
            fmt.Println("Object/Group", tokens[1:])
        case "v":
            fmt.Println("Vertex", tokens[1:])
        case "vn":
            fmt.Println("Vertex Normal", tokens[1:])
        case "vt":
            fmt.Println("Texture Coordinate", tokens[1:])
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
}
