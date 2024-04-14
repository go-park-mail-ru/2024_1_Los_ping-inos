package usecase

import (
	"context"
	"errors"
	"github.com/emirpasic/gods/sets/hashset"
	"golang.org/x/crypto/bcrypt"
	requests "main.go/internal/pkg"
	"main.go/internal/profile"
	"main.go/internal/types"
	"time"
)

type UseCase struct {
	personStorage   profile.PersonStorage
	interestStorage profile.InterestStorage
	imageStorage    profile.ImageStorage
	likeStorage     profile.LikeStorage
}

func NewProfileUseCase(pStorage profile.PersonStorage, intStorage profile.InterestStorage, imgStorage profile.ImageStorage, lStorage profile.LikeStorage) *UseCase {
	return &UseCase{
		personStorage:   pStorage,
		interestStorage: intStorage,
		imageStorage:    imgStorage,
		likeStorage:     lStorage,
	}
}

func (service *UseCase) GetProfile(params profile.ProfileGetParams, ctx context.Context) ([]profile.Card, error) {
	persons, err := service.personStorage.Get(ctx, &profile.PersonGetFilter{SessionID: params.SessionID, ID: params.ID})
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

	prof := combineToCards(persons, interests, images)
	if err != nil {
		return nil, err
	}

	prof[0].Email = persons[0].Email

	return prof, err
}

func (service *UseCase) UpdateProfile(SID string, prof requests.ProfileUpdateRequest, ctx context.Context) error {
	persons, err := service.personStorage.Get(ctx, &profile.PersonGetFilter{SessionID: []string{SID}})
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

func (service *UseCase) DeleteProfile(sessionID string, ctx context.Context) error {
	err := service.personStorage.Delete(ctx, sessionID)
	return err
}

func (service *UseCase) handleInterests(interests []string, personID types.UserID, ctx context.Context) error {
	interestsBefore, err := service.interestStorage.GetPersonInterests(ctx, personID)
	if err != nil {
		return err
	}

	interestsAfter, err := service.interestStorage.Get(ctx, &profile.InterestGetFilter{Name: interests})
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

func (service *UseCase) GetMatches(prof types.UserID, ctx context.Context) ([]profile.Card, error) {
	ids, err := service.likeStorage.GetMatch(ctx, prof)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	return service.GetProfile(profile.ProfileGetParams{ID: ids}, ctx)
}

func (service *UseCase) getUserCards(persons []*profile.Person, ctx context.Context) ([][]*profile.Interest, [][]profile.Image, error) {
	var err error
	interests := make([][]*profile.Interest, len(persons))
	images := make([][]profile.Image, len(persons))
	for j := range persons {
		interests[j], err = service.interestStorage.GetPersonInterests(ctx, persons[j].ID)
		if err != nil {
			return nil, nil, err
		}
		images[j], err = service.imageStorage.Get(ctx, int64(persons[j].ID))
		if err != nil {
			return nil, nil, err
		}
	}
	return interests, images, nil
}

func combineToCards(persons []*profile.Person, interests [][]*profile.Interest, images [][]profile.Image) []profile.Card {
	if len(persons) != len(interests) || len(persons) != len(images) {
		return nil
	}

	photos := make([][]profile.ImageToSend, len(persons))
	for i := range images {
		photos[i] = make([]profile.ImageToSend, len(images[i]))
		for j, image := range images[i] {
			photos[i][j] = profile.ImageToSend{
				Cell: image.CellNumber,
				Url:  image.Url,
			}
		}
	}

	res := make([]profile.Card, len(persons))
	for i := range persons {
		res[i] = profile.Card{Name: persons[i].Name, Birthday: persons[i].Birthday, Description: persons[i].Description,
			Interests: interests[i], Photos: photos[i]}
	}
	return res
}

func normalizeFromSet(input []interface{}) []types.InterestID {
	res := make([]types.InterestID, len(input))
	for i := range input {
		res[i] = input[i].(types.InterestID)
	}
	return res
}

func getInterestIDs(interests []*profile.Interest) []interface{} {
	res := make([]interface{}, len(interests))
	for i := range interests {
		res[i] = interests[i].ID
	}
	return res
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword - принимает hash - захэшированный пароль из базы и проверяет, соответствует ли ему password
func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
