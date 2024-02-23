package service

import (
	"main.go/internal/storage"
	"main.go/internal/types"
)

func GetCoolIdsList() ([]string, error) {
	res := make([]string, 0)
	eStorage := storage.ExampleStorage{}
	var i types.UserID
	for ; i < 5; i++ {
		res = append(res, eStorage.Get(i))
	}
	return res, nil
}
