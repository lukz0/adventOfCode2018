package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fr := newFileReader("input.txt")
	defer fr.closeFile()
	letterList := fr.createLetterList()
	sl := letterList.createSimilarityList()

	var (
		id1num        int
		id2num        int
		id1           []rune
		id2           []rune
		commonLetters []rune = make([]rune, 0)
	)

	for i := range sl {
		for j := range sl[i] {
			if sl[i][j] == sl[i][i]-1 {
				id1num = i
				id2num = j
			}
		}
	}
	id1 = letterList[id1num]
	id2 = letterList[id2num]
	for i := range id1 {
		if id1[i] == id2[i] {
			commonLetters = append(commonLetters, id1[i])
		}
	}
	fmt.Println(string(commonLetters))
}

type idList [][]rune

func (idl *idList) createSimilarityList() [][]int {
	var (
		returnedList [][]int = make([][]int, len(*idl))
	)

	for i := range returnedList {
		returnedList[i] = make([]int, len(*idl))

		for j := range returnedList[i] {

			for k := range (*idl)[i] {
				if (*idl)[i][k] == (*idl)[j][k] {
					returnedList[i][j]++
				}
			}
		}
	}

	return returnedList
}

type fileReader struct {
	file   *os.File
	reader *bufio.Reader
}

func (fr *fileReader) createLetterList() idList {
	fr.file.Seek(0, 0)

	var (
		returnedRunes      idList = make(idList, 0)
		lineCurrentlyEmpty bool   = true
	)

	for r, _, err := fr.reader.ReadRune(); err == nil; r, _, err = fr.reader.ReadRune() {
		if r == '\n' {
			lineCurrentlyEmpty = true
		} else {
			if lineCurrentlyEmpty {
				returnedRunes = append(returnedRunes, make([]rune, 0))
				lineCurrentlyEmpty = false
			}
			returnedRunes[len(returnedRunes)-1] = append(returnedRunes[len(returnedRunes)-1], r)
		}
	}

	return returnedRunes
}

func (fr *fileReader) closeFile() {
	fr.file.Close()
}

func newFileReader(path string) fileReader {
	file, err := os.Open(path)
	if err != nil {
		panic("Error opening file: " + err.Error())
	}
	return fileReader{file: file, reader: bufio.NewReader(file)}
}
