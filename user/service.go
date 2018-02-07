package user

import "github.com/zamedic/multipoolminerbot/multipoolminer"

type Service interface {
	setToken(user, token string) error
}

type service struct {
	db Store
	userService multipoolminer.Service
}

func NewService(db Store,userService multipoolminer.Service) Service {
	return &service{db:db, userService:userService}
}

func (s *service)setToken(user, token string) error {
	_, err := s.userService.CheckAPIToken(token)
	if err != nil {
		return err
	}
	s.db.saveKey(user,token)
	return nil
}


