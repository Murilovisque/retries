package retries

import (
	"fmt"
	"testing"
)

func TestShouldWorksFirstTime(t *testing.T) {
	r := &Retryer{Retries: 3, TimeBetweenRetries: 0}
	valueToChange := false
	const valueExpected = true
	err := r.Do(func() (bool, error) {
		valueToChange = valueExpected
		return false, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !valueToChange {
		t.Fatalf("Value should be %v", valueExpected)
	}
}

func TestShouldWorksSecondTime(t *testing.T) {
	var attempts int
	const expectedAttempts = 2
	jobCounter := 1
	r := &Retryer{Retries: 3, TimeBetweenRetries: 0}
	r.ExecBeforeTry = func(a int) {
		attempts = a
	}

	err := r.Do(func() (bool, error) {
		if jobCounter < expectedAttempts {
			jobCounter++
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if attempts != expectedAttempts {
		t.Fatalf("Expected %v attempts but %v", expectedAttempts, attempts)
	}
}

func TestShouldExceedAttempts(t *testing.T) {
	r := &Retryer{Retries: 3, TimeBetweenRetries: 0}
	err := r.Do(func() (bool, error) {
		return true, nil
	})
	if err != errExceededAttempts {
		t.Fatal(err)
	}
}

func TestShouldInformAndReturnErrors(t *testing.T) {
	var attempts int
	var err error
	r := &Retryer{Retries: 3, TimeBetweenRetries: 0}
	r.ExecWhenWorkFailed = func(errWork error) {
		if errWork != err {
			t.Fatal(errWork)
		}
	}
	lastErr := r.Do(func() (bool, error) {
		attempts++
		err = fmt.Errorf("Attempt %d failed", attempts)
		return true, err
	})
	if lastErr != err {
		t.Fatal(err)
	}
}
