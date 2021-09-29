package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	errChan := Run(ctx,
		func(c chan<- error) {
			cancel()
			c <- fmt.Errorf("error 1 %s", time.Now().Format("2006-01-02 15:04:05.000000"))

			fmt.Println(stringifyUserIDs("1", "2"))
		},
		func(c chan<- error) {
			c <- fmt.Errorf("error 2 %s", time.Now().Format("2006-01-02 15:04:05.000000"))
			fmt.Println(stringifyUserIDs("3", "4"))
			c <- fmt.Errorf("error 3 %s", time.Now().Format("2006-01-02 15:04:05.000000"))
		})
	for err := range errChan {
		fmt.Println(err.Error())
	}
}

func stringifyUserIDs(IDs ...string) string {
	return "[\"" + strings.Join(IDs, "\",\"") + "\"]"
}

//Run parallel functions
func Run(ctx context.Context, tasks ...func(chan<- error)) <-chan error {
	n := len(tasks)
	errCh := make(chan error, n)
	w := &sync.WaitGroup{}
	w.Add(n)

	go func() {
		defer close(errCh)
		for _, task := range tasks {
			go func(wg *sync.WaitGroup, t func(chan<- error)) {
				if ctx.Err() != nil {
					errCh <- fmt.Errorf("Cancelled at %s", time.Now().Format("2006-01-02 15:04:05.000000"))
				} else {
					t(errCh)
				}
				wg.Done()
			}(w, task)
		}

		w.Wait()
	}()

	return errCh
}
