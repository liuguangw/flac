package flac

import (
	"errors"
	"github.com/liuguangw/flac/frame"
	"github.com/liuguangw/flac/meta"
	"io"
	"os"
)

type Stream struct {
	handle    io.ReadSeekCloser
	BlockList []*meta.Block
	FrameData *frame.BinaryData
}

// Close 调用Close会关闭底层的输入流
func (stream *Stream) Close() error {
	return stream.handle.Close()
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
	blockCount := len(stream.BlockList)
	for blockIndex, blockInfo := range stream.BlockList {
		//计算header
		blockInfo.Header.SetLastBlock(blockIndex == blockCount-1)
		blockInfo.Header.SetBlockLength(len(blockInfo.Data))
		//写入Header
		n, err = out.Write(blockInfo.Header.Bytes())
		if err != nil {
			return err
		}
		if n != 4 {
			return errors.New("write Block.Header failed")
		}
		//写入Data
		n, err = out.Write(blockInfo.Data)
		if err != nil {
			return err
		}
		if n != blockInfo.Header.BlockLength() {
			return errors.New("write Block.Data failed")
		}
	}
	if _, err := stream.handle.Seek(stream.FrameData.Offset(), io.SeekStart); err != nil {
		return err
	}
	if _, err := io.Copy(out, stream.handle); err != nil {
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
