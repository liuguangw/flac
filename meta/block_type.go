package meta

const (
	BlockTypeStreamInfo byte = iota
	BlockTypePadding
	BlockTypeApplication
	BlockTypeSeekTable
	BlockTypeVorbisComment
	BlockTypeCueSheet
	BlockTypePicture
	BlockTypeReserved
	BlockTypeInvalid byte = 127
)
