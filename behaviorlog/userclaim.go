package behaviorlog

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

type UserClaim struct {
	SessionID  string
	Aud        string
	TenantCode string
	Username   string
}

func NewUserClaimFromJwtToken(token string) (userClaim UserClaim) {
	ss := strings.Split(token, ".")
	if len(ss) != 3 {
		return
	}

	payload, err := decodeSegment(ss[1])
	if err != nil {
		return
	}

	json.Unmarshal(payload, &userClaim)
	userClaim.SessionID = ss[2]

	return
}

func decodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}
