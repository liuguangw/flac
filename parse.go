package flac

import (
	"errors"
	"github.com/liuguangw/flac/frame"
	"io"
	"os"
)

func Parse(r io.ReadSeekCloser) (*Stream, error) {
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
	stream := &Stream{
		handle: r,
	}
	for {
		blockInfo, offset, err := parseMeta(r)
		if err != nil {
			return nil, err
		}
		seekPos += offset
		stream.BlockList = append(stream.BlockList, blockInfo)
		if blockInfo.Header.IsLastBlock() {
			break
		}
	}
	stream.FrameData = frame.NewBinaryData(seekPos)
	return stream, nil
}

func ParseFile(filename string) (*Stream, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}
	return Parse(f)
}
