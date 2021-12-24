package http

import (
	"io"
	"net/http"
	"time"

	retries "github.com/Murilovisque/retries/internal"
)

type HttpRetryer struct {
	Client                *http.Client
	ExecBeforeRequest     func(attempt int, req *http.Request)
	ExecWhenRequestFailed func(error, *http.Request, *http.Response)
	Retries               int
	TimeBetweenRetries    time.Duration
}

func (hr *HttpRetryer) RequestWithExpectedStatus(req *http.Request, expectedstatus ...int) (*http.Response, error) {
	return hr.RequestWithExpectedStatusAndBody(req, nil, expectedstatus...)
}

func (hr *HttpRetryer) RequestWithExpectedStatusAndBody(req *http.Request, bodyHandler func(io.ReadCloser) (bool, error), expectedstatus ...int) (*http.Response, error) {
	attempt := 1
	for {
		if hr.ExecBeforeRequest != nil {
			hr.ExecBeforeRequest(attempt, req)
		}
		res, err := hr.Client.Do(req)
		attempt++
		if err == nil {
			for _, expStatus := range expectedstatus {
				if expStatus == res.StatusCode {
					if bodyHandler == nil {
						return res, nil
					}
					shouldRetry, err := bodyHandler(res.Body)
					res.Body.Close()
					if !shouldRetry && err == nil {
						return res, nil
					}
					break
				}
			}
		}
		if (err != nil || res != nil) && hr.ExecWhenRequestFailed != nil {
			hr.ExecWhenRequestFailed(err, req, res)
		}
		if attempt > hr.Retries {
			if err == nil {
				return res, retries.ErrExceededAttempts
			}
			return res, err
		}
		time.Sleep(hr.TimeBetweenRetries)
	}
}

func (hr *HttpRetryer) SetFuncExecBeforeRequest(f func(attempt int, req *http.Request)) {
	hr.ExecBeforeRequest = f
}

func (hr *HttpRetryer) SetFuncExecWhenRequestFailed(f func(error, *http.Request, *http.Response)) {
	hr.ExecWhenRequestFailed = f
}
