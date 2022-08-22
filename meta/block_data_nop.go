package meta

import (
	"errors"
	"io"
)

// BlockDataNop 代表不能解析的BlockData
type BlockDataNop struct {
	blockType byte
	data      []byte
}

func NewBlockDataNop(blockType byte, data []byte) *BlockDataNop {
	return &BlockDataNop{
		blockType: blockType,
		data:      data,
	}
}

func (b *BlockDataNop) BlockType() byte {
	return b.blockType
}

func (b *BlockDataNop) Length() int {
	return len(b.data)
}

func (b *BlockDataNop) Marshal(w io.Writer) error {
	n, err := w.Write(b.data)
	if err != nil {
		return err
	}
	if b.Length() != n {
		return errors.New("marshal BlockDataNop failed")
	}
	return nil
}
