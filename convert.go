package atlas_claims

import (
	"github.com/golang-jwt/jwt/v4"
)

// SVIDToGoJWT converts secure jwt to RegisteredClaims
// it uses raw string since some data is lost in SVID
func SVIDToGoJWT(c *SVID) (*jwt.Token, jwt.RegisteredClaims, error) {

	s := c.Marshal()
	parser := jwt.NewParser()

	var claims jwt.RegisteredClaims
	tok, _, err := parser.ParseUnverified(s, &claims)

	return tok, claims, err
}
