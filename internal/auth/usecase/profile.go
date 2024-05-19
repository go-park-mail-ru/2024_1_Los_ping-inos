package usecase

import (
	"context"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"main.go/internal/auth"
	"main.go/internal/types"
)

func (service *UseCase) GetProfile(params auth.ProfileGetParams, ctx context.Context) ([]auth.Profile, error) {
	persons, err := service.personStorage.Get(ctx, &auth.PersonGetFilter{SessionID: params.SessionID, ID: params.ID, Name: params.Name})
	if err != nil {
		return nil, err
	}

	if len(persons) == 0 {
		return []auth.Profile{}, nil
	}

	interests, images, err := service.getUserCards(persons, ctx)
	if err != nil {
		return nil, err
	}

	if params.ID != nil { // сокрытие email'a чужого профиля
		persons[0].Email = ""
	}

	prof := combineToCards(persons, interests, images)
	if err != nil {
		return nil, err
	}

	prof[0].Email = persons[0].Email

	return prof, err
}

func (service *UseCase) UpdateProfile(UID types.UserID, prof auth.ProfileUpdateRequest, ctx context.Context) error {
	persons, err := service.personStorage.Get(ctx, &auth.PersonGetFilter{ID: []types.UserID{UID}})
	if err != nil {
		return err
	}
	person := persons[0]
	if prof.Name != "" {
		person.Name = prof.Name
	}
	if prof.Email != "" {
		if err = checkPassword(person.Password, prof.OldPassword); err != nil {
			return err
		}
		person.Email = prof.Email
	}
	if prof.Birthday != "" {
		person.Birthday, err = time.Parse("01.02.2006", prof.Birthday)
		if err != nil {
			return err
		}
	}
	if prof.Description != "" {
		person.Description = prof.Description
	}
	if prof.Password != "" {
		if err = checkPassword(person.Password, prof.OldPassword); err != nil {
			return err
		}
		person.Password, err = hashPassword(prof.Password)
		if err != nil {
			return err
		}
	}
	if prof.Interests != nil {
		err = service.handleInterests(prof.Interests, person.ID, ctx)
		if err != nil {
			return err
		}
	}
	err = service.personStorage.Update(ctx, *person)

	return err
}

func (service *UseCase) DeleteProfile(UID types.UserID, ctx context.Context) error {
	return service.personStorage.Delete(ctx, UID)
}

func (service *UseCase) handleInterests(interests []string, personID types.UserID, ctx context.Context) error {
	interestsBefore, err := service.interestStorage.GetPersonInterests(ctx, personID)
	if err != nil {
		return err
	}

	interestsAfter, err := service.interestStorage.GetInterest(ctx, &auth.InterestGetFilter{Name: interests})
	if err != nil {
		return err
	}

	beforeIDs := getInterestSetIDs(interestsBefore)
	afterIDs := getInterestSetIDs(interestsAfter)

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

func getInterestSetIDs(interests []*auth.Interest) []interface{} {
	res := make([]interface{}, len(interests))
	for i := range interests {
		res[i] = interests[i].ID
	}
	return res
}

func (service *UseCase) GetMatches(prof types.UserID, nameFilter string, ctx context.Context) ([]auth.Profile, error) {
	ids, err := service.personStorage.GetMatch(ctx, prof)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return make([]auth.Profile, 0), nil
	}

	return service.GetProfile(auth.ProfileGetParams{ID: ids, Name: nameFilter}, ctx)
}

func normalizeFromSet(input []interface{}) []types.InterestID {
	res := make([]types.InterestID, len(input))
	for i := range input {
		res[i] = input[i].(types.InterestID)
	}
	return res
}
