package service

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) SignUp() error {
	passwordHash, err :=
}
