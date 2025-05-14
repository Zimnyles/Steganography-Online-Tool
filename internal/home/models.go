package home


type UserCreateForm struct {
	Login    string
	Email    string
	Password string
}

type LoginForm struct {
	Login    string
	Email    string
	Password string
}

type UserCredentials struct {
	Login        string
	PasswordHash string
}



