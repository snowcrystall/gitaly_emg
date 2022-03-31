module Gitlab
  module Git
    module RepositoryMirroring
      GITLAB_PROJECTS_TIMEOUT = 10800

      RemoteError = Class.new(StandardError)

      REFMAPS = {
        # With `:all_refs`, the repository is equivalent to the result of `git clone --mirror`
        all_refs: '+refs/*:refs/*',
        heads: '+refs/heads/*:refs/heads/*',
        tags: '+refs/tags/*:refs/tags/*'
      }.freeze

      def remote_branches(remote_name, env:)
        list_remote_refs(remote_name, env: env).map do |line|
          target, refname = line.strip.split("\t")

          if target.nil? || refname.nil?
            Rails.logger.info("Empty or invalid list of heads for remote: #{remote_name}")
            break []
          end

          next unless refname.start_with?('refs/heads/')

          target_commit = Gitlab::Git::Commit.find(self, target)
          Gitlab::Git::Branch.new(self, refname, target, target_commit)
        end.compact
      end

      def push_remote_branches(remote_name, branch_names, forced: true, env: {})
        success = @gitlab_projects.push_branches(remote_name, GITLAB_PROJECTS_TIMEOUT, forced, branch_names, env: env)

        success || gitlab_projects_error
      end

      def delete_remote_branches(remote_name, branch_names, env: {})
        success = @gitlab_projects.delete_remote_branches(remote_name, branch_names, env: env)

        success || gitlab_projects_error
      end

      def set_remote_as_mirror(remote_name, refmap: :all_refs)
        set_remote_refmap(remote_name, refmap)

        rugged.config["remote.#{remote_name}.mirror"] = true
        rugged.config["remote.#{remote_name}.prune"] = true
      end

      def remote_tags(remote, env: {})
        # Each line has this format: "dc872e9fa6963f8f03da6c8f6f264d0845d6b092\trefs/tags/v1.10.0\n"
        # We want to convert it to: [{ 'v1.10.0' => 'dc872e9fa6963f8f03da6c8f6f264d0845d6b092' }, ...]
        list_remote_refs(remote, env: env).map do |line|
          target, refname = line.strip.split("\t")

          # When the remote repo does not have tags.
          if target.nil? || refname.nil?
            Rails.logger.info "Empty or invalid list of tags for remote: #{remote}"
            break []
          end

          next unless refname.start_with?('refs/tags/')

          # We're only interested in tag references
          # See: http://stackoverflow.com/questions/15472107/when-listing-git-ls-remote-why-theres-after-the-tag-name
          next if refname.end_with?('^{}')

          target_commit = Gitlab::Git::Commit.find(self, target)
          Gitlab::Git::Tag.new(self,
                               name: refname,
                               target: target,
                               target_commit: target_commit)
        end.compact
      end

      private

      def set_remote_refmap(remote_name, refmap)
        Array(refmap).each_with_index do |refspec, i|
          refspec = REFMAPS[refspec] || refspec

          # We need multiple `fetch` entries, but Rugged only allows replacing a config, not adding to it.
          # To make sure we start from scratch, we set the first using rugged, and use `git` for any others
          if i == 0
            rugged.config["remote.#{remote_name}.fetch"] = refspec
          else
            run_git(%W[config --add remote.#{remote_name}.fetch #{refspec}])
          end
        end
      end

      def list_remote_refs(remote, env:)
        @list_remote_refs ||= {}
        @list_remote_refs[remote] ||= begin
          ref_list, exit_code, error = nil

          # List heads and tags, ignoring stuff like `refs/merge-requests` and `refs/pull`
          cmd = %W[#{Gitlab.config.git.bin_path} --git-dir=#{path} ls-remote --heads --tags #{remote}]

          Open3.popen3(env, *cmd) do |_stdin, stdout, stderr, wait_thr|
            ref_list  = stdout.read
            error     = stderr.read
            exit_code = wait_thr.value.exitstatus
          end

          raise RemoteError, error unless exit_code.zero?

          ref_list.each_line(chomp: true)
        end
      end
    end
  end
end
