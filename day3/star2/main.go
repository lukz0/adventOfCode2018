package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	claims := readFrom("input.txt")
	inchMap := mapClaims(claims)
	fmt.Println(inchMap.findIDWithoutOverlap())
}

func readFrom(path string) (claims []claim) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	fr := bufio.NewReader(file)

	var line string

	for line, err = fr.ReadString('\n'); err == nil; line, err = fr.ReadString('\n') {
		line = strings.Replace(line, "#", "", -1)
		line = strings.Replace(line, "@ ", "", -1)
		line = strings.Replace(line, ",", " ", -1)
		line = strings.Replace(line, ":", "", -1)
		line = strings.Replace(line, "x", " ", -1)
		line = strings.Replace(line, "\n", "", -1)

		slicedLine := strings.Fields(line)
		if len(slicedLine) > 5 {
			break
		}

		fromLeft, err := strconv.Atoi(slicedLine[1])
		fromTop, err := strconv.Atoi(slicedLine[2])
		width, err := strconv.Atoi(slicedLine[3])
		height, err := strconv.Atoi(slicedLine[4])
		id, err := strconv.Atoi(slicedLine[0])

		if err != nil {
			panic(err)
		}

		claims = append(claims, claim{
			fromLeft: fromLeft,
			fromTop:  fromTop,
			width:    width,
			height:   height,
			id:       id,
		})
	}

	return
}

type claim struct {
	fromLeft int
	fromTop  int
	width    int
	height   int
	id       int
}

type squareInch struct {
	claimedByID []int
}

type squareInchDoubleMap map[int]map[int]squareInch

func mapClaims(claims []claim) (retMap squareInchDoubleMap) {
	retMap = make(squareInchDoubleMap)

	for i := range claims {
		for x := 0; x < claims[i].width; x++ {
			for y := 0; y < claims[i].height; y++ {
				_, ok := retMap[claims[i].fromLeft+x]
				if !ok {
					retMap[claims[i].fromLeft+x] = make(map[int]squareInch)
				}
				_, ok = retMap[claims[i].fromLeft+x][claims[i].fromTop+y]
				if !ok {
					retMap[claims[i].fromLeft+x][claims[i].fromTop+y] = squareInch{
						claimedByID: make([]int, 0),
					}
				}
				retMap[claims[i].fromLeft+x][claims[i].fromTop+y] = squareInch{
					claimedByID: append(retMap[claims[i].fromLeft+x][claims[i].fromTop+y].claimedByID, claims[i].id),
				}
			}
		}
	}
	return
}

func (sidm *squareInchDoubleMap) findIDWithoutOverlap() int {
	var (
		allIDs         []int = make([]int, 0)
		overlappingIDs []int = make([]int, 0)
	)

	isInSlice := func(slice []int, element int) bool {
		for i := range slice {
			if slice[i] == element {
				return true
			}
		}
		return false
	}

	for x := range *sidm {
		for y := range (*sidm)[x] {
			if len((*sidm)[x][y].claimedByID) > 1 {
				for i := range (*sidm)[x][y].claimedByID {
					if !isInSlice(overlappingIDs, (*sidm)[x][y].claimedByID[i]) {
						overlappingIDs = append(overlappingIDs, (*sidm)[x][y].claimedByID[i])
					}
					if !isInSlice(allIDs, (*sidm)[x][y].claimedByID[i]) {
						allIDs = append(allIDs, (*sidm)[x][y].claimedByID[i])
					}
				}
			} else {
				if !isInSlice(allIDs, (*sidm)[x][y].claimedByID[0]) {
					allIDs = append(allIDs, (*sidm)[x][y].claimedByID[0])
				}
			}
		}
	}

	for _, v := range allIDs {
		if !isInSlice(overlappingIDs, v) {
			return v
		}
	}

	return 0
}
