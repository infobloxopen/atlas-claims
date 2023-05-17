package atlas_claims

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
)

/*
	Setup instructions to create a spire JWT

-> % bin/spire-server token generate -spiffeID spiffe://example.org/myagent                                               env-4-env-4

Token: f92c3053-1943-4201-8df5-a15f8a200efe

-> % bin/spire-agent run -config conf/agent/agent.conf -joinToken f92c3053-1943-4201-8df5-a15f8a200efe

-> % bin/spire-agent api fetch jwt -spiffeID spiffe://example.org/myservice -audience ib-claims -timeout 78840000s

token(spiffe://example.org/myservice):

	eyJhbGciOiJFUzI1NiIsImtpZCI6Ikg5czlFRXg2RjZtamhxWUsxZ3Vjb3c2SXZSMXd1VlZaIiwidHlwIjoiSldUIn0.eyJhdWQiOlsiaWItY2xhaW1zIl0sImV4cCI6MTY4Mzc0MjQzMSwiaWF0IjoxNjgzNzQyMTMxLCJzdWIiOiJzcGlmZmU6Ly9leGFtcGxlLm9yZy9teXNlcnZpY2UifQ.znebAaE48jao0iqzq9O-UMS4GJxaGN3wScfGJ2unfEugzKsX7xuc5fMktpP6X6Y_nq8d9RK2mhd0ei9PoWJeLg

bundle(spiffe://example.org):

	        {
	    "keys": [
	        {
	            "kty": "EC",
	            "kid": "H9s9EEx6F6mjhqYK1gucow6IvR1wuVVZ",
	            "crv": "P-256",
	            "x": "tvkIptzDgghgQP-zU-6wUN57YZ-trnVQQsPqOzTaC-o",
	            "y": "QFA1EQN7Qdc9RAnuJqLyXCYrCEyfFzpDhSFLgJ6RkF0"
	        }
	    ]
	}
*/
func TestSpire_token(t *testing.T) {

	for _, tt := range []struct {
		raw string
		e   jwt.RegisteredClaims
	}{
		{
			raw: "eyJhbGciOiJFUzI1NiIsImtpZCI6Ikg5czlFRXg2RjZtamhxWUsxZ3Vjb3c2SXZSMXd1VlZaIiwidHlwIjoiSldUIn0.eyJhdWQiOlsiaWItY2xhaW1zIl0sImV4cCI6MTg0MTUyOTYwMCwiaWF0IjoxNjgzNzQyMTMxLCJzdWIiOiJzcGlmZmU6Ly9leGFtcGxlLm9yZy9teXNlcnZpY2UifQ.SphbHHorF3Z-IkJAavGE95BCX3vRkrigX5fO95RGTYymckFag9LOmMYLpu_78jiFCnaPh-CQxFf_FXTr_n8SRw",
			e: jwt.RegisteredClaims{
				Issuer:   "",
				Subject:  "spiffe://example.org/myservice",
				Audience: jwt.ClaimStrings{"ib-claims"},
				ID:       "spiffe://example.org/myservice",
			},
		},
	} {

		svid, err := ParseUnverifiedSpireJWT(tt.raw)
		if err != nil {
			t.Fatal(err)
		}

		future, _ := time.Parse(time.RFC3339, "2028-04-01T15:04:05Z")

		if !future.Before(svid.Expiry) {
			t.Error("Expires check failed")
		}
	}
}

func TestSpireToGoJWT(t *testing.T) {

	for _, tt := range []struct {
		raw string
		err error
	}{
		{
			raw: "eyJhbGciOiJFUzI1NiIsImtpZCI6Ikg5czlFRXg2RjZtamhxWUsxZ3Vjb3c2SXZSMXd1VlZaIiwidHlwIjoiSldUIn0.eyJhdWQiOlsiaWItY2xhaW1zIl0sImV4cCI6MTg0MTUyOTYwMCwiaWF0IjoxNjgzNzQyMTMxLCJzdWIiOiJzcGlmZmU6Ly9leGFtcGxlLm9yZy9teXNlcnZpY2UifQ.SphbHHorF3Z-IkJAavGE95BCX3vRkrigX5fO95RGTYymckFag9LOmMYLpu_78jiFCnaPh-CQxFf_FXTr_n8SRw",
		},
	} {

		svid, err := jwtsvid.ParseInsecure(tt.raw, []string{})
		if err != nil {
			t.Fatal(err)
		}

		_, claims, err := SVIDToGoJWT(&SVID{SVID: *svid})
		if err != tt.err {
			t.Error(err)
		}

		if !claims.VerifyAudience("ib-claims", true) {
			t.Error(claims.Issuer)
		}

		future, _ := time.Parse(time.RFC3339, "2028-04-01T15:04:05Z")

		if !claims.VerifyExpiresAt(future, true) {
			t.Error("token is expired")
		}

	}

}
