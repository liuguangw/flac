package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/liuguangw/flac/meta"
)

func parseBlockDataPicture(data []byte, dataTotalLength int) (*meta.BlockDataPicture, error) {
	var (
		tmpLength        uint32
		offset           int
		binaryOrder      = binary.BigEndian
		blockDataPicture = new(meta.BlockDataPicture)
	)
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binaryOrder, &blockDataPicture.PictureType); err != nil {
		return nil, err
	}
	offset += 4
	//MimeType
	if err := binary.Read(buf, binaryOrder, &tmpLength); err != nil {
		return nil, err
	}
	offset += 4
	if tmpLength > 0 {
		binaryData := make([]byte, tmpLength)
		n, err := buf.Read(binaryData)
		if err != nil {
			return nil, err
		}
		if uint32(n) != tmpLength {
			return nil, errors.New("read BlockDataPicture.MimeType failed")
		}
		offset += int(tmpLength)
		blockDataPicture.MimeType = string(binaryData)
	}
	//Description
	if err := binary.Read(buf, binaryOrder, &tmpLength); err != nil {
		return nil, err
	}
	offset += 4
	if tmpLength > 0 {
		binaryData := make([]byte, tmpLength)
		n, err := buf.Read(binaryData)
		if err != nil {
			return nil, err
		}
		if uint32(n) != tmpLength {
			return nil, errors.New("read BlockDataPicture.Description failed")
		}
		offset += int(tmpLength)
		blockDataPicture.Description = string(binaryData)
	}
	//
	if err := binary.Read(buf, binaryOrder, &blockDataPicture.Width); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binaryOrder, &blockDataPicture.Height); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binaryOrder, &blockDataPicture.ColorDepth); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binaryOrder, &blockDataPicture.IndexedColor); err != nil {
		return nil, err
	}
	offset += 16
	//ImageData
	if err := binary.Read(buf, binaryOrder, &tmpLength); err != nil {
		return nil, err
	}
	offset += 4
	if tmpLength > 0 {
		binaryData := make([]byte, tmpLength)
		n, err := buf.Read(binaryData)
		if err != nil {
			return nil, err
		}
		if uint32(n) != tmpLength {
			return nil, errors.New("read BlockDataPicture.ImageData failed")
		}
		offset += int(tmpLength)
		blockDataPicture.ImageData = binaryData
	}
	//长度校验
	if offset != dataTotalLength {
		return nil, errors.New("BlockDataPicture check dataTotalLength failed")
	}
	return blockDataPicture, nil
}
