package logic

import (
	"context"
	"strings"
	"sync"
)

type WordCount struct {
	ctx context.Context
	wg  sync.WaitGroup
	sb  strings.Builder
}
