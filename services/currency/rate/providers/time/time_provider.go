package time

import "time"

type TimeProvider interface {
	Now() time.Time
}
