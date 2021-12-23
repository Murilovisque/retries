package retries

import (
	"errors"
	"time"
)

var (
	ErrExceededAttempts = errors.New("exceeded attempts")
)

type Retryer struct {
	ExecBeforeTry      func(attempt int)
	ExecWhenWorkFailed func(error)
	Retries            int
	TimeBetweenRetries time.Duration
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
		if attempt > r.Retries {
			if err == nil {
				return ErrExceededAttempts
			}
			return err
		}
		time.Sleep(r.TimeBetweenRetries)
	}
}

func (r *Retryer) SetFuncExecBeforeTry(f func(attempt int)) {
	r.ExecBeforeTry = f
}

func (r *Retryer) SetFuncExecWhenWorkFailed(f func(error)) {
	r.ExecWhenWorkFailed = f
}
