package types

import "errors"

var SeveralEmailsError = errors.New("pq: повторяющееся значение ключа нарушает ограничение уникальности \"person_email_key\"")
var DifferentPasswordsError = errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password")
