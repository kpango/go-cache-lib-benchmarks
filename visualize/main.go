package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
	fileName := flag.String("f", "../README.md", "Filename to read benchmark results from")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", *fileName, err)
		return
	}
	defer file.Close()

	var results []BenchmarkResult

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
		fmt.Printf("Error reading %s: %v\n", *fileName, err)
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

	err = os.MkdirAll("../images", 0755)
	if err != nil {
		fmt.Printf("Error creating images directory: %v\n", err)
		return
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
	outPath := filepath.Join("../images", filename)

	f, err := os.Create(outPath)
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
