package retries

import (
	"net/http"
	"time"

	retries "github.com/Murilovisque/retries/internal"
	rhttp "github.com/Murilovisque/retries/internal/http"
)

func NewRetryer(qtdRetries int, timeBetweenRetries time.Duration) Retryer {
	return &retries.Retryer{Retries: qtdRetries, TimeBetweenRetries: timeBetweenRetries}
}

func NewHttpRetryer(client *http.Client, qtdRetries int, timeBetweenRetries time.Duration) HttpRetryer {
	r := &retries.Retryer{Retries: qtdRetries, TimeBetweenRetries: timeBetweenRetries}
	return &rhttp.HttpRetryer{Client: client, Retryer: r}
}
