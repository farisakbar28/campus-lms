# Domain Model — campus-lms

## 0. Tujuan dan batas domain

`campus-lms` adalah **Learning Management System (LMS) multi-tenant untuk perguruan tinggi**, bukan Sistem Informasi Akademik (SIAKAD).

LMS bertanggung jawab atas proses pembelajaran setelah data akademik dasar tersedia: penyelenggaraan kelas daring/hybrid, RPS dan capaian pembelajaran, materi, aktivitas belajar, tugas, kuis/ujian, diskusi, presensi, progres belajar, penilaian, feedback, gradebook, serta pelaporan aktivitas pembelajaran.

Data akademik inti tetap dimiliki oleh SIAKAD dan disinkronkan ke LMS.

### 0.1 Source of truth

**SIAKAD adalah source of truth** untuk:

* periode/semester akademik;
* program studi sebagai referensi;
* master mata kuliah;
* kelas/perkuliahan yang ditawarkan pada semester tertentu;
* dosen yang ditugaskan mengajar;
* mahasiswa peserta kelas;
* enrollment/KRS mahasiswa.

Data tersebut bersifat **read-only dari sisi LMS**.

LMS tidak boleh menjadi jalur alternatif untuk mengubah data akademik yang seharusnya dimiliki SIAKAD.

**LMS adalah source of truth** untuk:

* struktur konten course;
* RPS/learning plan yang dikelola dalam LMS;
* CPMK dan Sub-CPMK versi LMS;
* module dan lesson/session;
* materials;
* learning activities;
* assignments dan submissions;
* quizzes/exams dan attempts;
* question banks;
* discussion forums;
* attendance;
* activity completion/progress;
* rubrics;
* gradebook;
* feedback;
* final grade sebelum/sesudah dikirim kembali ke SIAKAD;
* announcements;
* audit trail aktivitas LMS.

Final grade dan attendance dapat dikirim kembali ke SIAKAD melalui integration adapter.

### 0.2 Yang berada di luar domain

Fitur berikut **bukan tanggung jawab `campus-lms`**:

* membuat atau mengubah fakultas;
* membuat atau mengubah program studi;
* penerimaan mahasiswa baru;
* registrasi akademik;
* KRS/perwalian;
* pembayaran UKT/SPP;
* penjadwalan ruang;
* manajemen ruang kuliah;
* pengelolaan kurikulum institusional secara penuh;
* transkrip akademik;
* IP/IPK;
* kelulusan dan yudisium;
* ijazah;
* kepegawaian dosen;
* direct message/chat 1:1;
* implementasi aktif SPADA/PDDikti pada scope saat ini.

Data referensi yang dibutuhkan dari domain-domain tersebut dapat diterima melalui integrasi SIAKAD.

### 0.3 Istilah utama

* **Tenant**: satu institusi/perguruan tinggi yang menggunakan SaaS.
* **User**: identitas login global pada platform SaaS.
* **Membership**: hubungan seorang user dengan satu tenant dan role-nya di tenant tersebut.
* **Course**: master mata kuliah, misalnya `IF101 — Algoritma dan Pemrograman`.
* **Course Offering**: pelaksanaan sebuah course pada term/semester dan kelas tertentu, misalnya `IF101 Kelas A — 2026/2027 Ganjil`.
* **Course Staff**: dosen atau Teaching Assistant yang ditugaskan pada suatu course offering.
* **Module**: kelompok/topik besar dalam struktur course.
* **Lesson**: unit/babak/session pembelajaran di dalam module.
* **Learning Activity**: aktivitas yang harus/ dapat dilakukan mahasiswa, misalnya assignment, quiz, forum, atau aktivitas eksternal.
* **Enrollment**: keikutsertaan mahasiswa pada suatu course offering.
* **Assessment**: aktivitas yang menghasilkan atau dapat menghasilkan penilaian.
* **Grade Item**: komponen dalam gradebook.
* **Grade**: skor mahasiswa pada satu grade item.
* **Final Grade**: hasil akhir seluruh grade item pada suatu course offering.
* **Activity Completion**: status penyelesaian aktivitas belajar. Bukan presensi.
* **Attendance**: catatan kehadiran mahasiswa pada suatu sesi presensi.

---

# 1. Aktor dan perannya

Role dibagi menjadi dua tingkat:

1. **platform/tenant role**, melalui `memberships`;
2. **course-scoped role**, melalui `course_staff`.

Seseorang dapat mempunyai satu identitas `user` global dan menjadi anggota lebih dari satu tenant. Seluruh data dan kewenangannya tetap terisolasi berdasarkan tenant.

## 1.1 Super Admin

Super Admin adalah pengelola global SaaS, bukan pengelola akademik sebuah kampus.

### Boleh

* membuat, mengaktifkan, menonaktifkan, atau menangguhkan tenant;
* melihat metadata tenant untuk operasional SaaS;
* mengelola konfigurasi platform global;
* mengelola paket/kapabilitas SaaS jika kelak terdapat subscription;
* melihat health/status integrasi secara global;
* melakukan tindakan support lintas tenant melalui mekanisme **break-glass**;
* mengelola akun administrator tenant dalam keadaan yang diperlukan;
* melihat `platform_audit_logs`;
* menjalankan proses administrasi platform yang tidak menyentuh proses akademik.

### Tidak boleh

* menjadi dosen hanya karena mempunyai role Super Admin;
* memberikan nilai mahasiswa;
* mengubah submission mahasiswa;
* mengubah attendance mahasiswa;
* mengubah enrollment akademik;
* mengubah RPS/course content tanpa otorisasi tenant;
* membaca data tenant secara bebas tanpa kebutuhan operasional yang sah;
* menonaktifkan audit trail;
* melakukan perubahan data tenant tanpa jejak audit.

Setiap akses break-glass ke data tenant wajib memiliki:

* alasan;
* actor;
* tenant target;
* timestamp;
* tindakan yang dilakukan;
* audit trail.

---

## 1.2 Admin Kampus

Admin Kampus adalah administrator LMS pada satu tenant.

Role ini **bukan operator SIAKAD** dan tidak memiliki kewenangan akademik otomatis.

### Boleh

* mengelola konfigurasi LMS tenant;
* mengatur branding tenant;
* mengatur timezone default;
* mengatur kebijakan LMS tenant;
* mengelola membership dan role LMS;
* menentukan konfigurasi storage;
* menentukan kebijakan file upload;
* mengatur integrasi tenant;
* mengatur metode attendance yang diperbolehkan;
* mengatur konfigurasi grade scheme default;
* melihat status sinkronisasi;
* melihat audit log tenant sesuai kebijakan akses;
* membantu lifecycle operasional course;
* melakukan archive/unarchive berdasarkan kewenangan administratif;
* mengelola feature flags tenant.

### Tidak boleh

* membuat fakultas/prodi/mata kuliah akademik yang seharusnya berasal dari SIAKAD;
* mengubah enrollment akademik secara manual untuk mengakali SIAKAD;
* mengubah grade mahasiswa;
* mempublikasikan final grade menggantikan Lead Instructor;
* mengubah submission mahasiswa;
* menyamar sebagai mahasiswa/dosen tanpa mekanisme support yang diaudit;
* menghapus audit log;
* hard-delete histori pembelajaran.

---

## 1.3 Academic Operator

Academic Operator bertanggung jawab atas **boundary antara SIAKAD dan LMS**.

### Boleh

* menjalankan sinkronisasi SIAKAD;
* melihat sync job dan sync error;
* melakukan retry terhadap sinkronisasi gagal;
* melakukan mapping external ID;
* menyelesaikan konflik mapping;
* melihat academic term;
* melihat course;
* melihat course offering;
* melihat instructor assignment;
* melihat enrollment;
* memvalidasi hasil sinkronisasi;
* mengirim final grade yang telah dipublikasikan ke SIAKAD;
* mengirim attendance ke SIAKAD;
* melakukan retry outbound synchronization;
* menjalankan atau memantau copy/import course apabila diberi kewenangan;
* mengarsipkan course offering sesuai proses operasional tenant.

### Tidak boleh

* mengubah fakta akademik dari SIAKAD langsung di LMS;
* menambah mahasiswa ke kelas secara permanen dengan bypass SIAKAD;
* menghapus mahasiswa dari enrollment yang masih aktif di SIAKAD;
* mengubah dosen pengampu yang berasal dari SIAKAD;
* mengubah skor mahasiswa;
* mempublikasikan final grade;
* mengubah submission;
* mengubah jawaban quiz;
* mengubah attendance tanpa workflow koreksi yang sah;
* menggunakan integration mapping untuk memindahkan data antar-tenant.

---

## 1.4 Dosen

Seorang user dengan membership dosen belum otomatis berhak mengelola semua course.

Kewenangan course berasal dari `course_staff`.

Course role untuk dosen:

* `lead_instructor`
* `instructor`

### Lead Instructor boleh

Semua kemampuan Instructor, ditambah:

* mengatur struktur pembelajaran;
* mengelola RPS versi course offering;
* mengelola CPMK/Sub-CPMK;
* mengatur gradebook;
* mengatur grade scheme course;
* mengatur bobot penilaian;
* mempublikasikan nilai;
* melakukan finalisasi final grade;
* melakukan lock final grade;
* membuka kembali final grade dengan alasan yang tercatat;
* menentukan hak Teaching Assistant;
* mempublikasikan course ke mahasiswa;
* menutup proses pembelajaran course.

### Instructor boleh

* membuat dan mengedit module;
* membuat dan mengedit lesson;
* mengunggah materials;
* membuat assignment;
* membuat quiz/exam;
* menggunakan question bank sesuai akses;
* membuat rubrics;
* membuat discussion forum;
* membuat announcement;
* membuat attendance session;
* melihat peserta kelas;
* melihat submissions;
* melakukan grading;
* memberikan feedback;
* mengatur individual override;
* mengatur group activities;
* melihat learning progress;
* melihat laporan course.

### Instructor tidak boleh

* mengubah enrollment SIAKAD;
* mengubah identitas akademik mata kuliah;
* mengubah academic term;
* mengubah instructor assignment yang berasal dari SIAKAD;
* memindahkan mahasiswa antar-course offering;
* menembus isolasi tenant;
* hard-delete grade/submission/audit history;
* publish atau lock final grade jika bukan `lead_instructor`.

---

## 1.5 Teaching Assistant

Teaching Assistant adalah role **course-scoped**, bukan hak administratif seluruh tenant.

Hak TA dapat dibatasi oleh konfigurasi course.

### Boleh

Jika diberikan izin:

* melihat daftar mahasiswa course;
* melihat materials;
* mengunggah/mengedit materials;
* membantu mengelola lesson;
* membantu forum diskusi;
* membuat atau memoderasi discussion;
* membuat attendance session;
* mengelola attendance;
* melihat submission;
* memberikan feedback;
* melakukan grading sebagai **draft score**;
* membantu penilaian rubric;
* melihat learning progress.

### Tidak boleh

* publish final grade;
* lock final grade;
* mengubah grade scheme;
* mengubah bobot gradebook;
* mengubah enrollment;
* mengubah course ownership;
* mengubah Lead Instructor;
* mengubah data akademik dari SIAKAD;
* mengirim final grade ke SIAKAD;
* mengubah nilai final yang sudah dikunci;
* melakukan hard-delete histori akademik.

Nilai yang diberikan TA harus menyimpan:

* `graded_by`;
* status `draft`;
* timestamp;
* rubric detail jika digunakan.

Dosen pengampu tetap bertanggung jawab terhadap publikasi nilai.

---

## 1.6 Mahasiswa

### Boleh

* melihat course offering tempat dirinya enrolled;
* membaca course overview dan RPS yang dipublikasikan;
* melihat module/lesson yang telah tersedia;
* membaca atau mengunduh material yang diperbolehkan;
* melakukan activity completion;
* mengumpulkan assignment;
* mengirim text submission;
* mengunggah file submission;
* melakukan quiz/exam;
* melihat feedback yang sudah dipublikasikan;
* melihat grade yang telah dirilis;
* mengikuti discussion forum;
* membuat thread/reply jika forum mengizinkan;
* melakukan attendance self check-in;
* melakukan attendance melalui QR/PIN;
* melihat attendance miliknya;
* melihat progress pembelajaran miliknya;
* menerima announcement.

### Tidak boleh

* mengakses course yang bukan enrollment-nya;
* melihat submission mahasiswa lain kecuali peer activity secara eksplisit mendukungnya di masa depan;
* melihat grade mahasiswa lain;
* melihat attendance mahasiswa lain;
* mengubah nilai;
* mengubah enrollment;
* mengubah deadline;
* mengubah attendance session;
* mengedit submission setelah submission dikunci;
* melakukan attempt quiz di luar window yang diizinkan;
* menggunakan individual override milik mahasiswa lain;
* mengakses unpublished material/activity;
* mengakses draft grade.

---

# 2. Boundary multi-tenant dan identity

## 2.1 Global identity

`users` adalah identitas global platform.

Satu user dapat:

* menjadi dosen pada Tenant A;
* menjadi mahasiswa pada Tenant B;
* memiliki role berbeda di masing-masing tenant.

Role tidak disimpan secara global di `users`.

Hubungan user dengan tenant wajib melalui:

`users -> memberships -> tenants`

## 2.2 Tenant isolation

Semua data pembelajaran wajib mempunyai hubungan deterministik ke satu `tenant_id`.

Aturan mutlak:

> Query terhadap data tenant tidak boleh bergantung hanya pada ID object. Tenant context wajib ikut divalidasi.

Contoh:

```text
SALAH:
SELECT * FROM submissions WHERE id = :submission_id

BENAR secara domain:
tenant -> course_offering -> enrollment -> submission
```

RLS harus memastikan object Tenant A tidak dapat dibaca, direferensikan, atau dimodifikasi dari Tenant B.

## 2.3 Global tables

Entitas global hanya digunakan jika memang tidak dimiliki tenant.

Global:

* `tenants`
* `users`
* `auth_identities`
* `platform_admins`
* `platform_audit_logs`

Semua domain pembelajaran lainnya tenant-scoped.

---

# 3. Entitas inti dan relasinya

## 3.1 Tenant dan identity

| Entity                | Scope  | Pemilik       | Atribut penting                                                                                  |
| --------------------- | ------ | ------------- | ------------------------------------------------------------------------------------------------ |
| `tenants`             | Global | Platform      | `id`, `slug`, `name`, `status`, `default_timezone`, `created_at`, `suspended_at`                 |
| `tenant_settings`     | Tenant | Admin Kampus  | `tenant_id`, branding, locale, upload policy, attendance policy, grading defaults, feature flags |
| `users`               | Global | User/Platform | `id`, `email`, `display_name`, `status`, `created_at`                                            |
| `auth_identities`     | Global | User/Platform | provider, provider subject, `user_id`, verified state                                            |
| `memberships`         | Tenant | Tenant        | `tenant_id`, `user_id`, `role`, `status`, `joined_at`                                            |
| `external_identities` | Tenant | Integration   | `tenant_id`, `user_id`, `source`, `external_user_id`, `external_type`                            |

Constraint:

```text
UNIQUE (tenant_id, user_id)
```

pada `memberships`.

Role membership:

```text
tenant_admin
academic_operator
lecturer
student
```

`teaching_assistant` tidak harus menjadi tenant-wide role; status TA ditentukan melalui `course_staff`.

---

# 4. Academic reference dari SIAKAD

Entitas di bagian ini tenant-scoped tetapi **externally authoritative**.

## 4.1 `academic_terms`

Merepresentasikan semester/periode akademik.

Atribut:

* `id`
* `tenant_id`
* `external_id`
* `code`
* `name`
* `starts_at`
* `ends_at`
* `status`
* `synced_at`

Owner: **SIAKAD**

LMS tidak boleh mengubah fakta akademiknya.

---

## 4.2 `academic_program_refs`

Referensi program studi untuk keperluan filtering/reporting.

Atribut:

* `id`
* `tenant_id`
* `external_id`
* `code`
* `name`
* `status`

Owner: **SIAKAD**

LMS tidak mengelola struktur fakultas/prodi.

---

## 4.3 `courses`

Master mata kuliah.

Contoh:

```text
IF101
Algoritma dan Pemrograman
3 SKS
```

Atribut:

* `id`
* `tenant_id`
* `external_id`
* `code`
* `name`
* `credits`
* `academic_program_ref_id`
* `status`
* `synced_at`

Owner: **SIAKAD**

`courses` bukan kelas semester.

---

## 4.4 `course_offerings`

Instance sebuah `course` pada term tertentu.

Contoh:

```text
IF101
Kelas A
Semester Ganjil 2026/2027
```

Atribut:

* `id`
* `tenant_id`
* `external_id`
* `course_id`
* `academic_term_id`
* `external_section_code`
* `display_name`
* `lms_status`
* `published_at`
* `closed_at`
* `archived_at`
* `course_plan_version_id`
* `created_at`

Owner terbagi:

**SIAKAD memiliki:**

* identitas course;
* term;
* kode section/class.

**LMS memiliki:**

* `lms_status`;
* publish state;
* content;
* learning plan;
* archive state.

Lifecycle:

```text
draft
  ↓
published
  ↓
active
  ↓
closed
  ↓
archived
```

`archived` bersifat read-only secara default.

---

## 4.5 `course_staff`

Menghubungkan dosen/TA dengan `course_offering`.

Atribut:

* `id`
* `tenant_id`
* `course_offering_id`
* `user_id`
* `role`
* `source`
* `permissions`
* `active`

Role:

```text
lead_instructor
instructor
teaching_assistant
```

Untuk `lead_instructor` dan `instructor`, assignment utama berasal dari SIAKAD.

TA dapat berasal dari integrasi atau assignment lokal sesuai kebijakan tenant.

Constraint:

```text
UNIQUE (course_offering_id, user_id)
```

---

## 4.6 `enrollments`

Hubungan mahasiswa dengan `course_offering`.

Atribut:

* `id`
* `tenant_id`
* `course_offering_id`
* `student_user_id`
* `external_id`
* `status`
* `enrolled_at`
* `withdrawn_at`
* `synced_at`

Owner: **SIAKAD**

Constraint wajib:

```text
UNIQUE (course_offering_id, student_user_id)
```

Mahasiswa tidak boleh terdaftar dua kali pada offering yang sama.

Enrollment yang hilang dari sync tidak langsung di-hard-delete jika sudah mempunyai aktivitas pembelajaran.

Gunakan status seperti:

```text
active
withdrawn
completed
inactive
```

---

# 5. RPS dan capaian pembelajaran

## 5.1 Prinsip

CPL berasal dari kurikulum/SIAKAD dan bersifat reference/read-only.

CPMK, Sub-CPMK, serta learning plan dapat dikelola LMS tetapi harus **versioned**.

Course offering lama tidak boleh berubah ketika RPS untuk semester berikutnya diperbarui.

---

## 5.2 `learning_outcomes`

Atribut:

* `id`
* `tenant_id`
* `course_id`
* `type`
* `code`
* `title`
* `description`
* `external_id`
* `source`
* `version`
* `active`

`type`:

```text
CPL
CPMK
SUB_CPMK
```

Rules:

* `CPL` → external/read-only;
* `CPMK` → LMS-managed dan versioned;
* `SUB_CPMK` → LMS-managed dan versioned.

---

## 5.3 `outcome_mappings`

Merepresentasikan relasi:

```text
CPL -> CPMK
CPMK -> Sub-CPMK
```

Atribut:

* `id`
* `tenant_id`
* `parent_outcome_id`
* `child_outcome_id`
* `weight` nullable

---

## 5.4 `course_plan_versions`

Versi RPS/learning plan.

Atribut:

* `id`
* `tenant_id`
* `course_id`
* `version_number`
* `title`
* `description`
* `learning_guidance`
* `status`
* `created_by`
* `published_by`
* `published_at`
* `created_at`

Status:

```text
draft
published
retired
```

Setelah sebuah version digunakan oleh course offering aktif, historinya tidak boleh ditimpa.

Perubahan besar menghasilkan version baru.

---

## 5.5 `course_plan_outcomes`

Menghubungkan `course_plan_versions` dengan outcome version yang digunakan.

Tujuannya agar course offering lama mempertahankan snapshot capaian pembelajaran yang benar.

---

# 6. Struktur pembelajaran

Struktur default:

```text
Course
└── Course Offering
    ├── Course Plan / RPS
    ├── Module
    │   └── Lesson
    │       ├── Material
    │       └── Learning Activity
    └── Gradebook
```

## 6.1 `modules`

Atribut:

* `id`
* `tenant_id`
* `course_offering_id`
* `title`
* `description`
* `position`
* `status`
* `available_from`
* `available_until`
* `created_by`

Status:

```text
draft
published
hidden
```

Module boleh merepresentasikan:

* minggu;
* topik;
* bab;
* unit;
* blok pembelajaran.

Domain tidak memaksakan bahwa satu module = satu minggu.

---

## 6.2 `lessons`

Unit/babak/session pembelajaran di dalam module.

Atribut:

* `id`
* `tenant_id`
* `module_id`
* `title`
* `description`
* `position`
* `learning_mode`
* `estimated_minutes`
* `available_from`
* `available_until`
* `status`

`learning_mode` dapat berupa:

```text
asynchronous
synchronous
blended
onsite
```

---

## 6.3 `lesson_outcomes`

Mapping:

```text
Lesson -> CPMK/Sub-CPMK
```

Digunakan untuk learning analytics dan OBE reporting.

---

# 7. Materials dan file

## 7.1 `files`

Metadata file tenant.

Atribut:

* `id`
* `tenant_id`
* `storage_key`
* `original_filename`
* `mime_type`
* `size_bytes`
* `checksum`
* `uploaded_by`
* `malware_scan_status`
* `created_at`

Binary object tidak disimpan langsung pada row domain.

Semua file harus mempunyai tenant ownership.

---

## 7.2 `materials`

Learning resource.

Atribut:

* `id`
* `tenant_id`
* `lesson_id`
* `title`
* `description`
* `type`
* `file_id`
* `external_url`
* `content`
* `position`
* `published`
* `created_by`

Material type minimal:

```text
text
file
link
video
audio
embed
learning_package
external_tool
```

File tidak dibatasi hanya PDF.

PDF dapat didukung sebagai salah satu format submission/material, bukan satu-satunya format.

---

# 8. Learning activities

## 8.1 `learning_activities`

Base entity untuk aktivitas pembelajaran.

Atribut:

* `id`
* `tenant_id`
* `lesson_id`
* `type`
* `title`
* `description`
* `position`
* `available_from`
* `due_at`
* `cutoff_at`
* `published`
* `completion_required`
* `created_by`

Type minimal:

```text
assignment
quiz
discussion
attendance
external_tool
```

Entity spesifik menyimpan konfigurasi masing-masing aktivitas.

---

## 8.2 `activity_outcomes`

Mapping antara learning activity dan CPMK/Sub-CPMK.

Atribut:

* `activity_id`
* `learning_outcome_id`
* `weight` nullable

Assessment dapat mengukur lebih dari satu CPMK.

---

# 9. Assignment dan submission

## 9.1 `assignments`

Atribut:

* `id`
* `tenant_id`
* `learning_activity_id`
* `instructions`
* `submission_mode`
* `submission_types`
* `max_attempts`
* `max_score`
* `group_mode`
* `allow_resubmission`
* `require_submission_statement`
* `created_by`

`submission_types` dapat mengizinkan:

```text
text
file
text_and_file
```

`group_mode`:

```text
individual
group
```

---

## 9.2 `assignment_overrides`

Pengecualian untuk mahasiswa atau group tertentu.

Atribut:

* `id`
* `tenant_id`
* `assignment_id`
* `enrollment_id` nullable
* `group_id` nullable
* `available_from`
* `due_at`
* `cutoff_at`
* `max_attempts`
* `reason`
* `created_by`

Exactly one target:

```text
enrollment_id XOR group_id
```

Override tidak mengubah deadline mahasiswa lain.

---

## 9.3 `submissions`

Logical submission milik mahasiswa/group.

Atribut:

* `id`
* `tenant_id`
* `assignment_id`
* `enrollment_id` nullable
* `group_id` nullable
* `status`
* `current_version`
* `first_submitted_at`
* `last_submitted_at`
* `locked_at`

Status:

```text
draft
submitted
returned
resubmission_allowed
locked
```

---

## 9.4 `submission_versions`

Setiap submit/resubmit menghasilkan version immutable.

Atribut:

* `id`
* `tenant_id`
* `submission_id`
* `attempt_number`
* `text_content`
* `submitted_at`
* `submitted_by`

Submission lama tidak ditimpa.

---

## 9.5 `submission_files`

Mapping submission version ke `files`.

---

# 10. Group learning

## 10.1 `course_groups`

Atribut:

* `id`
* `tenant_id`
* `course_offering_id`
* `name`
* `description`
* `created_by`

## 10.2 `course_group_members`

Atribut:

* `group_id`
* `enrollment_id`
* `joined_at`

Constraint:

```text
UNIQUE (group_id, enrollment_id)
```

Kelompok dapat digunakan untuk:

* group assignment;
* group discussion;
* aktivitas kolaboratif.

---

# 11. Quiz dan exam

Quiz dan exam menggunakan engine domain yang sama, dengan konfigurasi berbeda.

## 11.1 `question_banks`

Atribut:

* `id`
* `tenant_id`
* `course_id` nullable
* `course_offering_id` nullable
* `name`
* `description`
* `visibility`
* `created_by`

Bank dapat bersifat:

* khusus offering;
* reusable pada course yang sama.

Sharing lintas tenant dilarang kecuali ada fitur library eksplisit di masa depan.

---

## 11.2 `questions`

Logical question identity.

Atribut:

* `id`
* `tenant_id`
* `question_bank_id`
* `type`
* `created_by`
* `status`

Question type minimal:

```text
multiple_choice_single
multiple_choice_multiple
true_false
short_answer
essay
numeric
```

Domain harus extensible untuk tipe soal tambahan.

---

## 11.3 `question_versions`

Isi soal immutable/versioned.

Atribut:

* `id`
* `tenant_id`
* `question_id`
* `version_number`
* `prompt`
* `explanation`
* `default_points`
* `configuration`
* `created_at`

Jawaban quiz lama harus tetap merujuk ke question version yang digunakan saat attempt dimulai.

---

## 11.4 `question_options`

Untuk question type yang mempunyai pilihan.

Atribut:

* `id`
* `question_version_id`
* `content`
* `position`
* `is_correct`
* `score_fraction`

Correct-answer metadata tidak boleh pernah dikirim ke mahasiswa sebelum waktunya.

---

## 11.5 `quizzes`

Atribut:

* `id`
* `tenant_id`
* `learning_activity_id`
* `mode`
* `time_limit_seconds`
* `attempt_limit`
* `shuffle_questions`
* `shuffle_answers`
* `grade_method`
* `max_score`
* `feedback_release_policy`
* `created_by`

`mode`:

```text
quiz
exam
```

---

## 11.6 `quiz_question_rules`

Mengatur komposisi soal.

Mendukung:

* fixed question;
* random question;
* random from category/pool;
* points;
* order.

Randomization harus resolved ketika attempt dimulai dan hasilnya disimpan.

---

## 11.7 `quiz_overrides`

Override mahasiswa/group untuk:

* open time;
* close time;
* attempt limit;
* time limit.

---

## 11.8 `quiz_attempts`

Atribut:

* `id`
* `tenant_id`
* `quiz_id`
* `enrollment_id`
* `attempt_number`
* `started_at`
* `expires_at`
* `submitted_at`
* `status`
* `score`

Status:

```text
in_progress
submitted
auto_submitted
graded
invalidated
```

---

## 11.9 `quiz_attempt_questions`

Snapshot daftar question version yang diberikan kepada suatu attempt.

Hal ini penting untuk randomization dan histori ujian.

---

## 11.10 `quiz_responses`

Atribut:

* `id`
* `quiz_attempt_question_id`
* `response`
* `answered_at`
* `auto_score`
* `manual_score`
* `feedback`
* `graded_by`

Essay/manual response dapat menunggu grading dosen/TA.

---

# 12. Rubrics

## 12.1 `rubrics`

Atribut:

* `id`
* `tenant_id`
* `course_offering_id`
* `name`
* `description`
* `status`
* `created_by`

## 12.2 `rubric_criteria`

Atribut:

* `id`
* `rubric_id`
* `title`
* `description`
* `weight`
* `position`

## 12.3 `rubric_levels`

Atribut:

* `id`
* `rubric_criterion_id`
* `label`
* `description`
* `score`

Rubric yang sudah digunakan untuk grading harus mempertahankan version/snapshot historis.

---

# 13. Gradebook

Penilaian **tidak dikunci ke range 1–100**.

UI boleh menampilkan 1–100 sebagai default, tetapi domain menggunakan:

* raw score;
* maximum score;
* weight;
* normalized score;
* grade scheme.

---

## 13.1 `grade_categories`

Contoh:

```text
Tugas       30%
Kuis        20%
UTS         20%
UAS         30%
```

Atribut:

* `id`
* `tenant_id`
* `course_offering_id`
* `name`
* `weight`
* `position`

---

## 13.2 `grade_items`

Satu komponen gradebook.

Atribut:

* `id`
* `tenant_id`
* `course_offering_id`
* `grade_category_id`
* `source_type`
* `source_id`
* `name`
* `max_score`
* `weight`
* `grading_type`
* `rubric_id`
* `published`

`source_type` dapat menunjuk:

```text
assignment
quiz
manual
```

---

## 13.3 `grades`

Nilai mahasiswa untuk satu grade item.

Atribut:

* `id`
* `tenant_id`
* `grade_item_id`
* `enrollment_id`
* `raw_score`
* `normalized_score`
* `status`
* `feedback`
* `graded_by`
* `graded_at`
* `published_by`
* `published_at`

Status:

```text
draft
published
overridden
```

Constraint:

```text
UNIQUE (grade_item_id, enrollment_id)
```

TA hanya dapat menghasilkan/mengubah `draft`.

---

## 13.4 `grade_schemes`

Konversi score ke label/grade kampus.

Atribut:

* `id`
* `tenant_id`
* `name`
* `type`
* `configuration`
* `active`

Contoh output:

```text
A
AB
B
BC
C
D
E
```

Domain tidak mengasumsikan skema tersebut berlaku pada semua kampus.

---

## 13.5 `final_grades`

Atribut:

* `id`
* `tenant_id`
* `course_offering_id`
* `enrollment_id`
* `numeric_result`
* `grade_label`
* `status`
* `calculated_at`
* `published_by`
* `published_at`
* `locked_by`
* `locked_at`
* `sync_status`
* `synced_at`

Status:

```text
draft
published
locked
```

Constraint:

```text
UNIQUE (course_offering_id, enrollment_id)
```

Hanya `lead_instructor` yang boleh melakukan:

```text
draft -> published
published -> locked
```

---

# 14. Grade change dan audit

Perubahan nilai wajib mempunyai jejak.

Setiap perubahan `grades` atau `final_grades` harus menghasilkan `audit_logs`.

Minimal audit payload:

* actor;
* role;
* tenant;
* course offering;
* student;
* grade item;
* previous value;
* new value;
* reason jika override/reopen;
* timestamp;
* request/correlation ID.

Final grade yang sudah `locked` tidak dapat diedit langsung.

Flow:

```text
locked
  ↓
Lead Instructor requests reopen + reason
  ↓
audit
  ↓
published/draft
  ↓
correction
  ↓
re-publish
  ↓
lock
  ↓
outbound sync
```

Academic Operator tidak boleh mengubah score sebagai jalan pintas ketika sync gagal.

---

# 15. Attendance

Attendance dan activity completion adalah dua domain berbeda.

## 15.1 `attendance_sessions`

Atribut:

* `id`
* `tenant_id`
* `course_offering_id`
* `lesson_id` nullable
* `learning_activity_id` nullable
* `title`
* `method`
* `opens_at`
* `closes_at`
* `late_after`
* `created_by`
* `status`

Method:

```text
manual
self_check_in
pin
qr
activity_completion
external
```

---

## 15.2 QR/PIN

QR/PIN adalah **credential untuk check-in**, bukan attendance record.

Credential harus:

* mempunyai masa berlaku;
* terkait hanya dengan satu attendance session;
* tidak menyimpan PIN plaintext jika dapat dihindari;
* dapat dirotasi/revoke;
* tidak reusable setelah session selesai;
* tidak dapat digunakan lintas tenant/course.

QR sebaiknya menggunakan token yang cukup acak dan berumur pendek.

---

## 15.3 `attendance_credentials`

Atribut:

* `id`
* `tenant_id`
* `attendance_session_id`
* `type`
* `credential_hash`
* `valid_from`
* `valid_until`
* `revoked_at`

---

## 15.4 `attendance_records`

Atribut:

* `id`
* `tenant_id`
* `attendance_session_id`
* `enrollment_id`
* `status`
* `check_in_method`
* `checked_in_at`
* `recorded_by`
* `override_reason`
* `updated_at`

Status:

```text
present
late
absent
excused
```

Constraint:

```text
UNIQUE (attendance_session_id, enrollment_id)
```

Perubahan manual setelah check-in wajib diaudit.

---

# 16. Activity completion dan progress

## 16.1 `activity_completions`

Atribut:

* `id`
* `tenant_id`
* `learning_activity_id`
* `enrollment_id`
* `status`
* `progress`
* `completed_at`
* `completion_source`

Status:

```text
not_started
in_progress
completed
```

Constraint:

```text
UNIQUE (learning_activity_id, enrollment_id)
```

Activity completion dapat menjadi salah satu input attendance bila tenant/course mengaktifkannya, tetapi kedua record tetap terpisah.

---

# 17. Discussion forum

Tidak ada direct chat/DM 1:1 di core domain.

## 17.1 `discussion_forums`

Atribut:

* `id`
* `tenant_id`
* `learning_activity_id`
* `mode`
* `group_mode`
* `students_can_create_threads`
* `available_from`
* `available_until`

---

## 17.2 `discussion_threads`

Atribut:

* `id`
* `tenant_id`
* `forum_id`
* `group_id` nullable
* `title`
* `created_by`
* `created_at`
* `locked_at`
* `pinned_at`

---

## 17.3 `discussion_posts`

Atribut:

* `id`
* `tenant_id`
* `thread_id`
* `parent_post_id` nullable
* `author_user_id`
* `content`
* `created_at`
* `edited_at`
* `deleted_at`

Soft-delete digunakan untuk moderation.

Isi asli yang relevan dengan audit/moderation tidak langsung dimusnahkan.

---

# 18. Announcements

## 18.1 `announcements`

Atribut:

* `id`
* `tenant_id`
* `course_offering_id`
* `title`
* `content`
* `published_by`
* `published_at`
* `expires_at`

Audience default adalah seluruh active enrollment dan course staff.

---

## 18.2 `announcement_receipts`

Opsional tetapi disediakan dalam domain untuk read tracking.

Atribut:

* `announcement_id`
* `user_id`
* `read_at`

---

# 19. Copy/import course antarsemester

Course lama dapat digunakan sebagai sumber struktur untuk course offering baru.

## 19.1 Boleh disalin

* RPS sebagai draft version baru;
* CPMK/Sub-CPMK lokal;
* module;
* lesson;
* materials;
* learning activities;
* assignment configuration;
* quiz configuration;
* question bank sesuai pilihan;
* rubric;
* discussion forum configuration;
* gradebook structure.

## 19.2 Tidak boleh disalin

* enrollment;
* attendance records;
* attendance credentials;
* submissions;
* quiz attempts;
* quiz responses;
* student grades;
* final grades;
* activity completion;
* discussion post mahasiswa;
* audit logs.

---

## 19.3 `course_copy_jobs`

Atribut:

* `id`
* `tenant_id`
* `source_course_offering_id`
* `target_course_offering_id`
* `options`
* `status`
* `requested_by`
* `started_at`
* `completed_at`
* `error`

Copy harus menghasilkan entity ID baru.

Target tidak boleh menunjuk record mutable milik semester lama.

---

# 20. Integration boundary

Integrasi dibuat sebagai adapter di pinggir domain.

Core LMS tidak boleh mengetahui detail vendor/protokol SIAKAD tertentu.

```text
SIAKAD
   │
   ▼
Integration Adapter
   │
   ▼
Canonical LMS Domain
```

---

## 20.1 `integrations`

Atribut:

* `id`
* `tenant_id`
* `type`
* `provider`
* `status`
* `configuration_reference`
* `last_success_at`

Type:

```text
siakad
spada
pddikti
scorm
xapi
lti
other
```

Secrets tidak boleh disimpan dalam configuration JSON plaintext.

Gunakan secret manager/reference.

---

## 20.2 `external_mappings`

Mapping canonical ID dengan external ID.

Atribut:

* `id`
* `tenant_id`
* `integration_id`
* `entity_type`
* `internal_id`
* `external_id`
* `external_version`
* `last_synced_at`

Constraint minimal:

```text
UNIQUE (integration_id, entity_type, external_id)
```

---

## 20.3 `sync_jobs`

Atribut:

* `id`
* `tenant_id`
* `integration_id`
* `direction`
* `entity_type`
* `status`
* `started_at`
* `completed_at`
* `cursor`
* `summary`

Direction:

```text
inbound
outbound
```

---

## 20.4 `sync_errors`

Atribut:

* `id`
* `tenant_id`
* `sync_job_id`
* `external_id`
* `entity_type`
* `error_code`
* `error_message`
* `retryable`
* `resolved_at`
* `resolved_by`

---

## 20.5 Inbound SIAKAD

Minimal menerima:

```text
academic_terms
academic_program_refs
users/external identities
courses
course_offerings
course_staff
enrollments
CPL references
```

Sync harus idempotent.

Mengirim payload yang sama dua kali tidak boleh membuat duplicate entity.

---

## 20.6 Outbound SIAKAD

Minimal boundary untuk:

```text
final_grades
attendance
```

Hanya data yang sudah memenuhi kondisi publikasi/finalisasi boleh dikirim.

---

# 21. SPADA dan PDDikti

SPADA/PDDikti adalah **future integration**, bukan dependency core domain saat ini.

Rules:

* jangan menambahkan field SPADA ke setiap tabel core;
* jangan membuat schema LMS mengikuti payload API SPADA;
* gunakan `integrations` + `external_mappings` + adapter;
* perubahan API eksternal tidak boleh memaksa perubahan model pembelajaran internal.

SPADA/PDDikti adapter dapat dikembangkan kemudian tanpa mengubah konsep course, offering, enrollment, assessment, atau gradebook.

---

# 22. SCORM, xAPI, dan LTI

Ketiganya disiapkan sebagai integration boundary.

Implementasi protokol penuh tidak wajib pada MVP pertama.

## 22.1 `learning_packages`

Untuk package seperti SCORM.

Atribut:

* `id`
* `tenant_id`
* `title`
* `standard`
* `version`
* `file_id`
* `created_by`
* `status`

---

## 22.2 `external_tools`

Untuk integrasi tool seperti LTI.

Atribut:

* `id`
* `tenant_id`
* `name`
* `protocol`
* `configuration_reference`
* `status`

Credentials/secrets tidak disimpan secara plaintext.

---

## 22.3 `learning_events`

Boundary untuk event pembelajaran/xAPI-like tracking.

Atribut konseptual:

* `tenant_id`
* `user_id`
* `course_offering_id`
* `activity_id`
* `event_type`
* `occurred_at`
* `payload`

Retention dan volume event perlu dikelola terpisah dari transactional domain jika skala sudah besar.

---

# 23. Audit logs

## 23.1 `audit_logs`

Tenant-scoped immutable log.

Atribut:

* `id`
* `tenant_id`
* `actor_user_id`
* `actor_role`
* `action`
* `entity_type`
* `entity_id`
* `course_offering_id` nullable
* `before_data`
* `after_data`
* `reason`
* `ip_address`
* `request_id`
* `occurred_at`

Audit wajib untuk tindakan sensitif, termasuk:

* perubahan role;
* perubahan course staff;
* publish/unpublish course;
* perubahan deadline;
* individual override;
* grading;
* perubahan grade;
* publish grade;
* lock/unlock final grade;
* attendance override;
* archive/unarchive course;
* sync conflict resolution;
* break-glass access;
* perubahan konfigurasi tenant/integration.

Audit log tidak boleh diedit oleh user aplikasi biasa.

---

## 23.2 `platform_audit_logs`

Global.

Untuk:

* pembuatan/suspension tenant;
* perubahan Super Admin;
* break-glass;
* konfigurasi global;
* tindakan support lintas tenant.

---

# 24. Aturan bisnis yang tidak boleh dilanggar

## 24.1 Multi-tenancy

1. Data Tenant A tidak boleh terbaca oleh Tenant B.
2. Semua foreign key antar-entitas tenant harus menunjuk entity pada tenant yang sama.
3. Course Tenant A tidak boleh mempunyai material Tenant B.
4. Enrollment Tenant A tidak boleh menunjuk user membership yang tidak aktif pada Tenant A.
5. Integration mapping tidak boleh memetakan object lintas tenant.
6. Background job wajib membawa `tenant_id`.
7. File access wajib memvalidasi tenant.
8. Cache key untuk tenant data wajib membawa tenant scope.

---

## 24.2 Identity

1. `users` tidak menyimpan tenant role.
2. User harus mempunyai active `membership` untuk mengakses tenant.
3. Satu user boleh mempunyai banyak membership.
4. Hak course tidak boleh disimpulkan hanya dari membership `lecturer`.
5. Dosen harus tercatat di `course_staff`.
6. Mahasiswa harus mempunyai active `enrollment`.

---

## 24.3 SIAKAD

1. Data academic master tidak diedit manual di LMS.
2. Sync harus idempotent.
3. External ID harus unik dalam integration context.
4. Data tidak boleh hilang hanya karena satu sync gagal.
5. Sync failure harus menghasilkan `sync_error`.
6. Academic Operator tidak boleh mengubah data nilai untuk menyelesaikan error integrasi.
7. Conflict resolution harus diaudit.
8. Historical learning data tidak boleh ikut terhapus ketika enrollment berubah di SIAKAD.

---

## 24.4 Course offering

1. `course` dan `course_offering` adalah entity berbeda.
2. Course offering harus terkait satu academic term.
3. Content mahasiswa hanya terlihat setelah course/activity dipublikasikan.
4. Course offering `archived` read-only secara default.
5. Reopen course archived harus diaudit.
6. Data historis semester lama tidak boleh berubah akibat edit semester baru.

---

## 24.5 Enrollment

1. Mahasiswa tidak boleh mempunyai dua enrollment pada offering yang sama.
2. Mahasiswa hanya dapat mengakses course tempat dirinya enrolled.
3. Enrollment dari SIAKAD tidak boleh dibuat/diubah manual sebagai operasi normal.
4. Withdrawal tidak menghapus submission/grade/history.
5. Hak akses student harus mengikuti enrollment status dan kebijakan tenant.

---

## 24.6 Deadline

Semua waktu disimpan sebagai timestamp absolut dan ditampilkan menggunakan timezone tenant/course.

Jangan menggunakan server-local time sebagai makna deadline.

Assignment mempunyai:

```text
available_from
due_at
cutoff_at
```

Interpretasi:

* `due_at` = deadline akademik;
* `cutoff_at` = hard stop penerimaan submission.

Default tenant dapat mengatur:

```text
cutoff_at = due_at
```

untuk kebijakan strict deadline.

Jika kampus mengizinkan late submission:

```text
cutoff_at > due_at
```

Submission setelah `due_at` tetapi sebelum `cutoff_at` diberi status late.

Individual extension menggunakan `assignment_overrides`.

---

## 24.7 Submission

1. Submission harus dimiliki enrollment/group yang valid.
2. Submit di luar allowed window ditolak kecuali ada override.
3. Setiap resubmission membuat version baru.
4. Version lama tidak ditimpa.
5. Dosen/TA tidak boleh mengganti file mahasiswa secara diam-diam.
6. Setelah submission locked, mahasiswa tidak boleh mengedit.
7. Group submission hanya dapat dibuat oleh member group yang sah.

---

## 24.8 Quiz/exam

1. Attempt tidak boleh dimulai di luar window tanpa override.
2. Attempt limit harus ditegakkan.
3. Timer dihitung server-side.
4. Question version untuk attempt harus frozen.
5. Random question selection harus disimpan saat attempt dimulai.
6. Correct answer tidak boleh dikirim ke client sebelum diizinkan feedback policy.
7. Submitted attempt tidak dapat dimodifikasi mahasiswa.
8. Manual regrade harus diaudit.
9. Invalidation attempt harus mempunyai alasan.

---

## 24.9 Grade

1. TA hanya boleh membuat/mengubah draft score.
2. Instructor boleh grading tetapi tidak publish final grade jika bukan Lead Instructor.
3. Hanya Lead Instructor yang boleh publish/lock final grade.
4. Grade tidak diasumsikan 1–100.
5. Score tidak boleh melebihi rule grade item tanpa explicit override.
6. Setiap perubahan nilai harus menyimpan actor dan timestamp.
7. Perubahan final grade wajib diaudit.
8. Final grade locked tidak dapat diedit langsung.
9. Outbound sync hanya mengirim final grade yang telah dipublikasikan sesuai kebijakan tenant.
10. Sync error tidak mengubah nilai akademik internal.

---

## 24.10 Attendance

1. Satu mahasiswa hanya mempunyai satu attendance record per attendance session.
2. QR/PIN hanya valid untuk session dan window yang benar.
3. Attendance credential tidak boleh reusable lintas session.
4. Attendance correction wajib menyimpan actor dan reason.
5. Activity completion tidak otomatis sama dengan hadir kecuali rule course menyatakannya.
6. Attendance tidak boleh hilang ketika course diarsipkan.
7. Outbound attendance sync harus idempotent.

---

## 24.11 Discussion

1. User hanya dapat melihat forum course yang dapat diaksesnya.
2. Group forum hanya terlihat oleh group yang tepat dan staff berwenang.
3. Edit/delete post harus mengikuti moderation policy.
4. Soft-deleted content tetap dapat disimpan untuk audit/moderation sesuai retention policy.
5. Forum tidak memberi akses ke private DM.

---

## 24.12 Copy course

1. Copy hanya mengambil instructional structure.
2. Student data tidak boleh ikut.
3. ID baru harus dibuat pada target.
4. Grade/submission/attendance tidak boleh ikut.
5. Copy operation harus diaudit.
6. Target course tetap mempunyai tenant yang sama.
7. Copy lintas tenant dilarang pada core domain.

---

## 24.13 Deletion dan retention

Entity historis berikut tidak boleh hard-delete melalui aplikasi normal:

* submissions;
* submission versions;
* quiz attempts;
* quiz responses;
* grades;
* final grades;
* attendance records;
* audit logs;
* sync history yang relevan;
* published course-plan history.

Gunakan:

* archive;
* soft delete;
* retention/anonymization workflow khusus jika diwajibkan kebijakan privasi.

---

# 25. RLS classification

## Global — tidak memakai tenant RLS biasa

```text
tenants
users
auth_identities
platform_admins
platform_audit_logs
```

Akses tetap dibatasi berdasarkan role/platform policy.

## Tenant-scoped — wajib RLS

```text
tenant_settings
memberships
external_identities

academic_terms
academic_program_refs
courses
course_offerings
course_staff
enrollments

learning_outcomes
outcome_mappings
course_plan_versions
course_plan_outcomes

modules
lessons
lesson_outcomes
materials
files

learning_activities
activity_outcomes

assignments
assignment_overrides
submissions
submission_versions
submission_files

course_groups
course_group_members

question_banks
questions
question_versions
question_options
quizzes
quiz_question_rules
quiz_overrides
quiz_attempts
quiz_attempt_questions
quiz_responses

rubrics
rubric_criteria
rubric_levels

grade_categories
grade_items
grades
grade_schemes
final_grades

attendance_sessions
attendance_credentials
attendance_records

activity_completions

discussion_forums
discussion_threads
discussion_posts

announcements
announcement_receipts

course_copy_jobs

integrations
external_mappings
sync_jobs
sync_errors
learning_packages
external_tools
learning_events

audit_logs
```

Prinsip:

> Jika entity merupakan data tenant tetapi tidak memiliki `tenant_id` langsung karena child relation, tenant tetap harus dapat dibuktikan secara deterministik melalui parent dan tidak boleh menghasilkan jalur cross-tenant.

Untuk implementasi RLS yang lebih sederhana dan aman, entity sensitif sebaiknya tetap menyimpan `tenant_id` secara eksplisit meskipun tenant dapat diturunkan melalui parent.

---

# 26. Authorization hierarchy

Urutan validasi akses:

```text
Authenticated User
      │
      ▼
Active Membership?
      │
      ▼
Correct Tenant?
      │
      ▼
Tenant Role Permission?
      │
      ▼
Course Enrollment / Course Staff?
      │
      ▼
Object-level Permission?
      │
      ▼
Lifecycle/State Allows Action?
```

Contoh grading:

```text
User authenticated
   ↓
membership tenant active
   ↓
course_staff exists
   ↓
role = instructor / lead_instructor / permitted TA
   ↓
grade item belongs to same course offering
   ↓
grade not final-locked
   ↓
action permitted
```

Role saja tidak cukup.

State object juga harus divalidasi.

---

# 27. Lifecycle penting

## Course

```text
draft
  ↓
published
  ↓
active
  ↓
closed
  ↓
archived
```

## Assignment submission

```text
draft
  ↓
submitted
  ↓
returned ─────────┐
  ↓               │
resubmission_allowed
  ↓
submitted
  ↓
locked
```

## Quiz attempt

```text
in_progress
  ↓
submitted / auto_submitted
  ↓
graded
```

Exceptional:

```text
invalidated
```

## Grade

```text
draft
  ↓
published
```

## Final grade

```text
draft
  ↓
published
  ↓
locked
```

Correction harus melalui audited reopen flow.

---

# 28. ERD sederhana

```text
                              GLOBAL PLATFORM
┌──────────────────┐
│      users       │
└────────┬─────────┘
         │
         │ 1..*
         ▼
┌──────────────────┐        ┌──────────────────┐
│   memberships    │ *────1 │     tenants      │
└──────────────────┘        └────────┬─────────┘
                                    │
                     TENANT BOUNDARY│
                                    ▼

                    ┌──────────────────────┐
                    │    academic_terms    │
                    └──────────┬───────────┘
                               │
                               │
┌───────────────┐       ┌──────▼──────────────┐
│    courses    │ 1───* │  course_offerings  │
└───────┬───────┘       └──────┬──────┬──────┘
        │                       │      │
        │                       │      ├──────────────┐
        │                       │      │              │
        │                 ┌─────▼────┐ │       ┌──────▼──────┐
        │                 │course_   │ │       │ enrollments │
        │                 │staff     │ │       └──────┬──────┘
        │                 └──────────┘ │              │
        │                              │              │
        │                              │              │
        │             ┌────────────────▼───┐          │
        │             │course_plan_versions│          │
        │             └──────────┬─────────┘          │
        │                        │                    │
        ▼                        ▼                    │
┌────────────────┐      ┌────────────────┐           │
│learning_       │      │    modules     │           │
│outcomes        │      └───────┬────────┘           │
└───────┬────────┘              │                    │
        │                       ▼                    │
        │               ┌───────────────┐            │
        │               │    lessons    │            │
        │               └───────┬───────┘            │
        │                       │                    │
        │            ┌──────────┴──────────┐         │
        │            ▼                     ▼         │
        │     ┌─────────────┐      ┌────────────────┐│
        │     │  materials  │      │learning_       ││
        │     └─────────────┘      │activities      ││
        │                          └───────┬────────┘│
        │                                  │         │
        │        ┌─────────────────────────┼─────────┼─────────┐
        │        │                         │         │         │
        │        ▼                         ▼         │         ▼
        │ ┌─────────────┐          ┌────────────┐    │  ┌────────────┐
        │ │ assignments │          │  quizzes   │    │  │discussion_ │
        │ └──────┬──────┘          └─────┬──────┘    │  │forums      │
        │        │                       │           │  └────────────┘
        │        ▼                       ▼           │
        │ ┌─────────────┐        ┌───────────────┐   │
        │ │ submissions │        │quiz_attempts  │   │
        │ └──────┬──────┘        └──────┬────────┘   │
        │        │                      │            │
        │        ▼                      ▼            │
        │ ┌──────────────────┐  ┌────────────────┐   │
        │ │submission_versions│ │quiz_responses  │   │
        │ └──────────────────┘  └────────────────┘   │
        │                                             │
        │                         ┌───────────────────┘
        │                         │
        │                         ▼
        │                 ┌──────────────────┐
        │                 │attendance_records│
        │                 └────────┬─────────┘
        │                          │
        │                    ┌─────▼──────────────┐
        │                    │attendance_sessions│
        │                    └────────────────────┘
        │
        │
        │       ASSESSMENT / GRADEBOOK
        │
        │        ┌─────────────────┐
        └───────►│ activity_outcomes│
                 └─────────────────┘

┌──────────────────┐
│ grade_categories │
└────────┬─────────┘
         ▼
┌──────────────────┐
│   grade_items    │
└────────┬─────────┘
         ▼
┌──────────────────┐
│      grades      │◄──────── enrollments
└──────────────────┘

course_offerings
       │
       ▼
┌──────────────────┐
│   final_grades   │◄──────── enrollments
└────────┬─────────┘
         │
         │ outbound
         ▼
┌─────────────────────────────┐
│ SIAKAD Integration Adapter  │
└─────────────────────────────┘


                    INTEGRATION BOUNDARY

             inbound                 outbound
SIAKAD ───────────────► campus-lms ───────────────► SIAKAD
      courses                               final_grades
      offerings                             attendance
      staff
      enrollments
      CPL


SPADA / PDDikti
       │
       │ future adapter
       ▼
┌──────────────────┐
│   integrations   │
│ external_mappings│
│    sync_jobs     │
│    sync_errors   │
└──────────────────┘
```

---

# 29. Invariant database yang harus dipaksakan

Sebisa mungkin aturan berikut tidak hanya berada di application layer.

Database constraint:

```text
UNIQUE (tenant_id, user_id)
    memberships

UNIQUE (course_offering_id, student_user_id)
    enrollments

UNIQUE (course_offering_id, user_id)
    course_staff

UNIQUE (grade_item_id, enrollment_id)
    grades

UNIQUE (course_offering_id, enrollment_id)
    final_grades

UNIQUE (attendance_session_id, enrollment_id)
    attendance_records

UNIQUE (learning_activity_id, enrollment_id)
    activity_completions

UNIQUE (integration_id, entity_type, external_id)
    external_mappings
```

Gunakan check constraints untuk:

* score valid;
* time ranges;
* exactly-one owner untuk individual/group submission;
* status transition yang dapat divalidasi;
* non-negative weights/scores;
* tenant consistency jika memungkinkan.

Business invariant sensitif juga harus divalidasi di service layer karena tidak semuanya dapat dinyatakan sebagai constraint sederhana.

---

# 30. Prinsip waktu

Semua timestamp bisnis disimpan sebagai timezone-aware timestamp.

Contoh:

```text
2026-08-17T02:00:00+07:00
```

atau canonical UTC dengan timezone context tersimpan.

Tenant mempunyai `default_timezone`.

Course offering boleh mempunyai timezone override bila benar-benar diperlukan.

Deadline tidak boleh dihitung berdasarkan timezone process/server.

Entity yang membutuhkan waktu presisi:

* assignment;
* quiz;
* override;
* attendance session;
* announcement;
* material availability;
* module availability;
* course lifecycle;
* final grade publication.

---

# 31. Prinsip keamanan

Minimum domain requirements:

* deny-by-default authorization;
* RLS untuk tenant data;
* object-level authorization;
* signed/short-lived file access;
* secrets di secret manager;
* hashed attendance credential;
* immutable audit trail;
* idempotent integration;
* optimistic/concurrency control untuk grading bila diperlukan;
* server-side quiz timer;
* no correct-answer leakage;
* file malware scanning boundary;
* soft-delete/archive untuk academic records;
* correlation/request ID untuk operasi sensitif.

Super Admin bukan bypass otomatis terhadap semua RLS pada request aplikasi biasa.

Break-glass harus merupakan explicit privileged path.

---

# 32. Prinsip reporting dan learning analytics

Domain harus mampu menjawab minimal:

### Student level

* aktivitas apa yang belum selesai;
* submission apa yang belum dikumpulkan;
* quiz apa yang belum diambil;
* attendance;
* nilai yang telah dipublikasikan;
* progress course.

### Instructor level

* enrollment aktif;
* mahasiswa belum submit;
* distribution score;
* completion per activity;
* attendance;
* progress mahasiswa;
* mapping assessment terhadap CPMK;
* gradebook.

### Tenant level

* course offering aktif;
* adoption/usage;
* sync health;
* completion;
* assessment activity;
* attendance export;
* final-grade sync status.

Reporting tidak boleh membuat user melihat data course yang tidak menjadi haknya.

---

# 33. Future extensibility

Domain sengaja menyediakan boundary untuk kemampuan berikut tanpa menjadikannya dependency MVP:

* SPADA;
* PDDikti;
* SCORM;
* xAPI;
* LTI;
* video conference provider;
* plagiarism checker;
* online proctoring;
* object storage provider;
* notification provider;
* analytics warehouse.

Core entities tidak boleh didesain berdasarkan schema API satu provider tertentu.

---

# 34. Ringkasan ownership

```text
GLOBAL PLATFORM
---------------
tenants
users
auth identity
platform audit


SIAKAD-AUTHORITATIVE
--------------------
academic terms
program references
course master
course offering identity
lecturer assignment
student enrollment
CPL reference


LMS-AUTHORITATIVE
-----------------
course publication/lifecycle
RPS version
CPMK / Sub-CPMK
modules
lessons
materials
learning activities
assignments
submissions
quizzes/exams
question banks
rubrics
forums
attendance
activity completion
gradebook
grades
final grade
announcements
learning progress
audit trail


OUTBOUND TO SIAKAD
------------------
final grades
attendance
```

---

# 35. Keputusan domain yang telah dikunci

1. Produk adalah **LMS murni**, bukan SIAKAD.
2. SIAKAD adalah source of truth untuk data akademik.
3. `courses` dan `course_offerings` adalah entity berbeda.
4. LMS dibangun campus-grade sejak awal.
5. Grading configurable dan tidak hard-coded 1–100.
6. Role tambahan adalah Teaching Assistant dan Academic Operator.
7. TA boleh memberikan draft score.
8. Assessment mencakup assignment, quiz/exam, question bank, randomization, attempt rules, override, group work, rubric, formative/summative, outcome mapping, dan weighted gradebook.
9. Attendance mendukung manual, self check-in, QR/PIN, activity completion, dan external boundary.
10. Attendance terpisah dari activity completion.
11. Core interaction menggunakan announcement dan discussion forum.
12. Direct message/chat 1:1 tidak masuk core.
13. Copy course antarsemester didukung tanpa menyalin student-generated data.
14. SCORM/xAPI/LTI mempunyai integration boundary.
15. SPADA/PDDikti adalah future integration.
16. `users` global dan `memberships` tenant-scoped.
17. Course lama diarsipkan, bukan di-hard-delete.
18. Submission, grade, attendance, dan audit history tidak boleh di-hard-delete melalui aplikasi normal.
19. Hanya Lead Instructor yang dapat publish/lock final grade.
20. LMS mendukung outbound final grade dan attendance ke SIAKAD.

---

# 36. Dasar desain yang diverifikasi

Panduan Penyelenggaraan Pembelajaran Jarak Jauh di Perguruan Tinggi yang diterbitkan melalui SPADA/Kemdiktisaintek pada April 2026 menempatkan LMS sebagai sistem untuk mengelola mahasiswa/dosen, sumber belajar, aktivitas belajar, interaksi, asesmen/feedback, monitoring dan reporting. Panduan yang sama juga memasukkan integrasi sistem akademik serta dukungan standar seperti SCORM, xAPI, dan LTI, serta membahas CPL, CPMK, RPS, learning resources dan learning activities.

Struktur `course -> module/lesson -> materials/activities`, quiz/exam, assignment, forum, presensi, serta keberadaan asisten akademik konsisten dengan fitur yang didokumentasikan oleh LMS Edunex ITB. Dokumentasi Edunex juga secara eksplisit menyediakan import/duplikasi data course.

Model question bank, random question, individual override, activity completion untuk pelaporan kehadiran, attendance export, backup/restore dan import resource/activity konsisten dengan dokumentasi resmi CeLOE Telkom University.

Pemisahan LMS dari source data akademik juga didukung praktik LMS UNIKOM, yang menyinkronkan data perwalian/kelas dari sistem akademik dan mendukung user override untuk tugas atau ujian susulan.

UGM pernah mendokumentasikan integrasi SIMASTER dengan eLOK untuk data registrasi mata kuliah dan hasil nilai, sehingga pola inbound academic data + outbound result mempunyai preseden implementasi pada perguruan tinggi Indonesia.

SPADA mengumumkan perubahan API pada 18 Desember 2025 yang berdampak pada integrasi LMS perguruan tinggi. Karena itu SPADA sengaja ditempatkan di belakang adapter dan tidak dijadikan bentuk schema core LMS.

Untuk konteks regulasi, desain ini memperhatikan Permendiktisaintek Nomor 39 Tahun 2025 tentang Penjaminan Mutu Pendidikan Tinggi beserta perubahan melalui Permendiktisaintek Nomor 10 Tahun 2026.

---

# 37. Implementation Tiers

`domain.md` adalah **dokumen arsitektur**, bukan backlog rilis. Seluruh isinya
tetap berlaku sebagai desain. Bagian ini menyatakan secara eksplisit apa yang
dibangun dalam 12 minggu pertama dan apa yang sengaja hanya didokumentasikan.

Menyatakan batas ini adalah keputusan teknik, bukan pengurangan ambisi.
Membangun 68 tabel dalam 12 minggu berarti tidak pernah sampai ke deployment,
CI/CD, observability, dan lapisan AI — bagian yang justru membuktikan kemampuan
rekayasa.

## 37.1 Tier 1 — dibangun (26 tabel)

| Minggu | Tabel | Skill yang dibuktikan |
|---|---|---|
| 3 | `tenants`, `users`, `auth_identities`, `memberships`, `membership_roles`, `audit_logs`, `academic_terms`, `courses`, `course_offerings`, `course_staff`, `enrollments` | Multi-tenancy, RLS, composite FK, migrasi versioned, index & EXPLAIN |
| 4 | `auth_sessions` | Refresh token rotation, revocation |
| 5 | `modules`, `lessons`, `materials`, `files` | Relasi bertingkat, object storage, N+1 |
| 6 | `learning_activities`, `assignments`, `submissions`, `submission_versions`, `submission_files`, `grade_items`, `grades` | Supertype/subtype, versioning immutable, keyset pagination, caching |
| 7 | `ai_interactions`, `ai_artifacts` | Cost tracking, structured output, human review |
| 8 | `material_chunks` | pgvector, hybrid search, reranking, sitasi |

Tier 1 sudah cukup membuktikan seluruh kemampuan yang dinilai pasar:
isolasi multi-tenant yang teruji, versioning, penanganan file, background job,
tuning query, dan RAG produksi.

## 37.2 Tier 2 — dibangun bila waktu memungkinkan (Minggu 9–12)

```text
agent_actions
ai_quotas
attendance_sessions, attendance_credentials, attendance_records
activity_completions
announcements
grade_categories
tenant_settings
```

Diambil sesuai sisa waktu. `agent_actions` diprioritaskan karena menopang
deliverable agent Minggu 10.

## 37.3 Tier 3 — didokumentasikan, tidak dibangun

```text
learning_outcomes, outcome_mappings, course_plan_versions, course_plan_outcomes
lesson_outcomes, activity_outcomes
question_banks, questions, question_versions, question_options
quizzes, quiz_question_rules, quiz_overrides, quiz_attempts,
quiz_attempt_questions, quiz_responses
rubrics, rubric_criteria, rubric_levels
course_groups, course_group_members
assignment_overrides
discussion_forums, discussion_threads, discussion_posts
announcement_receipts
course_copy_jobs
integrations, external_mappings, sync_jobs, sync_errors
learning_packages, external_tools, learning_events
academic_program_refs
grade_schemes, final_grades
platform_admins, platform_audit_logs
external_identities
```

Tetap menjadi bagian desain. Saat ditanya di wawancara, jawaban yang benar
adalah: *"sudah dimodelkan, sengaja tidak diimplementasikan pada rilis pertama
karena tidak menambah bukti kemampuan rekayasa yang belum dibuktikan Tier 1."*

## 37.4 Aturan tier

1. Tabel Tier 3 tidak boleh muncul di migrasi. Tidak ada tabel kosong "untuk nanti".
2. Naik tier butuh keputusan tertulis di laporan mingguan.
3. Skema Tier 1 harus dirancang agar Tier 2 dan 3 dapat ditambahkan tanpa
   merombak — itulah gunanya `learning_activities` sebagai supertype.
4. Bila Tier 1 selesai lebih cepat, prioritas berikutnya adalah **memperdalam**
   (observability, eval, keamanan), bukan menambah tabel.

---

# 38. Amendemen terhadap keputusan yang telah dikunci

Delapan perubahan berikut mengoreksi atau melengkapi bagian sebelumnya.

## A1 — Multi-role dalam satu tenant

**Masalah:** `UNIQUE (tenant_id, user_id)` pada `memberships` memaksa satu peran
per tenant. Kaprodi yang juga mengajar (`academic_operator` + `lecturer`) adalah
kombinasi yang lazim di kampus Indonesia.

**Keputusan:** peran dipindahkan ke tabel terpisah.

```text
memberships
    id, tenant_id, user_id, status, joined_at
    UNIQUE (tenant_id, user_id)      -- tetap: satu membership per tenant

membership_roles
    id, tenant_id, membership_id, role, granted_by, granted_at, revoked_at
    UNIQUE (membership_id, role) WHERE revoked_at IS NULL
```

Konsekuensi:

* §24.2.1–24.2.3 tetap berlaku.
* Otorisasi memakai **union** dari seluruh peran aktif.
* Pencabutan peran adalah `revoked_at`, bukan `DELETE` — jejak audit terjaga.
* `teaching_assistant` tetap tidak menjadi peran tenant; tetap lewat `course_staff`.

## A2 — Penegakan konsistensi tenant di database

**Masalah:** §24.1.2 mewajibkan setiap FK antar-entitas tenant menunjuk tenant
yang sama, tanpa menyebut mekanismenya.

**Keputusan:** pakai composite foreign key.

```sql
-- setiap tabel tenant-scoped
ALTER TABLE course_offerings ADD UNIQUE (tenant_id, id);

-- setiap child
ALTER TABLE modules
  ADD CONSTRAINT modules_tenant_consistency
  FOREIGN KEY (tenant_id, course_offering_id)
  REFERENCES course_offerings (tenant_id, id);
```

Dengan pola ini referensi lintas tenant menjadi **mustahil secara struktural**,
bukan sekadar dilarang di service layer. Wajib untuk seluruh tabel Tier 1.

## A3 — `tenant_id` eksplisit pada tabel anak

**Masalah:** beberapa tabel di daftar RLS §25 tidak mencantumkan `tenant_id`
pada atributnya.

**Keputusan:** setiap tabel tenant-scoped **wajib** menyimpan `tenant_id`
langsung, meski dapat diturunkan lewat parent. Terdampak:

```text
announcement_receipts   course_group_members   rubric_criteria
rubric_levels           question_options       quiz_attempt_questions
quiz_responses          submission_files       lesson_outcomes
activity_outcomes       course_plan_outcomes
```

Alasan: RLS pada join berantai mahal dan mudah salah. Redundansi kecil ini
membeli kebenaran dan performa.

## A4 — Kehadiran sebagai komponen nilai

**Masalah:** `grade_items.source_type` hanya `assignment | quiz | manual`,
padahal "kehadiran 10%" adalah praktik standar di kampus Indonesia.

**Keputusan:** tambahkan `attendance` sebagai source type. Nilainya dihitung
dari `attendance_records` pada satu course offering menurut aturan yang
dikonfigurasi course (mis. persentase hadir). Tetap tunduk pada §24.9 —
publikasi hanya oleh Lead Instructor.

## A5 — Referensi polimorfik pada `grade_items`

**Masalah:** pasangan `source_type` + `source_id` tidak dapat dijaga foreign key.

**Keputusan:** ganti dengan kolom bertipe yang nullable plus check.

```sql
assignment_id  uuid NULL REFERENCES assignments(id),
quiz_id        uuid NULL REFERENCES quizzes(id),
attendance_ref uuid NULL,
CHECK (num_nonnulls(assignment_id, quiz_id, attendance_ref) <= 1)
```

`manual` ditandai dengan seluruh kolom bernilai NULL. Menambah kolom saat tipe
baru muncul lebih murah daripada kehilangan integritas referensial.

## A6 — Keunikan dua arah pada `external_mappings`

**Masalah:** hanya ada `UNIQUE (integration_id, entity_type, external_id)`,
sehingga satu entity internal bisa dipetakan ke dua external ID.

**Keputusan:** tambahkan

```text
UNIQUE (integration_id, entity_type, internal_id)
```

## A7 — Enrollment wajib menunjuk membership aktif

**Masalah:** §24.1.4 hanya berupa aturan tertulis.

**Keputusan:** tegakkan di database.

```sql
FOREIGN KEY (tenant_id, student_user_id)
  REFERENCES memberships (tenant_id, user_id)
```

Status aktif tetap divalidasi di service layer, tetapi keberadaan membership
dalam tenant yang sama dijamin database.

## A8 — Entitas sesi autentikasi

**Masalah:** §31 menuntut akses berumur pendek dan revocation, tetapi tidak ada
tempat menyimpan refresh token yang dapat dicabut.

**Keputusan:** tambahkan tabel global.

```text
auth_sessions
    id, user_id, refresh_token_hash, issued_at, expires_at,
    rotated_from, revoked_at, revoked_reason,
    user_agent, ip_address, last_seen_at
```

Token disimpan sebagai hash. Rotasi membuat baris baru dan mengisi
`rotated_from`. Pemakaian ulang token yang sudah dirotasi diperlakukan sebagai
indikasi pencurian: seluruh sesi milik user dicabut.

---

## 38.9 Tambahan pada §35 — keputusan yang dikunci

```text
21. Satu user dapat memegang lebih dari satu peran dalam satu tenant,
    melalui membership_roles.
22. Konsistensi tenant ditegakkan di database melalui composite foreign key,
    bukan hanya di service layer.
23. Setiap tabel tenant-scoped menyimpan tenant_id secara eksplisit.
24. Kehadiran dapat menjadi komponen nilai melalui grade_items.
25. Lapisan AI berada di dokumen terpisah (docs/domain-ai.md) dan tidak boleh
    menjadi dependency domain pembelajaran inti.
26. Implementasi mengikuti tier pada §37; tabel Tier 3 tidak boleh muncul
    di migrasi.
```
