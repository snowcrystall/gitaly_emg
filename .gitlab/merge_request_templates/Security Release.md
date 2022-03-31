<!--
# README first!
This MR should be created on `gitlab.com/gitlab-org/security/gitaly`.

See [the general developer security release guidelines](https://gitlab.com/gitlab-org/release/docs/blob/master/general/security/developer.md).

-->

## Related issues

<!-- Mention the GitLab Security issue this MR is related to -->

## Developer checklist

- [ ] **On "Related issues" section, write down the [GitLab Security] issue it belongs to (i.e. `Related to <issue_id>`).**
- [ ] Merge request targets `master`, or `X-Y-stable` for backports.
- [ ] Milestone is set for the version this merge request applies to. A closed milestone can be assigned via [quick actions].
- [ ] Title of this merge request is the same as for all backports.
- [ ] A [CHANGELOG entry] has been included, with `Changelog` trailer set to `security`.
- [ ] Assign to a reviewer and maintainer, per our [Code Review process].
- [ ] For the MR targeting `master`:
  - [ ] Ensure it's approved according to our [Approval Guidelines].
- [ ] Merge request _must not_ close the corresponding security issue, _unless_ it targets `master`.

**Note:** Reviewer/maintainer should not be a Release Manager

## Maintainer checklist
- [ ] Correct milestone is applied and the title is matching across all backports
- [ ] Assigned to `@gitlab-release-tools-bot` with passing CI pipelines and **when all backports including the MR targeting master are ready.**

/label ~security

[GitLab Security]: https://gitlab.com/gitlab-org/security/gitlab
[approval guidelines]: https://docs.gitlab.com/ee/development/code_review.html#approval-guidelines
[Code Review process]: https://docs.gitlab.com/ee/development/code_review.html
[quick actions]: https://docs.gitlab.com/ee/user/project/quick_actions.html#quick-actions-for-issues-merge-requests-and-epics
[CHANGELOG entry]: https://docs.gitlab.com/ee/development/changelog.html#overview
