package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/wcharczuk/go-chart/v2"
)

type BenchmarkResult struct {
	LibName      string
	DataPattern  string
	TTL          string
	Parallelism  string
	NsPerOp      float64
	BytesPerOp   float64
	AllocsPerOp  float64
}

type GroupKey struct {
	DataPattern string
	TTL         string
	Parallelism string
}

func main() {
	file, err := os.Open("../README.md")
	if err != nil {
		fmt.Println("Error opening README.md:", err)
		return
	}
	defer file.Close()

	var results []BenchmarkResult

	// Benchmark{Cache Lib Name}SetGet{Data Pattern}{TTL on/off}/P{Parallelism}-{CPU Cores}
	// Example: BenchmarkDefaultMapSetGetSmallDataNoTTL/P100-128
	// Let's use more generic non-greedy capture groups.
	// 1: LibName (.+?)
	// 2: DataPattern (.+?)
	// 3: TTL (NoTTL|WithTTL|.*TTL) -> actually just (.+?)
	// It's a bit tricky because both DataPattern and TTL are contiguous without delimiters.
	// But usually it ends with "NoTTL" or "WithTTL".
	// Let's assume TTL ends with "TTL". So (.*)(NoTTL|WithTTL|.*TTL)
	// Or we can just capture everything between SetGet and /P as one group, or split it if we know the suffixes.
	// Since the user asked for {Data Pattern}{TTL on/off}, let's try:
	// ^Benchmark(.+?)SetGet(.+?)(NoTTL|WithTTL|OffTTL|OnTTL|TTL.*)/P(\d+)-\d+\s+\d+\s+([\d.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op
	// Wait, we don't know the exact TTL string. We can just use (.+?)(NoTTL|WithTTL|On|Off|Yes|No) if they follow a pattern, but "NoTTL|WithTTL" is standard here.
	// Let's use (.+?)(NoTTL|WithTTL) to be safe with the current data, OR we can capture everything and not split them.
	// "Benchの命名には規則があり、 Benchmark{Cache Lib Name}SetGet{Data Pattern}{TTL on/off}/P{Parallelism}-{CPU Cores}"
	// If the TTL is arbitrary, (.+?) might not know where DataPattern ends.
	// Let's capture `(.+?)` for DataPattern and `((?:With|No)TTL)` or similar. Let's look at the data again.
	// The data has SmallDataNoTTL, SmallDataWithTTL, BigDataNoTTL, BigDataWithTTL.
	re := regexp.MustCompile(`^Benchmark(.+?)SetGet(.+?)(NoTTL|WithTTL)/P(\d+)-\d+\s+\d+\s+([\d.]+)\s+ns/op\s+([\d.]+)\s+B/op\s+([\d.]+)\s+allocs/op`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 8 {
			nsOp, _ := strconv.ParseFloat(matches[5], 64)
			bOp, _ := strconv.ParseFloat(matches[6], 64)
			allocsOp, _ := strconv.ParseFloat(matches[7], 64)

			results = append(results, BenchmarkResult{
				LibName:     matches[1],
				DataPattern: matches[2],
				TTL:         matches[3],
				Parallelism: matches[4],
				NsPerOp:     nsOp,
				BytesPerOp:  bOp,
				AllocsPerOp: allocsOp,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading README.md:", err)
		return
	}

	grouped := make(map[GroupKey][]BenchmarkResult)
	for _, res := range results {
		key := GroupKey{
			DataPattern: res.DataPattern,
			TTL:         res.TTL,
			Parallelism: res.Parallelism,
		}
		grouped[key] = append(grouped[key], res)
	}

	for key, group := range grouped {
		generateCharts(key, group)
	}
	fmt.Println("Charts generated successfully.")
}

func generateCharts(key GroupKey, results []BenchmarkResult) {
	generateChart(key, results, "NsPerOp", func(r BenchmarkResult) float64 { return r.NsPerOp })
	generateChart(key, results, "BytesPerOp", func(r BenchmarkResult) float64 { return r.BytesPerOp })
	generateChart(key, results, "AllocsPerOp", func(r BenchmarkResult) float64 { return r.AllocsPerOp })
}

func generateChart(key GroupKey, results []BenchmarkResult, metricName string, getValue func(BenchmarkResult) float64) {
	sortedResults := make([]BenchmarkResult, len(results))
	copy(sortedResults, results)
	sort.Slice(sortedResults, func(i, j int) bool {
		return getValue(sortedResults[i]) < getValue(sortedResults[j])
	})

	var values []chart.Value
	for _, res := range sortedResults {
		values = append(values, chart.Value{
			Value: getValue(res),
			Label: res.LibName,
		})
	}

	// Calculate a suitable width based on the number of bars
	width := 200 + len(values)*80

	barChart := chart.BarChart{
		Title: fmt.Sprintf("%s %s P%s - %s", key.DataPattern, key.TTL, key.Parallelism, metricName),
		Background: chart.Style{
			Padding: chart.Box{
				Top:    40,
				Left:   20,
				Right:  20,
				Bottom: 20,
			},
		},
		Height:   600,
		Width:    width,
		BarWidth: 50,
		Bars:     values,
	}

	filename := fmt.Sprintf("%s_%s_P%s_%s.svg", key.DataPattern, key.TTL, key.Parallelism, metricName)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	err = barChart.Render(chart.SVG, f)
	if err != nil {
		fmt.Println("Error rendering chart:", err)
	}
}
