# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: internal.proto

require 'google/protobuf'

require 'lint_pb'
Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("internal.proto", :syntax => :proto3) do
    add_message "gitaly.WalkReposRequest" do
      optional :storage_name, :string, 1
    end
    add_message "gitaly.WalkReposResponse" do
      optional :relative_path, :string, 1
    end
  end
end

module Gitaly
  WalkReposRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitaly.WalkReposRequest").msgclass
  WalkReposResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitaly.WalkReposResponse").msgclass
end