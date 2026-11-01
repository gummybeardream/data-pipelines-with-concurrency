package main

import "testing"

var benchmarkImagePaths = []string{
	"images/corgi.jpeg",
	"images/bernesemountain.jpeg",
	"images/icelandicsheep.jpeg",
	"images/norwegianelkhound.jpeg",
}

func runPipelineBenchmark(
	b *testing.B,
	processor func([]string) []Result,
) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()

	for iteration := 0; iteration < b.N; iteration++ {
		results := processor(benchmarkImagePaths)

		if len(results) != len(benchmarkImagePaths) {
			b.Fatalf(
				"pipeline returned %d results; want %d",
				len(results),
				len(benchmarkImagePaths),
			)
		}

		for _, result := range results {
			if result.Err != nil {
				b.Fatalf(
					"pipeline failed for %q: %v",
					result.InputPath,
					result.Err,
				)
			}
		}
	}
}

func BenchmarkProcessSequential(b *testing.B) {
	runPipelineBenchmark(b, processSequential)
}

func BenchmarkProcessConcurrent(b *testing.B) {
	runPipelineBenchmark(b, processConcurrent)
}
