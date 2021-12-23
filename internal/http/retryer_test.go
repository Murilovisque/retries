package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	retries "github.com/Murilovisque/retries/internal"
)

func TestShouldWorksWithStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()
	ht := &HttpRetryer{Client: ts.Client(), Retries: 3, TimeBetweenRetries: 0}
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ht.RequestWithExpectedStatus(req, 200)
	if err != nil {
		t.Fatal(err)
	}
}

func TestShouldWorksWithStatusSecondTime(t *testing.T) {
	var attempts int
	const expectedAttempts = 2
	mustFail := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mustFail {
			mustFail = false
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer ts.Close()
	ht := &HttpRetryer{Client: ts.Client(), Retries: 3, TimeBetweenRetries: 0}
	ht.ExecBeforeRequest = func(a int, req *http.Request) {
		attempts = a
	}
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ht.RequestWithExpectedStatus(req, 200)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != expectedAttempts {
		t.Fatalf("Expected %v attempts but %v", expectedAttempts, attempts)
	}
}

func TestShouldExceedAttemptsWithStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()
	ht := &HttpRetryer{Client: ts.Client(), Retries: 3, TimeBetweenRetries: 0}
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ht.RequestWithExpectedStatus(req, 200)
	if err != retries.ErrExceededAttempts {
		t.Fatal(err)
	}
}

func TestShouldWorksWithBody(t *testing.T) {
	const expectedBody = "ok"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedBody)
	}))
	defer ts.Close()
	ht := &HttpRetryer{Client: ts.Client(), Retries: 3, TimeBetweenRetries: 0}
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	bodyHandler := func(body io.ReadCloser) (bool, error) {
		bytesBody, err := ioutil.ReadAll(body)
		if err != nil {
			t.Fatal(err)
			return true, err
		}
		if string(bytesBody) != expectedBody {
			t.Fatalf("Expected %s, but %s", expectedBody, string(bytesBody))
		}
		return false, nil
	}
	_, err = ht.RequestWithExpectedStatusAndBody(req, bodyHandler, 200)
	if err != nil {
		t.Fatal(err)
	}
}

func TestShouldExceedAttemptsWithBody(t *testing.T) {
	const expectedBody = "falhou"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedBody)
	}))
	defer ts.Close()
	ht := &HttpRetryer{Client: ts.Client(), Retries: 3, TimeBetweenRetries: 0}
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	bodyHandler := func(body io.ReadCloser) (bool, error) {
		bytesBody, err := ioutil.ReadAll(body)
		if err != nil {
			t.Fatal(err)
			return true, err
		}
		if string(bytesBody) != expectedBody {
			t.Fatalf("Expected %s, but %s", expectedBody, string(bytesBody))
		}
		return true, nil
	}
	_, err = ht.RequestWithExpectedStatusAndBody(req, bodyHandler, 200)
	if err != retries.ErrExceededAttempts {
		t.Fatal(err)
	}
}

func TestShouldWorksWithBodySecondTime(t *testing.T) {
	expectedBodies := []string{"fail", "ok"}
	expectedRetries := []bool{true, false}
	var attempts int
	const expectedAttempts = 2
	responseIndex := -1
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseIndex++
		fmt.Fprint(w, expectedBodies[responseIndex])
	}))
	defer ts.Close()
	ht := &HttpRetryer{Client: ts.Client(), Retries: 3, TimeBetweenRetries: 0}
	ht.ExecBeforeRequest = func(a int, req *http.Request) {
		attempts = a
	}
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	bodyHandler := func(body io.ReadCloser) (bool, error) {
		bytesBody, err := ioutil.ReadAll(body)
		if err != nil {
			return true, err
		}
		if string(bytesBody) != expectedBodies[responseIndex] {
			t.Fatalf("Expected %s, but %s", expectedBodies[responseIndex], string(bytesBody))
		}
		return expectedRetries[responseIndex], nil
	}
	_, err = ht.RequestWithExpectedStatusAndBody(req, bodyHandler, 200)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != expectedAttempts {
		t.Fatalf("Expected %v attempts but %v", expectedAttempts, attempts)
	}
}
