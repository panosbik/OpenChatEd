package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"OpenChatEd/constants"
)

// JWT struct holds a private key and a public key for generating and decoding JWT tokens.
type JWT struct {
	privateKey []byte
	publicKey  []byte
}

// NewJWT creates a new JWT object using the provided private and public keys.
func NewJWT(privateKey []byte, publicKey []byte) JWT {
	return JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// EncodeJWToken generates a JWT token with a given time-to-live and payload (uint).
// Returns the token string, expiration timestamp, and any error that occurred.
func (j *JWT) EncodeJWToken(ttl time.Duration, payload uint) (string, int64, error) {
	now := time.Now().UTC()

	// Set the token claims with the provided payload and expiration times
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(payload), 10),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	// Parse the private key and create a signed token using the claims and RSA-256 signing method
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		log.Panicln("validate: parse key: %w", err)
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", 0, err
	}

	// Return the token string and expiration timestamp
	return token, claims.ExpiresAt.Unix(), nil
}

// DecodeJWToken decodes a given JWT token and returns its payload (uint) or an error.
func (j *JWT) DecodeJWToken(token string) (*uint, error) {
	// Parse the public key and decode the token
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		log.Panicln("validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (any, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf(constants.ErrInvalidToken.Error())
		}

		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidToken.Error())
	}

	// Get the claims from the token and return the payload (uint) or an error
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf(constants.ErrInvalidToken.Error())
	}

	// Convert the used ID string to uint
	u64, _ := strconv.ParseUint(claims["sub"].(string), 10, 32)
	id := uint(u64)

	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidToken.Error())
	}

	return &id, nil
}

func GenerateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 804)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(tokenBytes)
	return token, nil
}
