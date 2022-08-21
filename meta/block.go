package meta

type BlockData []byte

type Block struct {
	Header BlockHeader
	Data   BlockData
}

// ParseBlockDataVorbisComment 解析标签信息
func (b BlockData) ParseBlockDataVorbisComment() (*BlockDataVorbisComment, error) {
	return parseBlockDataVorbisComment(b)
}

// ParseBlockDataPicture 解析图片信息
func (b BlockData) ParseBlockDataPicture() (*BlockDataPicture, error) {
	return parseBlockDataPicture(b)
}
