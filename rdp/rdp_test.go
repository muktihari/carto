package rdp

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSimplify(t *testing.T) {
	tt := []struct {
		name    string
		epsilon float64
		points  []Point
		expect  []Point
	}{
		{
			name:    "one points",
			epsilon: 0.1,
			points: []Point{
				{X: 0.1, Y: 0.2},
			},
			expect: []Point{
				{X: 0.1, Y: 0.2},
			},
		},
		{
			name:    "two points",
			epsilon: 0.1,
			points: []Point{
				{X: 0.0, Y: 0.0},
				{X: 0.1, Y: 0.2},
			},
			expect: []Point{
				{X: 0.0, Y: 0.0},
				{X: 0.1, Y: 0.2},
			},
		},
		{
			name:    "two same points",
			epsilon: 0.1,
			points: []Point{
				{X: 0.0, Y: 0.0},
				{X: 0.0, Y: 0.0},
			},
			expect: []Point{
				{X: 0.0, Y: 0.0},
				{X: 0.0, Y: 0.0},
			},
		},
		{
			name:    "small points",
			epsilon: 0.1,
			points: []Point{
				{X: 0.0, Y: 0.0},
				{X: 0.1, Y: 0.2},
				{X: 0.2, Y: 0.3},
				{X: 0.3, Y: 0.4},
				{X: 0.4, Y: 0.8},
				{X: 0.5, Y: 0.1},
			},
			expect: []Point{
				{X: 0.0, Y: 0.0},
				{X: 0.3, Y: 0.4},
				{X: 0.4, Y: 0.8},
				{X: 0.5, Y: 0.1},
			},
		},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("[%d] %s", i, tc.name), func(t *testing.T) {
			result := Simplify(tc.points, tc.epsilon)
			if diff := cmp.Diff(result, tc.expect); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestPerpendicularDistance(t *testing.T) {
	tt := []struct {
		name  string
		p     Point
		start Point
		end   Point
		d     float64
		prec  float64
	}{
		{
			name:  "valid result",
			p:     Point{X: 5, Y: 6},
			start: Point{X: 0, Y: -1.3333333333333333},
			end:   Point{X: 2, Y: 0},
			d:     3.328,
			prec:  1000,
		},
		{
			name:  "zero result",
			p:     Point{X: 5, Y: 6},
			start: Point{X: 2, Y: 0},
			end:   Point{X: 2, Y: 0},
			d:     6.708, // euclidean((5,6), (2,0))
			prec:  1000,
		},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("[%d] %s", i, tc.name), func(t *testing.T) {
			d := perpendicularDistance(tc.p, tc.start, tc.end)
			d = math.Round(d*tc.prec) / tc.prec
			tc.d = math.Round(tc.d*tc.prec) / tc.prec
			if d != tc.d {
				t.Fatalf("expected: %g, got: %g", tc.d, d)
			}
		})
	}
}

var update = flag.Bool("update", false, "update the test file")

func BenchmarkSimplify(b *testing.B) {
	flag.Parse()

	const filename = "rdp_bench.txt"
	const n = 1000
	if *update {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			b.Fatalf("update: could not open file: %v", err)
		}
		w := bufio.NewWriter(f)
		points := makeRandomPoints(n)
		for _, p := range points {
			fmt.Fprintf(w, "%g,%g\n", p.X, p.Y)
		}
		if err := w.Flush(); err != nil {
			f.Close()
			b.Fatalf("update: could not flush: %v", err)
		}
		f.Close()
	}

	bb, err := os.ReadFile(filename)
	if err != nil {
		b.Fatalf("could not open file: %v", err)
	}

	buf := bytes.NewBuffer(bb)

	points := make([]Point, 0, n)
	for {
		line, err := buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			b.Fatalf("could not read bytes: %v", err)
		}
		parts := bytes.Split(bytes.Trim(line, "\n"), []byte{','})
		x, err := strconv.ParseFloat(string(parts[0]), 64)
		if err != nil {
			b.Fatalf("could not parse x: %v", err)
		}
		y, err := strconv.ParseFloat(string(parts[1]), 64)
		if err != nil {
			b.Fatalf("could not parse y: %v", err)
		}
		points = append(points, Point{X: x, Y: y})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Simplify(points, 0.5)
	}
}

func makeRandomPoints(n int) []Point {
	p := make([]Point, n)
	for i := range p {
		if i == 0 {
			p[i].X = rand.Float64()
			p[i].Y = rand.Float64()
		} else {
			p[i].X = p[i-1].X + rand.Float64()
		}
		p[i].Y = p[i].X + 10*rand.Float64()
	}
	return p
}
