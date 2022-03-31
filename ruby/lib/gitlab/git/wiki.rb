module Gitlab
  module Git
    class Wiki
      DuplicatePageError = Class.new(StandardError)
      OperationError = Class.new(StandardError)
      PageNotFound = Class.new(GRPC::NotFound)

      CommitDetails = Struct.new(:user_id, :username, :name, :email, :message) do
        def to_h
          { user_id: user_id, username: username, name: name, email: email, message: message }
        end
      end
      PageBlob = Struct.new(:name)

      attr_reader :repository

      def self.default_ref
        'master'
      end

      # Initialize with a Gitlab::Git::Repository instance
      def initialize(repository)
        @repository = repository
      end

      def repository_exists?
        @repository.exists?
      end

      def write_page(name, format, content, commit_details)
        gollum_write_page(name, format, content, commit_details)
      end

      def update_page(page_path, title, format, content, commit_details)
        gollum_update_page(page_path, title, format, content, commit_details)
      end

      def pages(limit: nil, sort: nil, direction_desc: false)
        gollum_get_all_pages(limit: limit, sort: sort, direction_desc: direction_desc)
      end

      def page(title:, version: nil, dir: nil)
        gollum_find_page(title: title, version: version, dir: dir)
      end

      def count_page_versions(page_path)
        @repository.count_commits(ref: 'HEAD', path: page_path)
      end

      def preview_slug(title, format)
        # Adapted from gollum gem (Gollum::Wiki#preview_page) to avoid
        # using Rugged through a Gollum::Wiki instance
        page_class = Gollum::Page
        page = page_class.new(nil)
        ext = page_class.format_to_ext(format.to_sym)
        name = page_class.cname(title) + '.' + ext
        blob = PageBlob.new(name)
        page.populate(blob)
        page.url_path
      end

      def gollum_wiki
        options = {}
        options[:ref] = gollum_default_ref if gollum_default_ref

        @gollum_wiki ||= Gollum::Wiki.new(@repository.path, options)
      end

      private

      def gollum_default_ref
        @gollum_default_ref ||= @repository.root_ref || @repository.head_symbolic_ref
      end

      def new_page(gollum_page)
        Gitlab::Git::WikiPage.new(gollum_page, new_version(gollum_page, gollum_page.version.id))
      end

      def new_version(gollum_page, commit_id)
        Gitlab::Git::WikiPageVersion.new(version(commit_id), gollum_page&.format)
      end

      def version(commit_id)
        Gitlab::Git::Commit.find(@repository, commit_id)
      end

      def assert_type!(object, klass)
        raise ArgumentError, "expected a #{klass}, got #{object.inspect}" unless object.is_a?(klass)
      end

      def committer_with_hooks(commit_details)
        Gitlab::Git::CommitterWithHooks.new(self, commit_details.to_h)
      end

      def with_committer_with_hooks(commit_details)
        committer = committer_with_hooks(commit_details)

        yield committer

        committer.commit

        nil
      end

      # options:
      #  :page     - The Integer page number.
      #  :per_page - The number of items per page.
      #  :limit    - Total number of items to return.
      def commits_from_page(gollum_page, options = {})
        unless options[:limit]
          options[:offset] = ([1, options.delete(:page).to_i].max - 1) * Gollum::Page.per_page
          options[:limit] = (options.delete(:per_page) || Gollum::Page.per_page).to_i
        end

        @repository.log(ref: gollum_page.last_version.id,
                        path: gollum_page.path,
                        limit: options[:limit],
                        offset: options[:offset])
      end

      # Retrieve the page at that `page_path`, raising an error if it does not exist
      def gollum_page_by_path(page_path)
        page_name = Gollum::Page.canonicalize_filename(page_path)
        page_dir = File.split(page_path).first

        gollum_wiki.paged(page_name, page_dir) || (raise PageNotFound, page_path)
      end

      def gollum_write_page(name, format, content, commit_details)
        assert_type!(format, Symbol)
        assert_type!(commit_details, CommitDetails)

        with_committer_with_hooks(commit_details) do |committer|
          filename = File.basename(name)
          dir = (tmp_dir = File.dirname(name)) == '.' ? '' : tmp_dir

          gollum_wiki.write_page(filename, format, content, { committer: committer }, dir)
        end
      rescue Gollum::DuplicatePageError => e
        raise Gitlab::Git::Wiki::DuplicatePageError, e.message
      end

      def gollum_update_page(page_path, title, format, content, commit_details)
        assert_type!(format, Symbol)
        assert_type!(commit_details, CommitDetails)

        with_committer_with_hooks(commit_details) do |committer|
          page = gollum_page_by_path(page_path)
          # Instead of performing two renames if the title has changed,
          # the update_page will only update the format and content and
          # the rename_page will do anything related to moving/renaming
          gollum_wiki.update_page(page, page.name, format, content, committer: committer)
          gollum_wiki.rename_page(page, title, committer: committer)
        end
      end

      def gollum_find_page(title:, version: nil, dir: nil)
        if version
          version = Gitlab::Git::Commit.find(@repository, version)&.id
          return unless version
        end

        gollum_page = gollum_wiki.page(title, version, dir)
        return unless gollum_page

        new_page(gollum_page)
      end

      def gollum_get_all_pages(limit: nil, sort: nil, direction_desc: false)
        gollum_wiki.pages(
          limit: limit, sort: sort, direction_desc: direction_desc
        ).map do |gollum_page|
          new_page(gollum_page)
        end
      end
    end
  end
end
