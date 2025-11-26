package config

import (
	"flag"
	"runtime"
)

type OptionFlag int

const (
	LINES = 1 << iota
	WORDS
	BYTES
)

type Config struct {
	Options    OptionFlag
	Files      []string
	NumWorkers uint8
}

func ReadConfig() *Config {
	cfg := &Config{NumWorkers: uint8(runtime.NumCPU())}

	var (
		IsLines bool
		IsWords bool
		IsBytes bool
	)
	flag.BoolVar(&IsLines, "l", false, "count lines")
	flag.BoolVar(&IsWords, "w", false, "count words")
	flag.BoolVar(&IsBytes, "b", false, "count bytes")

	flag.Parse()

	var options OptionFlag = 0

	if !IsWords && !IsLines && !IsBytes {
		options = BYTES | LINES | WORDS
	} else {
		if IsBytes {
			options = options | BYTES
		}
		if IsLines {
			options = options | LINES
		}
		if IsWords {
			options = options | WORDS
		}
	}

	cfg.Options = options
	cfg.Files = flag.Args()

	return cfg
}
