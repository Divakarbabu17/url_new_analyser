package utils

import "time"

// Semaphore is a simple concurrency limiter
type Semaphore chan struct{}

// NewSemaphore creates a semaphore with max concurrency
func NewSemaphore(max int) Semaphore {
    return make(chan struct{}, max)
}

// Acquire acquires the semaphore
func (s Semaphore) Acquire() {
    s <- struct{}{}
}

// Release releases the semaphore
func (s Semaphore) Release() {
    <-s
}

// TimeoutChannel returns a channel that fires after given duration
func TimeoutChannel(d time.Duration) <-chan struct{} {
    ch := make(chan struct{})
    go func() {
        time.Sleep(d)
        close(ch)
    }()
    return ch
}