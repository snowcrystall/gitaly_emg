module Gitlab
  #
  # In production, gitaly-ruby configuration is derived from environment
  # variables set by the Go gitaly parent process. During the rspec test
  # suite this parent process is not there so we need to get configuration
  # values from somewhere else. We used to work around this by setting
  # variables in ENV during the rspec boot but that turned out to be
  # fragile because Bundler.with_clean_env resets changes to ENV made
  # after Bundler was loaded. Instead of changing ENV, the TestSetup
  # module gives us a hacky way to set instance variables on the config
  # objects, bypassing the ENV lookups.
  #
  module TestSetup
    def test_global_ivar_override(name, value)
      instance_variable_set("@#{name}".to_sym, value)
    end
  end

  class Config
    class Git
      include TestSetup

      def bin_path
        @bin_path ||= ENV['GITALY_RUBY_GIT_BIN_PATH']
      end

      def hooks_directory
        @hooks_directory ||= ENV['GITALY_GIT_HOOKS_DIR']
      end

      def write_buffer_size
        @write_buffer_size ||= ENV['GITALY_RUBY_WRITE_BUFFER_SIZE'].to_i
      end

      def max_commit_or_tag_message_size
        @max_commit_or_tag_message_size ||= ENV['GITALY_RUBY_MAX_COMMIT_OR_TAG_MESSAGE_SIZE'].to_i
      end

      def rugged_git_config_search_path
        @rugged_git_config_search_path ||= ENV['GITALY_RUGGED_GIT_CONFIG_SEARCH_PATH']
      end
    end

    class Gitaly
      include TestSetup

      def bin_dir
        @bin_dir ||= ENV['GITALY_RUBY_GITALY_BIN_DIR']
      end

      def internal_socket
        @internal_socket ||= ENV['GITALY_SOCKET']
      end

      def rbtrace_enabled?
        @rbtrace_enabled ||= enabled?(ENV['GITALY_RUBY_RBTRACE_ENABLED'])
      end

      def objspace_trace_enabled?
        @objspace_trace_enabled ||= enabled?(ENV['GITALY_RUBY_OBJSPACE_TRACE_ENABLED'])
      end

      def enabled?(value)
        %w[true yes 1].include?(value&.downcase)
      end
    end

    class Logging
      def dir
        @dir ||= ENV['GITALY_LOG_DIR']
      end
    end

    def git
      @git ||= Git.new
    end

    def logging
      @logging ||= Logging.new
    end

    def gitaly
      @gitaly ||= Gitaly.new
    end
  end

  def self.config
    @config ||= Config.new
  end
end
