# Generated by the protocol buffer compiler.  DO NOT EDIT!
# Source: internal.proto for package 'gitaly'

require 'grpc'
require 'internal_pb'

module Gitaly
  module InternalGitaly
    # InternalGitaly is a gRPC service meant to be served by a Gitaly node, but
    # only reachable by Praefect or other Gitalies
    class Service

      include GRPC::GenericService

      self.marshal_class_method = :encode
      self.unmarshal_class_method = :decode
      self.service_name = 'gitaly.InternalGitaly'

      # WalkRepos walks the storage and streams back all known git repos on the
      # requested storage
      rpc :WalkRepos, Gitaly::WalkReposRequest, stream(Gitaly::WalkReposResponse)
    end

    Stub = Service.rpc_stub_class
  end
end