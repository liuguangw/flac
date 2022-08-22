package meta

import (
	"encoding/binary"
	"errors"
	"io"
	"strconv"
)

type BlockDataVorbisComment struct {
	Vendor    string
	FieldList [][2]string
}

func (vorbisComment *BlockDataVorbisComment) BlockType() byte {
	return BlockTypeVorbisComment
}

func (vorbisComment *BlockDataVorbisComment) Length() int {
	dataLength := 4 + len(vorbisComment.Vendor) + 4
	for fieldIndex := range vorbisComment.FieldList {
		dataLength += 4 + len(vorbisComment.FieldList[fieldIndex][0]) +
			1 + len(vorbisComment.FieldList[fieldIndex][1])
	}
	return dataLength
}

// Marshal 序列化
func (vorbisComment *BlockDataVorbisComment) Marshal(w io.Writer) error {
	var (
		tmpLength   uint32
		binaryOrder = binary.LittleEndian
	)
	tmpLength = uint32(len(vorbisComment.Vendor))
	if err := binary.Write(w, binaryOrder, tmpLength); err != nil {
		return err
	}
	if tmpLength > 0 {
		n, err := w.Write([]byte(vorbisComment.Vendor))
		if err != nil {
			return err
		}
		if uint32(n) != tmpLength {
			return errors.New("write vorbisComment.Vendor failed")
		}
	}
	tmpLength = uint32(len(vorbisComment.FieldList))
	if err := binary.Write(w, binaryOrder, tmpLength); err != nil {
		return err
	}
	for fieldIndex := range vorbisComment.FieldList {
		tmpStr := vorbisComment.FieldList[fieldIndex][0] + "=" + vorbisComment.FieldList[fieldIndex][1]
		tmpLength = uint32(len(tmpStr))
		if err := binary.Write(w, binaryOrder, tmpLength); err != nil {
			return err
		}
		if tmpLength > 0 {
			n, err := w.Write([]byte(tmpStr))
			if err != nil {
				return err
			}
			if uint32(n) != tmpLength {
				return errors.New("write vorbisComment.FieldList[" + strconv.Itoa(fieldIndex) + "] failed")
			}
		}
	}
	return nil
}
