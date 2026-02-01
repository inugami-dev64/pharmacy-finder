package web

import (
	"pharmafinder/db"
	"pharmafinder/service"
	"strings"
)

type RulePredicate[B interface{}] func(*HttpRequestDetails[B]) bool

// SecurityChain rules are checked with AND clause
// which means that all security chain rules must be permitted
// for the secured endpoint to be accessed
type SecurityChain[B interface{}] interface {
	// Allow all users that are authenticated
	RuleAuthenticated(repo db.ModeratorUserRepository, tokenManager service.SessionManager) SecurityChain[B]

	// Allow all users that are administrators
	RuleAdmin(repo db.ModeratorUserRepository, tokenManager service.SessionManager) SecurityChain[B]

	// Allow only when no admin user is present in the data layer
	RuleWhenNoAdminUser(repo db.ModeratorUserRepository) SecurityChain[B]

	// Returns true if all chain rules are satisfied
	// false otherwise
	ShouldPermitAccess(details *HttpRequestDetails[B]) bool
}

type SecurityChainImpl[B interface{}] struct {
	predicates []RulePredicate[B]
}

func NewSecurityChain[B interface{}]() SecurityChain[B] {
	return SecurityChainImpl[B]{
		predicates: make([]RulePredicate[B], 0),
	}
}

func (chain SecurityChainImpl[B]) RuleAuthenticated(repo db.ModeratorUserRepository, tokenManager service.SessionManager) SecurityChain[B] {
	chain.predicates = append(chain.predicates, func(details *HttpRequestDetails[B]) bool {
		bearer := details.Header.Get("Authorization")
		parts := strings.Split(bearer, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return false
		} else {
			data := tokenManager.VerifyToken(parts[1])
			if data == nil {
				return false
			}

			user, err := repo.FindUserByID(data.ID).Query()
			if err != nil || user == nil {
				return false
			}

			details.AuthenticatedUser = user
		}

		return true
	})

	return chain
}

func (chain SecurityChainImpl[B]) RuleAdmin(repo db.ModeratorUserRepository, tokenManager service.SessionManager) SecurityChain[B] {
	chain.predicates = append(chain.predicates, func(details *HttpRequestDetails[B]) bool {
		bearer := details.Header.Get("Authorization")
		parts := strings.Split(bearer, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return false
		} else {
			data := tokenManager.VerifyToken(parts[1])
			if data == nil || !data.Admin {
				return false
			}

			user, err := repo.FindUserByID(data.ID).Query()
			if err != nil || user == nil || !user.Administrator {
				return false
			}

			details.AuthenticatedUser = user
		}

		return true
	})

	return chain
}

func (chain SecurityChainImpl[B]) RuleWhenNoAdminUser(repo db.ModeratorUserRepository) SecurityChain[B] {
	chain.predicates = append(chain.predicates, func(details *HttpRequestDetails[B]) bool {
		hasAdmin, err := repo.HasAdministrator().Query()
		if err != nil || hasAdmin == nil {
			return false
		}

		return !*hasAdmin
	})

	return chain
}

func (chain SecurityChainImpl[B]) ShouldPermitAccess(details *HttpRequestDetails[B]) bool {
	allowAccess := true
	for _, predicate := range chain.predicates {
		allowAccess = allowAccess && predicate(details)
	}

	return allowAccess
}
