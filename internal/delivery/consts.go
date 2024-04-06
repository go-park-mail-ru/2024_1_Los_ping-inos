package delivery

import "time"

var expiredYear = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

func oneDayExpiration() time.Time { return time.Now().Add(24 * time.Hour) }
