package infra

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
)

type JWTIssuer struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTIssuer(secret, issuer string, accessTTL, refreshTTL time.Duration) *JWTIssuer {
	return &JWTIssuer{
		secret:     []byte(secret),
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (j *JWTIssuer) AccessTTL() time.Duration  { return j.accessTTL }
func (j *JWTIssuer) RefreshTTL() time.Duration { return j.refreshTTL }

func (j *JWTIssuer) IssueAccessToken(userID int64, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(j.accessTTL)
	claims := jwt.RegisteredClaims{
		Issuer:    j.issuer,
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}
	return signed, expiresAt, nil
}

func (j *JWTIssuer) ParseAccessToken(token string) (int64, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secret, nil
	}, jwt.WithIssuer(j.issuer), jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return 0, errors.Join(domain.ErrInvalidCredentials, err)
	}
	claims, ok := parsed.Claims.(*jwt.RegisteredClaims)
	if !ok || !parsed.Valid {
		return 0, domain.ErrInvalidCredentials
	}
	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, domain.ErrInvalidCredentials
	}
	return userID, nil
}

// GenerateRefreshToken returns the opaque token (to ship to the client exactly
// once) and its SHA-256 hex hash (to persist). The token is 32 bytes of CSPRNG
// output, base64url-encoded without padding.
func (j *JWTIssuer) GenerateRefreshToken() (string, string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("read random: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(buf)
	return token, j.HashRefreshToken(token), nil
}

func (j *JWTIssuer) HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
