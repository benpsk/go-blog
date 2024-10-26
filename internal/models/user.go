package models

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type UserResponse struct {
	Id    int
	Name  string
	Email string
}

type UserInput struct {
	Id              int
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
}

type Session struct {
	Id     int
	Token  string
	UserId int
}

type AuthUser struct {
	Id    int
	Name  string
	Email string
}

type LoginError struct {
	Email    string
	Password string
}

type RegisterError struct {
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
}
