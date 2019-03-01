package vegeta

// MetaInfo provide map to store meta information for Format
type MetaInfo map[string]string

// Format defines a type for the format query param
type Format interface {
	// fmt.Stringer
	String() string

	Meta() MetaInfo
	SetMeta(key, value string)
}

// NewFormat returns a new Format in accordance with the provided type string, by default it is JSON
func NewFormat(typ string) Format {
	switch typ {
	case "json":
		return NewJSONFormat()
	case "text":
		return NewTextFormat()
	case "binary":
		return NewBinaryFormat()
	case "histogram":
		return NewHistogramFormat()
	}
	return NewJSONFormat() // default
}

// JSONFormat typedef for query param "json"
type JSONFormat string

// NewJSONFormat returns a new Format of JSON type
func NewJSONFormat() *JSONFormat {
	f := JSONFormat("json")
	return &f
}

// SetMeta will set the meta information
func (j *JSONFormat) SetMeta(key, value string) {
}

// String implements Stringer for JSONFormat
func (j *JSONFormat) String() string {
	return string(*j)
}

// Meta returns the meta information stored in the Format
func (j *JSONFormat) Meta() (m MetaInfo) {
	return nil
}

// TextFormat typedef for query param "text"
type TextFormat string

// NewTextFormat returns a new Format of Text type
func NewTextFormat() *TextFormat {
	f := TextFormat("text")
	return &f
}

// SetMeta will set the meta information
func (j *TextFormat) SetMeta(key, value string) {
}

// String implements Stringer for TextFormat
func (j *TextFormat) String() string {
	return string(*j)
}

// Meta returns the meta information stored in the Format
func (j *TextFormat) Meta() (m MetaInfo) {
	return nil
}

// BinaryFormat typedef for query param "binary"
type BinaryFormat string

// NewBinaryFormat returns a new Format of Binary type
func NewBinaryFormat() *BinaryFormat {
	f := BinaryFormat("binary")
	return &f
}

// SetMeta will set the meta information
func (b *BinaryFormat) SetMeta(key, value string) {
}

// String implements Stringer for BinaryFormat
func (b *BinaryFormat) String() string {
	return string(*b)
}

// Meta returns the meta information stored in the Format
func (b *BinaryFormat) Meta() (m MetaInfo) {
	return
}

// HistogramFormat typedef for query param "histogram"
type HistogramFormat struct {
	repr string
	meta MetaInfo
}

// NewHistogramFormat returns a new Format of Histogram type
func NewHistogramFormat() *HistogramFormat {
	return &HistogramFormat{
		repr: "histogram",
		meta: make(MetaInfo),
	}
}

// SetMeta will set the meta information
func (h *HistogramFormat) SetMeta(key, value string) {
	h.meta[key] = value
}

// String implements Stringer for HistogramFormat
func (h *HistogramFormat) String() string {
	return h.repr
}

// Meta returns the meta information stored in the Format
func (h *HistogramFormat) Meta() MetaInfo {
	return h.meta
}
