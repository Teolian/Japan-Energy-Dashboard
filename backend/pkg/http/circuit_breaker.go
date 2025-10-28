// Package http provides circuit breaker implementation.
package http

import (
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the state of the circuit breaker.
type CircuitState int

const (
	// StateClosed: circuit is closed, requests flow normally
	StateClosed CircuitState = iota
	// StateOpen: circuit is open, requests are blocked
	StateOpen
	// StateHalfOpen: circuit is testing if service recovered
	StateHalfOpen
)

// CircuitBreaker prevents cascading failures by tracking consecutive errors.
type CircuitBreaker struct {
	maxFailures     int           // Max consecutive failures before opening
	timeout         time.Duration // Time to wait before trying again (half-open)
	state           CircuitState
	failures        int
	lastFailureTime time.Time
	mu              sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker.
// maxFailures: number of consecutive failures before opening the circuit
// timeout: duration to wait before attempting recovery (half-open state)
func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures: maxFailures,
		timeout:     timeout,
		state:       StateClosed,
	}
}

// Call wraps a function call with circuit breaker logic.
// Returns error if circuit is open.
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Check if we should transition from Open to HalfOpen
	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = StateHalfOpen
			cb.failures = 0
		} else {
			return fmt.Errorf("circuit breaker is open (last failure: %v ago)",
				time.Since(cb.lastFailureTime).Round(time.Second))
		}
	}

	// Execute the function
	err := fn()

	if err != nil {
		// Record failure
		cb.failures++
		cb.lastFailureTime = time.Now()

		// Open circuit if max failures reached
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
			return fmt.Errorf("circuit breaker opened after %d failures: %w", cb.failures, err)
		}

		return err
	}

	// Success - reset circuit
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
	}
	cb.failures = 0

	return nil
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Reset manually resets the circuit breaker to closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failures = 0
}

// Failures returns the current number of consecutive failures.
func (cb *CircuitBreaker) Failures() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}
