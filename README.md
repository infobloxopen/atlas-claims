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

