# ADR-0002c: Azure Naming, Region, and Tagging Conventions

- **Date:** 2026-08-16
- **Status:** 🟢 **Accepted**
- **Roadmap week:** 0

> **This file demonstrates the ADR ownership rule** (`agent/policy.md`):
> an agent may draft **Context** and **Options**; the **Decision** and
> **Consequences** sections must be written by the human, because you are the
> one who will have to defend the trade-off in an interview.
>
> Sections 1 and 2 below were drafted. Sections 3 and 4 record the human
> decision made for the project.

---

## 1. Context (drafted)

The project runs on an **Azure for Students** subscription: $100 credit valid
for 12 months, no credit card attached, with a hard spending limit. Production
is a single small VM; the database and frontend live on other providers' free
tiers.

Constraints that shape this decision:

- **Cost visibility matters more than usual.** With a fixed $100 for 12 months,
  an untracked resource is not a rounding error — it is a meaningful fraction of
  the entire budget.
- **Latency to users.** Primary users are in Indonesia (developer in Bali).
- **Region availability is not guaranteed.** Student subscriptions sometimes
  have zero vCPU quota in a given region, which is why quota is verified in
  Week 0 rather than discovered in Week 4.
- **The database is on Neon, not Azure.** Cross-network latency between the API
  and the database is a real cost, so their regions should be close.
- **Solo developer.** Conventions must be simple enough to follow consistently
  without tooling to enforce them.

## 2. Options considered (drafted)

### Region

| Option                         | Pros                                                                            | Cons                                                                    |
| ------------------------------ | ------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| **Southeast Asia** (Singapore) | Lowest latency from Indonesia; Neon also offers Singapore, keeping API↔DB close | Popular region, quota sometimes constrained                             |
| **East Asia** (Hong Kong)      | Usually good availability                                                       | Higher latency; fewer matching free-tier regions elsewhere in the stack |
| **Australia East**             | Often has spare quota                                                           | Noticeably higher latency from Indonesia                                |

### Resource grouping

| Option                                                   | Pros                                                     | Cons                                                       |
| -------------------------------------------------------- | -------------------------------------------------------- | ---------------------------------------------------------- |
| **Single resource group** (`rg-campuslms-prod`)          | One command to tear everything down; simple mental model | No isolation between environments if you later add staging |
| **Per-environment groups** (`rg-campuslms-prod`, `-dev`) | Cleaner separation, per-environment cost view            | More overhead for a solo project with one real environment |

### Tagging

| Option                                         | Pros                                                              | Cons                                                            |
| ---------------------------------------------- | ----------------------------------------------------------------- | --------------------------------------------------------------- |
| **Mandatory tags** (`project`, `env`, `owner`) | Cost reports become readable; orphaned resources are identifiable | Requires discipline on every resource                           |
| No tags                                        | Nothing to remember                                               | Cost analysis becomes guesswork — a real risk on a fixed budget |

## 3. Decision

We will use **Southeast Asia (Singapore)** as the primary Azure region for the
production environment because it provides the lowest expected latency for
users in Indonesia and keeps the Azure API geographically close to the Neon
database region. Production resources will be grouped under a single resource
group named `rg-campuslms-prod`. All Azure resources that support tagging must
use the mandatory tags `project=campus-lms`, `env=prod`, and `owner=faris`.
The Week 0 quota check confirmed that the **Standard BS Family vCPUs** quota
in Southeast Asia is **0 of 4 vCPUs currently in use**, meaning 4 vCPUs are
currently available and no regional fallback is required at this stage.

**Region:** Southeast Asia (Singapore)

**Resource group:** `rg-campuslms-prod`

**Mandatory tags:**

- `project=campus-lms`
- `env=prod`
- `owner=faris`

## 4. Consequences

**Positive:**

- Low expected latency for users in Indonesia.
- Keeps the Azure API geographically close to the Neon database region.
- A single production resource group keeps resource management simple for a
  solo developer.
- A single resource group also makes it straightforward to remove the
  production resources when the project is no longer needed.
- Mandatory tags make resource ownership and cost analysis easier.
- The current Southeast Asia quota provides sufficient Standard BS Family
  capacity for the planned small production VM.
- The convention is simple enough to apply consistently without additional
  infrastructure or policy tooling.

**Negative / accepted technical debt:**

- Southeast Asia is a popular Azure region, so quota availability may become
  constrained in the future.
- A single production resource group does not provide environment isolation if
  staging or development resources are introduced later.
- Mandatory tagging relies on manual discipline because there is currently no
  automated enforcement.
- If the required Standard BS Family vCPU quota becomes unavailable, the
  project may need to use another Azure region or request additional quota.

**When to revisit this decision:**

- If a second environment such as staging is introduced.
- If the required VM quota becomes unavailable in Southeast Asia.
- If latency between the Azure API and Neon becomes a measurable performance
  bottleneck.
- If the project grows beyond the capacity or cost assumptions of the current
  Azure for Students setup.
- If the selected Azure region no longer provides sufficient quota for the
  required VM configuration.

## 5. Notes

Evidence for the quota check performed in Week 0:

`docs/progress/evidence/week-00/<filename>`

The quota check was performed on **2026-08-16**.

The Southeast Asia quota check showed:

- **Quota family:** Standard BS Family vCPUs
- **Region:** Southeast Asia
- **Current usage:** 0 vCPUs
- **Quota limit:** 4 vCPUs
- **Result:** Sufficient quota for the planned small production VM; no regional
  fallback is required at this stage.
