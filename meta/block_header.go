package meta

type BlockHeader [4]byte

func (b *BlockHeader) IsLastBlock() bool {
	data := b[0] >> 7
	return data == 1
}
func (b *BlockHeader) SetLastBlock(isLast bool) {
	if b.IsLastBlock() == isLast {
		return
	}
	b[0] &= 0x7F
	if isLast {
		b[0] |= 0x80
	}
}

func (b *BlockHeader) BlockType() byte {
	return b[0] & 0x7F
}

func (b *BlockHeader) SetBlockType(blockType byte) {
	b[0] &= 0x80
	b[0] |= blockType & 0x7F
}

func (b *BlockHeader) BlockLength() int {
	return (int(b[1]) << 16) + (int(b[2]) << 8) + int(b[3])
}

func (b *BlockHeader) SetBlockLength(blockLength int) {
	b[1] = byte((blockLength >> 16) & 0xFF)
	b[2] = byte((blockLength >> 8) & 0xFF)
	b[3] = byte(blockLength & 0xFF)
}

func (b *BlockHeader) Bytes() []byte {
	return b[:]
}
