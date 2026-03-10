package user

import (
	"daymark/internal/models"
	"daymark/pkg/utils"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) SignUp(name, email, password string) (*models.User, error) {

	hash, err:=utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Provider:     "email",
	}

	err = s.repo.Create(user)
	return user, err
}

func (s *Service) SignIn(email, password string) (*models.User, error) {

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	flag:=utils.ComparePassword(user.PasswordHash,password)

	if !flag {
		return nil, err
	}

	return user, nil
}

func (s *Service) OAuthLogin(provider, providerID, name, email, avatar string) (*models.User, error) {

	user, err := s.repo.GetByProvider(provider, providerID)

	if err == nil {
		return user, nil
	}

	newUser := &models.User{
		Name:       name,
		Email:      email,
		Provider:   provider,
		ProviderID: providerID,
		AvatarURL:  avatar,
	}

	err = s.repo.Create(newUser)

	return newUser, err
}