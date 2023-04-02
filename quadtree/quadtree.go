package main

import (
  "fmt"
  "math/rand"
  "github.com/fogleman/gg"
)

const ROOT_BOUNDARY_WIDTH = 1000
const ROOT_BOUNDARY_HEIGHT = 1000
const NUM_POINTS = 100

type Point struct {
  X int
  Y int
}

type Rectangle struct {
  X int
  Y int
  Width int
  Height int
}

type Node struct {
  Pos Point
  Data int
}

type QuadTree struct {
  Boundary Rectangle
  Capacity int
  Depth int
  MaxDepth int
  Points []Point

  Children []*QuadTree
}

func (re Rectangle) contains(point Point) bool {
  return (re.X <= point.X && re.Y <= point.Y && point.X <= (re.X + re.Width) && point.Y <= (re.Y + re.Height))
}

func (qt *QuadTree) insert(point Point) bool {
  if !qt.Boundary.contains(point) {
    return false
  }

  if(len(qt.Children) == 0) {
    if (len(qt.Points) < qt.Capacity && qt.Depth < qt.MaxDepth) {
      fmt.Printf("point (%dx %dy) assigned to quadtree w/ boundary %v\n", point.X, point.Y, qt.Boundary)
      qt.Points = append(qt.Points, point)
      return true
    } else {
      if len(qt.Children) == 0 {
        qt.subdivide()
      }
    }
  }
  inserted := false
  for _, child := range qt.Children {
    if child.insert(point) {
      inserted = true
      break
    }
  }

  return inserted
}

func (qt *QuadTree) subdivide() {
  for i := 0; i < 4; i++ {
    boundary := Rectangle{
      X: qt.Boundary.X + (i % 2) * (qt.Boundary.Width / 2),
      Y: qt.Boundary.Y + (i / 2) * (qt.Boundary.Height / 2),
      Width: qt.Boundary.Width / 2,
      Height:  qt.Boundary.Height / 2,
    }

    child := QuadTree{
      Boundary: boundary,
      Capacity: qt.Capacity,
      MaxDepth: qt.MaxDepth,
      Depth: qt.Depth + 1,
    }

    qt.Children = append(qt.Children, &child)
  }

  // Move existing points to child nodes
  for _, point := range qt.Points {
    for _, child := range qt.Children {
      if child.Boundary.contains(point) && len(child.Points) < child.Capacity {
        child.insert(point)
        break
      }
    }
  }
  qt.Points = nil
}

func drawQuadTree(dc *gg.Context, qt *QuadTree, depth int) {
    fmt.Printf("Depth %d, Points: %d\n", depth, len(qt.Points))
    for i, point := range qt.Points {
      fmt.Printf("  point %d: %dx %dy\n", i, point.X, point.Y) 
    }
    drawBoundary(dc, qt.Boundary, qt.Depth)
    drawPoints(dc, qt.Points)
    for _, child := range qt.Children {
        drawQuadTree(dc, child, depth+1)
    }
}

func drawBoundary(dc *gg.Context, boundary Rectangle, depth int) {
    fmt.Printf("Depth %d, Boundary: %v\n", depth, boundary)
    lineWidth := 1.0 + (float64(depth) * 0.5)
    dc.SetLineWidth(lineWidth)

    dc.DrawRectangle(float64(boundary.X), float64(boundary.Y), float64(boundary.Width), float64(boundary.Height))
    dc.Stroke()
    dc.SetRGB(0, 0, 0)
    dc.DrawString(fmt.Sprintf("(%d, %d)", boundary.X, boundary.Y), float64(boundary.X), float64(boundary.Y))
}

func drawPoints(dc *gg.Context, points []Point) {
    for _, point := range points {
        dc.SetRGB(1, 0, 0) // red
        radius := 2.0
        x := float64(point.X)
        y := float64(point.Y)
        dc.DrawCircle(x, y, radius)
        dc.Fill()
        dc.SetRGB(0, 0, 0)
        dc.DrawString(fmt.Sprintf("(%d, %d)", point.X, point.Y), x, y)
    }
}

func main() {
  tree := QuadTree{
    Boundary: Rectangle{X: 0, Y: 0, Width: ROOT_BOUNDARY_WIDTH, Height: ROOT_BOUNDARY_HEIGHT},
    Capacity: 4,
    Depth: 0,
    MaxDepth: 10,
  }
  points := []Point{}
  for i := 0; i < NUM_POINTS; i++ {
    points = append(points, Point{X:rand.Intn(ROOT_BOUNDARY_WIDTH), Y: rand.Intn(ROOT_BOUNDARY_HEIGHT)})
  }
  for i, point := range(points) {
    fmt.Printf("inserting point %d: %dx %dy\n", i, point.X, point.Y) 
    tree.insert(point)
  }

  // Draw the QuadTree
  dc := gg.NewContext(ROOT_BOUNDARY_WIDTH, ROOT_BOUNDARY_HEIGHT)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	drawQuadTree(dc, &tree, 0)
	err := dc.SavePNG("quadtree.png")
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}
	fmt.Println("Saved quadtree visualization to quadtree.png")
}
