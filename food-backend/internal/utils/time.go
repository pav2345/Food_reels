package utils

import "time"

// FormatMongoDate formats a time value the same way Mongoose serializes dates in JSON.
func FormatMongoDate(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}
