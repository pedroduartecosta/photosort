package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {

	// TODO: Test is using multiple CPUs actually improves perfomance since  this is a sequential process
	runtime.GOMAXPROCS(runtime.NumCPU())

	startFolder := os.Args[1]
	archiveFolder := os.Args[2]

	var totalSize int64 = 0

	err := filepath.Walk(startFolder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			copiedBytes, err := processFile(path, archiveFolder)
			totalSize += copiedBytes
			fmt.Println(ByteCountSI(totalSize))

			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

//  Converts a size in bytes to a human-readable string in SI (decimal)
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func processFile(filePath string, archiveFolder string) (int64, error) {
	filename := filepath.Base(filePath)

	date, err := getDate(filePath)

	createDir(fmt.Sprintf("%s/%s", archiveFolder, "others"))

	if err != nil {
		extension := filepath.Ext(filePath)
		if isImageOrVideo(extension) {
			finalPath := fmt.Sprintf("%s/%s/%s", archiveFolder, "others", filename)
			return copy(filePath, finalPath)
		}
		return 0, fmt.Errorf("not a media file")
	}

	finalPath := newPath(archiveFolder, filename, date)

	return copy(filePath, finalPath)
}

func isImageOrVideo(extension string) bool {

	imageExtensions := map[string]bool{".tiff": true, ".tif": true, ".gif": true, ".jpeg": true, ".jpg": true, ".png": true, ".raw": true, ".webm": true, ".mkv": true, ".avi": true, ".mov": true, ".wmv": true, ".mp4": true, ".m4v": true, ".mpg": true, ".mp2": true, ".mpeg": true}

	return imageExtensions[extension]
}

func getDate(filepath string) (time.Time, error) {
	var dt time.Time
	file, err := os.Open(filepath)
	if err != nil {
		return dt, err
	}

	data, err := exif.Decode(file)
	if err != nil {
		return dt, err
	}

	return data.DateTime()
}

// Returns true if a dir/file already exists
func Exists(filepath string) bool {
	_, err := os.Stat(filepath)

	if err != nil && os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Generates the entire new path based on all the data, checks for collisions (and rename if needed)
func newPath(archive string, oldName string, date time.Time) string {
	dir := fmt.Sprintf("%s/%d/%d", archive, date.Year(), date.Month())
	createDir(dir)

	return fmt.Sprintf("%s/%d/%d/%s", archive, date.Year(), date.Month(), oldName)
}

// Creates a directory if it doesn't exist yet
func createDir(dir string) {
	if !Exists(dir) {
		os.MkdirAll(dir, 0644)
	}
}

func copy(src, dst string) (int64, error) {

	if Exists(dst) {
		return 0, fmt.Errorf("File already exists")
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
