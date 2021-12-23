package retries

import (
	"io"
	"net/http"
)

type Retryer interface {
	Do(work func() (bool, error)) error
	SetFuncExecBeforeTry(func(attempt int))
	SetFuncExecWhenWorkFailed(func(error))
}

type HttpRetryer interface {
	RequestWithExpectedStatus(req *http.Request, expectedstatus ...int) (*http.Response, error)
	RequestWithExpectedStatusAndBody(req *http.Request, bodyHandler func(io.ReadCloser) (bool, error), expectedstatus ...int) (*http.Response, error)
	SetFuncExecBeforeRequest(func(attempt int, req *http.Request))
	SetFuncExecWhenRequestFailed(func(error, *http.Response))
}
