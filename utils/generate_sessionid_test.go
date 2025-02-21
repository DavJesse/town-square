package utils

import (
	"github.com/gofrs/uuid"
	"testing"
)

/*
* Tests the GenerateSessionID function. These are the cases:
* 1. Generation of a valid UUID
* 2. Check if the session ID is a valid UUID
* 3. Ensure that multiple calls generate different UUIDs
 */
func TestGenerateSessionID(t *testing.T) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = uuid.FromString(sessionID)
	if err != nil {
		t.Errorf("expected valid UUID, got error: %v", err)
	}

	sessionID2, err := GenerateSessionID()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if sessionID == sessionID2 {
		t.Errorf("expected different session IDs, but got the same")
	}
}
