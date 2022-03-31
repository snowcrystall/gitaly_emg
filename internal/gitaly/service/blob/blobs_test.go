package blob

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/quarantine"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testassert"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"gitlab.com/gitlab-org/gitaly/v14/streamio"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	gitattributesOID  = "b680596c9f3a3c834b933aef14f94a0ab9fa604a"
	gitattributesData = "/custom-highlighting/*.gitlab-custom gitlab-language=ruby\n*.lfs filter=lfs diff=lfs merge=lfs -text\n"
	gitattributesSize = int64(len(gitattributesData))

	gitignoreOID  = "dfaa3f97ca337e20154a98ac9d0be76ddd1fcc82"
	gitignoreSize = 241
)

func TestListBlobs(t *testing.T) {
	// In order to get deterministic behaviour with regards to how streamio splits up values,
	// we're going to limit the write batch size to something smallish. Otherwise, tests will
	// depend on both the batch size of streamio and the batch size of io.Copy.
	defer func(oldValue int) {
		streamio.WriteBufferSize = oldValue
	}(streamio.WriteBufferSize)
	streamio.WriteBufferSize = 200

	ctx, cancel := testhelper.Context()
	defer cancel()

	cfg, repoProto, repoPath, client := setup(t)

	bigBlobContents := bytes.Repeat([]byte{1}, streamio.WriteBufferSize*2+1)
	bigBlobOID := gittest.WriteBlob(t, cfg, repoPath, bigBlobContents)

	for _, tc := range []struct {
		desc          string
		revisions     []string
		limit         uint32
		bytesLimit    int64
		expectedErr   error
		expectedBlobs []*gitalypb.ListBlobsResponse_Blob
	}{
		{
			desc:        "missing revisions",
			revisions:   []string{},
			expectedErr: status.Error(codes.InvalidArgument, "missing revisions"),
		},
		{
			desc: "invalid revision",
			revisions: []string{
				"--foobar",
			},
			expectedErr: status.Error(codes.InvalidArgument, "invalid revision: \"--foobar\""),
		},
		{
			desc: "single blob",
			revisions: []string{
				lfsPointer1,
			},
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: lfsPointer1, Size: lfsPointers[lfsPointer1].Size},
			},
		},
		{
			desc: "multiple blobs",
			revisions: []string{
				lfsPointer1,
				lfsPointer2,
				lfsPointer3,
			},
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: lfsPointer1, Size: lfsPointers[lfsPointer1].Size},
				{Oid: lfsPointer2, Size: lfsPointers[lfsPointer2].Size},
				{Oid: lfsPointer3, Size: lfsPointers[lfsPointer3].Size},
			},
		},
		{
			desc: "tree",
			revisions: []string{
				"b95c0fad32f4361845f91d9ce4c1721b52b82793",
			},
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: "93e123ac8a3e6a0b600953d7598af629dec7b735", Size: 59},
			},
		},
		{
			desc: "revision range",
			revisions: []string{
				"master",
				"^master~2",
			},
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: "723c2c3f4c8a2a1e957f878c8813acfc08cda2b6", Size: 1219696},
			},
		},
		{
			desc: "pseudorevisions",
			revisions: []string{
				"master",
				"--not",
				"--all",
			},
			expectedBlobs: nil,
		},
		{
			desc: "revision with limit",
			revisions: []string{
				"master",
			},
			limit: 2,
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: gitattributesOID, Size: gitattributesSize},
				{Oid: gitignoreOID, Size: gitignoreSize},
			},
		},
		{
			desc: "revision with limit and bytes limit",
			revisions: []string{
				"master",
			},
			limit:      2,
			bytesLimit: 5,
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: gitattributesOID, Size: gitattributesSize, Data: []byte("/cust")},
				{Oid: gitignoreOID, Size: gitignoreSize, Data: []byte("*.rbc")},
			},
		},
		{
			desc: "revision with path",
			revisions: []string{
				"master:.gitattributes",
			},
			bytesLimit: -1,
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: gitattributesOID, Size: gitattributesSize, Data: []byte(gitattributesData)},
			},
		},
		{
			desc: "complete contents via negative bytes limit",
			revisions: []string{
				lfsPointer1,
				lfsPointer2,
				lfsPointer3,
			},
			bytesLimit: -1,
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: lfsPointer1, Size: lfsPointers[lfsPointer1].Size, Data: lfsPointers[lfsPointer1].Data},
				{Oid: lfsPointer2, Size: lfsPointers[lfsPointer2].Size, Data: lfsPointers[lfsPointer2].Data},
				{Oid: lfsPointer3, Size: lfsPointers[lfsPointer3].Size, Data: lfsPointers[lfsPointer3].Data},
			},
		},
		{
			desc: "contents truncated by bytes limit",
			revisions: []string{
				lfsPointer1,
				lfsPointer2,
				lfsPointer3,
			},
			bytesLimit: 10,
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: lfsPointer1, Size: lfsPointers[lfsPointer1].Size, Data: lfsPointers[lfsPointer1].Data[:10]},
				{Oid: lfsPointer2, Size: lfsPointers[lfsPointer2].Size, Data: lfsPointers[lfsPointer2].Data[:10]},
				{Oid: lfsPointer3, Size: lfsPointers[lfsPointer3].Size, Data: lfsPointers[lfsPointer3].Data[:10]},
			},
		},
		{
			desc: "bytes limit exceeding total blob content size",
			revisions: []string{
				lfsPointer1,
				lfsPointer2,
				lfsPointer3,
			},
			bytesLimit: 9000,
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: lfsPointer1, Size: lfsPointers[lfsPointer1].Size, Data: lfsPointers[lfsPointer1].Data},
				{Oid: lfsPointer2, Size: lfsPointers[lfsPointer2].Size, Data: lfsPointers[lfsPointer2].Data},
				{Oid: lfsPointer3, Size: lfsPointers[lfsPointer3].Size, Data: lfsPointers[lfsPointer3].Data},
			},
		},
		{
			desc: "bytes limit partially exceeding limit",
			revisions: []string{
				lfsPointer1,
				lfsPointer2,
				lfsPointer3,
			},
			bytesLimit: 128,
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: lfsPointer1, Size: lfsPointers[lfsPointer1].Size, Data: lfsPointers[lfsPointer1].Data[:128]},
				{Oid: lfsPointer2, Size: lfsPointers[lfsPointer2].Size, Data: lfsPointers[lfsPointer2].Data},
				{Oid: lfsPointer3, Size: lfsPointers[lfsPointer3].Size, Data: lfsPointers[lfsPointer3].Data},
			},
		},
		{
			desc: "blob with content bigger than a gRPC message",
			revisions: []string{
				bigBlobOID.String(),
				lfsPointer1,
			},
			bytesLimit: -1,
			expectedBlobs: []*gitalypb.ListBlobsResponse_Blob{
				{Oid: bigBlobOID.String(), Size: int64(len(bigBlobContents)), Data: bigBlobContents[:streamio.WriteBufferSize]},
				{Data: bigBlobContents[streamio.WriteBufferSize : 2*streamio.WriteBufferSize]},
				{Data: bigBlobContents[2*streamio.WriteBufferSize:]},
				{Oid: lfsPointer1, Size: lfsPointers[lfsPointer1].Size, Data: lfsPointers[lfsPointer1].Data},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			stream, err := client.ListBlobs(ctx, &gitalypb.ListBlobsRequest{
				Repository: repoProto,
				Revisions:  tc.revisions,
				Limit:      tc.limit,
				BytesLimit: tc.bytesLimit,
			})
			require.NoError(t, err)

			var blobs []*gitalypb.ListBlobsResponse_Blob
			for {
				resp, err := stream.Recv()
				if err != nil {
					if !errors.Is(err, io.EOF) {
						testassert.GrpcEqualErr(t, tc.expectedErr, err)
					}
					break
				}

				blobs = append(blobs, resp.Blobs...)
			}

			testassert.ProtoEqual(t, tc.expectedBlobs, blobs)
		})
	}
}

func TestListAllBlobs(t *testing.T) {
	cfg, repo, _, client := setup(t)

	ctx, cancel := testhelper.Context()
	defer cancel()

	quarantine, err := quarantine.New(ctx, repo, config.NewLocator(cfg))
	require.NoError(t, err)

	quarantineRepoWithoutAlternates := proto.Clone(quarantine.QuarantinedRepo()).(*gitalypb.Repository)
	quarantineRepoWithoutAlternates.GitAlternateObjectDirectories = []string{}

	emptyRepo, _ := gittest.InitRepo(t, cfg, cfg.Storages[0])

	singleBlobRepo, singleBlobRepoPath := gittest.InitRepo(t, cfg, cfg.Storages[0])
	blobID := gittest.WriteBlob(t, cfg, singleBlobRepoPath, []byte("foobar"))

	for _, tc := range []struct {
		desc    string
		request *gitalypb.ListAllBlobsRequest
		verify  func(*testing.T, []*gitalypb.ListAllBlobsResponse_Blob)
	}{
		{
			desc: "empty repo",
			request: &gitalypb.ListAllBlobsRequest{
				Repository: emptyRepo,
			},
			verify: func(t *testing.T, blobs []*gitalypb.ListAllBlobsResponse_Blob) {
				require.Empty(t, blobs)
			},
		},
		{
			desc: "repo with single blob",
			request: &gitalypb.ListAllBlobsRequest{
				Repository: singleBlobRepo,
				BytesLimit: -1,
			},
			verify: func(t *testing.T, blobs []*gitalypb.ListAllBlobsResponse_Blob) {
				require.Equal(t, []*gitalypb.ListAllBlobsResponse_Blob{{
					Oid:  blobID.String(),
					Size: 6,
					Data: []byte("foobar"),
				}}, blobs)
			},
		},
		{
			desc: "repo with single blob and bytes limit",
			request: &gitalypb.ListAllBlobsRequest{
				Repository: singleBlobRepo,
				BytesLimit: 1,
			},
			verify: func(t *testing.T, blobs []*gitalypb.ListAllBlobsResponse_Blob) {
				require.Equal(t, []*gitalypb.ListAllBlobsResponse_Blob{{
					Oid:  blobID.String(),
					Size: 6,
					Data: []byte("f"),
				}}, blobs)
			},
		},
		{
			desc: "normal repo",
			request: &gitalypb.ListAllBlobsRequest{
				Repository: repo,
			},
			verify: func(t *testing.T, blobs []*gitalypb.ListAllBlobsResponse_Blob) {
				require.Greater(t, len(blobs), 300)
			},
		},
		{
			desc: "normal repo with limit",
			request: &gitalypb.ListAllBlobsRequest{
				Repository: repo,
				Limit:      2,
			},
			verify: func(t *testing.T, blobs []*gitalypb.ListAllBlobsResponse_Blob) {
				require.Len(t, blobs, 2)
			},
		},
		{
			desc: "normal repo with bytes limit",
			request: &gitalypb.ListAllBlobsRequest{
				Repository: repo,
				BytesLimit: 1,
			},
			verify: func(t *testing.T, blobs []*gitalypb.ListAllBlobsResponse_Blob) {
				require.Greater(t, len(blobs), 300)
				for _, blob := range blobs {
					emptyBlobID := "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391"
					if blob.Oid == emptyBlobID {
						require.Empty(t, blob.Data)
						require.Equal(t, int64(0), blob.Size)
					} else {
						require.Len(t, blob.Data, 1)
					}
				}
			},
		},
		{
			desc: "quarantine repo with alternates",
			request: &gitalypb.ListAllBlobsRequest{
				Repository: quarantine.QuarantinedRepo(),
			},
			verify: func(t *testing.T, blobs []*gitalypb.ListAllBlobsResponse_Blob) {
				require.Greater(t, len(blobs), 300)
			},
		},
		{
			desc: "quarantine repo without alternates",
			request: &gitalypb.ListAllBlobsRequest{
				Repository: quarantineRepoWithoutAlternates,
			},
			verify: func(t *testing.T, blobs []*gitalypb.ListAllBlobsResponse_Blob) {
				require.Empty(t, blobs)
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			stream, err := client.ListAllBlobs(ctx, tc.request)
			require.NoError(t, err)

			var blobs []*gitalypb.ListAllBlobsResponse_Blob
			for {
				resp, err := stream.Recv()
				if err != nil {
					if !errors.Is(err, io.EOF) {
						require.NoError(t, err)
					}
					break
				}

				blobs = append(blobs, resp.Blobs...)
			}

			tc.verify(t, blobs)
		})
	}
}
