package data

import (
	"embed"
	"io"
	"os"
	"time"
)

//go:embed bus-stops.csv
var DataFS embed.FS

type FileInfo struct {
	ModTime time.Time
	IsOS    bool
}

type BusStopReader struct {
	reader io.ReadCloser
	info   FileInfo
}

func NewBusStopReader() (*BusStopReader, error) {
	// try local file
	if file, err := os.Open("bus-stops.csv"); err == nil {
		stat, err := file.Stat()
		if err != nil {
			file.Close()
			return nil, err
		}
		return &BusStopReader{
			reader: file,
			info: FileInfo{
				ModTime: stat.ModTime(),
				IsOS:    true,
			},
		}, nil
	}

	// fallback: embedded
	embeddedFile, err := DataFS.Open("bus-stops.csv")
	if err != nil {
		return nil, err
	}
	return &BusStopReader{
		reader: embeddedFile,
		info: FileInfo{
			IsOS: false,
		},
	}, nil
}

func (r *BusStopReader) Close() error {
	return r.reader.Close()
}

func (r *BusStopReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *BusStopReader) Info() FileInfo {
	return r.info
}
