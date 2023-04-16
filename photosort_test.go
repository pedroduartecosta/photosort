package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIsImageOrVideo(t *testing.T) {
	tests := []struct {
		extension string
		expected  bool
	}{
		{".jpg", true},
		{".mp4", true},
		{".txt", false},
		{".JPG", true},
		{".MP4", true},
	}

	for _, test := range tests {
		result := isImageOrVideo(test.extension)
		if result != test.expected {
			t.Errorf("isImageOrVideo(%s) returned %t, expected %t", test.extension, result, test.expected)
		}
	}
}

func TestCreateDir(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "photosort_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testDir := filepath.Join(tempDir, "testdir")

	err = createDir(testDir)
	if err != nil {
		t.Errorf("createDir(%s) returned error: %v", testDir, err)
	}

	exists, err := Exists(testDir)
	if err != nil {
		t.Errorf("Exists(%s) returned error: %v", testDir, err)
	}

	if !exists {
		t.Errorf("Directory %s was not created", testDir)
	}
}

func createTestImageWithExifDate(t *testing.T, tempDir string, date time.Time) string {
	// You may use an actual image file with EXIF data or mock the exif.Decode function to return the desired date

	imgPath := filepath.Join(tempDir, "test_image.jpg")
	err := ioutil.WriteFile(imgPath, []byte("fake_image_data"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	return imgPath
}

func TestGetDate(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "photosort_test_get_date")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	currentDate := time.Now()
	expectedDate := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), currentDate.Hour(), currentDate.Minute(), 0, 0, time.Local)
	testImagePath := createTestImageWithExifDate(t, tempDir, expectedDate)

	date, err := getDate(testImagePath)
	if err != nil {
		t.Errorf("getDate(%s) returned error: %v", testImagePath, err)
	}

	// Truncate date to remove any extra precision
	date = date.UTC().Truncate(time.Minute)
	currentDate = currentDate.UTC().Truncate(time.Minute)

	if !date.Equal(expectedDate) {
		t.Errorf("getDate(%s) returned date %v, expected %v", testImagePath, date, expectedDate)
	}
}

func TestNewPath(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "photosort_test_new_path")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	date := time.Date(2022, 10, 10, 0, 0, 0, 0, time.UTC)
	fileName := "test_image.jpg"

	newPath, err := newPath(tempDir, fileName, date)
	if err != nil {
		t.Errorf("newPath(%s, %s, %v) returned error: %v", tempDir, fileName, date, err)
	}

	expectedPath := filepath.Join(tempDir, "2022", "10", fileName)
	if newPath != expectedPath {
		t.Errorf("newPath(%s, %s, %v) returned path %s, expected %s", tempDir, fileName, date, newPath, expectedPath)
	}
}
