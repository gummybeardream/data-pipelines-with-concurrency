# Data Pipelines with Concurrency

[Episode link](https://www.codeheim.io/courses/Episode-21-Concurrency-in-Go-Pipeline-Pattern-65c3ca14e4b0628a4e002201)

Requires Golang 1.20 or higher.

## Input Images

Three source images were converted from WebP to JPEG using macOS Preview at approximately 90% quality. The conversion was necessary because the original pipeline did not support WebP decoding. The Norwegian Elkhound image was created by Sannse and is used under CC BY-SA 3.0. The filename was changed for use in this project. The image remains licensed under CC BY-SA 3.0 and is not relicensed under Apache 2.0.

| Local filename | Source | Creator | License | Original format | Dimensions |
|---|---|---|---|---|---|
| `bernesemountain.jpeg` | [Kevin The Berner](https://www.kevintheberner.com/blog/6-reasons-you-need-a-bernese-mountain-dog) | Kevin The Berner website; photographer not stated | Not stated | WebP | 2500 × 3124 |
| `corgi.jpeg` | [Modern Dog Magazine](https://moderndogmagazine.com/articles/breeds/the-pembroke-welsh-corgi/) | Liudmila Bohush; attribution inferred from source filename | Not stated | WebP | 819 × 1024 |
| `icelandicsheep.jpeg` | [Rover: Dogs That Look Like Foxes](https://www.rover.com/blog/dogs-that-look-like-foxes/) | Veronica Druk, via Wikimedia Commons | [CC BY-SA 3.0](https://creativecommons.org/licenses/by-sa/3.0/) | WebP | 700 × 700 |
| `norwegianelkhound.jpeg` | [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Norwegian_Elkhound.jpg) | Sannse | [CC BY-SA 3.0](https://creativecommons.org/licenses/by-sa/3.0/) | JPEG | 600 × 629 |


## Input Error Handling

The image-loading helper function returns errors that indicate what has been affected instead of stopping with `panic`. Errors can occur when an input file cannot be opened, decoded, or closed. Each error displays the affected file path and is carried through the pipeline without attempting later image transformations.

The program continues processing other valid images and reports every result. If any input fails, the program exits with status 1 after all jobs have been handled. Manual checks confirmed that missing files and files containing invalid images produce descriptive errors without causing a panic.

## Output Error Handling

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


## Image-Processing Helper Tests

Six unit tests verify image resizing, grayscale conversion, successful image reading and writing, missing input files, invalid image data, and invalid output paths. The tests use synthetic images and temporary directories so they do not depend on the assignment input files.

The helper package achieved **91.3% statement coverage**. Running the complete test suite with Go's race detector reported no data races.

## Results

### Original Baseline Results

The original pipeline was compiled before runtime was measured so that compilation time was excluded. Using the `time` command, the program was run 6 times with the first time being a warm-up and the next five runs were measured. The average of the `total` runtime was .1492 seconds or 149.2 milliseconds.

This result represents the original concurrent implementation and is used only as a baseline. The final comparison will measure sequential and concurrent processing with the same four replacement images.

| Run | User time | System time | CPU usage | Elapsed (`total`) time |
|---:|---:|---:|---:|---:|
| 1 | 0.43 s | 0.02 s | 270% | 0.164 s |
| 2 | 0.44 s | 0.01 s | 320% | 0.141 s |
| 3 | 0.45 s | 0.02 s | 317% | 0.145 s |
| 4 | 0.44 s | 0.02 s | 299% | 0.154 s |
| 5 | 0.44 s | 0.01 s | 321% | 0.142 s |
| **Mean** | **0.440 s** | **0.016 s** | **305.4%** | **0.1492 s** |

## Future Comparison
| Version | Mode | Measured runs | Mean elapsed time |
|---|---|---:|---:|
| Original | Concurrent | 5 | 0.1492 s (149.2 ms) |
| Final | Sequential | TBD | TBD |
| Final | Concurrent | TBD | TBD |


## Benchmark Methodology
The benchmarks processed the same four assignment images through the complete pipeline: reading, resizing, grayscale conversion and output writing. The sequential and concurrent benchmarks used identical inputs and image-processing helpers. Repeated iterations overwrote the same generated output files.

Benchmarks were run on macOS with an Apple M1 processor using Go's ARM64 architecture. Each benchmark processed the same four input images through the pipeline. Five samples were collected for each mode, with each sample running for approximately three seconds.

The benchmark command was:

```bash
go test -run='^$' -bench='BenchmarkProcess' -benchmem -benchtime=3s -count=5 .
```

Go reported elapsed nanoseconds, allocated bytes, and memory allocations per pipeline run. The raw output is saved in `results/benchmark_results.txt`.

## Benchmark Results

| Mode | Samples | Mean time/run | Median time/run | Mean throughput | Mean bytes/run | Mean allocations/run |
|---|---:|---:|---:|---:|---:|---:|
| Sequential | 5 | 183.975 ms | 183.965 ms | 21.74 images/sec | 61,797,420 B | 1,000,558 |
| Concurrent | 5 | 149.500 ms | 149.300 ms | 26.76 images/sec | 61,797,714 B | 1,000,568 |

The concurrent pipeline ran **1.23 times faster** than the sequential pipeline, reducing mean elapsed time by **18.7%**. Its estimated throughput increased from 21.74 to 26.76 images per second. Memory use and allocation counts were nearly identical between the two modes.