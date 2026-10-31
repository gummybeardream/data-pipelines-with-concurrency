package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"

	imageprocessing "github.com/gummybeardream/data-pipelines-with-concurrency/image_processing"
)

func TestNewJob(t *testing.T) {
	inputPath := filepath.Join("testdata", "images", "dog.jpeg")
	expectedOutput := filepath.Join(
		"testdata",
		"images",
		"output",
		"dog.jpeg",
	)

	job := newJob(inputPath)

	if job.InputPath != inputPath {
		t.Errorf(
			"newJob() input path = %q; want %q",
			job.InputPath,
			inputPath,
		)
	}

	if job.OutPath != expectedOutput {
		t.Errorf(
			"newJob() output path = %q; want %q",
			job.OutPath,
			expectedOutput,
		)
	}

	if job.Err != nil {
		t.Errorf("newJob() error = %v; want nil", job.Err)
	}

	if job.Image != nil {
		t.Error("newJob() image is not nil; want nil before loading")
	}
}

func TestProcessImagesInvalidMode(t *testing.T) {
	results, err := processImages("parallel", nil)

	if err == nil {
		t.Fatal("processImages() error = nil; want invalid-mode error")
	}

	if !strings.Contains(err.Error(), "invalid mode") {
		t.Errorf(
			"processImages() error = %q; want text containing %q",
			err.Error(),
			"invalid mode",
		)
	}

	if results != nil {
		t.Errorf("processImages() results = %v; want nil", results)
	}
}

func createTestImages(t *testing.T) []string {
	t.Helper()

	inputDir := filepath.Join(t.TempDir(), "images")
	paths := make([]string, 3)

	for index := range paths {
		path := filepath.Join(
			inputDir,
			fmt.Sprintf("image%d.jpeg", index+1),
		)

		img := image.NewRGBA(
			image.Rect(0, 0, 20+index, 10+index),
		)

		img.Set(
			0,
			0,
			color.RGBA{
				R: uint8(50 * (index + 1)),
				G: 100,
				B: 150,
				A: 255,
			},
		)

		if err := imageprocessing.WriteImage(path, img); err != nil {
			t.Fatalf("create test image %q: %v", path, err)
		}

		paths[index] = path
	}

	return paths
}

func TestProcessImagesMissingInput(t *testing.T) {
	modes := []string{"sequential", "concurrent"}

	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			missingPath := filepath.Join(
				t.TempDir(),
				"images",
				"missing.jpeg",
			)

			results, err := processImages(
				mode,
				[]string{missingPath},
			)
			if err != nil {
				t.Fatalf(
					"processImages(%q) mode error: %v",
					mode,
					err,
				)
			}

			if len(results) != 1 {
				t.Fatalf(
					"processImages(%q) returned %d results; want 1",
					mode,
					len(results),
				)
			}

			result := results[0]

			if result.Err == nil {
				t.Errorf(
					"processImages(%q) result error = nil; want missing-file error",
					mode,
				)
			}

			if _, err := os.Stat(result.OutPath); !os.IsNotExist(err) {
				t.Errorf(
					"processImages(%q) unexpectedly created output %q",
					mode,
					result.OutPath,
				)
			}
		})
	}
}

func TestProcessImagesValidModes(t *testing.T) {
	modes := []string{"sequential", "concurrent"}

	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			paths := createTestImages(t)

			results, err := processImages(mode, paths)
			if err != nil {
				t.Fatalf(
					"processImages(%q) returned an error: %v",
					mode,
					err,
				)
			}

			if len(results) != len(paths) {
				t.Fatalf(
					"processImages(%q) returned %d results; want %d",
					mode,
					len(results),
					len(paths),
				)
			}

			for _, result := range results {
				if result.Err != nil {
					t.Errorf(
						"processImages(%q) failed for %q: %v",
						mode,
						result.InputPath,
						result.Err,
					)
					continue
				}

				output, err := imageprocessing.ReadImage(result.OutPath)
				if err != nil {
					t.Errorf(
						"read output %q: %v",
						result.OutPath,
						err,
					)
					continue
				}

				bounds := output.Bounds()
				if bounds.Dx() != 500 || bounds.Dy() != 500 {
					t.Errorf(
						"output %q dimensions = %dx%d; want 500x500",
						result.OutPath,
						bounds.Dx(),
						bounds.Dy(),
					)
				}
			}
		})
	}
}
