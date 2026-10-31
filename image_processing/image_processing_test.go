package imageprocessing

import (
	"image"
	"image/color"
	"os"
	"testing"

	"path/filepath"
)

func TestResize(t *testing.T) {
	input := image.NewRGBA(image.Rect(0, 0, 100, 50))

	output := Resize(input)
	bounds := output.Bounds()

	if bounds.Dx() != 500 {
		t.Errorf("Resize() width = %d; want 500", bounds.Dx())
	}

	if bounds.Dy() != 500 {
		t.Errorf("Resize() height = %d; want 500", bounds.Dy())
	}
}

func TestGrayscale(t *testing.T) {
	input := image.NewRGBA(image.Rect(0, 0, 2, 2))
	input.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	output := Grayscale(input)

	if output.Bounds() != input.Bounds() {
		t.Errorf(
			"Grayscale() bounds = %v; want %v",
			output.Bounds(),
			input.Bounds(),
		)
	}

	if output.ColorModel() != color.GrayModel {
		t.Errorf(
			"Grayscale() color model = %v; want grayscale",
			output.ColorModel(),
		)
	}
}

func TestWriteAndReadImage(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "nested", "test.jpeg")

	input := image.NewRGBA(image.Rect(0, 0, 20, 10))

	if err := WriteImage(outputPath, input); err != nil {
		t.Fatalf("WriteImage() returned an error: %v", err)
	}

	output, err := ReadImage(outputPath)
	if err != nil {
		t.Fatalf("ReadImage() returned an error: %v", err)
	}

	if output.Bounds().Dx() != 20 || output.Bounds().Dy() != 10 {
		t.Errorf(
			"written image dimensions = %dx%d; want 20x10",
			output.Bounds().Dx(),
			output.Bounds().Dy(),
		)
	}
}

func TestReadImageMissingFile(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "missing.jpeg")

	_, err := ReadImage(missingPath)
	if err == nil {
		t.Fatal("ReadImage() error = nil; want a missing-file error")
	}
}

func TestReadImageInvalidData(t *testing.T) {
	invalidPath := filepath.Join(t.TempDir(), "invalid.jpeg")

	err := os.WriteFile(
		invalidPath,
		[]byte("this is not an image"),
		0600,
	)
	if err != nil {
		t.Fatalf("create invalid test file: %v", err)
	}

	_, err = ReadImage(invalidPath)
	if err == nil {
		t.Fatal("ReadImage() error = nil; want a decoding error")
	}
}

func TestWriteImageInvalidOutputPath(t *testing.T) {
	tempDir := t.TempDir()
	parentFile := filepath.Join(tempDir, "not-a-directory")

	if err := os.WriteFile(parentFile, []byte("file"), 0600); err != nil {
		t.Fatalf("create parent file: %v", err)
	}

	outputPath := filepath.Join(parentFile, "output.jpeg")
	input := image.NewRGBA(image.Rect(0, 0, 10, 10))

	err := WriteImage(outputPath, input)
	if err == nil {
		t.Fatal("WriteImage() error = nil; want an invalid-path error")
	}
}
