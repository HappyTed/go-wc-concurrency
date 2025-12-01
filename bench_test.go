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

	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(testData)
		ThreeLoops(r)
	}
}

func BenchmarkCount_ThreeGorutines(b *testing.B) {

	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(testData)
		ThreeGorutines(r)
	}
}

func BenchmarkCount_LineByLine(b *testing.B) {

	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(testData)
		LineByLine(r)
	}
}
