package atlas_claims

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestGetAccountID(t *testing.T) {
	var accountIDTests = []struct {
		claims   *Claims
		expected string
		err      error
	}{
		{
			claims: &Claims{
				AccountId: "id-abc-123",
			},
			expected: "id-abc-123",
			err:      nil,
		},
		{
			claims: &Claims{
				AccountId:     "id-abc-123",
				CompartmentID: "cmp-1",
			},
			expected: "id-abc-123",
			err:      nil,
		},
		{
			claims:   &Claims{},
			expected: "",
			err:      errMissingField,
		},
	}
	for _, test := range accountIDTests {
		token := makeToken(test.claims, t)
		ctx := contextWithToken(token, DefaultSubjectAuthType)
		actual, err := GetAccountID(ctx)
		if err != test.err {
			t.Errorf("Invalid error value: %v - expected %v", err, test.err)
		}
		if actual != test.expected {
			t.Errorf("Invalid AccountID: %v - expected %v", actual, test.expected)
		}
	}
}

func TestGetCompartmentID(t *testing.T) {
	var compartmentIDTests = []struct {
		claims   *Claims
		expected string
	}{
		{
			claims: &Claims{
				AccountId:     "id-abc-123",
				CompartmentID: "",
			},
			expected: "",
		},
		{
			claims: &Claims{
				AccountId:     "id-abc-123",
				CompartmentID: "cmp-1",
			},
			expected: "cmp-1",
		},
		{
			claims:   &Claims{},
			expected: "",
		},
	}
	for _, test := range compartmentIDTests {
		token := makeToken(test.claims, t)
		ctx := contextWithToken(token, DefaultSubjectAuthType)
		actual, _ := GetCompartmentID(ctx)
		if actual != test.expected {
			t.Errorf("Invalid CompartmentID: %v - expected %v", actual, test.expected)
		}
	}
}

// contextWithToken creates a context with a JWT
func contextWithToken(token, tokenType string) context.Context {
	md := metadata.Pairs(
		"authorization", fmt.Sprintf("%s %s", tokenType, token),
	)
	return metadata.NewIncomingContext(context.Background(), md)
}

// makeToken creates a JWT
func makeToken(claims *Claims, t *testing.T) string {
	cfgStandardClaimsExpires := time.Hour * 24 * 365
	testToken, err := BuildJwt(claims, "hmackey", cfgStandardClaimsExpires)
	if err != nil {
		t.Fatalf("Error when building token: %v", err)
	}

	return testToken
}
