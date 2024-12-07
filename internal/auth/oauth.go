package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/utility"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
	"golang.org/x/oauth2"
)

type Client struct {
	Provider *oidc.Provider
	OIDC     *oidc.IDTokenVerifier
	Oauth    oauth2.Config
}
type Config struct {
	BaseURL      string // Authorization base url
	ClientID     string // client id oauth
	RedirectURI  string // valid redirect uri
	ClientSecret string // optional
}

func New(ctx context.Context, config *Config) (*Client, error) {
	// Construct the provider URI oauth
	providerURL := config.BaseURL
	// Google's OAuth 2.0 endpoint
	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}
	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})
	// Configure an OpenID Connect aware OAuth2 client
	oauth2 := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURI,
		Endpoint:     provider.Endpoint(),
		Scopes: []string{
			oidc.ScopeOpenID,
			"profile",
			"email",
		},
	}
	return &Client{
		Oauth:    oauth2,
		OIDC:     verifier,
		Provider: provider,
	}, nil
}

func (a *Client) HandlerRedirect(ctx *gin.Context, sqlDB *sql.DB, stateAuthStore store.State) {
	stateID, err := utility.GenerateRandomBase64Str()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	err = sqldb.WithinTx(ctx, sqlDB, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		return stateAuthStore.Insert(txContext, stateID)
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.Redirect(http.StatusFound, a.Oauth.AuthCodeURL(stateID))
}

func (a *Client) HandlerCallback(ctx *gin.Context,
	db *sql.DB,
	stateAuthStore store.State, sessionStore store.Session) {
	sessionData, err := a.processCallback(ctx, db, stateAuthStore, sessionStore)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Set session cookie
	// Parameters:
	// 1. name: cookie name
	// 2. value: cookie value (userID/session ID)
	// 3. maxAge: cookie duration in seconds
	// 4. path: cookie path
	// 5. domain: cookie domain
	// 6. secure: only send cookie over HTTPS
	// 7. httpOnly: prevent JavaScript access to cookie
	ctx.SetCookie(
		"session_id",          // name
		sessionData.SessionID, // value
		3600,                  // maxAge (1 hour)
		"/",                   // path
		"",                    // domain
		true,                  // secure
		true,                  // httpOnly
	)
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (a *Client) processCallback(
	ctx *gin.Context, db *sql.DB,
	stateAuthStore store.State, sessionStore store.Session) (*store.SessionData, error) {
	processorCallback := &processorCallback{stateAuthStore: stateAuthStore}
	stateData, err := processorCallback.verifyState(ctx)
	if err != nil {
		return nil, err
	}
	authorizationCode := ctx.Query("code")
	if authorizationCode == "" {
		return nil, errors.New("authorization code is required")
	}
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("grant_type", "authorization_code"),
	}
	var oauth2Token *oauth2.Token
	oauth2Token, err = a.Oauth.Exchange(ctx, authorizationCode, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange the token, got : %w", err)
	}
	idTokenString, err := processorCallback.getIDToken(oauth2Token)
	if err != nil {
		return nil, err
	}
	userInfoData, err := processorCallback.verifyIDTokenAndGetClaims(ctx, idTokenString, a.OIDC)
	if err != nil {
		return nil, err
	}
	insertedSessionData := &store.SessionData{
		SessionID:    uuid.NewString(),
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
		UserEmail:    userInfoData.Email,
		UserFullName: userInfoData.Name,
		CreatedAt:    time.Now().UTC(),
		UserID:       userInfoData.Sub,
		ExpiresAt:    time.Now().UTC().Add(5 * time.Minute),
	}
	err = sqldb.WithinTx(ctx, db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		if err = stateAuthStore.Delete(txContext, stateData.ID); err != nil {
			return fmt.Errorf("failed to reset state data, got : %w", err)
		}
		return sessionStore.Insert(txContext, insertedSessionData)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert session, got : %w", err)
	}
	return insertedSessionData, nil
}

type processorCallback struct {
	stateAuthStore store.State
}

func (a *processorCallback) verifyState(ctx *gin.Context) (*store.StateData, error) {
	stateID := ctx.Query("state")
	if stateID == "" {
		return nil, errors.New("stateID is required")
	}
	_, err := a.stateAuthStore.FindOne(ctx, stateID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("stateID is invalid or not found")
		}
		return nil, fmt.Errorf("failed to get state data, got : %w", err)
	}
	return &store.StateData{ID: stateID}, nil
}
func (a *processorCallback) getIDToken(oauth2Token *oauth2.Token) (string, error) {
	// Get raw ID token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("no id_token in response")
	}
	return rawIDToken, nil
}

type userInfoClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Sub   string `json:"sub"`
}

func (a *processorCallback) verifyIDTokenAndGetClaims(
	ctx *gin.Context,
	rawIDToken string,
	verifier *oidc.IDTokenVerifier,
) (
	*userInfoClaims, error) {
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify id token, got : %w", err)
	}
	userInfoClaims := &userInfoClaims{}
	if err := idToken.Claims(&userInfoClaims); err != nil {
		return nil, fmt.Errorf("failed to extract claims, got : %w", err)
	}
	return userInfoClaims, nil
}
