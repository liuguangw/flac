package main

import (
	"fmt"
	"github.com/liuguangw/flac"
	"github.com/liuguangw/flac/meta"
	"os"
)

func main() {
	stream, err := flac.ParseFile("E:\\go_projects\\src\\hello\\q_org.flac")
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	for i, item := range stream.BlockList {
		header := &item.Header
		fmt.Printf("BlockList[%d] is_last=%v, type=%d, length=%d\n",
			i, header.IsLastBlock(), header.BlockType(), header.BlockLength())
		if header.BlockType() == meta.BlockTypeVorbisComment {
			dataVorbisComment, err := item.Data.ParseBlockDataVorbisComment()
			if err != nil {
				panic(err)
			}
			fmt.Println("\tvendor=" + dataVorbisComment.Vendor)
			hasDate := false
			for fieldIndex, fieldInfo := range dataVorbisComment.FieldList {
				fmt.Printf("\t\tFieldList[%d] %s=%s\n", fieldIndex, fieldInfo[0], fieldInfo[1])
				if fieldInfo[0] == "DATE" {
					hasDate = true
				}
			}
			if !hasDate {
				dataVorbisComment.FieldList = append(dataVorbisComment.FieldList, [2]string{
					"DATE", "2022-08-21",
				})
				dataVorbisComment.FieldList = append(dataVorbisComment.FieldList, [2]string{
					"YEAR", "2022",
				})
				newData, err := dataVorbisComment.Marshal()
				if err != nil {
					panic(err)
				}
				item.Data = newData
			}
		} else if header.BlockType() == meta.BlockTypePicture {
			picture, err := item.Data.ParseBlockDataPicture()
			if err != nil {
				panic(err)
			}
			fmt.Println(picture.PictureType, picture.MimeType)
		}
	}
	pictureBlock, err := createFlacMetaPicture("C:\\Users\\liuguang\\Pictures\\TEHeTtW2_400x400.jpg")
	if err != nil {
		panic(err)
	}
	stream.BlockList = append(stream.BlockList, pictureBlock)
	if err := stream.SaveFile("E:\\go_projects\\src\\hello\\qmm.flac"); err != nil {
		panic(err)
	}
}

// createFlacMetaPicture 写入封面图
func createFlacMetaPicture(fPath string) (*meta.Block, error) {
	picture := new(meta.BlockDataPicture)
	picture.MimeType = "image/jpeg"
	picture.PictureType = meta.PictureTypeFrontCover
	f, err := os.OpenFile(fPath, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := picture.FillFromReader(f); err != nil {
		return nil, err
	}
	//构造block
	block := new(meta.Block)
	block.Header.SetBlockType(meta.BlockTypePicture)
	newData, err := picture.Marshal()
	if err != nil {
		return nil, err
	}
	block.Data = newData
	return block, nil
}
