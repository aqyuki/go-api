package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/aqyuki/jwt-demo/account"
	"github.com/aqyuki/jwt-demo/logging"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const tokenContextKey = "jwt_token"

// accountJWTClaims is a jwt claims for account
type accountJWTClaims struct {
	ID string `json:"user_id"`
	jwt.RegisteredClaims
}

// AccountServer provides account management API
type AccountServer struct {
	server *echo.Echo
	secret []byte

	logger  *slog.Logger
	service account.Service
}

func (x *AccountServer) SignIn(c echo.Context) error {
	// try to convert request body to go structure (SignInRequest)
	reqBody := new(SignInRequest)
	if err := c.Bind(reqBody); err != nil {
		x.logger.Info("Failed to extract request body", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Failed to extract request body",
			Reason:  err.Error(),
		})
	}

	x.logger.Info("Sign in request", slog.Any("request", reqBody))

	ctx := logging.ContextWithLogger(context.Background(), x.logger)
	// account verification
	account, err := x.service.SignIn(ctx, reqBody.ID, reqBody.Password)
	if err != nil {
		x.logger.Info("Failed to fetch account", slog.Any("error", err))
		return c.JSON(http.StatusOK, ErrorResponse{
			Message: "Failed to fetch account",
			Reason:  err.Error(),
		})
	}

	token, err := x.generateJWT(account.ID)
	if err != nil {
		x.logger.Error("Failed to generate jwt token", slog.Any("error", err))
		return c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{
				Message: "Failed to generate jwt token",
				Reason:  err.Error(),
			})
	}

	return c.JSON(http.StatusOK, AccountVerifyResponse{
		Token: token,
	})
}

func (x *AccountServer) SignUp(c echo.Context) error {
	// try to convert request body to go structure (SignUpRequest)
	reqBody := new(SignUpRequest)
	if err := c.Bind(reqBody); err != nil {
		x.logger.Info("Failed to extract request body", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Failed to extract request body",
			Reason:  err.Error(),
		})
	}
	x.logger.Info("Sign up request", slog.Any("request", reqBody))

	ctx := logging.ContextWithLogger(context.Background(), x.logger)
	account, err := x.service.SignUp(ctx, reqBody.ID, reqBody.Password, reqBody.Name, reqBody.Bio)
	if err != nil {
		x.logger.Info("Failed to register account", slog.Any("error", err))
		return c.JSON(http.StatusOK, ErrorResponse{
			Message: "Failed to register account",
			Reason:  err.Error(),
		})
	}

	token, err := x.generateJWT(account.ID)
	if err != nil {
		x.logger.Error("Failed to generate jwt token", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to generate jwt token",
			Reason:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, AccountVerifyResponse{
		Token: token,
	})
}

func (x *AccountServer) GetAccountInfo(c echo.Context) error {
	// try to get jwt token from context, and extract user id from it.
	user, ok := c.Get(tokenContextKey).(*jwt.Token)
	if !ok {
		x.logger.Info("Failed to get jwt token from context")
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "Failed to get jwt token from context",
		})
	}
	claims := user.Claims.(*accountJWTClaims)
	tokenID := claims.ID

	// try to get user id from url parameter
	urlID := c.Param("user_id")

	if tokenID != urlID {
		slog.Info("User id in token and url parameter are not matched", slog.String("tokenID", tokenID), slog.String("urlID", urlID))
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "User id in token and url parameter are not matched",
			Reason:  "",
		})
	}

	ctx := logging.ContextWithLogger(c.Request().Context(), x.logger)
	account, err := x.service.FetchAccountInfo(ctx, tokenID)
	if err != nil {
		x.logger.Info("Failed to fetch account", slog.Any("error", err))
		return c.JSON(http.StatusOK, ErrorResponse{
			Message: "Failed to fetch account",
			Reason:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, AccountInfoResponse{
		ID:   account.ID,
		Name: account.Name,
		Bio:  account.Bio,
	})
}

func (x *AccountServer) Start(addr string) error {
	return x.server.Start(addr)
}

func (x *AccountServer) Shutdown(ctx context.Context) error {
	return x.server.Shutdown(ctx)
}

func (x *AccountServer) generateJWT(id string) (string, error) {
	claims := &accountJWTClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(x.secret)
	if err != nil {
		return "", err
	}
	return t, nil
}

// NewAccountServer creates new AccountServer
func NewAccountServer(key []byte, l *slog.Logger, s account.Service) *AccountServer {
	server := new(AccountServer)

	e := echo.New()

	// set group
	apiGroup := e.Group("/api/v1")
	userGroup := apiGroup.Group("/users")

	// set common middleware
	e.Use(middleware.Recover())
	e.Use(NewCustomLogger(l))

	// set jwt middleware on group
	config := echojwt.Config{
		ContextKey: tokenContextKey,
		SigningKey: key,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(accountJWTClaims)
		},
	}
	userGroup.Use(echojwt.WithConfig(config))

	// register API
	apiGroup.POST("/signin", server.SignIn)
	apiGroup.POST("/signup", server.SignUp)
	userGroup.GET("/:user_id", server.GetAccountInfo)

	server.server = e
	server.secret = key
	server.logger = l
	server.service = s
	return server
}
