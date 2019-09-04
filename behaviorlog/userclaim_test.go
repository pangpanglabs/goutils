package behaviorlog

import (
	"testing"

	"github.com/pangpanglabs/goutils/test"
)

func TestUserClaim(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhdWQiLCJ0ZW5hbnRDb2RlIjoidGVuYW50Q29kZSIsInN1YiI6IjEyMzQ1Njc4OTAiLCJ1c2VybmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.fktZVBVAJHaZfySgZ1ameBHRhCKw5asDDkpJrbwOXKc"
	userClaim := NewUserClaimFromJwtToken(token)
	test.Equals(t, len(userClaim.SessionID) != 0, true)
	test.Equals(t, len(userClaim.Aud) != 0, true)
	test.Equals(t, len(userClaim.Username) != 0, true)
	test.Equals(t, len(userClaim.TenantCode) != 0, true)
}
