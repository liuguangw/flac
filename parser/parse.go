package parser

import (
	"errors"
	"github.com/liuguangw/flac/common"
	"io"
)

// Parse 输入流解析
func Parse(r io.Reader) (*common.Stream, error) {
	var buffData [4]byte
	n, err := r.Read(buffData[:])
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, errors.New("read flac tag failed")
	}
	if string(buffData[:]) != "fLaC" {
		return nil, errors.New("invalid file, not flac format")
	}
	stream := common.NewStream(r)
	for {
		blockInfo, err := parseMeta(r)
		if err != nil {
			return nil, err
		}
		stream.BlockList = append(stream.BlockList, blockInfo)
		if blockInfo.IsLastBlock() {
			break
		}
	}
	return stream, nil
}
