package users

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/persistent/users_repository"
)

func NewUserService(repo *users_repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}