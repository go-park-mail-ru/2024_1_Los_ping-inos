package types

import "errors"

var ErrSeveralEmails = errors.New("pq: повторяющееся значение ключа нарушает ограничение уникальности \"person_email_key\"")
