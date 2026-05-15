package user

import (
	"daymark/internal/models"
	"daymark/pkg/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) SignUp(name, email, password string) (*models.User, string, error) {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	otp := utils.GenerateOTP()
	expiresAt := time.Now().Add(10 * time.Minute)

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user = &models.User{
				Name:              name,
				Email:             email,
				PasswordHash:      string(hash),
				Provider:          "email",
				EmailVerified:     false,
				EmailOTP:          otp,
				EmailOTPExpiresAt: expiresAt,
			}
			if err := s.repo.Create(user); err != nil {
				return nil, "", err
			}
		} else {
			return nil, "", err
		}
	} else {
		// Existing user with this email
		if user.Provider != "email" {
			return nil, "", fmt.Errorf("email already used with %s login", user.Provider)
		}
		user.Name = name
		user.PasswordHash = string(hash)
		user.EmailVerified = false
		user.EmailOTP = otp
		user.EmailOTPExpiresAt = expiresAt
		if err := s.repo.Update(user); err != nil {
			return nil, "", err
		}
	}

	return user, otp, nil
}

func (s *Service) SignIn(email, password string) (*models.User, error) {

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user.Provider == "email" && !user.EmailVerified {
		return nil, fmt.Errorf("email not verified")
	}

	if ok := utils.ComparePassword(user.PasswordHash, password); !ok {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *Service) OAuthLogin(provider, providerID, name, email, avatar string) (*models.User, error) {

	user, err := s.repo.GetByProvider(provider, providerID)

	if err == nil {
		// Update avatar if it's missing or different
		if avatar != "" && user.AvatarURL != avatar {
			user.AvatarURL = avatar
			s.repo.Update(user) // assuming we need an update method
		}
		return user, nil
	}

	newUser := &models.User{
		Name:          name,
		Email:         email,
		Provider:      provider,
		ProviderID:    providerID,
		AvatarURL:     avatar,
		EmailVerified: true,
	}

	err = s.repo.Create(newUser)

	return newUser, err
}

func (s *Service) GetByEmail(email string) (*models.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *Service) GetByID(id uint) (*models.User, error) {
	return s.repo.GetByID(id)
}

// VerifyEmail checks the OTP for a given email and marks the user as verified.
func (s *Service) VerifyEmail(email, otp string) (*models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user.Provider != "email" {
		return nil, fmt.Errorf("email is associated with %s login", user.Provider)
	}
	if user.EmailVerified {
		return user, nil
	}
	if time.Now().After(user.EmailOTPExpiresAt) {
		return nil, fmt.Errorf("otp expired")
	}
	if otp == "" || otp != user.EmailOTP {
		return nil, fmt.Errorf("invalid otp")
	}

	user.EmailVerified = true
	user.EmailOTP = ""
	user.EmailOTPExpiresAt = time.Time{}
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
