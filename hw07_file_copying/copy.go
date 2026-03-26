package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileOpen              = errors.New("error while opening the file")
	ErrReadStat              = errors.New("error while reading file's stat")
	ErrFileSeek              = errors.New("error while seeking")
	ErrFileWrite             = errors.New("error writing file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return ErrFileOpen
	}
	defer closeFile(fromFile)

	fromFileInfo, err := fromFile.Stat()
	if err != nil {
		return ErrReadStat
	}
	if !fromFileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if offset >= fromFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return ErrFileOpen
	}
	defer closeFile(toFile)

	toFileInfo, err := toFile.Stat()
	if err != nil {
		return ErrReadStat
	}
	if !toFileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return ErrFileSeek
	}

	total := limit
	if total == 0 {
		total = fromFileInfo.Size()
	}
	bar := pb.Full.Start64(total)
	barReader := bar.NewProxyReader(fromFile)
	defer bar.Finish()

	if limit == 0 {
		_, err = io.Copy(toFile, barReader)
	} else {
		_, err = io.CopyN(toFile, barReader, limit)
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return ErrFileWrite
	}

	return nil
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Println(err)
	}
}
