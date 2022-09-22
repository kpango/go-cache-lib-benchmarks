package main

import (
	"math/rand"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/kpango/glg"
)

// generate random data for bar chart
func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}

func main() {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "Go's cache libraries benchmarking",
	}))

	// Put data into instance
	bar.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("SmallData", generateBarItems()).
		AddSeries("BigData", generateBarItems())
	// Where the magic happens
	f, err := os.Create("benchmark_result.html")
	if err != nil {
		glg.Fatal(err)
	}
	defer f.Close()
	bar.Render(f)
}
