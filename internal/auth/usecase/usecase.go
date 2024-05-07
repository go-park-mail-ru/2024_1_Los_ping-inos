package usecase

import (
	"cmp"
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"main.go/internal/auth"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
	"slices"
	"time"
)

type UseCase struct {
	personStorage   auth.PersonStorage
	sessionStorage  auth.SessionStorage
	interestStorage auth.InterestStorage
	imageStorage    auth.ImageStorage
}

func NewAuthUseCase(dbReader auth.PersonStorage, sstore auth.SessionStorage, istore auth.InterestStorage, imgStore auth.ImageStorage) *UseCase {
	return &UseCase{
		personStorage:   dbReader,
		sessionStorage:  sstore,
		interestStorage: istore,
		imageStorage:    imgStore,
	}
}

func (service *UseCase) IsAuthenticated(sessionID string, ctx context.Context) (types.UserID, bool, error) {
	defer requests.TrackContextTimings(ctx, "IsAuthUc", time.Now())

	person, err := service.sessionStorage.GetBySID(ctx, sessionID)
	if err != nil {
		return -1, false, err
	}
	return person.UID, true, nil
}

// Login - принимает email, пароль; возвращает ID сессии и ошибку
func (service *UseCase) Login(email, password string, ctx context.Context) (*auth.Profile, string, error) {
	defer requests.TrackContextTimings(ctx, "LoginUc", time.Now())

	users, ok := service.personStorage.Get(ctx, &auth.PersonGetFilter{Email: []string{email}})
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
	err = service.sessionStorage.CreateSession(ctx, auth.Session{
		UID: user.ID,
		SID: SID,
	})
	if err != nil {
		return nil, "", err
	}

	interests, images, err := service.getUserCards(users, ctx)
	if err != nil {
		return nil, "", err
	}
	profiles := combineToCards(users, interests, images)
	return &profiles[0], SID, nil
}

func (service *UseCase) GetAllInterests(ctx context.Context) ([]*auth.Interest, error) {
	defer requests.TrackContextTimings(ctx, "GetAllInterestsUc", time.Now())

	return service.interestStorage.Get(ctx, nil)
}

func (service *UseCase) Registration(body auth.RegitstrationBody, ctx context.Context) (*auth.Profile, string, error) {
	defer requests.TrackContextTimings(ctx, "RegistrationUc", time.Now())

	hashedPassword, err := hashPassword(body.Password)
	if err != nil {
		return nil, "", err
	}

	err = service.personStorage.AddAccount(ctx, body.Name, body.Birthday, body.Gender, body.Email, hashedPassword)
	if err != nil {
		return nil, "", err
	}

	prof, SID, err := service.Login(body.Email, body.Password, ctx)
	if err != nil {
		return nil, "", err
	}

	interests, err := service.interestStorage.Get(ctx, &auth.InterestGetFilter{Name: body.Interests})
	if err != nil {
		return nil, "", err
	}
	err = service.interestStorage.CreatePersonInterests(ctx, prof.ID, getInterestIDs(interests))
	if err != nil {
		return nil, "", err
	}

	prof.Interests = interests
	return prof, SID, nil
}

func (service *UseCase) GetName(userID types.UserID, ctx context.Context) (string, error) {
	defer requests.TrackContextTimings(ctx, "GetNameUc", time.Now())

	person, err := service.personStorage.Get(ctx, &auth.PersonGetFilter{ID: []types.UserID{userID}})
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

func (service *UseCase) Logout(sessionID string, ctx context.Context) error {
	defer requests.TrackContextTimings(ctx, "LogoutUc", time.Now())

	return service.sessionStorage.DeleteSession(ctx, sessionID)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword - принимает hash - захэшированный пароль из базы и проверяет, соответствует ли ему password
func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (service *UseCase) getUserCards(persons []*auth.Person, ctx context.Context) ([][]*auth.Interest, [][]auth.Image, error) {
	var err error
	interests := make([][]*auth.Interest, len(persons))
	images := make([][]auth.Image, len(persons))

	// TODO
	//interests, images, err := service.personStorage.GetUserCards(ctx, getUserIDs(persons))
	//
	//if err != nil {
	//	return nil, nil, err
	//}

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

func getUserIDs(persons []*auth.Person) []types.UserID {
	res := make([]types.UserID, len(persons))
	for i := range persons {
		res[i] = persons[i].ID
	}
	return res
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
		slices.SortFunc(photos[i], func(a, b auth.ImageToSend) int { // TODO
			return cmp.Compare(a.Cell, b.Cell)
		})
		res[i] = auth.Profile{ID: persons[i].ID, Name: persons[i].Name, Birthday: persons[i].Birthday, Description: persons[i].Description,
			Email: persons[i].Email, Interests: interests[i], Photos: photos[i]}
	}
	return res
}
