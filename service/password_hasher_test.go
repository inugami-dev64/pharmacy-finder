package service_test

import (
	"pharmafinder/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedPassword = "Password123"

func TestPasswordHashingAndVerification_ExpectSuccess(t *testing.T) {
	hasher := service.ProvidePasswordHasher()
	hash, err := hasher.CreatePasswordHash(expectedPassword)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)

	assert.True(t, hasher.Verify(expectedPassword, hash))
}

func TestPasswordHashingAndVerification_ExpectMismatch(t *testing.T) {
	hasher := service.ProvidePasswordHasher()
	hash, err := hasher.CreatePasswordHash(expectedPassword)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)

	assert.False(t, hasher.Verify("WarezWarehouse123", hash))
}
