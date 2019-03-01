package vegeta

// Format defines a type for the format query param
type Format interface {
	GetFormat() string
	SetFormat(meta ...string)
	GetMetaInfo() []string
}

// NewFormat returns a new Format in accordance with the provided type string, by default it is JSON
func NewFormat(typ string) Format {
	switch typ {
	case "json":
		f := JSONFormat("json")
		return &f
	case "text":
		f := TextFormat("text")
		return &f
	case "binary":
		f := BinaryFormat("binary")
		return &f
	case "histogram":
		return &HistogramFormat{
			data: "histogram",
		}
	}
	return new(JSONFormat) // default
}

// JSONFormat typedef for query param "json"
type JSONFormat string

// SetFormat will set the meta information strings of the JSONFormat
func (j *JSONFormat) SetFormat(_ ...string) {
	*j = "json"
}

// GetFormat will return the type of the JSONFormat
func (j *JSONFormat) GetFormat() string {
	return string(*j)
}

// GetMetaInfo will return the meta information of the JSONFormat
func (j *JSONFormat) GetMetaInfo() (m []string) {
	return
}

// TextFormat typedef for query param "text"
type TextFormat string

// SetFormat will set the meta information strings of the TextFormat
func (t *TextFormat) SetFormat(_ ...string) {
	*t = "text"
}

// GetFormat will return the type of the TextFormat
func (t *TextFormat) GetFormat() string {
	return string(*t)
}

// GetMetaInfo will return the meta information of the TextFormat
func (t *TextFormat) GetMetaInfo() (m []string) {
	return
}

// BinaryFormat typedef for query param "binary"
type BinaryFormat string

// SetFormat will set the meta information strings of the BinaryFormat
func (b *BinaryFormat) SetFormat(_ ...string) {
	*b = "binary"
}

// GetFormat will return the type of the BinaryFormat
func (b *BinaryFormat) GetFormat() string {
	return string(*b)
}

// GetMetaInfo will return the meta information of the BinaryFormat
func (b *BinaryFormat) GetMetaInfo() (m []string) {
	return
}

// HistogramFormat typedef for query param "histogram"
type HistogramFormat struct {
	data, metadata string
}

// SetFormat will set the meta information strings of the HistogramFormat
func (h *HistogramFormat) SetFormat(meta ...string) {
	h.data = "histogram"
	h.metadata = meta[0]
}

// GetFormat will return the type of the HistogramFormat
func (h *HistogramFormat) GetFormat() string {
	return h.data
}

// GetMetaInfo will return the meta information of the HistogramFormat
func (h *HistogramFormat) GetMetaInfo() []string {
	return []string{h.metadata}
}
