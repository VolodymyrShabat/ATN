package business

import (
	"fmt"
	"github.com/VolodymyrShabat/Test_ATN/src/storage/models"
	"github.com/golang-jwt/jwt"
	"time"
)

type Token struct {
	TokenString string
	ExpiresAt   time.Time
}

type CustomClaims struct {
	UserId int
	jwt.StandardClaims
}

func (u *UserService) CreateJWT(userId int, isReset bool) (response *models.Token, err error) {
	accessClaims := &CustomClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
	}
	if isReset {
		accessClaims.StandardClaims.ExpiresAt = time.Now().Add(10 * time.Minute).Unix()
	} else {
		accessClaims.StandardClaims.ExpiresAt = time.Now().Add(30 * time.Minute).Unix()
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(u.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	refreshClaims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(8 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(u.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	return &models.Token{AccessToken: accessToken, RefreshToken: refreshToken}, err
}

func (u *UserService) VerifyToken(signedToken string) (CustomClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return 0, fmt.Errorf("unexpected signing method")
			}
			return u.JwtSecretKey, nil
		})
	if err != nil {
		return CustomClaims{}, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return CustomClaims{}, fmt.Errorf("parse claims: %w", err)
	}

	return CustomClaims{
		UserId: claims.UserId,
	}, nil
}

func (u *UserService) RefreshToken(refreshToken string) (*models.Token, error) {
	claims, err := u.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	user, err := u.userDAO.GetUserById(claims.UserId)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	tokenPair, err := u.CreateJWT(user.Id, false)
	if err != nil {
		return nil, fmt.Errorf("create token pair: %w", err)
	}

	return tokenPair, nil
}
