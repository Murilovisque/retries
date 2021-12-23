package http

import (
	"fmt"
	"io"
	"net/http"

	retries "github.com/Murilovisque/retries/internal"
)

type HttpRetryer struct {
	*retries.Retryer
	Client *http.Client
}

func (hr *HttpRetryer) RequestWithExpectedStatus(req *http.Request, expectedstatus ...int) (*http.Response, error) {
	var resToReturn *http.Response
	err := hr.Retryer.Do(func() (bool, error) {
		res, err := hr.Client.Do(req)
		if err != nil {
			return true, err
		}
		for _, expStatus := range expectedstatus {
			if expStatus == res.StatusCode {
				resToReturn = res
				return false, nil
			}
		}
		return true, fmt.Errorf("expected a status between %v , but received %d", expectedstatus, res.StatusCode)
	})
	return resToReturn, err
}

func (hr *HttpRetryer) RequestWithExpectedStatusAndBody(req *http.Request, bodyHandler func(io.ReadCloser) (bool, error), expectedstatus ...int) (*http.Response, error) {
	var resToReturn *http.Response
	err := hr.Retryer.Do(func() (bool, error) {
		res, err := hr.Client.Do(req)
		if err != nil {
			return true, err
		}
		for _, expStatus := range expectedstatus {
			if expStatus == res.StatusCode {
				shouldRetry, err := bodyHandler(res.Body)
				res.Body.Close()
				if !shouldRetry && err == nil {
					resToReturn = res
					return false, nil
				}
				return shouldRetry, err
			}
		}
		return true, fmt.Errorf("expected a status between %v , but received %d", expectedstatus, res.StatusCode)
	})
	return resToReturn, err
}
