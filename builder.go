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
	if claims.StandardClaims.Issuer == "" {
		claims.StandardClaims.Issuer = DefaultIssuer
	}
	if claims.StandardClaims.Audience == "" {
		claims.StandardClaims.Audience = DefaultAudience
	}
	if claims.StandardClaims.IssuedAt == 0 {
		claims.StandardClaims.IssuedAt = time.Now().Unix()
	}
	if claims.StandardClaims.NotBefore == 0 {
		claims.StandardClaims.NotBefore = claims.StandardClaims.IssuedAt
	}
	if claims.StandardClaims.ExpiresAt == 0 {
		claims.StandardClaims.ExpiresAt = time.Unix(claims.StandardClaims.IssuedAt, 0).Add(expires_duration).Unix()
	}

	// non-standard claims
	if claims.AccountId == "" && claims.StandardClaims.Audience != STKAudience {
		claims.AccountId = "0"
	}
	if claims.Service == "" {
		claims.Service = DefaultService
	}
	if claims.Subject.Id == "" {
		claims.Subject.Id = fmt.Sprintf("service.%s.%d", claims.Service, claims.StandardClaims.IssuedAt)
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
