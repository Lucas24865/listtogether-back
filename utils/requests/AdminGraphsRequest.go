package requests

import "time"

type AdminGraphRequest struct {
	From time.Time
	To   time.Time
}
