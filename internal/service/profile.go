package service

import (
	"errors"
	"github.com/emirpasic/gods/sets/hashset"
	models "main.go/db"
	"main.go/internal/types"
	"time"
)

func (service *Service) GetProfile(sessionID string, requestID int64) (string, error) {
	person, err := service.personStorage.Get(requestID, &models.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil {
		return "", err
	}

	if len(person) == 0 {
		return "", errors.New("no person with such sessionID")
	}

	interests, err := service.interestStorage.GetPersonInterests(requestID, person[0].ID)
	if err != nil {
		return "", err
	}

	res, err := personsToJSON(person, [][]*models.Interest{interests})
	if err != nil {
		return "", err
	}
	return res, err
}

func (service *Service) UpdateProfile(profile ProfileUpdate, requestID int64) error {
	persons, err := service.personStorage.Get(requestID, &models.PersonGetFilter{SessionID: []string{profile.SessionID}})
	if err != nil {
		return err
	}
	person := persons[0]
	if profile.Name != "" {
		person.Name = profile.Name
	}
	if profile.Email != "" {
		person.Email = profile.Email
	}
	if profile.Birthday != "" {
		person.Birthday, err = time.Parse("01.02.2006", profile.Birthday)
		if err != nil {
			return err
		}
	}
	if profile.Description != "" {
		person.Description = profile.Description
	}
	if profile.Password != "" {
		person.Password, err = hashPassword(profile.Password)
		if err != nil {
			return err
		}
	}
	if profile.Interests != nil {
		err = service.handleInterests(profile.Interests, person.ID, requestID)
		if err != nil {
			return err
		}
	}
	err = service.personStorage.Update(requestID, *person)

	return err
}

func (service *Service) DeleteProfile(sessionID string, requestID int64) error {
	err := service.personStorage.Delete(requestID, sessionID)
	return err
}

func (service *Service) handleInterests(interests []string, personID types.UserID, requestID int64) error {
	interestsBefore, err := service.interestStorage.GetPersonInterests(requestID, personID)
	if err != nil {
		return err
	}

	interestsAfter, err := service.interestStorage.Get(requestID, &models.InterestGetFilter{Name: interests})
	if err != nil {
		return err
	}

	beforeIDs := getInterestIDs(interestsBefore)
	afterIDs := getInterestIDs(interestsAfter)

	setBefore := hashset.New(beforeIDs...)
	setAfter := hashset.New(afterIDs...)

	ad := setAfter.Difference(setBefore).Values()
	add := normalizeFromSet(ad)
	err = service.interestStorage.CreatePersonInterests(requestID, personID, add)
	if err != nil {
		return err
	}

	del := setBefore.Difference(setAfter).Values()
	delet := normalizeFromSet(del)
	err = service.interestStorage.DeletePersonInterests(requestID, personID, delet)

	return err
}

func normalizeFromSet(input []interface{}) []types.InterestID {
	res := make([]types.InterestID, len(input))
	for i := range input {
		res[i] = input[i].(types.InterestID)
	}
	return res
}

func getInterestIDs(interests []*models.Interest) []interface{} {
	res := make([]interface{}, len(interests))
	for i := range interests {
		res[i] = interests[i].ID
	}
	return res
}
