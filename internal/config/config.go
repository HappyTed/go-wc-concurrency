package config

import (
	"flag"
	"runtime"
)

type Flags int

const (
	LINES = 1 << iota
	WORDS
	BYTES
)

type Config struct {
	Flag       Flags
	FilesPath  []string
	NumWorkers uint8
}

func ReadConfig() *Config {
	cfg := &Config{NumWorkers: uint8(runtime.NumCPU())}
	IsLines := *flag.Bool("l", false, "count lines")
	IsWords := *flag.Bool("w", false, "count words")
	IsBytes := *flag.Bool("b", false, "count bytes")

	flag.Parse()

	var flags Flags

	if IsBytes {
		flags = flags | BYTES
	}
	if IsLines {
		flags = flags | LINES
	}
	if IsWords {
		flags = flags | WORDS
	}

	cfg.Flag = flags
	cfg.FilesPath = flag.Args()

	return cfg
}
