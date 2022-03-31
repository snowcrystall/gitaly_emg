package gitalyssh

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/metadata/featureflag"
	gitalyx509 "gitlab.com/gitlab-org/gitaly/v14/internal/x509"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"gitlab.com/gitlab-org/labkit/correlation"
	"gitlab.com/gitlab-org/labkit/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	// GitalyInternalURL is a special URL that indicates Gitaly wants to
	// push or fetch to another Gitaly instance
	GitalyInternalURL = "ssh://gitaly/internal.git"
)

var (
	envInjector = tracing.NewEnvInjector()
)

// UploadPackEnv returns a list of the key=val pairs required to set proper configuration options for upload-pack command.
func UploadPackEnv(ctx context.Context, cfg config.Cfg, req *gitalypb.SSHUploadPackRequest) ([]string, error) {
	env, err := commandEnv(ctx, cfg, req.Repository.StorageName, "upload-pack", req)
	if err != nil {
		return nil, err
	}
	return envInjector(ctx, env), nil
}

func commandEnv(ctx context.Context, cfg config.Cfg, storageName, command string, message proto.Message) ([]string, error) {
	payload, err := protojson.Marshal(message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "commandEnv: marshalling payload failed: %v", err)
	}

	serversInfo, err := helper.ExtractGitalyServers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "commandEnv: extracting Gitaly servers: %v", err)
	}

	storageInfo, ok := serversInfo[storageName]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "commandEnv: no storage info for %s", storageName)
	}

	if storageInfo.Address == "" {
		return nil, status.Errorf(codes.InvalidArgument, "commandEnv: empty gitaly address")
	}

	featureFlagPairs := featureflag.AllFlags(ctx)

	return []string{
		fmt.Sprintf("GITALY_PAYLOAD=%s", payload),
		fmt.Sprintf("GIT_SSH_COMMAND=%s %s", filepath.Join(cfg.BinDir, "gitaly-ssh"), command),
		fmt.Sprintf("GITALY_ADDRESS=%s", storageInfo.Address),
		fmt.Sprintf("GITALY_TOKEN=%s", storageInfo.Token),
		fmt.Sprintf("GITALY_FEATUREFLAGS=%s", strings.Join(featureFlagPairs, ",")),
		fmt.Sprintf("CORRELATION_ID=%s", correlation.ExtractFromContextOrGenerate(ctx)),
		// please see https://github.com/git/git/commit/0da0e49ba12225684b75e86a4c9344ad121652cb for mote details
		"GIT_SSH_VARIANT=simple",
		// Pass through the SSL_CERT_* variables that indicate which
		// system certs to trust
		fmt.Sprintf("%s=%s", gitalyx509.SSLCertDir, os.Getenv(gitalyx509.SSLCertDir)),
		fmt.Sprintf("%s=%s", gitalyx509.SSLCertFile, os.Getenv(gitalyx509.SSLCertFile)),
	}, nil
}
