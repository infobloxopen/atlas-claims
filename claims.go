package atlas_claims

import (
	"github.com/golang-jwt/jwt/v4"
)

// Subject describes the authenticated entity making a request.
type Subject struct {
	Id                 string `json:"id,omitempty"`                // user_email if authentication_type=(user_jwt|api_key)
	SubjectType        string `json:"subject_type"`                //valid values: user/s2s
	AuthenticationId   string `json:"authentication_id,omitempty"` // the id of the api_key if authentication_type=token
	AuthenticationType string `json:"authentication_type"`         //valid values: bearer/token
}

// Claims models the claims that atlas authz cares about.
type Claims struct {
	UserId            string   `json:"user_id,omitempty"`
	IdentityUserId    string   `json:"identity_user_id,omitempty"`
	AccountId         string   `json:"account_id,omitempty"`
	IdentityAccountId string   `json:"identity_account_id,omitempty"`
	Service           string   `json:"service,omitempty"`
	Groups            []string `json:"groups,omitempty"`
	Subject           Subject  `json:"subject"`
	jwt.StandardClaims
}

func (c Claims) Valid() error {
	return c.StandardClaims.Valid()
}
