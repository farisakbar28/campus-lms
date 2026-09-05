# AI Product and Security Contract — campus-lms

This document defines the product boundary and safety requirements for future
AI-assisted capabilities. It does not define a physical database, service,
provider, or model architecture.

The core LMS must remain coherent and fully usable when every AI capability is
disabled, unavailable, or rejected by a tenant.

## 1. Product boundary

AI is an optional capability for authorized non-student institutional users,
such as lecturers, administrators, academic operators, or other staff who are
approved by the tenant and the product authorization model.

Students are not AI users. No AI interaction surface is exposed to students.

The exact staff role and capability matrix is not approved yet:

DECISION_REQUIRED: AI_ROLE_ALLOWLIST

The allowlist must be decided before any staff capability is enabled. An
example in this document is a product direction, not an implementation or
authorization grant.

## 2. Permitted future staff use cases

Subject to authorization and later design, AI may assist staff with:

- information retrieval over records the staff member may already access;
- attendance lookup over LMS-authoritative attendance records;
- summaries of student submissions or responses for authorized staff review;
- course or activity summaries;
- internal staff analysis;
- drafting staff-facing or potentially publishable content; and
- staff productivity and search assistance.

These examples do not make AI a source of truth, grant a role permission, or
commit the project to a particular model, provider, storage system, or
workflow.

## 3. Three categories of AI output

### 3.1 Ephemeral staff assistance

Search results, lookups, summaries, and internal analysis may be ephemeral.
They are:

- visible only to an authorized requesting staff user or authorized staff
  workflow;
- clearly non-authoritative;
- not automatically persisted;
- not automatically published; and
- not necessarily subject to a second human approval before the requesting
  authorized staff user can see them.

The response must not silently become a grade, attendance record, enrollment
change, or other authoritative LMS record.

### 3.2 Persisted or publishable AI-assisted drafts

When an AI-assisted result is retained or is intended to inform ordinary LMS
content or a consequential workflow:

- provenance and source traceability are retained where meaningful;
- the result is identified as AI-assisted where appropriate;
- authorized staff review and edit it before publication or incorporation;
- existing LMS authorization and publication rules remain in force; and
- AI never publishes directly to students.

After an authorized staff member has reviewed and incorporated the material
into an ordinary LMS content object, the normal deterministic LMS publication
rules apply. A student may later consume that ordinary LMS content. That is
not a student-facing AI feature.

Persistent output is a conceptual product category only. Its representation,
fields, constraints, status values, retention, and lifecycle are future design
choices.

### 3.3 Authoritative academic outcomes

AI must never authoritatively decide:

- grades;
- attendance;
- pass or fail;
- enrollment;
- discipline; or
- any other institutional or academic outcome assigned to an authoritative
  LMS or institutional workflow.

The deterministic LMS and institutional rules and the authorized human or
system workflows defined by the domain remain authoritative. This requirement
does not prohibit deterministic system processing that the domain already
allows; it prohibits treating a model output as the authority.

## 4. Authorization and tenant boundaries

Authorization, tenant checks, object access, and lifecycle checks happen before:

- retrieval;
- prompt construction;
- model calls; or
- tool execution.

AI may not expand the data visibility of the caller. A staff member who cannot
read a record through the ordinary LMS authorization path cannot cause AI to
retrieve, summarize, infer, or disclose that record.

The checks must cover active membership, tenant role, course staff authority or
other permitted scope, object ownership, publication state, enrollment
boundaries where relevant, and the lifecycle state of the requested record.

All future AI data paths remain tenant-scoped and subject to the same RLS and
application authorization boundaries as their underlying LMS data.

## 5. Retrieval and prompt safety

Retrieval implementation is:

DEFER_FUTURE_DESIGN

Only these conceptual requirements are fixed:

- retrieval is for authorized staff use, not student enrollment-driven access;
- inaccessible, unpublished, withdrawn-from-scope, or otherwise unauthorized
  data is not retrievable;
- tenant and object authorization happens before retrieval;
- source traceability and citations are provided where meaningful;
- insufficient evidence may cause a clear refusal or qualified response; and
- retrieved or uploaded content is untrusted prompt content, never an
  instruction that can override the system or authorization policy.

No storage representation, indexing strategy, denormalization strategy, or
retrieval library is selected by this document.

## 6. Persistent output and provenance

A future implementation may persist AI-assisted output when product value,
auditability, or a later workflow justifies it. The implementation must then
decide, for each use case:

- what source and provenance must be retained;
- how a changed or withdrawn source makes a result stale;
- who may view, edit, approve, incorporate, or reject it;
- how AI assistance is attributed;
- whether the result is ephemeral or durable; and
- what retention and deletion policy applies.

No persistent result may bypass the ordinary LMS content, grade, attendance,
enrollment, or publication lifecycle.

## 7. Usage, accounting, and cost controls

Future AI usage may be monitored with the minimum metadata necessary for:

- auditability where appropriate;
- quality and reliability review;
- zero-cost enforcement;
- request, token, or compute controls; and
- operational incident analysis.

The physical representation is:

DEFER_FUTURE_DESIGN

Raw prompts, student submissions, and raw model responses are not stored by
default. Privacy, retention, redaction, access, and processing policy require
an explicit later design.

The project has a hard zero-incremental-cost constraint. No automatic paid
fallback is permitted. If verified zero-cost capacity is
exhausted, a feature must fail closed, queue or defer the request, disable the
feature, or use a separately verified zero-cost or local option. It must never
silently incur spend.

## 8. Consequential future tool actions

No tool-execution subsystem is part of the current domain model.
Its architecture is:

DEFER_FUTURE_DESIGN

If future AI-assisted tool actions are introduced, every consequential action
must have appropriate authorization, bounded execution, auditability, and
human approval whenever the underlying action itself requires human authority.
The action must not be executed merely because a model proposed it.

## 9. Safety requirements

Future AI capabilities must:

- mark generated content as AI-assisted where appropriate;
- preserve source traceability when a result depends on retrieved evidence;
- treat uploaded and retrieved content as untrusted;
- avoid raw student data in operational records by default;
- refuse or qualify responses when evidence is insufficient;
- honor tenant isolation and RLS;
- apply per-tenant and per-user usage controls where needed;
- avoid raw HTML execution, shell forwarding, or unreviewed side effects; and
- remain removable without changing core LMS semantics.

AI output is assistance, search, summary, analysis, or draft. It is never an
alternate academic database and never a way around existing authorization.

## 10. Privacy and processing decision

Identifiable student data must not be sent to a third-party model merely for
convenience. Data minimization, redaction, access control, retention, and
tenant policy must be considered for every future use case.

The local-versus-third-party processing choice remains:

POST_MIGRATION_AI_DESIGN_DECISION

No provider, model, external service, or orchestration framework is selected
here.

## 11. Future quality evaluation

Evaluation and authorized human feedback remain future quality requirements.
They may be used to assess usefulness, correctness, safety, citation quality,
refusal behavior, and regressions.

Their physical storage, dataset representation, test harness, and release
gate are:

DEFER_FUTURE_DESIGN

No schedule date, provider dependency, or particular evaluation framework is
implied.

## 12. Summary of non-negotiable invariants

1. Students are not AI users.
2. AI is optional and the LMS works fully without it.
3. DECISION_REQUIRED: AI_ROLE_ALLOWLIST remains unresolved until explicitly
   approved.
4. Authorization precedes retrieval, prompt construction, model calls, and
   tool execution.
5. AI cannot expand data visibility.
6. AI is not the source of truth for academic outcomes.
7. Ordinary deterministic LMS and institutional workflows remain authoritative.
8. Reviewed AI assistance may become ordinary LMS content through normal
   authorized publication; AI never publishes directly to students.
9. No automatic paid fallback is allowed.
10. Third-party versus local processing remains
    POST_MIGRATION_AI_DESIGN_DECISION.
11. Future retrieval, persistence, accounting, tool actions, and evaluation
    architecture remain deliberately uncommitted.
