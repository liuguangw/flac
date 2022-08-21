package flac

import (
	"errors"
	"github.com/liuguangw/flac/meta"
	"io"
)

func parseMeta(handle io.ReadSeeker) (*meta.Block, int64, error) {
	var (
		block  meta.Block
		offset int64
	)
	n, err := handle.Read(block.Header.Bytes())
	if err != nil {
		return nil, 0, err
	}
	if n != 4 {
		return nil, 0, errors.New("read block header failed")
	}
	offset += int64(n)
	//
	blockLength := block.Header.BlockLength()
	block.Data = make([]byte, blockLength)
	n, err = handle.Read(block.Data)
	if err != nil {
		return nil, 0, err
	}
	if n != blockLength {
		return nil, 0, errors.New("read block data failed")
	}
	offset += int64(n)
	return &block, offset, nil
}
