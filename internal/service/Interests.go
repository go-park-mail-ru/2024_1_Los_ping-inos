package service

import (
	"encoding/json"
)

func (service *Service) GetAllInterests(requestID int64) (string, error) {
	interests, err := service.interestStorage.Get(requestID, nil)
	if err != nil {
		return "", err
	}

	res, err := json.Marshal(interests)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
