package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/zhashkevych/go-clean-architecture/auth"
)

type AuthClaims struct {
	jwt.StandardClaims
	User *auth.User `json:"user"`
}

type AuthUseCase struct {
	userRepo   auth.UserRepository
	hashSalt   string
	signingKey []byte
}

func NewAuthUseCase(userRepo auth.UserRepository, hashSalt string, signingKey []byte) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		hashSalt:   hashSalt,
		signingKey: signingKey,
	}
}

func (a *AuthUseCase) SignUp(ctx context.Context, username, password string) error {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))

	user := &auth.User{
		Username: username,
		Password: fmt.Sprintf("%x", pwd.Sum(nil)),
	}

	return a.userRepo.CreateUser(ctx, user)
}

func (a *AuthUseCase) SignIn(ctx context.Context, username, password string) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.userRepo.GetUser(ctx, username, password)
	if err != nil {
		return "", auth.ErrUserNotFound
	}

	claims := AuthClaims{
		User: user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(a.signingKey)
}

func (a *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (*auth.User, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(AuthClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, auth.ErrInvalidAccessToken
}
