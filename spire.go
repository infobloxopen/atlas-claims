package atlas_claims

import "github.com/spiffe/go-spiffe/v2/svid/jwtsvid"

// ParseUnverifiedSpireJWT parses a JWT token and returns the SVID and error
func ParseUnverifiedSpireJWT(token string) (*SVID, error) {

	svid, err := jwtsvid.ParseInsecure(token, []string{})
	if err != nil {
		return nil, err
	}

	return &SVID{
		SVID: *svid,
	}, err
}
