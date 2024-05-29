package types

import (
	"errors"
	"testing"
)

func TestSeveralEmailsError(t *testing.T) {
	expectedErrorMessage := "pq: повторяющееся значение ключа нарушает ограничение уникальности \"person_email_key\""

	if SeveralEmailsError.Error() != expectedErrorMessage {
		t.Errorf("expected error message: %v, got: %v", expectedErrorMessage, SeveralEmailsError.Error())
	}
}

func TestDifferentPasswordsError(t *testing.T) {
	expectedErrorMessage := "crypto/bcrypt: hashedPassword is not the hash of the given password"

	if DifferentPasswordsError.Error() != expectedErrorMessage {
		t.Errorf("expected error message: %v, got: %v", expectedErrorMessage, DifferentPasswordsError.Error())
	}
}

func TestMyErr_Error(t *testing.T) {
	customErr := errors.New("custom error message")
	myErr := MyErr{Err: customErr}

	expectedErrorMessage := "custom error message"
	if myErr.Error() != expectedErrorMessage {
		t.Errorf("expected error message: %v, got: %v", expectedErrorMessage, myErr.Error())
	}
}
