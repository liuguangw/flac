package parser

import (
	"errors"
	"github.com/liuguangw/flac/meta"
	"io"
)

func parseMeta(reader io.Reader) (*meta.Block, error) {
	header := make([]byte, 4)
	n, err := reader.Read(header)
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, errors.New("read block header failed")
	}
	blockType := header[0] & 0x7F
	block := meta.NewBlock(blockType)
	block.SetLastBlock((header[0] >> 7) == 1)
	blockDataLength := (int(header[1]) << 16) + (int(header[2]) << 8) + int(header[3])
	blockData := make([]byte, blockDataLength)
	if n, err = reader.Read(blockData); err != nil {
		return nil, err
	}
	if n != blockDataLength {
		return nil, errors.New("read block data failed")
	}
	//根据类型解析
	var blockDataInfo meta.BlockData
	if blockType == meta.BlockTypeVorbisComment {
		if blockDataInfo, err = parseBlockDataVorbisComment(blockData, blockDataLength); err != nil {
			return nil, err
		}
	} else if blockType == meta.BlockTypePicture {
		if blockDataInfo, err = parseBlockDataPicture(blockData, blockDataLength); err != nil {
			return nil, err
		}
	} else {
		blockDataInfo = meta.NewBlockDataNop(blockType, blockData)
	}
	block.SetData(blockDataInfo)
	return block, nil
}
