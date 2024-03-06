package service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func (service *Service) GetAllInterests() (string, error) {
	interests, err := service.interestStorage.Get()
	if err != nil {
		logrus.Info("can't get interests")
		return "", err
	}

	res, err := json.Marshal(interests)
	if err != nil {
		logrus.Info("can't marshal interests")
		return "", err
	}

	return string(res), nil
}
