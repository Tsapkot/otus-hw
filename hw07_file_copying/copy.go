package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileInfo              = errors.New("failed to get file info")
	ErrLimit                 = errors.New("limit might not be lower then 0")
	ErrFileOverlap           = errors.New("source and destination files might not be the same")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Stat(fromPath)
	if err != nil {
		return ErrFileInfo
	}
	dstFile, err := os.Stat(toPath)
	if os.SameFile(srcFile, dstFile) {
		return ErrFileOverlap
	}

	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("source file error: %w", err)
	}
	defer fromFile.Close()

	fromInfo, err := fromFile.Stat()
	if err != nil {
		return ErrFileInfo
	}

	if !fromInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fromInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit < 0 {
		return ErrLimit
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer toFile.Close()

	if _, err := fromFile.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek input file: %w", err)
	}

	bytesToCopy := fromInfo.Size() - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	bar := pb.Full.Start64(bytesToCopy)
	barReader := bar.NewProxyReader(fromFile)
	defer bar.Finish()

	_, err = io.CopyN(toFile, barReader, bytesToCopy)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	return nil
}
