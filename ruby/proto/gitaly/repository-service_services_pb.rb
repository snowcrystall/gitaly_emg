# Generated by the protocol buffer compiler.  DO NOT EDIT!
# Source: repository-service.proto for package 'gitaly'

require 'grpc'
require 'repository-service_pb'

module Gitaly
  module RepositoryService
    class Service

      include GRPC::GenericService

      self.marshal_class_method = :encode
      self.unmarshal_class_method = :decode
      self.service_name = 'gitaly.RepositoryService'

      rpc :RepositoryExists, Gitaly::RepositoryExistsRequest, Gitaly::RepositoryExistsResponse
      rpc :RepackIncremental, Gitaly::RepackIncrementalRequest, Gitaly::RepackIncrementalResponse
      rpc :RepackFull, Gitaly::RepackFullRequest, Gitaly::RepackFullResponse
      rpc :MidxRepack, Gitaly::MidxRepackRequest, Gitaly::MidxRepackResponse
      rpc :GarbageCollect, Gitaly::GarbageCollectRequest, Gitaly::GarbageCollectResponse
      rpc :WriteCommitGraph, Gitaly::WriteCommitGraphRequest, Gitaly::WriteCommitGraphResponse
      rpc :RepositorySize, Gitaly::RepositorySizeRequest, Gitaly::RepositorySizeResponse
      rpc :ApplyGitattributes, Gitaly::ApplyGitattributesRequest, Gitaly::ApplyGitattributesResponse
      # FetchRemote fetches references from a remote repository into the local
      # repository.
      rpc :FetchRemote, Gitaly::FetchRemoteRequest, Gitaly::FetchRemoteResponse
      rpc :CreateRepository, Gitaly::CreateRepositoryRequest, Gitaly::CreateRepositoryResponse
      rpc :GetArchive, Gitaly::GetArchiveRequest, stream(Gitaly::GetArchiveResponse)
      rpc :HasLocalBranches, Gitaly::HasLocalBranchesRequest, Gitaly::HasLocalBranchesResponse
      # FetchSourceBranch fetches a branch from a second (potentially remote)
      # repository into the given repository.
      rpc :FetchSourceBranch, Gitaly::FetchSourceBranchRequest, Gitaly::FetchSourceBranchResponse
      rpc :Fsck, Gitaly::FsckRequest, Gitaly::FsckResponse
      rpc :WriteRef, Gitaly::WriteRefRequest, Gitaly::WriteRefResponse
      rpc :FindMergeBase, Gitaly::FindMergeBaseRequest, Gitaly::FindMergeBaseResponse
      rpc :CreateFork, Gitaly::CreateForkRequest, Gitaly::CreateForkResponse
      rpc :IsSquashInProgress, Gitaly::IsSquashInProgressRequest, Gitaly::IsSquashInProgressResponse
      rpc :CreateRepositoryFromURL, Gitaly::CreateRepositoryFromURLRequest, Gitaly::CreateRepositoryFromURLResponse
      # CreateBundle creates a bundle from all refs
      rpc :CreateBundle, Gitaly::CreateBundleRequest, stream(Gitaly::CreateBundleResponse)
      # CreateBundleFromRefList creates a bundle from a stream of ref patterns
      rpc :CreateBundleFromRefList, stream(Gitaly::CreateBundleFromRefListRequest), stream(Gitaly::CreateBundleFromRefListResponse)
      rpc :CreateRepositoryFromBundle, stream(Gitaly::CreateRepositoryFromBundleRequest), Gitaly::CreateRepositoryFromBundleResponse
      # GetConfig reads the target repository's gitconfig and streams its contents
      # back. Returns a NotFound error in case no gitconfig was found.
      rpc :GetConfig, Gitaly::GetConfigRequest, stream(Gitaly::GetConfigResponse)
      rpc :SetConfig, Gitaly::SetConfigRequest, Gitaly::SetConfigResponse
      rpc :DeleteConfig, Gitaly::DeleteConfigRequest, Gitaly::DeleteConfigResponse
      rpc :FindLicense, Gitaly::FindLicenseRequest, Gitaly::FindLicenseResponse
      rpc :GetInfoAttributes, Gitaly::GetInfoAttributesRequest, stream(Gitaly::GetInfoAttributesResponse)
      rpc :CalculateChecksum, Gitaly::CalculateChecksumRequest, Gitaly::CalculateChecksumResponse
      rpc :Cleanup, Gitaly::CleanupRequest, Gitaly::CleanupResponse
      rpc :GetSnapshot, Gitaly::GetSnapshotRequest, stream(Gitaly::GetSnapshotResponse)
      rpc :CreateRepositoryFromSnapshot, Gitaly::CreateRepositoryFromSnapshotRequest, Gitaly::CreateRepositoryFromSnapshotResponse
      rpc :GetRawChanges, Gitaly::GetRawChangesRequest, stream(Gitaly::GetRawChangesResponse)
      rpc :SearchFilesByContent, Gitaly::SearchFilesByContentRequest, stream(Gitaly::SearchFilesByContentResponse)
      rpc :SearchFilesByName, Gitaly::SearchFilesByNameRequest, stream(Gitaly::SearchFilesByNameResponse)
      rpc :RestoreCustomHooks, stream(Gitaly::RestoreCustomHooksRequest), Gitaly::RestoreCustomHooksResponse
      rpc :BackupCustomHooks, Gitaly::BackupCustomHooksRequest, stream(Gitaly::BackupCustomHooksResponse)
      rpc :GetObjectDirectorySize, Gitaly::GetObjectDirectorySizeRequest, Gitaly::GetObjectDirectorySizeResponse
      rpc :CloneFromPool, Gitaly::CloneFromPoolRequest, Gitaly::CloneFromPoolResponse
      rpc :CloneFromPoolInternal, Gitaly::CloneFromPoolInternalRequest, Gitaly::CloneFromPoolInternalResponse
      # RemoveRepository will move the repository to `+gitaly/tmp/<relative_path>_removed` and
      # eventually remove it. This ensures that even on networked filesystems the
      # data is actually removed even if there's someone still handling the data.
      rpc :RemoveRepository, Gitaly::RemoveRepositoryRequest, Gitaly::RemoveRepositoryResponse
      rpc :RenameRepository, Gitaly::RenameRepositoryRequest, Gitaly::RenameRepositoryResponse
      rpc :ReplicateRepository, Gitaly::ReplicateRepositoryRequest, Gitaly::ReplicateRepositoryResponse
      rpc :OptimizeRepository, Gitaly::OptimizeRepositoryRequest, Gitaly::OptimizeRepositoryResponse
      # SetFullPath writes the "gitlab.fullpath" configuration into the
      # repository's gitconfig. This is mainly to help debugging purposes in case
      # an admin inspects the repository's gitconfig such that he can easily see
      # what the repository name is.
      rpc :SetFullPath, Gitaly::SetFullPathRequest, Gitaly::SetFullPathResponse
    end

    Stub = Service.rpc_stub_class
  end
end