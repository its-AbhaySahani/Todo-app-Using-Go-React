package users

import (
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/domain"
)

func NewUserService(repo domain.UserRepository) *UserService {
    return &UserService{repo: repo}
}