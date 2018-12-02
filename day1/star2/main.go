package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	fileReader := bufio.NewReader(file)

	var (
		repeatedFrequency    int
		hasRepeatedFrequency bool
		wasPresentFrequency  map[int]bool = make(map[int]bool)
		currentFrequency     int
		fileLine             string
	)

	for fileLine, err = fileReader.ReadString('\n'); !hasRepeatedFrequency; fileLine, err = fileReader.ReadString('\n') {
		fileLine = strings.Replace(fileLine, "\n", "", -1)
		number, atoierr := strconv.Atoi(fileLine)
		if atoierr != nil {
			if fileLine == "" {
				file.Seek(0, 0)
				continue
			} else {
				panic(atoierr)
			}
		}
		currentFrequency += number
		_, present := wasPresentFrequency[currentFrequency]
		if present {
			hasRepeatedFrequency = true
			repeatedFrequency = currentFrequency
		} else {
			wasPresentFrequency[currentFrequency] = true
		}
	}
	fmt.Println(repeatedFrequency)

	file.Close()
}
