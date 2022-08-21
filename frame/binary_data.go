package frame

type BinaryData struct {
	offset int64
}

func NewBinaryData(offset int64) *BinaryData {
	return &BinaryData{offset: offset}
}
func (b *BinaryData) Offset() int64 {
	return b.offset
}
