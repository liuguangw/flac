package common

import (
	"errors"
	"github.com/liuguangw/flac/meta"
	"io"
	"os"
)

type Stream struct {
	reader    io.Reader
	BlockList []*meta.Block
}

func NewStream(reader io.Reader) *Stream {
	return &Stream{reader: reader}
}

// Save 保存输出，输入流和输出流不能相同
func (stream *Stream) Save(out io.Writer) error {
	flacTag := []byte("fLaC")
	n, err := out.Write(flacTag)
	if err != nil {
		return err
	}
	if n != 4 {
		return errors.New("write fLaC tag failed")
	}
	lastBlockIndex := len(stream.BlockList) - 1
	for blockIndex, blockInfo := range stream.BlockList {
		//是否为最后一个Block
		blockInfo.SetLastBlock(blockIndex == lastBlockIndex)
		blockData, err := blockInfo.Marshal()
		if err != nil {
			return err
		}
		if _, err := out.Write(blockData); err != nil {
			return err
		}
	}
	if _, err := io.Copy(out, stream.reader); err != nil {
		return err
	}
	return nil
}

func (stream *Stream) SaveFile(filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return stream.Save(f)
}
