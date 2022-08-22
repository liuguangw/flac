package main

import (
	"fmt"
	"github.com/liuguangw/flac"
	"github.com/liuguangw/flac/meta"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("E:\\go_projects\\src\\hello\\q_org.flac", os.O_RDONLY, 0400)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	if err := processMusic(f); err != nil {
		log.Fatalln(err)
	}
}

func processMusic(r io.Reader) error {
	stream, err := flac.Parse(r)
	if err != nil {
		return err
	}
	for i, block := range stream.BlockList {
		fmt.Printf("BlockList[%d] is_last=%v, type=%d, length=%d\n",
			i, block.IsLastBlock(), block.BlockType(), block.Data().Length())
		switch blockData := block.Data().(type) {
		case *meta.BlockDataVorbisComment:
			fmt.Println("\tvendor=" + blockData.Vendor)
			hasDate := false
			for fieldIndex, fieldInfo := range blockData.FieldList {
				fmt.Printf("\t\tFieldList[%d] %s=%s\n", fieldIndex, fieldInfo[0], fieldInfo[1])
				if fieldInfo[0] == "DATE" {
					hasDate = true
				}
			}
			if !hasDate {
				blockData.FieldList = append(blockData.FieldList, [2]string{
					"DATE", "2022-08-21",
				})
				blockData.FieldList = append(blockData.FieldList, [2]string{
					"YEAR", "2022",
				})
			}
		case *meta.BlockDataPicture:
			fmt.Println(blockData.PictureType, blockData.MimeType)
		default:
			fmt.Println("<nop>")
		}
	}
	pictureBlock, err := createFlacMetaPicture("C:\\Users\\liuguang\\Pictures\\TEHeTtW2_400x400.jpg")
	if err != nil {
		return err
	}
	stream.BlockList = append(stream.BlockList, pictureBlock)
	if err := stream.SaveFile("E:\\go_projects\\src\\hello\\qmm.flac"); err != nil {
		return err
	}
	return nil
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
	block := meta.NewBlock()
	block.SetData(picture)
	return block, nil
}
