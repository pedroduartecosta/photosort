package main

import (
	"flag"
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
	// TODO: Test if using multiple CPUs actually improves performance since  this is a sequential process
	runtime.GOMAXPROCS(runtime.NumCPU())

	var sourceFolder, destinationFolder string
	flag.StringVar(&sourceFolder, "source-folder", "", "Source folder with photos")
	flag.StringVar(&destinationFolder, "destination-folder", "", "Destination folder with archived sorted photos")
	flag.Parse()

	if sourceFolder == "" {
		log.Fatal("source must be specified")
	}
	if destinationFolder == "" {
		log.Fatal("archive must be specified")
	}

	var totalSize int64 = 0

	err := filepath.Walk(sourceFolder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			copiedBytes, err := processFile(path, destinationFolder)
			if err != nil {
				// we won't return nil just to not quit the walk function
				log.Println(err)
			} else {
				totalSize += copiedBytes
				log.Println(ByteCountSI(totalSize))
			}

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
	if err != nil {
		if err := createDir(fmt.Sprintf("%s/%s", archiveFolder, "others")); err != nil {
			return 0, err
		}
		extension := filepath.Ext(filePath)
		if isImageOrVideo(extension) {
			finalPath := fmt.Sprintf("%s/%s/%s", archiveFolder, "others", filename)
			return copyFile(filePath, finalPath)
		}
		return 0, &ErrNotMediaFile{extension, filePath}
	}

	finalPath, err := newPath(archiveFolder, filename, date)
	if err != nil {
		return 0, err
	}

	return copyFile(filePath, finalPath)
}

func isImageOrVideo(extension string) bool {

	imageExtensions := map[string]bool{".tiff": true, ".tif": true, ".gif": true, ".jpeg": true, ".jpg": true, ".png": true, ".raw": true, ".webm": true, ".mkv": true, ".avi": true, ".mov": true, ".wmv": true, ".mp4": true, ".MP4": true, ".m4v": true, ".mpg": true, ".mp2": true, ".mpeg": true}

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
func Exists(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// Generates the entire new path based on all the data, checks for collisions (and rename if needed)
func newPath(archive string, oldName string, date time.Time) (string, error) {
	dir := fmt.Sprintf("%s/%d/%d", archive, date.Year(), date.Month())
	if err := createDir(dir); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%d/%d/%s", archive, date.Year(), date.Month(), oldName), nil
}

// Creates a directory if it doesn't exist yet
func createDir(dir string) error {
	if exists, err := Exists(dir); err != nil {
		return err
	} else if !exists {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) (int64, error) {

	if exists, err := Exists(dst); err != nil {
		return 0, err
	} else if exists {
		return 0, &ErrFileExists{filePath: dst}
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, &ErrFileNotRegular{filePath: src}
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func() { _ = source.Close() }()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() { _ = destination.Close() }()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
