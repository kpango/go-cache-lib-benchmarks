package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// rankingChartFile returns the file name (relative to images/) of the ranking chart.
func rankingChartFile(key GroupKey) string {
	return fmt.Sprintf("%s_%s_ranking.svg", key.DataPattern, key.TTL)
}

// rankLabel renders the multiple of v relative to best (the lowest value in the
// column), e.g. "12x". The best entry is labelled "best".
func rankLabel(v, best float64) string {
	switch {
	case v <= best:
		return "best"
	case best <= 0:
		return "best is 0"
	case v/best >= 10:
		return fmt.Sprintf("%.0fx", v/best)
	}
	return fmt.Sprintf("%.1fx", v/best)
}

// logBarLen maps v onto [0, width] on a log10 scale spanning [lo, hi] (log10 units).
func logBarLen(v, lo, hi, width float64) float64 {
	if v <= 0 || hi <= lo {
		return 2
	}
	return math.Max(2, (math.Log10(v)-lo)/(hi-lo)*width)
}

// generateRankingSVGChart draws, for each parallelism level, horizontal log-scale
// bars of latency and allocations ranked best-first, annotated with the multiple
// relative to the best entry of each column. A log axis keeps libraries that
// differ by orders of magnitude on one chart while the labels state the real
// multiple; no library is special-cased.
func generateRankingSVGChart(key GroupKey, results []BenchmarkResult) {
	byLibP := make(map[string]map[int]BenchmarkResult)
	for _, r := range results {
		if byLibP[r.LibName] == nil {
			byLibP[r.LibName] = make(map[int]BenchmarkResult)
		}
		byLibP[r.LibName][r.ProcessCount] = r
	}
	if len(byLibP) == 0 {
		return
	}
	ps := sortedProcessCounts(results)

	// Shared log ranges so panels are visually comparable.
	loNs, hiNs := math.Inf(1), math.Inf(-1)
	loAl, hiAl := math.Inf(1), math.Inf(-1)
	for _, r := range results {
		if r.NsPerOp > 0 {
			loNs, hiNs = math.Min(loNs, r.NsPerOp), math.Max(hiNs, r.NsPerOp)
		}
		if r.AllocsPerOp > 0 {
			loAl, hiAl = math.Min(loAl, r.AllocsPerOp), math.Max(hiAl, r.AllocsPerOp)
		}
	}
	loNs, hiNs = math.Log10(loNs)-0.3, math.Log10(hiNs)
	loAl, hiAl = math.Log10(loAl)-0.3, math.Log10(hiAl)

	const (
		w        = 1160.0
		nameX    = 120.0
		col1X    = 135.0
		col2X    = 645.0
		barW     = 270.0
		rowH     = 26.0
		headerH  = 90.0
		panelGap = 22.0
	)
	libCount := len(byLibP)
	panelH := 30 + float64(libCount)*rowH
	h := headerH + float64(len(ps))*(panelH+panelGap) + 30

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f" font-family="sans-serif">`+"\n", w, h, w, h)
	fmt.Fprintf(&sb, `<rect width="%.0f" height="%.0f" fill="#ffffff"/>`+"\n", w, h)
	fmt.Fprintf(&sb, `<text x="%.0f" y="32" text-anchor="middle" font-size="20" font-weight="bold" fill="#222">%s %s — ranking by latency</text>`+"\n", w/2, key.DataPattern, key.TTL)
	fmt.Fprintf(&sb, `<text x="%.0f" y="54" text-anchor="middle" font-size="12" fill="#555">log10 scale, shorter bar = better. Labels show the multiple relative to the best in each column.</text>`+"\n", w/2)
	fmt.Fprintf(&sb, `<text x="%.0f" y="80" font-size="13" font-weight="bold" fill="#222">Latency (ns/op)</text>`+"\n", col1X)
	fmt.Fprintf(&sb, `<text x="%.0f" y="80" font-size="13" font-weight="bold" fill="#222">Allocations (allocs/op)</text>`+"\n", col2X)

	y := headerH
	for _, p := range ps {
		type row struct {
			lib string
			r   BenchmarkResult
		}
		var rows []row
		for lib, m := range byLibP {
			if r, ok := m[p]; ok {
				rows = append(rows, row{lib, r})
			}
		}
		if len(rows) == 0 {
			continue
		}
		sort.Slice(rows, func(i, j int) bool { return rows[i].r.NsPerOp < rows[j].r.NsPerOp })
		bestNs, bestAl := math.Inf(1), math.Inf(1)
		for _, rw := range rows {
			bestNs, bestAl = math.Min(bestNs, rw.r.NsPerOp), math.Min(bestAl, rw.r.AllocsPerOp)
		}

		fmt.Fprintf(&sb, `<text x="10" y="%.1f" font-size="14" font-weight="bold" fill="#222">P%d</text>`+"\n", y+14, p)
		fmt.Fprintf(&sb, `<line x1="10" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#ddd"/>`+"\n", y+22, w-10, y+22)
		ry := y + 30
		for rank, rw := range rows {
			isBase := rank == 0
			color := libColor(rw.lib)
			weight, textFill := "normal", "#333"
			if isBase {
				weight, textFill = "bold", "#111"
				fmt.Fprintf(&sb, `<rect x="4" y="%.1f" width="%.0f" height="%.1f" fill="#fff7d6"/>`+"\n", ry-2, w-8, rowH-2)
			}
			mid := ry + rowH/2 + 3
			fmt.Fprintf(&sb, `<text x="%.0f" y="%.1f" text-anchor="end" font-size="12" font-weight="%s" fill="%s">#%d %s</text>`+"\n", nameX, mid, weight, textFill, rank+1, rw.lib)
			nsL := logBarLen(rw.r.NsPerOp, loNs, hiNs, barW)
			alL := logBarLen(rw.r.AllocsPerOp, loAl, hiAl, barW)
			fmt.Fprintf(&sb, `<rect x="%.0f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`+"\n", col1X, ry+2, nsL, rowH-8, color)
			fmt.Fprintf(&sb, `<rect x="%.0f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`+"\n", col2X, ry+2, alL, rowH-8, color)
			nsTxt := fmtVal(rw.r.NsPerOp) + " ns  (" + rankLabel(rw.r.NsPerOp, bestNs) + ")"
			alTxt := fmtVal(rw.r.AllocsPerOp) + "  (" + rankLabel(rw.r.AllocsPerOp, bestAl) + ")"
			fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" font-size="11" font-weight="%s" fill="%s">%s</text>`+"\n", col1X+nsL+6, mid, weight, textFill, nsTxt)
			fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" font-size="11" font-weight="%s" fill="%s">%s</text>`+"\n", col2X+alL+6, mid, weight, textFill, alTxt)
			ry += rowH
		}
		y += panelH + panelGap
	}
	sb.WriteString("</svg>\n")

	outPath := filepath.Join("../images", rankingChartFile(key))
	if err := os.WriteFile(outPath, []byte(sb.String()), 0o644); err != nil {
		fmt.Printf("Error writing ranking chart %s: %v\n", outPath, err)
		return
	}
	fmt.Printf("Generated: %s\n", outPath)
}
