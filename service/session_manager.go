package service

import (
	"crypto/rand"
	"encoding/hex"
	"pharmafinder/db/entity"
	"pharmafinder/types"
	"pharmafinder/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// By default session tokens are valid for 30 days
const SESSION_TTL = time.Hour * 24 * 30

type SessionUserData struct {
	ID       types.UUID `json:"sub"`
	Username string     `json:"name"`
	Admin    bool       `json:"admin"`
	IssuedAt types.Time `json:"iat"`
}

type SessionManager interface {
	// Creates a new authentication session token for provided user
	NewSessionToken(user *entity.ModeratorUser) string

	// Verifies if provided token is valid and if it is,
	// returns an instance of SessionUserData
	//
	// If the token is not valid, then the returned data is nil
	VerifyToken(token string) *SessionUserData
}

// Initializes the HS256 key to use for JWT token signing
//
// If provided key value is nil then the key is randomly generated.
// Otherwise at most 32 bytes of key are used as the HS256 key with
// zero padding if the key length is less than 32 bytes.
func initializeHS256Key(key []byte) []byte {
	hs256Key := make([]byte, 32)
	if key == nil {
		rand.Read(hs256Key)
	} else {
		for i := range min(len(key), 32) {
			hs256Key[i] = key[i]
		}
	}

	return hs256Key
}

func ProvideSessionManager() SessionManager {
	keyStr := utils.Getenv("JWT_KEY", "")
	key, err := hex.DecodeString(keyStr)
	if len(key) == 0 || err != nil {
		key = initializeHS256Key(nil)
	}

	return &JWTSessionManagerImpl{
		hs256Key: initializeHS256Key(key),
		logger:   utils.GetLogger("SERVICE"),
	}
}

type JWTSessionManagerImpl struct {
	hs256Key []byte
	logger   zerolog.Logger
}

func (sess *JWTSessionManagerImpl) NewSessionToken(user *entity.ModeratorUser) string {
	sess.logger.Debug().Msgf("Issuing a new token to moderator '%s'", user.Username)
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":   user.ID,
			"name":  user.Username,
			"admin": user.Administrator,
			"iat":   jwt.NumericDate{Time: time.Now().UTC()},
		})

	s, err := t.SignedString(sess.hs256Key)
	if err != nil {
		sess.logger.Error().Msgf("Failed to issue a new token to moderator '%s'", user.Username)
		return ""
	}

	return s
}

func (sess *JWTSessionManagerImpl) VerifyToken(token string) *SessionUserData {
	jwtToken, err := jwt.Parse(
		token,
		func(t *jwt.Token) (any, error) {
			return sess.hs256Key, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuedAt())

	if err != nil {
		sess.logger.Warn().Msg("Failed to parse JWT token")
		return nil
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		var ok [4]bool
		var subVal, nameVal, adminVal, iatVal any
		subVal, ok[0] = claims["sub"]
		nameVal, ok[1] = claims["name"]
		adminVal, ok[2] = claims["admin"]
		iatVal, ok[3] = claims["iat"]

		if !ok[0] || !ok[1] || !ok[2] || !ok[3] {
			return nil
		}

		userData := SessionUserData{}
		if sub, ok := subVal.(string); ok {
			id, _ := uuid.Parse(sub)
			userData.ID = types.UUID(id)
		} else {
			return nil
		}

		if name, ok := nameVal.(string); ok {
			userData.Username = name
		} else {
			return nil
		}

		if admin, ok := adminVal.(bool); ok {
			userData.Admin = admin
		} else {
			return nil
		}

		if iat, ok := iatVal.(float64); ok {
			userData.IssuedAt = types.Time(time.Unix(int64(iat), 0))
		} else {
			return nil
		}

		// validate that the token has not expired
		if time.Now().UTC().Sub(time.Time(userData.IssuedAt)).Seconds() < SESSION_TTL.Seconds() {
			return &userData
		}
	}

	return nil
}
