package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	p := readPolymerFrom("input.txt")

	m, lr := mapReactedLenWithoutUnit(p)

	fmt.Printf("The reacted length is shortest without %s, without it the reacted length is %d\n", strconv.QuoteRune(lr), m[lr])
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

func removeUnit(from []rune, unit rune) (result []rune) {
	result = make([]rune, 0)

	for i := range from {
		if r := from[i]; unicode.ToUpper(r) != unit {
			result = append(result, r)
		}
	}

	return
}

func mapReactedLenWithoutUnit(p []rune) (m map[rune]int, lowestUnit rune) {
	m = make(map[rune]int)
	l := make(map[rune]bool)

	lowDefined := false

	for _, r := range p {
		_, ok := l[unicode.ToUpper(r)]
		if !ok {
			l[unicode.ToUpper(r)] = true
		}
	}

	for r := range l {
		pLen := len(fullyReactPolymer(removeUnit(p, r)))

		if !lowDefined || pLen < m[lowestUnit] {
			lowestUnit = r
			lowDefined = true
		}
		m[r] = pLen
	}

	return
}
