package main

import "fmt"

type ErrNotMediaFile struct{ extension, filePath string }

func (e *ErrNotMediaFile) Error() string {
	return fmt.Sprintf("extension `%s` is not a media file (file: %s)", e.extension, e.filePath)
}

type ErrFileExists struct{ filePath string }

func (e *ErrFileExists) Error() string {
	return fmt.Sprintf("file `%s` already exists", e.filePath)
}

type ErrFileNotRegular struct{ filePath string }

func (e *ErrFileNotRegular) Error() string {
	return fmt.Sprintf("file `%s` is not a regular file", e.filePath)
}
