# Domain Model — campus-lms

**MINGGU 1 — TULIS SENDIRI.** Sebelum satu baris SQL pun ditulis.

Menulis ini duluan adalah kebiasaan engineer senior: memahami masalah
sebelum memilih solusi. Kalau kamu langsung membuat tabel, kamu akan
merombaknya tiga kali di Minggu 3.

## Yang harus kamu isi

### 1. Aktor dan perannya
- Super admin (pengelola platform SaaS)
- Admin kampus
- Dosen
- Mahasiswa
- (ada lagi? asisten dosen? orang tua? auditor akreditasi?)

Untuk tiap aktor: apa yang boleh dan tidak boleh dilakukan.

### 2. Entitas inti dan relasinya
Minimal: `tenants, users, memberships, courses, modules, lessons,
materials, enrollments, assignments, submissions, grades, announcements,
audit_logs`.

Untuk tiap entitas: atribut penting, siapa pemiliknya, dan **apakah
ber-tenant atau global** (ini menentukan RLS di Minggu 3).

### 3. Aturan bisnis yang tidak boleh dilanggar
Contoh untuk memancing pikiranmu:
- Mahasiswa hanya boleh submit sebelum deadline — kecuali dosen memberi
  perpanjangan individual.
- Nilai final hanya boleh diubah dosen pengampu, dan setiap perubahan
  wajib tercatat di audit log.
- Satu mahasiswa tidak boleh terdaftar dua kali di kursus yang sama.
- Data antar-tenant TIDAK BOLEH bocor dalam kondisi apa pun.

### 4. Pertanyaan yang belum terjawab
Tulis kebingunganmu. Ini bagian paling jujur dan paling berguna — nanti
kamu akan senang bisa melihat bagaimana pemahamanmu berkembang.

## Diagram

TODO: ERD sederhana. Boleh ASCII dulu, diperbagus di Minggu 12.
