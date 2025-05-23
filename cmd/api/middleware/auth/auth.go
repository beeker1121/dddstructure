package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	apictx "dddstructure/cmd/api/context"
	"dddstructure/cmd/api/errors"
	"dddstructure/proto"
	serverrors "dddstructure/service/errors"

	"github.com/dgrijalva/jwt-go"
)

// key is the key type used by this package for the request context.
type key int

// AuthKey is the key used for storing and retrieving the user data from the
// request context.
var AuthKey key = 1

// TokenClaims defines the custom claims we use for the JWT.
type TokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// NewJWT creates and returns a new signed JWT.
func NewJWT(ac *apictx.Context, userPassword string, uid uint) (string, error) {
	// Set expiry time.
	issued := time.Now()
	expires := issued.Add(time.Minute * ac.Config.JWTExpiryTime)

	// Create the claims.
	claims := &TokenClaims{
		uid,
		jwt.StandardClaims{
			IssuedAt:  issued.Unix(),
			ExpiresAt: expires.Unix(),
		},
	}

	// Create and sign the token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(GetJWTSigningKey(ac.Config.JWTSecret, userPassword))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// AuthenticateEndpoint is the middleware for authenticating API requests.
//
// This function will first try to determine the type of authorization being
// requested, and then either authorize via a JWT or an API key.
//
// JWTs are passed via the Authorization header as a Bearer token.
//
// API keys should be passed via the Authorization header using Basic Auth.
//
// Currently, the only supported authorization method is via JWTs.
func AuthenticateEndpoint(ac *apictx.Context, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &proto.User{}
		var err error

		// Get the Authorization header.
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")

		// Check for either Bearer or Basic authoriation type.
		if authHeader[0] != "Bearer" && authHeader[0] != "Basic" {
			errors.Default(ac.Logger, w, errors.New(http.StatusUnauthorized, "", ErrUnauthorized.Error()))
			return
		}

		if len(authHeader) == 2 && authHeader[0] == "Bearer" {
			// Try authorization via JWT Authorization Bearer header first.
			u, err = GetUserFromJWT(ac, authHeader[1])
			if err == ErrJWTUnauthorized {
				ac.Logger.Error("API authorization via JWT failure")
				errors.Default(ac.Logger, w, errors.New(http.StatusUnauthorized, "", err.Error()))
				return
			} else if err != nil {
				ac.Logger.Error("auth.GetUserFromJWT() error",
					slog.Any("error", err))
				errors.Default(ac.Logger, w, errors.ErrInternalServerError)
				return
			}
		} else {
			// Get the user from the API key.
			ac.Logger.Error("API key authorization not implemented")
			errors.Default(ac.Logger, w, errors.New(http.StatusUnauthorized, "", "API key authorization not implemented"))
			return
		}

		// Pass user to request context and call next handler.
		ctx := context.WithValue(r.Context(), AuthKey, u)
		h(w, r.WithContext(ctx))
	}
}

// GetUserFromJWT retrieves the user from the given JWT.
func GetUserFromJWT(ac *apictx.Context, headerToken string) (*proto.User, error) {
	// Get the signing key for this user from the JWT claims.
	signingKey, err := GetUserSigningKey(ac, headerToken)
	if err != nil {
		return nil, err
	}

	// Parse the token.
	token, err := jwt.ParseWithClaims(headerToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrJWTUnauthorized
		}

		return signingKey, nil
	})
	if err != nil {
		return nil, ErrJWTUnauthorized
	}

	// Get token claims and check token validity.
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrJWTUnauthorized
	}

	// Get the user using the UserID claim.
	u, err := ac.Service.User.GetByID(claims.UserID)
	switch {
	case err == serverrors.ErrUserNotFound:
		return nil, ErrJWTUnauthorized
	case err != nil:
		return nil, err
	}

	return u, nil
}

// GetUserSigningKey creates the unique JWT signing key for the given user
// using the JWT secret and their current hashed password.
//
// The claims are parsed from the payload portion of the token to get the
// user ID, which is then used to retrieve the hashed user password.
func GetUserSigningKey(ac *apictx.Context, token string) ([]byte, error) {
	// Split token.
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return []byte{}, ErrJWTUnauthorized
	}

	// Parse claims.
	claimBytes, err := jwt.DecodeSegment(parts[1])
	if err != nil {
		return []byte{}, err
	}

	// Unmarshal into TokenClaims type.
	var claims TokenClaims
	if err := json.Unmarshal(claimBytes, &claims); err != nil {
		return []byte{}, err
	}

	// Get the user from the UserID claim.
	u, err := ac.Service.User.GetByID(claims.UserID)
	switch {
	case err == serverrors.ErrUserNotFound:
		return []byte{}, ErrJWTUnauthorized
	case err != nil:
		return []byte{}, err
	}

	return GetJWTSigningKey(ac.Config.JWTSecret, u.Password), nil
}

// GetJWTSigningKey returns the JWT signing key.
//
// It is constructed using the user's hashed password and the application JWT
// secret.
func GetJWTSigningKey(jwtSecret, password string) []byte {
	return []byte(jwtSecret + password)
}

// GetUserFromRequest retrieves the authenticated user from the request
// context.
func GetUserFromRequest(r *http.Request) (*proto.User, error) {
	u, ok := r.Context().Value(AuthKey).(*proto.User)
	if !ok {
		return nil, fmt.Errorf("could not type assert user from request context")
	}
	return u, nil
}
