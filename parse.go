package flac

import (
	"github.com/liuguangw/flac/common"
	"github.com/liuguangw/flac/parser"
	"io"
)

// Parse 从输入流中解析Stream,输入流的生命周期必须>=Stream的生命周期
func Parse(r io.Reader) (*common.Stream, error) {
	return parser.Parse(r)
}
