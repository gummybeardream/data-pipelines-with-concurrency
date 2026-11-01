# Data Pipelines with Concurrency

Requires Go 1.20 or higher.

## Management Problem

A technology start-up would like to reduce the processing time associated with data engineering pipelines. This project will replicate a pipeline to evaluate whether Go concurrency can reduce the elapsed time required to process a collection of images.

## Project Overview

The application reads four images, resizes them to 500 x 500 pixels, converts them to grayscale, and writes the processed images to an output directory. It supports a sequential mode and a concurrent pipeline mode so their performance can be compared under the same workload.

## Original Project Attribution

This project is based on Amrit Singh's CODEHEIM image-processing pipeline
example:

- [CODEHEIM Episode 21](https://www.codeheim.io/courses/Episode-21-Concurrency-in-Go-Pipeline-Pattern-65c3ca14e4b0628a4e002201)
- [Original GitHub repository](https://github.com/code-heim/go_21_goroutines_pipeline)

## Pipeline Design

The pipeline performs four stages:

1. Read each input image.
2. Resize it to 500 × 500 pixels.
3. Convert it to grayscale.
4. Write it to the output directory.

### Concurrent Processing

Concurrent mode runs the four pipeline stages in separate goroutines. Channels pass image jobs and errors between stages, allowing different images to occupy different stages at the same time.

### Sequential Processing

Sequential mode completes all four stages for one image before starting the next image. It uses the same image-processing helper functions but does not start pipeline goroutines.

## Input Images

Three source images were converted from WebP to JPEG using macOS Preview at approximately 90% quality. The conversion was necessary because the original pipeline did not support WebP decoding. The Norwegian Elkhound image was created by Sannse and is used under CC BY-SA 3.0. The filename was changed for use in this project. The image remains licensed under CC BY-SA 3.0 and is not relicensed under Apache 2.0.

| Local filename | Source | Creator | License | Original format | Dimensions |
|---|---|---|---|---|---|
| `bernesemountain.jpeg` | [Kevin The Berner](https://www.kevintheberner.com/blog/6-reasons-you-need-a-bernese-mountain-dog) | Kevin The Berner website; photographer not stated | Not stated | WebP | 2500 × 3124 |
| `corgi.jpeg` | [Modern Dog Magazine](https://moderndogmagazine.com/articles/breeds/the-pembroke-welsh-corgi/) | Liudmila Bohush; attribution inferred from source filename | Not stated | WebP | 819 × 1024 |
| `icelandicsheep.jpeg` | [Rover: Dogs That Look Like Foxes](https://www.rover.com/blog/dogs-that-look-like-foxes/) | Veronica Druk, via Wikimedia Commons | [CC BY-SA 3.0](https://creativecommons.org/licenses/by-sa/3.0/) | WebP | 700 × 700 |
| `norwegianelkhound.jpeg` | [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Norwegian_Elkhound.jpg) | Sannse | [CC BY-SA 3.0](https://creativecommons.org/licenses/by-sa/3.0/) | JPEG | 600 × 629 |

## Error Handling

### Input Error Handling

The image-loading helper function returns errors that indicate what has been affected instead of stopping with `panic`. Errors can occur when an input file cannot be opened, decoded, or closed. Each error displays the affected file path and is carried through the pipeline without attempting later image transformations.

The program continues processing other valid images and reports every result. If any input fails, the program exits with status 1 after all jobs have been handled. Manual checks confirmed that missing files and files containing invalid images produce descriptive errors without causing a panic.

### Output Error Handling

The image-writing helper function creates the output directory and returns the errors instead of stopping the program with `panic`. Errors can occur while creating the directory, creating the output file, encoding the JPEG image, or closing the file. Each error includes the affected output path.

Output errors are structured similarly to the input error reports. The program reports every result and exits with status 1 if any output fails. A manual check using an invalid output directory confirmed that failures produce descriptive errors without causing a panic.

## How to Run the Image Processor

It will run in default concurrent mode with this command:

```bash
go run .
```

You can also use this command to display `-mode` options and their descriptions:

```bash
go run . -help
```

A manual test with an invalid mode was used to show that the program will prompt to call either concurrent or sequential and end with exit status 2.

### Run with goroutines

```bash
go run . -mode=concurrent
```

### Run without goroutines

```bash
go run . -mode=sequential
```

## Testing

The project includes unit tests for image-processing helpers and complete pipeline behavior. Synthetic images and temporary directories keep the tests independent of the assignment input files.

The tests cover:

- Image resizing and grayscale conversion
- Successful image reading and writing
- Missing and invalid input files
- Invalid output paths
- Job input and output paths
- Invalid processing modes
- Missing inputs in both processing modes
- Successful sequential and concurrent processing
- Output image dimensions

Run the complete test suite:

```bash
go test ./...
```

Run the tests with Go's race detector:

```bash
go test -race ./...
```

Generate statement coverage:

```bash
go test -cover ./...
```

The final test run achieved **75.0% statement coverage** for the main pipeline package and **91.3% statement coverage** for the image-processing package. The race detector reported no data races.

## Performance Evaluation

### Test Environment

| Component | Value |
|---|---|
| Operating system | macOS 26.6.2 |
| Go version | Go 1.27.1 |
| Architecture | ARM64 |
| Model identifier | MacBookPro17,1 |
| Processor | Apple M1 |
| Memory | 16 GiB |

### Original Baseline

The original pipeline was compiled before runtime was measured so that compilation time was excluded. Using the `time` command, the program was run 6 times with the first time being a warm-up and the next five runs were measured. The average of the `total` runtime was 0.1492 seconds or 149.2 milliseconds.

This result represents the original concurrent implementation and is retained as a historical baseline. It is not a controlled comparison with the final results because the replacement images and implementation differ.

| Run | User time | System time | CPU usage | Elapsed (`total`) time |
|---:|---:|---:|---:|---:|
| 1 | 0.43 s | 0.02 s | 270% | 0.164 s |
| 2 | 0.44 s | 0.01 s | 320% | 0.141 s |
| 3 | 0.45 s | 0.02 s | 317% | 0.145 s |
| 4 | 0.44 s | 0.02 s | 299% | 0.154 s |
| 5 | 0.44 s | 0.01 s | 321% | 0.142 s |
| **Mean** | **0.440 s** | **0.016 s** | **305.4%** | **0.1492 s** |

### Benchmark Methodology

The benchmarks processed the same four assignment images through the complete pipeline: reading, resizing, grayscale conversion and output writing. The sequential and concurrent benchmarks used identical inputs and image-processing helpers. Repeated iterations overwrote the same generated output files.

Benchmarks were run on macOS with an Apple M1 processor using Go's ARM64 architecture. Each benchmark processed the same four input images through the pipeline. Five samples were collected for each mode, with each sample running for approximately three seconds.

The benchmark command was:

```bash
go test -run='^$' -bench='BenchmarkProcess' -benchmem -benchtime=3s -count=5 .
```

Go reported elapsed nanoseconds, allocated bytes, and memory allocations per pipeline run. The raw output is saved in `results/benchmark_results.txt`.

### Benchmark Results

| Mode | Samples | Mean time/run | Median time/run | Mean throughput | Mean bytes/run | Mean allocations/run |
|---|---:|---:|---:|---:|---:|---:|
| Sequential | 5 | 183.975 ms | 183.965 ms | 21.74 images/sec | 61,797,420 B | 1,000,558 |
| Concurrent | 5 | 149.500 ms | 149.300 ms | 26.76 images/sec | 61,797,714 B | 1,000,568 |

The concurrent pipeline ran **1.23 times faster** than the sequential pipeline, reducing mean elapsed time by **18.7%**. Its estimated throughput increased from 21.74 to 26.76 images per second. Memory use and allocation counts were nearly identical between the two modes.

- [Raw benchmark results](results/benchmark_results.txt)

### Final Timing Methodology

The final application was compiled once before measurement so compilation time was excluded. Both modes processed the same four input images and performed the same reading, resizing, grayscale conversion, and output-writing work.

Each mode was run once as an unrecorded warm-up. Five paired sequential and concurrent runs were then measured using the zsh `time` command. The elapsed `total` time was used for comparison.

### Final Executable Timing Results

| Mode | Runs | Mean elapsed time | Median elapsed time | Range | Mean throughput | Mean CPU usage |
|---|---:|---:|---:|---:|---:|---:|
| Sequential | 5 | 0.1976 s | 0.198 s | 0.193–0.202 s | 20.24 images/sec | 197.2% |
| Concurrent | 5 | 0.1532 s | 0.153 s | 0.151–0.156 s | 26.11 images/sec | 254.8% |

The concurrent pipeline completed the four-image workload approximately **1.29 times faster** than the sequential pipeline. Its mean elapsed time was approximately **22.5%** lower, and estimated throughput increased from 20.24 to 26.11 images per second.

Concurrent processing used more CPU capacity, averaging 254.8% compared with 197.2% for sequential processing. The additional CPU used allowed pipeline stages to make progress at the same time, reducing the total time required to complete the workload.

- [Raw executable timing results](results/manual_timing_results.txt)

### Performance Summary

| Measurement method | Sequential | Concurrent | Speedup | Time reduction |
|---|---:|---:|---:|---:|
| Go benchmark | 183.975 ms | 149.500 ms | 1.23× | 18.7% |
| Built executable | 197.600 ms | 153.200 ms | 1.29× | 22.5% |

Both measurement methods found that concurrent processing reduced elapsed time. The Go benchmarks measured repeated pipeline operations inside the testing framework, while the executable measurements included application startup and terminal output.

### Interpretation

The results support the management hypothesis for this four-image workload. Concurrent processing reduced mean executable elapsed time by approximately 22.5%, while the Go benchmarks measured an 18.7% reduction. Concurrent mode used more CPU capacity but completed the workload sooner, indicating that the goroutine-based pipeline made better use of available processing resources.

These findings apply to this computer, these four images, and this pipeline design. Larger and more varied workloads should be tested before assuming that the same performance improvement will occur in other data-engineering
applications.

## Building the Executable

Build the macOS executable from the repository root:

```bash
go build -o image-pipeline .
```

Run concurrent mode:

```bash
./image-pipeline -mode=concurrent
```

Run sequential mode:

```bash
./image-pipeline -mode=sequential
```

The executable is a command-line application rather than a macOS `.app`
bundle.

To cross-compile a Windows executable from macOS:

```bash
GOOS=windows GOARCH=amd64 go build -o image-pipeline.exe .
```

The Windows executable must be run on Windows:

```powershell
.\image-pipeline.exe -mode=concurrent
```

## GenAI Tools

I used OpenAI Codex as a programming tutor and review assistant during this
project. It helped me:

- Break the assignment into manageable tasks.
- Understand Go errors, goroutines, channels, and command-line flags.
- Diagnose compiler errors and runtime panics.
- Design input and output error handling.
- Plan unit tests and benchmarks.
- Interpret coverage, benchmark, and executable timing results.
- Review code organization and README documentation.

I reviewed, entered, tested, and revised the project code and documentation.
AI-generated suggestions were validated with Go formatting, tests, race
detection, vetting, benchmarks, and manual application runs.

## License

The project source code is licensed under the Apache License 2.0. Third-party
input images remain governed by the individual licenses listed in the Input
Images table and are not relicensed under Apache 2.0.
