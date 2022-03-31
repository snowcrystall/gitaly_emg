<!--
# Read me first!

Create this issue under https://gitlab.com/gitlab-org/security

Set the title to: `Description of the original issue`
-->

## Prior to starting the security release work

- [ ] Read the [security process for developers] if you are not familiar with it.
- [ ] Mark this [issue as related] to the Security Release tracking issue. You can find it on the topic of the `#releases` Slack channel.
- [ ] Run `scripts/security-harness` in your local repository to prevent accidentally pushing to any remote besides `gitlab.com/gitlab-org/security`.
- Fill out the [Links section](#links):
  - [ ] Next to **Issue on Gitaly**, add a link to the `gitlab-org/gitaly` issue that describes the security vulnerability.
  - [ ] Next to **Security Release tracking issue**, add a link to the security release issue that will include this security issue.

## Development

- [ ] Create a new branch prefixing it with `security-`.
- [ ] Create a merge request targeting `master` on `gitlab.com/gitlab-org/security` and use the [Security Release merge request template].
- [ ] Follow the same [code review] process: Assign to a reviewer, then to a maintainer.

After your merge request has been approved according to our [approval guidelines], you're ready to prepare the backports

## Backports

- [ ] Once the MR is ready to be merged, create MRs targeting the latest 3 stable branches
   * At this point, it might be easy to squash the commits from the MR into one
   * You can use the script `bin/secpick` instead of the following steps, to help you cherry-picking. See the [secpick documentation]
- [ ] Create each MR targeting the stable branch `X-Y-stable`, using the [Security Release merge request template].
   * Every merge request will have its own set of TODOs, so make sure to complete those.
- [ ] On the "Related merge requests" section, ensure all MRs are linked to this issue.
   * This section should only list the merge requests created for this issue: One targeting `master` and the 3 backports.

## Documentation and final details

- [ ] Ensure the [Links section](#links) is completed.
- [ ] Add the GitLab [versions](https://gitlab.com/gitlab-org/release/docs/-/blob/master/general/security/developer.md#versions-affected) and editions affected to the [details section](#details)
  * The Git history of the files affected may help you associate the issue with a [release](https://about.gitlab.com/releases/)
- [ ] Fill in any upgrade notes that users may need to take into account in the [details section](#details)
- [ ] Add the nickname of the external user who found the issue (and/or HackerOne profile) to the Thanks row in the [details section](#details)

## Summary

### Links

| Description | Link |
| -------- | -------- |
| Issue on [Gitaly](https://gitlab.com/gitlab-org/gitaly/issues) | #TODO  |
| Security Release tracking issue | #TODO  |

### Details

| Description | Details | Further details|
| -------- | -------- | -------- |
| Versions affected | X.Y  | |
| Upgrade notes | | |
| Thanks | | |

[gitaly-ce-version]: https://gitlab.com/gitlab-org/gitlab-ce/blob/master/GITALY_SERVER_VERSION
[security process for developers]: https://gitlab.com/gitlab-org/release/docs/blob/master/general/security/developer.md
[secpick documentation]: https://gitlab.com/gitlab-org/release/docs/-/blob/master/general/security/utilities/secpick_script.md
[security Release merge request template]: https://gitlab.com/gitlab-org/security/gitaly/blob/master/.gitlab/merge_request_templates/Security%20Release.md
[code review process]: https://docs.gitlab.com/ee/development/code_review.html
[approval guidelines]: https://docs.gitlab.com/ee/development/code_review.html#approval-guidelines
[issue as related]: https://docs.gitlab.com/ee/user/project/issues/related_issues.html#adding-a-related-issue

/label ~"devops::create" ~"group::gitaly" ~"security"
