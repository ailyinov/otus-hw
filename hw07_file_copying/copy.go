package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

const buffSize = 64

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	s, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if s.Size() == 0 || s.IsDir() {
		return ErrUnsupportedFile
	}
	if offset > s.Size() {
		return ErrOffsetExceedsFileSize
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

	buf := make([]byte, buffSize)
	writtenCnt := 0

	pbStart := limit
	if pbStart == 0 || (limit+offset) > s.Size() {
		pbStart = s.Size() - offset
	}
	bar := pb.Start64(pbStart)
	defer func() {
		bar.Finish()
	}()

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
		bar.Add(n)
		if int(limit) > 0 && writtenCnt == int(limit) {
			return nil
		}
	}

	return nil
}
