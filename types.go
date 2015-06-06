package go3dm

// Interfaces

type Mesh interface {
    Vertices() []float32
    Normals() []float32
    TextureCoords() []float32
    Elements() []uint32
    Objects() []MeshObject
    VTN() ([]float32, []float32, []float32)
}

type MeshObject interface {
    Name() string
    ElementOffset() int32
    ElementCount() int32
    MaterialRef() string
    Smooth() bool
}

type Material interface {
    Name() string
    Ka() []float32
    Kd() []float32
    Ks() []float32
    Ns() float32
    Tr() float32
    KaMapName() string
    KdMapName() string
    KsMapName() string
}

// Public Structs

type BasicMesh struct {
    vertices []float32
    normals []float32
    textureCoords []float32
    elements []uint32
    objects []MeshObject
}

func (m *BasicMesh) Vertices() []float32 {
    return m.vertices
}

func (m *BasicMesh) Normals() []float32 {
    return m.normals
}

func (m *BasicMesh) TextureCoords() []float32 {
    return m.textureCoords
}

func (m *BasicMesh) Elements() []uint32 {
    return m.elements
}

func (m *BasicMesh) Objects() []MeshObject {
    return m.objects
}

func (m *BasicMesh) VTN() ([]float32, []float32, []float32) {
    return m.vertices, m.textureCoords, m.normals
}

type BasicMeshObject struct {
    name string
    elementOffset int
    elementCount int
    materialRef string
    smooth bool
}

func (mo *BasicMeshObject) Name() string {
    return mo.name
}

func (mo *BasicMeshObject) ElementOffset() int32 {
    return int32(mo.elementOffset)
}

func (mo *BasicMeshObject) ElementCount() int32 {
    return int32(mo.elementCount)
}

func (mo *BasicMeshObject) MaterialRef() string {
    return mo.materialRef
}

func (mo *BasicMeshObject) Smooth() bool {
    return mo.smooth
}

type BasicMaterial struct {
    name string
    ka []float32
    kd []float32
    ks []float32
    ns float32
    tr float32
    kaMapName string
    kdMapName string
    ksMapName string
}

func NewBasicMaterial(name string) Material {
    mat := new(BasicMaterial)
    mat.name = name
    return mat
}

func (m *BasicMaterial) Name() string {
    return m.name
}

func (m *BasicMaterial) Ka() []float32 {
    return m.ka
}

func (m *BasicMaterial) Kd() []float32 {
    return m.kd
}

func (m *BasicMaterial) Ks() []float32 {
    return m.ks
}

func (m *BasicMaterial) Ns() float32 {
    return m.ns
}

func (m *BasicMaterial) Tr() float32 {
    return m.tr
}

func (m *BasicMaterial) KaMapName() string {
    return m.kaMapName
}

func (m *BasicMaterial) KdMapName() string {
    return m.kdMapName
}

func (m *BasicMaterial) KsMapName() string {
    return m.ksMapName
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

