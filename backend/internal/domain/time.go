//go:generate mockgen -source=time.go -destination=time_mocks.go -package=domain -mock_names=TimeProvider=MockTimeProvider

package domain

import "time"

type TimeProvider interface {
	Now() time.Time
}

type TimeNow func() time.Time

func (tn TimeNow) Now() time.Time {
	return tn().UTC()
}
