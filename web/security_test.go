package web_test

import (
	"fmt"
	"net/http"
	"pharmafinder/db/entity"
	"pharmafinder/mock"
	"pharmafinder/service"
	"pharmafinder/types"
	"pharmafinder/web"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiaWF0IjoxNTY5ODYyNDIwLCJuYW1lIjoia2FybWVuIiwic3ViIjoiYjYyYTAxNjMtNWQzYi00N2MwLWI1ZTctNDE1MTdiYjVjZmMyIn0"

var userID = types.UUID(uuid.MustParse("b62a0163-5d3b-47c0-b5e7-41517bb5cfc2"))

const username = "karmen"

func ptr[T any](val T) *T {
	return &val
}

/*******************************
 *** RuleAuthenticated tests ***
 *******************************/
func TestSecurityChain_RuleAuthenticatedAllow(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryMock := mock.NewMockQuery[entity.ModeratorUser](ctrl)
	queryMock.EXPECT().
		Query().
		Return(&entity.ModeratorUser{}, nil)

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	repoMock.EXPECT().
		FindUserByID(gomock.Eq(userID)).
		Return(queryMock)

	sessionManagerMock := mock.NewMockSessionManager(ctrl)
	sessionManagerMock.EXPECT().
		VerifyToken(gomock.Eq(token)).
		Return(&service.SessionUserData{
			Admin:    true,
			IssuedAt: types.Time(time.Unix(1569862420, 0)),
			Username: username,
			ID:       userID,
		})

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer "+token)
	assert.True(t, chain.ShouldPermitAccess(details))
}

func TestSecurityChain_RuleAuthenticated_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := mock.NewMockModeratorUserRepository(ctrl)

	sessionManagerMock := mock.NewMockSessionManager(ctrl)
	sessionManagerMock.EXPECT().
		VerifyToken(gomock.Eq(token)).
		Return(nil)

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer "+token)
	assert.False(t, chain.ShouldPermitAccess(details))
}

func TestSecurityChain_RuleAuthenticated_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryMock := mock.NewMockQuery[entity.ModeratorUser](ctrl)
	queryMock.EXPECT().
		Query().
		Return(nil, fmt.Errorf("db error"))

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	repoMock.EXPECT().
		FindUserByID(gomock.Eq(userID)).
		Return(queryMock)

	sessionManagerMock := mock.NewMockSessionManager(ctrl)
	sessionManagerMock.EXPECT().
		VerifyToken(gomock.Eq(token)).
		Return(&service.SessionUserData{
			Admin:    true,
			IssuedAt: types.Time(time.Unix(1569862420, 0)),
			Username: username,
			ID:       userID,
		})

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer "+token)
	assert.False(t, chain.ShouldPermitAccess(details))
}

func TestSecurityChain_RuleAuthenticated_MalformedAuthorizationHeader(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	sessionManagerMock := mock.NewMockSessionManager(ctrl)

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer     "+token)
	assert.False(t, chain.ShouldPermitAccess(details))
	details.Header.Del("Authorization")
	details.Header.Set("Authorization", token)
	assert.False(t, chain.ShouldPermitAccess(details))
}

/***********************
 *** RuleAdmin tests ***
 ***********************/

func TestSecurityChain_RuleAdminAllow(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryMock := mock.NewMockQuery[entity.ModeratorUser](ctrl)
	queryMock.EXPECT().
		Query().
		Return(&entity.ModeratorUser{Administrator: true}, nil)

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	repoMock.EXPECT().
		FindUserByID(gomock.Eq(userID)).
		Return(queryMock)

	sessionManagerMock := mock.NewMockSessionManager(ctrl)
	sessionManagerMock.EXPECT().
		VerifyToken(gomock.Eq(token)).
		Return(&service.SessionUserData{
			Admin:    true,
			IssuedAt: types.Time(time.Unix(1569862420, 0)),
			Username: username,
			ID:       userID,
		})

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAdmin(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer "+token)
	assert.True(t, chain.ShouldPermitAccess(details))
}

func TestSecurityChain_RuleAdmin_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockModeratorUserRepository(ctrl)

	sessionManagerMock := mock.NewMockSessionManager(ctrl)
	sessionManagerMock.EXPECT().
		VerifyToken(gomock.Eq(token)).
		Return(nil)

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAdmin(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer "+token)
	assert.False(t, chain.ShouldPermitAccess(details))
}

func TestSecurityChain_RuleAdmin_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryMock := mock.NewMockQuery[entity.ModeratorUser](ctrl)
	queryMock.EXPECT().
		Query().
		Return(nil, fmt.Errorf("db error"))

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	repoMock.EXPECT().
		FindUserByID(gomock.Eq(userID)).
		Return(queryMock)
	sessionManagerMock := mock.NewMockSessionManager(ctrl)
	sessionManagerMock.EXPECT().
		VerifyToken(gomock.Eq(token)).
		Return(&service.SessionUserData{
			Admin:    true,
			IssuedAt: types.Time(time.Unix(1569862420, 0)),
			Username: username,
			ID:       userID,
		})

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAdmin(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer "+token)
	assert.False(t, chain.ShouldPermitAccess(details))
}

func TestSecurityChain_RuleAdmin_JWTUserNotAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := mock.NewMockModeratorUserRepository(ctrl)

	sessionManagerMock := mock.NewMockSessionManager(ctrl)
	sessionManagerMock.EXPECT().
		VerifyToken(gomock.Eq(token)).
		Return(&service.SessionUserData{
			Admin:    false,
			IssuedAt: types.Time(time.Unix(1569862420, 0)),
			Username: username,
			ID:       userID,
		})

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAdmin(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer "+token)
	assert.False(t, chain.ShouldPermitAccess(details))
}

func TestSecurityChain_RuleAdmin_UserNotAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryMock := mock.NewMockQuery[entity.ModeratorUser](ctrl)
	queryMock.EXPECT().
		Query().
		Return(&entity.ModeratorUser{Administrator: false}, nil)

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	repoMock.EXPECT().
		FindUserByID(gomock.Eq(userID)).
		Return(queryMock)

	sessionManagerMock := mock.NewMockSessionManager(ctrl)
	sessionManagerMock.EXPECT().
		VerifyToken(gomock.Eq(token)).
		Return(&service.SessionUserData{
			Admin:    true,
			IssuedAt: types.Time(time.Unix(1569862420, 0)),
			Username: username,
			ID:       userID,
		})

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAdmin(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer "+token)
	assert.False(t, chain.ShouldPermitAccess(details))
}

func TestSecurityChain_RuleAdmin_MalformedAuthorizationHeader(t *testing.T) {
	ctrl := gomock.NewController(t)

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	sessionManagerMock := mock.NewMockSessionManager(ctrl)

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAdmin(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer     "+token)
	assert.False(t, chain.ShouldPermitAccess(details))
	details.Header.Del("Authorization")
	details.Header.Set("Authorization", token)
	assert.False(t, chain.ShouldPermitAccess(details))
}

/*********************************
 *** RuleWhenNoAdminUser tests ***
 *********************************/
func TestSecurityChain_RuleWhenNoAdminUserAllow(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryMock := mock.NewMockQuery[bool](ctrl)
	queryMock.EXPECT().
		Query().
		Return(ptr(false), nil)

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	repoMock.EXPECT().
		HasAdministrator().
		Return(queryMock)

	chain := web.NewSecurityChain[web.EmptyBody]().RuleWhenNoAdminUser(repoMock)
	assert.True(t, chain.ShouldPermitAccess(nil))
}

func TestSecurityChain_RuleWhenNoAdminUser_AdminExists(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryMock := mock.NewMockQuery[bool](ctrl)
	queryMock.EXPECT().
		Query().
		Return(ptr(true), nil)

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	repoMock.EXPECT().
		HasAdministrator().
		Return(queryMock)

	chain := web.NewSecurityChain[web.EmptyBody]().RuleWhenNoAdminUser(repoMock)
	assert.False(t, chain.ShouldPermitAccess(nil))
}

func TestSecurityChain_RuleWhenNoAdminUser_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryMock := mock.NewMockQuery[bool](ctrl)
	queryMock.EXPECT().
		Query().
		Return(nil, fmt.Errorf("db error"))

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	repoMock.EXPECT().
		HasAdministrator().
		Return(queryMock)

	chain := web.NewSecurityChain[web.EmptyBody]().RuleWhenNoAdminUser(repoMock)
	assert.False(t, chain.ShouldPermitAccess(nil))
}

/**************************
 *** Chained rule tests ***
 **************************/
func TestSecurityChain_ChainedRules_TrueAndFalse(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryMock := mock.NewMockQuery[entity.ModeratorUser](ctrl)
	queryMock.EXPECT().
		Query().
		Times(2).
		Return(&entity.ModeratorUser{Administrator: false}, nil)

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	repoMock.EXPECT().
		FindUserByID(gomock.Eq(userID)).
		Times(2).
		Return(queryMock)

	sessionManagerMock := mock.NewMockSessionManager(ctrl)
	sessionManagerMock.EXPECT().
		VerifyToken(gomock.Eq(token)).
		Times(2).
		Return(&service.SessionUserData{
			Admin:    true,
			IssuedAt: types.Time(time.Unix(1569862420, 0)),
			Username: username,
			ID:       userID,
		})

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(repoMock, sessionManagerMock).RuleAdmin(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer "+token)
	assert.False(t, chain.ShouldPermitAccess(details))
}

func TestSecurityChain_ChainedRules_TrueAndTrue(t *testing.T) {
	ctrl := gomock.NewController(t)

	queryMock := mock.NewMockQuery[entity.ModeratorUser](ctrl)
	queryMock.EXPECT().
		Query().
		Times(2).
		Return(&entity.ModeratorUser{Administrator: true}, nil)

	repoMock := mock.NewMockModeratorUserRepository(ctrl)
	repoMock.EXPECT().
		FindUserByID(gomock.Eq(userID)).
		Times(2).
		Return(queryMock)

	sessionManagerMock := mock.NewMockSessionManager(ctrl)
	sessionManagerMock.EXPECT().
		VerifyToken(gomock.Eq(token)).
		Times(2).
		Return(&service.SessionUserData{
			Admin:    true,
			IssuedAt: types.Time(time.Unix(1569862420, 0)),
			Username: username,
			ID:       userID,
		})

	chain := web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(repoMock, sessionManagerMock).RuleAdmin(repoMock, sessionManagerMock)
	details := &web.HttpRequestDetails[web.EmptyBody]{
		Header: make(http.Header),
	}
	details.Header.Set("Authorization", "Bearer "+token)
	assert.True(t, chain.ShouldPermitAccess(details))
}
