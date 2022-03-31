require 'spec_helper'

describe Gitlab::Git::RemoteRepository do
  include TestRepo

  let(:repository) { gitlab_git_from_gitaly(git_test_repo_read_only) }
  let(:non_existing_gitaly_repo) do
    Gitaly::Repository.new(storage_name: DEFAULT_STORAGE_NAME, relative_path: 'does-not-exist.git')
  end

  subject { described_class.new(repository) }

  describe '#empty?' do
    using RSpec::Parameterized::TableSyntax

    where(:repository, :result) do
      repository                                       | false
      gitlab_git_from_gitaly(non_existing_gitaly_repo) | true
    end

    with_them do
      it { expect(subject.empty?).to eq(result) }
    end
  end

  describe '#commit_id' do
    it 'returns an OID if the revision exists' do
      expect(subject.commit_id('v1.0.0')).to eq('6f6d7e7ed97bb5f0054f2b1df789b39ca89b6ff9')
    end

    it 'is nil when the revision does not exist' do
      expect(subject.commit_id('does-not-exist')).to be_nil
    end
  end

  describe '#branch_exists?' do
    using RSpec::Parameterized::TableSyntax

    where(:branch, :result) do
      'master'         | true
      'does-not-exist' | false
    end

    with_them do
      it { expect(subject.branch_exists?(branch)).to eq(result) }
    end
  end

  describe '#same_repository?' do
    using RSpec::Parameterized::TableSyntax

    where(:other_repository, :result) do
      repository                                               | true
      repository_from_relative_path(repository.relative_path)  | true
      repository_from_relative_path('wrong/relative-path.git') | false
    end

    with_them do
      it { expect(subject.same_repository?(other_repository)).to eq(result) }
    end
  end

  describe '#fetch_env' do
    let(:remote_repository) { described_class.new(repository) }

    let(:gitaly_client) { double(:gitaly_client) }
    let(:address) { 'fake-address' }
    let(:shared_secret) { 'fake-secret' }

    subject { remote_repository.fetch_env }

    before do
      allow(remote_repository).to receive(:gitaly_client).and_return(gitaly_client)

      expect(gitaly_client).to receive(:address).with(repository.storage).and_return(address)
      expect(gitaly_client).to receive(:shared_secret).with(repository.storage).and_return(shared_secret)
    end

    it { expect(subject).to be_a(Hash) }
    it { expect(subject['GITALY_ADDRESS']).to eq(address) }
    it { expect(subject['GITALY_TOKEN']).to eq(shared_secret) }
    it { expect(subject['GITALY_WD']).to eq(Dir.pwd) }

    it 'creates a plausible GIT_SSH_COMMAND' do
      git_ssh_command = subject['GIT_SSH_COMMAND']

      expect(git_ssh_command).to start_with('/')
      expect(git_ssh_command).to end_with('/gitaly-ssh upload-pack')
    end

    it 'creates a plausible GITALY_PAYLOAD' do
      req = Gitaly::SSHUploadPackRequest.decode_json(subject['GITALY_PAYLOAD'])

      expect(remote_repository.gitaly_repository).to eq(req.repository)
    end

    context 'when the token is blank' do
      let(:shared_secret) { '' }

      it { expect(subject.keys).not_to include('GITALY_TOKEN') }
    end
  end
end
