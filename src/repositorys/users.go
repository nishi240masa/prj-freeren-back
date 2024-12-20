package repositorys

import (
	"md2s/models"
)

func GetUsers() ([]models.User, error) {
	var users []models.User

	query := db.Table("users")
	result := query.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func GetUser(userId string) (*models.User, error) {
	var user models.User

	query := db.Table("users").Where("id = ?", userId)
	result := query.First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func GetUserByName(name string) (*models.User, error) {
	var user models.User

	query := db.Table("users").Where("name = ?", name)
	result := query.First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func CreateUser(newUser *models.User) error {
	result := db.Create(newUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateUser(user *models.User) error {
	result := db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}