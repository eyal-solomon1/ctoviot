package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// HashedPassword generates bcrypt hashed password from a string
func TestPassword(t *testing.T) {
	pass := RandomString(10)
	hashedPassword, err := HashPassword(pass)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(pass, hashedPassword)
	require.NoError(t, err)

	wrongPass := RandomString(6)

	err = CheckPassword(wrongPass, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword(pass)

	require.NoError(t, err)
	require.NotEqual(t, hashedPassword, hashedPassword2) // as each hashed password has a unieq SALT added to it

}
