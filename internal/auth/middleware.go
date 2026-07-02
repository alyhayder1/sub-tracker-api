package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/jackc/pgx/v5"

	db "github.com/sub-tracker-hq/sub-tracker-api/db/generated"
	"github.com/sub-tracker-hq/sub-tracker-api/internal/common"
)

type contextKey int

const userContextKey contextKey = iota

// UserFromContext returns the authenticated local user resolved by Middleware.
func UserFromContext(ctx context.Context) (db.User, bool) {
	u, ok := ctx.Value(userContextKey).(db.User)
	return u, ok
}

// Middleware verifies the Clerk bearer session token on every request, then
// resolves the corresponding local users row — provisioning one on first
// login — and stores it in the request context for downstream handlers.
func Middleware(queries *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerToken(r)
			if !ok {
				common.WriteError(w, http.StatusUnauthorized, "missing or malformed Authorization header")
				return
			}

			claims, err := jwt.Verify(r.Context(), &jwt.VerifyParams{Token: token})
			if err != nil {
				common.WriteError(w, http.StatusUnauthorized, "invalid session token")
				return
			}

			localUser, err := resolveUser(r.Context(), queries, claims.Subject)
			if err != nil {
				common.WriteError(w, http.StatusUnauthorized, "could not resolve user")
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, localUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}
	return token, true
}

// resolveUser looks up the local users row for a verified Clerk identity,
// provisioning one on first login via an idempotent upsert (so concurrent
// first logins for the same user can't race into a duplicate/failed insert).
func resolveUser(ctx context.Context, queries *db.Queries, clerkID string) (db.User, error) {
	existing, err := queries.GetUserByClerkID(ctx, clerkID)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return db.User{}, err
	}

	clerkUser, err := user.NewClient(&clerk.ClientConfig{}).Get(ctx, clerkID)
	if err != nil {
		return db.User{}, err
	}

	email := primaryEmail(clerkUser)
	if email == "" {
		return db.User{}, errors.New("auth: clerk user has no primary email address")
	}

	return queries.UpsertUserByClerkID(ctx, db.UpsertUserByClerkIDParams{
		ClerkID:         clerkID,
		Email:           email,
		FullName:        fullName(clerkUser),
		DefaultCurrency: "EUR",
	})
}

func primaryEmail(u *clerk.User) string {
	if u.PrimaryEmailAddressID != nil {
		for _, addr := range u.EmailAddresses {
			if addr.ID == *u.PrimaryEmailAddressID {
				return addr.EmailAddress
			}
		}
	}
	if len(u.EmailAddresses) > 0 {
		return u.EmailAddresses[0].EmailAddress
	}
	return ""
}

func fullName(u *clerk.User) *string {
	first, last := "", ""
	if u.FirstName != nil {
		first = *u.FirstName
	}
	if u.LastName != nil {
		last = *u.LastName
	}
	name := strings.TrimSpace(first + " " + last)
	if name == "" {
		return nil
	}
	return &name
}
