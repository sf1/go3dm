package go3dm

type TriangleMesh interface {
    Vertices() []float32
    Normals() []float32
    TextureCoords() []float32
    VTN() ([]float32, []float32, []float32)
}

type OBJMesh struct {
    vertices *f32VA
    normals *f32VA
    texCoords *f32VA
    objects []*OBJMeshObject
    mtllib string
}

func (om *OBJMesh) Vertices() []float32 {
    return om.vertices.Values
}

func (om *OBJMesh) Normals() []float32 {
    if om.normals == nil { return nil }
    return om.normals.Values
}

func (om *OBJMesh) TextureCoords() []float32 {
    if om.texCoords == nil { return nil }
    return om.texCoords.Values
}

func (om *OBJMesh) VTN() ([]float32, []float32, []float32) {
    normals := om.Normals()
    texc := om.TextureCoords()
    return om.vertices.Values, texc, normals
}

func (om *OBJMesh) MTLLib() string {
    return om.mtllib
}

type OBJMeshObject struct {
    Name string
    FirstFloat int
    FloatCount int
    Smooth bool
    MaterialRef string
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

