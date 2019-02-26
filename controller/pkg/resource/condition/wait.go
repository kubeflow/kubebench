// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package condition

import (
	"errors"
	"math/rand"
	"time"
)

// ConditionFunc returns true if the condition is satisfied, or an error
// if the loop should be aborted.
type ConditionFunc func() (done bool, err error)

// ErrWaitTimeout is returned when the condition exited without success.
var ErrWaitTimeout = errors.New("timed out waiting for the condition")

// Backoff holds parameters applied to a Backoff function.
type Backoff struct {
	// The initial duration.
	Duration time.Duration
	// Duration is multiplied by factor each iteration. Must be greater
	// than or equal to zero.
	Factor float64
	// The amount of jitter applied each iteration. Jitter is applied after
	// cap.
	Jitter float64
	// The number of steps before duration stops changing. If zero, initial
	// duration is always used. Used for exponential backoff in combination
	// with Factor.
	Steps int
	// The returned duration will never be greater than cap *before* jitter
	// is applied. The actual maximum cap is `cap * (1.0 + jitter)`.
	Cap time.Duration
}

// Step returns the next interval in the exponential backoff. This method
// will mutate the provided backoff.
func (b *Backoff) Step() time.Duration {
	if b.Steps < 1 {
		return b.Duration
	}
	b.Steps--

	duration := b.Duration

	// calculate the next step
	if b.Factor != 0 {
		b.Duration = time.Duration(float64(b.Duration) * b.Factor)
		if b.Cap > 0 && b.Duration > b.Cap {
			b.Duration = b.Cap
			b.Steps = 0
		}
	}
	if b.Jitter > 0 {
		duration = Jitter(duration, b.Jitter)
	}

	return duration
}

// Jitter returns a time.Duration between duration and duration + maxFactor *
// duration.
func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	wait := duration + time.Duration(rand.Float64()*maxFactor*float64(duration))
	return wait
}

// ExponentialBackoff repeats a condition check with exponential backoff.
func ExponentialBackoff(backoff Backoff, timeout time.Duration, condition ConditionFunc) error {
	totalDuration := 0.0
	for {
		if ok, err := condition(); err != nil || ok {
			return err
		}
		// if timeout <= 0 then no timeout is applied
		if timeout > 0 && float64(totalDuration) >= float64(timeout) {
			break
		}
		time.Sleep(backoff.Step())
		totalDuration += float64(backoff.Duration)
	}
	return ErrWaitTimeout
}
