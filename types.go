package go3dm

type TriangleMesh interface {
    Vertices() []float32
    Normals() []float32
    TextureCoords() []float32
    Objects() []*MeshObject
    VTN() ([]float32, []float32, []float32)
}

type MeshObject struct {
    Name string
    FirstFloat int
    FloatCount int
    Smooth bool
}

type triangleMesh struct {
    vertices *f32VA
    normals *f32VA
    texCoords *f32VA
    objects []*MeshObject
}

func (tm *triangleMesh) Vertices() []float32 {
    return tm.vertices.Values
}

func (tm *triangleMesh) Normals() []float32 {
    if tm.Normals == nil { return nil }
    return tm.normals.Values
}

func (tm *triangleMesh) TextureCoords() []float32 {
    if tm.texCoords == nil { return nil }
    return tm.texCoords.Values
}

func (tm *triangleMesh) VTN() ([]float32, []float32, []float32) {
    normals := tm.Normals()
    texc := tm.TextureCoords()
    return tm.vertices.Values, texc, normals
}

func (tm *triangleMesh) Objects() []*MeshObject {
    return tm.objects
}

type f32VA struct {
    Values []float32
    VectorSize int
}

func NewF32VA(vectorSize int) *f32VA {
    return &f32VA{make([]float32, 0, vectorSize*50), vectorSize}
}

func (va *f32VA) AppendVector(vector []float32) {
    for _, v := range vector {
        va.Values = append(va.Values, v)
    }
}

func (va *f32VA) Append(value float32) {
    va.Values = append(va.Values, value)
}

func (va *f32VA) GetVector(idx int) []float32 {
    aIdx := idx * va.VectorSize
    return va.Values[aIdx:aIdx+va.VectorSize]
}

func (va *f32VA) SetVector(idx int, vector []float32) {
    aIdx := idx * va.VectorSize
    for _, val := range vector {
        va.Values[aIdx] = val
        aIdx++
    }
}

func (va *f32VA) VectorCount() int {
    return len(va.Values) / va.VectorSize
}

