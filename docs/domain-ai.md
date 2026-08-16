# Domain Model — Lapisan AI (`campus-lms`)

> Dokumen pendamping `docs/domain.md`. Sengaja dipisah agar batas antara
> **domain pembelajaran** dan **lapisan AI** tetap tegas.
>
> Aturan pemisahan: `domain.md` harus tetap utuh dan bermakna **tanpa** dokumen
> ini. LMS wajib berfungsi penuh kalau seluruh fitur AI dimatikan.

---

## 0. Tujuan dan batas

Lapisan AI menambahkan kemampuan berbasis LLM di atas domain pembelajaran yang
sudah ada: peringkasan materi, tanya-jawab berbasis materi kuliah dengan sitasi,
pembuatan draft kuis, bantuan penilaian, dan agent perencana belajar.

### 0.1 Prinsip yang tidak bisa ditawar

1. **AI tidak pernah menjadi source of truth akademik.** Tidak ada nilai,
   kehadiran, atau status kelulusan yang ditentukan model. AI hanya
   menghasilkan *draft* atau *bantuan* yang harus melewati manusia.
2. **Semua keluaran AI wajib ditandai sebagai AI-generated**, tersimpan bersama
   `model`, `provider`, dan `prompt_version` yang menghasilkannya.
3. **Konten AI tidak terlihat mahasiswa sebelum ditinjau course staff.** Ini isu
   integritas akademik, bukan preferensi UX.
4. **Konten dokumen yang di-retrieve adalah data, bukan instruksi.** PDF yang
   diunggah dosen dapat memuat prompt injection dan harus diperlakukan sebagai
   input tak tepercaya.
5. **Isolasi tenant berlaku penuh pada embedding.** Kebocoran retrieval lintas
   tenant adalah pelanggaran data, bukan bug fitur.
6. **Setiap panggilan model tercatat biayanya.** Tanpa pencatatan, kuota gratis
   habis tanpa jejak dan tidak ada dasar optimasi.
7. **Fitur AI dapat dimatikan per tenant** melalui feature flag di
   `tenant_settings`. Kampus berhak menolak pemrosesan AI.

### 0.2 Di luar lingkup

* Fine-tuning atau training model.
* Menyimpan bobot model.
* Proctoring otomatis / deteksi kecurangan berbasis AI.
* Deteksi plagiarisme (disediakan sebagai integration boundary di `domain.md` §33).
* Keputusan akademik otomatis tanpa manusia.

---

## 1. Entitas

Seluruh entitas di bawah **tenant-scoped** dan wajib RLS.

### 1.1 `ai_artifacts`

Keluaran AI yang berumur panjang dan ditampilkan di produk: ringkasan materi,
learning objective, draft kuis, draft feedback penilaian.

| Atribut | Catatan |
|---|---|
| `id` | |
| `tenant_id` | |
| `course_offering_id` | nullable, untuk artefak level course |
| `subject_type` | `material` \| `lesson` \| `submission` \| `course_offering` |
| `subject_id` | id entitas sumber |
| `artifact_type` | `summary` \| `learning_objectives` \| `quiz_draft` \| `grading_feedback` |
| `content` | JSON tervalidasi skema |
| `model` | mis. `gemini-2.5-flash` |
| `provider` | mis. `google`, `groq`, `cerebras` |
| `prompt_version` | menunjuk file prompt versioned |
| `input_hash` | hash konten sumber; deteksi kadaluarsa saat sumber berubah |
| `review_status` | `draft` \| `approved` \| `rejected` |
| `reviewed_by` | user_id course staff |
| `reviewed_at` | |
| `published_at` | hanya setelah `approved` |
| `created_at` | |

Constraint:

```text
UNIQUE (tenant_id, subject_type, subject_id, artifact_type, prompt_version)
CHECK  (published_at IS NULL OR review_status = 'approved')
```

Lifecycle:

```text
draft ──► approved ──► published
  └────► rejected
```

Jika `input_hash` tidak lagi cocok dengan sumbernya, artefak ditandai kadaluarsa
dan tidak boleh ditampilkan sebagai representasi materi terkini.

---

### 1.2 `material_chunks`

Potongan materi beserta embedding-nya untuk RAG.

| Atribut | Catatan |
|---|---|
| `id` | |
| `tenant_id` | wajib eksplisit meski dapat diturunkan dari `material_id` |
| `material_id` | FK ke `materials` |
| `course_offering_id` | denormalisasi untuk filter retrieval yang cepat dan aman |
| `chunk_index` | urutan dalam dokumen |
| `content` | teks potongan |
| `token_count` | |
| `page` | nullable — dasar sitasi |
| `section` | nullable — heading/slide |
| `content_hash` | untuk re-embed inkremental |
| `embedding_model` | mis. `multilingual-e5-small` |
| `embedding` | `vector(384)` — dimensi mengikuti model |
| `created_at` | |

Constraint:

```text
UNIQUE (material_id, chunk_index, embedding_model)
```

Aturan:

* Index HNSW; parameter `m`, `ef_construction`, `ef_search` dicatat di ADR.
* RLS wajib. Retrieval **tidak boleh** hanya mengandalkan filter aplikasi.
* Ganti model embedding = baris baru, bukan menimpa. Dua model boleh hidup
  berdampingan selama migrasi.
* Materi diubah → hanya chunk dengan `content_hash` berbeda yang di-embed ulang.
* Materi dihapus/di-unpublish → chunk-nya tidak boleh lagi ter-retrieve.

---

### 1.3 `ai_interactions`

Satu baris per pemanggilan model. Sumber kebenaran biaya, latensi, dan kuota.

| Atribut | Catatan |
|---|---|
| `id` | |
| `tenant_id` | |
| `user_id` | nullable untuk job sistem |
| `course_offering_id` | nullable |
| `feature` | `summary` \| `ask_material` \| `quiz_gen` \| `grading_assist` \| `agent` |
| `provider` / `model` | |
| `prompt_version` | |
| `tokens_in` / `tokens_out` | |
| `cost_estimate` | `numeric` — estimasi, bukan tagihan |
| `latency_ms` | |
| `cache_hit` | `none` \| `exact` \| `semantic` |
| `status` | `ok` \| `error` \| `rate_limited` \| `fallback` |
| `fallback_from` | provider yang gagal, nullable |
| `trace_id` | korelasi ke OpenTelemetry |
| `created_at` | |

Aturan:

* Isi prompt dan jawaban **tidak** disimpan di tabel ini. Konten mengalir ke
  Langfuse; tabel ini menyimpan metadata saja.
* Retensi bergilir — data operasional, bukan catatan akademik.
* Semua angka di laporan mingguan yang menyangkut biaya/latensi AI harus dapat
  ditelusuri ke tabel ini.

---

### 1.4 `agent_actions`

Audit setiap pemanggilan tool oleh agent. Bukan log, melainkan jejak audit.

| Atribut | Catatan |
|---|---|
| `id` | |
| `tenant_id` | |
| `user_id` | atas nama siapa agent bertindak |
| `session_id` | mengelompokkan langkah dalam satu tugas |
| `step_index` | |
| `tool_name` | |
| `arguments` | JSON, PII teredaksi |
| `result_summary` | ringkas, bukan payload penuh |
| `status` | `proposed` \| `approved` \| `executed` \| `rejected` \| `failed` |
| `requires_approval` | boolean |
| `approved_by` / `approved_at` | |
| `tokens_used` / `cost_estimate` | |
| `trace_id` | |
| `created_at` | |

Aturan:

* Tool berdampak (mengubah nilai, mengirim pengumuman massal, menghubungi
  mahasiswa) **wajib** `requires_approval = true` dan tidak boleh
  `executed` tanpa `approved_by`.
* Batas keras per sesi: maksimum langkah, maksimum token, maksimum waktu.
  Terlampaui → berhenti dengan pesan yang jelas, bukan diam.
* Immutable. Sama seperti `audit_logs`, tidak boleh diedit aplikasi.

---

### 1.5 `ai_quotas`

Anggaran token per tenant per periode. Mencegah satu tenant menghabiskan kuota
gratis bersama.

| Atribut | Catatan |
|---|---|
| `id`, `tenant_id` | |
| `period_start` / `period_end` | |
| `token_budget` / `tokens_used` | |
| `request_budget` / `requests_used` | |
| `hard_limit` | `true` = tolak; `false` = degradasi ke model murah |
| `updated_at` | |

Constraint:

```text
UNIQUE (tenant_id, period_start)
```

---

### 1.6 `ai_feedback`

Umpan balik pengguna terhadap keluaran AI. Sumber bahan golden dataset di
Minggu 9 — kasus yang dinilai buruk oleh manusia adalah kandidat test terbaik.

| Atribut | Catatan |
|---|---|
| `id`, `tenant_id`, `user_id` | |
| `ai_interaction_id` | |
| `rating` | `helpful` \| `not_helpful` \| `wrong` \| `unsafe` |
| `comment` | nullable |
| `created_at` | |

---

## 2. Aturan bisnis

### 2.1 Integritas akademik

1. Keluaran AI selalu ditandai AI-generated di UI.
2. Ringkasan, learning objective, dan kuis tidak terlihat mahasiswa sebelum
   `review_status = approved`.
3. Grading assistant hanya menghasilkan **draft**. Nilai final tetap tunduk pada
   `domain.md` §24.9 — hanya Lead Instructor yang publish/lock.
4. Setiap nilai yang berasal dari usulan AI menyimpan jejak: artefak sumber,
   siapa yang menyetujui, dan apakah nilainya diubah dari usulan.
5. AI tidak boleh mengakses submission mahasiswa lain saat menilai satu
   submission.

### 2.2 Retrieval dan sitasi

1. Jawaban wajib membawa sitasi yang dapat diverifikasi (`material_id`, `page`).
2. Skor retrieval di bawah ambang → **menolak menjawab**. Dalam produk
   pendidikan, jawaban salah yang meyakinkan lebih berbahaya daripada
   "tidak ditemukan di materi".
3. Retrieval hanya menjangkau materi yang: satu tenant, course offering yang
   di-enroll pengguna, dan sudah `published`.
4. Materi yang di-unpublish langsung hilang dari hasil retrieval.

### 2.3 Keamanan

1. Konten dokumen adalah data, bukan instruksi. System prompt menyatakan
   eksplisit bahwa isi materi tidak dapat menimpa instruksi.
2. PII mahasiswa diredaksi sebelum dikirim ke provider, kecuali ada alasan
   terdokumentasi.
3. Keluaran model tidak pernah dieksekusi, tidak dirender sebagai HTML mentah,
   dan tidak diteruskan ke shell.
4. Tool agent tunduk allow-list per peran.
5. Rate limit dan kuota berlaku per tenant **dan** per pengguna.

### 2.4 Biaya

1. Setiap panggilan tercatat di `ai_interactions`. Tidak tercatat = dianggap bug.
2. Cache diperiksa sebelum memanggil model.
3. Kuota habis → degradasi ke model murah atau tolak dengan pesan jelas; jangan
   pernah menggantung diam-diam.

---

## 3. Klasifikasi RLS

Tenant-scoped, wajib RLS:

```text
ai_artifacts
material_chunks
ai_interactions
agent_actions
ai_quotas
ai_feedback
```

Berlaku pola composite FK yang sama seperti `domain.md` §38-A2, sehingga
referensi lintas tenant mustahil secara struktural.

---

## 4. Invariant database

```text
UNIQUE (tenant_id, subject_type, subject_id, artifact_type, prompt_version)
    ai_artifacts

UNIQUE (material_id, chunk_index, embedding_model)
    material_chunks

UNIQUE (tenant_id, period_start)
    ai_quotas

CHECK (published_at IS NULL OR review_status = 'approved')
    ai_artifacts

CHECK (status <> 'executed' OR requires_approval = false OR approved_by IS NOT NULL)
    agent_actions

CHECK (tokens_in >= 0 AND tokens_out >= 0 AND cost_estimate >= 0)
    ai_interactions
```

---

## 5. Pemetaan ke roadmap

| Minggu | Entity yang dibangun | Fitur |
|---|---|---|
| 7 | `ai_interactions`, `ai_artifacts` | Auto-summary materi, LLM router, cost tracking |
| 8 | `material_chunks` | RAG "Tanya Materi" dengan sitasi, draft kuis |
| 9 | `ai_feedback` | Golden dataset, eval harness, quality gate di CI |
| 10 | `agent_actions` | Study-plan agent, MCP server, guardrail, red-team |
| 12 | `ai_quotas` | Rate limit & kuota token per tenant |

---

## 6. Yang sengaja tidak dilakukan

* Tidak menyimpan isi prompt/jawaban di database transaksional — itu tugas
  Langfuse, dan memisahkannya menjaga tabel akademik tetap ramping.
* Tidak membuat tabel eval di database. Golden dataset hidup di git supaya
  ter-version bersama kode dan dapat di-review lewat PR.
* Tidak ada fitur AI yang menulis langsung ke `grades`, `final_grades`,
  `attendance_records`, atau `enrollments`. Jalur satu-satunya adalah usulan
  yang disetujui manusia melalui alur `domain.md` yang sudah ada.
