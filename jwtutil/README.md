# jwtutil

## Getting Started

Init
```golang
func init() {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		jwtutil.SetJwtSecret(s)
	}
}
```

Generate JWT token
```golang
token, err := jwtutil.NewToken(map[string]string{
        "iss":      "account",
        "aud":      "pangpanglabs",
        "username": "jack",
        "tenant":   "github",
})
```

Generate JWT token with secret
```golang
token, err := jwtutil.NewTokenWithSecret(map[string]string{
        "iss":      "account",
        "aud":      "pangpanglabs",
        "username": "jack",
        "tenant":   "github",
}, myJwtSecret)
```

Extract claim info
```golang
claim, err := jwtutil.Extract(token)
```

Extract claim info with secret
```golang
claim, err := jwtutil.Extract(token, secret)
```

Renew token
```golang
newToken, err := jwtutil.Renew(oldToken)
```