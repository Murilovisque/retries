package retries

import (
	"errors"
	"time"
)

var (
	errExceededAttempts = errors.New("exceeded attempts")
)

func NewRetryer(retries int, timeBetweenRetries time.Duration) *Retryer {
	return &Retryer{retries: retries, timeBetweenRetries: timeBetweenRetries}
}

type Retryer struct {
	ExecBeforeTry      func(attempt int)
	ExecWhenWorkFailed func(error)
	retries            int
	timeBetweenRetries time.Duration
}

func (r *Retryer) Do(work func() (bool, error)) error {
	attempt := 1
	for {
		if r.ExecBeforeTry != nil {
			r.ExecBeforeTry(attempt)
		}
		shouldRetry, err := work()
		attempt++
		if !shouldRetry && err == nil {
			return nil
		}
		if err != nil && r.ExecWhenWorkFailed != nil {
			r.ExecWhenWorkFailed(err)
		}
		if attempt > r.retries {
			if err == nil {
				return errExceededAttempts
			}
			return err
		}
		time.Sleep(r.timeBetweenRetries)
	}
}
