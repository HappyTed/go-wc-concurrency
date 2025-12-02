package logic_test

import (
	"go-wc-concurrency/internal/entity"
	"go-wc-concurrency/internal/logic"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func real(input string, t *testing.T) *entity.OutputData {
	counter, err := logic.NewJob(strings.NewReader(input))
	require.NoError(t, err, "failed to create Job")
	data, err := counter.Calculate()
	require.NoError(t, err, "failed to calculate")
	return data
}

func TestZeroLine(t *testing.T) {
	input := ""
	expected := &entity.OutputData{"", 0, 0, 0}

	real := real(input, t)

	assert.Equal(t, expected, real, "Actual counter values ​​do not match expected values")
}

func TestOneSpace(t *testing.T) {
	input := " "
	expected := &entity.OutputData{"", 0, 0, 1}

	real := real(input, t)

	assert.Equal(t, expected, real, "Actual counter values ​​do not match expected values")
}

func TestOneLine(t *testing.T) {
	input := "hello world!"
	expected := &entity.OutputData{"", 0, 2, 12}

	real := real(input, t)

	assert.Equal(t, expected, real, "Actual counter values ​​do not match expected values")
}

func TestMoreLines(t *testing.T) {
	input := `hello
world`
	expected := &entity.OutputData{"", 1, 2, 11}

	real := real(input, t)

	assert.Equal(t, expected, real, "Actual counter values ​​do not match expected values")

}

// пограничные ситуации

func TestOneChar(t *testing.T) {
	input := "a"
	expected := &entity.OutputData{"", 0, 1, 1}

	real := real(input, t)

	assert.Equal(t, expected, real, "Actual counter values ​​do not match expected values")
}

func TestOneLineWithALineBreak(t *testing.T) {
	input := "hello\n"
	expected := &entity.OutputData{"", 1, 1, 6}

	real := real(input, t)

	assert.Equal(t, expected, real, "Actual counter values ​​do not match expected values")
}

func TestOneLineWithMoreSpaces(t *testing.T) {
	input := "hello world\tfoo bar"
	expected := &entity.OutputData{"", 0, 4, 19}

	real := real(input, t)

	assert.Equal(t, expected, real, "Actual counter values ​​do not match expected values")
}

func TestMoreLinesWithBlankLines(t *testing.T) {
	input := `a

b
`
	expected := &entity.OutputData{"", 3, 2, 5}

	real := real(input, t)

	assert.Equal(t, expected, real, "Actual counter values ​​do not match expected values")
}

func TestUnicodeLine(t *testing.T) {
	input := "привет мир"
	expected := &entity.OutputData{"", 0, 2, 19}

	real := real(input, t)

	assert.Equal(t, expected, real, "Actual counter values ​​do not match expected values")
}

func TestMegaString(t *testing.T) {
	input := strings.Repeat("a", 10000000)
	expected := &entity.OutputData{"", 0, 1, 10000000}

	real := real(input, t)

	assert.Equal(t, expected, real, "Actual counter values ​​do not match expected values")
}
