# atlas-claims

This repo models the JTW that is used for authentication in Atlas based applications.

1. import

```
import atlas_claims "github.com/infobloxopen/atlas-claims"
```


2. Getting the claims

```
claims, ok := atlas_claims.ClaimsFromContext(ctx)
```

3. Getting claims from command line
```bash
$ go run cmd/jwt-builder/main.go -h
Usage of /var/folders/49/3pxhjsps4fx4q21nkbj1j6n00000gp/T/go-build270617767/b001/exe/main:
  -account_id string
        account id (int)
  -audience string
        audience (default "ib-stk")
  -expires string
        expires duration, i.e. 48h (default "8760h0m0s")
  -groups string
        groups (comma separated list)
  -hmac_key string
        default hmac key (default "swordfish")
  -identity_account_id string
        identity account id (guid)
  -identity_user_id string
        identity user id
  -issuer string
        issuer (default "atlas-claims")
  -service string
        service (default "all")
  -subject_auth_id string
        subject authentication id
  -subject_auth_type string
        subject authentication type
  -subject_id string
        subject id
  -subject_type string
        subject type (default "s2s")
  -user_email string
        user email
  -user_id string
        user id
```

```bash
go run cmd/jwt-builder/main.go -account_id=123-456-789
```
