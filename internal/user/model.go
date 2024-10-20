package user

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
