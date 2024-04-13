package service

import (
	"context"
	"errors"
	"github.com/emirpasic/gods/sets/hashset"
	models "main.go/db"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
	"time"
)

func (service *Service) GetProfile(params ProfileGetParams, ctx context.Context) ([]models.Card, error) {
	persons, err := service.personStorage.Get(ctx, &models.PersonGetFilter{SessionID: params.SessionID, ID: params.ID})
	if err != nil {
		return nil, err
	}

	if len(persons) == 0 {
		return nil, errors.New("no such person")
	}

	interests, images, err := service.getUserCards(persons, ctx)
	if err != nil {
		return nil, err
	}

	if params.ID != nil { // сокрытие email'a чужого профиля
		persons[0].Email = ""
	}

	profile := combineToCards(persons, interests, images)
	if err != nil {
		return nil, err
	}

	profile[0].Email = persons[0].Email

	return profile, err
}

func (service *Service) UpdateProfile(SID string, profile requests.ProfileUpdateRequest, ctx context.Context) error {
	persons, err := service.personStorage.Get(ctx, &models.PersonGetFilter{SessionID: []string{SID}})
	if err != nil {
		return err
	}
	person := persons[0]
	if profile.Name != "" {
		person.Name = profile.Name
	}
	if profile.Email != "" {
		if err = checkPassword(person.Password, profile.OldPassword); err != nil {
			return err
		}
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
		if err = checkPassword(person.Password, profile.OldPassword); err != nil {
			return err
		}
		person.Password, err = hashPassword(profile.Password)
		if err != nil {
			return err
		}
	}
	if profile.Interests != nil {
		err = service.handleInterests(profile.Interests, person.ID, ctx)
		if err != nil {
			return err
		}
	}
	err = service.personStorage.Update(ctx, *person)

	return err
}

func (service *Service) DeleteProfile(sessionID string, ctx context.Context) error {
	err := service.personStorage.Delete(ctx, sessionID)
	return err
}

func (service *Service) handleInterests(interests []string, personID types.UserID, ctx context.Context) error {
	interestsBefore, err := service.interestStorage.GetPersonInterests(ctx, personID)
	if err != nil {
		return err
	}

	interestsAfter, err := service.interestStorage.Get(ctx, &models.InterestGetFilter{Name: interests})
	if err != nil {
		return err
	}

	beforeIDs := getInterestIDs(interestsBefore)
	afterIDs := getInterestIDs(interestsAfter)

	setBefore := hashset.New(beforeIDs...)
	setAfter := hashset.New(afterIDs...)

	ad := setAfter.Difference(setBefore).Values()
	add := normalizeFromSet(ad)
	err = service.interestStorage.CreatePersonInterests(ctx, personID, add)
	if err != nil {
		return err
	}

	del := setBefore.Difference(setAfter).Values()
	delet := normalizeFromSet(del)
	err = service.interestStorage.DeletePersonInterests(ctx, personID, delet)

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
