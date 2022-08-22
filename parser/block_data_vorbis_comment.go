package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/liuguangw/flac/meta"
	"strconv"
	"strings"
)

func parseBlockDataVorbisComment(blockData []byte, dataTotalLength int) (*meta.BlockDataVorbisComment, error) {
	var (
		tmpLength   uint32
		offset      int
		binaryOrder = binary.LittleEndian
	)
	buf := bytes.NewBuffer(blockData)
	if err := binary.Read(buf, binaryOrder, &tmpLength); err != nil {
		return nil, err
	}
	offset += 4
	binaryData := make([]byte, tmpLength)
	n, err := buf.Read(binaryData)
	if err != nil {
		return nil, err
	}
	if uint32(n) != tmpLength {
		return nil, errors.New("read BlockDataVorbisComment.Vendor failed")
	}
	offset += int(tmpLength)
	vorbisComment := &meta.BlockDataVorbisComment{
		Vendor: string(binaryData),
	}
	if err := binary.Read(buf, binaryOrder, &tmpLength); err != nil {
		return nil, err
	}
	offset += 4
	vorbisComment.FieldList = make([][2]string, tmpLength)
	for fieldIndex := range vorbisComment.FieldList {
		if err := binary.Read(buf, binaryOrder, &tmpLength); err != nil {
			return nil, err
		}
		offset += 4
		binaryData = make([]byte, tmpLength)
		n, err := buf.Read(binaryData)
		if err != nil {
			return nil, err
		}
		if uint32(n) != tmpLength {
			return nil, errors.New("read BlockDataVorbisComment.FieldList[" + strconv.Itoa(fieldIndex) + "] length failed")
		}
		offset += int(tmpLength)
		tmpStr := string(binaryData)
		pos := strings.IndexByte(tmpStr, '=')
		if pos > 0 {
			vorbisComment.FieldList[fieldIndex][0] = tmpStr[:pos]
			vorbisComment.FieldList[fieldIndex][1] = tmpStr[pos+1:]
		}
	}
	//长度校验
	if offset != dataTotalLength {
		return nil, errors.New("BlockDataVorbisComment check dataTotalLength failed")
	}
	return vorbisComment, nil
}
