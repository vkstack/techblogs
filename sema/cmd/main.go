package main

import (
	"fmt"
	"sync"

	"github.com/vkstack/techblogs/sema"
	"github.com/vkstack/techblogs/sema/logger"
)

func main() {
	l := logger.NewLogger()
	s := sema.NewIO(5, l)
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			s.LimitedProcess(fmt.Sprintf("file-%d", x))
		}(i)
	}
	wg.Wait()
}
