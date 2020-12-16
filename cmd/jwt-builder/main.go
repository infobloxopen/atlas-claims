package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	atlas_claims "github.com/infobloxopen/atlas-claims"
)

const (
	s2s_audience = "ib-stk"
)

var (
	// default claims
	cfgUserId                 = ""
	cfgIdentityUserId         = ""
	cfgAccountId              = ""
	cfgIdentityAccountId      = ""
	cfgService                = atlas_claims.DefaultService
	strGroups                 = ""
	cfgSubjectId              = ""
	cfgSubjectType            = atlas_claims.DefaultSubjectType
	cfgSubjectAuthId          = ""
	cfgSubjectAuthType        = ""
	cfgStandardClaimsAudience = atlas_claims.DefaultAudience
	cfgStandardClaimsExpires  = time.Duration(time.Hour * 24 * 365).String()
	cfgStandardClaimsIssuer   = atlas_claims.DefaultIssuer
	cfgHmacKey                = "swordfish"
)

func main() {
	// define flag overrides
	flag.StringVar(&cfgUserId, "user_id", cfgUserId, "user id")
	flag.StringVar(&cfgIdentityUserId, "identity_user_id", cfgIdentityUserId, "identity user id")
	flag.StringVar(&cfgAccountId, "account_id", cfgAccountId, "account id (int)")
	flag.StringVar(&cfgIdentityAccountId, "identity_account_id", cfgIdentityAccountId, "identity account id (guid)")
	flag.StringVar(&cfgService, "service", cfgService, "service")
	flag.StringVar(&strGroups, "groups", strGroups, "groups (comma separated list)")
	flag.StringVar(&cfgSubjectId, "subject_id", cfgSubjectId, "subject id")
	flag.StringVar(&cfgSubjectType, "subject_type", cfgSubjectType, "subject type")
	flag.StringVar(&cfgSubjectAuthId, "subject_auth_id", cfgSubjectAuthId, "subject authentication id")
	flag.StringVar(&cfgSubjectAuthType, "subject_auth_type", cfgSubjectAuthType, "subject authentication type")
	flag.StringVar(&cfgStandardClaimsAudience, "audience", cfgStandardClaimsAudience, "audience")
	flag.StringVar(&cfgStandardClaimsExpires, "expires", cfgStandardClaimsExpires, "expires duration, i.e. 48h")
	flag.StringVar(&cfgStandardClaimsIssuer, "issuer", cfgStandardClaimsIssuer, "issuer")
	flag.StringVar(&cfgHmacKey, "hmac_key", cfgHmacKey, "default hmac key")
	flag.Parse()

	// define claims from cli
	cfgGroups := strings.Split(strGroups, ",")
	claims := &atlas_claims.Claims{
		UserId:            cfgUserId,
		IdentityUserId:    cfgIdentityUserId,
		AccountId:         cfgAccountId,
		IdentityAccountId: cfgIdentityAccountId,
		Service:           cfgService,
		Groups:            cfgGroups,
		Subject: atlas_claims.Subject{
			Id:                 cfgSubjectId,
			SubjectType:        cfgSubjectType,
			AuthenticationId:   cfgSubjectAuthId,
			AuthenticationType: cfgSubjectType,
		},
		StandardClaims: jwt.StandardClaims{
			Audience: cfgStandardClaimsAudience,
			Issuer:   cfgStandardClaimsIssuer,
		},
	}

	// parse expires
	dur, err := time.ParseDuration(cfgStandardClaimsExpires)
	if err != nil {
		log.Fatalf("failed to parse expires duration: %s", err)
	}

	// build jwt
	tokenString, err := atlas_claims.BuildJwt(claims, cfgHmacKey, dur)
	if err != nil {
		log.Fatalf("failed to build jwt due to error: %s", err)
	}

	// print jwt
	fmt.Printf("%v\n", tokenString)
}
