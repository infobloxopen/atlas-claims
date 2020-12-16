package atlas_claims

import (
	"context"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
)

const (
	SetJwtHeader = "set-authorization"
	JwtName      = "bearer"
)

func UnverifiedClaimsFromContext(ctx context.Context) (*Claims, bool) {
	bearer, newBearer := AuthBearersFromCtx(ctx)
	validClaim, _ := UnverifiedClaimFromBearers([]string{bearer}, []string{newBearer})
	return validClaim, validClaim != nil
}

func AuthBearersFromCtx(ctx context.Context) (string, string) {
	var newBearer string
	bearer, _ := grpc_auth.AuthFromMD(ctx, JwtName)
	val := metautils.ExtractIncoming(ctx).Get(SetJwtHeader)
	if val != "" {
		splits := strings.SplitN(val, " ", 2)
		if len(splits) >= 2 && strings.ToLower(splits[0]) == strings.ToLower(JwtName) {
			newBearer = splits[1]
		}
	}
	return bearer, newBearer
}

func UnverifiedClaimFromBearers(bearer, newBearer []string) (*Claims, []error) {
	validBearerClaim, bearerErrorList := ParseUnverifiedClaimsFromJwtStrings(bearer)
	validNewBearerClaim, newBearerErrorList := ParseUnverifiedClaimsFromJwtStrings(newBearer)
	if len(newBearerErrorList) > 0 || len(bearerErrorList) > 0 {
		//fishy Should not have multiple newBearers
	}
	// Take the new bearer if possible.
	if validNewBearerClaim != nil {
		return validNewBearerClaim, nil
	} else if validBearerClaim != nil {
		return validBearerClaim, nil
	} else {
		return nil, append(bearerErrorList, newBearerErrorList...)
	}
}

func ParseUnverifiedClaimsFromJwtStrings(jwtStrings []string) (validClaim *Claims, errList []error) {
	validClaim, _, errList = ParseUnverifiedClaimsFromJwtStringsRaw(jwtStrings)
	return
}

// ParseUnverifiedClaimsFromJwtStringsRaw will return the raw (unmarshaled) jwt in addition to the valid claim.
func ParseUnverifiedClaimsFromJwtStringsRaw(jwtStrings []string) (validClaim *Claims, raw string, errList []error) {
	for _, jwtString := range jwtStrings {
		claims := &Claims{}
		parser := &jwt.Parser{}
		_, _, err := parser.ParseUnverified(jwtString, claims)

		// We use the most recent token
		if err != nil {
			errList = append(errList, err)
		} else {
			if validClaim == nil || (claims.IssuedAt > validClaim.IssuedAt) {
				validClaim = claims
				raw = jwtString
			}
		}
	}
	return
}
