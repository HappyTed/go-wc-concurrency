package main

import (
	"bytes"
	"testing"
)

var testData = []byte(`in this text
words 9
lines 4
bytes хз`)

func BenchmarkCount_ThreeLoops(b *testing.B) {

	r := bytes.NewReader(testData)

	for i := 0; i < b.N; i++ {
		ThreeLoops(r)
	}
}

func BenchmarkCount_ThreeGorutines(b *testing.B) {

	r := bytes.NewReader(testData)

	for i := 0; i < b.N; i++ {
		ThreeGorutines(r)
	}
}

func BenchmarkCount_LineByLine(b *testing.B) {

	r := bytes.NewReader(testData)

	for i := 0; i < b.N; i++ {
		LineByLine(r)
	}
}
