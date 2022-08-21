package meta

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

//图片类型定义
const (
	PictureTypeOther             uint32 = iota //Other
	PictureTypePngIcon                         //PNG file icon of 32x32 pixels
	PictureTypeGeneralIcon                     //General file icon
	PictureTypeFrontCover                      //Front cover
	PictureTypeBackCover                       //Back cover
	PictureTypeLinerNotesPage                  //Liner notes page
	PictureTypeMediaLabel                      //Media label (e.g. CD, Vinyl or Cassette label)
	PictureTypeLead                            //Lead artist, lead performer or soloist
	PictureTypeArtist                          //Artist or performer
	PictureTypeConductor                       //Conductor
	PictureTypeBand                            //Band or orchestra
	PictureTypeComposer                        //Composer
	PictureTypeLyricist                        //Lyricist or text writer
	PictureTypeRecording                       //Recording location
	PictureTypeDuringRecording                 //During recording
	PictureTypeDuringPerformance               //During performance
	PictureTypeMovie                           //Movie or video screen capture
	PictureTypeFish                            //A bright colored fish
	PictureTypeIllustration                    //Illustration
	PictureTypeBandOrLogo                      //Band or artist logotype
	PictureTypePublisher                       //Publisher or studio logotype
)

type BlockDataPicture struct {
	PictureType  uint32
	MimeType     string
	Description  string
	Width        uint32
	Height       uint32
	ColorDepth   uint32
	IndexedColor uint32
	ImageData    []byte
}

func (blockDataPicture *BlockDataPicture) BlockType() byte {
	return BlockTypePicture
}

func (blockDataPicture *BlockDataPicture) Length() int {
	dataLength := 4 + 4 + len(blockDataPicture.MimeType) +
		4 + len(blockDataPicture.Description) +
		16 + 4 +
		len(blockDataPicture.ImageData)
	return dataLength
}

func (blockDataPicture *BlockDataPicture) FillFromReader(r io.Reader) error {
	buff := new(bytes.Buffer)
	if _, err := buff.ReadFrom(r); err != nil {
		return err
	}
	blockDataPicture.ImageData = buff.Bytes()
	//fmt.Println("FillFromReader", n)
	return nil
}

// Marshal 序列化
func (blockDataPicture *BlockDataPicture) Marshal(buf *bytes.Buffer) error {
	var (
		tmpLength   uint32
		binaryOrder = binary.BigEndian
	)
	if err := binary.Write(buf, binaryOrder, blockDataPicture.PictureType); err != nil {
		return err
	}
	//MimeType
	tmpLength = uint32(len(blockDataPicture.MimeType))
	if err := binary.Write(buf, binaryOrder, tmpLength); err != nil {
		return err
	}
	if tmpLength > 0 {
		n, err := buf.WriteString(blockDataPicture.MimeType)
		if err != nil {
			return err
		}
		if uint32(n) != tmpLength {
			return errors.New("write BlockDataPicture.MimeType failed")
		}
	}
	//Description
	tmpLength = uint32(len(blockDataPicture.Description))
	if err := binary.Write(buf, binaryOrder, tmpLength); err != nil {
		return err
	}
	if tmpLength > 0 {
		n, err := buf.WriteString(blockDataPicture.Description)
		if err != nil {
			return err
		}
		if uint32(n) != tmpLength {
			return errors.New("write BlockDataPicture.Description failed")
		}
	}
	//
	if err := binary.Write(buf, binaryOrder, blockDataPicture.Width); err != nil {
		return err
	}
	if err := binary.Write(buf, binaryOrder, blockDataPicture.Height); err != nil {
		return err
	}
	if err := binary.Write(buf, binaryOrder, blockDataPicture.ColorDepth); err != nil {
		return err
	}
	if err := binary.Write(buf, binaryOrder, blockDataPicture.IndexedColor); err != nil {
		return err
	}
	//ImageData
	tmpLength = uint32(len(blockDataPicture.ImageData))
	if err := binary.Write(buf, binaryOrder, tmpLength); err != nil {
		return err
	}
	if tmpLength > 0 {
		n, err := buf.Write(blockDataPicture.ImageData)
		if err != nil {
			return err
		}
		if uint32(n) != tmpLength {
			return errors.New("write BlockDataPicture.ImageData failed")
		}
	}
	return nil
}
