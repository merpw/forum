package database_test

import (
	"testing"
	"time"
)

func TestOpsSessions(t *testing.T) {
	// add long session
	t.Run("AddLongSession", func(t *testing.T) {
		DB.AddSession("longToken", 9999999999, 1)
	})

	// add short session
	t.Run("AddShortSession", func(t *testing.T) {
		DB.AddSession("shortToken", 1, 2)
	})

	// pause code for 2 milliseconds
	t.Run("Pause 2 milliseconds", func(t *testing.T) {
		time.Sleep(2 * time.Millisecond)
	})

	// remove expired sessions
	t.Run("RemoveExpiredSessions", func(t *testing.T) {
		DB.RemoveExpiredSessions()
	})

	// check if long session is valid
	t.Run("CheckLongSessionIsStillValid", func(t *testing.T) {
		userId := DB.CheckSession("longToken")
		if userId != 1 {
			t.Fatalf("Expected userId 1, got %d", userId)
		}
	})

	// remove long session
	t.Run("RemoveLongSession", func(t *testing.T) {
		DB.RemoveSession("longToken")
	})

	// check if short session is already expired and removed
	t.Run("CheckShortSessionIsExpired", func(t *testing.T) {
		userId := DB.CheckSession("shortToken")
		if userId != -1 {
			t.Fatalf("Expected userId -1, got %d", userId)
		}
	})
}
