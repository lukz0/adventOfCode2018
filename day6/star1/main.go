package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	points := coordinatesFromFile("input.txt")
	minX, minY, maxX, maxY := areaSize(points)
	minX--
	minY--
	maxX++
	maxY++
	cMap := createEmptyPointMap(minX, minY, maxX, maxY)
	addPointsToCoordMap(cMap, points)
	findClosestPoints(cMap, points)

	isInfMap := findInfiniteAreas(cMap)
	sizes := findAreaSizes(cMap, points)

	index, size := findBiggestNonInfiniteArea(sizes, isInfMap)
	fmt.Printf("The answer is \"%d\" which is the size of area #%d\n", size, index)
}

type point struct {
	x int
	y int
}

func coordinatesFromFile(path string) []point {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	fileReader := bufio.NewReader(file)
	retCoord := make([]point, 0)

	for line, err := fileReader.ReadString('\n'); err == nil; line, err = fileReader.ReadString('\n') {
		line = strings.TrimSpace(line)
		line = strings.Replace(line, ",", "", -1)
		stringCoordinate := strings.Fields(line)
		if len(stringCoordinate) != 2 {
			break
		}
		x, err := strconv.Atoi(stringCoordinate[0])
		y, err := strconv.Atoi(stringCoordinate[1])
		if err != nil {
			panic(err)
		}
		retCoord = append(retCoord, point{x: x, y: y})
	}

	return retCoord
}

func areaSize(coords []point) (minX, minY, maxX, maxY int) {
	if len(coords) == 0 {
		panic("the point slice is empty")
	}

	minX, minY, maxX, maxY = coords[0].x, coords[0].y, coords[0].x, coords[0].y
	for i := 1; i < len(coords); i++ {
		if coords[i].x < minX {
			minX = coords[i].x
		} else if coords[i].x > maxX {
			maxX = coords[i].x
		}

		if coords[i].y < minY {
			minY = coords[i].y
		} else if coords[i].y > maxY {
			maxY = coords[i].y
		}
	}

	return
}

type coordinate struct {
	coordType           int
	pointIndex          int
	closestToPointIndex int
	hasClosestPoint     bool
}

const (
	COORD_POINT = iota
	COORD_EMPTY
)

func createEmptyPointMap(minX, minY, maxX, maxY int) map[int]map[int]coordinate {
	coordMap := make(map[int]map[int]coordinate)
	for x := minX; x <= maxX; x++ {
		coordMap[x] = make(map[int]coordinate)
		for y := minY; y <= maxY; y++ {
			coordMap[x][y] = coordinate{coordType: COORD_EMPTY}
		}
	}
	return coordMap
}

func addPointsToCoordMap(coordMap map[int]map[int]coordinate, points []point) {
	for i, p := range points {
		coordMap[p.x][p.y] = coordinate{
			coordType:  COORD_POINT,
			pointIndex: i,
		}
	}
}

func findClosestPoints(coordMap map[int]map[int]coordinate, points []point) {
	if len(points) == 0 {
		panic("Point slice empty")
	}
	for x := range coordMap {
		for y := range coordMap[x] {
			if coordMap[x][y].coordType == COORD_POINT {
				continue
			}

			distanceX := points[0].x - x
			distanceY := points[0].y - y
			if distanceX < 0 {
				distanceX = -distanceX
			}
			if distanceY < 0 {
				distanceY = -distanceY
			}
			closestPointDistance := distanceX + distanceY
			closestPoint := 0
			hasMultipleClosestDistances := false

			for i := 1; i < len(points); i++ {
				distanceX = points[i].x - x
				distanceY = points[i].y - y
				if distanceX < 0 {
					distanceX = -distanceX
				}
				if distanceY < 0 {
					distanceY = -distanceY
				}
				if closestPointDistance > distanceX+distanceY {
					hasMultipleClosestDistances = false
					closestPointDistance = distanceX + distanceY
					closestPoint = i
				} else if closestPointDistance == distanceX+distanceY {
					hasMultipleClosestDistances = true
				}
			}

			if hasMultipleClosestDistances {
				coordMap[x][y] = coordinate{
					coordType:       COORD_EMPTY,
					hasClosestPoint: false,
				}
			} else {
				coordMap[x][y] = coordinate{
					coordType:           COORD_EMPTY,
					hasClosestPoint:     true,
					closestToPointIndex: closestPoint,
				}
			}
		}
	}
}

func findInfiniteAreas(coordMap map[int]map[int]coordinate) map[int]bool {
	infiniteAreas := make(map[int]bool)
	for x := range coordMap {
		if coordMap[x][0].coordType == COORD_POINT {
			infiniteAreas[coordMap[x][0].pointIndex] = true
		} else if coordMap[x][0].hasClosestPoint {
			infiniteAreas[coordMap[x][0].closestToPointIndex] = true
		}

		if coordMap[x][len(coordMap)-1].coordType == COORD_POINT {
			infiniteAreas[coordMap[x][len(coordMap)-1].pointIndex] = true
		} else if coordMap[x][len(coordMap)-1].hasClosestPoint {
			infiniteAreas[coordMap[x][len(coordMap)-1].closestToPointIndex] = true
		}
	}

	for y := range coordMap[0] {
		if coordMap[0][y].coordType == COORD_POINT {
			infiniteAreas[coordMap[0][y].pointIndex] = true
		} else if coordMap[0][y].hasClosestPoint {
			infiniteAreas[coordMap[0][y].closestToPointIndex] = true
		}
	}
	for y := range coordMap[len(coordMap)-1] {
		if coordMap[len(coordMap)-1][y].coordType == COORD_POINT {
			infiniteAreas[coordMap[len(coordMap)-1][y].pointIndex] = true
		} else if coordMap[len(coordMap)-1][y].hasClosestPoint {
			infiniteAreas[coordMap[len(coordMap)-1][y].closestToPointIndex] = true
		}
	}

	return infiniteAreas
}

func findAreaSizes(coordMap map[int]map[int]coordinate, points []point) []int {
	areaSizes := make([]int, len(points))
	for x := range coordMap {
		for y := range coordMap[x] {
			if coord := coordMap[x][y]; coord.coordType == COORD_POINT {
				areaSizes[coord.pointIndex]++
			} else if coord.hasClosestPoint {
				areaSizes[coord.closestToPointIndex]++
			}
		}
	}

	return areaSizes
}

func findBiggestNonInfiniteArea(sizes []int, isInfMap map[int]bool) (index, size int) {
	for i := range sizes {
		isInf, _ := isInfMap[i]
		if sizes[i] > size && !isInf {
			size = sizes[i]
			index = i
		}
	}

	return
}
