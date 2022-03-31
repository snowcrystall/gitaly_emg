package auth

import (
	"context"
	"time"

	grpcmwauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	gitalyauth "gitlab.com/gitlab-org/gitaly/v14/auth"
	gitalycfgauth "gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	authCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitaly_authentications_total",
			Help: "Counts of of Gitaly request authentication attempts",
		},
		[]string{"enforced", "status"},
	)
)

// StreamServerInterceptor checks for Gitaly bearer tokens.
func StreamServerInterceptor(conf gitalycfgauth.Config) grpc.StreamServerInterceptor {
	return grpcmwauth.StreamServerInterceptor(checkFunc(conf))
}

// UnaryServerInterceptor checks for Gitaly bearer tokens.
func UnaryServerInterceptor(conf gitalycfgauth.Config) grpc.UnaryServerInterceptor {
	return grpcmwauth.UnaryServerInterceptor(checkFunc(conf))
}

func checkFunc(conf gitalycfgauth.Config) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		if len(conf.Token) == 0 {
			countStatus("server disabled authentication", conf.Transitioning).Inc()
			return ctx, nil
		}

		err := gitalyauth.CheckToken(ctx, conf.Token, time.Now())
		switch status.Code(err) {
		case codes.OK:
			countStatus(okLabel(conf.Transitioning), conf.Transitioning).Inc()
		case codes.Unauthenticated:
			countStatus("unauthenticated", conf.Transitioning).Inc()
		case codes.PermissionDenied:
			countStatus("denied", conf.Transitioning).Inc()
		default:
			countStatus("invalid", conf.Transitioning).Inc()
		}

		if conf.Transitioning {
			err = nil
		}

		return ctx, err
	}
}

func okLabel(transitioning bool) string {
	if transitioning {
		// This special value is an extra warning sign to administrators that
		// authentication is currently not enforced.
		return "would be ok"
	}
	return "ok"
}

func countStatus(status string, transitioning bool) prometheus.Counter {
	enforced := "true"
	if transitioning {
		enforced = "false"
	}
	return authCount.WithLabelValues(enforced, status)
}
