# ENG-017 — Simplify the engineering workflow to v2

Contract version: `campus-lms-work/v1`

Status: `APPROVED`

Issue: `ENG-017` / GitHub issue `#21`

Issue URL: <https://github.com/farisakbar28/campus-lms/issues/21>

Issue specification digest:
`sha256:4ceb385a4201e13b3322b0446f93fa60847dab6465ab0aa2f014787e1918b310`

Base Git revision:
`3cc3a32bbf339816a98ac0c8dc865a2f64d065bb`

Plan author actor label: `Codex planner`

Plan author session label: `ENG-017-plan-r3-2026-09-21`

Plan revision: `3`

Plan hash: `sha256:60716163fcececeaf5ba1a9bce73d0570aa494967f68d6ae5a5b6b0862226f96`

Risk: `CRITICAL`

Release required: `false`

Human plan-approval comment: `https://github.com/farisakbar28/campus-lms/issues/21#issuecomment-5761260917`

Implementation author actor label: `Codex implementer`

Implementation author session label: `ENG-017-in-plan-fix-2026-09-21`

Candidate Git commit SHA: `NONE`

This is the temporary v1 execution binding for the governance migration. The
v1 contract remains authoritative until the v2 Pull Request is human-merged.

<!-- PLAN-NORMATIVE-BEGIN -->

## Identity and input checkpoint

- Canonical task: open GitHub Issue #21, `ENG-017`, "Replace the v1 engineering
  workflow harness with a simpler repository-native v2 workflow"; title/body
  digest `sha256:4ceb385a4201e13b3322b0446f93fa60847dab6465ab0aa2f014787e1918b310`.
- Revision 3 is based on merged `master` commit
  `3cc3a32bbf339816a98ac0c8dc865a2f64d065bb`. PR #20 merged the
  `AUTHCTX-001` cleanup. At planning time `work/` contains only its README and
  templates; no other active issue or phase artifacts exist. Recheck the Issue,
  base, and active tree at implementation start and immediately before cutover.
- Current authorities: `docs/domain.md`, `docs/domain-ai.md`, `SECURITY.md`,
  accepted ADRs, `docs/roadmap.md`, root and scoped `AGENTS.md`, and the v1
  contract in `docs/engineering/ai-workflow.md`. Current source and native
  Git/GitHub/CI records govern implementation facts and evidence. ADR-0001
  remains proposed. No new ADR or domain decision is in scope.

## Goal and boundaries

Replace the v1 repository-local workflow harness with the Issue's simpler
repository-native workflow and clear source routing. Preserve independent plan
and implementation review, explicit human task choice and plan approval, human
merge, separate production approval, fail-closed security decisions, and
explicit authority for remote operations. Destructive database/data and
external/cloud operations require both explicit task authority and human
confirmation; CI success or plan approval supplies neither.
The v2 documents become authoritative together only when the reviewed PR is
human-merged into `master`.

This is a dedicated engineering-governance change. Do not change LMS runtime,
API/authentication/authorization implementation, schema, applied migrations,
dependency files, deployment topology, ingress, cloud resources, branch
protection, provider or production-readiness claims, or accepted ADRs. Do not
add a replacement validator, workflow engine, task-state scheme, digest system,
receipt ledger, archive, service, or paid tooling. Retain the zero incremental
paid infrastructure/tooling constraint. `release_required=false`; no production
readiness review, deployment, or production approval applies to this work.

## Implementation plan

1. Before implementation, re-fetch Issue #21 and its human approval; verify
   the exact digest and plan revision/hash, fresh independent review of this
   revision with all prior plan findings resolved, including the final-diff
   review gap, current base, and absence of unrelated active work or phases.
   If title/body changes, authority is
   unclear, or another v1 item becomes active, stop and return to planning or
   finish/explicitly migrate that item before removing `work/`.
2. On a short-lived branch, draft `docs/README.md` as source-of-truth index;
   `docs/architecture.md` as a concise facts/decisions/gaps overview; and
   `docs/engineering/history.md` as curated, verifiable milestones with a
   clearly identified v2 transition. Keep history out of task-archive and
   command-receipt territory. Trace implementation statements to current code
   and keep proposals distinct from accepted ADRs and observed facts. Make
   `docs/architecture.md` the current repository implementation snapshot,
   distinguishing code facts, accepted decisions, proposed decisions, known
   gaps, and future work. The roadmap keeps directional implementation
   classifications and dependency order, with a clear link to that snapshot;
   remove any competing detailed current-state snapshot. Move the roadmap's
   dated `CI / API` evidence checkpoint into living history with its date and
   exact SHA, labeled as historical evidence rather than current CI status.
   Verify those milestones against native Git/GitHub records; keep current
   hosted CI status in GitHub Actions, not in a static document.
3. Add `docs/engineering/workflow.md`, `agent-playbook.md`, and
   `production-readiness.md`. Document the complete Issue target sequence,
   native evidence, independent reviews, material-change return to planning,
   finding gates, and authority stops. Require explicit task authority and
   human confirmation before destructive database/data or external/cloud
   operations. For production-impacting work, require a separate readiness
   review and directly human-authored approval bound to the exact revision or
   immutable artifact, target environment, and approved scope; a material
   change invalidates that approval. Document rechecks before deployment,
   safe recovery and stop conditions, and post-deploy verification without
   secrets or student personal data. Specify that a post-merge history update
   uses a separately reviewed short-lived branch/PR, CI, and human merge, or
   is prepared in the task PR when the milestone is already known; no direct
   agent write to protected `master` is authorized.
4. Add `.github/pull_request_template.md` with task link, summary/scope, risk,
   security and tenancy impact, verification, independent review,
   documentation/history, and release/production impact. Simplify root
   `AGENTS.md`, `CONTRIBUTING.md`, `docs/roadmap.md`, and `README.md` to route
   to v2 and preserve their product, safety, sequencing, classification, and
   honest implementation boundaries. Strengthen root `AGENTS.md` so its
   destructive database/data rule also requires human confirmation, matching
   Issue #21; retain the existing human-confirmation rule for external/cloud
   operations. Change scoped `AGENTS.md` files only for link maintenance if
   needed; preserve their semantic constraints. Leave `docs/domain.md`,
   `docs/domain-ai.md`, `SECURITY.md`, and accepted ADRs intact.
5. Prepare the cutover while v1 still exists: check the new documents and
   cross-links, run `make workflow-test` and `make workflow-check` on the
   pre-retirement state, and capture only concise results in native review
   records. Then create the v1 implementation candidate commit that removes
   `docs/engineering/ai-workflow.md`, `work/README.md`, `work/templates/`, both
   validator scripts, both workflow Make targets, and exactly the two validator
   CI steps. Retain only `work/active/ENG-017/WORK.md` and `REVIEWS.md`
   temporarily for v1 transition. Do not run the retired validator against
   this candidate or claim it passed there. Keep CI name `CI`, job `api`/`API`,
   triggers, permissions, concurrency, actions, Go setup and all application
   checks otherwise identical.
6. Obtain a fresh independent v1 implementation review of that exact candidate
   commit, its base, candidate diff, and observed checks. Record the verdict
   and findings in `REVIEWS.md` while the active files still exist. Resolve all
   blocking findings; agents do not accept residual risk. A relevant candidate
   change requires a new review; a material deviation restarts planning and
   human approval.
7. With separate explicit remote-mutation authority, hand the ENG-017 plan
   approval, Issue digest, v1 candidate SHA, independent review verdict,
   release choice, and applicable verification references to the canonical
   Issue record. Only after that handoff, make an authorized local cleanup
   commit deleting exactly the remaining ENG-017 `WORK.md` and `REVIEWS.md`;
   its v1-candidate-to-cleanup diff must be those two deletions only. This
   cleanup commit is the final PR candidate; its base-to-head diff must remove
   the entire `work/` tree and contain every intended v2 change.
8. Before creating any PR, obtain a fresh independent implementation review
   of the exact final PR candidate commit, its base-to-head diff including
   both active-file deletions and all other harness removals, the cleanup-only
   comparison, and applicable verification results. With separate explicit
   Issue-mutation authority, record this second verdict and findings on the
   canonical Issue, bound to the final commit SHA, plan revision/hash, and
   Issue digest. If that authority or a passing independent verdict is absent,
   stop before PR creation. Any relevant change after this review requires a
   new review of the new complete final diff; blocking findings stop progress,
   and material deviations return to planning. For an in-plan fix, restore the
   transitional v1 artifacts on the branch without rewriting history and
   repeat the v1 candidate, handoff, cleanup, and complete-diff review gates.
   Do not treat the earlier v1 review or a proposed cleanup as review of the
   final diff.
9. With separate explicit push/PR authority, create the PR only from the exact
   reviewed final commit, record/link its independent pre-PR review, observe
   hosted `CI / API` for its head, and leave squash merge to the human. Revert,
   if needed, via an ordinary reviewed PR without rewriting history.

## Acceptance criteria

- All Issue #21 required documents and PR template exist, with usable links
  from `README.md` and `docs/README.md`. Workflow, playbook, readiness, and
  history have the distinct roles and authority boundaries specified above.
- Architecture and roadmap preserve source ownership, classifications,
  dependency order, unresolved decisions, and current limitations. Root
  guidance is internally consistent with v2; scoped guidance, domain sources,
  `SECURITY.md`, and ADR semantics are preserved, including proposed ADR-0001.
- Current-state ownership is explicit: architecture holds the factual
  implementation snapshot, roadmap retains product direction and high-level
  classifications without a competing detailed snapshot, history holds dated
  and verified milestones including the old CI checkpoint, and GitHub Actions
  remains the source for current hosted CI results. Cross-links and claims do
  not imply that a dated run validates today's head.
- Root guidance and v2 workflow/playbook require both explicit task authority
  and human confirmation for destructive database/data and external/cloud
  operations. Production guidance requires a separate readiness review and
  human approval bound to exact revision or immutable artifact, target
  environment, and approved scope, with material changes invalidating it.
  Deployment must stop if any binding or authority check fails; approved safe
  recovery and sanitized post-deploy verification are documented.
- Final PR diff removes the v1 contract, whole `work/` tree, both validator
  scripts, Make targets, and only the two workflow-validator CI steps. No live
  reference to deleted surfaces remains; historical mentions are labeled.
- `CI / API` identity, triggers, permissions, concurrency, checkout, Go setup,
  module download/verification, formatting, vet, test, build, and unchanged
  module-file verification remain. No runtime, schema, authentication,
  deployment, dependency, or cloud change enters the diff.
- The v1 candidate has its own independent review and the cleanup-only
  descendant rule holds. A second independent review covers the exact final
  cleanup commit and complete base-to-head diff, including both active-file
  deletions, before PR creation; its verdict is bound in the canonical Issue
  and linked from the PR. Hosted `CI / API` passes before human merge. Release
  remains unnecessary.

## Verification and documentation impact

Inspect the complete final diff for source precedence and authority conflicts;
run `git diff --check`, `make help`, focused obsolete-reference searches,
`make test`, and `make build`; compare CI against the preserved step list;
confirm no Go module/dependency or prohibited file change; and manually check
Markdown navigation and relative links. Specifically compare the root and
scoped authority wording against Issue #21, check production approval binding
and invalidation examples, and cross-check architecture/roadmap/history status
ownership and the dated CI evidence against native records. Run v1 workflow
tests/checks before retirement as specified above. Re-read the canonical Issue
and approval at required gates. Before PR creation, compare the exact final
commit with the v1-reviewed candidate to confirm cleanup-only descent and
have the independent reviewer inspect the complete final base-to-head diff.
Observe hosted PR `CI / API` before merge.
Report actual results and environmental limits, never assumed success.
Documentation impact is the new v2 source map and engineering guidance,
updated root/roadmap/README, new PR template, and deletion of superseded v1
guidance and harness.

## Risks and deviations

Risk is `CRITICAL` because replacing governance can silently weaken human
authority, review independence, or production safety. Review the exact text
of each retained boundary, CI identity and application checks, current-state
claims, and cutover ordering. Security-sensitive ambiguity fails closed.
Changes to Issue specification, scope, behavior, architecture, security,
authority, acceptance, verification, risk, or production impact are material:
record a concise deviation, revise this normative plan with a new revision and
hash, obtain fresh independent plan review and human approval, then resume.

<!-- PLAN-NORMATIVE-END -->

- Material deviations: `NONE`
- Concise verification result: `git diff --check, make help, make test,
  make build, go mod verify, go vet, Markdown-link check, and focused obsolete
  reference search passed for this in-plan fix; hosted PR CI is pending.`
- Concise completion summary: `Documentation contradictions corrected and the
  ENG-017 transitional files restored. Independent v1 candidate review, Issue
  handoff, cleanup, and final-diff review remain pending.`
