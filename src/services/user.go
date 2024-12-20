package services

import (
	"errors"
	"md2s/dto"
	"md2s/models"
	"md2s/repositorys"
)

func GetUsers() ([]models.User, error) {
	return repositorys.GetUsers()
}

func GetUser(userId string) (*models.User, error) {
	return repositorys.GetUser(userId)
}


func CreateUser(input dto.CreateUserData) (*models.User, error) {

    // 既に登録されているか確認
    user, err := repositorys.GetUserByName(input.Name)
    if err != nil {
        return nil, err
    }

    if user != nil {
        // ユーザーが既に存在する場合の処理
        return nil, errors.New("user already exists")
    }

    // 新しいユーザーを作成
    newUser := &models.User{
        Name:    input.Name,
        IconURL: input.IconURL,
    }

	err = repositorys.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
