package workwx

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
)

const mediaFieldName = "media"

// Media 欲上传的素材
//
// NOTE: 由于 Go `mime/multipart` 包的实现细节原因，
// 暂时不开放 Content-Type 定制，全部传 `application/octet-stream`。
// 如有需求，请去 GitHub 提 issue。
type Media struct {
	filename string
	filesize int64
	stream   io.Reader
}

// NewMediaFromFile 从操作系统级文件创建一个欲上传的素材对象
func NewMediaFromFile(f *os.File) (*Media, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &Media{
		filename: stat.Name(),
		filesize: stat.Size(),
		stream:   f,
	}, nil
}

// NewMediaFromBuffer 从内存创建一个欲上传的素材对象
func NewMediaFromBuffer(filename string, buf []byte) (*Media, error) {
	stream := bytes.NewReader(buf)
	return &Media{
		filename: filename,
		filesize: int64(len(buf)),
		stream:   stream,
	}, nil
}

func (m *Media) writeTo(w *multipart.Writer) error {
	wr, err := w.CreateFormFile(mediaFieldName, m.filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(wr, m.stream)
	if err != nil {
		return err
	}

	return nil
}
