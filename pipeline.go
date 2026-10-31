package main

import (
	"fmt"
	"image"
	"strings"

	imageprocessing "github.com/gummybeardream/data-pipelines-with-concurrency/image_processing"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
	Err       error
}

type Result struct {
	InputPath string
	OutPath   string
	Err       error
}

// newJob creates a job with matching input and output paths
func newJob(inputPath string) Job {
	return Job{
		InputPath: inputPath,
		OutPath: strings.Replace(
			inputPath,
			"images/",
			"images/output/",
			1,
		),
	}
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		// Create one job per input path and send it to the next pipeline stage
		for _, path := range paths {
			job := newJob(path)

			img, err := imageprocessing.ReadImage(path)
			job.Image = img
			job.Err = err
			out <- job
		}
		close(out)
	}()
	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		// Resize successful jobs and pass every job to the next pipeline stage
		for job := range input { 
			if job.Err == nil {
				job.Image = imageprocessing.Resize(job.Image)
			}

			out <- job
		}

		close(out)
	}()

	return out
}

// convertToGrayscale converts successful jobs and passes all jobs onward
func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			if job.Err == nil {
				job.Image = imageprocessing.Grayscale(job.Image)
			}
			out <- job
		}
		close(out)
	}()
	return out
}

// saveImage writes successful jobs and returns the result of every job
func saveImage(input <-chan Job) <-chan Result {
	out := make(chan Result)
	go func() {
		for job := range input {
			if job.Err == nil {
				job.Err = imageprocessing.WriteImage(job.OutPath, job.Image)
			}

			out <- Result{
				InputPath: job.InputPath,
				OutPath:   job.OutPath,
				Err:       job.Err,
			}

		}
		close(out)

	}()
	return out
}

// processConcurrent runs the four goroutine-based pipeline stages
func processConcurrent(paths []string) []Result {
	channel1 := loadImage(paths)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	writeResults := saveImage(channel3)

	var results []Result

	for result := range writeResults {
		results = append(results, result)
	}

	return results
}

// processSequential processes each image without starting goroutines
func processSequential(paths []string) []Result {
	results := make([]Result, 0, len(paths))

	for _, path := range paths {
		job := newJob(path)

		job.Image, job.Err = imageprocessing.ReadImage(job.InputPath)

		if job.Err == nil {
			job.Image = imageprocessing.Resize(job.Image)
		}

		if job.Err == nil {
			job.Image = imageprocessing.Grayscale(job.Image)
		}

		if job.Err == nil {
			job.Err = imageprocessing.WriteImage(job.OutPath, job.Image)
		}

		results = append(results, Result{
			InputPath: job.InputPath,
			OutPath:   job.OutPath,
			Err:       job.Err,
		})
	}

	return results
}

func processImages(mode string, paths []string) ([]Result, error) {
	switch mode {
	case "concurrent":
		return processConcurrent(paths), nil
	case "sequential":
		return processSequential(paths), nil
	default:
		return nil, fmt.Errorf(
			"invalid mode %q: use concurrent or sequential",
			mode,
		)
	}
}
