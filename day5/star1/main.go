package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
)

func main() {
	start := time.Now()
	p := readPolymerFrom("input.txt")

	p1 := fullyReactPolymer(p)

	fmt.Printf("The polymer's size went from %d to %d\n", len(p), len(p1))
	fmt.Println(time.Since(start))
}

func readPolymerFrom(path string) []rune {
	f, e := os.Open(path)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	fr := bufio.NewReader(f)
	polymer, e := fr.ReadString('\n')
	if e != nil {
		panic(e)
	}
	return []rune(strings.TrimSpace(polymer))
}

func reactPolymer(p []rune) (n []rune) {
	n = make([]rune, 0)
	skipNext := false

	for i, r := range p {
		if skipNext {
			skipNext = false
			continue
		}

		if i != len(p)-1 {
			next := p[i+1]
			if unicode.IsUpper(r) && unicode.IsLower(next) {
				if r == unicode.ToUpper(next) {
					skipNext = true
				} else {
					n = append(n, r)
				}
			} else if unicode.IsLower(r) && unicode.IsUpper(next) {
				if r == unicode.ToLower(next) {
					skipNext = true
				} else {
					n = append(n, r)
				}
			} else {
				n = append(n, r)
			}
		} else {
			n = append(n, r)
		}
	}
	return
}

func fullyReactPolymer(p []rune) (n []rune) {
	n = reactPolymer(p)
	n2 := reactPolymer(n)

	for len(n) != len(n2) {
		if len(n) > len(n2) {
			n = reactPolymer(n2)
		} else {
			n2 = reactPolymer(n)
		}
	}
	return
}
