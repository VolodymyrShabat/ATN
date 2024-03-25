package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/VolodymyrShabat/Test_ATN/src/business"
	"github.com/VolodymyrShabat/Test_ATN/src/storage/models"
	"github.com/VolodymyrShabat/Test_ATN/src/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type UserAccountService interface {
	SignUp(user models.User) error
	GetUserById(id int) (*models.User, error)
	GetUserByLogin(login string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	SignIn(login, password string) (bool, error)
	CreateJWT(userId int, isReset bool) (response *models.Token, err error)
	VerifyToken(signedToken string) (business.CustomClaims, error)
	RefreshToken(refreshToken string) (*models.Token, error)
	UpdatePasswordResetToken(email, passwordResetToken string) error
	UpdatePasswordByResetToken(resetToken, password string) error
}

type UserAccount struct {
	UserAccountService UserAccountService
	Config             *models.Config
	validator          *validator.Validate
}

func NewUserAccountHandler(UserAccountService *business.UserService, config *models.Config, validator *validator.Validate) *UserAccount {
	return &UserAccount{
		UserAccountService: UserAccountService,
		Config:             config,
		validator:          validator,
	}
}

type LoginPayload struct {
	Login    string `json:"login,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

func (ua *UserAccount) SignUp() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var person models.User
		err := json.NewDecoder(r.Body).Decode(&person)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during unmarshalling request: %v", err))
			return
		}

		if validationErr := ua.validator.Struct(&person); validationErr != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during validation: %v", validationErr))
			return
		}

		id := uuid.New()
		person.Id = int(id.ID())
		err = ua.UserAccountService.SignUp(person)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during creating user: %v", err))
			return
		}

		SendResponse(rw, http.StatusCreated, "successfully registered")
	}
}

func (ua *UserAccount) SignIn() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		var payload LoginPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during unmarshalling: %v", err))
			return
		}

		if validationErr := ua.validator.Struct(&payload); validationErr != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during validation: %v", validationErr))
			return
		}

		isAllowed, err := ua.UserAccountService.SignIn(payload.Login, payload.Password)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during handling sign in: %v", err))
			return
		}

		if !isAllowed {
			SendResponse(rw, http.StatusForbidden, "wrong password or login")
			return
		}

		u, err := ua.UserAccountService.GetUserByLogin(payload.Login)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during getting user by login: %v", err))
			return
		}

		res, err := ua.UserAccountService.CreateJWT(u.Id, false)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during creating tokens :%v", err))
			return
		}

		rw.Header().Del("Authorization")
		rw.Header().Del("refresh_token")
		rw.Header().Set("Authorization", "Bearer "+res.AccessToken)
		rw.Header().Set("refresh_token", "Bearer "+res.RefreshToken)
		SendResponse(rw, http.StatusOK, fmt.Sprintf("logged in successfully\n access token - %v", res.AccessToken))
	}
}

func (ua *UserAccount) GetUserById() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strId, ok := vars["id"]
		if !ok {
			SendResponse(rw, http.StatusBadRequest, "id in url is missing")
			return
		}

		id, err := strconv.Atoi(strId)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during converting url id: %v", err))
			return
		}

		u, err := ua.UserAccountService.GetUserById(id)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during getting user by id: %v", err))
			return
		}

		SendResponse(rw, http.StatusOK, fmt.Sprintf("%v", u))
	}
}

func (ua *UserAccount) ForgotPassword() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var userCredential models.ForgotPasswordInput
		err := json.NewDecoder(r.Body).Decode(&userCredential)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during unmarshalling: %v", err))
			return
		}

		if validationErr := ua.validator.Struct(userCredential); validationErr != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during validation: %v", validationErr))
			return
		}

		user, err := ua.UserAccountService.FindUserByEmail(userCredential.Email)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during getting user by email: %v", err))
			return
		}

		resetToken, err := ua.UserAccountService.CreateJWT(int(uuid.New().ID()), true)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during creating tokens :%v", err))
			return
		}

		passwordResetToken := resetToken.AccessToken
		err = ua.UserAccountService.UpdatePasswordResetToken(user.Email, passwordResetToken)
		if err != nil {
			SendResponse(rw, http.StatusInternalServerError, fmt.Sprintf("error during reseting password token :%v", err))
			return
		}

		emailData := utils.EmailData{
			URL:       "http://localhost:" + ua.Config.Port + "/user/reset-password/" + resetToken.AccessToken,
			FirstName: user.Name,
			Subject:   "Your password reset token (valid for 10min)",
		}

		errs := make(chan error, 1)
		go func() {
			errs <- utils.SendEmail(user, &emailData, "resetPassword.html", *ua.Config)
		}()
		if err := <-errs; err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during sending recovery email :%v", err))
			return
		}

		SendResponse(rw, http.StatusOK, "Recovery email has been sent to your email")
	}
}

func (ua *UserAccount) ResetPassword() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		token, ok := vars["reset_token"]
		if !ok {
			SendResponse(rw, http.StatusBadRequest, "id in url is missing")
			return
		}

		var userCredential models.ResetPasswordInput
		err := json.NewDecoder(r.Body).Decode(&userCredential)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during unmarshalling: %v", err))
			return
		}

		if validationErr := ua.validator.Struct(userCredential); validationErr != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during validation: %v", validationErr))
			return
		}

		_, err = ua.UserAccountService.VerifyToken(token)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during verifying token: %v", err))
			return
		}

		err = ua.UserAccountService.UpdatePasswordByResetToken(token, userCredential.Password)
		if err != nil {
			SendResponse(rw, http.StatusBadRequest, fmt.Sprintf("error during updating password: %v", err))
			return
		}

		SendResponse(rw, http.StatusOK, "Password successfully recovered")
	}
}
