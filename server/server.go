package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aqyuki/jwt-demo/account"
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

	service account.Service
}

func (x *AccountServer) SignIn(c echo.Context) error {
	// try to convert request body to go structure (SignInRequest)
	reqBody := new(SignInRequest)
	if err := c.Bind(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Failed to extract request body",
			Reason:  err.Error(),
		})
	}

	ctx := c.Request().Context()

	// account verification
	account, err := x.service.SignIn(ctx, reqBody.ID, reqBody.Password)
	if err != nil {
		return c.JSON(http.StatusOK, ErrorResponse{
			Message: "Failed to fetch account",
			Reason:  err.Error(),
		})
	}

	token, err := x.generateJWT(account.ID)
	if err != nil {
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
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Failed to extract request body",
			Reason:  err.Error(),
		})
	}

	ctx := c.Request().Context()
	account, err := x.service.SignUp(ctx, reqBody.ID, reqBody.Password, reqBody.Name, reqBody.Bio)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to register account",
			Reason:  err.Error(),
		})
	}

	token, err := x.generateJWT(account.ID)
	if err != nil {
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
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "Failed to get jwt token from context",
		})
	}
	claims := user.Claims.(*accountJWTClaims)
	tokenID := claims.ID

	// try to get user id from url parameter
	urlID := c.Param("user_id")

	if tokenID != urlID {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "User id in token and url parameter are not matched",
			Reason:  "",
		})
	}

	ctx := c.Request().Context()
	account, err := x.service.FetchAccountInfo(ctx, tokenID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
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
	userGroup := apiGroup.Group("users")

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
	server.service = s
	return server
}
