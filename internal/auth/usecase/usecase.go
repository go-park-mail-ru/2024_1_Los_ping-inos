package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"main.go/internal/auth"
	"main.go/internal/types"
)

type UseCase struct {
	dbReader        auth.PostgresRepo
	interestStorage auth.InterestStorage
	imageStorage    auth.ImageStorage
}

func NewAuthUseCase(dbReader auth.PostgresRepo, istore auth.InterestStorage, imgStore auth.ImageStorage) *UseCase {
	return &UseCase{
		dbReader:        dbReader,
		interestStorage: istore,
		imageStorage:    imgStore,
	}
}

func (api *UseCase) IsAuthenticated(sessionID string, ctx context.Context) (types.UserID, bool) {
	person, err := api.dbReader.Get(ctx, &auth.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil || len(person) == 0 {
		return -1, false
	}
	return person[0].ID, true
}

// Login - принимает email, пароль; возвращает ID сессии и ошибку
func (api *UseCase) Login(email, password string, ctx context.Context) (*auth.Profile, string, error) {
	ems := make([]string, 1)
	ems[0] = email
	users, ok := api.dbReader.Get(ctx, &auth.PersonGetFilter{Email: ems})
	if ok != nil {
		return nil, "", ok
	}

	if len(users) == 0 {
		return nil, "", errors.New("no such person")
	}

	user := users[0]
	err := checkPassword(user.Password, password)

	if err != nil {
		return nil, "", err
	}

	SID := uuid.NewString()
	user.SessionID = SID
	err = api.dbReader.Update(ctx, *user)
	if err != nil {
		return nil, "", err
	}

	interests, images, err := api.getUserCards(users, ctx)
	if err != nil {
		return nil, "", err
	}
	profiles := combineToCards(users, interests, images)
	return &profiles[0], SID, nil
}

func (api *UseCase) GetAllInterests(ctx context.Context) ([]*auth.Interest, error) {
	return api.interestStorage.Get(ctx, nil)
}

func (api *UseCase) Registration(body auth.RegitstrationBody, ctx context.Context) (*auth.Profile, string, error) {
	hashedPassword, err := hashPassword(body.Password)
	if err != nil {
		return nil, "", err
	}

	err = api.dbReader.AddAccount(ctx, body.Name, body.Birthday, body.Gender, body.Email, hashedPassword)
	if err != nil {
		return nil, "", err
	}

	prof, SID, err := api.Login(body.Email, body.Password, ctx)
	if err != nil {
		return nil, "", err
	}

	interests, err := api.interestStorage.Get(ctx, &auth.InterestGetFilter{Name: body.Interests})
	if err != nil {
		return nil, "", err
	}
	err = api.interestStorage.CreatePersonInterests(ctx, prof.ID, getInterestIDs(interests))
	if err != nil {
		return nil, "", err
	}

	prof.Interests = interests
	return prof, SID, nil
}

func (api *UseCase) GetName(sessionID string, ctx context.Context) (string, error) {
	person, err := api.dbReader.Get(ctx, &auth.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil {
		return "", err
	}

	if len(person) == 0 {
		return "", errors.New("no person with such sessionID")
	}

	return person[0].Name, err
}

func getInterestIDs(interests []*auth.Interest) []types.InterestID {
	res := make([]types.InterestID, len(interests))
	for i := range interests {
		res[i] = interests[i].ID
	}
	return res
}

func (api *UseCase) Logout(sessionID string, ctx context.Context) error {
	err := api.dbReader.RemoveSession(ctx, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword - принимает hash - захэшированный пароль из базы и проверяет, соответствует ли ему password
func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (api *UseCase) getUserCards(persons []*auth.Person, ctx context.Context) ([][]*auth.Interest, [][]auth.Image, error) {
	var err error
	interests := make([][]*auth.Interest, len(persons))
	images := make([][]auth.Image, len(persons))
	for j := range persons {
		interests[j], err = api.interestStorage.GetPersonInterests(ctx, persons[j].ID)
		if err != nil {
			return nil, nil, err
		}
		images[j], err = api.imageStorage.Get(ctx, int64(persons[j].ID))
		if err != nil {
			return nil, nil, err
		}
	}
	return interests, images, nil
}

func combineToCards(persons []*auth.Person, interests [][]*auth.Interest, images [][]auth.Image) []auth.Profile {
	if len(persons) != len(interests) || len(persons) != len(images) {
		return nil
	}

	photos := make([][]auth.ImageToSend, len(persons))
	for i := range images {
		photos[i] = make([]auth.ImageToSend, len(images[i]))
		for j, image := range images[i] {
			photos[i][j] = auth.ImageToSend{
				Cell: image.CellNumber,
				Url:  image.Url,
			}
		}
	}

	res := make([]auth.Profile, len(persons))
	for i := range persons {
		res[i] = auth.Profile{ID: persons[i].ID, Name: persons[i].Name, Birthday: persons[i].Birthday, Description: persons[i].Description,
			Interests: interests[i], Photos: photos[i]}
	}
	return res
}
