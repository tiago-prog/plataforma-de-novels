package usecase

import (
	"fmt"
	"time"

	"github.com/tiago-prog/novels-api/internal/model"
	"github.com/tiago-prog/novels-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) Register(user model.User) (int, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	user.Password = string(hash)
	user.Role = model.RoleReader
	user.Verified = false

	fmt.Println("User to be registered:", user)

	id, err := u.userRepo.RegisterUser(user)
	if err != nil {
		fmt.Println("Error registering user:", err)
		return 0, err
	}

	return id, nil
}

func (u *UserUsecase) GetUserByEmail(email string) (model.User, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		fmt.Println("Error fetching user by email:", err)
		return model.User{}, err
	}

	// Set LastLoginAt to current time if not set
	if user.LastLoginAt == nil || user.LastLoginAt.IsZero() {
		now := time.Now()
		user.LastLoginAt = &now
	}

	return user, nil
}

func (u *UserUsecase) Login(email, password string) (model.User, error) {
	user, err := u.userRepo.Login(email, password)
	if err != nil {
		fmt.Println("Error logging in:", err)
		return model.User{}, err
	}
	return user, nil
}

func (u *UserUsecase) SuspendUser(executorID int, targetID int) error {
	return u.userRepo.SuspendUser(executorID, targetID)
}

func (u *UserUsecase) GetAllUsers() ([]model.User, error) {
	users, err := u.userRepo.GetAllUsers()
	if err != nil {
		fmt.Println("Error fetching all users:", err)
		return nil, err
	}
	return users, nil
}
