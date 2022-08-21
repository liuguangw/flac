package meta

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

type BlockDataVorbisComment struct {
	Vendor    string
	FieldList [][2]string
}

func parseBlockDataVorbisComment(data []byte) (*BlockDataVorbisComment, error) {
	var (
		tmpLength   uint32
		binaryOrder = binary.LittleEndian
	)
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binaryOrder, &tmpLength); err != nil {
		return nil, err
	}
	binaryData := make([]byte, tmpLength)
	n, err := buf.Read(binaryData)
	if err != nil {
		return nil, err
	}
	if uint32(n) != tmpLength {
		return nil, errors.New("read BlockDataVorbisComment.Vendor failed")
	}
	vorbisComment := &BlockDataVorbisComment{
		Vendor: string(binaryData),
	}
	if err := binary.Read(buf, binaryOrder, &tmpLength); err != nil {
		return nil, err
	}
	vorbisComment.FieldList = make([][2]string, tmpLength)
	for fieldIndex := range vorbisComment.FieldList {
		if err := binary.Read(buf, binaryOrder, &tmpLength); err != nil {
			return nil, err
		}
		binaryData = make([]byte, tmpLength)
		n, err := buf.Read(binaryData)
		if err != nil {
			return nil, err
		}
		if uint32(n) != tmpLength {
			return nil, errors.New("read BlockDataVorbisComment.FieldList[" + strconv.Itoa(fieldIndex) + "] length failed")
		}
		tmpStr := string(binaryData)
		pos := strings.IndexByte(tmpStr, '=')
		if pos > 0 {
			vorbisComment.FieldList[fieldIndex][0] = tmpStr[:pos]
			vorbisComment.FieldList[fieldIndex][1] = tmpStr[pos+1:]
		}
	}
	//fmt.Printf("#offset=%d\n", offset)
	return vorbisComment, nil
}

func (vorbisComment *BlockDataVorbisComment) Length() int64 {
	dataLength := 4 + int64(len(vorbisComment.Vendor)) + 4
	for fieldIndex := range vorbisComment.FieldList {
		dataLength += 4 + int64(len(vorbisComment.FieldList[fieldIndex][0])) +
			1 + int64(len(vorbisComment.FieldList[fieldIndex][1]))
	}
	return dataLength
}

// Marshal 序列化
func (vorbisComment *BlockDataVorbisComment) Marshal() ([]byte, error) {
	binaryData := make([]byte, 0, vorbisComment.Length())
	var (
		tmpLength   uint32
		binaryOrder = binary.LittleEndian
	)
	tmpLength = uint32(len(vorbisComment.Vendor))
	buf := bytes.NewBuffer(binaryData)
	if err := binary.Write(buf, binaryOrder, tmpLength); err != nil {
		return nil, err
	}
	if tmpLength > 0 {
		n, err := buf.WriteString(vorbisComment.Vendor)
		if err != nil {
			return nil, err
		}
		if uint32(n) != tmpLength {
			return nil, errors.New("write vorbisComment.Vendor failed")
		}
	}
	tmpLength = uint32(len(vorbisComment.FieldList))
	if err := binary.Write(buf, binaryOrder, tmpLength); err != nil {
		return nil, err
	}
	for fieldIndex := range vorbisComment.FieldList {
		tmpStr := vorbisComment.FieldList[fieldIndex][0] + "=" + vorbisComment.FieldList[fieldIndex][1]
		tmpLength = uint32(len(tmpStr))
		if err := binary.Write(buf, binaryOrder, tmpLength); err != nil {
			return nil, err
		}
		if tmpLength > 0 {
			n, err := buf.WriteString(tmpStr)
			if err != nil {
				return nil, err
			}
			if uint32(n) != tmpLength {
				return nil, errors.New("write vorbisComment.FieldList[" + strconv.Itoa(fieldIndex) + "] failed")
			}
		}
	}
	return buf.Bytes(), nil
}
