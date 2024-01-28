package api

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/GDGVIT/configgy-backend/api/pkg/auth"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/authsvc"
	"github.com/GDGVIT/configgy-backend/pkg/logger"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type EchoContext interface {
	context.Context
}

//go:embed openapi-spec.yaml
var openAPISpec string

type EchoServer interface {
	Start(string) error
	Shutdown(ctx context.Context) error
}

type Service struct {
	ctx    context.Context
	opts   *Options
	logger logger.Logger
	server EchoServer
	spec   map[string]interface{}

	Services    Services
	Middlewares Middlewares
}

// Dependencies - dependencies for Service constructor
type Dependencies struct {
	Logger        logger.Logger
	EchoServer    EchoServer
	MessageBroker MessageBroker
	GormDB        *gorm.DB
	DB            *tables.DB
	Services      Services
	Middlewares   Middlewares
}

type RabbitMqService struct {
	HostPort string
}

type MessageBroker interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
}

type Options struct {
	Path                string
	Port                int
	ShutdownGracePeriod time.Duration
}

type Middlewares interface {
	CheckIfUser(c context.Context, token string) (*authsvc.AuthData, error)
	CheckIfAdmin(c context.Context, token string) (*authsvc.AuthData, error)
}

type Services struct {
	AuthSvc       AuthService
	AuthzSvc      AuthzService
	UserSvc       UserService
	GroupSvc      GroupService
	VaultSvc      VaultService
	CredentialSvc CredentialService
}

type AuthService interface {
	GenerateToken(c context.Context, req authsvc.TokenReq) (authsvc.TokenRes, error)
}

// GetFlagSet returns flag set for Options
func (o *Options) GetFlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("apiOptions", pflag.ExitOnError)
	flags.StringVar(&o.Path, "path", o.Path, "path to serve API on")
	flags.IntVar(&o.Port, "port", o.Port, "port to serve API on")
	flags.DurationVar(&o.ShutdownGracePeriod, "shutdown-grace-period", o.ShutdownGracePeriod, "shutdown grace period")
	return flags
}

// NewService - constructor for Service
func NewService(ctx context.Context, opts *Options, deps *Dependencies) (*Service, error) {
	svc := &Service{
		ctx:         ctx,
		opts:        opts,
		logger:      deps.Logger,
		server:      deps.EchoServer,
		spec:        make(map[string]interface{}),
		Services:    deps.Services,
		Middlewares: deps.Middlewares,
	}
	svc.server = svc.createServer()
	err := yaml.Unmarshal([]byte(openAPISpec), svc.spec)
	return svc, err
}

// Start starts the API
func (svc *Service) Start() {
	go func() {
		addr := fmt.Sprintf(":%d", svc.opts.Port)
		if err := svc.server.Start(addr); err != nil {
			logger.Println(err)
		}
	}()
}

// Close closes the API
func (svc *Service) Close() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), svc.opts.ShutdownGracePeriod)
	defer cancel()
	return svc.server.Shutdown(ctx)
}

func (svc *Service) createServer() EchoServer {
	server := echo.New()
	server.JSONSerializer = &jsonSerializer{}
	server.GET(svc.opts.Path+"/openapi-spec.json", svc.GetOpenAPISpec)
	apiGroup := server.Group("")
	RegisterHandlersWithBaseURL(apiGroup, svc, svc.opts.Path)
	return server
}

// GetOpenAPISpec - (GET /openapi-spec.json)
func (svc *Service) GetOpenAPISpec(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, svc.spec, "  ")
}

func getUserPIDFromAuthToken(c echo.Context) (string, error) {
	authToken, err := auth.GetAuthorizationHeader(c.Request())
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	err = auth.ValidateToken(authToken)
	if err != nil {
		message := "invalid token"
		return "", errors.New(message)
	}
	authData, err := auth.GetAuthDataFromToken(authToken)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	return authData.UserPID, nil
}
