package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/montanaflynn/stats"
	"gonum.org/v1/gonum/stat/distuv"
)

//let's make a application that the user can decide the # of bootstraps and estimates standard error of the median from a normal distribution

// function for mean
// function for median
// function for

// Define a struct to hold the results of each iteration
type IterationResult struct {
	N          int
	SampleMean float64
	BootMean   float64
	BootMedian float64
}

func sampleWithReplacement(slice []float64, sampleSize int) []float64 {
	sampled := make([]float64, sampleSize)
	for i := 0; i < sampleSize; i++ {
		sampled[i] = slice[rand.Intn(len(slice))]
	}
	return sampled
}

type BootstrapResult struct {
	BootMean   float64
	BootMedian float64
}

func rnorm(n int, mean, stddev float64) []float64 {
	// Create a normal distribution with specified mean and standard deviation
	dist := distuv.Normal{
		Mu:    mean,
		Sigma: stddev,
	}

	// Slice to hold the generated random numbers
	samples := make([]float64, n)

	// Generate n random samples
	for i := 0; i < n; i++ {
		samples[i] = dist.Rand()
	}

	return samples
}

func filterResults(results []IterationResult, n int) []IterationResult {
	var filtered []IterationResult
	for _, result := range results { // Iterate through the slice
		if result.N == n { // Check if the N field matches the desired value
			filtered = append(filtered, result) // Add to the filtered slice
		}
	}
	return filtered
}

func Bootstrap(B int) {
	//run 100 iterations
	studySampleSizes := []int{25, 100, 225, 400}
	popMean := 100
	popSD := 10
	var studyResults []IterationResult

	for i := 0; i < 100; i++ {
		//iterate over sample sizes
		for _, index := range studySampleSizes {
			//generate a sample of size n from a normal distribution

			//thisSample generates a random sample of size n from a normal distribution with mean popMean and standard deviation popSD.
			thisSample := rnorm(index, float64(popMean), float64(popSD))

			thisSampleMean, _ := stats.Mean(thisSample)

			var bootstrapSampleResults []BootstrapResult

			//boostrap
			for b := 1; b <= B; b++ {
				thisBootstrapSample := sampleWithReplacement(thisSample, index)
				thisBootstrapMean, _ := stats.Mean(thisBootstrapSample)
				thisBootstrapMedian, _ := stats.Median(thisBootstrapSample)

				// Create a struct for the current bootstrap result
				thisBootstrapSampleResult := BootstrapResult{
					BootMean:   thisBootstrapMean,
					BootMedian: thisBootstrapMedian,
				}
				bootstrapSampleResults = append(bootstrapSampleResults, thisBootstrapSampleResult)
			}

			var bootMeans []float64
			var bootMedians []float64
			for _, result := range bootstrapSampleResults {
				bootMeans = append(bootMeans, result.BootMean)
				bootMedians = append(bootMedians, result.BootMedian)
			}
			bootMean, _ := stats.Mean(bootMeans)
			bootMedian, _ := stats.Mean(bootMedians)

			// Create a new IterationResult struct to hold the results of this iteration
			thisIterationResults := IterationResult{
				N:          index,
				SampleMean: thisSampleMean,
				BootMean:   bootMean,
				BootMedian: bootMedian,
			}
			// Append the results to the studyResults slice
			studyResults = append(studyResults, thisIterationResults)
		}
	}
	for _, index := range studySampleSizes {
		var finalsamplemeans []float64
		var finalbootmeans []float64
		var finalbootmedian []float64

		filterResults := filterResults(studyResults, index)

		for _, result := range filterResults {
			finalsamplemeans = append(finalsamplemeans, result.SampleMean)
			finalbootmeans = append(finalbootmeans, result.BootMean)
			finalbootmedian = append(finalbootmedian, result.BootMedian)
		}
		x, _ := stats.StandardDeviation(finalsamplemeans)
		y, _ := stats.StandardDeviation(finalbootmeans)
		z, _ := stats.StandardDeviation(finalbootmedian)

		fmt.Printf("\nSamples of size n = %d\n", index)
		fmt.Printf("  SE Mean from Central Limit Theorem: %.2f\n", float64(popSD)/math.Pow(float64(index), 0.5))
		fmt.Printf("  SE Mean from Samples: %.2f\n", x)
		fmt.Printf("  SE Mean from Bootstrap Samples: %.2f\n", y)
		fmt.Printf("  SE Median from Bootstrap Samples: %.2f\n", z)
	}
}

func main() {
	printTotalAlloc() // Print total memory allocation at the start

	// Simulate memory allocation
	_ = make([]byte, 50<<20) // Allocate 50 MB

	start := time.Now()
	fmt.Println("Enter number of boostraps")
	var input int
	fmt.Scanln(&input)
	if input > 1000 {
		fmt.Println("Please enter a number smaller than 500.")
	} else {
		Bootstrap(input)
	}
	duration := time.Since(start)
	fmt.Println(duration)
	printTotalAlloc() // Print total memory allocation after allocation
}

func printTotalAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Total memory allocated (in bytes)
	fmt.Printf("TotalAlloc = %v MiB\n", bToMb(m.TotalAlloc))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
