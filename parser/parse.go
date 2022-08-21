package parser

import (
	"errors"
	"github.com/liuguangw/flac/common"
	"github.com/liuguangw/flac/frame"
	"io"
)

func Parse(r io.ReadSeekCloser) (*common.Stream, error) {
	seekPos, err := r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	var buffData [4]byte
	n, err := r.Read(buffData[:])
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, errors.New("read flac tag failed")
	}
	seekPos += int64(n)
	if string(buffData[:]) != "fLaC" {
		return nil, errors.New("invalid file, not flac format")
	}
	stream := common.NewStream(r)
	for {
		blockInfo, offset, err := parseMeta(r)
		if err != nil {
			return nil, err
		}
		seekPos += offset
		stream.BlockList = append(stream.BlockList, blockInfo)
		if blockInfo.IsLastBlock() {
			break
		}
	}
	stream.FrameData = frame.NewBinaryData(seekPos)
	return stream, nil
}
