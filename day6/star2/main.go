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
	dMap := createEmptyDistanceMap(minX, minY, maxX, maxY)
	findDistances(dMap, points)
	fmt.Println(countSmallDistances(dMap, 10000))
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

func createEmptyDistanceMap(minX, minY, maxX, maxY int) map[int]map[int]int {
	coordMap := make(map[int]map[int]int)
	for x := minX; x <= maxX; x++ {
		coordMap[x] = make(map[int]int)
		for y := minY; y <= maxY; y++ {
			coordMap[x][y] = 0
		}
	}
	return coordMap
}

func findDistances(distanceMap map[int]map[int]int, points []point) {
	for x := range distanceMap {
		for y := range distanceMap[x] {
			for _, p := range points {
				distanceX := p.x - x
				distanceY := p.y - y
				if distanceX < 0 {
					distanceX = -distanceX
				}
				if distanceY < 0 {
					distanceY = -distanceY
				}
				distanceMap[x][y] += distanceX + distanceY
			}
		}
	}
}

func countSmallDistances(distanceMap map[int]map[int]int, lessThan int) int {
	var count int

	for x := range distanceMap {
		for y := range distanceMap[x] {
			if distanceMap[x][y] < lessThan {
				count++
			}
		}
	}

	return count
}
