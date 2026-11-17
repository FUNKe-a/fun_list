package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	jwt.RegisteredClaims
	Name string `json:"name"`
	Role string `json:"role,omitempty"`
}

type AuthService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
	accDur     time.Duration
	refDur     time.Duration
}

func NewAuthService(priv *rsa.PrivateKey, pub *rsa.PublicKey, issuer string, accDur time.Duration, refDur time.Duration) *AuthService {
	return &AuthService{
		privateKey: priv,
		publicKey:  pub,
		issuer:     issuer,
		accDur:     accDur,
		refDur:     refDur,
	}
}

func (s *AuthService) AccTTL() time.Duration { return s.accDur }
func (s *AuthService) RefTTL() time.Duration { return s.refDur }

func CreateRef() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	token := base64.RawURLEncoding.EncodeToString(buf)
	return token
}

// func (s *AuthService) RotateRef(token string) (string, error) {
// 	nTo := genRef()
//
// 	tokenHa := sha256.Sum256([]byte(token))
// 	nToHa := sha256.Sum256([]byte(nTo))
//
// 	_, err := s.Queries.RotateToken(nil, sqlc.RotateTokenParams{
// 		TokenHash:   nToHa[:],
// 		TokenHash_2: tokenHa[:],
// 	})
//
// 	return nTo, err
// }
//
// func (s *AuthService) IssueRef(userID int64) (string, error) {
// 	token := genRef()
// 	hash := sha256.Sum256([]byte(token))
//
// 	now := time.Now()
// 	_, err := s.Queries.CreateRefreshToken(nil, sqlc.CreateRefreshTokenParams{
// 		TokenHash: hash[:],
// 		UserID:    userID,
// 		IssuedAt:  now.Unix(),
// 		ExpiresAt: now.Add(s.refDur).Unix(),
// 	})
//
// 	return token, err
// }

func (s *AuthService) IssueAcc(userID, name, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accDur)),
		},
		Name: name,
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *AuthService) Validate(tokenString string) (*Claims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
		jwt.WithIssuer(s.issuer),
		jwt.WithExpirationRequired(),
	)

	claims := &Claims{}
	token, err := parser.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return s.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
