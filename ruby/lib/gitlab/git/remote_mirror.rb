module Gitlab
  module Git
    class RemoteMirror
      attr_reader :repository, :remote_name, :ssh_auth, :only_branches_matching

      # An Array of local refnames that have diverged on the remote
      #
      # Only populated when `keep_divergent_refs` is enabled
      attr_reader :divergent_refs

      def initialize(repository, remote_name, ssh_auth:, only_branches_matching:, keep_divergent_refs:)
        @repository = repository
        @remote_name = remote_name
        @ssh_auth = ssh_auth
        @only_branches_matching = only_branches_matching
        @keep_divergent_refs = keep_divergent_refs

        @divergent_refs = []
      end

      def update
        ssh_auth.setup do |env|
          # Retrieve the remote branches first since they may take a while to load,
          # and the local branches may have changed during this time.
          remote_branch_list = remote_branches(env: env)
          updated_branches = changed_refs(local_branches, remote_branch_list)
          push_refs(default_branch_first(updated_branches.keys), env: env)
          delete_refs(local_branches, remote_branches(env: env), env: env)

          local_tags = refs_obj(repository.tags)
          remote_tags = refs_obj(repository.remote_tags(remote_name, env: env))

          updated_tags = changed_refs(local_tags, remote_tags)
          push_refs(updated_tags.keys, env: env)
          delete_refs(local_tags, remote_tags, env: env)
        end
      end

      private

      def ref_matchers
        @ref_matchers ||= only_branches_matching.map do |ref|
          GitLab::RefMatcher.new(ref)
        end
      end

      def local_branches
        @local_branches ||= refs_obj(
          repository.local_branches,
          match_refs: true
        )
      end

      def remote_branches(env:)
        @remote_branches ||= refs_obj(
          repository.remote_branches(remote_name, env: env),
          match_refs: true
        )
      end

      def refs_obj(refs, match_refs: false)
        refs.each_with_object({}) do |ref, refs|
          next if match_refs && !include_ref?(ref.name)

          key = ref.is_a?(Gitlab::Git::Tag) ? ref.refname : ref.name
          refs[key] = ref
        end
      end

      def changed_refs(local_refs, remote_refs)
        local_refs.select do |ref_name, ref|
          remote_ref = remote_refs[ref_name]

          # Ref doesn't exist on the remote, it should be created
          next true if remote_ref.nil?

          local_target = ref.dereferenced_target
          remote_target = remote_ref.dereferenced_target

          if local_target == remote_target
            # Ref is identical on the remote, no point mirroring
            false
          elsif @keep_divergent_refs
            # Mirror the ref if its remote counterpart hasn't diverged
            if repository.ancestor?(remote_target&.id, local_target&.id)
              true
            else
              Gitlab::GitLogger.info("Divergent ref `#{ref_name}` in #{repository.path} due to ancestry -- remote: #{remote_target&.id}, local: #{local_target&.id}")
              @divergent_refs << ref.refname
              false
            end
          else
            # Attempt to overwrite whatever's on the remote; push rules and
            # protected branches may still prevent this
            true
          end
        end
      end

      # Put the default branch first so it works fine when remote mirror is empty.
      def default_branch_first(branches)
        return unless branches.present?

        default_branch, branches = branches.partition do |branch|
          repository.root_ref == branch
        end

        branches.unshift(*default_branch)
      end

      def push_refs(refs, env:)
        return unless refs.present?

        repository.push_remote_branches(remote_name, refs, env: env)
      end

      def delete_refs(local_refs, remote_refs, env:)
        return if @keep_divergent_refs

        refs = refs_to_delete(local_refs, remote_refs)

        return unless refs.present?

        repository.delete_remote_branches(remote_name, refs.keys, env: env)
      end

      def refs_to_delete(local_refs, remote_refs)
        default_branch_id = repository.commit.id

        remote_refs.select do |remote_ref_name, remote_ref|
          next false if local_refs[remote_ref_name] # skip if branch or tag exist in local repo

          remote_ref_id = remote_ref.dereferenced_target.try(:id)

          repository.ancestor?(remote_ref_id, default_branch_id)
        end
      end

      def include_ref?(ref_name)
        return true unless ref_matchers.present?

        ref_matchers.any? { |matcher| matcher.matches?(ref_name) }
      end
    end
  end
end
