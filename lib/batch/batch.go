package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	wg := new(sync.WaitGroup)
	jobs := make(chan int64)
	result := make(chan user)

	wg.Add(int(pool))
	for i := int64(0); i < pool; i++ {
		go getWorker(jobs, result, wg)
	}
	go func() {
		wg.Wait()
		close(result)
	}()
	go func() {
		for i := int64(0); i < n; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	for user := range result {
		res = append(res, user)
	}

	return
}

func getWorker(jobs <-chan int64, out chan<- user, wg *sync.WaitGroup) {
	for id := range jobs {
		user := getOne(id)
		out <- user
	}
	wg.Done()
}
