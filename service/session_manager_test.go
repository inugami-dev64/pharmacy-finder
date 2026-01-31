package service_test

import (
	"os"
	"pharmafinder/db/entity"
	"pharmafinder/service"
	"pharmafinder/types"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func unwrap[T any](val T, err error) T {
	return val
}

var key = "a-string-secret-at-least-256-bit"

var testUser = entity.ModeratorUser{
	ID:            types.UUID(uuid.MustParse("b62a0163-5d3b-47c0-b5e7-41517bb5cfc2")),
	Username:      "karmen",
	Password:      "",
	FirstName:     "Karmen",
	LastName:      "Ott",
	Administrator: true,
}

func TestJWTTokenManager_BasicCreationAndVerification(t *testing.T) {
	os.Setenv("JWT_KEY", string(key))
	sessionManager := service.ProvideSessionManager()

	token := sessionManager.NewSessionToken(&testUser)
	assert.NotEmpty(t, token)

	data := sessionManager.VerifyToken(token)
	assert.NotNil(t, data)
	assert.Equal(t, testUser.ID, data.ID)
	assert.Equal(t, testUser.Username, data.Username)
	assert.Equal(t, testUser.Administrator, data.Admin)
}

func TestJWTTokenManager_ExpiredToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiaWF0IjoxNTY5ODYyNDIwLCJuYW1lIjoia2FybWVuIiwic3ViIjoiYjYyYTAxNjMtNWQzYi00N2MwLWI1ZTctNDE1MTdiYjVjZmMyIn0.J_jRG9eBhNAgDrbyTxKpYrKoYoecv_30nfNjBIW6l10"
	os.Setenv("JWT_KEY", string(key))
	sessionManager := service.ProvideSessionManager()

	data := sessionManager.VerifyToken(token)
	assert.Nil(t, data)
}

func TestJWTTokenManager_MalformedPayloadTypes(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6InRydWUiLCJpYXQiOiIxNTY5ODYyNDIwIiwibmFtZSI6Imthcm1lbiIsInN1YiI6ImI2MmEwMTYzLTVkM2ItNDdjMC1iNWU3LTQxNTE3YmI1Y2ZjMiJ9.6On67oyJ6rWDddFnSLGlv21fOcsQf3dQ3kJsHKfVZSk"
	os.Setenv("JWT_KEY", string(key))
	sessionManager := service.ProvideSessionManager()

	data := sessionManager.VerifyToken(token)
	assert.Nil(t, data)
}

func TestJWTTokenManager_MalformedUserID(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiaWF0IjoxNTY5ODYyNDIwLCJuYW1lIjoia2FybWVuIiwic3ViIjoiYjYyYTAxNjMtNWQzYi00N2MwLWI1ZTctNDE1MTdiYjVjZmMyYSJ9.-KHyLhkbWUIffY7ckk-Dku95InTy2HuW8DSz9dxSvL4"
	os.Setenv("JWT_KEY", string(key))
	sessionManager := service.ProvideSessionManager()

	data := sessionManager.VerifyToken(token)
	assert.Nil(t, data)
}
