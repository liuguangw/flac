package meta

import (
	"bytes"
	"errors"
	"io"
)

type BlockData interface {
	BlockType() byte
	Length() int
	Marshal(w io.Writer) error
}

type Block struct {
	isLast    bool
	blockType byte
	data      BlockData
}

func NewBlock() *Block {
	return &Block{blockType: BlockTypeInvalid}
}

func (b *Block) IsLastBlock() bool {
	return b.isLast
}

func (b *Block) SetLastBlock(isLast bool) {
	b.isLast = isLast
}

func (b *Block) BlockType() byte {
	return b.blockType
}

func (b *Block) Data() BlockData {
	return b.data
}

func (b *Block) SetData(data BlockData) {
	b.data = data
	b.blockType = data.BlockType()
}
func (b *Block) Marshal() ([]byte, error) {
	if b.blockType == BlockTypeInvalid {
		return nil, errors.New("invalid block type")
	}
	dataLength := b.data.Length()
	binaryData := make([]byte, 0, 4+dataLength)
	buf := bytes.NewBuffer(binaryData)
	var data byte
	if b.isLast {
		data = 0x80
	}
	data |= b.blockType & 0x7F
	if err := buf.WriteByte(data); err != nil {
		return nil, err
	}
	//长度三个字节
	binaryData = []byte{
		byte((dataLength >> 16) & 0xFF),
		byte((dataLength >> 8) & 0xFF),
		byte(dataLength & 0xFF),
	}
	n, err := buf.Write(binaryData)
	if err != nil {
		return nil, err
	}
	if n != 3 {
		return nil, errors.New("marshal Block header length failed")
	}
	if err := b.data.Marshal(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
