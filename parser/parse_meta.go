package parser

import (
	"errors"
	"github.com/liuguangw/flac/meta"
	"io"
)

func parseMeta(handle io.ReadSeeker) (*meta.Block, int64, error) {
	var offset int64
	header := make([]byte, 4)
	n, err := handle.Read(header)
	if err != nil {
		return nil, 0, err
	}
	if n != 4 {
		return nil, 0, errors.New("read block header failed")
	}
	offset += int64(n)
	blockType := header[0] & 0x7F
	block := meta.NewBlock(blockType)
	block.SetLastBlock((header[0] >> 7) == 1)
	blockDataLength := (int(header[1]) << 16) + (int(header[2]) << 8) + int(header[3])
	blockData := make([]byte, blockDataLength)
	if n, err = handle.Read(blockData); err != nil {
		return nil, 0, err
	}
	if n != blockDataLength {
		return nil, 0, errors.New("read block data failed")
	}
	offset += int64(n)
	//根据类型解析
	var blockDataInfo meta.BlockData
	if blockType == meta.BlockTypeVorbisComment {
		if blockDataInfo, err = parseBlockDataVorbisComment(blockData, blockDataLength); err != nil {
			return nil, 0, err
		}
	} else if blockType == meta.BlockTypePicture {
		if blockDataInfo, err = parseBlockDataPicture(blockData, blockDataLength); err != nil {
			return nil, 0, err
		}
	} else {
		blockDataInfo = meta.NewBlockDataNop(blockType, blockData)
	}
	block.SetData(blockDataInfo)
	return block, offset, nil
}
