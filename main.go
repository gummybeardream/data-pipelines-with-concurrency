package main

import (
	"fmt"
	imageprocessing "github.com/gummybeardream/data-pipelines-with-concurrency/image_processing"
	"image"
	"os"
	"strings"
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

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		// For each input path create a job and add it to
		// the out channel
		for _, p := range paths {
			job := Job{InputPath: p,
				OutPath: strings.Replace(p, "images/", "images/output/", 1)}

			img, err := imageprocessing.ReadImage(p)
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
		// For each input job, create a new job after resize and add it to
		// the out channel
		for job := range input { // Read from the channel
			if job.Err == nil {
				job.Image = imageprocessing.Resize(job.Image)
			}

			out <- job
		}

		close(out)
	}()

	return out
}

func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input { // Read from the channel
			if job.Err == nil {
				job.Image = imageprocessing.Grayscale(job.Image)
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(input <-chan Job) <-chan Result {
	out := make(chan Result)
	go func() {
		for job := range input { // Read from the channel
			if job.Err == nil {
				imageprocessing.WriteImage(job.OutPath, job.Image)
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

func main() {

	imagePaths := []string{"images/corgi.jpeg",
		"images/bernesemountain.jpeg",
		"images/icelandicsheep.jpeg",
		"images/norwegianelkhound.jpeg",
	}

	channel1 := loadImage(imagePaths)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	writeResults := saveImage(channel3)

	hadError := false

	for result := range writeResults {
		if result.Err != nil {
			hadError = true
			fmt.Fprintf(os.Stderr, "Failed: %v\n", result.Err)
			continue
		}

		fmt.Printf("Success: %s\n", result.OutPath)
	}
	if hadError {
		os.Exit(1)
	}
}
