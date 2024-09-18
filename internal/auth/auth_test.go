package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	seckey   = "keiw8Gah"
	username = "test"
	password = "testpass"
	userid   = int64(10)
)

func TestAuth(t *testing.T) {
	jwtm := NewJWTManager(seckey)
	user, err := jwtm.PrepareUser(username, password)
	require.NoError(t, err)
	hashpass := user.HashedPassword
	// change for next hop
	user.HashedPassword = password
	isHashAndPassCorrect := user.IsCorrectPassword(hashpass)
	require.Equal(t, true, isHashAndPassCorrect)

	str, err := jwtm.Generate(userid)
	require.NoError(t, err)
	// extract token
	uc, err := jwtm.Verify(str)
	require.NoError(t, err)
	require.Equal(t, userid, uc.UserID)
}
