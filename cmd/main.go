package main

import (
	"fmt"
	"github.com/liuguangw/flac"
	"github.com/liuguangw/flac/meta"
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
				if err := stream.SaveFile("E:\\go_projects\\src\\hello\\qmm.flac"); err != nil {
					panic(err)
				}
			}
		}
	}
}
