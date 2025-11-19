package main

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

var (
	wordsCounter uint
	linesCounter uint
	bytesCounter uint
)

func ThreeLoops(data io.Reader) {

	scanner := bufio.NewScanner(data)

	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordsCounter++
	}
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		bytesCounter++
	}
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		bytesCounter++
	}
}

func ThreeGorutines(data io.Reader) {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(data)
		scanner.Split(bufio.ScanBytes)
		for scanner.Scan() {
			bytesCounter++
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(data)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			linesCounter++
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(data)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			wordsCounter++
		}
		wg.Done()
	}()
}

func LineByLine(data io.Reader) {

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		linesCounter++
		bytesCounter += uint(len(scanner.Bytes()))

		line := scanner.Text()
		wordsCounter += uint(len(strings.Split(line, " ")))
	}
}
