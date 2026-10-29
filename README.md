# Data Pipelines with Concurrency 

[Episode link](https://www.codeheim.io/courses/Episode-21-Concurrency-in-Go-Pipeline-Pattern-65c3ca14e4b0628a4e002201)

Requires Golang 1.20 or higher.


## Results 

### Original Baseline 

### Original Baseline Results

The original pipeline was compiled before runtime was measurement so that compliation time was excluded. Using the `time` command, the program was run 6 times with the first time being a warm-up and the next five runs were measured. The average of the `total` runtime was .1492 seconds or 149.2 millisecones. 

This result represents the original concurrent implementation and used only as a baseline. The final comparison will measure sequential and concurrent processing with the same four replacement images. 

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