package go3dm

// Public Structs

type TriangleMesh struct {
    Vertices []float32
    Normals []float32
    TextureCoords []float32
    Indicies []uint32
    Objects []*MeshObject
}

func (m *TriangleMesh) VTN() ([]float32, []float32, []float32) {
    return m.Vertices, m.TextureCoords, m.Normals
}

type MeshObject struct {
    Name string
    IndexOffset int32
    IndexCount int32
    MaterialRef string
    Smooth bool
}

type Material struct {
    Name string
    Ka []float32
    Kd []float32
    Ks []float32
    Ns float32
    Tr float32
    KaMapName string
    KdMapName string
    KsMapName string
}

// Internal Structs

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

