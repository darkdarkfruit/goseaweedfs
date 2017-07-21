package model

import (
	"io"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// FilePart ...
type FilePart struct {
	Reader     io.Reader
	FileName   string
	FileSize   int64
	IsGzipped  bool
	MimeType   string
	ModTime    int64 //in seconds
	Collection string

	// Ttl Time to live.
	// 3m: 3 minutes
	// 4h: 4 hours
	// 5d: 5 days
	// 6w: 6 weeks
	// 7M: 7 months
	// 8y: 8 years
	Ttl string

	Server string
	FileID string
}

// NewFilePartFromReader ...
func NewFilePartFromReader(reader io.Reader, fileName string, fileSize int64) *FilePart {
	ret := FilePart{
		Reader:   reader,
		FileSize: fileSize,
		FileName: fileName,
	}

	ext := strings.ToLower(path.Ext(fileName))
	if ext != "" {
		ret.MimeType = mime.TypeByExtension(ext)
	}
	ret.IsGzipped = ext == ".gz"

	return &ret
}

// NewFilePart ...
func NewFilePart(fullPathFilename string) (*FilePart, error) {
	ret := FilePart{}

	fh, openErr := os.Open(fullPathFilename)
	if openErr != nil {
		return nil, openErr
	}
	ret.Reader = fh
	ret.FileName = filepath.Base(fullPathFilename)

	if fi, fiErr := fh.Stat(); fiErr != nil {
		return nil, fiErr
	} else {
		ret.ModTime = fi.ModTime().UTC().Unix()
		ret.FileSize = fi.Size()
	}

	ext := strings.ToLower(path.Ext(ret.FileName))
	if ext != "" {
		ret.MimeType = mime.TypeByExtension(ext)
	}
	ret.IsGzipped = ext == ".gz"

	return &ret, nil
}

func NewFileParts(fullPathFilenames []string) (ret []*FilePart, err error) {
	ret = make([]*FilePart, len(fullPathFilenames))
	for index, file := range fullPathFilenames {
		if ret[index], err = NewFilePart(file); err != nil {
			return
		}
	}
	return
}