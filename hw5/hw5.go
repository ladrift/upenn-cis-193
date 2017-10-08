// Homework 5: Goroutines
package main

import "fmt"
import "sync"
import "time"

func main() {
	ints := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ints <- i
		}
		close(ints)
	}()
	isEven := func(i int) bool {
		return i%2 == 0
	}
	for i := range Filter(ints, isEven) {
		fmt.Println(i)
	}

	var t1 taskFunc = func(i int) (int, error) {
		time.Sleep(5 * time.Second)
		return i + 1, nil
	}
	var t2 taskFunc = func(i int) (int, error) {
		time.Sleep(2 * time.Second)
		return i + 2, nil
	}
	var t3 taskFunc = func(i int) (int, error) {
		time.Sleep(3 * time.Second)
		return i + 3, nil
	}
	result, _ := Fastest(1, t1, t2, t3)
	fmt.Printf("result = %v\n", result)
}

type taskFunc func(int) (int, error)

func (t taskFunc) Execute(i int) (int, error) {
	result, err := t(i)
	return result, err
}

// Filter copies values from the input channel into an output channel that match the filter function p
// The function p determines whether an int from the input channel c is sent on the output channel
func Filter(c <-chan int, p func(int) bool) <-chan int {
	output := make(chan int)

	go func() {
		for i := range c {
			if p(i) {
				output <- i
			}
		}
		close(output)
	}()

	return output
}

// Result is a type representing a single result with its index from a slice
type Result struct {
	index  int
	result string
}

// ConcurrentRetry runs all the tasks concurrently and sends the output in a Result channel
//
// concurrent is the limit on the number of tasks running in parallel. Your
// solution must not run more than `concurrent` number of tasks in parallel.
//
// retry is the number of times that the task should be attempted. If a task
// returns an error, the function should be retried immediately up to `retry`
// times. Only send the results of a task into the output channel if it does not error.
//
// Multiple instances of ConcurrentRetry should be able to run simultaneously
// without interfering with one another, so global variables should not be used.
// The function must return the channel without waiting for the tasks to
// execute, and all results should be sent on the output channel. Once all tasks
// have been completed, close the channel.
func ConcurrentRetry(tasks []func() (string, error), concurrent int, retry int) <-chan Result {
	type Task struct {
		index int
		task  func() (string, error)
	}

	taskCh := make(chan Task)
	go func() {
		for i, t := range tasks {
			taskCh <- Task{i, t}
		}
		close(taskCh)
	}()

	results := make(chan Result)
	var wg sync.WaitGroup
	wg.Add(concurrent)
	go func() {
		wg.Wait()
		close(results)
	}()

	worker := func() {
		for t := range taskCh {
			for i := 0; i < retry; i++ {
				result, err := t.task()
				if err == nil {
					results <- Result{t.index, result}
				}
			}
		}
		wg.Done()
	}

	go func() {
		for i := 0; i < concurrent; i++ {
			go worker()
		}
	}()

	return results
}

// Task is an interface for types that process integers
type Task interface {
	Execute(int) (int, error)
}

// Fastest returns the result of the fastest running task
// Fastest accepts any number of Task structs. If no tasks are submitted to
// Fastest(), it should return an error.
// You should return the result of a Task even if it errors.
// Do not leave any pending goroutines. Make sure all goroutines are cleaned up
// properly and any synchronizing mechanisms closed.
func Fastest(input int, tasks ...Task) (int, error) {
	if len(tasks) == 0 {
		return 0, fmt.Errorf("no task is given")
	}

	results := make(chan int, len(tasks))
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	go func() {
		wg.Wait()
		close(results)
	}()

	for _, t := range tasks {
		go func(t Task) {
			result, err := t.Execute(input)
			if err == nil {
				results <- result
			}
			wg.Done()
		}(t)
	}

	return <-results, nil
}

// MapReduce takes any number of tasks, and feeds their results through reduce
// If no tasks are supplied, return an error.
// If any of the tasks error during their execution, return an error immediately.
// Once all tasks have completed successfully, return the value of reduce on
// their results in any order.
// Do not leave any pending goroutines. Make sure all goroutines are cleaned up
// properly and any synchronizing mechanisms closed.
func MapReduce(input int, reduce func(results []int) int, tasks ...Task) (int, error) {
	if len(tasks) == 0 {
		return 0, fmt.Errorf("no task is given")
	}

	resultCh := make(chan int, len(tasks))

	var wg sync.WaitGroup
	wg.Add(len(tasks))
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for _, t := range tasks {
		go func(t Task) {
			result, err := t.Execute(input)
			if err == nil {
				resultCh <- result
			}
			wg.Done()
		}(t)
	}

	results := make([]int, len(tasks))
	for result := range resultCh {
		results = append(results, result)
	}

	return reduce(results), nil
}
