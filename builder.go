package atlas_claims

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	STKAudience = "ib-stk"

	DefaultIssuer          = "atlas-claims"
	DefaultAudience        = STKAudience
	DefaultService         = "all"
	DefaultSubjectType     = "s2s"
	DefaultSubjectAuthType = "bearer"
)

func BuildJwt(claims *Claims, hmacKey string, expires_duration time.Duration) (string, error) {
	if len(strings.TrimSpace(hmacKey)) < 1 {
		return "", fmt.Errorf("non-empty hmac key is required")
	}

	// standard claims
	if claims.Issuer == "" {
		claims.Issuer = DefaultIssuer
	}
	if claims.Audience == "" {
		claims.Audience = DefaultAudience
	}
	if claims.IssuedAt == 0 {
		claims.IssuedAt = time.Now().Unix()
	}
	if claims.NotBefore == 0 {
		claims.NotBefore = claims.IssuedAt
	}
	if claims.ExpiresAt == 0 {
		claims.ExpiresAt = time.Unix(claims.IssuedAt, 0).Add(expires_duration).Unix()
	}

	// non-standard claims
	if claims.AccountId == "" && claims.Audience != STKAudience {
		claims.AccountId = "0"
	}
	if claims.Service == "" {
		claims.Service = DefaultService
	}
	if claims.Subject.Id == "" {
		claims.Subject.Id = fmt.Sprintf("service.%s.%d", claims.Service, claims.IssuedAt)
	}
	if claims.Subject.SubjectType == "" {
		claims.Subject.SubjectType = DefaultSubjectType
	}
	if claims.Subject.AuthenticationType == "" {
		claims.Subject.AuthenticationType = DefaultSubjectAuthType
	}

	// sign the jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(hmacKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign claim: %v", err)
	}

	return tokenString, nil
}
