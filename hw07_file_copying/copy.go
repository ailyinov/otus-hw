package main

import (
	"errors"
	"io"
	"math"
	"os"
)

const buffSize = 64

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	if s, err := os.Stat(fromPath); err != nil {
		return err
	} else {
		if s.Size() == 0 || s.IsDir() {
			return ErrUnsupportedFile
		}
		if offset > s.Size() {
			return ErrOffsetExceedsFileSize
		}
	}
	source, err := os.Open(fromPath)
	defer func() {
		err = source.Close()
	}()
	if err != nil {
		return err
	}

	dest, err := os.Create(toPath)
	defer func() {
		err = dest.Close()
	}()
	if err != nil {
		return err
	}

	_, err = source.Seek(offset, 0)
	if err != nil {
		return err
	}
	buf := make([]byte, int(math.Min(float64(buffSize), float64(limit))))
	writtenCnt := 0
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if cnt := writtenCnt + n; cnt > int(limit) && limit > 0 {
			n = int(limit) - writtenCnt
		}

		if _, err := dest.Write(buf[:n]); err != nil {
			return err
		}
		writtenCnt += n
		if int(limit) > 0 && writtenCnt == int(limit) {
			return nil
		}
	}

	return nil
}
