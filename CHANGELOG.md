# Gitaly changelog

## 14.2.3 (2021-09-01)

### Performance (1 change)

- [Only activate Git pack-objects hook if cache is enabled](gitlab-org/gitaly@11e725e16c2f6c1fc93ceb3cc250ec25c6276f00) ([merge request](gitlab-org/gitaly!3815))

## 14.2.2 (2021-08-31)

### Fixed (1 change)

- [Downgrade grpc from 1.38.0 to 1.30.2](gitlab-org/security/gitaly@d8ccb23140704d941daca6a69ee3953997ac2cac)

## 14.2.1 (2021-08-23)

No changes.

## 14.2.0 (2021-08-20)

### Added (8 changes)

- [Add ListRef RPC to list all refs](gitlab-org/gitaly@cb8c7392b544ece7c27fdd2467b2afcab7c0b400) ([merge request](gitlab-org/gitaly!3731))
- [gitaly: Sort the tags returned by FindAllTags RPC](gitlab-org/gitaly@1a79647c0140fb94387859b8b042cedc312af7da) ([merge request](gitlab-org/gitaly!3707))
- [gitaly: Sort the tags returned by FindAllTags RPC](gitlab-org/gitaly@241290c104d58ebb88d01bc9704fd94d4e705ab2) ([merge request](gitlab-org/gitaly!3707))
- [Add CreateBundleFromRefList RPC](gitlab-org/gitaly@dae23dd0dbe37d4e1997fb9b1dfbea554be25f64) ([merge request](gitlab-org/gitaly!3594))
- [repository: Implement new `SetFullPath()` RPC](gitlab-org/gitaly@e7fa825668f2c57885d29f8ea77bd8a92102c0e0) ([merge request](gitlab-org/gitaly!3706))
- [backup: support of the Cloud storages](gitlab-org/gitaly@a9801375db1253f7413449e9dbc9a12b7325cf31) ([merge request](gitlab-org/gitaly!3687))
- [Add pagination to GetTreeEntries RPC](gitlab-org/gitaly@b458bbc55c56011d9e7b1df1092b9c008cea575f) ([merge request](gitlab-org/gitaly!3611))
- [Add sort capability for GetTreeEntries](gitlab-org/gitaly@09739e1096c078646fb5f1d2982b4d662d902907) ([merge request](gitlab-org/gitaly!3671))

### Fixed (23 changes)

- [conflicts: Always resolve conflicts with hooks](gitlab-org/gitaly@2dbee32be724c00d376c6ef44da429af1c7cd12f) ([merge request](gitlab-org/gitaly!3743))
- [Accept array of patterns in CreateBundleFromRefList](gitlab-org/gitaly@87b298809f40895437cdcdeec965b3fcd643a443) ([merge request](gitlab-org/gitaly!3742))
- [Correct the first page when paginating CommitService.CommitsBetween](gitlab-org/gitaly@be8444d9cac37544c52581faca3e2f7eadf25520) ([merge request](gitlab-org/gitaly!3737))
- [Set a timeout for Praefect's SQL metrics](gitlab-org/gitaly@ea7a08a8fc1c65fd0e360106c69327b6362207a8) ([merge request](gitlab-org/gitaly!3730))
- [operations: Fix error when rebase in UserSquash fails](gitlab-org/gitaly@682c59143e598c041b02a0100beac16c23bb5e4a) ([merge request](gitlab-org/gitaly!3727))
- [Fix create backup parallelism](gitlab-org/gitaly@ef7e85412b244326ddae77e5ccea1d37bf300eea) ([merge request](gitlab-org/gitaly!3690))
- [ref: Fix FindAllTags' pipeline code not returning nested tags](gitlab-org/gitaly@7344791e72998f36eef9d754a53c9ebdef0f4438) ([merge request](gitlab-org/gitaly!3728))
- [featureflag: Unconditionally enable transactional removal of repos](gitlab-org/gitaly@bd8740dae59716e5ecb1ae934243562901a01fa9) ([merge request](gitlab-org/gitaly!3725))
- [Increment generations of up to date storages only](gitlab-org/gitaly@dcd7dcfadd91de1039b84ee997c1daa4227b5ad5) ([merge request](gitlab-org/gitaly!3689))
- [repository: Remove feature flag for direct fetches](gitlab-org/gitaly@6490360d504c39bf6ebe9bf5414ba91ed4b7f6c5) ([merge request](gitlab-org/gitaly!3716))
- [repository: Fix SetFullPath writing multiple entries](gitlab-org/gitaly@0622b256ea3bdef6a0a52cf66ce9d3ecca08b338) ([merge request](gitlab-org/gitaly!3712))
- [repository: Allow voting on missing config](gitlab-org/gitaly@0663898ba5ba01894c7cbbdc38609edc02ed320b) ([merge request](gitlab-org/gitaly!3706))
- [Ignore symbolic refs in UpdateRemoteMirror](gitlab-org/gitaly@cb29419829d2af90cf8b24eac90d0ec18ec4b5d8) ([merge request](gitlab-org/gitaly!3699))
- [Fix UpdateRemoteMirror force pushing over diverged refs due to a race](gitlab-org/gitaly@dd160399735126da6a1757bf9ff3818059848c35) ([merge request](gitlab-org/gitaly!3699))
- [Fix mirroring empty repository in UpdateRemoteMirror](gitlab-org/gitaly@55488e5c3ac398359edf64aec613b6a31ca623c8) ([merge request](gitlab-org/gitaly!3699))
- [Fix mirroring branch and tag of same name in UpdateRemoteMirror](gitlab-org/gitaly@bdb785c8d2d7361d31a3f3649ca25cfae56b153c) ([merge request](gitlab-org/gitaly!3699))
- [Fix mirroring branch called 'tag' in UpdateRemoteMirror](gitlab-org/gitaly@6e2d217dd7c6720c22a65abe9a18b39216801fae) ([merge request](gitlab-org/gitaly!3699))
- [remote: Remove hardcoded branch name from error](gitlab-org/gitaly@ba882c2c9320ba354cd880983eb2f011b3086b4e) ([merge request](gitlab-org/gitaly!3701))
- [quarantine: Fix calling RPCs with manual object quarantine environments](gitlab-org/gitaly@e491f9266d25c43ac33f324cf1bc662254668f93) ([merge request](gitlab-org/gitaly!3697))
- [Remove orphan repository replicas](gitlab-org/gitaly@943fdd846af8ccb755b104dceb65bd8926f960b4) ([merge request](gitlab-org/gitaly!3674))
- [ref: Fix FindAllTags returning tag names instead of ref names](gitlab-org/gitaly@3ad398c3ef00869989f1eea38bc676dab9bf9f56) ([merge request](gitlab-org/gitaly!3694))
- [operations: Squash via patches if commits are not directly related](gitlab-org/gitaly@0cbd9335b99d97cdafd13a64bfc41476541713a4) ([merge request](gitlab-org/gitaly!3685))
- [Set restrictive permissions on backup directories](gitlab-org/gitaly@49cfca5e02964d48a1e3750f0bdce140bfeb7966) ([merge request](gitlab-org/gitaly!3677))

### Changed (12 changes)

- [Write ref list when taking backups](gitlab-org/gitaly@ed2445a6975f0b300b287e710f4b2710cea7068f) ([merge request](gitlab-org/gitaly!3747))
- [Share default prometheus configuration between Gitaly and Praefect](gitlab-org/gitaly@9d66dfdf91c77a6e75eeb49c8b78e6a0e9443abb) ([merge request](gitlab-org/gitaly!3730))
- [operations: Make UserMergeBranch return PreconditionFailed on conflicts](gitlab-org/gitaly@766e37dc7cc02a61225ee65bbd2ebe16c91e9663) ([merge request](gitlab-org/gitaly!3700))
- [featureflag: Enable CreateRepositoryFromBundleAtomicFetch by default](gitlab-org/gitaly@87b3e2aa402516ae92a976be7ea25e113eb2110b) ([merge request](gitlab-org/gitaly!3696))
- [operations: Fix UserCreateTag not respecting timezone](gitlab-org/gitaly@b9f35d555d6866380c3d433dc9eab21052b9e711) ([merge request](gitlab-org/gitaly!3694))
- [Bump gitlab-shell dependency](gitlab-org/gitaly@eaa96cecf20a08b3e95131b73b98be01acf14d8f) ([merge request](gitlab-org/gitaly!3695))
- [blob: Remove LFSPointersPipeline feature flag](gitlab-org/gitaly@e2dd22c4d58593f7e02febfb3113e18b6b07434d) ([merge request](gitlab-org/gitaly!3692))
- [Enable ResolveConflictsWithHooks feature flag by default](gitlab-org/gitaly@2211513fecbb62b0e42f8dffffcffa8618e3a443) ([merge request](gitlab-org/gitaly!3639))
- [Improve tracing instrumentation of catfile.Batch](gitlab-org/gitaly@a93d5db5daa89cf060eb1ec41e56c0cad88d406d) ([merge request](gitlab-org/gitaly!3688))
- [Bump grpc and grpc-tools to 1.38.0](gitlab-org/gitaly@5c2d0d08b9a788ab96e34060509e8f6ef51223df) ([merge request](gitlab-org/gitaly!3659))
- [operations: Use object quarantine directories in UserCreateTags](gitlab-org/gitaly@eca4b3c241ef5596501fc0a8eb7697c98c39b580) ([merge request](gitlab-org/gitaly!3655))
- [HEAD: default to main for new repositories](gitlab-org/gitaly@553513b4d65d4f1203ef1b0b3b71a0d5287df1c6) ([merge request](gitlab-org/gitaly!3537))

### Removed (3 changes)

- [proto: Remove the IsRebaseInProgress RPC](gitlab-org/gitaly@8ef37f161d6ddf5344e46e24bf4e00644e07b659) ([merge request](gitlab-org/gitaly!3724))
- [Remove the repository importer](gitlab-org/gitaly@dcc71730c2422d5a949979e2517a3207460d5582) ([merge request](gitlab-org/gitaly!3721))
- [Remove UpdateRemoteMirror Go port's feature flag](gitlab-org/gitaly@8a5e40ddae82af307f68b73aefa4c141316a6f48) ([merge request](gitlab-org/gitaly!3693))

### Security (2 changes)

- [Updates rdoc version](gitlab-org/gitaly@2470ed7a73b354e20709699d4bdef57f2787a8a0) ([merge request](gitlab-org/gitaly!3760))
- [catfile: Allow parsing of long git commit headers](gitlab-org/gitaly@b3e664dc0094a7fa25e8547a223607b9045e0030)

### Performance (14 changes)

- [operations: Support object quarantine in UserRebaseConfirmable](gitlab-org/gitaly@02976eb66065800455d12bd041d211f8a4f9059d) ([merge request](gitlab-org/gitaly!3753))
- [ref: Skip loading of peeled non-tag objects in FindAllTags](gitlab-org/gitaly@bf6e98c3646711d22a573df4a8eb85d8982a344e) ([merge request](gitlab-org/gitaly!3750))
- [git: Do not write commit graphs on fetches](gitlab-org/gitaly@863e5b0e57fb7ad1cb0b73abbcd411a2d1bbf59b) ([merge request](gitlab-org/gitaly!3748))
- [Enable Go port of UserApplyPatch](gitlab-org/gitaly@9d3f0698aabc2f8941e6fb91f90ffc413a2944dd) ([merge request](gitlab-org/gitaly!3745))
- [operations: Always use quarantine in UserCreateTag](gitlab-org/gitaly@3d6d75846bd5d3f0e03eacb06bbf461b01b5fc97) ([merge request](gitlab-org/gitaly!3744))
- [conflicts: Always use object quarantine for ResolveConflicts](gitlab-org/gitaly@7ddff642768fe71c9c9c21d144260ba7a1c8995b) ([merge request](gitlab-org/gitaly!3743))
- [localrepo: Always disable use of alternate refs in internal fetches](gitlab-org/gitaly@d5db37ec4eac68a2e635c80ef01b274fddc61ddf) ([merge request](gitlab-org/gitaly!3740))
- [git: Speed up fetches in repos with many refs](gitlab-org/gitaly@caf2cfabc79b21e8b4c5e0245eb4fbf5e7f7a493) ([merge request](gitlab-org/gitaly!3739))
- [localrepo: Speed up connectivity check when fetching into pooled repos](gitlab-org/gitaly@ddfe6e8191d8b24f0bcdd2da1ea1de071704cc3e) ([merge request](gitlab-org/gitaly!3720))
- [operations: Support worktreeless squashing for non-fast-forward merges](gitlab-org/gitaly@007b43e561f7851cbaa96dfa75a21eabd2508539) ([merge request](gitlab-org/gitaly!3714))
- [repository: Remove feature flag for atomic repo creation from bundles](gitlab-org/gitaly@c425e05bfeac31763ffecfbbacd93798f39a7186) ([merge request](gitlab-org/gitaly!3717))
- [conflitcs: Implement object quarantine for ResolveConflicts](gitlab-org/gitaly@7c1143da4de8fd0c4e120b9a2c7d24381f8d400e) ([merge request](gitlab-org/gitaly!3680))
- [operations: Implement squashing without worktrees](gitlab-org/gitaly@640b4de0970bd302cf178393755375b410985c7b) ([merge request](gitlab-org/gitaly!3657))
- [coordinator: Only schedule replication for differing error states](gitlab-org/gitaly@73839029f79d4ebdbc8d96475cf9bd0e2a599b2b) ([merge request](gitlab-org/gitaly!3660))

### Other (7 changes)

- [PostUploadPack tests: conform to test name style guide](gitlab-org/gitaly@9351b4ac62e01a81ddd892c617e3033ebde523b1) ([merge request](gitlab-org/gitaly!3704))
- [Move pack.stat logging to PackObjectsHook](gitlab-org/gitaly@84a84e95ab791803b9ea3f57bb48982f4359abcc) ([merge request](gitlab-org/gitaly!3749))
- [gitaly-git2go: cleanup backwards compatibility code](gitlab-org/gitaly@335e3d74b2a67d79f7dcb9bc668ca84178686de6) ([merge request](gitlab-org/gitaly!3683))
- [Improve streamrpc error messages](gitlab-org/gitaly@29bd70ca2e4d334645ccc7675bb5d1bc6aa7564a) ([merge request](gitlab-org/gitaly!3673))
- [Use service.Dependencies to inject pack-objects cache](gitlab-org/gitaly@5cc7430e77ec6c2d27b7c1ac60e7edc5dea5b9a6) ([merge request](gitlab-org/gitaly!3672))
- [Return close error in streamrpc.Call](gitlab-org/gitaly@ed0a0115476b1f0ec4748d04be7adfdbd6fb58ed) ([merge request](gitlab-org/gitaly!3670))
- [Add documentation for pack-object cache](gitlab-org/gitaly@4d6d21999ed4e70c84f8e754adacb4f6204ca0b0) ([merge request](gitlab-org/gitaly!3636))

### bug (1 change)

- [operations: Always enable squashing without worktrees](gitlab-org/gitaly@6e5bf3a4d062e1eb5d8e15784139f05a754a4a53) ([merge request](gitlab-org/gitaly!3766))

### feature (1 change)

- [blob: Implement new ListAllBlobs RPC](gitlab-org/gitaly@598fb165eff266513b9788ebcdc6aa900a65f083) ([merge request](gitlab-org/gitaly!3703))

### removal (1 change)

- [Remove `praefect reconcile` subcommand](gitlab-org/gitaly@0c48ca9e3f368f03d0d78de0a31540d2fe248242) ([merge request](gitlab-org/gitaly!3667))

## 14.1.3 (2021-08-17)

No changes.

## 14.1.2 (2021-08-03)

### Security (1 change)

- [catfile: Allow parsing of long git commit headers](gitlab-org/security/gitaly@c870566d92cdd9158272d3d0d2cedcbe587095c0) ([merge request](gitlab-org/security/gitaly!40))

## 14.1.1 (2021-07-28)

No changes.

## 14.1.0 (2021-07-21)

### Added (8 changes)

- [repository: Support transactional voting in `RemoveRepository()`](gitlab-org/gitaly@74895a79a782887d9f113d1748328e19d3d885ee) ([merge request](gitlab-org/gitaly!3603))
- [commit: Allow listing commits in reverse](gitlab-org/gitaly@e25c30d8b1a4e453209c08cae97871da300566a4) ([merge request](gitlab-org/gitaly!3650))
- [commit: Implement new `ListAllCommits()` RPC](gitlab-org/gitaly@c3dfb4ecb6b4f439ea2ed7685e9e46a7a31ae955) ([merge request](gitlab-org/gitaly!3628))
- [commit: Implement new `ListCommits()` RPC](gitlab-org/gitaly@a14bbcd6d026ba09d9277e49c7b380d95fe2fbbe) ([merge request](gitlab-org/gitaly!3628))
- [blackbox: Implement support for "push" probes](gitlab-org/gitaly@67783f933b9d9a9d5c7788311c5526f61e1d3a36) ([merge request](gitlab-org/gitaly!3613))
- [Retrieve ref from unborn HEAD for wikis](gitlab-org/gitaly@a1395884165c14f0dbaa21bb59209b026d60edf6) ([merge request](gitlab-org/gitaly!3610))
- [Log details about merge commit errors](gitlab-org/gitaly@9ccd5948b94b1ea5fbbe077e26aa9094d3da5bf1) ([merge request](gitlab-org/gitaly!3618))
- [blob: Implement new RPC to list blobs via revisions](gitlab-org/gitaly@ab1c4c6b67de7800d04c73beac02f916423b3e4d) ([merge request](gitlab-org/gitaly!3598))

### Fixed (10 changes)

- [coordinator: Create replication jobs if the primary cast a vote](gitlab-org/gitaly@cbf4ed8ffcb9e2600bdb20da420bc694c8523d8f) ([merge request](gitlab-org/gitaly!3642))
- [coordinator: Fix repo creation/removal race for up-to-date secondaries](gitlab-org/gitaly@a22bb54cce526c870600f6e2d77f57aa951ffe7c) ([merge request](gitlab-org/gitaly!3603))
- [replicator: Replicate `all_refs` parameter for PackRefs RPC](gitlab-org/gitaly@39fedeab73b9e1d9af25dfc867725388c5387c72) ([merge request](gitlab-org/gitaly!3638))
- [replicator: Replicate `prune` parameter for GarbageCollect RPC](gitlab-org/gitaly@159d50fbf912945f9bce7e8c0d5a8f9847bc6c49) ([merge request](gitlab-org/gitaly!3638))
- [fix: Set operations timezone by User](gitlab-org/gitaly@78c38deaa1c6d43c02c3d5ada735fcfd1296f8db) ([merge request](gitlab-org/gitaly!3586))
- [repository: Fix repo replication with transactions](gitlab-org/gitaly@62ed25098036ea988d3592ede7b72ac4220029a1) ([merge request](gitlab-org/gitaly!3630))
- [Fix ResolveConflicts RPC not calling Git hooks](gitlab-org/gitaly@11d8569fe62abcb19fdff014c6d1c142f430766a) ([merge request](gitlab-org/gitaly!3619))
- [Use ForceServerCodec() instead of CustomCodec()](gitlab-org/gitaly@c8bf3820f427d4dbf145238ffb6b4d44b22864f9) ([merge request](gitlab-org/gitaly!3624))
- [praefect: Fix incorrect error tracking for secondaries](gitlab-org/gitaly@b76154dd7dce4aa7b5d3fa7416d1697965b656d0) ([merge request](gitlab-org/gitaly!3620))
- [Bump Shell dep to break the dependency on the old Gitaly package](gitlab-org/gitaly@55d79bed009b8b3ac14dc178145120e0d0990206) ([merge request](gitlab-org/gitaly!3570))

### Changed (11 changes)

- [Set default Prometheus buckets for Gitalys RPC instrumentation](gitlab-org/gitaly@7c2c4253864e6efc9425a7884c4eaf0894a54a86) ([merge request](gitlab-org/gitaly!3669))
- [Support lazy failovers in `praefect dataloss`](gitlab-org/gitaly@e900df0947d188194f319b8afb9aa4fe857d357d) ([merge request](gitlab-org/gitaly!3549))
- [Update ffi gem to 1.15.3](gitlab-org/gitaly@0ef432a6ba8ab72fbf638f722516fb65faeb58b9) ([merge request](gitlab-org/gitaly!3664))
- [Bump actionpack, activesupport to 6.1](gitlab-org/gitaly@14084dc485cf798cfed3599463d4bc84ccfc3a0e) ([merge request](gitlab-org/gitaly!3661))
- [Update google-protobuf to v3.17.1 and labkit-ruby to v0.20.0](gitlab-org/gitaly@2f691bf85adfd20e66c046b883e11c694ec0ccc9) ([merge request](gitlab-org/gitaly!3656))
- [featureflag: Default-enable LFS pointers pipeline](gitlab-org/gitaly@5c43dfe6f4d8f97f0a7c33d6c8d9da0a0f2f3018) ([merge request](gitlab-org/gitaly!3653))
- [ListConflictFiles: Allow tree conflicts when specified](gitlab-org/gitaly@c5fefbdb306d0957815a9e808fc9638f7a98e10e) ([merge request](gitlab-org/gitaly!3648))
- [git: Accept commits and tags with malformed signatures](gitlab-org/gitaly@2da0b393998d394b743c70e7cf9cd0757a8f2733) ([merge request](gitlab-org/gitaly!3640))
- [Bump labkit dependency to v1.5.0, compile in stackdriver trace support](gitlab-org/gitaly@4961a32f5e13ca7b6105c3573619cffccefce065) ([merge request](gitlab-org/gitaly!3637))
- [Refactor update reference with hooks code](gitlab-org/gitaly@a17eb2f677f9710ee46ac4a9007b0cd828ca389a) ([merge request](gitlab-org/gitaly!3625))
- [featureflag: Default-enable GoUpdateRemoteMirror](gitlab-org/gitaly@c5b8403b7e5c46b0d60a5417b6142582c0a91340) ([merge request](gitlab-org/gitaly!3616))

### Removed (2 changes)

- [Remove support for virtual storage scoped primaries in read-only metrics](gitlab-org/gitaly@77c84dd5b41ca2bd7d22aa5c6409bd043499f73e) ([merge request](gitlab-org/gitaly!3548))
- [Remove support for virtual storage primaries in `praefect dataloss`](gitlab-org/gitaly@4d07d9b0fc94482ddd88a700624040a0005f031a) ([merge request](gitlab-org/gitaly!3549))

### Performance (9 changes)

- [Perform failovers lazily](gitlab-org/gitaly@3f09e462bd4dff39c54881acb1d8ce8e29cf2b21) ([merge request](gitlab-org/gitaly!3543))
- [Implement optimized RSS monitor based on /proc/[pid]/statm](gitlab-org/gitaly@82065592c3acfb92bef7aaebc8a33c35fbd58114) ([merge request](gitlab-org/gitaly!3646))
- [ref: Reimplement ListAllTags via object pipeline](gitlab-org/gitaly@77e9b5551ea01842218a5150dc9ec7292251f77f) ([merge request](gitlab-org/gitaly!3645))
- [replicator: Special-case replication of OptimizeRepository](gitlab-org/gitaly@1d47a0abfe386c62f15d8f49d71f6275f143aa92) ([merge request](gitlab-org/gitaly!3638))
- [replicator: Special-case replication of MidxRepack](gitlab-org/gitaly@45fcea2403eeb32987c84a7bf405a3a16d68bc32) ([merge request](gitlab-org/gitaly!3638))
- [replicator: Special-case replication of WriteCommitGraph](gitlab-org/gitaly@6f1601eaefcc247c3c7df6ee0a8dfdf579b0eb45) ([merge request](gitlab-org/gitaly!3638))
- [repository: Fix excessive voting in CreateRepositoryFromBundle](gitlab-org/gitaly@3d010219d28a5f5f49aeb60819e319d7101e99f1) ([merge request](gitlab-org/gitaly!3615))
- [blob: Speed up blob search via object type filters](gitlab-org/gitaly@d1dd987d926674aaa73315a99fc2ffe32d227c20) ([merge request](gitlab-org/gitaly!3590))
- [blob: Speed up LFS pointer search via object type filters](gitlab-org/gitaly@d2870b204c6801317a6e4c4fba09968fa6fd283d) ([merge request](gitlab-org/gitaly!3590))

### Other (4 changes)

- [Add StreamRPC library code](gitlab-org/gitaly@8a925b40a5e35848600600ee72441224f99af0fa) ([merge request](gitlab-org/gitaly!3601))
- [Refine metrics descriptions](gitlab-org/gitaly@d89639edbd90e1f0ca23a3bc63d8e70d15139004) ([merge request](gitlab-org/gitaly!3652))
- [Separate listenmux from backchannel](gitlab-org/gitaly@c95298c125a680a006153d5aca5d3dbb575ce352) ([merge request](gitlab-org/gitaly!3593))
- [Use upstream implementation of insecure credentials](gitlab-org/gitaly@7453f84b0bb385a958943b5e0910b8f6bb3906bb) ([merge request](gitlab-org/gitaly!3591))

## 14.0.7 (2021-08-03)

### Security (1 change)

- [catfile: Allow parsing of long git commit headers](gitlab-org/security/gitaly@f67249ab2dcd2b4c380d4dae5e889008229cd77c) ([merge request](gitlab-org/security/gitaly!37))

## 14.0.6 (2021-07-20)

### Added (1 change)

- [repository: Support transactional voting in `RemoveRepository()`](gitlab-org/gitaly@aa09579c973f8565a473fee8fc05c896ba8b1a9d) ([merge request](gitlab-org/gitaly!3675))

### Fixed (1 change)

- [coordinator: Fix repo creation/removal race for up-to-date secondaries](gitlab-org/gitaly@22e386d7cdecbb8023918918f7224a1791e975c6) ([merge request](gitlab-org/gitaly!3675))

## 14.0.5 (2021-07-08)

No changes.

## 14.0.4 (2021-07-07)

No changes.

## 14.0.3 (2021-07-06)

### Fixed (1 change)

- [repository: Fix repo replication with transactions](gitlab-org/gitaly@a483defd4600258e59ee21c358d71773dad50f58) ([merge request](gitlab-org/gitaly!3635))

### Performance (1 change)

- [repository: Fix excessive voting in CreateRepositoryFromBundle](gitlab-org/gitaly@22779fad58448e6dc4b0058fd4ff7486cd20cba7) ([merge request](gitlab-org/gitaly!3635))

## 14.0.2 (2021-07-01)

### Fixed (1 change)

- [repository: Fix repo replication with transactions](gitlab-org/security/gitaly@4d3ac6e8c20d88be0befed3b3966cc122288098f)

### Performance (1 change)

- [repository: Fix excessive voting in CreateRepositoryFromBundle](gitlab-org/security/gitaly@66832663b337c761e1948fe8dc2b385641e24655)

## 14.0.1 (2021-06-24)

No changes.

## 14.0.0 (2021-06-21)

### Added (8 changes)

- [remote: Allow for in-memory remotes in UpdateRemoteMirror](gitlab-org/gitaly@016625321a1c61c64fbd26655d0d8fb92f2f3bab) ([merge request](gitlab-org/gitaly!3566))
- [coordinator: Add replication metrics for transactions](gitlab-org/gitaly@bdf5df69bf32c8e736bd33bc4920d57e5ee3261e) ([merge request](gitlab-org/gitaly!3519))
- [remote: Vote when adding and removing remotes](gitlab-org/gitaly@bedf32172effabd52562fcf54df3600ff395d404) ([merge request](gitlab-org/gitaly!3507))
- [transactions: Fail early if the threshold cannot be reached anymore](gitlab-org/gitaly@7584c40e29ea3c1f5d2835f7622838d7d1b2c82a) ([merge request](gitlab-org/gitaly!3530))
- [featureflag: Remove LogCommandStats feature flag](gitlab-org/gitaly@745802f36733cf02d7195a3ccf8a7dfcac1a7296) ([merge request](gitlab-org/gitaly!3517))
- [remote: Vote when adding and removing remotes](gitlab-org/gitaly@9cd4553caf38ba381ca598cde1ada9ea7be13a11) ([merge request](gitlab-org/gitaly!3508))
- [repository: Enable transactional voting on the gitconfig](gitlab-org/gitaly@c3a9fd04c1c8ed3b3cb9c2d0b837c2e44b7cbb3c) ([merge request](gitlab-org/gitaly!3511))
- [repository: Replicate gitconfig](gitlab-org/gitaly@c28412f83bfb074908189c8bd42e3ea632a06efb) ([merge request](gitlab-org/gitaly!3511))

### Fixed (11 changes)

- [Fix Unix socket address handling following gRPC upgrade](gitlab-org/gitaly@d7a8c3abbb87dd7544b8a824fbbe0447ed5d841d) ([merge request](gitlab-org/gitaly!3592))
- [Do not track gRPC NotFound code as an error in Sentry for TreeEntry](gitlab-org/gitaly@49715b618e5cd3905ad743e152b88a5e1fe97190) ([merge request](gitlab-org/gitaly!3581))
- [Fix incorrect branchCreated result when startBranchName is provided](gitlab-org/gitaly@1de9f5fc18bcd92e47401ec3785159ff635f89c6) ([merge request](gitlab-org/gitaly!3562))
- [Don't create records in storage_repositories on generation increment](gitlab-org/gitaly@0134c6a540c90e3406848e7dd18b11739aea5ed6) ([merge request](gitlab-org/gitaly!3557))
- [Don't run housekeeping in Cleanup RPC](gitlab-org/gitaly@0849bcfa5ce3d734d535eed63b38667844401739) ([merge request](gitlab-org/gitaly!3507))
- [Disjoint request finalizer timeout from the RPC](gitlab-org/gitaly@be5fb6b267c1b0cae3ac18c707de39e48bf3624f) ([merge request](gitlab-org/gitaly!3515))
- [Vote when reference transaction has been committed](gitlab-org/gitaly@480dec51c438c89d4d9d20ac47307ea3f4311d4d) ([merge request](gitlab-org/gitaly!3514))
- [Consider primary modified only if a subtransaction was committed](gitlab-org/gitaly@d87747c82394e0ba0a2fd09a01e840dd9f6b6d27) ([merge request](gitlab-org/gitaly!3494))
- [repository: Remove housekeeping from Cleanup RPC](gitlab-org/gitaly@7a1d224d0ee6208350df12279b0f7319e9e0070c) ([merge request](gitlab-org/gitaly!3502))
- [Mark Repository service's Fsck as an accssor](gitlab-org/gitaly@c613e382f60ab261d7ce2b1c2ebf6531fd67641f) ([merge request](gitlab-org/gitaly!3499))
- [Respect failover disabled config option with per_repository elector](gitlab-org/gitaly@43ddbc3937304e8a22a201cb0fa105bb0af27a83) ([merge request](gitlab-org/gitaly!3491))

### Changed (17 changes)

- [featureflag: Remove reference transactions feature flag](gitlab-org/gitaly@9f296b8c194d10b777905ad8073445529ef3a503) ([merge request](gitlab-org/gitaly!3575))
- [Fix issues in tests](gitlab-org/gitaly@d32f88d89abe04a500ad64eabf8ea0d619f2259f) ([merge request](gitlab-org/gitaly!3403))
- [featureflag: Enable tx_config and tx_remote by default](gitlab-org/gitaly@332c5354b7a2f141f627d4a917d1df43e98d939b) ([merge request](gitlab-org/gitaly!3572))
- [operations: Skip precursory update of target ref in UserMergeToRef](gitlab-org/gitaly@907c03bf923e425274fd05129d3bf97ac7be33a4) ([merge request](gitlab-org/gitaly!3574))
- [Makefile: Upgrade Git to v2.32.0](gitlab-org/gitaly@8d4884c01f52eb8e3678de2ec837e50ae366f17e) ([merge request](gitlab-org/gitaly!3573))
- [Auto-resolve other conflict scenarios when AllowConflicts is true](gitlab-org/gitaly@8ea1987fd91810fe02bc3180f96729fddf5e74ee) ([merge request](gitlab-org/gitaly!3504))
- [repository: Relax URL check when fetching remotes](gitlab-org/gitaly@dc1a10393a12ffbe1bb850c287e6835f6058a0c1) ([merge request](gitlab-org/gitaly!3568))
- [UserRebaseConfirmable: Remove feature flag](gitlab-org/gitaly@8f55745792293c5fbe998c0794560a66a4db3766) ([merge request](gitlab-org/gitaly!3553))
- [Remove on-by-default gitaly_go_user_update_branch feature flag](gitlab-org/gitaly@c2cedf6d1a6131062fc58a43a11f4e85b6151bfc) ([merge request](gitlab-org/gitaly!3475))
- [Update default & secondary Go versions](gitlab-org/gitaly@a23fbc520eb1e82110a80ebe62c5e73712ff9fd5) ([merge request](gitlab-org/gitaly!3552))
- [logging: Drop topLevelGroup field](gitlab-org/gitaly@a2ba4b7c2b8745bd55fee19b2b62eecb0dc87e5d) ([merge request](gitlab-org/gitaly!3556))
- [Expand configuration of direct database connection](gitlab-org/gitaly@61e15d288450f6c8f4d242af496acc88a788a5ad) ([merge request](gitlab-org/gitaly!3495))
- [nodes: Mention gitaly in 'ErrPrimaryNotHealthy'](gitlab-org/gitaly@7f7273e5ff35d99feefec87260495d22bfc7682d) ([merge request](gitlab-org/gitaly!3551))
- [Do not fail over to outdated replicas](gitlab-org/gitaly@3309609bd8d38fb63a1c81638485af7725005618) ([merge request](gitlab-org/gitaly!3542))
- [Remove gitaly feature flag gitaly_go_user_revert](gitlab-org/gitaly@fcc18f919bcf259fbf259fb1358227e5d497716e) ([merge request](gitlab-org/gitaly!3507))
- [Cancel a vote associated with a node that stops waiting for a quorum](gitlab-org/gitaly@f58dd1af2f547ef959fe9dffa2f99e622f836936) ([merge request](gitlab-org/gitaly!3523))
- [Remove gitaly feature flag gitaly_go_user_revert](gitlab-org/gitaly@8949536f07581509949e3b37b4307937fcd42508) ([merge request](gitlab-org/gitaly!3516))

### Removed (2 changes)

- [Remove GrpcTreeEntryNotFound feature flag](gitlab-org/gitaly@7e47739dbebbb083316fadcb8874f833bb6bcd74) ([merge request](gitlab-org/gitaly!3567))
- [Prevent usage of other election strategies than per_repository](gitlab-org/gitaly@54948e21921527e0b455a56b601854d152e58ba2) ([merge request](gitlab-org/gitaly!3544))

### Security (1 change)

- [Update nokogiri from 1.11.1 to 1.11.5](gitlab-org/gitaly@4ad079594bd45eb51cc4ea3fc09a7cb8b2e2707d) ([merge request](gitlab-org/gitaly!3534))

### Performance (3 changes)

- [Avoid some allocations during diff parsing](gitlab-org/gitaly@f83247414d3ad4e3ec876644ea909328711c682e) ([merge request](gitlab-org/gitaly!3576))
- [blob: Improve latency and memory consumption for LFS pointers](gitlab-org/gitaly@8452f3daf6fa6d436cfed38a09fa2f76a54aa7dc) ([merge request](gitlab-org/gitaly!3507))
- [blob: Improve latency and memory consumption for LFS pointers](gitlab-org/gitaly@44678d2dfa47157c6a70594e66ba407f46e3a3b1) ([merge request](gitlab-org/gitaly!3483))

### Other (1 change)

- [Update gitlab-labkit to 0.17.1](gitlab-org/gitaly@904af72eb7c42124978370cef53681fa561b10f5) ([merge request](gitlab-org/gitaly!3395))

## 13.12.10 (2021-08-10)

No changes.

## 13.12.9 (2021-08-03)

### Security (1 change)

- [catfile: Allow parsing of long git commit headers](gitlab-org/security/gitaly@ae22c42551963b5fde79f55fd4f89c136a3d8dbb) ([merge request](gitlab-org/security/gitaly!38))

## 13.12.8 (2021-07-07)

No changes.

## 13.12.7 (2021-07-05)

### Fixed (1 change)

- [repository: Fix repo replication with transactions](gitlab-org/gitaly@74f3fd3902d7979845a063bb0bd6919316b58b4e) ([merge request](gitlab-org/gitaly!3632))

### Performance (1 change)

- [repository: Fix excessive voting in CreateRepositoryFromBundle](gitlab-org/gitaly@e324090114a8741ac966ef7b172463208463692c) ([merge request](gitlab-org/gitaly!3632))

## 13.12.6 (2021-07-01)

No changes.

## 13.12.5 (2021-06-21)

No changes.

## 13.12.4 (2021-06-14)

No changes.

## 13.12.3 (2021-06-07)

No changes.

## 13.12.2 (2021-06-01)

No changes.

## 13.12.1 (2021-05-25)

No changes.

## 13.12.0 (2021-05-22)

### Security (2 changes, 1 of them is from the community)

- Update golang.org/x/crypto to the latest to address CVE-2020-29652. !3400 (Takuya Noguchi)
- git: Always check fetched objects for consistency. !3458

### Removed (1 change)

- wiki: Remove FindFile RPC. !3454

### Fixed (8 changes)

- conflicts: Fix use of ambiguous refs in ResolveConflicts. !3386
- conflicts: Fix fetching from target repository. !3386
- conflicts: Fix segfault in case unresolved conflicts have no ancestor. !3386
- conflicts: Fix trailing newline handling when resolving conflicts. !3386
- Revert commit: Raise error if skipping commit offsets fails. !3390
- commit: Handle real errors when skipping commits failed. !3392
- Makefile: Unset PROFILE envvar before building git. !3414
- ssh: Fix secondaries being out-of-date if all refs are rejected. !3455

### Changed (1 change)

- Use the go implementation of UserRevert by default. !3438

### Performance (4 changes)

- Reduce memory usage in GetAllLFSPointers. !3379
- blob: Drop LFS pointer bitmap experiment. !3385
- conflicts: Drop ResolveConflicts feature flag. !3410
- featureflag: Activate Rebase implementation in Go. !3484

### Added (4 changes)

- Implement repository backups as per backup.rake. !3287
- gitlab-backup: Restore repositories as per backup.rake. !3383
- gitlab: Implement metric to measure latency of API calls. !3409
- remote: Add RemoteUrl parameter to FindRemoteRootRef. !3412

### Other (1 change)

- Multiplex connections between Praefect and Gitaly by default. !3360


## 13.11.7 (2021-07-07)

No changes.

## 13.11.6 (2021-07-01)

No changes.

## 13.11.5 (2021-06-01)

No changes.

## 13.11.4 (2021-05-14)

- No changes.

## 13.11.3 (2021-04-30)

- No changes.

## 13.11.2 (2021-04-27)

- No changes.

## 13.11.1 (2021-04-22)

- No changes.

## 13.11.0 (2021-04-22)

### Removed (1 change)

- Removal of the feature flag: distributed_reads. !3271

### Fixed (10 changes)

- repository: Fix default refspecs force-updating references. !3253
- Upgrade git version to v2.31.1. !3306
- Add CheckHostIP=no to SSH auth options for mirroring. !3312
- Close streamcache writer on all return paths. !3335
- Fail pipe writes when readers leave. !3341
- repository: Fix fetching in-memory remotes with SSH params. !3344
- remote: Fix UpdateRemoteMirror having transactional semantics. !3345
- ref: Fix missing votes for `DeleteRefs()` RPC. !3347
- operations: Fix UserRebaseConfirmable not using transactions. !3369
- proxy: Fix Goroutine leak in `forwardClientToServers()`. !3371

### Deprecated (1 change)

- Upgrade minimum required Go version to 1.15. !3352

### Changed (7 changes)

- CommitsBetween: learn to accept pagination params. !2484
- Enable gprc-go debug log messages with GRPC_GO_LOG_SEVERITY_LEVEL. !3266
- featureflag: Remove per-RPC transactional feature flags. !3284
- Turn UserUpdateBranch in Go on by default. !3286
- Remove upload_pack_gitaly_hooks feature flag. !3301
- featureflag: Default enable LogCommandStats. !3350
- Update activesupport to v6.0.3.6. !3373

### Performance (9 changes)

- Enable Go implementation for UserCherryPick. !3262
- Remove go_user_commit_files feature flag. !3281
- git: Generate reverse packfile indices. !3292
- blob: Remove feature flags for LFS pointer RPC ports. !3309
- repository: Allow fetching via in-memory remotes. !3321
- git: Use atomic fetches to allow for transactional behaviour. !3324
- featureflag: Remove reverse-packfile index feature flag. !3358
- featureflag: Remove `AtomicFetch` feature flag. !3359
- Makefile: Add custom patch to fix pathological perf with bitmap indices. !3362

### Added (5 changes)

- Add support for word-diff mode. !3086
- config: Allow injection of git config via Gitaly's config. !3279
- Integrate connection multiplexing into Gitaly and Praefect. !3293
- blob: Revamp interface for LFS pointers. !3316
- git: Bump minimum git version to git v2.31.0. !3340

### Other (1 change, 1 of them is from the community)

- Update gitlab-gollum-rugged_adapter to 0.4.4.4.gitlab.1. !3357 (Takuya Noguchi)


## 13.10.5 (2021-06-01)

No changes.

## 13.10.4 (2021-04-27)

- No changes.

## 13.10.3 (2021-04-13)

- No changes.

## 13.10.2 (2021-04-01)

- No changes.

## 13.10.1 (2021-03-31)

- No changes.

## 13.10.0 (2021-03-22)

### Removed (1 change)

- repository: Remove harmful `name` paramater in FetchRemote. !3227

### Fixed (13 changes)

- hooks: Fix inadvertent execution of hooks. !3119
- praefect: Stop creating replication jobs on cleanup. !3120
- operations: Fix transactions when deleting refs. !3144
- hook: Fix voting on pushes which delete packed references. !3146
- Reconcile missing repositories to assigned nodes. !3153
- coordinator: Fix replication and early failures when proxying transactional RPCs. !3158
- operations: Fix handling of update-ref failures in UserCommitFiles. !3159
- Fix UserCommitFiles index error handling. !3165
- blob: Fix filtering of new LFS pointers. !3175
- Prevent concurrent reconciliation. !3182
- coordinator: Fix inconsistent repository sizes when using reads distribution. !3209
- Fix text logging format erroring out. !3222
- operations: Fix deletion of branches with prefix. !3228

### Deprecated (1 change)

- Upgrading of the Go version. !3145

### Changed (8 changes)

- Remove Ruby code for old 100% go_user_create_{branch,tag} feature. !3056
- Gitaly config default for maintenance window. !3124
- Reconciliation sub-command performs as much as possible. !3142
- Remove repositories from unassigned storages. !3162
- featureflag: Remove UserFFBranch feature gate. !3180
- Add pktline side-band-64 writer. !3215
- gitaly-lfs-smudge: Validate OID parsed from LFS pointer. !3221
- Allow gitaly-ruby to run if log file cannot be written. !3224

### Performance (9 changes)

- Enable the Go port of UserCommitFiles by default. !3061
- blob: Port GetAllLFSPointers and GetLFSPointers to Go. !3173
- blob: Port GetNewLFSPointers to Go. !3195
- blob: Optimize Go implementation of LFS pointer lookup. !3210
- replicator: Do not bump repository generation for PackRefs. !3216
- featureflag: Enable Go port of Get{All,}LFSPointers. !3234
- blob: Enable use of bitmap indices when searching LFS pointers. !3238
- blob: Buffer output of git-catfile to speed up reading LFS pointers. !3241
- featureflag: Enable Go implementation of GetNewLFSPointers. !3252

### Added (6 changes)

- Add unique index for delete_replica replication events. !3183
- featureflag: Enable first batch of transactional RPCs. !3189
- featureflag: Enable transactional behaviour for all repository-scoped mutators. !3214
- housekeeping: Move cleanup of empty refs into housekeeping. !3246
- repository: Call housekeeping tasks in `OptimizeRepository()`. !3246
- housekeeping: Move cleanup of git config into housekeeping. !3246

### Other (2 changes)

- Streamio: remove custom ReadFrom and WriteTo. !3201
- Remove unused Ruby code. !3203


## 13.9.7 (2021-04-27)

- No changes.

## 13.9.6 (2021-04-13)

- No changes.

## 13.9.5 (2021-03-31)

- No changes.

## 13.9.4 (2021-03-17)

### Changed (1 change)

- Allow gitaly-ruby to run if log file cannot be written. !3236


## 13.9.3 (2021-03-08)

- No changes.

## 13.9.2 (2021-03-04)

- No changes.

## 13.9.1 (2021-02-23)

### Fixed (1 change)

- Fix 500 errors in Wiki pages with trailers containing UTF-8. !3170


## 13.9.0 (2021-02-22)

### Fixed (9 changes)

- operations: Fix hooks running on secondaries when creating annotated tags. !3022
- transactions: Optionally use timestamps for deterministic results. !3036
- repository: Fix regressions in FetchRemote. !3043
- Make gitaly_ruby_json.log work again. !3052
- coordinator: Fix outdated repos not getting repljobs with transactions. !3055
- hook: Stop transactions when post-receive and update hooks fail. !3094
- Fix premature cgroups cleanup. !3098
- hook: Increase the timeout when casting votes. !3115
- localrepo: Fix lookup of wrong ref if requesting prefix. !3127

### Changed (10 changes)

- Port UserUpdateBranch to Go. !3013
- Remove Ruby code for on-by-default go_user_delete_{branch,tag} feature flags. !3033
- Enable feature flag go_user_create_{branch,tag} by default. !3035
- Enable go implementation of UserFFBranch by default. !3057
- gitaly-lfs-smudge: Clean up URL building. !3058
- Intercept RepositoryExists calls in Praefect. !3075
- ruby: Upgrade to Rugged 1.0. !3076
- Standardize Praefect and Gitaly log formats. !3121
- Ignore SIGURG in gitaly-wrapper. !3131
- Standardize Praefect and Gitaly timestamps. !3133

### Performance (3 changes)

- featureflags: Remove GoUserMergeBranch feature flag. !3049
- featureflag: Remove GoFetchSourceBranch feature flag. !3050
- Restrict number of threads for a full repack. !3108

### Added (7 changes)

- Track feature flags used for RPC call. !2971
- operations: Wire up AllowConflicts handling for Go. !2997
- git: Add support for options which always get injected. !3028
- repository: Cleanup stale lockfiles when running housekeeping. !3051
- repository: Use transactions when writing gitattributes. !3064
- Configurable default replication factor for virtual storages. !3091
- hook: Use proper error codes when transactions fail. !3097

### Other (1 change)

- Upgrade labkit-ruby to v0.15.0. !3118


## 13.8.8 (2021-04-13)

- No changes.

## 13.8.7 (2021-03-31)

- No changes.

## 13.8.6 (2021-03-17)

- No changes.

## 13.8.5 (2021-03-04)

### Fixed (1 change)

- Fix 500 errors in Wiki pages with trailers containing UTF-8. !3169


## 13.8.4 (2021-02-11)

- No changes.

## 13.8.3 (2021-02-05)

- No changes.

## 13.8.2 (2021-02-01)

- No changes.

## 13.8.1 (2021-01-26)

- No changes.

## 13.8.0 (2021-01-22)

### Security (2 changes)

- Bump actionpack gem to 6.0.3.4. !2982
- grpc: raise minimum TLS version to 1.2. !2985

### Removed (2 changes)

- Removal of ruby implementation of the FetchRemote. !2967
- Removal of ruby implementation of the UserSquash. !2968

### Fixed (8 changes)

- operations: Fix Go UserMergeBranch failing with ambiguous references. !2921
- hooks: Correctly filter environment variables for custom hooks. !2933
- Fix wrongly labeled prometheus metrics for limithandler. !2955
- Fix internal API errors not being passed back to UserMergeBranch. !2987
- repository: Silence progress meter of FetchSourceBranch. !2991
- operations: Fix UserFFBranch if called on an ambiguous reference. !2992
- Fix ResolveConflicts file limit error. !3004
- repository: Fix ReplicateRepository returning before RPCs have finished. !3011

### Changed (6 changes)

- praefect: intercept CreateRepository* RPCs to populate database. !2873
- Support repository specific primaries and host assignments in dataloss. !2890
- Port UserCreateTag to Go. !2911
- Drop unused assigned column. !2972
- Enable feature flag go_user_delete_{branch,tag} by default. !2994
- FindCommit[s]: add a Trailers boolean flag to do %(trailers) work. !2999

### Performance (5 changes)

- Don't query for primary for read operations. !2909
- Feature flag gitaly_distributed_reads enabled by default. !2960
- featureflags: Enable Go implementation of UserMergeBranch by default. !2976
- repository: Short-circuit fetches when objects exist already. !2980
- Disable ref tx hooks for FetchRemote calls. !3002

### Added (3 changes)

- Parse Git commit trailers when processing commits. !2842
- Add information about whether tags were updated to the FetchRemote RPC. !2901
- objectpool: Count normal references when collecting stats. !2993

### Other (1 change)

- Make command stats logging concurrency-safe. !2956


## 13.7.9 (2021-03-17)

- No changes.

## 13.7.8 (2021-03-04)

- No changes.

## 13.7.7 (2021-02-11)

- No changes.

## 13.7.6 (2021-02-01)

- No changes.

## 13.7.5 (2021-01-25)

### Performance (1 change)

- Disable ref tx hooks for FetchRemote calls. !3006


## 13.7.4 (2021-01-13)

- No changes.

## 13.7.3 (2021-01-08)

- No changes.

## 13.7.2 (2021-01-07)

- No changes.

## 13.7.1 (2020-12-23)

- No changes.

## 13.7.0 (2020-12-22)

### Removed (1 change)

- Remove MemoryRepositoryStore. !2845

### Fixed (22 changes)

- command: Fix panics and segfaults caused by LogCommandStats. !2791
- Praefect reconcile hangs and fails in case of an error during processing. !2795
- nodes: Use context to perform database queries. !2796
- Fix `updateReferenceWithHooks()` not forwarding stderr. !2804
- Discard git-rev-parse error messages in. !2809
- hooks: Improved validation and testing. !2824
- Remove records of invalid repositories. !2833
- User{Branch,Submodule}: remove erroneously copy/pasted error handling. !2841
- Evict broken connections from the pool. !2849
- UserCreateBranch: unify API responses between Go & Ruby response paths. !2857
- UserDeleteBranch: unify API responses between Go & Ruby response paths. !2864
- operations: Fix wrong ordering of merge parents for UserMergeBranch. !2868
- Run sql electors checks in a transaction. !2869
- hooks: Fix ambiguous envvars conflicting with gitaly-ssh. !2874
- UserCreateTag: stop dying when a tag pointing to tree or blob is created + test fixes. !2881
- CreateFork recovers when encountering existing empty directory target. !2886
- Handle nil index entries in resolve conflicts. !2895
- Update github-linguist to v7.12.1. !2897
- Update resolve conflict command to use gob over stdin. !2903
- Fix missing cgroups after upgrading Gitaly. !2914
- Run housekeeping on pool repository. !2916
- User{Branch,Tag,Submodule}: ferry update-ref errors upwards. !2926

### Changed (9 changes)

- Make git gc --prune more aggressive. !2758
- featureflag: Enable Go implementation of UserSquash. !2807
- Port UserDeleteTag to Go. !2839
- Print host assignments and primary per repository in `praefect dataloss`. !2843
- transactions: Allow disabling with an env var. !2853
- No longer compare checksums after replication. !2861
- Reintroduce assignment schema change without dropping the old column. !2867
- Revert featureflag: Remove reference transaction feature flag. !2884
- Cleanup redundant data from notification events. !2893

### Performance (3 changes)

- git: Speed up creation of packfiles via pack window memory limit. !2856
- git2go: Restrict number of computed virtual merge bases. !2860
- Disable hooks when fetching. !2923

### Added (11 changes)

- Introduction of in-memory cache for reads distribution. !2738
- Support for logging propagated client identity. !2802
- Add initial implementation of spawning git inside cgroups. !2819
- Tell Git where to find reference-transaction hooks. !2834
- Conditionally enable use of transactions for all reference-modifying RPCs. !2850
- Set replication factor for a repository. !2851
- hooks: Remove the Ruby reference-transaction hook feature flag. !2866
- Enable feature flag gitaly_go_fetch_remote by default. !2872
- conflicts: Remove GoListConflictFiles feature flag. !2878
- operations: Remove GoUserMergeToRef feature flag. !2879
- Perform housekeeping for object pools. !2885

### Other (5 changes)

- Instrument git-cat-file's batch commands for more granular tracing. !2687
- Update LabKit to v1.0.0. !2827
- Update Rouge gem to v3.25.0. !2829
- Support Golang v1.15.5 in CI. !2858
- Update Rouge gem to v3.26.0. !2927


## 13.6.7 (2021-02-11)

- No changes.

## 13.6.6 (2021-02-01)

- No changes.

## 13.6.5 (2021-01-13)

- No changes.

## 13.6.4 (2021-01-07)

- No changes.

## 13.6.3 (2020-12-10)

- No changes.

## 13.6.2 (2020-12-07)

- No changes.

## 13.6.1 (2020-11-23)

- No changes.

## 13.6.0 (2020-11-22)

### Security (1 change)

- Configure the GitLab internal API client to provide client certificates when using TLS. !2794

### Fixed (12 changes)

- config: Fix check for executability on systems with strict permissions. !2668
- Create missing directories in CreateRepositoryFromSnapshot. !2683
- operations: Fix feature flag for UserMergeToRef. !2689
- operations: Return correct error code when merge fails. !2690
- Fall back to /dev/null when opening gitaly_ruby_json.log. !2708
- git: Recognize "vX.Y.GIT" versions. !2714
- hooks: Always consume stdin for reference transaction hook. !2719
- gitaly: Fix deadlock when writing to gRPC streams concurrently. !2723
- Use new correlation ID generator. !2746
- operations: Always set GL_PROTOCOL in hooks. !2753
- operations: Fix error message when UserMergeToRef conflicts. !2756
- Fix handling of symlinks in custom hooks directory. !2790

### Changed (6 changes)

- hooks: Check command ported to Go. !2650
- Expose ancestor path in Conflicts RPC. !2672
- Remove primary-wins and reference-transaction-hook feature flags. !2681
- featureflag: Enable Ruby reference transaction hooks by default. !2717
- featureflag: Remove reference transaction feature flag. !2725
- git: Upgrade minimum required version to v2.29.0. !2727

### Performance (4 changes)

- Port UserCommitFiles to Go. !2655
- Port ResolveConflicts from Ruby to Go. !2693
- featureflag: Enable Go implementation of ListConflictFiles. !2782
- featureflag: Enable Go implementation of UserMergeToRef. !2783

### Added (8 changes)

- Port UserCreateBranch to go. !2613
- Receiving notifications on changes in database. !2631
- hooks: Set Gitaly as user agent for API calls. !2663
- Add JSON request logging for gitaly-ruby. !2678
- Enforce minimum required Git version. !2701
- proto: Add Tree ID to GitCommit structure. !2703
- Log LFS smudge activity to gitaly_lfs_smudge.log. !2734
- Create database records for repositories missing them. !2749

### Other (10 changes)

- Ensure reference hooks are used in valid commands. !2583
- Update Ruby to v2.7.2. !2633
- docs: remove feature flag references for hooks. !2658
- Instrument git commands for tracing. !2685
- Update ruby parser for Ruby v2.7.2. !2699
- Update to bundler v2.1.4. !2733
- Store repository host node assignments. !2737
- Use labkit-ruby 0.13.2. !2743
- Remove the RepositoryService.FetchHTTPRemote RPC. !2744
- Improve logging in ReplicateRepository. !2767


## 13.5.7 (2021-01-13)

- No changes.

## 13.5.6 (2021-01-07)

- No changes.

## 13.5.5 (2020-12-07)

- No changes.

## 13.5.4 (2020-11-13)

- No changes.

## 13.5.3 (2020-11-03)

- No changes.

## 13.5.2 (2020-11-02)

### Security (1 change)

- Removal of all http.*.extraHeader config values.


## 13.5.1 (2020-10-22)

- No changes.

## 13.5.0 (2020-10-22)

### Fixed (9 changes)

- Pass correlation ID to hooks. !2576
- transactions: Correctly handle cancellation of empty transactions. !2594
- Ensure local branches are current when syncing remote mirrors. !2606
- Fix injection of gitaly servers info. !2615
- Verification of gitaly-ssh runs. !2616
- git: Fix parsing of dashed -rc versions. !2639
- hook: Stop transactions on pre-receive hook failure. !2643
- Doubled invocation of gitaly-ssh on upload pack cmd. !2645
- operations: Fix PostReceive hook receiving no input. !2653

### Changed (6 changes)

- transactions: Only vote when reftx hook is in prepared state. !2571
- transactions: Remove voting via pre-receive hook. !2578
- linguist: Bump version for better detection. !2591
- transactions: Remove service-specific feature flags. !2599
- Disabling of reads distribution feature. !2619
- Send CORRELATION_ID to gitaly-lfs-smudge filter. !2662

### Performance (2 changes)

- Port operations.UserMergeToRef to Go. !2580
- conflicts: Port ListConflictFiles to Go. !2598

### Added (9 changes)

- transactions: Implement RPC to gracefully stop transactions. !2532
- Add Git LFS smudge filter. !2577
- git2go: Implement new command to list merge conflicts. !2587
- PostgreSQL notifications listener. !2603
- Add include_lfs_blobs flag in GetArchiveRequest RPC. !2607
- Add support for using LFS smudge filter. !2621
- Port UserSquash to Go. !2623
- Add option to include conflict markers to merge commit. !2626
- Per-connection gRPC options in client. !2630

### Other (6 changes)

- Reference hook option for Git command DSL. !2596
- Update json gem to v2.3.1. !2610
- Fix Ruby 2.7 keyword deprecation deprecation warnings. !2611
- Refactor server metadata to be more type safe. !2624
- Remote repository abstraction for resolving refish. !2629
- Upgrade Rubocop to 0.86.0. !2634


## 13.4.7 (2020-12-07)

- No changes.

## 13.4.6 (2020-11-03)

- No changes.

## 13.4.5 (2020-11-02)

### Security (1 change)

- Removal of all http.*.extraHeader config values.


## 13.4.4 (2020-10-15)

- No changes.

## 13.4.3 (2020-10-06)

- No changes.

## 13.4.2

- No changes.

## 13.4.1

- No changes.

## 13.4.0

- No changes.
### Removed (1 change)

- Remove server scoped handling from coordinator. !2546

### Fixed (13 changes)

- Change timeformat to not omit trailing 0s. !256
- Fix sparse checkout file list for rebase onto remote branch. !2447
- Fix Git hooks when GitLab relative URL path and UNIX socket in use. !2485
- Relabel smarthttp.InfoRefsReceivePack as accessor. !2487
- Move fake git path to test target. !2490
- Fix potentially executing linguist with wrong version of bundle. !2495
- Fixup reference-transaction hook name based on arguments. !2506
- Fix GIT_VERSION build flag overriding Git's version. !2507
- Fix stale connections to Praefect due to keepalive policy. !2511
- Fix downgrade error handling. !2522
- Makefile: Avoid Git2Go being linked against stale libgit2. !2525
- Pass correlation_id over to gitaly-ssh. !2530
- Fix logging of replication events processing. !2547

### Changed (6 changes)

- hooks: Remove update feature flag. !2501
- hooks: Remove prereceive hook Ruby implementation. !2519
- transactions: Enable majority-wins voting strategy by default. !2529
- Export GL_REPOSITORY and GL_PROJECT_PATH in git archive call. !2557
- Enable voting via reference-transaction hook by default. !2558
- Introduce Locator abstraction to diff service. !2559

### Performance (3 changes)

- Improved SQL to get up to date storages for repository. !2514
- Port OperationService.UserMergeBranch to Go. !2540
- git: Optimize check for reference existence. !2549

### Added (12 changes)

- Use error tracker to determine if node is healthy. !2341
- Daily maintenance scheduler. !2423
- Bump default Git version to v2.28.0. !2432
- In-memory merges via Git2Go. !2433
- Add Git2Go integration. !2438
- Replication job acknowledge removes 'completed' and 'dead' events. !2457
- Automatic repository reconciliation. !2462
- Implement majority-wins transaction voting strategy. !2476
- Provide generic "git" Makefile target. !2488
- Rebuild targets only if Makefile content changes. !2492
- Transactional voting via reference-transaction hook. !2509
- hooks: Call reference-transaction hook from Ruby HooksService. !2566

### Other (6 changes)

- Add fuzz testing to objectinfo parser. !2481
- Update rbtrace gem to v0.4.14. !2491
- Upgrade sentry-raven and Faraday gems to v1.0.1. !2533
- Include grpc_service in gitaly_service_client_requests_total metric. !2536
- Bump labkit dependency to get mutex profiling. !2562
- Update Nokogiri gem to v1.10.10. !2567


## 13.3.9 (2020-11-02)

### Security (1 change)

- Removal of all http.*.extraHeader config values.


## 13.3.8 (2020-10-21)

- No changes.

## 13.3.6

- No changes.

## 13.3.5

- No changes.

## 13.3.4

- No changes.

## 13.3.3 (2020-09-02)

### Security (1 change)

- Don't expand filesystem paths of wiki pages.


## 13.3.2 (2020-08-28)

### Fixed (1 change)

- Fix hanging info refs cache when error occurs. !2497


## 13.3.1 (2020-08-25)

- No changes.

## 13.3.0 (2020-08-22)

### Removed (1 change)

- Remove Praefect primary from config. !2392

### Fixed (12 changes)

- Praefect: storage scoped RPCs not handled properly. !2234
- Fix parsing of Git release-candidate versions. !2389
- Fix push options not working with Gitaly post-receive hook. !2394
- Fix detection of context cancellation for health updater. !2397
- lines.Send() only support byte delimiter. !2402
- Fix Praefect not starting with unhealthy Gitalys. !2422
- Nodes elector for configuration with disabled failover. !2444
- Fix connecting to Praefect service from Ruby sidecar. !2451
- Fix transaction voting delay metric for pre-receive hook. !2458
- Fix accessors mislabeled as mutators. !2466
- Fix registration of Gitaly metrics dependent on config. !2467
- Fix post-receive hooks with reference transactions. !2471

### Changed (12 changes)

- Improve query to identify up to date storages for reads distribution. !2372
- Add old path to NumStats protobuf output. !2395
- Generate data loss report from repository generation info. !2403
- Enforce read-only mode per repository. !2405
- Remove remote_branches_ls_remote feature flag. !2417
- Remove virtual storage wide read-only mode. !2431
- Update gRPC to v1.30.2 and google-protobuf to v3.12.4. !2442
- Report only read-only repositories by default in dataloss. !2449
- Configurable replication queue batch size. !2450
- Use repository generations to determine the best leader to elect. !2459
- Default-enable primary-wins reference transaction. !2460
- Enable distributed_reads feature flag by default. !2470

### Performance (1 change)

- Log cumulative per-request rusage ("command stats"). !2368

### Added (13 changes, 1 of them is from the community)

- GetArchive: Support path elision. !2342 (Ethan Reesor (@firelizzard))
- Praefect: include PgBouncer in CI. !2378
- Support dry-run cherry-picks and reverts. !2382
- Add subtransactions histogram metric. !2390
- Export connection pool and support setting DialOptions. !2401
- Queue replication jobs in case a transaction is unused. !2404
- Add support for primary-wins voting strategy. !2408
- Add accept-dataloss sub-command. !2415
- Metric for the number of read-only repositories. !2426
- Improve transaction metrics. !2441
- Praefect: replication processing should omit unhealthy nodes. !2464
- Log transaction state when cancelling them. !2465
- Prune objects older than 24h on garbage collection. !2469

### Other (2 changes)

- Update mime-types for Ruby 2.7. !2456
- Pass CORRELATION_ID env variable to spawned git subprocesses. !2478


## 13.2.9

### Fixed (1 change)

- Fix hanging info refs cache when error occurs. !2497


## 13.2.8

- No changes.

## 13.2.7 (2020-09-02)

### Security (1 change)

- Don't expand filesystem paths of wiki pages.


## 13.2.6

- No changes.

## 13.2.5

- No changes.

## 13.2.4

- No changes.

## 13.2.3

### Security (1 change)

- Fix injection of arbitrary `http.*` options.


## 13.2.2

- No changes.

## 13.2.1

- No changes.

## 13.2.0

- No changes.
### Fixed (13 changes)

- Set default branch to match remote for FetchInternalRemote. !2265
- Make --literal-pathspecs optional for LastCommitForPath. !2285
- Only execute hooks on primary nodes. !2294
- Avoid duplicated primary when getting synced nodes. !2312
- Force gitaly-ruby to use UTF-8 encoding. !2316
- Always git fetch after snapshot repository creation. !2320
- Return nil for missing wiki page versions. !2323
- Pass env vars correctly into custom hooks. !2331
- Fix getting default branch of remote repo in. !2340
- Fix casting wrong votes on secondaries using Ruby pre-receive hook. !2347
- Improve logging on the replication manager. !2351
- Only let healthy secondaries take part in transactions. !2365
- Fix pre-receive hooks not working with symlinked paths variable field. !2381

### Changed (7 changes, 1 of them is from the community)

- Support literal path lookups in LastCommitsForTreeRequest. !2301
- GetArchive: Support excluding paths. !2304 (Ethan Reesor (@firelizzard))
- Turn on go prereceive and go update by default. !2329
- Gather remote branches via ls-remote, rather than fetch, by default. !2330
- Include change type as a label on in-flight replication jobs gauge. !2373
- Scale transaction registration metric by registered voters. !2375
- Allow multiple votes per transaction. !2386

### Added (12 changes)

- Multi node write. !2208
- Allow pagination for FindAllLocalBranches. !2251
- Add TLS support to Praefect. !2276
- PostReceiveHook in Go. !2290
- Praefect: replication jobs health ping. !2321
- Praefect: handling of stale replication jobs. !2322
- Implement weighted voting. !2334
- Schedule replication jobs for failed transaction voters. !2355
- Praefect: Collapse duplicate replication jobs. !2357
- Start using transactions for UserCreateBranch. !2364
- Transaction support for modification of branches and tags via operations service. !2374
- Remove upload-pack feature flags to enable partial clone by default.

### Other (3 changes)

- Support literal path lookups in other commit RPCs. !2303
- Ensure Praefect replication queue tables exist. !2309
- Error forwarded mutator RPCs on replication job enqueue failure. !2332


## 13.1.11

### Fixed (1 change)

- Fix hanging info refs cache when error occurs. !2497


## 13.1.10

- No changes.

## 13.1.9 (2020-09-02)

### Security (1 change)

- Don't expand filesystem paths of wiki pages.


## 13.1.8

- No changes.

## 13.1.7

- No changes.

## 13.1.6

### Security (1 change)

- Fix injection of arbitrary `http.*` options.


## 13.1.5

### Fixed (1 change)

- Fix pre-receive hooks not working with symlinked paths variable field. !2381


## 13.1.3

### Fixed (1 change)

- Fix HTTP proxies not working in Gitaly hooks. !2325

### Changed (1 change)

- Add GL_PROJECT_PATH for custom hooks. !2313


## 13.1.2

### Security (1 change)

- Add random suffix to worktree paths to obstruct path traversal.


## 13.1.1

- No changes.

## 13.1.0

### Fixed (9 changes)

- Stale packed-refs.new files prevents branches from being deleted. !2206
- Praefect graceful stop. !2210
- Check auth before limit handler. !2221
- Warn if SO_REUSEPORT cannot be set. !2236
- Praefect: unhealthy nodes must not be used for read distribution. !2250
- Pass git push options into pre receive. !2266
- Allow more frequent keep-alive checking on server. !2272
- Adjust Praefect server address based on peer info. !2273
- Replication not working on Praefect. !2281

### Deprecated (1 change)

- Upgrade to Git 2.27. !2237

### Changed (8 changes)

- PreReceive in go. !2155
- Remote branches via ls-remote is now a toggle. !2183
- Run replication jobs against read-only primaries. !2190
- Praefect: Enable database replication queue by default. !2193
- Remove feature flag go_fetch_internal_remote. !2203
- Skip creation of gitlab api if GITALY_TESTING_NO_GIT_HOOKS is set. !2245
- Set a stable signature for .patch endpoints to create reproducible patches. !2253
- Improved dataloss subcommand. !2278

### Performance (2 changes)

- OptimizeRepository will remove empty ref directories. !2204
- Decrease memory consuption when parsing tree objects. !2241

### Added (9 changes)

- Allow transaction manager to handle multi-node transactions. !2182
- Add end-of-options to supported commands. !2192
- Praefect gauge for replication jobs scoped by storage. !2207
- failover: Default to enabling SQL strategy. !2218
- Only log relevant storages in Praefect dial-nodes. !2229
- Add support for filter-repo commit map to cleaner. !2247
- How to handle proxying FindRemoteRepository. !2260
- Distribute reads between all shards, including primaries. !2267
- Expose ref names in list commits by ref name response. !2269

### Other (2 changes)

- Bump Ruby to v2.6.6. !2231
- danger: Suggest merge request ID in the changelog. !2254


## 13.0.14

- No changes.

## 13.0.13

- No changes.

## 13.0.12

### Security (1 change)

- Fix injection of arbitrary `http.*` options.


## 13.0.11

This version has been skipped due to packaging problems.

## 13.0.10

- No changes.

## 13.0.9

- No changes.

## 13.0.8

### Security (1 change)

- Add random suffix to worktree paths to obstruct path traversal.


## 13.0.7

- No changes.

## 13.0.6

- No changes.

## 13.0.5

- No changes.

## 13.0.4

### Fixed (1 change)

- Clean configured storage paths. !2223


## 13.0.3

- No changes.

## 13.0.2

- No changes.

## 13.0.1

- No changes.

## 13.0.0

- No changes.
### Security (1 change)

- Improved path traversal protection. !2132

### Fixed (16 changes, 1 of them is from the community)

- Delete tags by canonical reference name when mirroring repository. !2026
- Fix signature detection for tags. !2045 (Roger Meier)
- Ignore repositories deleted during a file walk. !2083
- Revert gRPC upgrade to v1.27.0 to fix issues on multiple platforms. !2088
- Do not enable SQL elector if failover is disabled. !2091
- cleanup commit-graph-chain.lock file after crash. !2099
- Provide consistent view of primary and secondaries. !2105
- Fix rebase when diff contains only deleted files. !2109
- Praefect: proper multi-virtual storage support. !2117
- Use tableflip with praefect prometheus listener. !2122
- Praefect: configuration verification. !2130
- HTTPSettings to handle bools and ints. !2142
- Bump gitlab-markup gem to v1.7.1. !2143
- Revert charlock holmes bump. !2154
- Configure logging before running sub-commands. !2169
- Allow port reuse in tableflip. !2175

### Changed (9 changes)

- Allow Praefect's ServerInfo RPC to succeed even if the internal gitaly node calls fail. !2067
- Choose primary with smallest replication queue size. !2078
- Use go update hook when feature flag is toggled. !2095
- Modify chunker to send message based on size. !2096
- Use separate http settings config object. !2104
- Extract reference transaction manager from transaction service. !2114
- Enable feature flag for go update hooks through operations service. !2120
- Block in NewManager to wait for nodes to be up. !2134
- Include Praefect usage in the usage ping. !2180

### Added (10 changes)

- Add DivergentRefs to UpdateRemoteMirrorResponse. !2028
- Write gitlab shell config from gitaly to gitlab shell as env var. !2066
- Implement reference transaction service. !2077
- Reconciliation should report progress and warn user. !2084
- Improve error messages for repository creation RPCs. !2118
- Upgrade github-linguist to version 7.9.0. !2145
- Single-node transactions via pre-receive hook. !2147
- Enforce read-only status for virtual storages. !2148
- Extract client name from Go GRPC client. !2152
- Praefect enable-writes subcommand. !2157

### Other (3 changes)

- Upgrade activesupport and related Ruby gems to v6.0.2.2. !2110
- Update ffi gem to v1.12.2. !2111
- Update activesupport to v6.0.3 and gitlab-labkit to v0.12.0. !2178

## 12.10.14

- No changes.

## 12.10.13

### Security (1 change)

- Add random suffix to worktree paths to obstruct path traversal.


## 12.10.12

- No changes.

## 12.10.11

- No changes.

## 12.10.10

- No changes.

## 12.10.7

- No changes.

## 12.10.6

- No changes.

## 12.10.4

- No changes.

## 12.10.3

- No changes.

## 12.10.2

### Security (1 change)

- gems: Upgrade nokogiri to > 1.10.7. !2128


## 12.10.1

- No changes.

## 12.10.0

#### Added
- Praefect: Postgres queue implementation in use
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1989
- Adding metrics to track which nodes are up and down
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2019
- RPC ConsistencyCheck and Praefect reconcile subcommand
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1903
- Add gitaly-blackbox prometheus exporter
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1860
- Add histogram to keep track of node healthcheck latency
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1921
- Praefect dataloss subcommand
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2057
- Add metric for replication delay
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1997
- Add a metric counting the number of negotiated packfiles
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2011
- Praefect: Enable Postgres binary protocol
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1995
- Add repository profile
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1959
- Don't push divergent remote branches with keep_divergent_refs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1915
- Add metric counter for mismatched checksums after replication
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1943
- Praefect: replication event queue as a primary storage of events
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1948
- Add SQL-based election for shard primaries
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1979
- Propagate GarbageCollect, RepackFull, RepackIncremental to secondary nodes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1970
- Call hook rpcs from operations service
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2034
- Enable client prometheus histogram
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1987
- Add Praefect command to show migration status
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2041

#### Changed
- Support ignoring unknown Praefect migrations
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2039
- Make Praefect sql-migrate ignore unknown migrations by default
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2058
- Add EnvironmentVariables field in hook rpcs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1969
- Explicitly check for existing repository in CreateRepositoryFromBundle
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1980

#### Deprecated
- Drop support for Gitaly v1 authentication
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2024
- Upgrade to Git 2.26
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1983

#### Fixed
- Commit signature parsing must consume all data
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1953
- Allow commits with invalid timezones to be pushed
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1947
- Check for git-linguist error code
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1923
- Praefect: avoid early request cancellation when queueing replication jobs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2062
- Race condition on ProxyHeaderWhitelist
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2025
- UserCreateTag: pass tag object to hooks when creating annotated tag
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1956
- Fix flaky test TestLimiter
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1965
- Exercise Operations service tests with auth
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2059
- Fix localElector locking issues
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2030
- Pass gitaly token into gitaly-hooks
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2035
- Validate offset/limit for ListLastCommitsForTree
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1996
- Use reference counting in limithandler middleware
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1984
- Modify Praefect's server info implementation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1991

#### Other
- Refactor Praefect node manager
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1940
- update ruby gems grpc/grpc-tools to 1.27.0 and google-protobuf to 3.11.4
  https://gitlab.com/gitlab-org/gitaly/merge_requests/
- Update parser and unparser Ruby gems
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1961
- Remove feature flags for InfoRef cache
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2038
- Static code analysis: unconvert
  https://gitlab.com/gitlab-org/gitaly/merge_requests/2046

#### Performance
- Call Hook RPCs from gitaly-hooks binary
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1740

#### Removed
- Drop go 1.12 support
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1976
- Remove gitaly-remote command
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1992

#### Security
- Validate content of alternates file
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1946

## 12.9.10

- No changes.

## 12.9.9

- No changes.

## 12.9.8

- No changes.

## 12.9.7

- No changes.

## 12.9.6

- No changes.

## 12.9.5

### Security (1 change)

- gems: Upgrade nokogiri to > 1.10.7. !2129


## 12.9.4

- No changes.

## 12.9.3

- No changes.

## 12.9.2

- No changes.

## 12.9.1

- No changes.

## 12.9.0

### Security (1 change)

- Validate object pool relative paths to prevent path traversal attacks. !1900

### Fixed (11 changes, 1 of them is from the community)

- Handle malformed PID file. !1825 (maxmati)
- Handle ambiguous refs in CommitLanguages. !1829
- Enforce diff.noprefix=false for generating Git diffs. !1854
- Fix expected porcelain output for PushResults. !1862
- Properly account for tags in PushResults. !1874
- ReplicateRepository error when result from FetchInternalRemote is false. !1879
- Praefect should not emit Gitaly errors to Sentry. !1880
- Task proto has dependency to already generated source code. !1884
- Explicit error what type of path can't be read. !1891
- Allow filters when advertising refs. !1894
- Fix gitaly-ruby not starting on case-sensitive filesystems. !1939

### Changed (6 changes)

- Change ListRepositories RPC to RepostoryReplicas. !1692
- Remove deprecated UserRebase RPC. !1851
- Replication: propagate RenameRepository RPC to Praefect secondaries. !1853
- Add node gauge that keeps track of node status. !1904
- Praefect: use enum values for job states. !1906
- Use millisecond precision for time in JSON logs.

### Performance (1 change)

- Use Rugged::Repository#bare over #new. !1920

### Added (11 changes)

- Praefect: add sql-migrate-down subcommand. !1770
- Praefect SQL: support of transactions. !1815
- Optionally keep divergent refs when mirroring. !1828
- Push with the --porcelain flag and parse output of failed pushes. !1845
- Internal RPC for walking Gitaly repos. !1855
- Praefect: Move replication queue to database. !1865
- Add basic auth support to clone analyzer. !1866
- Praefect ping-node must verify storage locations are served. !1881
- Support partial clones with SSH transports. !1893
- Add storage name to healthcheck error log. !1934
- Always use V2 tokens in gitaly auth client.

### Other (7 changes)

- Bypass praefect server in tests that check the error message. !1799
- Set default concurrency limit for ReplicateRepository. !1822
- Fix example Praefect config file for virtual storage changes. !1856
- Add correlation ID to Praefect replication jobs. !1869
- Remove dependency on the outdated golang.org/x/net package. !1882
- Upgrade parser gem to v2.7.0.4. !1935
- Simplify loading of required Ruby files. !1942


## 12.8.10

### Security (1 change)

- gems: Upgrade nokogiri to > 1.10.7. !2127


## 12.8.9

- No changes.

## 12.8.7

- No changes.

## 12.8.6

- No changes.

## 12.8.5

- No changes.

## 12.8.4

- No changes.

## 12.8.3

- No changes.

## 12.8.2

- No changes.

## 12.8.1

- No changes.

## 12.8.0

- No changes.
### Added (1 change)

- Add praefect client prometheus interceptor. !1836

### Other (1 change)

- Praefect sub-commands: avoid garbage in logs. !1819


## 12.7.9

- No changes.

## 12.7.8

- No changes.

## 12.7.7

- No changes.

## v1.87.0

#### Added
- Logging of repository objects statistics after repack
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1800
- Support of basic interaction with Postgres database
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1768
- Pass Ruby-specific feature flags to the Ruby server
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1818

#### Changed
- Wire up coordinator to use node manager
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1806
- Enable toggling on the node manager through a config value
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1803
- Praefect replicator to mark job as failed and retry failed jobs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1804

#### Fixed
- Calculate praefect replication latency using seconds with float64
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1820
- Pass in virtual storage name to repl manager
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1831
- Return full GitCommit message as part of FindLocalBranch gRPC response.
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1827
- Use token auth for Praefect connection checker
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1816
- Decode user info for request authorization
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1805

#### Other
- Remove protocol v2 feature flag
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1841

#### Performance
- Skip the batch check for commits and tags
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1812

## v1.86.0

#### Added
- Support FindCommitsRequest with order (--topo-order)
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1791
- Add Node manager
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1779

#### Changed
- simplify praefect routing to primary and replication jobs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1760
- PostReceiveHook: add support for Git push options
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1756

#### Fixed
- Incorrect changelogs should be caught by Danger
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1811
- Cache and reuse client connection in ReplicateRepository
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1801
- Fix cache walker to only walk each path once
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1769
- UpdateRemoteMirror: handle large number of branches
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1745

#### Other
- Update activesupport, gitlab-labkit, and other Ruby dependencies
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1794
- Add deadline_type prometheus label
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1737
- Use golangci-lint for static code analysis
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1722
- Remove unused rubyserver in structs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1807
- Reenable git wire protocol v2 behind feature flag
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1797
- Add grpc tag interceptor
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1795

#### Security
- Validate bad branches for UserRebase and UserRebaseConfirmable
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1735

## v1.85.0

#### Deprecated
- Revert branch field removal in UserSquashRequest message for RPC operations.UserSquash
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1792

#### Other
- Add praefect as a transparent pass through for tests
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1736

## v1.84.0

#### Added
- Use core delta islands to increase opportunity of pack reuse
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1775
- Feature flag: look up commits without batch-check
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1711

#### Fixed
- Include stderr in output of worktree operations
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1787
- Fix gitaly-hooks check command which was broken due to an incorrect implemention
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1761
- Fix squash when target repository has a renamed file
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1786
- Make parent directories before snapshot replication
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1767

#### Other
- Update rouge gem to 3.15.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1778
- Register praefect server for grpc prom metrics
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1782
- Update tzinfo gem to v1.2.6
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1762
- Explicitly mention gitaly-ruby in error messages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1776
- Remove revision from GetAllLFSPointers request
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1771

#### Security
- Do not log entire node struct because it includes tokens
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1766
- Replace CommandWithoutRepo usage with safe version
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1783
- Bump go-yaml and the rack gem dependencies
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

## v1.83.0

#### Other
- Refine telemetry for cache walker
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1759

## v1.82.0

#### Added
- Praefect subcommand for checking node connections
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1700

## v1.81.0

#### Added
- Allow git_push_options to be passed to UserRebaseConfirmable RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1748
- Add sql-migrate subcommand to Praefect
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1738

#### Changed
- Add snapshot replication to ReplicateRepository
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1717
- Add exponential backoff to replication manager
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1746

#### Fixed
- Fix middleware to stop panicking from bad requests
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1747
- Call client.Dial in ClientConnection helper
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1749

#### Other
- Log statistics of pack re-use on clone and fetch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1743
- Change signature of hook RPCs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1741
- Update loofah and crass gems to match GitLab CE/EE
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1752

#### Security
- Do not leak sensitive urls
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1710

## v1.80.0

#### Fixed
- Fix DiskStatistics on OpenBSD
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1728
  Contributed by bitgestalt

#### Other
- File walker deletes empty directories
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1721
- FetchIntoObjectPool: log pool object and ref directory sizes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1614

#### Performance
- Add hook rpcs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1686

## v1.79.0

#### Changed
- praefect replicator links object pool if it exists
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1718
- Use configurable buckets for praefect replication metrics
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1719

#### Fixed
- Strip invalid characters in signatures on commit
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1709

#### Other
- Fix order of branches in git diff when preparing sparse checkout rebase
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1716
- Fix call to testhelper.TempDir
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1730
- Upgrade Nokogiri to 1.10.7
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1731
- Bump Ruby to 2.6.5
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1727

## v1.78.0


## vv1.78.0

#### Changed
- Use ReplicateRepository in replicator
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1713

## v1.77.0

#### Changed
- Add author to FindCommits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1702
- Remove get_tag_messages_go feature flag
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1698

#### Fixed
- Add back feature flag for cache invalidation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1706

#### Other
- Sync info attributes in ReplicateRepository
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1693

#### Security
- Upgrade Rugged to v0.28.4.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1701

## v1.76.0

#### Added
- ReplicateRepository RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1605
- add signature type to GitCommit
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1635
  Contributed by bufferoverflow

#### Deprecated
- PreFetch: remove unused RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1675
- Upgrade to Git 2.24
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1653

#### Fixed
- Fix forking with custom CA in RPC CreateFork
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1658

#### Other
- Start up log messages are now using structured logging too
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1674
- Log an error if praefect's server info fails to connect to a node
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1651
- Update msgpack-ruby to v1.3.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1677
- Log all diskcache state changes and stream access
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1682
- Move prometheus config to its own package
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1676
- Remove ruby script approach to GetAllLFSPointers
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1695
- StreamDirector returns StreamParams
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1679
- Move bootstrap env vars into bootstrap constructor
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1670

#### Performance
- Filter collection of SHAs which has signatures and return those SHAs: go implementation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1672

#### Security
- Update loofah gem to v2.3.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1678

## v1.75.0

#### Changed
- Praefect multiple virtual storage
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1606

#### Fixed
- Gitaly feature flags are broken: convert underscores to dashes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1659
- Allow internal fetches to see all hidden references
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1640
- SSHUpload{Pack,Archive}: fix timeout tests
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1664
- Restore gitaly_connections_total prometheus metric
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1657

#### Other
- Add labkit healthcheck
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1646
- Configure logging as early as possible
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1666
- Use internal socket dir for internal gitaly socket
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1642
- Leverage the bootstrap package to support Praefect zero downtime deploys
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1638

#### Security
- Limit the negotiation phase for certain Gitaly RPCs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

## v1.74.0

#### Added
- Add DiskStatistics grpc method to ServerService
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1620
  Contributed by maxmati
- Add Praefect service
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1628

#### Fixed
- UpdateRemoteMirror: fix default branch resolution
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1641

#### Other
- Refactor datastore to its own package
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1627
- Validate that hook files are reachable and executable
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1574
- Upgrade charlock_holmes Ruby gem to v0.7.7
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1647

## v1.73.0

#### Added
- Label storage name in all storage scoped requests
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1617
  Contributed by maxmati
- Allow socket dir for Gitaly-Ruby to be configured
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1184

#### Fixed
- Fix client keep alive for all network types
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1637

#### Other
- Move over to Labkit Healthcheck endpoint
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1448

## v1.72.0

#### Added
- Propagate repository removal to Praefect backends
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1586
- Add GetObjectPool RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1583
- Add ReplicateRepository protocol
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1601

#### Fixed
- Fix CreateBundle scope and target_repository_field
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1597

#### Performance
- Changed files used for sparse checkout
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1611

## v1.71.0

#### Added
- Fishy config log warning
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1570
- Add gRPC intercept loggers to Praefect
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1573
- Add check subcommand in gitaly-hooks
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1587

#### Deprecated
- Remove StorageService
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1544

#### Other
- Count dangling refs before/after FetchIntoObjectPool
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1585
- Count v2 auth error return paths
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1568
- Upgrade gRPC Ruby library to v1.24.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/
- Adding sentry config to praefect
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1546
- Add HookService RPCs and methods
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1553
- Update rails-html-sanitizer and loofah gems in gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1566

## v1.70.0


## vv1.70.0

#### Added
- Lower gRPC server inactivity ping timeout
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1294

#### Other
- Remove RepositoryService WriteConfig
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1558
- Refactor praefect server tests
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1554

## v1.69.0

#### Fixed
- Fix praefect server info to include git version, server version
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1550

#### Other
- Enable second repository to have its storage re-written
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1548

## v1.68.0

#### Added
- Add virtual storage name to praefect config
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1525
- Support health checks in Praefect
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1541

#### Changed
- Provide specifics about why a cherry-pick or revert has failed
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1543

#### Other
- Upgrade GRPC to 1.24.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1539

## v1.67.0

#### Added
- Adding auth to praefect
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1517
- Don't panic, go retry
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1523
- Allow praefect to handle ServerInfoRequest
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1527

#### Fixed
- Support configurable Git config search path for Rugged
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1526

#### Other
- Create go binary to execute hooks
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1328

#### Performance
- Allow upload pack filters with feature flag
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1411

## v1.66.0

#### Changed
- Include file count and bytes in CommitLanguage response
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1482

#### Fixed
- Leave stderr alone when passed into command
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1456
- Fix cache invalidator for Create* RPCs and health checks
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1494
- Set split index to false
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1521
- Ensure temp dir exists when removing a repository
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1509

#### Other
- Use safe command in HasLocalBranches
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1499
- Explicitly designate primary and replica nodes in praefect config
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1483
- Nested command for DSL
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1498
- Use SafeCmd in WriteRef
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1506
- Move cache state files to +gitaly directory
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1515

#### Security
- ConfigPair option for DSL
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1507

## v1.65.0

#### Fixed
- Replicator fixes from demo
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1487
- Prevent nil panics in housekeeping.Perform
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1492
- Upgrade Rouge to v3.11.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1493

#### Other
- Git Command DSL
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1476
- Measure replication latency
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1481

## v1.64.0

#### Added
- Confirm checksums after replication
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1479
- FindCommits/CommitCounts: Add support for `first_parent` parameter
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1463
  Contributed by jhenkens

#### Other
- Update Rouge to v3.10.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1475

#### Security
- Add dedicated CI job for deprecation warnings
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1480

## v1.63.0

#### Added
- Set permission of attributes file to `0644`
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1466
  Contributed by njkevlani

#### Other
- Add RemoveRepository RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1470

#### Performance
- FetchIntoObjectPool: pack refs after fetch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1464

## v1.62.0

#### Added
- Praefect Realtime replication
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1423

#### Fixed
- GetAllLFSPointers: use binmode in inline Ruby script
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1454
- Fix catfile metrics counting bug
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1462
- main: start ruby server after opening network listeners
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1455

#### Other
- Create parent directory in RenameNamespace
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1452
- Update Ruby gitlab-labkit to 0.5.2
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1453
- Update logging statements to leverage structured logging better
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

#### Performance
- Modify GetBlobs RPC to return type
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1445

## v1.61.0

#### Security
- Do not follow redirect when cloning repo from URL
  https://gitlab.com/gitlab-org/gitaly/merge_requests/
- Add http.followRedirects directive to `git fetch` command
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

## v1.60.0

#### Added
- Praefect data model changes with replication
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1399
- Include process PID in log messages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1422

#### Changed
- Make it easier to add new kinds of internal post receive messages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1410

#### Fixed
- Validate commitIDs parameter for get_commit_signatures RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1428

#### Other
- Update ffi gem to 1.11.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/
- Add back old env vars for backwards compatibility
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1367
- Pass through GOPATH to control cache location
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1436

#### Performance
- Port GetAllLFSPointers to go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1414
- FindTag RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1409
- Disk cache object directory initial clear
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1424

#### Security
- Upgrade Rugged to 0.28.3.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1427

## v1.59.0

#### Added
- Port GetCommitSignatures to Go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1283

#### Other
- Update gitlab-labkit to 0.4.2
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1412

#### Performance
- Enable disk cache invalidator gRPC interceptors
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1373

#### Security
- Bump nokogiri to 1.10.4
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1415
- Fix FindCommits flag injection exploit
  https://gitlab.com/gitlab-org/gitaly/merge_requests/31

## v1.58.0

#### Fixed
- Properly clean up worktrees after commit operations
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1383

## v1.57.0

#### Fixed
- Fix Praefect's mock service protobuf Go-stub file generation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1400

#### Performance
- Add gRPC method scope to protoregistry
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1397

## v1.56.0

#### Added
- Add capability to replace certain frames in a stream
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1382
- Gitaly proto method request factories
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1375

#### Fixed
- Unable to extract target repo when nested in oneOf field
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1377
- Update Rugged to 0.28.2
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1384

#### Other
- Remove rescue of Gitlab::Git::CommitError at UserMergeToRef RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1372
- Handle failover by catching SIGUSR1 signal
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1346
- Update rouge to v3.7.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1387
- Update msgpack to 1.3.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1381
- Upgrade Ruby gitaly-proto to 1.37.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1374

#### Performance
- Unary gRPC interceptor for cache invalidation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1371

## 1.55.0

### Fixed (1 change)

- Remove context from safe file. !1369

### Added (5 changes)

- Cache invalidation via gRPC interceptor. !1268
- Add CloneFromPool RPC. !1301
- Add support for Git 2.22. !1359
- Disk cache object walker. !1362
- Add ListCommitsByRefName RPC. !1365

### Other (1 change)

- Remove catfile cache feature flag. !1357


## 1.54.1

- No changes.

## 1.54.0 (2019-07-22)

- No changes.

## v1.53.0

#### Added
- Expose the Praefect server version
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1358

#### Changed
- Support start_sha parameter in UserCommitFiles
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1308

## v1.52.0

#### Other
- Do not add linked repos as remotes on pool repository
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1356
- Upgrade rouge to 3.5.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1355

## v1.51.0

#### Changed
- Add support first_parent_ref in UserMergeToRef
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1210

#### Fixed
- More informative error states for missing pages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1340
- Use guard in fetch_legacy_config
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1345

#### Other
- Add HTTP clone analyzer
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1338

## v1.50.0

#### Added
- Use datastore to store the primary node
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1335

#### Fixed
- Fix default lookup of global custom hooks
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1336

#### Other
- Add filesystem metadata file on startup
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1289
- Pass down gitlab-shell log config through env vars
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1293
- Remove duplication of receive-pack config
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1332
- Update Prometheus client library
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1329

#### Performance
- Hide object pools .have refs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1323

## v1.49.0

#### Fixed
- Cleanup RPC now uses proper prefix for worktree clean up
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1325
- GetRawChanges RPC uses both string and byte path fields
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1207

#### Other
- Add lock file package for atomic file writes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1298

#### Performance
- Maintain pool packfiles
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1316

## v1.48.0

#### Fixed
- Fix praefect not listening to the correct socket path
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1313

#### Other
- Skip hooks for UserMergeToRef RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1312

## v1.47.0

#### Changed
- Remove member bitmaps when linking to objectpool
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1311
- Un-dangle dangling objects in object pools
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1297

#### Fixed
- Fix ignored registerNode error
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1307
- Fix Prometheus metric naming errors
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1292
- Cast FsStat syscall to int64
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1306

#### Other
- Upgrade protobuf, prometheus and logrus
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1290
- Replace govendor with 'go mod'
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1286

#### Removed
- Remove ruby code to create a repository
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1302

## v1.46.0

#### Added
- Add GetObjectDirectorySize RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1263

#### Changed
- Make catfile cache size configurable
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1271

#### Fixed
- Wait for all the socket to terminate during a graceful restart
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1190

#### Performance
- Enable bitmap hash cache
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1282

## v1.45.0

#### Performance
- Enable splitIndex for repositories in GarbageCollect rpc
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1247

#### Security
- Fix GetArchive injection vulnerability
  https://gitlab.com/gitlab-org/gitaly/merge_requests/26

## v1.44.0

#### Added
- Expose the FileSystem name on the storage info
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1261

#### Changed
- DisconnectGitAlternates: bail out more often
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1266

#### Fixed
- Created repository directories have FileMode 0770
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1274
- Fix UserRebaseConfirmable not sending PreReceiveError and GitError responses to client
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1270
- Fix stderr log writer
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1275

#### Other
- Speed up 'make assemble' using rsync
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1272

## v1.43.0

#### Added
- Stop symlinking hooks on repository creation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1052
- Replication logic
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1219
- gRPC proxy stream peeking capability
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1260
- Introduce ps package
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1258

#### Changed
- Remove delta island feature flags
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1267

#### Fixed
- Fix class name of Labkit tracing inteceptor
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1269
- Fix replication job state changing
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1236
- Remove path field in ListLastCommitsForTree response
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1240
- Check if PID belongs to Gitaly before adopting an existing process
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1249

#### Other
- Absorb grpc-proxy into Gitaly
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1248
- Add git2go dependency
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1061
  Contributed by maxmati
- Upgrade Rubocop to 0.69.0 with other dependencies
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1250
- LabKit integration with Gitaly-Ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1083

#### Performance
- Fix catfile N+1 in ListLastCommitsForTree
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1253
- Use --perl-regexp for code search
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1241
- Upgrade to Ruby 2.6
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1228
- Port repository creation to Golang
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1245

## v1.42.0

#### Other
- Use simpler data structure for cat-file cache
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1233

## v1.41.0

#### Added
- Implement the ApplyBfgObjectMapStream RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1199

## v1.40.0

#### Fixed
- Do not close the TTL channel twice
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1235

## v1.39.0

#### Added
- Add option to add Sentry environment
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1216
  Contributed by Roger Meier

#### Fixed
- Fix CacheItem pointer in cache
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1234

## v1.38.0

#### Added
- Add cache for batch files
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1203

#### Other
- Upgrade Rubocop to 0.68.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1229

## v1.37.0

#### Added
- Add DisconnectGitAlternates RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1141

## v1.36.0

#### Added
- adding ProtoRegistry
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1188
- Adding FetchIntoObjectPool RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1172
- Add new two-step UserRebaseConfirmable RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1208

#### Fixed
- Include stderr in err returned by git.Command Wait()
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1167
- Use 3-way merge for squashing commits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1214
- Close logrus writer when command finishes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1225

#### Other
- Bump Ruby bundler version to 1.17.3
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1215
- Upgrade Ruby gRPC 1.19.0 and protobuf to 3.7.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1066
- Ensure pool exists in LinkRepositoryToObjectPool rpc
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1222
- Update FetchRemote ruby to write http auth as well as add remote
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1126

#### Performance
- GarbageCollect RPC writes commit graph and enables via config
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1218

#### Security
- Bump Nokogiri to 1.10.3
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1217
- Loosen regex for exception sanitization
  https://gitlab.com/gitlab-org/gitaly/merge_requests/25

## v1.35.1

The v1.35.1 tag points to a release that was made on the wrong branch, please
ignore.

## v1.35.0

#### Added
- Return path data in ListLastCommitsForTree
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1168

## v1.34.0

#### Added
- Add PackRefs RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1161
- Implement ListRemotes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1019
- Test and require Git 2.21
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1205
- Add WikiListPages RPC call
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1194

#### Fixed
- Cleanup RPC prunes disconnected work trees
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1189
- Fix FindAllTags to dereference tags that point to other tags
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1193

#### Other
- Datastore pattern for replication jobs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1147
- Remove find all tags ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1163
- Delete SSH frontend code from ruby/gitlab-shell
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1179

## v1.33.0

#### Added
- Zero downtime deployment
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1133

#### Changed
- Move gitlab-shell out of ruby/vendor
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1173

#### Other
- Bump Ruby gitaly-proto to v1.19.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1186
- Bump sentry-raven to 2.9.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1183
- Bump gitlab-markup to 1.7.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1182

#### Performance
- Improve GetBlobs performance for fetching lots of files
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1165

#### Security
- Bump activesupport to 5.0.2.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1185

## v1.32.0

#### Fixed
- Remove test dependency in main binaries
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1171

#### Other
- Vendor gitlab-shell at 433cc96551a6d1f1621f9e10
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1175

## v1.31.0

#### Added
- Accept Path option for GetArchive RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1142

#### Changed
- UnlinkRepositoryFromObjectPool: stop removing objects/info/alternates
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1151

#### Other
- Always use overlay2 storage driver on Docker build
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1148
  Contributed by Takuya Noguchi
- Remove unused Ruby implementation of GetRawChanges
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1169
- Remove Ruby implementation of remove remote
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1164

#### Removed
- Drop support for Golang 1.10
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1149

## v1.30.0

#### Added
- WikiGetAllPages RPC - Add params for sorting
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1081

#### Changed
- Keep origin remote and refs when creating an object pool
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1136

#### Fixed
- Bump github-linguist to 6.4.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1153
- Fix too lenient ref wildcard matcher
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1158

#### Other
- Bump Rugged to 0.28.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/154
- Remove FindAllTags feature flag
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1155

#### Performance
- Use delta islands in RepackFull and GarbageCollect
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1110

## v1.29.0

#### Fixed
- FindAllTags: Handle edge case of annotated tags without messages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1134
- Fix "bytes written" count in pktline.WriteString
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1129
- Prevent clobbering existing Git alternates
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1132
- Revert !1088 "Stop using gitlab-shell hooks -- but keep using gitlab-shell config"
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1117

#### Other
- Introduce text.ChompBytes helper
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1144
- Re-apply MR 1088 (Git hooks change)
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1130

## v1.28.0

Should not be used as it [will break gitlab-rails](https://gitlab.com/gitlab-org/gitlab-ce/issues/58855).

#### Changed
- RenameNamespace RPC creates parent directories for the new location
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1090

## v1.27.0

#### Added
- Support socket paths for praefect
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1115

#### Fixed
- Fix bug in FindAllTags when commit shas are used as tag names
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1119

## v1.26.0

#### Added
- PreFetch RPC: to optimize a full fetch by doing a local clone from the fork source
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1073

## v1.25.0

#### Added
- Add prometheus listener to Praefect
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1108

#### Changed
- Stop using gitlab-shell hooks -- but keep using gitlab-shell config
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1088

#### Fixed
- Fix undefined logger panicing
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1114

#### Other
- Stop running tests on Ruby 2.4
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1113
- Add feature flag for FindAllTags
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1106

#### Performance
- Rewrite remove remote in go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1051
  Contributed by maxmati

## v1.24.0

#### Added
- Accept Force option for UserCommitFiles to overwrite branch on commit
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1077

#### Fixed
- Fix missing SEE_DOC constant in Danger
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1109

#### Other
- Increase Praefect unit test coverage
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1103
- Use GitLab for License management
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1076

## v1.23.0

#### Added
- Allow debugging ruby tests with pry
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1102
- Create Praefect binary for proxy server execution
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1068

#### Fixed
- Try to resolve flaky TestRemoval balancer test
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1094
- Bump Rugged to 0.28.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

#### Other
- Remove unused Ruby implementation for CommitStats
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1092
- GitalyBot will apply labels to merge requests
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1105
- Remove non-chunked code path for SearchFilesByContent
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1100
- Remove ruby implementation of find commits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1099
- Update Gitaly-Proto with protobuf go compiler 1.2
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1084
- removing deprecated ruby write-ref
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1098

## v1.22.0

#### Fixed
- Pass GL_PROTOCOL and GL_REPOSITORY to update hook
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1082

#### Other
- Support distributed tracing in child processes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1085

#### Removed
- Removing find_branch ruby implementation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1096

## v1.21.0

#### Added
- Support merge ref writing (without merging to target branch)
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1057

#### Fixed
- Use VERSION file to detect version as fallback
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1056
- Fix GetSnapshot RPC to work with repos with object pools
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1045

#### Other
- Remove another test that exercises gogit feature flag
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1086

#### Performance
- Rewrite FindAllTags in Go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1036
- Reimplement DeleteRefs in Go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1069

## v1.20.0

#### Fixed
- Bring back a custom dialer for Gitaly Ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1072

#### Other
- Initial design document for High Availability
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1058
- Reverse proxy pass thru for HA
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1064

## v1.19.1

#### Fixed
- Use empty tree if initial commit
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1075

## v1.19.0

#### Fixed
- Return blank checksum for git repositories with only invalid refs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1065

#### Other
- Use chunker in GetRawChanges
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1043

## v1.18.0

#### Other
- Make clear there is no []byte reuse bug in SearchFilesByContent
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1055
- Use chunker in FindCommits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1059
- Statically link jaeger into Gitaly by default
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1063

## v1.17.0

#### Other
- Add glProjectPath to logs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1049
- Switch from commitsSender to chunker
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1060

#### Security
- Disable git protocol v2 temporarily
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

## v1.16.0


## v1.15.0

#### Added
- Support rbtrace and ObjectSpace via environment flags
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1046

#### Changed
- Add CountDivergingCommits RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1023

#### Fixed
- Add chunking support to SearchFilesByContent RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1015
- Avoid unsafe use of scanner.Bytes() in ref name RPC's
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1054
- Fix tests that used long unix socket paths
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1039

#### Other
- Use chunker for ListDirectories RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1042
- Stop using nil internally to signal "commit not found"
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1050
- Refactor refnames RPC's to use chunker
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1041

#### Performance
- Rewrite CommitStats in Go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1048

## v1.14.1

#### Security
- Disable git protocol v2 temporarily
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

## v1.14.0

#### Fixed
- Ensure that we kill ruby Gitlab::Git::Popen reader threads
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1040

#### Other
- Vendor gitlab-shell at 6c5b195353a632095d7f6
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1037

## v1.13.0

#### Fixed
- Fix 503 errors when Git outputs warnings to stderr
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1024
- Fix regression for https_proxy and unix socket connections
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1032
- Fix flaky rebase test
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1028
- Rewrite GetRawChanges and fix quoting bug
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1026
- Fix logging of RenameNamespace RPC parameters
  https://gitlab.com/gitlab-org/gitaly/merge_requests/847

#### Other
- Small refactors to gitaly/client
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1034
- Prepare for vendoring gitlab-shell hooks
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1020
- Replace golang.org/x/net/context with context package
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1038
- Migrate writeref from using the ruby implementation to go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1008
  Contributed by johncai
- Switch from honnef.co/go/tools/megacheck to staticcheck
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1021
- Add distributed tracing support with LabKit
  https://gitlab.com/gitlab-org/gitaly/merge_requests/976
- Simplify error wrapping in service/ref
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1009
- Remove dummy RequestStore
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1007
- Simplify error handling in ssh package
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1029
- Add response chunker abstraction
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1031
- Use go implementation of FindCommits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1025
- Rewrite get commit message
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1012
- Update docs about monitoring and README
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1016
  Contributed by Takuya Noguchi
- Remove unused Ruby rebase/squash code
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1033

## v1.12.2

#### Security
- Disable git protocol v2 temporarily
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

## v1.12.1

#### Fixed
- Fix flaky rebase test
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1028
- Fix regression for https_proxy and unix socket connections
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1032

## v1.12.0

#### Fixed
- Fix wildcard protected branches not working with remote mirrors
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1006

## v1.11.0

#### Fixed
- Fix incorrect tree entries retrieved with directories that have curly braces
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1013
- Deduplicate CA in gitaly tls
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1005

## v1.10.0

#### Added
- Allow repositories to be reduplicated
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1003

#### Fixed
- Add GIT_DIR to hook environment
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1001

#### Performance
- Re-implemented FindBranch in Go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/981

## v1.9.0

#### Changed
- Improve Linking and Unlink object pools RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1000

#### Other
- Fix tests failing due to test-repo change
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1004

## v1.8.0

#### Other
- Log correlation_id field in structured logging output
  https://gitlab.com/gitlab-org/gitaly/merge_requests/995
- Add explicit null byte check in internal/command.New
  https://gitlab.com/gitlab-org/gitaly/merge_requests/997
- README cleanup
  https://gitlab.com/gitlab-org/gitaly/merge_requests/996

## v1.7.2

#### Other
- Fix tests failing due to test-repo change
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1004

#### Security
- Disable git protocol v2 temporarily
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

## v1.7.1

#### Other
- Log correlation_id field in structured logging output
  https://gitlab.com/gitlab-org/gitaly/merge_requests/995

## v1.7.0

#### Added
- Add an RPC that allows repository size to be reduced by bulk-removing internal references
  https://gitlab.com/gitlab-org/gitaly/merge_requests/990

## v1.6.0

#### Other
- Clean up invalid keep-around refs when performing housekeeping
  https://gitlab.com/gitlab-org/gitaly/merge_requests/992

## v1.5.0

#### Added
- Add tls configuration to gitaly golang server
  https://gitlab.com/gitlab-org/gitaly/merge_requests/932

#### Fixed
- Fix TLS client code on macOS
  https://gitlab.com/gitlab-org/gitaly/merge_requests/994

#### Other
- Update to latest goimports formatting
  https://gitlab.com/gitlab-org/gitaly/merge_requests/993

## v1.4.0

#### Added
- Link and Unlink RPCs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/986

## v1.3.0

#### Other
- Remove unused bridge_exceptions method
  https://gitlab.com/gitlab-org/gitaly/merge_requests/987
- Clean up process documentation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/984

## v1.2.0

#### Added
- Upgrade proto to v1.2
  https://gitlab.com/gitlab-org/gitaly/merge_requests/985
- Allow moved files to infer their content based on the source
  https://gitlab.com/gitlab-org/gitaly/merge_requests/980

#### Other
- Add connectivity tests
  https://gitlab.com/gitlab-org/gitaly/merge_requests/968

## v1.1.0

#### Other
- Remove grpc dependency from catfile
  https://gitlab.com/gitlab-org/gitaly/merge_requests/983
- Don't use rugged when calling write-ref
  https://gitlab.com/gitlab-org/gitaly/merge_requests/982

## v1.0.0

#### Added
- Add gitaly-debug production debugging tool
  https://gitlab.com/gitlab-org/gitaly/merge_requests/967

#### Fixed
- Bump gitlab-markup to 1.6.5
  https://gitlab.com/gitlab-org/gitaly/merge_requests/975
- Fix to reallow tcp URLs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/974

#### Other
- Upgrade minimum required Git version to 2.18.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/958
- Bump tzinfo to 1.2.5
  https://gitlab.com/gitlab-org/gitaly/merge_requests/977
- Bump activesupport gem to 5.0.7
  https://gitlab.com/gitlab-org/gitaly/merge_requests/978
- Propagate correlation-ids in from upstream services and out to Gitaly-Ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/970

#### Security
- Bump nokogiri to 1.8.5
  https://gitlab.com/gitlab-org/gitaly/merge_requests/979

## v0.133.0

#### Other
- Upgrade gRPC-go from v1.9.1 to v1.16
  https://gitlab.com/gitlab-org/gitaly/merge_requests/972

## v0.132.0

#### Other
- Upgrade to Ruby 2.5.3
  https://gitlab.com/gitlab-org/gitaly/merge_requests/942
- Remove dead code post 10.8
  https://gitlab.com/gitlab-org/gitaly/merge_requests/964

## v0.131.0

#### Fixed
- Fixed bug with wiki operations enumerator when content nil
  https://gitlab.com/gitlab-org/gitaly/merge_requests/962

## v0.130.0

#### Added
- Support SSH credentials for push mirroring
  https://gitlab.com/gitlab-org/gitaly/merge_requests/959

## v0.129.1

#### Other
- Fix tests failing due to test-repo change
  https://gitlab.com/gitlab-org/gitaly/merge_requests/1004

#### Security
- Disable git protocol v2 temporarily
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

## v0.129.0

#### Added
- Add submodule reference update operation in the repository
  https://gitlab.com/gitlab-org/gitaly/merge_requests/936

#### Fixed
- Improve wiki hook error message
  https://gitlab.com/gitlab-org/gitaly/merge_requests/963
- Fix encoding bug in User{Create,Delete}Tag
  https://gitlab.com/gitlab-org/gitaly/merge_requests/952

#### Other
- Expand Gitlab::Git::Repository unit specs with examples from rails
  https://gitlab.com/gitlab-org/gitaly/merge_requests/945
- Update vendoring
  https://gitlab.com/gitlab-org/gitaly/merge_requests/954

## v0.128.0

#### Fixed
- Fix incorrect committer when committing patches
  https://gitlab.com/gitlab-org/gitaly/merge_requests/947
- Fix makefile 'find ruby/vendor' bug
  https://gitlab.com/gitlab-org/gitaly/merge_requests/946

## v0.127.0

#### Added
- Make git hooks self healing
  https://gitlab.com/gitlab-org/gitaly/merge_requests/886
- Add an endpoint to apply patches to a branch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/926

#### Fixed
- Use $(MAKE) when re-invoking make
  https://gitlab.com/gitlab-org/gitaly/merge_requests/933

#### Other
- Bump google-protobuf gem to 3.6.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/941
- Bump Rouge gem to 3.3.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/943
- Upgrade Ruby version to 2.4.5
  https://gitlab.com/gitlab-org/gitaly/merge_requests/944

## v0.126.0

#### Added
- Add support for closing Rugged/libgit2 file descriptors
  https://gitlab.com/gitlab-org/gitaly/merge_requests/903

#### Changed
- Require storage directories to exist at startup
  https://gitlab.com/gitlab-org/gitaly/merge_requests/675

#### Fixed
- Don't confuse govendor license with ruby gem .go files
  https://gitlab.com/gitlab-org/gitaly/merge_requests/935
- Rspec and bundler setup fixes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/901
- Fix git protocol prometheus metrics
  https://gitlab.com/gitlab-org/gitaly/merge_requests/908
- Fix order in config.toml.example
  https://gitlab.com/gitlab-org/gitaly/merge_requests/923

#### Other
- Standardize git command invocation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/915
- Update grpc to v1.15.x in gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/918
- Add package tests for internal/git/pktline
  https://gitlab.com/gitlab-org/gitaly/merge_requests/909
- Make Makefile more predictable by bootstrapping
  https://gitlab.com/gitlab-org/gitaly/merge_requests/913
- Force english output on git commands
  https://gitlab.com/gitlab-org/gitaly/merge_requests/898
- Restore notice check
  https://gitlab.com/gitlab-org/gitaly/merge_requests/902
- Prevent stale packed-refs file when Gitaly is running on top of NFS
  https://gitlab.com/gitlab-org/gitaly/merge_requests/924

#### Performance
- Update Prometheus vendoring
  https://gitlab.com/gitlab-org/gitaly/merge_requests/922
- Free Rugged open file descriptors in gRPC middleware
  https://gitlab.com/gitlab-org/gitaly/merge_requests/911

#### Removed
- Remove deprecated methods
  https://gitlab.com/gitlab-org/gitaly/merge_requests/910

#### Security
- Bump Rugged to 0.27.5 for security fixes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/907

## v0.125.0

#### Added
- Support Git protocol v2
  https://gitlab.com/gitlab-org/gitaly/merge_requests/844

#### Other
- Remove test case that exercises gogit feature flag
  https://gitlab.com/gitlab-org/gitaly/merge_requests/899

## v0.124.0

#### Deprecated
- Remove support for Go 1.9
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

#### Fixed
- Fix panic in git pktline splitter
  https://gitlab.com/gitlab-org/gitaly/merge_requests/893

#### Other
- Rename gitaly proto import to gitalypb
  https://gitlab.com/gitlab-org/gitaly/merge_requests/895

## v0.123.0

#### Added
- Add ListLastCommitsForTree to retrieve the last commits for every entry in the current path
  https://gitlab.com/gitlab-org/gitaly/merge_requests/881

#### Other
- Wait for gitaly to boot in rspec integration tests
  https://gitlab.com/gitlab-org/gitaly/merge_requests/890

## v0.122.0

#### Added
- Implements CHMOD action of UserCommitFiles API
  https://gitlab.com/gitlab-org/gitaly/merge_requests/884
  Contributed by Jacopo Beschi @jacopo-beschi

#### Changed
- Use CommitDiffRequest.MaxPatchBytes instead of hardcoded limit for single diff patches
  https://gitlab.com/gitlab-org/gitaly/merge_requests/880
- Implement new credentials scheme on gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/873

#### Fixed
- Export HTTP proxy environment variables to Gitaly
  https://gitlab.com/gitlab-org/gitaly/merge_requests/885

#### Security
- Sanitize sentry events' logentry messages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/

## v0.121.0

#### Changed
- CalculateChecksum: Include keep-around and other references in the checksum calculation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/731

#### Other
- Stop vendoring Gitlab::Git
  https://gitlab.com/gitlab-org/gitaly/merge_requests/883

## v0.120.0

#### Added
- Server implementation ListDirectories
  https://gitlab.com/gitlab-org/gitaly/merge_requests/868

#### Changed
- Return old and new modes on RawChanges
  https://gitlab.com/gitlab-org/gitaly/merge_requests/878

#### Other
- Allow server to receive an hmac token with the client timestamp for auth
  https://gitlab.com/gitlab-org/gitaly/merge_requests/872

## v0.119.0

#### Added
- Add server implementation for FindRemoteRootRef
  https://gitlab.com/gitlab-org/gitaly/merge_requests/874

#### Changed
- Allow merge base to receive more than 2 revisions
  https://gitlab.com/gitlab-org/gitaly/merge_requests/869
- Stop vendoring some Gitlab::Git::* classes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/865

#### Fixed
- Support custom_hooks being a symlink
  https://gitlab.com/gitlab-org/gitaly/merge_requests/871
- Prune large patches by default when enforcing limits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/858
- Fix diffs being collapsed unnecessarily
  https://gitlab.com/gitlab-org/gitaly/merge_requests/854
- Remove stale HEAD.lock if it exists
  https://gitlab.com/gitlab-org/gitaly/merge_requests/861
- Fix patch size calculations to not include headers
  https://gitlab.com/gitlab-org/gitaly/merge_requests/859

#### Other
- Vendor Gitlab::Git at c87ca832263
  https://gitlab.com/gitlab-org/gitaly/merge_requests/860
- Bump gitaly-proto to 0.112.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/857

#### Security
- Bump rugged to 0.27.4 for security fixes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/856
- Update the sanitize gem to at least 4.6.6
  https://gitlab.com/gitlab-org/gitaly/merge_requests/876
- Bump rouge to 3.2.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/862

## v0.118.0

#### Added
- Add ability to support custom options for git-receive-pack
  https://gitlab.com/gitlab-org/gitaly/merge_requests/834

## v0.117.2

#### Fixed
- Fix diffs being collapsed unnecessarily
  https://gitlab.com/gitlab-org/gitaly/merge_requests/854
- Fix patch size calculations to not include headers
  https://gitlab.com/gitlab-org/gitaly/merge_requests/859
- Prune large patches by default when enforcing limits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/858

## v0.117.1

#### Security
- Bump rouge to 3.2.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/862
- Bump rugged to 0.27.4 for security fixes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/856

## v0.117.0

#### Performance
- Only load Wiki formatted data upon request
  https://gitlab.com/gitlab-org/gitaly/merge_requests/839

## v0.116.0

#### Added
- Add ListNewBlobs RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/849

## v0.115.0

#### Added
- Implement DiffService.DiffStats RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/808
- Update gitaly-proto to 0.109.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/843

#### Changed
- Stop vendoring Gitlab::VersionInfo
  https://gitlab.com/gitlab-org/gitaly/merge_requests/840

#### Fixed
- Check errors and fix chunking in ListNewCommits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/852
- Fix reStructuredText not working on Gitaly nodes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/838

#### Other
- Add auth to the config.toml.example file
  https://gitlab.com/gitlab-org/gitaly/merge_requests/851
- Remove the Dockerfile for Danger since the image is now built by https://gitlab.com/gitlab-org/gitlab-build-images
  https://gitlab.com/gitlab-org/gitaly/merge_requests/836
- Vendor Gitlab::Git at 2ca8219a20f16
  https://gitlab.com/gitlab-org/gitaly/merge_requests/841
- Move diff parser test to own package
  https://gitlab.com/gitlab-org/gitaly/merge_requests/837

## v0.114.0

#### Added
- Remove stale config.lock files
  https://gitlab.com/gitlab-org/gitaly/merge_requests/832

#### Fixed
- Handle nil commit in buildLocalBranch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/822
- Handle non-existing branch on UserDeleteBranch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/826
- Handle non-existing tags on UserDeleteTag
  https://gitlab.com/gitlab-org/gitaly/merge_requests/827

#### Other
- Lower gitaly-ruby default max_rss to 200MB
  https://gitlab.com/gitlab-org/gitaly/merge_requests/833
- Vendor gitlab-git at 92802e51
  https://gitlab.com/gitlab-org/gitaly/merge_requests/825
- Bump Linguist version to match Rails
  https://gitlab.com/gitlab-org/gitaly/merge_requests/821
- Stop vendoring gitlab/git/index.rb
  https://gitlab.com/gitlab-org/gitaly/merge_requests/824
- Bump rspec from 3.6.0 to 3.7.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/830

#### Performance
- Bump nokogiri to 1.8.4 and sanitize to 4.6.6
  https://gitlab.com/gitlab-org/gitaly/merge_requests/831

#### Security
- Update to gitlab-gollum-lib v4.2.7.5 and make Gemfile consistent with GitLab versions
  https://gitlab.com/gitlab-org/gitaly/merge_requests/828

## v0.113.0

#### Added
- Update Git to 2.18.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/795
- Implement RefService.FindAllRemoteBranches RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/799

#### Fixed
- Fix lines.Sender message chunking
  https://gitlab.com/gitlab-org/gitaly/merge_requests/817
- Fix nil commit author dereference
  https://gitlab.com/gitlab-org/gitaly/merge_requests/800

#### Other
- Vendor gitlab_git at 740ae2d194f3833e224
  https://gitlab.com/gitlab-org/gitaly/merge_requests/819
- Vendor gitlab-git at 49d7f92fd7476b4fb10e44f
  https://gitlab.com/gitlab-org/gitaly/merge_requests/798
- Vendor gitlab_git at 555afe8971c9ab6f9
  https://gitlab.com/gitlab-org/gitaly/merge_requests/803
- Move git/wiki*.rb out of vendor
  https://gitlab.com/gitlab-org/gitaly/merge_requests/804
- Clean up CI matrix
  https://gitlab.com/gitlab-org/gitaly/merge_requests/811
- Stop vendoring some gitlab_git files we don't need
  https://gitlab.com/gitlab-org/gitaly/merge_requests/801
- Vendor gitlab_git at 16b867d8ce6246ad8
  https://gitlab.com/gitlab-org/gitaly/merge_requests/810
- Vendor gitlab-git at e661896b54da82c0327b1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/814
- Catch SIGINT in gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/818
- Fix diff path logging
  https://gitlab.com/gitlab-org/gitaly/merge_requests/812
- Exclude more gitlab_git files from vendoring
  https://gitlab.com/gitlab-org/gitaly/merge_requests/815
- Improve ListError message
  https://gitlab.com/gitlab-org/gitaly/merge_requests/809

#### Performance
- Add limit parameter for WikiGetAllPagesRequest
  https://gitlab.com/gitlab-org/gitaly/merge_requests/807

#### Removed
- Remove implementation of Notifications::PostReceive
  https://gitlab.com/gitlab-org/gitaly/merge_requests/806

## v0.112.0

#### Fixed
- Translate more ListConflictFiles errors into FailedPrecondition
  https://gitlab.com/gitlab-org/gitaly/merge_requests/797
- Implement fetch keep-around refs in create from bundle
  https://gitlab.com/gitlab-org/gitaly/merge_requests/790
- Remove unnecessary commit size calculations
  https://gitlab.com/gitlab-org/gitaly/merge_requests/791

#### Other
- Add validation for config keys
  https://gitlab.com/gitlab-org/gitaly/merge_requests/788
- Vendor gitlab-git at b14b31b819f0f09d73e001
  https://gitlab.com/gitlab-org/gitaly/merge_requests/792

#### Performance
- Rewrite ListCommitsByOid in Go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/787

## v0.111.3

#### Security
- Update to gitlab-gollum-lib v4.2.7.5 and make Gemfile consistent with GitLab versions
  https://gitlab.com/gitlab-org/gitaly/merge_requests/828

## v0.111.2

#### Fixed
- Handle nil commit in buildLocalBranch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/822

## v0.111.1

#### Fixed
- Fix nil commit author dereference
  https://gitlab.com/gitlab-org/gitaly/merge_requests/800
- Remove unnecessary commit size calculations
  https://gitlab.com/gitlab-org/gitaly/merge_requests/791

## v0.111.0

#### Added
- Implement DeleteConfig and SetConfig
  https://gitlab.com/gitlab-org/gitaly/merge_requests/786
- Add OperationService.UserUpdateBranch RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/778

#### Other
- Vendor gitlab-git at 7e9f46d0dc1ed34d7
  https://gitlab.com/gitlab-org/gitaly/merge_requests/783
- Vendor gitlab-git at bdb64ac0a1396a7624
  https://gitlab.com/gitlab-org/gitaly/merge_requests/784
- Remove unnecessary existence check in AddNamespace
  https://gitlab.com/gitlab-org/gitaly/merge_requests/785

## v0.110.0

#### Added
- Server implementation ListNewCommits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/779

#### Fixed
- Fix encoding bug in UserCommitFiles
  https://gitlab.com/gitlab-org/gitaly/merge_requests/782

#### Other
- Tweak spawn token defaults and add logging
  https://gitlab.com/gitlab-org/gitaly/merge_requests/781

#### Performance
- Use 'git cat-file' to retrieve commits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/771

#### Security
- Sanitize paths when importing repositories
  https://gitlab.com/gitlab-org/gitaly/merge_requests/780

## v0.109.0

#### Added
- Reject nested storage paths
  https://gitlab.com/gitlab-org/gitaly/merge_requests/773

#### Fixed
- Bump rugged to 0.27.2
  https://gitlab.com/gitlab-org/gitaly/merge_requests/769
- Fix TreeEntry relative path bug
  https://gitlab.com/gitlab-org/gitaly/merge_requests/776

#### Other
- Vendor Gitlab::Git at 292cf668
  https://gitlab.com/gitlab-org/gitaly/merge_requests/777
- Vendor Gitlab::Git at f7b59b9f14
  https://gitlab.com/gitlab-org/gitaly/merge_requests/768
- Vendor Gitlab::Git at 7c11ed8c
  https://gitlab.com/gitlab-org/gitaly/merge_requests/770

## v0.108.0

#### Added
- Server info performs read and write checks
  https://gitlab.com/gitlab-org/gitaly/merge_requests/767

#### Changed
- Remove GoGit
  https://gitlab.com/gitlab-org/gitaly/merge_requests/764

#### Other
- Use custom log levels for grpc-go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/765
- Vendor Gitlab::Git at 2a82179e102159b8416f4a20d3349ef208c58738
  https://gitlab.com/gitlab-org/gitaly/merge_requests/766

## v0.107.0

#### Added
- Add BackupCustomHooks
  https://gitlab.com/gitlab-org/gitaly/merge_requests/760

#### Other
- Try to fix flaky rubyserver.TestRemovals test
  https://gitlab.com/gitlab-org/gitaly/merge_requests/759
- Vendor gitlab_git at a20d3ff2b004e8ab62c037
  https://gitlab.com/gitlab-org/gitaly/merge_requests/761
- Bumping gitlab-gollum-rugged-adapter to version 0.4.4.1 and gitlab-gollum-lib to 4.2.7.4
  https://gitlab.com/gitlab-org/gitaly/merge_requests/762

## v0.106.0

#### Changed
- Colons are not allowed in refs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/747

#### Fixed
- Reraise UnsupportedEncodingError as FailedPrecondition
  https://gitlab.com/gitlab-org/gitaly/merge_requests/718

#### Other
- Vendor gitlab_git at 930ad88a87b0814173989
  https://gitlab.com/gitlab-org/gitaly/merge_requests/752
- Upgrade vendor to d2aa3e3d5fae1017373cc047a9403cfa111b2031
  https://gitlab.com/gitlab-org/gitaly/merge_requests/755

## v0.105.1

#### Other
- Bumping gitlab-gollum-rugged-adapter to version 0.4.4.1 and gitlab-gollum-lib to 4.2.7.4
  https://gitlab.com/gitlab-org/gitaly/merge_requests/762

## v0.105.0

#### Added
- RestoreCustomHooks
  https://gitlab.com/gitlab-org/gitaly/merge_requests/741

#### Changed
- Rewrite Repository::Fsck in Go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/738

#### Fixed
- Fix committer bug in go-git adapter
  https://gitlab.com/gitlab-org/gitaly/merge_requests/748

## v0.104.0

#### Added
- Use Go-Git for the FindCommit RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/691

#### Fixed
- Ignore ENOENT when cleaning up lock files
  https://gitlab.com/gitlab-org/gitaly/merge_requests/740
- Fix rename similarity in CommitDiff
  https://gitlab.com/gitlab-org/gitaly/merge_requests/727
- Use grpc 1.11.0 in gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/732

#### Other
- Tests: only match error strings we create
  https://gitlab.com/gitlab-org/gitaly/merge_requests/743
- Use gitaly-proto 0.101.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/745
- Upgrade to Ruby 2.4.4
  https://gitlab.com/gitlab-org/gitaly/merge_requests/725
- Use the same faraday gem version as gitlab-ce
  https://gitlab.com/gitlab-org/gitaly/merge_requests/733

#### Performance
- Rewrite IsRebase/SquashInProgress in Go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/698

#### Security
- Use rugged 0.27.1 for security fixes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/744

## v0.103.0

#### Added
- Add StorageService::DeleteAllRepositories RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/726

#### Other
- Fix Dangerfile bad changelog detection
  https://gitlab.com/gitlab-org/gitaly/merge_requests/724

## v0.102.0

#### Changed
- Unvendor Repository#add_branch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/717

#### Fixed
- Fix matching bug in SearchFilesByContent
  https://gitlab.com/gitlab-org/gitaly/merge_requests/722

## v0.101.0

#### Changed
- Add gitaly-ruby installation debug log messages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/710

#### Fixed
- Use round robin load balancing instead of 'pick first' for gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/700

#### Other
- Generate changelog when releasing a tag to prevent merge conflicts
  https://gitlab.com/gitlab-org/gitaly/merge_requests/719
- Unvendor Repository#create implementation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/713

## v0.100.0

- Fix WikiFindPage when the page has invalidly-encoded content
  https://gitlab.com/gitlab-org/gitaly/merge_requests/712
- Add danger container to the Gitaly project
  https://gitlab.com/gitlab-org/gitaly/merge_requests/711
- Remove ruby concurrency limiter
  https://gitlab.com/gitlab-org/gitaly/merge_requests/708
- Drop support for Golang 1.8
  https://gitlab.com/gitlab-org/gitaly/merge_requests/715
- Introduce src-d/go-git as dependency
  https://gitlab.com/gitlab-org/gitaly/merge_requests/709
- Lower spawn log level to 'debug'
  https://gitlab.com/gitlab-org/gitaly/merge_requests/714

## v0.99.0

- Improve changelog entry checks using Danger
  https://gitlab.com/gitlab-org/gitaly/merge_requests/705
- GetBlobs: don't create blob reader if limit is zero
  https://gitlab.com/gitlab-org/gitaly/merge_requests/706
- Implement SearchFilesBy{Content,Name}
  https://gitlab.com/gitlab-org/gitaly/merge_requests/677
- Introduce feature flag package based on gRPC metadata
  https://gitlab.com/gitlab-org/gitaly/merge_requests/704
- Return DataLoss error for non-valid git repositories when calculating the checksum
  https://gitlab.com/gitlab-org/gitaly/merge_requests/697

## v0.98.0

- Server implementation for repository raw_changes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/699
- Add 'large request' test case to ListCommitsByOid
  https://gitlab.com/gitlab-org/gitaly/merge_requests/703
- Vendor gitlab_git at gitlab-org/gitlab-ce@3fcb9c115d776feb
  https://gitlab.com/gitlab-org/gitaly/merge_requests/702
- Limit concurrent gitaly-ruby requests from the client side
  https://gitlab.com/gitlab-org/gitaly/merge_requests/695
- Allow configuration of the log level in `config.toml`
  https://gitlab.com/gitlab-org/gitaly/merge_requests/696
- Copy Gitlab::Git::Repository#exists? implementation for internal method calls
  https://gitlab.com/gitlab-org/gitaly/merge_requests/693
- Upgrade Licensee gem to match the CE gem
  https://gitlab.com/gitlab-org/gitaly/merge_requests/693
- Vendor gitlab_git at 8b41c40674273d6ee
  https://gitlab.com/gitlab-org/gitaly/merge_requests/684
- Make wiki commit fields backwards compatible
  https://gitlab.com/gitlab-org/gitaly/merge_requests/685
- Catch CommitErrors while rebasing
  https://gitlab.com/gitlab-org/gitaly/merge_requests/680

## v0.97.0

- Use gitaly-proto 0.97.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/683
- Make gitaly-ruby's grpc server log at level WARN
  https://gitlab.com/gitlab-org/gitaly/merge_requests/681
- Add health checks for gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/678
- Add config option to point to languages.json
  https://gitlab.com/gitlab-org/gitaly/merge_requests/652

## v0.96.1

- Vendor gitlab_git at 7e3bb679a92156304
  https://gitlab.com/gitlab-org/gitaly/merge_requests/669
- Make it a fatal error if gitaly-ruby can't start
  https://gitlab.com/gitlab-org/gitaly/merge_requests/667
- Tag log entries with repo.GlRepository
  https://gitlab.com/gitlab-org/gitaly/merge_requests/663
- Add {Get,CreateRepositoryFrom}Snapshot RPCs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/644

## v0.96.0

Skipped. We cut and pushed the wrong tag.

## v0.95.0
- Fix fragile checksum test
  https://gitlab.com/gitlab-org/gitaly/merge_requests/661
- Use rugged 0.27.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/660

## v0.94.0

- Send gitaly-ruby exceptions to their own DSN
  https://gitlab.com/gitlab-org/gitaly/merge_requests/656
- Run Go test suite with '-race' in CI
  https://gitlab.com/gitlab-org/gitaly/merge_requests/654
- Ignore more grpc codes in sentry
  https://gitlab.com/gitlab-org/gitaly/merge_requests/655
- Implement Get{Tag,Commit}Messages RPCs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/646
- Fix directory permission walker for Go 1.10
  https://gitlab.com/gitlab-org/gitaly/merge_requests/650

## v0.93.0

- Fix concurrency limit handler stream interceptor
  https://gitlab.com/gitlab-org/gitaly/merge_requests/640
- Vendor gitlab_git at 9b76d8512a5491202e5a953
  https://gitlab.com/gitlab-org/gitaly/merge_requests/647
- Add handling for large commit and tag messages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/635
- Update gitaly-proto to v0.91.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/643

## v0.92.0

- Server Implementation GetInfoAttributes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/641
- Fix encoding error in ListConflictFiles
  https://gitlab.com/gitlab-org/gitaly/merge_requests/639
- Add catfile convenience methods
  https://gitlab.com/gitlab-org/gitaly/merge_requests/638
- Server implementation FindRemoteRepository
  https://gitlab.com/gitlab-org/gitaly/merge_requests/636
- Log process PID in 'spawn complete' entry
  https://gitlab.com/gitlab-org/gitaly/merge_requests/637
- Vendor gitlab_git at 79aa00321063da
  https://gitlab.com/gitlab-org/gitaly/merge_requests/633

## v0.91.0

- Rewrite RepositoryService.HasLocalBranches in Go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/629
- Rewrite RepositoryService.MergeBase in Go
  https://gitlab.com/gitlab-org/gitaly/merge_requests/632
- Encode OperationsService errors in UTF-8 before sending them
  https://gitlab.com/gitlab-org/gitaly/merge_requests/627
- Add param logging in NamespaceService RPCs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/626
- Sanitize URLs before sending gitaly-ruby exceptions to Sentry
  https://gitlab.com/gitlab-org/gitaly/merge_requests/625

## v0.90.0

- Implement SSHService.SSHUploadArchive RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/621
- Sanitize URLs before logging them
  https://gitlab.com/gitlab-org/gitaly/merge_requests/624
- Clean stale worktrees before performing garbage collection
  https://gitlab.com/gitlab-org/gitaly/merge_requests/622

## v0.89.0

- Report original exceptions to Sentry instead of wrapped ones by the exception bridge
  https://gitlab.com/gitlab-org/gitaly/merge_requests/623
- Upgrade grpc gem to 1.10.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/620
- Fix FetchRemote throwing "Host key verification failed"
  https://gitlab.com/gitlab-org/gitaly/merge_requests/617
- Use only 1 gitaly-ruby process in test
  https://gitlab.com/gitlab-org/gitaly/merge_requests/615
- Bump github-linguist to 5.3.3
  https://gitlab.com/gitlab-org/gitaly/merge_requests/613

## v0.88.0

- Add support for all field to {Find,Count}Commits RPCs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/611
- Vendor gitlab_git at de454de9b10f
  https://gitlab.com/gitlab-org/gitaly/merge_requests/611

## v0.87.0

- Implement GetCommitSignatures RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/609

## v0.86.0

- Implement BlobService.GetAllLfsPointers
  https://gitlab.com/gitlab-org/gitaly/merge_requests/562
- Implement BlobService.GetNewLfsPointers
  https://gitlab.com/gitlab-org/gitaly/merge_requests/562
- Use gitaly-proto v0.86.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/606

## v0.85.0

- Implement recursive tree entries fetching
  https://gitlab.com/gitlab-org/gitaly/merge_requests/600

## v0.84.0

- Send gitaly-ruby exceptions to Sentry
  https://gitlab.com/gitlab-org/gitaly/merge_requests/598
- Detect License type for repositories
  https://gitlab.com/gitlab-org/gitaly/merge_requests/601

## v0.83.0

- Delete old lock files before performing Garbage Collection
  https://gitlab.com/gitlab-org/gitaly/merge_requests/587

## v0.82.0

- Implement RepositoryService.IsSquashInProgress RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/593
- Added test to prevent wiki page duplication
  https://gitlab.com/gitlab-org/gitaly/merge_requests/539
- Fixed bug in wiki_find_page method
  https://gitlab.com/gitlab-org/gitaly/merge_requests/590

## v0.81.0

- Vendor gitlab_git at 7095c2bf4064911
  https://gitlab.com/gitlab-org/gitaly/merge_requests/591
- Vendor gitlab_git at 9483cbab26ad239
  https://gitlab.com/gitlab-org/gitaly/merge_requests/588

## v0.80.0

- Lock protobuf to 3.5.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/589

## v0.79.0

- Update the activesupport gem
  https://gitlab.com/gitlab-org/gitaly/merge_requests/584
- Update the grpc gem to 1.8.7
  https://gitlab.com/gitlab-org/gitaly/merge_requests/585
- Implement GetBlobs RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/582
- Check info split size in catfile parser
  https://gitlab.com/gitlab-org/gitaly/merge_requests/583

## v0.78.0

- Vendor gitlab_git at 498d32363aa61d679ff749b
  https://gitlab.com/gitlab-org/gitaly/merge_requests/579
- Convert inputs to UTF-8 before passing them to Gollum
  https://gitlab.com/gitlab-org/gitaly/merge_requests/575
- Implement OperationService.UserSquash RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/548
- Update recommended and minimum git versions to 2.14.3 and 2.9.0 respectively
  https://gitlab.com/gitlab-org/gitaly/merge_requests/548
- Handle binary commit messages better
  https://gitlab.com/gitlab-org/gitaly/merge_requests/577
- Vendor gitlab_git at a03ea19332736c36ecb9
  https://gitlab.com/gitlab-org/gitaly/merge_requests/574

## v0.77.0

- Implement RepositoryService.WriteConfig RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/554

## v0.76.0

- Add support for PreReceiveError in UserMergeBranch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/573
- Add support for Refs field in DeleteRefs RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/565
- Wait between ruby worker removal from pool and graceful shutdown
  https://gitlab.com/gitlab-org/gitaly/merge_requests/567
- Register the ServerService
  https://gitlab.com/gitlab-org/gitaly/merge_requests/572
- Vendor gitlab_git at f8dd398a21b19cb7d56
  https://gitlab.com/gitlab-org/gitaly/merge_requests/571
- Vendor gitlab_git at 4376be84ce18cde22febc50
  https://gitlab.com/gitlab-org/gitaly/merge_requests/570

## v0.75.0

- Implement WikiGetFormattedData RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/564
- Implement ServerVersion and ServerGitVersion
  https://gitlab.com/gitlab-org/gitaly/merge_requests/561
- Vendor Gitlab::Git @ f9b946c1c9756533fd95c8735803d7b54d6dd204
  https://gitlab.com/gitlab-org/gitaly/merge_requests/563
- ListBranchNamesContainingCommit server implementation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/537
- ListTagNamesContainingCommit server implementation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/537

## v0.74.0

- Implement CreateRepositoryFromBundle RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/557
- Use gitaly-proto v0.77.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/556
- Automatically remove tempdir when context is over
  https://gitlab.com/gitlab-org/gitaly/merge_requests/555
- Add automatic tempdir cleaner
  https://gitlab.com/gitlab-org/gitaly/merge_requests/540

## v0.73.0

- Implement CreateBundle RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/546

## v0.72.0

- Implement RemoteService.UpdateRemoteMirror RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/544
- Implement OperationService.UserCommitFiles RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/516
- Use grpc-go 1.9.1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/547

## v0.71.0

- Implement GetLfsPointers RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/543
- Add tempdir package
  https://gitlab.com/gitlab-org/gitaly/merge_requests/538
- Fix validation for Repositoryservice::WriteRef
  https://gitlab.com/gitlab-org/gitaly/merge_requests/542

## v0.70.0

- Handle non-existent commits in ExtractCommitSignature
  https://gitlab.com/gitlab-org/gitaly/merge_requests/535
- Implement RepositoryService::WriteRef
  https://gitlab.com/gitlab-org/gitaly/merge_requests/526

## v0.69.0

- Fix handling of paths ending with slashes in TreeEntry
  https://gitlab.com/gitlab-org/gitaly/merge_requests/532
- Implement CreateRepositoryFromURL RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/529

## v0.68.0

- Check repo existence before passing to gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/528
- Implement ExtractCommitSignature RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/521
- Update Gitlab::Git vendoring to b10ea6e386a025759aca5e9ef0d23931e77d1012
  https://gitlab.com/gitlab-org/gitaly/merge_requests/525
- Use gitlay-proto 0.71.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/524
- Fix pagination bug in GetWikiPageVersions
  https://gitlab.com/gitlab-org/gitaly/merge_requests/524
- Use gitaly-proto 0.70.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/522

## v0.67.0

- Implement UserRebase RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/511
- Implement IsRebaseInProgress RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/519
- Update to gitaly-proto v0.67.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/520
- Fix an error in merged all branches logic
  https://gitlab.com/gitlab-org/gitaly/merge_requests/517
- Allow RemoteService.AddRemote to receive several mirror_refmaps
  https://gitlab.com/gitlab-org/gitaly/merge_requests/513
- Update vendored gitlab_git to 33cea50976
  https://gitlab.com/gitlab-org/gitaly/merge_requests/518
- Update vendored gitlab_git to bce886b776a
  https://gitlab.com/gitlab-org/gitaly/merge_requests/515
- Update vendored gitlab_git to 6eeb69fc9a2
  https://gitlab.com/gitlab-org/gitaly/merge_requests/514
- Add support for MergedBranches in FindAllBranches RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/510

## v0.66.0

- Implement RemoteService.FetchInternalRemote RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/508

## v0.65.0

- Add support for MaxCount in CountCommits RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/507
- Implement CreateFork RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/497

## v0.64.0

- Update vendored gitlab_git to b98c69470f52185117fcdb5e28096826b32363ca
  https://gitlab.com/gitlab-org/gitaly/merge_requests/506

## v0.63.0

- Handle failed merge when branch gets updated
  https://gitlab.com/gitlab-org/gitaly/merge_requests/505

## v0.62.0

- Implement ConflictsService.ResolveConflicts RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/470
- Implement ConflictsService.ListConflictFiles RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/470
- Implement RemoteService.RemoveRemote RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/490
- Implement RemoteService.AddRemote RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/490

## v0.61.1

- gitaly-ruby shutdown improvements
  https://gitlab.com/gitlab-org/gitaly/merge_requests/500
- Use go 1.9
  https://gitlab.com/gitlab-org/gitaly/merge_requests/496

## v0.61.0

- Add rdoc to gitaly-ruby's Gemfile
  https://gitlab.com/gitlab-org/gitaly/merge_requests/487
- Limit the number of concurrent process spawns
  https://gitlab.com/gitlab-org/gitaly/merge_requests/492
- Update vendored gitlab_git to 858edadf781c0cc54b15832239c19fca378518ad
  https://gitlab.com/gitlab-org/gitaly/merge_requests/493
- Eagerly close logrus writer pipes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/489
- Panic if a command has no Done() channel
  https://gitlab.com/gitlab-org/gitaly/merge_requests/485
- Update vendored gitlab_git to 31fa9313991881258b4697cb507cfc8ab205b7dc
  https://gitlab.com/gitlab-org/gitaly/merge_requests/486

## v0.60.0

- Implement FindMergeBase RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/477
- Update vendored gitlab_git to 359b65beac43e009b715c2db048e06b6f96b0ee8
  https://gitlab.com/gitlab-org/gitaly/merge_requests/481

## v0.59.0

- Restart gitaly-ruby when it uses too much memory
  https://gitlab.com/gitlab-org/gitaly/merge_requests/465

## v0.58.0

- Implement RepostoryService::Fsck
  https://gitlab.com/gitlab-org/gitaly/merge_requests/475
- Increase default gitaly-ruby connection timeout to 40s
  https://gitlab.com/gitlab-org/gitaly/merge_requests/476
- Update vendored gitlab_git to f3a3bd50eafdcfcaeea21d6cfa0b8bbae7720fec
  https://gitlab.com/gitlab-org/gitaly/merge_requests/478

## v0.57.0

- Implement UserRevert RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/471
- Fix commit message encoding and support alternates in CatFile
  https://gitlab.com/gitlab-org/gitaly/merge_requests/469
- Raise an exception when Git::Env.all is called
  https://gitlab.com/gitlab-org/gitaly/merge_requests/474
- Update vendored gitlab_git to c594659fea15c6dd17b
  https://gitlab.com/gitlab-org/gitaly/merge_requests/473
- More logging in housekeeping
  https://gitlab.com/gitlab-org/gitaly/merge_requests/435

## v0.56.0

- Implement UserCherryPick RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/457
- Use grpc-go 1.8.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/466
- Fix a panic in ListFiles RPC when git process is killed abruptly
  https://gitlab.com/gitlab-org/gitaly/merge_requests/460
- Implement CommitService::FilterShasWithSignatures
  https://gitlab.com/gitlab-org/gitaly/merge_requests/461
- Implement CommitService::ListCommitsByOid
  https://gitlab.com/gitlab-org/gitaly/merge_requests/438

## v0.55.0

- Include pprof debug access in the Prometheus listener
  https://gitlab.com/gitlab-org/gitaly/merge_requests/442
- Run gitaly-ruby in the same directory as gitaly
  https://gitlab.com/gitlab-org/gitaly/merge_requests/458

## v0.54.0

- Implement RefService.DeleteRefs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/453
- Use --deployment flag for bundler and force `bundle install` on `make assemble`
  https://gitlab.com/gitlab-org/gitaly/merge_requests/448
- Update License as requested in: gitlab-com/organization#146
- Implement RepositoryService::FetchSourceBranch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/434

## v0.53.0

- Update vendored gitlab_git to f7537ce03a29b
  https://gitlab.com/gitlab-org/gitaly/merge_requests/449
- Update vendored gitlab_git to 6f045671e665e42c7
  https://gitlab.com/gitlab-org/gitaly/merge_requests/446
- Implement WikiGetPageVersions RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/430

## v0.52.1

- Include pprof debug access in the Prometheus listener
  https://gitlab.com/gitlab-org/gitaly/merge_requests/442

## v0.52.0

- Implement WikiUpdatePage RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/422

## v0.51.0

- Implement OperationService.UserFFMerge
  https://gitlab.com/gitlab-org/gitaly/merge_requests/426
- Implement WikiFindFile RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/425
- Implement WikiDeletePage RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/414
- Implement WikiFindPage RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/419
- Update gitlab_git to b3ba3996e0bd329eaa574ff53c69673efaca6eef and set
  `GL_USERNAME` env variable for hook excecution
  https://gitlab.com/gitlab-org/gitaly/merge_requests/423
- Enable logging in client-streamed and bidi GRPC requests
  https://gitlab.com/gitlab-org/gitaly/merge_requests/429

## v0.50.0

- Pass repo git alternate dirs to gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/421
- Remove old temporary files from repositories after GC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/411

## v0.49.0

- Use sentry fingerprinting to group exceptions
  https://gitlab.com/gitlab-org/gitaly/merge_requests/417
- Use gitlab_git c23c09366db610c1
  https://gitlab.com/gitlab-org/gitaly/merge_requests/415

## v0.48.0

- Implement WikiWritePage RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/410

## v0.47.0

- Pass full BranchUpdate result on successful merge
  https://gitlab.com/gitlab-org/gitaly/merge_requests/406
- Deprecate implementation of RepositoryService.Exists
  https://gitlab.com/gitlab-org/gitaly/merge_requests/408
- Use gitaly-proto 0.42.0
  https://gitlab.com/gitlab-org/gitaly/merge_requests/407


## v0.46.0

- Add a Rails logger to ruby-git
  https://gitlab.com/gitlab-org/gitaly/merge_requests/405
- Add `git version` to `gitaly_build_info` metrics
  https://gitlab.com/gitlab-org/gitaly/merge_requests/400
- Use relative paths for git object dir attributes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/393

## v0.45.1

- Implement OperationService::UserMergeBranch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/394
- Add client feature logging and metrics
  https://gitlab.com/gitlab-org/gitaly/merge_requests/392
- Implement RepositoryService.HasLocalBranches RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/397
- Fix Commit Subject parsing in rubyserver
  https://gitlab.com/gitlab-org/gitaly/merge_requests/388

## v0.45.0

Skipped. We cut and pushed the wrong tag.

## v0.44.0

- Update gitlab_git to 4a0f720a502ac02423
  https://gitlab.com/gitlab-org/gitaly/merge_requests/389
- Fix incorrect parsing of diff chunks starting with ++ or --
  https://gitlab.com/gitlab-org/gitaly/merge_requests/385
- Implement Raw{Diff,Patch} RPCs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/381

## v0.43.0

- Pass details of Gitaly-Ruby's Ruby exceptions back to
  callers in the request trailers
  https://gitlab.com/gitlab-org/gitaly/merge_requests/358
- Allow individual endpoints to be rate-limited per-repository
  https://gitlab.com/gitlab-org/gitaly/merge_requests/376
- Implement OperationService.UserDeleteBranch RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/377
- Fix path bug in CommitService::FindCommits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/364
- Fail harder during startup, fix version string
  https://gitlab.com/gitlab-org/gitaly/merge_requests/379
- Implement RepositoryService.GetArchive RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/370
- Add `gitaly-ssh` command
  https://gitlab.com/gitlab-org/gitaly/merge_requests/368

## v0.42.0

- Implement UserCreateTag RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/374
- Return pre-receive errors in UserDeleteTag response
  https://gitlab.com/gitlab-org/gitaly/merge_requests/378
- Check if we don't overwrite a namespace moved to gitaly
  https://gitlab.com/gitlab-org/gitaly/merge_requests/375

## v0.41.0

- Wait for monitor goroutine to return during supervisor shutdown
  https://gitlab.com/gitlab-org/gitaly/merge_requests/341
- Use grpc 1.6.0 and update all the things
  https://gitlab.com/gitlab-org/gitaly/merge_requests/354
- Update vendored gitlab_git to 4c6c105909ea610eac7
  https://gitlab.com/gitlab-org/gitaly/merge_requests/360
- Implement UserDeleteTag RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/366
- Implement RepositoryService::CreateRepository
  https://gitlab.com/gitlab-org/gitaly/merge_requests/361
- Fix path bug for gitlab-shell. gitlab-shell path is now required
  https://gitlab.com/gitlab-org/gitaly/merge_requests/365
- Remove support for legacy services not ending in 'Service'
  https://gitlab.com/gitlab-org/gitaly/merge_requests/363
- Implement RepositoryService.UserCreateBranch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/344
- Make gitaly-ruby config mandatory
  https://gitlab.com/gitlab-org/gitaly/merge_requests/373

## v0.40.0
- Use context cancellation instead of command.Close
  https://gitlab.com/gitlab-org/gitaly/merge_requests/332
- Fix LastCommitForPath handling of tree root
  https://gitlab.com/gitlab-org/gitaly/merge_requests/350
- Don't use 'bundle show' to find Linguist
  https://gitlab.com/gitlab-org/gitaly/merge_requests/339
- Fix diff parsing when the last 10 bytes of a stream contain newlines
  https://gitlab.com/gitlab-org/gitaly/merge_requests/348
- Consume diff binary notice as a patch
  https://gitlab.com/gitlab-org/gitaly/merge_requests/349
- Handle git dates larger than golang's and protobuf's limits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/353

## v0.39.0
- Reimplement FindAllTags RPC in Ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/334
- Re-use gitaly-ruby client connection
  https://gitlab.com/gitlab-org/gitaly/merge_requests/330
- Fix encoding-bug in GitalyServer#gitaly_commit_from_rugged
  https://gitlab.com/gitlab-org/gitaly/merge_requests/337

## v0.38.0

- Update vendor/gitlab_git to b58c4f436abaf646703bdd80f266fa4c0bab2dd2
  https://gitlab.com/gitlab-org/gitaly/merge_requests/324
- Add missing cmd.Close in log.GetCommit
  https://gitlab.com/gitlab-org/gitaly/merge_requests/326
- Populate `flat_path` field of `TreeEntry`s
  https://gitlab.com/gitlab-org/gitaly/merge_requests/328

## v0.37.0

- Implement FindBranch RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/315

## v0.36.0

- Terminate commands when their context cancels
  https://gitlab.com/gitlab-org/gitaly/merge_requests/318
- Implement {Create,Delete}Branch RPCs
  https://gitlab.com/gitlab-org/gitaly/merge_requests/311
- Use git-linguist to implement CommitLanguages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/316

## v0.35.0

- Implement CommitService.CommitStats
  https://gitlab.com/gitlab-org/gitaly/merge_requests/312
- Use bufio.Reader instead of bufio.Scanner for lines.Send
  https://gitlab.com/gitlab-org/gitaly/merge_requests/303
- Restore support for custom environment variables
  https://gitlab.com/gitlab-org/gitaly/merge_requests/319

## v0.34.0

- Export environment variables for git debugging
  https://gitlab.com/gitlab-org/gitaly/merge_requests/306
- Fix bugs in RepositoryService.FetchRemote
  https://gitlab.com/gitlab-org/gitaly/merge_requests/300
- Respawn gitaly-ruby when it crashes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/293
- Use a fixed order when auto-loading Ruby files
  https://gitlab.com/gitlab-org/gitaly/merge_requests/302
- Add signal handler for ruby socket cleanup on shutdown
  https://gitlab.com/gitlab-org/gitaly/merge_requests/304
- Use grpc 1.4.5 in gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/308
- Monitor gitaly-ruby RSS via Prometheus
  https://gitlab.com/gitlab-org/gitaly/merge_requests/310

## v0.33.0

- Implement DiffService.CommitPatch RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/279
- Use 'bundle config' for gitaly-ruby in source production installations
  https://gitlab.com/gitlab-org/gitaly/merge_requests/298

## v0.32.0

- RefService::RefExists endpoint
  https://gitlab.com/gitlab-org/gitaly/merge_requests/275

## v0.31.0

- Implement CommitService.FindCommits
  https://gitlab.com/gitlab-org/gitaly/merge_requests/266
- Log spawned process metrics
  https://gitlab.com/gitlab-org/gitaly/merge_requests/284
- Implement RepositoryService.ApplyGitattributes RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/278
- Implement RepositoryService.FetchRemote RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/276

## v0.30.0

- Add a middleware for handling Git object dir attributes
  https://gitlab.com/gitlab-org/gitaly/merge_requests/273

## v0.29.0

- Use BUNDLE_PATH instead of --path for gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/271
- Add GitLab-Shell Path to config
  https://gitlab.com/gitlab-org/gitaly/merge_requests/267
- Don't count on PID 1 to be the reaper
  https://gitlab.com/gitlab-org/gitaly/merge_requests/270
- Log top level project group for easier analysis
  https://gitlab.com/gitlab-org/gitaly/merge_requests/272

## v0.28.0

- Increase gitaly-ruby connection timeout to 20s
  https://gitlab.com/gitlab-org/gitaly/merge_requests/265
- Implement RepositorySize RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/262
- Implement CommitsByMessage RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/263

## v0.27.0

- Support `git -c` options in SSH upload-pack
  https://gitlab.com/gitlab-org/gitaly/merge_requests/242
- Add storage dir existence check to repo lookup
  https://gitlab.com/gitlab-org/gitaly/merge_requests/259
- Implement RawBlame RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/257
- Implement LastCommitForPath RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/260
- Deprecate Exists RPC in favor of RepositoryExists
  https://gitlab.com/gitlab-org/gitaly/merge_requests/260
- Install gems into vendor/bundle
  https://gitlab.com/gitlab-org/gitaly/merge_requests/264

## v0.26.0

- Implement CommitService.CommitLanguages, add gitaly-ruby
  https://gitlab.com/gitlab-org/gitaly/merge_requests/210
- Extend CountCommits RPC to support before/after/path arguments
  https://gitlab.com/gitlab-org/gitaly/merge_requests/252
- Fix a bug in FindAllTags parsing lightweight tags
  https://gitlab.com/gitlab-org/gitaly/merge_requests/256

## v0.25.0

- Implement FindAllTags RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/246

## v0.24.1

- Return an empty array on field `ParentIds` of `GitCommit`s if it has none
  https://gitlab.com/gitlab-org/gitaly/merge_requests/237

## v0.24.0

- Consume stdout during repack/gc
  https://gitlab.com/gitlab-org/gitaly/merge_requests/249
- Implement RefService.FindAllBranches RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/239

## v0.23.0

- Version without Build Time
  https://gitlab.com/gitlab-org/gitaly/merge_requests/231
- Implement CommitService.ListFiles
  https://gitlab.com/gitlab-org/gitaly/merge_requests/205
- Change the build process from copying to using symlinks
  https://gitlab.com/gitlab-org/gitaly/merge_requests/230
- Implement CommitService.FindCommit
  https://gitlab.com/gitlab-org/gitaly/merge_requests/217
- Register RepositoryService
  https://gitlab.com/gitlab-org/gitaly/merge_requests/233
- Correctly handle a non-tree path on CommitService.TreeEntries
  https://gitlab.com/gitlab-org/gitaly/merge_requests/234

## v0.22.0

- Various build file improvements
  https://gitlab.com/gitlab-org/gitaly/merge_requests/229
- Implement FindAllCommits RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/226
- Send full repository path instead of filename on field `path` of TreeEntry
  https://gitlab.com/gitlab-org/gitaly/merge_requests/232

## v0.21.2

- Config: do not start Gitaly without at least one storage
  https://gitlab.com/gitlab-org/gitaly/merge_requests/227
- Implement CommitService.GarbageCollect/Repack{Incremental,Full}
  https://gitlab.com/gitlab-org/gitaly/merge_requests/218

## v0.21.1

- Make sure stdout.Read has enough bytes buffered to read from
  https://gitlab.com/gitlab-org/gitaly/merge_requests/224

## v0.21.0

- Send an empty response for TreeEntry instead of nil
  https://gitlab.com/gitlab-org/gitaly/merge_requests/223

## v0.20.0

- Implement commit diff limiting logic
  https://gitlab.com/gitlab-org/gitaly/merge_requests/211
- Increase message size to 5 KB for Diff service
  https://gitlab.com/gitlab-org/gitaly/merge_requests/221

## v0.19.0

- Send parent ids and raw body on CommitService.CommitsBetween
  https://gitlab.com/gitlab-org/gitaly/merge_requests/216
- Streamio chunk size optimizations
  https://gitlab.com/gitlab-org/gitaly/merge_requests/206
- Implement CommitService.GetTreeEntries
  https://gitlab.com/gitlab-org/gitaly/merge_requests/208

## v0.18.0

- Add config to specify a git binary path
  https://gitlab.com/gitlab-org/gitaly/merge_requests/177
- CommitService.CommitsBetween fixes: Invert commits order, populates commit
  message bodies, reject suspicious revisions
  https://gitlab.com/gitlab-org/gitaly/merge_requests/204

## v0.17.0

- Rename auth 'unenforced' to 'transitioning'
  https://gitlab.com/gitlab-org/gitaly/merge_requests/209
- Also check for "refs" folder for repo existence
  https://gitlab.com/gitlab-org/gitaly/merge_requests/207

## v0.16.0

- Implement BlobService.GetBlob
  https://gitlab.com/gitlab-org/gitaly/merge_requests/202

## v0.15.0

- Ensure that sub-processes inherit TZ environment variable
  https://gitlab.com/gitlab-org/gitaly/merge_requests/201
- Implement CommitService::CommitsBetween
  https://gitlab.com/gitlab-org/gitaly/merge_requests/197
- Implement CountCommits RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/203

## v0.14.0

- Added integration test for SSH, and a client package
  https://gitlab.com/gitlab-org/gitaly/merge_requests/178/
- Override gRPC code to Canceled/DeadlineExceeded on requests with
  canceled contexts
  https://gitlab.com/gitlab-org/gitaly/merge_requests/199
- Add RepositoryExists Implementation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/200

## v0.13.0

- Added usage and version flags to the command line interface
  https://gitlab.com/gitlab-org/gitaly/merge_requests/193
- Optional token authentication
  https://gitlab.com/gitlab-org/gitaly/merge_requests/191

## v0.12.0

- Stop using deprecated field `path` in Repository messages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/179
- Implement TreeEntry RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/187

## v0.11.2

Skipping 0.11.1 intentionally, we messed up the tag.

- Add context to structured logging messages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/184
- Fix incorrect dependency in Makefile
  https://gitlab.com/gitlab-org/gitaly/merge_requests/189

## v0.11.0

- FindDefaultBranchName: decorate error
  https://gitlab.com/gitlab-org/gitaly/merge_requests/148
- Hide chatty logs behind GITALY_DEBUG=1. Log access times.
  https://gitlab.com/gitlab-org/gitaly/merge_requests/149
- Count accepted gRPC connections
  https://gitlab.com/gitlab-org/gitaly/merge_requests/151
- Disallow directory traversal in repository paths for security
  https://gitlab.com/gitlab-org/gitaly/merge_requests/152
- FindDefaultBranchName: Handle repos with non-existing HEAD
  https://gitlab.com/gitlab-org/gitaly/merge_requests/164
- Add support for structured logging via logrus
  https://gitlab.com/gitlab-org/gitaly/merge_requests/163
- Add support for exposing the Gitaly build information via Prometheus
  https://gitlab.com/gitlab-org/gitaly/merge_requests/168
- Set GL_PROTOCOL during SmartHTTP.PostReceivePack
  https://gitlab.com/gitlab-org/gitaly/merge_requests/169
- Handle server side errors from shallow clone
  https://gitlab.com/gitlab-org/gitaly/merge_requests/173
- Ensure that grpc server log messages are sent to logrus
  https://gitlab.com/gitlab-org/gitaly/merge_requests/174
- Add support for GRPC Latency Histograms in Prometheus
  https://gitlab.com/gitlab-org/gitaly/merge_requests/172
- Add support for Sentry exception reporting
  https://gitlab.com/gitlab-org/gitaly/merge_requests/171
- CommitDiff: Send chunks of patches over messages
  https://gitlab.com/gitlab-org/gitaly/merge_requests/170
- Upgrade gRPC and its dependencies
  https://gitlab.com/gitlab-org/gitaly/merge_requests/180

## v0.10.0

- CommitDiff: Parse a typechange diff correctly
  https://gitlab.com/gitlab-org/gitaly/merge_requests/136
- CommitDiff: Implement CommitDelta RPC
  https://gitlab.com/gitlab-org/gitaly/merge_requests/139
- PostReceivePack: Set GL_REPOSITORY env variable when provided in request
  https://gitlab.com/gitlab-org/gitaly/merge_requests/137
- Add SSHUpload/ReceivePack Implementation
  https://gitlab.com/gitlab-org/gitaly/merge_requests/132

## v0.9.0

- Add support ignoring whitespace diffs in CommitDiff
  https://gitlab.com/gitlab-org/gitaly/merge_requests/126
- Add support for path filtering in CommitDiff
  https://gitlab.com/gitlab-org/gitaly/merge_requests/126

## v0.8.0

- Don't error on invalid ref in CommitIsAncestor
  https://gitlab.com/gitlab-org/gitaly/merge_requests/129
- Don't error on invalid commit in FindRefName
  https://gitlab.com/gitlab-org/gitaly/merge_requests/122
- Return 'Not Found' gRPC code when repository is not found
  https://gitlab.com/gitlab-org/gitaly/merge_requests/120

## v0.7.0

- Use storage configuration data from config.toml, if possible, when
  resolving repository paths.
  https://gitlab.com/gitlab-org/gitaly/merge_requests/119
- Add CHANGELOG.md
