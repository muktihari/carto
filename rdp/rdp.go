// Package rdp host the implementation of [Ramer–Douglas–Peucker algorithm](https://w.wiki/B6U3) algorithm.
package rdp

import "math"

// Point represents coordinates in 2D plane.
type Point struct {
	X, Y  float64 // X and Y represent coordinates (e.g Latitude and Longitude)
	Index int     // [Optional] Index reference to the origin data.
}

// Simplify simplifies a curve by reducing the number of points in a curve composed of line segments,
// resulting in a similar curve with fewer points. Ref: [Ramer–Douglas–Peucker algorithm](https://w.wiki/B6U3).
//
// Note: The resulting slice is a reslice of given points (it shares the same underlying array) for efficiency.
// It works similar to append, so the input points should not be used after this call, use only the returned value.
func Simplify(points []Point, epsilon float64) []Point {
	if len(points) <= 2 {
		return points
	}

	var (
		index   int
		maxDist float64
		first   = points[0]
		last    = points[len(points)-1]
	)

	for i := range points {
		d := perpendicularDistance(points[i], first, last)
		if d > maxDist {
			maxDist = d
			index = i
		}
	}

	if maxDist <= epsilon {
		return append(points[:0], first, last)
	}

	// Move index to avoids infinite recursive as slice input
	// for next operation is never changed if we keep it as is.
	if index == 0 || index == len(points) {
		index++
	}

	left, right := points[:index], points[index:]

	return append(Simplify(left, epsilon), Simplify(right, epsilon)...)
}

// perpendicularDistance calculates the perpendicular distance from a point to a line segment
func perpendicularDistance(p, start, end Point) float64 {
	if start.X == end.X && start.Y == end.Y {
		// Find distance between p and (start or end)
		return euclidean(p, start)
	}

	// Standard Form: Ax + Bx + C = 0
	A := end.Y - start.Y
	B := start.X - end.X
	C := (end.X * start.Y) - (start.X * end.Y)

	// d = | Ax + By + C = 0 | / ✓(A²+B²)
	return math.Abs(A*p.X+B*p.Y+C) / math.Sqrt(A*A+B*B)
}

// euclidean calculates the distance between two points.
func euclidean(p1, p2 Point) float64 {
	x := p2.X - p1.X
	y := p2.Y - p1.Y
	return math.Sqrt(x*x + y*y)
}
