package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// BenchmarkResult holds parsed data from a single benchmark line.
type BenchmarkResult struct {
	LibName      string
	DataPattern  string
	TTL          string
	ProcessCount int
	NsPerOp      float64
	BytesPerOp   float64
	AllocsPerOp  float64
}

// GroupKey identifies a (DataPattern, TTL) combination.
type GroupKey struct {
	DataPattern string
	TTL         string
}

// colorPalette holds full hex colors used as fallback for unknown libraries.
var colorPalette = []string{
	"#1f77b4", "#ff7f0e", "#2ca02c", "#d62728", "#9467bd",
	"#8c564b", "#e377c2", "#7f7f7f", "#bcbd22", "#17becf",
}

// libColorMap assigns a fixed color to each known library so that every chart
// uses the same color for the same library, regardless of which subset of
// libraries appears in a given chart group.
var libColorMap = map[string]string{
	"BigCache":   "#1f77b4", // Blue
	"DefaultMap": "#ff7f0e", // Orange
	"Gache":      "#2ca02c", // Green
	"GacheV2":    "#d62728", // Red
	"GCacheARC":  "#9467bd", // Purple
	"GCacheLFU":  "#8c564b", // Brown
	"GCacheLRU":  "#e377c2", // Pink
	"GoCache":    "#7f7f7f", // Gray
	"SyncMap":    "#bcbd22", // Yellow-green
	"TTLCache":   "#17becf", // Cyan
}

// libColor returns the fixed color for a library name.  If the library is not
// in the static map it falls back to the palette (cycled by index).
func libColor(name string, index int) string {
	if c, ok := libColorMap[name]; ok {
		return c
	}
	return colorPalette[index%len(colorPalette)]
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
			parallelism, _ := strconv.Atoi(matches[4])
			results = append(results, BenchmarkResult{
				LibName:      matches[1],
				DataPattern:  matches[2],
				TTL:          matches[3],
				ProcessCount: parallelism,
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

	// Group results by (DataPattern, TTL); parallelism levels become the x-axis.
	grouped := make(map[GroupKey][]BenchmarkResult)
	for _, res := range results {
		key := GroupKey{DataPattern: res.DataPattern, TTL: res.TTL}
		grouped[key] = append(grouped[key], res)
	}

	if err := os.MkdirAll("../images", 0755); err != nil {
		fmt.Printf("Error creating images directory: %v\n", err)
		return
	}

	var groupKeys []GroupKey
	for key, group := range grouped {
		groupKeys = append(groupKeys, key)
		// Interactive HTML 3D chart (go-echarts / ECharts GL).
		generateTrue3DChart(key, group)
		// Static SVG 3D chart (oblique projection) — for README embedding.
		generate3DSVGChart(key, group)
	}
	fmt.Println("Charts generated successfully.")

	// Update README.md with embedded chart images.
	updateREADME(*fileName, groupKeys)
}

// sortedLibs returns a deterministically sorted list of unique library names.
func sortedLibs(results []BenchmarkResult) []string {
	libSet := make(map[string]bool)
	for _, r := range results {
		libSet[r.LibName] = true
	}
	libs := make([]string, 0, len(libSet))
	for lib := range libSet {
		libs = append(libs, lib)
	}
	sort.Strings(libs)
	return libs
}

// sortedProcessCounts returns a sorted list of unique process count values.
func sortedProcessCounts(results []BenchmarkResult) []int {
	pSet := make(map[int]bool)
	for _, r := range results {
		pSet[r.ProcessCount] = true
	}
	ps := make([]int, 0, len(pSet))
	for p := range pSet {
		ps = append(ps, p)
	}
	sort.Ints(ps)
	return ps
}

// generateTrue3DChart creates an interactive HTML 3D line chart using go-echarts (ECharts GL):
//   - X-axis = Process Count (category)
//   - Y-axis = Latency/Op (ns/op)
//   - Z-axis = Alloc/Op
//
// Each cache library is a separate 3D line series, colour-coded.
// The chart is fully interactive: rotate, zoom, and tooltip on hover.
func generateTrue3DChart(key GroupKey, results []BenchmarkResult) {
	processCounts := sortedProcessCounts(results)
	libs := sortedLibs(results)
	if len(processCounts) == 0 || len(libs) == 0 {
		return
	}

	xLabels := make([]string, len(processCounts))
	for i, p := range processCounts {
		xLabels[i] = fmt.Sprintf("P%d", p)
	}

	line3d := charts.NewLine3D()
	line3d.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "980px", Height: "700px"}),
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("%s %s — ProcessCount × Latency/Op × Alloc/Op", key.DataPattern, key.TTL),
		}),
		charts.WithGrid3DOpts(opts.Grid3D{
			BoxWidth:  200,
			BoxHeight: 100,
			BoxDepth:  80,
			ViewControl: &opts.ViewControl{
				AutoRotate:      opts.Bool(true),
				AutoRotateSpeed: 10,
			},
		}),
		charts.WithXAxis3DOpts(opts.XAxis3D{
			Name: "Process Count",
			Type: "category",
			Data: xLabels,
		}),
		charts.WithYAxis3DOpts(opts.YAxis3D{
			Name: "Latency/Op (ns/op)",
			Type: "value",
		}),
		charts.WithZAxis3DOpts(opts.ZAxis3D{
			Name: "Alloc/Op",
			Type: "value",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
	)

	for i, lib := range libs {
		color := libColor(lib, i)
		var data3d []opts.Chart3DData
		for j, p := range processCounts {
			for _, r := range results {
				if r.LibName == lib && r.ProcessCount == p {
					data3d = append(data3d, opts.Chart3DData{
						Value: []interface{}{xLabels[j], r.NsPerOp, r.AllocsPerOp},
						Name:  fmt.Sprintf("%s P%d", lib, p),
					})
					break
				}
			}
		}
		if len(data3d) == 0 {
			continue
		}
		line3d.AddSeries(lib, data3d,
			charts.WithItemStyleOpts(opts.ItemStyle{Color: color}),
		)
	}

	outPath := filepath.Join("../images", fmt.Sprintf("%s_%s_3d_chart.html", key.DataPattern, key.TTL))
	f, err := os.Create(outPath)
	if err != nil {
		fmt.Printf("Error creating 3D chart %s: %v\n", outPath, err)
		return
	}
	defer f.Close()
	if err := line3d.Render(f); err != nil {
		fmt.Printf("Error rendering 3D chart %s: %v\n", outPath, err)
	}
	fmt.Printf("Generated: %s\n", outPath)
}

// fmtVal formats a float64 compactly with K/M suffixes for SVG labels.
func fmtVal(v float64) string {
	if v >= 1e6 {
		return fmt.Sprintf("%.1fM", v/1e6)
	}
	if v >= 1e3 {
		return fmt.Sprintf("%.1fK", v/1e3)
	}
	if v >= 10 {
		return fmt.Sprintf("%.0f", v)
	}
	return fmt.Sprintf("%.1f", v)
}

// generate3DSVGChart creates a static SVG using oblique 3D projection:
//   - x-axis = Process Count
//   - y-axis = Latency/Op (ns/op)
//   - z-axis = Alloc/Op
//
// Each cache library is drawn as a colour-coded 3D polyline with drop lines,
// back-wall and floor projections, and per-point value labels.
// The SVG is suitable for embedding in README.md on GitHub.
func generate3DSVGChart(key GroupKey, results []BenchmarkResult) {
	processCounts := sortedProcessCounts(results)
	libs := sortedLibs(results)
	if len(processCounts) == 0 || len(libs) == 0 {
		return
	}

	// Value ranges.
	minNs, maxNs := math.MaxFloat64, -math.MaxFloat64
	minAlloc, maxAlloc := math.MaxFloat64, -math.MaxFloat64
	for _, r := range results {
		if r.NsPerOp < minNs {
			minNs = r.NsPerOp
		}
		if r.NsPerOp > maxNs {
			maxNs = r.NsPerOp
		}
		if r.AllocsPerOp < minAlloc {
			minAlloc = r.AllocsPerOp
		}
		if r.AllocsPerOp > maxAlloc {
			maxAlloc = r.AllocsPerOp
		}
	}
	if maxNs <= minNs {
		maxNs = minNs + 1
	}
	if maxAlloc <= minAlloc {
		maxAlloc = minAlloc + 1
	}

	// Oblique projection geometry (30° angle).
	xLen := 340.0 // ProcessCount axis pixel length
	yLen := 310.0 // Latency axis pixel length
	zLen := 210.0 // Alloc axis projected pixel length
	phi := math.Pi / 6
	ox, oy := 210.0, 555.0 // 3D origin in screen coordinates
	cosPhi := math.Cos(phi)
	sinPhi := math.Sin(phi)

	// project converts normalised [0,1]³ coords to 2D screen coords.
	project := func(nx, ny, nz float64) (float64, float64) {
		return ox + nx*xLen + nz*zLen*cosPhi,
			oy - ny*yLen - nz*zLen*sinPhi
	}

	nPC := len(processCounts)
	pcIdx := make(map[int]int, nPC)
	for i, p := range processCounts {
		pcIdx[p] = i
	}
	normX := func(i int) float64 {
		if nPC <= 1 {
			return 0
		}
		return float64(i) / float64(nPC-1)
	}
	normY := func(v float64) float64 { return (v - minNs) / (maxNs - minNs) }
	normZ := func(v float64) float64 { return (v - minAlloc) / (maxAlloc - minAlloc) }

	outPath := filepath.Join("../images", fmt.Sprintf("%s_%s_3d_chart.svg", key.DataPattern, key.TTL))
	f, err := os.Create(outPath)
	if err != nil {
		fmt.Printf("Error creating 3D SVG chart %s: %v\n", outPath, err)
		return
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	defer bw.Flush()
	wr := func(s string, a ...interface{}) { fmt.Fprintf(bw, s+"\n", a...) }

	wr(`<?xml version="1.0" encoding="UTF-8"?>`)
	wr(`<svg xmlns="http://www.w3.org/2000/svg" width="980" height="700">`)
	wr(`  <rect width="980" height="700" fill="white"/>`)

	// Title.
	titleX := ox + xLen/2 + zLen*cosPhi*0.5
	wr(`  <text x="%.0f" y="36" text-anchor="middle" font-family="sans-serif" font-size="15" font-weight="bold">%s %s — ProcessCount × Latency/Op × Alloc/Op</text>`,
		titleX, key.DataPattern, key.TTL)

	// ── Walls and floor (drawn back-to-front) ───────────────────────────────

	// Back wall: XY plane at z=1.
	{
		ax, ay := project(0, 0, 1)
		bx, by := project(1, 0, 1)
		cx, cy := project(1, 1, 1)
		dx, dy := project(0, 1, 1)
		wr(`  <polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="#eef5ee" stroke="#aaa" stroke-width="0.8"/>`,
			ax, ay, bx, by, cx, cy, dx, dy)
	}
	// Left wall: YZ plane at x=0.
	{
		ax, ay := project(0, 0, 0)
		bx, by := project(0, 0, 1)
		cx, cy := project(0, 1, 1)
		dx, dy := project(0, 1, 0)
		wr(`  <polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="#eeeff8" stroke="#aaa" stroke-width="0.8"/>`,
			ax, ay, bx, by, cx, cy, dx, dy)
	}
	// Floor: XZ plane at y=0.
	{
		ax, ay := project(0, 0, 0)
		bx, by := project(1, 0, 0)
		cx, cy := project(1, 0, 1)
		dx, dy := project(0, 0, 1)
		wr(`  <polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="#f5f5ee" stroke="#aaa" stroke-width="0.8"/>`,
			ax, ay, bx, by, cx, cy, dx, dy)
	}

	// ── Grid lines ──────────────────────────────────────────────────────────

	const (
		nGrid  = 4
		nZTick = 3
	)
	// Floor grid.
	for i := 0; i <= nGrid; i++ {
		t := float64(i) / float64(nGrid)
		x0, y0 := project(0, 0, t)
		x1, y1 := project(1, 0, t)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#ccc" stroke-width="0.7" stroke-dasharray="4,3"/>`, x0, y0, x1, y1)
		x0, y0 = project(t, 0, 0)
		x1, y1 = project(t, 0, 1)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#ccc" stroke-width="0.7" stroke-dasharray="4,3"/>`, x0, y0, x1, y1)
	}
	// Left wall grid.
	for i := 1; i <= nGrid; i++ {
		t := float64(i) / float64(nGrid)
		x0, y0 := project(0, t, 0)
		x1, y1 := project(0, t, 1)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#ccc" stroke-width="0.7" stroke-dasharray="4,3"/>`, x0, y0, x1, y1)
		x0, y0 = project(0, 0, t)
		x1, y1 = project(0, 1, t)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#ccc" stroke-width="0.7" stroke-dasharray="4,3"/>`, x0, y0, x1, y1)
	}
	// Z-axis grid lines.
	for i := 1; i < nZTick; i++ {
		t := float64(i) / float64(nZTick)
		x0, y0 := project(0, 0, t)
		x1, y1 := project(1, 0, t)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#b0b8d0" stroke-width="1.1" stroke-dasharray="5,3"/>`, x0, y0, x1, y1)
		x0, y0 = project(0, 0, t)
		x1, y1 = project(0, 1, t)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#b0b8d0" stroke-width="1.1" stroke-dasharray="5,3"/>`, x0, y0, x1, y1)
	}

	// ── Axes ────────────────────────────────────────────────────────────────

	// X axis (ProcessCount).
	{
		x0, y0 := project(0, 0, 0)
		x1, y1 := project(1.05, 0, 0)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#333" stroke-width="2"/>`, x0, y0, x1, y1)
	}
	// Y axis (Latency).
	{
		x0, y0 := project(0, 0, 0)
		x1, y1 := project(0, 1.05, 0)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#333" stroke-width="2"/>`, x0, y0, x1, y1)
	}
	// Z axis (Alloc) — dashed.
	{
		x0, y0 := project(0, 0, 0)
		x1, y1 := project(0, 0, 1.05)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#333" stroke-width="2" stroke-dasharray="6,3"/>`, x0, y0, x1, y1)
	}

	// ── X axis ticks and labels ──────────────────────────────────────────────

	for i, p := range processCounts {
		nx := normX(i)
		tx, ty := project(nx, 0, 0)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#555"/>`, tx, ty, tx, ty+5)
		wr(`  <text x="%.1f" y="%.1f" text-anchor="middle" font-family="sans-serif" font-size="12">P%d</text>`, tx, ty+20, p)
	}
	{
		lx, ly := project(0.5, 0, 0)
		wr(`  <text x="%.1f" y="%.1f" text-anchor="middle" font-family="sans-serif" font-size="13">Process Count</text>`, lx, ly+38)
	}

	// ── Y axis ticks and label ───────────────────────────────────────────────

	const nYTick = 4
	for i := 0; i <= nYTick; i++ {
		t := float64(i) / float64(nYTick)
		val := minNs + t*(maxNs-minNs)
		tx, ty := project(0, t, 0)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#555"/>`, tx-5, ty, tx, ty)
		wr(`  <text x="%.1f" y="%.1f" text-anchor="end" font-family="sans-serif" font-size="11" dominant-baseline="middle">%s</text>`, tx-8, ty, fmtVal(val))
	}
	{
		lx, ly := project(0, 0.5, 0)
		wr(`  <text transform="rotate(-90,%.0f,%.0f)" x="%.0f" y="%.0f" text-anchor="middle" font-family="sans-serif" font-size="13">Latency/Op (ns/op)</text>`,
			lx-60, ly, lx-60, ly)
	}

	// ── Z axis ticks and label ───────────────────────────────────────────────

	for i := 0; i <= nZTick; i++ {
		t := float64(i) / float64(nZTick)
		val := minAlloc + t*(maxAlloc-minAlloc)
		tx, ty := project(0, 0, t)
		wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#555"/>`, tx-4, ty+2, tx, ty)
		wr(`  <text x="%.1f" y="%.1f" text-anchor="end" font-family="sans-serif" font-size="11" dominant-baseline="middle">%s</text>`, tx-7, ty+3, fmtVal(val))
	}
	{
		lx, ly := project(0, 0, 0.5)
		wr(`  <text x="%.1f" y="%.1f" text-anchor="middle" font-family="sans-serif" font-size="13">Alloc/Op</text>`, lx, ly+22)
	}

	// ── Per-library 3D data (painter's algorithm: back → front) ─────────────

	type pt3 struct {
		nx, ny, nz float64
		p          int
		ns, alloc  float64
	}
	type libEntry struct {
		name  string
		color string
		avgNZ float64
		pts   []pt3
	}

	var entries []libEntry
	for i, lib := range libs {
		colorHex := libColor(lib, i)
		var e libEntry
		e.name = lib
		e.color = colorHex
		var sumNZ float64
		for j, p := range processCounts {
			for _, r := range results {
				if r.LibName == lib && r.ProcessCount == p {
					nz := normZ(r.AllocsPerOp)
					e.pts = append(e.pts, pt3{normX(j), normY(r.NsPerOp), nz, p, r.NsPerOp, r.AllocsPerOp})
					sumNZ += nz
					break
				}
			}
		}
		if len(e.pts) > 0 {
			e.avgNZ = sumNZ / float64(len(e.pts))
			entries = append(entries, e)
		}
	}
	// Highest avgNZ = furthest back = drawn first.
	sort.Slice(entries, func(i, j int) bool { return entries[i].avgNZ > entries[j].avgNZ })

	// Back-wall projections (Latency/Op vs Process Count at z=1).
	for _, e := range entries {
		for k := 1; k < len(e.pts); k++ {
			x0, y0 := project(e.pts[k-1].nx, e.pts[k-1].ny, 1)
			x1, y1 := project(e.pts[k].nx, e.pts[k].ny, 1)
			wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1.8" stroke-dasharray="10,5" opacity="0.65"/>`,
				x0, y0, x1, y1, e.color)
		}
	}
	// Floor projections (Alloc/Op vs Process Count at y=0).
	for _, e := range entries {
		for k := 1; k < len(e.pts); k++ {
			x0, y0 := project(e.pts[k-1].nx, 0, e.pts[k-1].nz)
			x1, y1 := project(e.pts[k].nx, 0, e.pts[k].nz)
			wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1.8" stroke-dasharray="5,4" opacity="0.65"/>`,
				x0, y0, x1, y1, e.color)
		}
	}
	// Drop lines to floor.
	for _, e := range entries {
		for _, pt := range e.pts {
			sx, sy := project(pt.nx, pt.ny, pt.nz)
			fx, fy := project(pt.nx, 0, pt.nz)
			wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="0.9" stroke-dasharray="3,2" opacity="0.5"/>`,
				sx, sy, fx, fy, e.color)
		}
	}
	// 3D polylines (dashed).
	for _, e := range entries {
		for k := 1; k < len(e.pts); k++ {
			x0, y0 := project(e.pts[k-1].nx, e.pts[k-1].ny, e.pts[k-1].nz)
			x1, y1 := project(e.pts[k].nx, e.pts[k].ny, e.pts[k].nz)
			wr(`  <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="2.5" stroke-dasharray="8,4"/>`, x0, y0, x1, y1, e.color)
		}
	}
	// Data point markers + value labels.
	for _, e := range entries {
		for _, pt := range e.pts {
			sx, sy := project(pt.nx, pt.ny, pt.nz)
			wr(`  <circle cx="%.1f" cy="%.1f" r="5" fill="%s" stroke="white" stroke-width="1.5"><title>%s P%d — %s ns/op, %s allocs/op</title></circle>`,
				sx, sy, e.color, e.name, pt.p, fmtVal(pt.ns), fmtVal(pt.alloc))
			wr(`  <text x="%.1f" y="%.1f" text-anchor="middle" font-family="sans-serif" font-size="9" fill="%s" stroke="white" stroke-width="3" paint-order="stroke">%s ns/op</text>`,
				sx, sy-14, e.color, fmtVal(pt.ns))
			wr(`  <text x="%.1f" y="%.1f" text-anchor="middle" font-family="sans-serif" font-size="9" fill="%s" stroke="white" stroke-width="3" paint-order="stroke">%s a/op</text>`,
				sx, sy+20, e.color, fmtVal(pt.alloc))
		}
	}

	// ── Legend ──────────────────────────────────────────────────────────────

	legendX := int(ox+xLen+zLen*cosPhi) + 20
	// Dash-style key.
	wr(`  <text x="%d" y="60" font-family="sans-serif" font-size="12" font-weight="bold">Dash style</text>`, legendX)
	wr(`  <line x1="%d" y1="76" x2="%d" y2="76" stroke="#555" stroke-width="2" stroke-dasharray="8,4"/>`, legendX, legendX+22)
	wr(`  <text x="%d" y="76" font-family="sans-serif" font-size="11" dominant-baseline="middle">3D path</text>`, legendX+28)
	wr(`  <line x1="%d" y1="94" x2="%d" y2="94" stroke="#555" stroke-width="1.8" stroke-dasharray="10,5"/>`, legendX, legendX+22)
	wr(`  <text x="%d" y="94" font-family="sans-serif" font-size="11" dominant-baseline="middle">Latency/Op (back wall)</text>`, legendX+28)
	wr(`  <line x1="%d" y1="112" x2="%d" y2="112" stroke="#555" stroke-width="1.8" stroke-dasharray="5,4"/>`, legendX, legendX+22)
	wr(`  <text x="%d" y="112" font-family="sans-serif" font-size="11" dominant-baseline="middle">Alloc/Op (floor)</text>`, legendX+28)
	// Per-library colour entries.
	wr(`  <text x="%d" y="136" font-family="sans-serif" font-size="13" font-weight="bold">Library</text>`, legendX)
	sort.Slice(entries, func(i, j int) bool { return entries[i].name < entries[j].name })
	for i, e := range entries {
		ly := 156 + i*22
		wr(`  <line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2.5" stroke-dasharray="8,4"/>`,
			legendX, ly, legendX+18, ly, e.color)
		wr(`  <circle cx="%d" cy="%d" r="4" fill="%s"/>`, legendX+9, ly, e.color)
		wr(`  <text x="%d" y="%d" font-family="sans-serif" font-size="12" dominant-baseline="middle">%s</text>`,
			legendX+24, ly, e.name)
	}

	wr(`</svg>`)
	fmt.Printf("Generated: %s\n", outPath)
}

// chartSectionMarkerStart is the marker that delimits the auto-generated chart
// section in README.md.
const chartSectionMarkerStart = "<!-- benchmark-chart-section-start -->"

// chartSectionMarkerEnd is the closing marker for the auto-generated section.
const chartSectionMarkerEnd = "<!-- benchmark-chart-section-end -->"

// htmlPreviewBaseURL is the GitHub Pages base URL used to serve interactive HTML charts.
// The custom domain kpango.com is served via Cloudflare and proxies to kpango.github.io.
const htmlPreviewBaseURL = "https://kpango.com/go-cache-lib-benchmarks/"

// updateREADME inserts (or replaces) an auto-generated chart section at the end
// of the README.md file.  The section embeds all generated SVG charts and adds
// preview links to the interactive HTML versions.
func updateREADME(readmePath string, keys []GroupKey) {
	// Sort keys for deterministic order: BigData before SmallData, NoTTL before WithTTL.
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].DataPattern != keys[j].DataPattern {
			return keys[i].DataPattern < keys[j].DataPattern
		}
		return keys[i].TTL < keys[j].TTL
	})

	data, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Printf("Error reading README for update: %v\n", err)
		return
	}
	content := string(data)

	// Build the new chart section.
	var sb strings.Builder
	sb.WriteString(chartSectionMarkerStart + "\n")
	sb.WriteString("\n## Benchmark Charts\n\n")
	for _, key := range keys {
		svgFile := fmt.Sprintf("images/%s_%s_3d_chart.svg", key.DataPattern, key.TTL)
		htmlFile := fmt.Sprintf("%s_%s_3d_chart.html", key.DataPattern, key.TTL)
		title := fmt.Sprintf("%s %s", key.DataPattern, key.TTL)
		sb.WriteString(fmt.Sprintf("### %s\n\n", title))
		sb.WriteString(fmt.Sprintf("![%s](%s)\n\n", title, svgFile))
		sb.WriteString(fmt.Sprintf("[📊 View Interactive 3D Chart](%s%s)\n\n", htmlPreviewBaseURL, htmlFile))
	}
	sb.WriteString(chartSectionMarkerEnd + "\n")

	newSection := sb.String()

	// Replace existing section or append.
	startIdx := strings.Index(content, chartSectionMarkerStart)
	endIdx := strings.Index(content, chartSectionMarkerEnd)
	if startIdx >= 0 && endIdx >= 0 {
		tailStart := endIdx + len(chartSectionMarkerEnd)
		if tailStart < len(content) && content[tailStart] == '\n' {
			tailStart++
		}
		content = content[:startIdx] + newSection + content[tailStart:]
	} else {
		content = strings.TrimRight(content, "\n") + "\n\n" + newSection
	}

	if err := os.WriteFile(readmePath, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing README: %v\n", err)
		return
	}
	fmt.Printf("Updated: %s\n", readmePath)
}
