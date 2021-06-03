// +build linux

package store_test

import (
	"os"
	"testing"

	"github.com/99designs/keyring"
	"github.com/joemiller/vault-token-helper/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestSecretServiceStore(t *testing.T) {
	// TODO: get this working in CI. The current blocker is needing to have a dbus prompter service that
	// can be driven automatically and headless:
	// https://github.com/99designs/keyring/blob/d9b6b92e219ff56ce753cf84d4956f823d431651/libsecret_test.go#L13-L22
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}

	st, err := store.New(keyring.Config{
		ServiceName:              "test",
		KeychainTrustApplication: true,
		LibSecretCollectionName:  "test",
		AllowedBackends:          []keyring.BackendType{keyring.SecretServiceBackend},
	})
	assert.Nil(t, err)
	assert.NotNil(t, st)

	// should be empty
	tokens, err := st.List()
	assert.Nil(t, err)
	assert.Empty(t, tokens)

	// Get of a missing item should not return an error
	_, err = st.Get("https://localhost:8200")
	assert.Nil(t, err)

	// Store a token
	token1 := store.Token{
		VaultAddr: "https://localhost:8200",
		Token:     "token-foo",
	}

	err = st.Store(token1)
	assert.Nil(t, err)

	// GetAll tokens
	tokens, err = st.List()
	assert.Nil(t, err)
	assert.NotEmpty(t, tokens)

	// Get a token by addr. Mixed case addr should be normalized for a successful lookup
	v1, err := st.Get("httpS://LOCALhost:8200/")
	assert.Nil(t, err)
	assert.Equal(t, token1, v1)

	// Erase
	err = st.Erase("https://localhost:8200")
	assert.Nil(t, err)

	// empty token store
	tokens, err = st.List()
	assert.Nil(t, err)
	assert.Empty(t, tokens)
}
