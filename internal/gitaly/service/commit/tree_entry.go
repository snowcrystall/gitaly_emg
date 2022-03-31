package commit

import (
	"fmt"
	"io"
	"strings"

	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/catfile"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"gitlab.com/gitlab-org/gitaly/v14/streamio"
)

func sendTreeEntry(stream gitalypb.CommitService_TreeEntryServer, c catfile.Batch, revision, path string, limit, maxSize int64) error {
	ctx := stream.Context()

	treeEntry, err := NewTreeEntryFinder(c).FindByRevisionAndPath(ctx, revision, path)
	if err != nil {
		return err
	}

	if treeEntry == nil || len(treeEntry.Oid) == 0 {
		return helper.ErrNotFoundf("not found: %s", path)
	}

	if treeEntry.Type == gitalypb.TreeEntry_COMMIT {
		response := &gitalypb.TreeEntryResponse{
			Type: gitalypb.TreeEntryResponse_COMMIT,
			Mode: treeEntry.Mode,
			Oid:  treeEntry.Oid,
		}
		if err := stream.Send(response); err != nil {
			return helper.ErrUnavailablef("TreeEntry: send: %v", err)
		}

		return nil
	}

	if treeEntry.Type == gitalypb.TreeEntry_TREE {
		treeInfo, err := c.Info(ctx, git.Revision(treeEntry.Oid))
		if err != nil {
			return err
		}

		response := &gitalypb.TreeEntryResponse{
			Type: gitalypb.TreeEntryResponse_TREE,
			Oid:  treeEntry.Oid,
			Size: treeInfo.Size,
			Mode: treeEntry.Mode,
		}
		return helper.ErrUnavailable(stream.Send(response))
	}

	objectInfo, err := c.Info(ctx, git.Revision(treeEntry.Oid))
	if err != nil {
		return helper.ErrInternalf("TreeEntry: %v", err)
	}

	if strings.ToLower(treeEntry.Type.String()) != objectInfo.Type {
		return helper.ErrInternalf(
			"TreeEntry: mismatched object type: tree-oid=%s object-oid=%s entry-type=%s object-type=%s",
			treeEntry.Oid, objectInfo.Oid, treeEntry.Type.String(), objectInfo.Type,
		)
	}

	dataLength := objectInfo.Size

	if maxSize > 0 && dataLength > maxSize {
		return helper.ErrFailedPreconditionf(
			"TreeEntry: object size (%d) is bigger than the maximum allowed size (%d)",
			dataLength, maxSize,
		)
	}

	if limit > 0 && dataLength > limit {
		dataLength = limit
	}

	response := &gitalypb.TreeEntryResponse{
		Type: gitalypb.TreeEntryResponse_BLOB,
		Oid:  objectInfo.Oid.String(),
		Size: objectInfo.Size,
		Mode: treeEntry.Mode,
	}
	if dataLength == 0 {
		return helper.ErrUnavailable(stream.Send(response))
	}

	blobObj, err := c.Blob(ctx, git.Revision(objectInfo.Oid))
	if err != nil {
		return err
	}

	sw := streamio.NewWriter(func(p []byte) error {
		response.Data = p

		if err := stream.Send(response); err != nil {
			return helper.ErrUnavailablef("TreeEntry: send: %v", err)
		}

		// Use a new response so we don't send other fields (Size, ...) over and over
		response = &gitalypb.TreeEntryResponse{}

		return nil
	})

	_, err = io.CopyN(sw, blobObj.Reader, dataLength)
	return err
}

func (s *server) TreeEntry(in *gitalypb.TreeEntryRequest, stream gitalypb.CommitService_TreeEntryServer) error {
	if err := validateRequest(in); err != nil {
		return helper.ErrInvalidArgumentf("TreeEntry: %v", err)
	}

	repo := s.localrepo(in.GetRepository())

	requestPath := string(in.GetPath())
	// filepath.Dir("api/docs") => "api" Correct!
	// filepath.Dir("api/docs/") => "api/docs" WRONG!
	if len(requestPath) > 1 {
		requestPath = strings.TrimRight(requestPath, "/")
	}

	c, err := s.catfileCache.BatchProcess(stream.Context(), repo)
	if err != nil {
		return err
	}

	return sendTreeEntry(stream, c, string(in.GetRevision()), requestPath, in.GetLimit(), in.GetMaxSize())
}

func validateRequest(in *gitalypb.TreeEntryRequest) error {
	if err := git.ValidateRevision(in.Revision); err != nil {
		return err
	}

	if len(in.GetPath()) == 0 {
		return fmt.Errorf("empty Path")
	}

	if in.GetLimit() < 0 {
		return fmt.Errorf("negative Limit")
	}
	if in.GetMaxSize() < 0 {
		return fmt.Errorf("negative MaxSize")
	}

	return nil
}
