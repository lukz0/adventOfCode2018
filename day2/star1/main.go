package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fr := newFileReader("input.txt")
	fmt.Println(fr.countAmountLetters(2) * fr.countAmountLetters(3))
	fr.CloseFile()
}

type fileReader struct {
	file   *os.File
	reader *bufio.Reader
}

func (fr *fileReader) countAmountLetters(amount int) (count int) {
	fr.file.Seek(0, 0)
	letterMap := make(map[rune]int)

	for r, _, err := fr.reader.ReadRune(); err == nil; r, _, err = fr.reader.ReadRune() {
		if r == '\n' {
			hasDoubleLetter := false

			for _, v := range letterMap {
				if v == amount {
					hasDoubleLetter = true
				}
			}

			if hasDoubleLetter {
				count++
			}
			letterMap = make(map[rune]int)
		} else {
			v, _ := letterMap[r]
			letterMap[r] = v + 1
		}
	}

	return
}

func (fr *fileReader) CloseFile() {
	fr.file.Close()
}

func newFileReader(path string) fileReader {
	file, err := os.Open(path)
	if err != nil {
		panic("Error opening file: " + err.Error())
	}
	return fileReader{file: file, reader: bufio.NewReader(file)}
}
