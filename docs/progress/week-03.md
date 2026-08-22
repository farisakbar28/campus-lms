# Laporan Minggu 3 — PostgreSQL Mendalam & Multi-Tenancy

> Template. Salin ke `docs/progress/week-<NN>.md`. Disusun agent, ditandatangani manusia.
> Jangan hapus satu section pun. Section yang tidak relevan diisi "Tidak ada minggu ini".

- **Periode:** 2026-08-13 s/d 2026-08-20
- **Fokus roadmap:** PostgreSQL Mendalam & Multi-Tenancy (11 tabel Tier 1, RLS, composite FK, A1/A2/A7).
- **Total jam:** NOT MEASURED / to be filled by human.
- **Commit/history reference:** Fixed Task-2 implementation commits: `9c5e2e0`, `2c3a31f`, `e52d702`, `23a8058`, and `0081188`. This deliberately does not use a mutable `HEAD` range and does not include historical orphaned commit `af6a0d3`.

---

## 1. Ringkasan

**FACT:** Closeout Task-2 menambahkan tooling `golang-migrate` yang dipin pada v4.18.3 dan SQL seeder deterministik; implementasinya tercatat secara tetap pada `9c5e2e0`, `2c3a31f`, `e52d702`, `23a8058`, dan `0081188`. Migration 0001–0004 tidak berubah. **FACT:** Bukti v4 membuktikan fresh disposable database melalui transisi 0 → 4 → 3 → 4, strict structural check, dan cleanup; main database yang sudah dibaseline tetap version 4, `dirty=false`, melewati `force`, dan `migrate up` tidak mengubah apa pun. **FACT:** Seeder transaksional melakukan UPSERT/reconciliation tanpa `TRUNCATE` atau delete/rebuild destruktif; run kedua mempertahankan 14 jumlah fixture yang sama dan 14 metrik pelanggaran invariant bernilai 0. **FACT:** `make test` dengan race detector lulus (EXIT 0); `make lint` belum terverifikasi karena `golangci-lint` tidak tersedia (make EXIT 2, underlying Error 127). Week 3 belum selesai karena repository/endpoint, query tuning, backup/restore, dokumentasi/ADR, integrasi RLS pada path Go, lint final, quiz, explain-back, dan human sign-off masih terbuka.

---

## 2. Dikerjakan Agent

| # | Pekerjaan | File | Commit | Bukti |
|---|---|---|---|---|
| 1 | Create migration 0004: CHECK constraint academic_terms_valid_time_range | `apps/api/migrations/0004_academic_term_time_range_check.up.sql`, `.down.sql` | 68ab28f | [0004-up-v2.txt](evidence/week-03/0004-up-v2.txt) |
| 2 | Code snapshot commit: complete Week 3 migration set | 7 migration files | 68ab28f | git log |
| 3 | Create committed verification harness (19 SQL files) | `apps/api/testdata/week-03/*.sql` | 6e8e6d8 | git log |
| 4 | Create disposable verification DB `campus_lms_week03_verify_v3` | — | 6e8e6d8 | (inline) |
| 5 | Apply 0001, 0002, 0003, 0004 atomically | — | 6e8e6d8 | [0001-up-v2.txt](evidence/week-03/0001-up-v2.txt), [0002-up-v2.txt](evidence/week-03/0002-up-v2.txt), [0003-up-v2.txt](evidence/week-03/0003-up-v2.txt), [0004-up-v2.txt](evidence/week-03/0004-up-v2.txt) |
| 6 | Fixture setup v3 (deterministic UUIDs, all INSERTs succeed) | — | 6e8e6d8 | [fixtures-v3.txt](evidence/week-03/fixtures-v3.txt) |
| 7 | A1 v3: two distinct active roles, duplicate rejected, revoke/regrant | — | 6e8e6d8 | [a1-two-active-v3.txt](evidence/week-03/a1-two-active-v3.txt), [a1-duplicate-active-v3.txt](evidence/week-03/a1-duplicate-active-v3.txt), [a1-revoke-regrant-v3.txt](evidence/week-03/a1-revoke-regrant-v3.txt) |
| 8 | A2 v3: 5 composite FK cross-tenant violations (each names intended FK) | — | 6e8e6d8 | [a2-course-v3.txt](evidence/week-03/a2-course-v3.txt), [a2-membership-role-v3.txt](evidence/week-03/a2-membership-role-v3.txt), [a2-course-staff-v3.txt](evidence/week-03/a2-course-staff-v3.txt), [a2-audit-log-v3.txt](evidence/week-03/a2-audit-log-v3.txt), [a2-enrollment-offering-v3.txt](evidence/week-03/a2-enrollment-offering-v3.txt) |
| 9 | A7 v3: valid enrollment succeeds, cross-tenant membership fails on A7 FK | — | 6e8e6d8 | [a7-valid-v3.txt](evidence/week-03/a7-valid-v3.txt), [a7-cross-tenant-membership-v3.txt](evidence/week-03/a7-cross-tenant-membership-v3.txt) |
| 10 | Create and verify rls_verifier role (NOLOGIN, NOSUPERUSER, NOBYPASSRLS) | — | 6e8e6d8 | [rls-verifier-role-v3.txt](evidence/week-03/rls-verifier-role-v3.txt) |
| 11 | RLS behavioral v3: SELECT isolation, INSERT WITH CHECK, UPDATE USING (0), UPDATE WITH CHECK, audit immutability | — | 6e8e6d8 | [rls-select-isolation-v3.txt](evidence/week-03/rls-select-isolation-v3.txt), [rls-insert-with-check-v3.txt](evidence/week-03/rls-insert-with-check-v3.txt), [rls-update-hidden-v3.txt](evidence/week-03/rls-update-hidden-v3.txt), [rls-update-with-check-v3.txt](evidence/week-03/rls-update-with-check-v3.txt), [audit-immutability-v3.txt](evidence/week-03/audit-immutability-v3.txt) |
| 12 | RLS negative control v3: DISABLE RLS → cross-tenant visible → ROLLBACK → RLS enabled | — | 6e8e6d8 | [rls-negative-control-v3.txt](evidence/week-03/rls-negative-control-v3.txt) |
| 13 | Cleanup rls_verifier role → count 0 | — | 6e8e6d8 | [rls-verifier-cleanup-v3.txt](evidence/week-03/rls-verifier-cleanup-v3.txt) |
| 14 | Drop verification DB | — | 6e8e6d8 | (inline) |
| 15 | Main dev DB final state (historical v3, superseded) | — | 6e8e6d8 | [main-db-final-v3.txt](evidence/week-03/main-db-final-v3.txt) |
| 16 | Lint check (v3, Go toolchain not available in WSL) | — | 6e8e6d8 | [lint-v3.txt](evidence/week-03/lint-v3.txt) |
| 17 | Test check (v3, Go toolchain not available in WSL) | — | 6e8e6d8 | [test-v3.txt](evidence/week-03/test-v3.txt) |
| 18 | Reproducibility fix: add committed main-db-final-check.sql | `apps/api/testdata/week-03/main-db-final-check.sql` | d1d6adf | git log |
| 19 | Main dev DB final state (historical v4, reproducible from d1d6adf, superseded by v5) | — | d1d6adf | [main-db-final-v4.txt](evidence/week-03/main-db-final-v4.txt) |
| 20 | Lint check (v4, make lint attempted: EXIT 2; underlying golangci-lint unavailable: Error 127) | — | d1d6adf | [lint-v4.txt](evidence/week-03/lint-v4.txt) |
| 21 | Test check (v4, make test attempted: EXIT 2; underlying go unavailable: Error 127) | — | d1d6adf | [test-v4.txt](evidence/week-03/test-v4.txt) |
| 22 | Tighten audit_logs policy assertion: prove exact policy counts (total=2, select=1, insert=1, update=0, delete=0, all=0) | `apps/api/testdata/week-03/main-db-final-check.sql` | ef9ad0d | git log |
| 23 | Main dev DB final state (authoritative v5, reproducible from ef9ad0d, strict audit policy proof) | — | ef9ad0d | [main-db-final-v5.txt](evidence/week-03/main-db-final-v5.txt) |
| 24 | Lint check (v5, make lint attempted: EXIT 2; underlying golangci-lint unavailable: Error 127) | — | ef9ad0d | [lint-v5.txt](evidence/week-03/lint-v5.txt) |
| 25 | Test check (v5, make test attempted: EXIT 2; underlying go unavailable: Error 127) | — | ef9ad0d | [test-v5.txt](evidence/week-03/test-v5.txt) |
| 26 | Add pinned `golang-migrate` tooling with PostgreSQL build support and Make targets for up/down/version | `Makefile`, `.env.example`, `apps/api/testdata/week-03/` | 9c5e2e0 | [migration-cli-v3.txt](evidence/week-03/migration-cli-v3.txt), [migration-main-current-v3.txt](evidence/week-03/migration-main-current-v3.txt) |
| 27 | Add and harden deterministic transactional SQL seed fixture plus count/invariant harness | `apps/api/testdata/week-03/seed.sql`, `verify-seed-counts.sql` | 2c3a31f, e52d702, 23a8058, 0081188 | [seeder-run-1-v4.txt](evidence/week-03/seeder-run-1-v4.txt), [seeder-counts-1-v4.txt](evidence/week-03/seeder-counts-1-v4.txt), [seeder-run-2-v4.txt](evidence/week-03/seeder-run-2-v4.txt), [seeder-counts-2-v4.txt](evidence/week-03/seeder-counts-2-v4.txt) |
| 28 | Verify fresh migration lifecycle and manually-created schema baseline transition on disposable targets | `apps/api/testdata/week-03/test-fresh-migrations.sh`, `test-baseline-transition.sh` | 0081188 | [migration-fresh-cycle-v4.txt](evidence/week-03/migration-fresh-cycle-v4.txt), [migration-baseline-transition-v4.txt](evidence/week-03/migration-baseline-transition-v4.txt) |
| 29 | Verify focused negative constraints, lint state, and race-test state | `apps/api/testdata/week-03/fail_*.sql` | 0081188 | [expected-fail-a7-v3.txt](evidence/week-03/expected-fail-a7-v3.txt), [expected-fail-term-v3.txt](evidence/week-03/expected-fail-term-v3.txt), [lint-task2-v4.txt](evidence/week-03/lint-task2-v4.txt), [test-task2-v4.txt](evidence/week-03/test-task2-v4.txt) |

**Catatan implementasi:**

- **FACT:** Migration 0001, 0002, 0003, 0004 tidak diedit — sudah benar sejak f4b7ddb dan 68ab28f.
- **FACT:** Verification harness committed di `apps/api/testdata/week-03/` agar reproducible dari SHA `6e8e6d8`.
- **FACT:** Semua v3 evidence COMMAND merujuk path committed (`apps/api/testdata/week-03/*.sql` atau migration files).
- **FACT:** Fixture setup v3 EXIT 0, no ERROR.
- **FACT:** RLS scripts pakai `psql -v ON_ERROR_STOP=1` (tanpa `-1`) karena script sendiri punya BEGIN/COMMIT — tidak ada WARNING transaction noise.
- **FACT:** Old v2 evidence retained as-is; v3 supersedes non-reproducible v2 behavioral receipts.
- **FACT:** Root `tmp_*.sql` files removed.
- **FACT:** Task-2 current authoritative evidence is v4 for changed migration/seed harness, lint, and test; v3 remains authoritative only for unchanged CLI, main-current, and focused negative checks. v1/v2 receipts remain historical/superseded and are not used for current Task-2 claims.
- **FACT:** The v4 seed harness is deterministic and transactional, uses UPSERT/reconciliation, contains no `TRUNCATE`, and does not delete/rebuild the fixture state.

---

## 3. Dikerjakan Manusia

| # | Pekerjaan | Kenapa harus manual | Hasil |
|---|---|---|---|
| 1 | Sign-off laporan minggu ini | Hanya manusia yang boleh men-tick DoD dan menandatangani | **Belum — waiting for human review** |
| 2 | Jawab quiz minggu 3 | Verifikasi pemahaman tidak bisa diotomasi | **Belum — quiz week-03 belum digenerate** |
| 3 | Rekam explain-back 3 menit | Bukti pemahaman sendiri | **Belum** |
| 4 | Isi total jam aktual minggu ini | Agent tidak bisa mengukur waktu manusia | **TODO: isi di sini** |
| 5 | Spot-check 3 file bukti secara acak | Human gate / independent evidence validation | **Selesai — sampel cocok dengan klaim** |

> Catatan: Bagian ini **wajib diisi manusia**. Agent tidak menebak aktivitas manual.

---

## 4. Keputusan yang Diambil

| Keputusan | Alternatif yang ditolak | Alasan | ADR |
|---|---|---|---|
| A7 FK column order: `REFERENCES memberships(tenant_id, user_id)` | Biarkan `(user_id, tenant_id)` | **FACT:** Domain §38-A7 eksplisit: `enrollments.tenant_id → memberships.tenant_id`, `enrollments.student_user_id → memberships.user_id`. The original referenced-column ordering produced the wrong child-to-parent semantic mapping and therefore did not enforce the intended A7 invariant. | — |
| Composite FK tenant-first normalization pada semua relasi tenant→tenant | Biarkan mixed order | **FACT:** Konsistensi: tenant-first ordering keeps composite keys consistent with the domain's tenant-first relationship model and makes the tenant pairing explicit. Migration 0002 sudah pakai tenant-first. | — |
| audit_logs RLS: hanya `FOR SELECT` + `FOR INSERT` policy | Single policy `USING ... WITH CHECK ...` | **FACT:** Domain: audit_logs immutable log. Normal app roles tidak boleh punya RLS path UPDATE/DELETE. Owner/admin tetap bisa UPDATE/DELETE via bypass RLS (bukan via policy). | — |
| Migration execution atomic (`-1` flag) | Direct multi-statement psql tanpa transaction wrapper | **FACT:** Upaya sebelumnya gagal parsial karena statement tengah error tapi yang sebelumnya sudah commit. Atomic execution memastikan all-or-nothing. | — |
| Disposable verification DB sebelum main DB | Langsung test di main dev DB | **FACT:** Isolasi: verifikasi fixtures dan rls_verifier role tidak mencemari dev DB. Bisa drop/restart verification DB tanpa risiko data dev. | — |
| Migration 0004 untuk CHECK constraint (bukan edit 0002) | Edit 0002 yang sudah applied | **FACT:** Rule: never edit an already-applied migration. Write a new one. | — |
| Verification harness di `apps/api/testdata/week-03/` (bukan root atau migrations/) | Simpan di root atau di migrations/ | **FACT:** Root = untracked; migrations/ = bisa discan migration tool. testdata/ aman dan versioned. | — |

---

## 5. Angka & Bukti

> Setiap baris WAJIB punya file bukti. Tidak ada bukti → tulis `NOT MEASURED`.

| Metrik | Nilai | Cara diukur | File bukti |
|---|---|---|---|
| 0001 up v2 | 3 tables created | `CREATE TABLE` ×3 | [0001-up-v2.txt](evidence/week-03/0001-up-v2.txt) |
| 0002 up v2 | 3 tables + 3 policies | `CREATE TABLE` ×3, `CREATE POLICY` ×3 | [0002-up-v2.txt](evidence/week-03/0002-up-v2.txt) |
| 0003 up v2 | 5 tables + 1 index + 6 policies | `CREATE TABLE` ×5, `CREATE INDEX`, `CREATE POLICY` ×6 | [0003-up-v2.txt](evidence/week-03/0003-up-v2.txt) |
| 0004 up v2 | CHECK constraint added | `ALTER TABLE` | [0004-up-v2.txt](evidence/week-03/0004-up-v2.txt) |
| Total tables after up cycle | 11 | `pg_tables` count | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| Primary keys all tables | 11/11 | `pg_constraint` contype='p' | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| UNIQUE (tenant_id,id) | 8/8 tenant-scoped tables | `pg_constraint` | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| Domain unique constraints | 3/3 (memberships, course_staff, enrollments) | `pg_constraint` | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| Composite FKs tenant→tenant | 7/7 (tenant_id first) | `pg_constraint` conkey/confkey | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| A7 FK mapping verified | `conkey={2,4} confkey={2,3}` | `pg_constraint` | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| membership_roles partial index | `membership_roles_active_role_idx` ON `(membership_id, role) WHERE revoked_at IS NULL` | `pg_get_indexdef` | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| RLS enabled tenant-scoped | 8/8 tables | `relrowsecurity` | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| audit_logs policies | 2 (SELECT r, INSERT a) — no UPDATE/DELETE | `pg_policy` polcmd | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| academic_terms CHECK | `academic_terms_valid_time_range` (starts_at < ends_at) | `pg_constraint` contype='c' | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| A1-1: 2 distinct active roles | SUCCESS (EXIT 0) | INSERT lecturer + student | [a1-two-active-v3.txt](evidence/week-03/a1-two-active-v3.txt) |
| A1-2: duplicate active role | FAILED (EXIT 3, unique_violation) | INSERT lecturer again → ERROR 23505 on `membership_roles_active_role_idx` | [a1-duplicate-active-v3.txt](evidence/week-03/a1-duplicate-active-v3.txt) |
| A1-3: revoke via revoked_at | SUCCESS (EXIT 0) | UPDATE revoked_at → re-INSERT same role OK | [a1-revoke-regrant-v3.txt](evidence/week-03/a1-revoke-regrant-v3.txt) |
| A2-A: course_offerings cross-tenant | FK violation (EXIT 3) | `course_offerings_tenant_id_course_id_fkey` | [a2-course-v3.txt](evidence/week-03/a2-course-v3.txt) |
| A2-B: membership_roles cross-tenant | FK violation (EXIT 3) | `membership_roles_tenant_id_membership_id_fkey` | [a2-membership-role-v3.txt](evidence/week-03/a2-membership-role-v3.txt) |
| A2-C: course_staff cross-tenant | FK violation (EXIT 3) | `course_staff_tenant_id_course_offering_id_fkey` | [a2-course-staff-v3.txt](evidence/week-03/a2-course-staff-v3.txt) |
| A2-D: audit_logs cross-tenant | FK violation (EXIT 3) | `audit_logs_tenant_id_course_offering_id_fkey` | [a2-audit-log-v3.txt](evidence/week-03/a2-audit-log-v3.txt) |
| A2-E: enrollments cross-tenant | FK violation (EXIT 3) | `enrollments_tenant_id_course_offering_id_fkey` | [a2-enrollment-offering-v3.txt](evidence/week-03/a2-enrollment-offering-v3.txt) |
| A7-1: valid enrollment | SUCCESS (EXIT 0) | Tenant A + User A (has Tenant A membership) | [a7-valid-v3.txt](evidence/week-03/a7-valid-v3.txt) |
| A7-2: invalid enrollment | FK violation (EXIT 3) | `enrollments_tenant_id_student_user_id_fkey` | [a7-cross-tenant-membership-v3.txt](evidence/week-03/a7-cross-tenant-membership-v3.txt) |
| rls_verifier role | NOLOGIN, NOSUPERUSER, NOBYPASSRLS, not owner | `pg_roles` + `pg_class.relowner` | [rls-verifier-role-v3.txt](evidence/week-03/rls-verifier-role-v3.txt) |
| RLS SELECT isolation | Tenant A sees 1 own, 0 other (EXIT 0) | `SET LOCAL ROLE/tenant_id` + SELECT | [rls-select-isolation-v3.txt](evidence/week-03/rls-select-isolation-v3.txt) |
| RLS INSERT WITH CHECK | Cross-tenant INSERT → RLS violation (EXIT 3) | `SET LOCAL` + INSERT tenant B row as tenant A | [rls-insert-with-check-v3.txt](evidence/week-03/rls-insert-with-check-v3.txt) |
| RLS UPDATE USING | Hidden tenant B row → UPDATE 0 (EXIT 0) | `SET LOCAL` + UPDATE hidden row | [rls-update-hidden-v3.txt](evidence/week-03/rls-update-hidden-v3.txt) |
| RLS UPDATE WITH CHECK | tenant_id change A→B → RLS violation (EXIT 3) | `SET LOCAL` + UPDATE tenant_id | [rls-update-with-check-v3.txt](evidence/week-03/rls-update-with-check-v3.txt) |
| audit_logs immutability | UPDATE 0, DELETE 0, row unchanged (EXIT 0) | `SET LOCAL` + UPDATE/DELETE own row (after GRANT table privs) | [audit-immutability-v3.txt](evidence/week-03/audit-immutability-v3.txt) |
| RLS negative control | RLS disabled → cross-tenant visible (1) → ROLLBACK → RLS enabled (EXIT 0) | `ALTER TABLE DISABLE RLS` + SELECT | [rls-negative-control-v3.txt](evidence/week-03/rls-negative-control-v3.txt) |
| rls_verifier cleanup | Role count 0 (EXIT 0) | `DROP OWNED` + `DROP ROLE` | [rls-verifier-cleanup-v3.txt](evidence/week-03/rls-verifier-cleanup-v3.txt) |
| Main dev DB final state (authoritative v5, reproducible from ef9ad0d) | 11 tables, 8 RLS, A7 correct, CHECK exists, audit_policy_total=2, audit_policy_select=1, audit_policy_insert=1, audit_policy_update=0, audit_policy_delete=0, audit_policy_all=0, no rls_verifier, fixture_rows_remaining=0 | Full catalog check via stricter committed script | [main-db-final-v5.txt](evidence/week-03/main-db-final-v5.txt) |
| Lint (v5, literal `make lint`) | NOT VERIFIED — tool unavailable (EXIT 2 from make; underlying golangci-lint Error 127) | `golangci-lint` not installed in WSL | [lint-v5.txt](evidence/week-03/lint-v5.txt) |
| Test (v5, literal `make test`) | NOT VERIFIED — tool unavailable (EXIT 2 from make; underlying go Error 127) | `go` not installed in WSL | [test-v5.txt](evidence/week-03/test-v5.txt) |
| Migration CLI | `golang-migrate` v4.18.3, built with PostgreSQL support | `go version -m /go/bin/migrate` in API container | [migration-cli-v3.txt](evidence/week-03/migration-cli-v3.txt) |
| Main DB migration state | version 4, `dirty=false`, force skipped, `migrate up=no_change`, structural postcheck pass | guarded baseline script | [migration-main-current-v3.txt](evidence/week-03/migration-main-current-v3.txt) |
| Fresh disposable migration cycle | 0 → 4 → 3 → 4; CHECK absent at 3, restored at 4; structural check and cleanup pass | `test-fresh-migrations.sh` | [migration-fresh-cycle-v4.txt](evidence/week-03/migration-fresh-cycle-v4.txt) |
| Disposable baseline transition | manually-created v4 schema: `schema_migrations` absent → version 4 clean; no-change up; target-isolated cleanup | `test-baseline-transition.sh` | [migration-baseline-transition-v4.txt](evidence/week-03/migration-baseline-transition-v4.txt) |
| Deterministic fixture counts | tenants=3; lecturers=50; students=2000; users=2050; auth identities=2050; memberships=2050; active roles=2050; lecturer roles=50; student roles=2000; academic terms=6; courses=200; course offerings=400; course staff=400; enrollments=20000 | seed then strict count/invariant SQL | [seeder-counts-1-v4.txt](evidence/week-03/seeder-counts-1-v4.txt) |
| Second seed-run state | Same 14 fixture counts; all 14 reported invariant violation metrics=0, including active instructor `course_staff` and same-tenant active student/lecturer membership/role checks | second seed then strict count/invariant SQL | [seeder-run-2-v4.txt](evidence/week-03/seeder-run-2-v4.txt), [seeder-counts-2-v4.txt](evidence/week-03/seeder-counts-2-v4.txt) |
| Focused negative constraints | Cross-tenant enrollment rejected by `enrollments_tenant_id_student_user_id_fkey`; invalid academic term rejected by `academic_terms_valid_time_range` | expected-failure SQL | [expected-fail-a7-v3.txt](evidence/week-03/expected-fail-a7-v3.txt), [expected-fail-term-v3.txt](evidence/week-03/expected-fail-term-v3.txt) |
| Lint Task-2 v4 | NOT VERIFIED — `make lint` EXIT 2; underlying `golangci-lint` unavailable / Make Error 127 | `make lint` | [lint-task2-v4.txt](evidence/week-03/lint-task2-v4.txt) |
| Test Task-2 v4 | PASS — `make test` EXIT 0; `go test -race -count=1 ./...` passes | `make test` | [test-task2-v4.txt](evidence/week-03/test-task2-v4.txt) |

**Perbandingan sebelum/sesudah (koreksi migration 0003 dari f4b7ddb):**

| Metrik | Sebelum (migration 0003 original) | Sesudah (migration 0003 corrected) | Perubahan | Bukti |
|---|---|---|---|---|
| A7 FK referenced columns | `(user_id, tenant_id)` | `(tenant_id, user_id)` | Column order fixed per domain §38-A7 | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| membership_roles FK | `(membership_id, tenant_id) REFERENCES (id, tenant_id)` | `(tenant_id, membership_id) REFERENCES (tenant_id, id)` | Tenant-first normalized | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| course_staff FK | `(course_offering_id, tenant_id) REFERENCES (id, tenant_id)` | `(tenant_id, course_offering_id) REFERENCES (tenant_id, id)` | Tenant-first normalized | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |
| audit_logs policies | 1 policy `*` (ALL commands) | 2 policies: `r` (SELECT), `a` (INSERT) | Immutable: no UPDATE/DELETE path | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) |

---

## 6. Konsep yang Dipelajari

### Composite Foreign Key + RLS — Defense in Depth

- **Apa:** Foreign key yang memasukkan `tenant_id` sebagai kolom pertama dalam referensi composite, dipadukan dengan Row Level Security policy menggunakan `current_setting('app.tenant_id')`.
- **Kenapa dipakai di sini:** Domain §38-A2 mewajibkan konsistensi tenant di level database, bukan hanya service layer. Composite FK mencegah cross-tenant reference struktural (DB menolak INSERT/UPDATE), RLS menyembunyikan row dari query tenant lain. Keduanya komplementer: FK = structural guarantee, RLS = query isolation.
- **Alternatif yang tidak dipilih:** Hanya RLS tanpa composite FK tenant-aware. Biaya: referensi silang antar tenant yang tidak valid bisa berhasil disimpan saat write time. RLS mungkin menyembunyikan row saat read normal, tapi tidak memberikan jaminan referensial struktural yang diberikan composite FK.
- **Cara membuktikan sendiri:** `cat apps/api/migrations/0003_auth_membership_schema.up.sql | grep -A2 "FOREIGN KEY (tenant_id"` — lihat semua composite FK pakai tenant-first. Atau jalankan A2 test v3.
- **Pertanyaan interview terkait:** "Kenapa butuh composite FK kalau sudah ada RLS? Apa yang terjadi kalau FK-nya composite tapi column order-nya salah?"

### A7 Enrollment → Membership Referential Integrity

- **Apa:** Foreign key dari `enrollments(tenant_id, student_user_id)` ke `memberships(tenant_id, user_id)` yang memastikan mahasiswa hanya bisa enroll ke course offering jika dia memiliki membership **di tenant yang sama**. FK ini menjamin **keberadaan** membership, **bukan** status aktif. Status aktif (revoked_at IS NULL) divalidasi service layer (ini adalah bagian dari membership_roles role-history semantics).
- **Kenapa dipakai di sini:** Domain §38-A7: "Enrollment wajib menunjuk membership". Keberadaan membership di tenant yang sama ditegakkan database via FK; status aktif divalidasi service layer kecuali ada constraint DB tambahan di masa depan.
- **Alternatif yang tidak dipilih:** Validasi keberadaan membership di application layer saja. Biaya: race condition, bypass via direct DB access, inconsistent state kalau migration/service bug.
- **Cara membuktikan sendiri:** Lihat `pg_constraint` untuk `enrollments_tenant_id_student_user_id_fkey` — `conkey={2,4} confkey={2,3}` berarti `enrollments.tenant_id(2)→memberships.tenant_id(2)`, `enrollments.student_user_id(4)→memberships.user_id(3)`. Original bug memiliki mapping yang salah (tenant_id dipetakan ke user_id dan sebaliknya), sehingga mapping yang diperbaiki (tenant_id ke tenant_id, student_user_id ke user_id) sekarang enforce domain invariant dengan benar.
- **Pertanyaan interview terkait:** "Bagaimana memastikan enrollment tidak bisa reference user yang bukan member tenant tsb? Apa bedanya FK ini dengan FK biasa? Apakah FK ini juga memvalidasi status aktif membership?"

### audit_logs Immutability via RLS Policy Scope

- **Apa:** Tabel `audit_logs` hanya memiliki RLS policy `FOR SELECT` dan `FOR INSERT` untuk normal tenant access. Tidak ada policy `FOR UPDATE` atau `FOR DELETE`. Owner/admin (bypass RLS) tetap bisa mutate, tapi normal app role tidak punya RLS path.
- **Kenapa dipakai di sini:** Domain: audit_logs immutable log. Normal application flow tidak boleh ubah/hapus audit trail. RLS policy scope (`FOR SELECT`/`FOR INSERT`) membatasi command yang diizinkan via RLS.
- **Alternatif yang tidak dipilih:** Trigger `BEFORE UPDATE/DELETE` yang `RAISE EXCEPTION`. Biaya: Sebuah trigger menambah satu lagi layer mekanisme penegakan, dan privileged maintenance akan membutuhkan strategi bypass yang diatur secara eksplisit dan hati-hati.
- **Cara membuktikan sendiri:** `SELECT polname, polcmd FROM pg_policy WHERE polrelid = 'audit_logs'::regclass;` → harus return `r` dan `a` saja. Test UPDATE/DELETE sebagai `rls_verifier` role (dengan GRANT table UPDATE/DELETE) → `UPDATE 0` / `DELETE 0`.
- **Pertanyaan interview terkait:** "Gimana cara bikin tabel audit log immutable di Postgres tanpa trigger? Apa bedanya `CREATE POLICY ... FOR SELECT` vs tanpa `FOR`?"

### SET LOCAL vs SET — Connection Pool Safety

- **Apa:** `SET LOCAL app.tenant_id = '...'` di dalam transaksi eksplisit (`BEGIN`...`COMMIT`), bukan `SET` di level session.
- **Kenapa dipakai di sini:** `agent/rules/20-database.md` §25: connection pooling menyebabkan state session-level (seperti yang di-set oleh `SET`) bertahan pada scope koneksi. Jika sebuah pooled connection digunakan ulang oleh request berikutnya, dan aplikasi gagal melakukan reset/overwrite setting tenant di semua code path, request berikutnya bisa mewarisi state tenant lama secara tidak sengaja. `SET LOCAL` otomatis berakhir saat transaksi selesai, mencegah kebocoran context ke request lain.
- **Alternatif yang tidak dipilih:** `SET app.tenant_id` tanpa `LOCAL`. Biaya: potensi kebocoran context tenant jika aplikasi mengalami bug saat menggunakan connection pool.
- **Cara membuktikan sendiri:** Semua verifikasi RLS di evidence pakai `BEGIN; SET LOCAL ROLE ...; SET LOCAL app.tenant_id ...; ... COMMIT;`. 
- **Pertanyaan interview terkait:** "Mengapa `SET LOCAL` wajib untuk multi-tenant RLS dengan connection pool? Apa yang terjadi kalau pakai `SET` biasa?"

### Table Owner Bypasses RLS — Need Non-Owner Verifier

- **Apa:** Role yang create table (owner) secara default **bypass RLS** untuk tabel tersebut (kecuali `FORCE ROW LEVEL SECURITY` dipakai). Testing RLS dengan role owner akan selalu pass (tidak detect misconfiguration).
- **Kenapa dipakai di sini:** Verifikasi RLS harus pakai role non-owner (`rls_verifier` NOLOGIN NOBYPASSRLS) yang di-GRANT privileges minimal. Ini memastikan RLS benar-benar enforce.
- **Alternatif yang tidak dipilih:** Test pakai role `campus` (owner). Biaya: false positive — RLS kelihatan kerja tapi sebenarnya bypassed.
- **Cara membuktikan sendiri:** Negative control test: `ALTER TABLE audit_logs DISABLE ROW LEVEL SECURITY; SET LOCAL ROLE rls_verifier; SELECT * FROM audit_logs WHERE tenant_id = 'tenant-b';` → row jadi visible. Rollback → RLS enabled lagi.
- **Pertanyaan interview terkait:** "Kalau test RLS pakai user superuser/owner, apa risikonya? Gimana cara bikin verifier role yang benar?"

---

## 7. Belum Terverifikasi

| Hal | Kenapa belum terverifikasi | Rencana verifikasi |
|---|---|---|
| Repository layer Go (pgx) + request/transaction tenant context | Belum diimplementasikan. Roadmap deliverable #3 dan #6 mensyaratkan context disetel per-request di transaksi. | Current Week 3 task loop |
| Course offering + participants endpoint/query dan bukti N+1 | Belum ada repository/endpoint atau query-count proof. | Current Week 3 task loop |
| EXPLAIN (ANALYZE, BUFFERS) query tuning (5 query) | Seeder sudah tersedia, tetapi analisis before/after dan tuning belum dikerjakan. | Current Week 3 task loop |
| Backup/restore drill (`make db-restore`) | Belum ada script backup/restore atau bukti restore actual. | Current Week 3 task loop |
| ADR-0002 multi-tenancy + catatan Supabase vs self-managed | Dokumentasi roadmap deliverable #8 belum ditulis. | Current Week 3 task loop |
| Integration-level RLS verification pada path Go | Bukti SQL RLS ada, tetapi integrasi request/transaction Go belum ada. | Setelah repository + tenant context tersedia, masih di Week 3 |
| Total jam aktual minggu ini | Agent tidak bisa mengukur waktu manusia | Human fill |
| Final lint verification | Bukti Task-2 v4: `make lint` EXIT 2 karena `golangci-lint` tidak tersedia (Make Error 127). | Sediakan/jalankan `golangci-lint`, lalu capture evidence pada task loop baru |
| Quiz, explain-back, dan sign-off manusia | Quiz dan explain-back masih kosong; laporan belum ditandatangani. | Human gate setelah pekerjaan Week 3 lain selesai |

**Asumsi yang dipakai tapi belum dibuktikan:**

- **INFERENCE:** Migration tooling dan deterministic seeder sudah diverifikasi secara lokal melalui evidence Task-2; perilaku pada database produksi belum diamati.
- **INFERENCE:** `auth_sessions` (Tier 1 Week 4) belum ada — migration 0005 belum dibuat (0004 dipakai untuk CHECK constraint).
- **INFERENCE:** Neon production database belum diprovision — Minggu 4 task.

---

## 8. Masalah & Cara Diselesaikan

### Masalah: Previous behavioral evidence (non-v2) was INVALID — false PASS claims
- **Gejala:** Evidence files `a1-verification.txt`, `a2-verification.txt`, `a7-verification.txt`, `rls-verification.txt`, `rls-negative-control.txt`, `cleanup-verifier.txt` mencetak "SUCCESS" / "UNEXPECTED SUCCESS" tapi RAW OUTPUT menunjukkan: fixture Tenant A missing (FK violation on ordinary `tenant_id`), `rls_verifier` role not exist (transaction aborted), psql continued after errors.
- **Hipotesis yang salah:** "Label SUCCESS di output SQL cukup sebagai bukti" — salah, command EXIT code dan RAW OUTPUT adalah bukti.
- **Akar masalah:** (1) Fixture setup gagal (Tenant A UUID `aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa` tidak di-INSERT dulu). (2) `rls_verifier` role tidak dibuat sebelum RLS test. (3) `psql` tanpa `ON_ERROR_STOP=1` → transaction aborted tapi statement berikutnya jalan → output kacau. (4) Evidence manual reconstruction bukan raw capture.
- **Solusi:** Correction pass v2: fixture setup v2 (semua INSERT berhasil), create rls_verifier role dulu, `psql -v ON_ERROR_STOP=1 -1`, evidence via `make evidence` raw capture. Bukti lama DISIMPAN sebagai historical failed evidence (tidak dihapus, tidak dikutip sebagai PASS).
- **Pencegahan:** Selalu verify fixture setup EXIT 0 sebelum behavioral test. Selalu create verifier role. Selalu `ON_ERROR_STOP=1`. Selalu `make evidence`.
- **Waktu terbuang:** NOT MEASURED (diagnosis, correction pass, re-verification).

### Masalah: Invalid evidence file `0001-up.txt` manually reconstructed (dari recovery pass sebelumnya)
- **Gejala:** File `docs/progress/evidence/week-03/0001-up.txt` berisi output yang tidak match format evidence protocol.
- **Solusi:** Sudah diganti dengan `0001-up-authentic.txt` (recovery pass) dan sekarang `0001-up-v2.txt` (correction pass).

### Masalah: Migration 0003 A7 FK column order reversed (diperbaiki di f4b7ddb)
- **Gejala:** `FOREIGN KEY (tenant_id, student_user_id) REFERENCES memberships (user_id, tenant_id)` — column order terbalik.
- **Solusi:** Sudah dikoreksi di commit f4b7ddb ke `REFERENCES memberships (tenant_id, user_id)`. Diverifikasi via `pg_constraint.conkey/confkey` di v2.

### Masalah: v2 behavioral evidence non-reproducible (references untracked `tmp_*.sql`)
- **Gejala:** Banyak v2 evidence COMMAND merujuk `tmp_a1_test.sql`, `tmp_rls_select.sql`, dll di root yang tidak di-commit. COMMIT recorded sebagai `68ab28f` tapi file tidak ada di SHA tersebut.
- **Akar masalah:** Test scripts dibuat ad-hoc di root selama v2 run, tidak di-commit sebelum evidence capture.
- **Solusi:** Pass ini memindahkan semua 19 script ke `apps/api/testdata/week-03/` dengan nama stabil, commit sebagai SHA `6e8e6d8`, lalu regenerate v3 evidence yang COMMAND-nya merujuk path committed. Root `tmp_*.sql` dihapus.
- **Pencegahan:** Selalu commit test harness SEBELUM capture evidence. Gunakan `make evidence` dengan CMD yang merujuk path tracked.

### Masalah: RLS v2 evidence transaction wrapper noise
- **Gejala:** v2 RLS evidence menggunakan `psql -1` (single-transaction) tapi SQL script juga punya BEGIN/COMMIT → WARNING "there is already a transaction in progress" / "there is no transaction in progress".
- **Solusi:** v3 RLS scripts menggunakan `psql -v ON_ERROR_STOP=1` (tanpa `-1`) karena script sendiri manage transaksi. Migration files tetap pakai `psql -v ON_ERROR_STOP=1 -1`.

### Masalah: Superseded final-main-DB evidence versions
- **Gejala:** Laporan sebelumnya merujuk `main-db-final-v3.txt` dan `main-db-final-v4.txt`, tetapi script awal tidak membuktikan ketiadaan UPDATE/DELETE/ALL policy secara eksplisit dan memiliki keterbatasan provenans reproduktibilitas terhadap snapshot commit.
- **Akar masalah:** Skrip verifikasi v3/v4 awal hanya menghitung policy dengan `polcmd IN ('r', 'a')` tanpa assertion eksplisit untuk breakdown setiap tipe `polcmd`.
- **Solusi:** Perbaiki skrip verifikasi (`main-db-final-check.sql`) untuk memecah count per `polcmd` (`r`, `a`, `w`, `d`, `*`) dan membuktikan total=2, select=1, insert=1, update=0, delete=0, all=0. Commit perbaikan sebagai SHA `ef9ad0d`. Regenerate evidence sebagai `main-db-final-v5.txt` (authoritative).
- **Status evidence:**
  - v3 = superseded / historical
  - v4 = superseded oleh assertion strict v5
  - v5 = authoritative current final-main-DB receipt (COMMIT=ef9ad0d)

### Masalah: Task-2 evidence memiliki beberapa generasi receipt
- **Gejala:** v1 dan v2 merekam implementasi/harness sebelum hardening akhir; beberapa header claim masih placeholder atau run lama melakukan delete/rebuild. Memakainya sebagai bukti current akan mencampur state yang sudah digantikan.
- **Akar masalah:** Tooling migrasi dan seeder diperbaiki bertahap melalui lima commit implementasi Task-2.
- **Solusi:** Raw receipt tidak diubah atau dihapus. Untuk klaim Task-2 current, v4 adalah authoritative untuk fresh cycle, baseline transition, seed run/count, lint, dan test; v3 authoritative hanya bagi source yang tidak berubah: CLI, main-current, serta focused negative checks.
- **Pencegahan:** Laporan menyebut generasi evidence dan commit yang tepat; receipt historical/superseded tetap dipertahankan sebagai audit trail.
- **Waktu terbuang:** NOT MEASURED.

---

## 9. Status Definition of Done

> Agent hanya MENGUSULKAN. Kolom "Dicentang manusia" diisi olehmu setelah melihat bukti.

| DoD dari roadmap | Usulan agent | Bukti | Dicentang manusia |
|---|---|---|---|
| Query sebagai Tenant A tidak dapat membaca atau memodifikasi row Tenant B, dibuktikan integration test RLS | ✅ terpenuhi | [rls-select-isolation-v3.txt](evidence/week-03/rls-select-isolation-v3.txt), [rls-insert-with-check-v3.txt](evidence/week-03/rls-insert-with-check-v3.txt), [rls-update-hidden-v3.txt](evidence/week-03/rls-update-hidden-v3.txt), [rls-update-with-check-v3.txt](evidence/week-03/rls-update-with-check-v3.txt) | ☐ |
| Composite FK benar-benar menolak upaya referensi silang tenant | ✅ terpenuhi | [a2-course-v3.txt](evidence/week-03/a2-course-v3.txt), [a2-membership-role-v3.txt](evidence/week-03/a2-membership-role-v3.txt), [a2-course-staff-v3.txt](evidence/week-03/a2-course-staff-v3.txt), [a2-audit-log-v3.txt](evidence/week-03/a2-audit-log-v3.txt), [a2-enrollment-offering-v3.txt](evidence/week-03/a2-enrollment-offering-v3.txt) | ☐ |
| Satu membership dapat memegang dua role aktif yang berbeda; duplikasi role aktif yang sama ditolak; pencabutan menggunakan `revoked_at`, bukan `DELETE` | ✅ terpenuhi | [a1-two-active-v3.txt](evidence/week-03/a1-two-active-v3.txt), [a1-duplicate-active-v3.txt](evidence/week-03/a1-duplicate-active-v3.txt), [a1-revoke-regrant-v3.txt](evidence/week-03/a1-revoke-regrant-v3.txt) | ☐ |
| Tidak ada tabel Minggu 5/6 atau Tier 3 yang muncul prematur pada migrasi Minggu 3 | ✅ terpenuhi | [catalog-checks-v2.txt](evidence/week-03/catalog-checks-v2.txt) — 11 tables exactly | ☐ |
| Minimal satu query yang terbukti lambat/inefisien pada dataset seeder diperbaiki dan mempunyai `EXPLAIN` before/after | ❌ belum | Seeder Task-2 tersedia, tetapi belum ada `EXPLAIN (ANALYZE, BUFFERS)` before/after atau tuning | ☐ |
| Endpoint daftar **course offering + peserta** bebas N+1, dibuktikan dengan hitungan query di test | ❌ belum | Belum ada repository/endpoint | ☐ |
| `make db-restore` berhasil memulihkan backup ke database kosong dengan data lengkap | ❌ belum | Belum ada backup script | ☐ |
| Kamu bisa menjelaskan read committed vs serializable memakai contoh dari skema LMS-mu sendiri | ⚠️ sebagian | Dijelaskan di konsep tapi belum demo konkret — **pending human explain-back / human verification** | ☐ |
| Kamu bisa menjelaskan kenapa `courses` tanpa `course_offerings` akan merusak histori saat masuk semester berikutnya | ⚠️ sebagian | Dijelaskan di konsep "Composite FK + RLS" — **pending human explain-back / human verification** | ☐ |

---

## 10. Untuk Minggu Depan

- **Carry-over:** Tidak ada keputusan carry-over; semua item berikut tetap pekerjaan current Week 3.
- **Remaining current Week 3 work:** Repository layer pgx; request/transaction tenant context; course-offering + participants endpoint/query; N+1 query-count proof; `EXPLAIN (ANALYZE, BUFFERS)` dan tuning; backup + restore drill / `make db-restore`; ADR/documentation roadmap; integration-level RLS pada path Go; final lint verification; quiz; explain-back; dan human gate/sign-off.
- **Utang teknis yang sengaja diambil:** Tidak ada utang migration tooling atau seeder dari Task-2; keduanya telah melalui closeout evidence. Lint tetap belum terverifikasi karena command dependency tidak tersedia.
- **Persiapan yang perlu dilakukan manusia lebih dulu:** Azure VM provisioning (Minggu 4), Neon database setup (Minggu 4), Cloudflare Pages untuk frontend (Minggu 4).
- **Week 3 tetap terbuka** sampai gate complete (laporan ditandatangani, quiz ≥ 70%, explain-back terekam). Jangan pindahkan DoD belum selesai ke Minggu 4.

---

## 11. Verifikasi Manusia

- [X] Saya sudah spot-check 3 file bukti migration-foundation secara acak dan isinya cocok dengan klaim
- [ ] Skor quiz: ____ / ____ (minimal 70% untuk lanjut) — **quiz week-03 belum digenerate**
- [ ] Explain-back 3 menit sudah direkam: `docs/progress/explain/week-03.<mp3|txt>`
- [ ] Saya bisa menjelaskan setiap keputusan di section 4 tanpa membuka catatan

**Self-audit agent (dari `agent/evidence-protocol.md` §8):** Sudah dijalankan untuk closeout dokumentasi: angka Task-2 merujuk receipt v3/v4 yang ada, raw evidence tidak diubah, v1/v2 dipertahankan sebagai historical/superseded, dan klaim lint/test mengikuti EXIT receipt terbaru (`make lint` NOT VERIFIED; `make test -race` PASS). Agent tidak menjalankan ulang command pada pass ini, tidak men-tick human box, tidak mengubah quiz, dan tidak menandatangani laporan. Spot-check yang sudah tercatat di atas berasal dari migration-foundation sebelumnya; belum ada klaim bahwa manusia telah spot-check evidence Task-2 baru.

**Ditandatangani:** ______________  **Tanggal:** __________
