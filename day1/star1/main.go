package main

import (
	"bufio"
	"fmt"
	"io"
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
		fileLine string
		sum      int
	)

	for fileLine, err = fileReader.ReadString('\n'); err == nil; fileLine, err = fileReader.ReadString('\n') {
		number, atoierr := strconv.Atoi(strings.Replace(fileLine, "\n", "", -1))
		if atoierr != nil {
			panic(atoierr)
		}
		sum += number
	}

	fmt.Printf("Result: %d\n", sum)
	if err != nil && err != io.EOF {
		fmt.Println("Error: " + err.Error())
	}

	file.Close()
}
