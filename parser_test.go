package atlas_claims

import (
	"context"
	"log"
	"testing"

	"google.golang.org/grpc/metadata"
)

var ClaimsTestData = []struct {
	name      string
	headers   map[string]string
	signature string
	valid     bool
	jti       string
}{
	{
		name:      "valid HTTP header, valid signature, valid claims",
		signature: "secret",
		valid:     true,
		jti:       "1",
		headers:   map[string]string{"authorization": "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFwaV90b2tlbiI6IjIxOThhYjBmZWMwMzI3MDFiYzMyMWUiLCJhY2NvdW50X2lkIjoiMSIsImdyb3VwcyI6WyJhY3RfYWRtaW4iXSwic2VydmljZSI6ImFsbCIsImF1ZCI6ImliLXN0ayIsImV4cCI6MjU3MDAwMDAwMCwianRpIjoiMSIsImlhdCI6MTUzNTMyMTQwNywiaXNzIjoiYXRoZW5hLWF1dGhuLXN2YyIsIm5iZiI6MTUzNTMyMTQwN30.GRtr1SrWkc85jXAdyW2ncTIKuHdwgOBIbQGKeI67ySzcyCoefsheX7i3KnCaGxNf1VeyPUMNDohOgT7sJP7r_Q"},
	},
	{
		name:      "invalid HTTP header, valid signature, valid claims",
		signature: "secret",
		valid:     false,
		jti:       "1",
		headers:   map[string]string{"authorization1": "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFwaV90b2tlbiI6IjIxOThhYjBmZWMwMzI3MDFiYzMyMWUiLCJhY2NvdW50X2lkIjoiMSIsImdyb3VwcyI6WyJhY3RfYWRtaW4iXSwic2VydmljZSI6ImFsbCIsImF1ZCI6ImliLXN0ayIsImV4cCI6MjU3MDAwMDAwMCwianRpIjoiMSIsImlhdCI6MTUzNTMyMTQwNywiaXNzIjoiYXRoZW5hLWF1dGhuLXN2YyIsIm5iZiI6MTUzNTMyMTQwN30.GRtr1SrWkc85jXAdyW2ncTIKuHdwgOBIbQGKeI67ySzcyCoefsheX7i3KnCaGxNf1VeyPUMNDohOgT7sJP7r_Q"},
	},
	{
		name:      "valid HTTP header, invalid signature, valid claims",
		signature: "secret",
		valid:     true,
		jti:       "1",
		headers:   map[string]string{"authorization": "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFjY291bnRfaWQiOiIxIiwiZ3JvdXBzIjpbImFjdF9hZG1pbiJdLCJzZXJ2aWNlIjoiYWxsIiwiYXVkIjoiaWItc3RrIiwiZXhwIjoyNTcwMDAwMDAwLCJqdGkiOiIxIiwiaWF0IjoxNTM1MzIxNDA3LCJpc3MiOiJhdGhlbmEtYXV0aG4tc3ZjIiwibmJmIjoxNTM1MzIxNDA3fQ.CrXv-6U2JXZsIJdLlbqPHVoyWJoCBIJvhRpjUcYoD24fipHa-YPL5vKlZpItq4Q9r3P3pHPBGJm_VEh_Rkjvog"},
	},
	{
		name:      "valid HTTP header, valid signature, invalid claims",
		signature: "secret",
		valid:     false,
		jti:       "1",
		headers:   map[string]string{"authorization": "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFjY291bnRfaWQiOjEsImdyb3VwcyI6WyJhY3RfYWRtaW4iXSwic2VydmljZSI6ImFsbCIsImF1ZCI6ImliLXN0ayIsImV4cCI6MjU3MDAwMDAwMCwianRpIjoiVE9ETyIsImlhdCI6MTUzNTMyMTQwNywiaXNzIjoiYXV0aG4tc3ZjIiwibmJmIjoxNTM1MzIxNDA3fQ.UJNMjAeAbuKOPGNj7k-RECb-5omgBfGDosTsNNDBAFcjiqYI2JTktrsDonNTT6Q16zw1tKgU6zavwLFOzmphaA"},
	},
	{
		name:      "Auth: valid HTTP header, valid signature, valid claims; Set-Auth: valid HTTP header, valid signature, valid claims",
		signature: "secret",
		valid:     true,
		jti:       "2",
		headers: map[string]string{
			"authorization":     "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFjY291bnRfaWQiOiIxIiwiZ3JvdXBzIjpbImFjdF9hZG1pbiJdLCJzZXJ2aWNlIjoiYWxsIiwiYXVkIjoiaWItc3RrIiwiZXhwIjoyNTcwMDAwMDAwLCJqdGkiOiIxIiwiaWF0IjoxNTM1MzIxNDA3LCJpc3MiOiJhdXRobi1zdmMiLCJuYmYiOjE1MzUzMjE0MDd9.jYG8Q1d3tlYfaacFoqmvfWJ_HQw_E1t9SESUx0VESM2smtLkD7sBK0x9NDJidUS-TPJNZtoeMWw_c_JKnm9_ag",
			"set-authorization": "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFjY291bnRfaWQiOiIxIiwiZ3JvdXBzIjpbImFjdF9hZG1pbiJdLCJzZXJ2aWNlIjoiYWxsIiwiYXVkIjoiaWItc3RrIiwiZXhwIjoyNTcwMDAwMDAwLCJqdGkiOiIyIiwiaWF0IjoxNTM1MzIzNDA3LCJpc3MiOiJhdXRobi1zdmMiLCJuYmYiOjE1MzUzMjM0MDd9.x5HGvONge1Y4TXI8nIg5LcQC_qwCKNkXvpOXE5pv4jGJhEmUQ_SvYKcoc8s8v9zv5wMA-2vlDV2LiqtnhQOkiQ",
		},
	},
	{
		name:      "Auth: valid HTTP header, valid signature, valid claims; Set-Auth: invalid HTTP header, valid signature, valid claims",
		signature: "secret",
		valid:     true,
		jti:       "1",
		headers: map[string]string{
			"authorization":     "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFjY291bnRfaWQiOiIxIiwiZ3JvdXBzIjpbImFjdF9hZG1pbiJdLCJzZXJ2aWNlIjoiYWxsIiwiYXVkIjoiaWItc3RrIiwiZXhwIjoyNTcwMDAwMDAwLCJqdGkiOiIxIiwiaWF0IjoxNTM1MzIxNDA3LCJpc3MiOiJhdXRobi1zdmMiLCJuYmYiOjE1MzUzMjE0MDd9.jYG8Q1d3tlYfaacFoqmvfWJ_HQw_E1t9SESUx0VESM2smtLkD7sBK0x9NDJidUS-TPJNZtoeMWw_c_JKnm9_ag",
			"set-authorization": "beaer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFjY291bnRfaWQiOiIxIiwiZ3JvdXBzIjpbImFjdF9hZG1pbiJdLCJzZXJ2aWNlIjoiYWxsIiwiYXVkIjoiaWItc3RrIiwiZXhwIjoyNTcwMDAwMDAwLCJqdGkiOiIyIiwiaWF0IjoxNTM1MzIzNDA3LCJpc3MiOiJhdXRobi1zdmMiLCJuYmYiOjE1MzUzMjM0MDd9.x5HGvONge1Y4TXI8nIg5LcQC_qwCKNkXvpOXE5pv4jGJhEmUQ_SvYKcoc8s8v9zv5wMA-2vlDV2LiqtnhQOkiQ",
		},
	},
	{
		name:      "Auth: valid HTTP header, valid signature, valid claims; Set-Auth: valid HTTP header, invalid signature, valid claims",
		signature: "secret",
		valid:     true,
		jti:       "2",
		headers: map[string]string{
			"authorization":     "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFwaV90b2tlbiI6IjIxOThhYjBmZWMwMzI3MDFiYzMyMWUiLCJhY2NvdW50X2lkIjoiMSIsImdyb3VwcyI6WyJhY3RfYWRtaW4iXSwic2VydmljZSI6ImFsbCIsImF1ZCI6ImliLXN0ayIsImV4cCI6MjU3MDAwMDAwMCwianRpIjoiMSIsImlhdCI6MTUzNTMyMTQwNywiaXNzIjoiYXRoZW5hLWF1dGhuLXN2YyIsIm5iZiI6MTUzNTMyMTQwN30.GRtr1SrWkc85jXAdyW2ncTIKuHdwgOBIbQGKeI67ySzcyCoefsheX7i3KnCaGxNf1VeyPUMNDohOgT7sJP7r_Q",
			"set-authorization": "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFwaV90b2tlbiI6IjIxOThhYjBmZWMwMzI3MDFiYzMyMWUiLCJhY2NvdW50X2lkIjoiMSIsImdyb3VwcyI6WyJhY3RfYWRtaW4iXSwic2VydmljZSI6ImFsbCIsImF1ZCI6ImliLXN0ayIsImV4cCI6MjU3MDAwMDAwMCwianRpIjoiMiIsImlhdCI6MTUzNTMyMzQwNywiaXNzIjoiYXRoZW5hLWF1dGhuLXN2YyIsIm5iZiI6MTUzNTMyMzQwN30.GZ0wYTBO6dFA5ufZlWza33luMTYmaChqTQxyI9Zv9akN0MujBAOyIiCkvSp5Wa2NsPUMCfRpaDK4qk1SukhCMg",
		},
	},
	{
		name: "Auth: valid HTTP header, valid signature, valid claims; Set-Auth: valid HTTP header, valid signature, invalid claims",
		headers: map[string]string{
			"authorization":     "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFjY291bnRfaWQiOiIxIiwiZ3JvdXBzIjpbImFjdF9hZG1pbiJdLCJzZXJ2aWNlIjoiYWxsIiwiYXVkIjoiaWItc3RrIiwiZXhwIjoyNTcwMDAwMDAwLCJqdGkiOiIxIiwiaWF0IjoxNTM1MzIxNDA3LCJpc3MiOiJhdXRobi1zdmMiLCJuYmYiOjE1MzUzMjE0MDd9.jYG8Q1d3tlYfaacFoqmvfWJ_HQw_E1t9SESUx0VESM2smtLkD7sBK0x9NDJidUS-TPJNZtoeMWw_c_JKnm9_ag",
			"set-authorization": "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RlckB0ZXN0LmNvbSIsImFjY291bnRfaWQiOiIxIiwiZ3JvdXBzIjoiYWN0X2FkbWluIiwic2VydmljZSI6ImFsbCIsImF1ZCI6ImliLXN0ayIsImV4cCI6MjU3MDAwMDAwMCwianRpIjoiMiIsImlhdCI6MTUzNTMyMzQwNywiaXNzIjoiYXV0aG4tc3ZjIiwibmJmIjoxNTM1MzIzNDA3fQ.kjBQ7-iLbIZ7nzFeMtCuRSamFcABVwz2sVIEH9ZPJy8XdHMyKkC57oxGGcBFHWmYJfdZ4SBy-xGDz6QKgDYhDA",
		},
		signature: "secret",
		valid:     true,
		jti:       "1",
	},
	{
		name: "svc-svc token",
		headers: map[string]string{
			"authorization": "bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzZXJ2aWNlIjoiYWxsIiwiYXVkIjoiaWItc3RrIiwiZXhwIjoyMzk4ODcyNzc4LCJqdGkiOiJUT0RPIiwiaWF0IjoxNTM1MzIxNDA3LCJpc3MiOiJhdXRobi1zdmMiLCJuYmYiOjE1MzUzMjE0MDd9.KkG5ho_GiSvxA-p5cOjEaBLe9pvs3HPmBmYCJ87iiV9-uMA_GMdYrC3UB-jtzcsk7TRO1P0TwOjp7FO5a-lYJA",
		},
		signature: "secret",
		valid:     true,
		jti:       "1",
	},
}

func TestUnverifiedAthenaClaimsFromContext(t *testing.T) {
	ctx := context.Background()
	for _, test := range ClaimsTestData {
		log.Println(test.name)
		ctx2 := metadata.NewIncomingContext(ctx, metadata.New(test.headers))
		claims, hasClaims := UnverifiedClaimsFromContext(ctx2)
		log.Println("Valid Expected: ", test.valid)
		log.Println("Valid Actual: ", hasClaims)
		if test.valid != (claims != nil) && claims.Id == test.jti {
			log.Println("Failed Test")
			log.Println(test)
			t.Fail()
		}
	}
}
