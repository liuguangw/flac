package flac

import (
	"github.com/liuguangw/flac/common"
	"github.com/liuguangw/flac/parser"
	"io"
	"os"
)

func Parse(r io.ReadSeekCloser) (*common.Stream, error) {
	return parser.Parse(r)
}

func ParseFile(filename string) (*common.Stream, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}
	return parser.Parse(f)
}
