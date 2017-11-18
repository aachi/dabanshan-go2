package endpoint

import (
	"context"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	rl "github.com/juju/ratelimit"
	m_user "github.com/laidingqing/dabanshan-go/svcs/user/model"
	"github.com/laidingqing/dabanshan-go/svcs/user/service"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	GetUserEndpoint  endpoint.Endpoint
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint    endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc service.Service, logger log.Logger, duration metrics.Histogram, trace stdopentracing.Tracer) Set {
	var (
		getUserEndpoint  endpoint.Endpoint
		registerEndpoint endpoint.Endpoint
		loginEndpoint    endpoint.Endpoint
	)
	{
		getUserEndpoint = MakeGetUserEndpoint(svc)
		getUserEndpoint = ratelimit.NewTokenBucketLimiter(rl.NewBucketWithRate(1, 1))(getUserEndpoint)
		getUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getUserEndpoint)
		getUserEndpoint = opentracing.TraceServer(trace, "GetUser")(getUserEndpoint)
		getUserEndpoint = LoggingMiddleware(log.With(logger, "method", "GetUser"))(getUserEndpoint)
		getUserEndpoint = InstrumentingMiddleware(duration.With("method", "GetUser"))(getUserEndpoint)
	}
	{
		registerEndpoint = MakeRegisterEndpoint(svc)
		registerEndpoint = ratelimit.NewTokenBucketLimiter(rl.NewBucketWithRate(1, 1))(registerEndpoint)
		registerEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(registerEndpoint)
		registerEndpoint = opentracing.TraceServer(trace, "Register")(registerEndpoint)
		registerEndpoint = LoggingMiddleware(log.With(logger, "method", "Register"))(registerEndpoint)
		registerEndpoint = InstrumentingMiddleware(duration.With("method", "Register"))(registerEndpoint)
	}
	{
		loginEndpoint = MakeLoginEndpoint(svc)
		loginEndpoint = ratelimit.NewTokenBucketLimiter(rl.NewBucketWithRate(1, 1))(loginEndpoint)
		loginEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(loginEndpoint)
		loginEndpoint = opentracing.TraceServer(trace, "Login")(loginEndpoint)
		loginEndpoint = LoggingMiddleware(log.With(logger, "method", "Login"))(loginEndpoint)
		loginEndpoint = InstrumentingMiddleware(duration.With("method", "Login"))(loginEndpoint)
	}

	return Set{
		GetUserEndpoint:  getUserEndpoint,
		RegisterEndpoint: registerEndpoint,
		LoginEndpoint:    loginEndpoint,
	}
}

// GetUser implements the service interface, so Set may be used as a service.
func (s Set) GetUser(ctx context.Context, a string) (m_user.GetUserResponse, error) {
	resp, err := s.GetUserEndpoint(ctx, m_user.GetUserRequest{A: a})
	if err != nil {
		return m_user.GetUserResponse{}, err
	}
	response := resp.(m_user.GetUserResponse)
	return response, response.Err
}

// Register implements the service interface,
func (s Set) Register(ctx context.Context, us m_user.RegisterRequest) (r m_user.RegisterUserResponse, err error) {
	resp, err := s.RegisterEndpoint(ctx, us)
	if err != nil {
		return m_user.RegisterUserResponse{ID: ""}, err
	}
	response := resp.(m_user.RegisterUserResponse)
	return response, err
}

// Login implements the service interface.
func (s Set) Login(ctx context.Context, login m_user.LoginRequest) (m_user.LoginResponse, error) {
	resp, err := s.LoginEndpoint(ctx, login)
	if err != nil {
		return m_user.LoginResponse{}, err
	}
	response := resp.(m_user.LoginResponse)
	return response, err
}

// MakeGetUserEndpoint constructs a GetUser endpoint wrapping the service.
func MakeGetUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(m_user.GetUserRequest)
		v, err := s.GetUser(ctx, req.A)
		return v, err
	}
}

// MakeRegisterEndpoint constructs a Register endpoint wrapping the service.
func MakeRegisterEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(m_user.RegisterRequest)
		v, err := s.Register(ctx, req)
		return v, err
	}
}

// MakeLoginEndpoint constructs a Login endpoint wrapping the service.
func MakeLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(m_user.LoginRequest)
		v, err := s.Login(ctx, req)
		return v, err
	}
}
