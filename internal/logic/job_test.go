package logic_test

import (
	"go-wc-concurrency/internal/entity"
	"go-wc-concurrency/internal/logic"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroLine(t *testing.T) {
	input := ""
	expected := entity.OutputData{"", 0, 0, 0}

	counter, err := logic.NewJob("", strings.NewReader(input))
	assert.NotNil(t, err)

	counter.Calculate()
}

func TestOneSpace(t *testing.T) {
	input := " "
	expected := entity.OutputData{"", 0, 0, 5}
}

func TestOneLine(t *testing.T) {
	input := "hello world!"
	expected := entity.OutputData{"", 1, 2, 11}
}

func TestMoreLines(t *testing.T) {
	input := `hello
world`
	expected := entity.OutputData{"", 2, 2, 11}

}

// пограничные ситуации

func TestOneChar(t *testing.T) {
	input := "a"
	expected := entity.OutputData{"", 1, 1, 1}
}

func TestOneLineWithALineBreak(t *testing.T) {
	input := "hello\n"
	expected := entity.OutputData{"", 1, 1, 6}
}

func TestOneLineWithMoreSpaces(t *testing.T) {
	input := "hello world\tfoo bar"
	expected := entity.OutputData{"", 1, 4, 19}
}

func TestMoreLinesWithBlankLines(t *testing.T) {
	input := `a

b
`
	expected := entity.OutputData{"", 3, 2, 4}

}

func TestUnicodeLine(t *testing.T) {
	input := "привет мир"
	expected := entity.OutputData{"", 1, 2, 19}
}

func TestMegaString(t *testing.T) {
	input := strings.Repeat("a", 10000000)
	expected := entity.OutputData{"", 1, 1, 10000000}
}

func TestBinaryFile(t *testing.T) {
	input := []byte("\x00\x01hello\nworld")
	expected := entity.OutputData{"", 1, 0, 13}
}
