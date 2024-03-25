package business

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/VolodymyrShabat/Test_ATN/src/storage"
	"github.com/VolodymyrShabat/Test_ATN/src/storage/models"
)

type UserService struct {
	userDAO      storage.UserDAO
	JwtSecretKey []byte
	Salt         []byte
}

func NewUserService(userDAO storage.UserDAO, JwtSecretKey, Salt []byte) *UserService {
	return &UserService{
		userDAO:      userDAO,
		JwtSecretKey: JwtSecretKey,
		Salt:         Salt,
	}
}

func (us UserService) SignUp(user models.User) error {
	user.Password = hashPassword(user.Password, us.Salt)
	err := us.userDAO.CreateUser(user)
	return err
}

func (us UserService) GetUserById(id int) (*models.User, error) {
	user, err := us.userDAO.GetUserById(id)
	return user, err
}

func (us UserService) GetUserByLogin(login string) (*models.User, error) {
	user, err := us.userDAO.GetUserByLogin(login)
	return user, err
}

func (us UserService) SignIn(login, password string) (bool, error) {
	user, err := us.userDAO.GetUserByLogin(login)
	if err != nil {
		return false, fmt.Errorf("error during getting user by login :%v", err)
	}
	if user == nil {
		return false, fmt.Errorf("user not found with such login")
	}
	return hashPassword(password, us.Salt) == user.Password, err
}

func (us UserService) FindUserByEmail(email string) (*models.User, error) {
	user, err := us.userDAO.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error during getting user by login :%v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found with such login")
	}
	return user, nil
}

func (us UserService) UpdatePasswordResetToken(email, passwordResetToken string) error {
	err := us.userDAO.UpdatePasswordResetToken(email, passwordResetToken)
	if err != nil {
		return fmt.Errorf("error during updating reset token %v", err)
	}
	return nil
}

func (us UserService) UpdatePasswordByResetToken(resetToken, password string) error {
	err := us.userDAO.UpdatePasswordByResetToken(resetToken, hashPassword(password, us.Salt))
	if err != nil {
		return fmt.Errorf("error during updating password by reset token %v", err)
	}
	return nil
}

func hashPassword(password string, salt []byte) string {
	var passwordBytes = []byte(password)

	var sha512Hasher = sha512.New()

	passwordBytes = append(passwordBytes, salt...)

	sha512Hasher.Write(passwordBytes)

	return hex.EncodeToString(sha512Hasher.Sum(nil))
}
