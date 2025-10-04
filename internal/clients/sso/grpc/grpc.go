package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	ssov1 "github.com/Sheridanlk/protos/gen/go/sso"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type Client struct {
	api ssov1.AuthClient
	log *log.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	// TODO: переделать с защищённым соединением
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Client{
		api: ssov1.NewAuthClient(cc),
	}, nil
}

func (c *Client) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "grpc.IsAdmin"

	resp, err := c.api.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: userID,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
			}
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.IsAdmin, nil
}

func (c *Client) Login(ctx context.Context, email string, password string, appID int32) (string, error) {
	const op = "grpc.Login"

	resp, err := c.api.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
			}
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return resp.Token, nil
}

func (c *Client) Regiser(ctx context.Context, email string, password string) (int64, error) {
	const op = "grpc.Register"

	resp, err := c.api.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return resp.UserId, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}
