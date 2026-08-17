# Laporan Minggu 02 — Docker & Compose sampai ke Tulang

> Disusun agent, ditandatangani manusia.

- **Periode:** 2026-08-18
- **Fokus roadmap:** Dockerfile API multi-stage, runtime minimal, build multi-arch, dan kebersihan build context.
- **Total jam:** NOT MEASURED
- **Commit range:** `a8c99bf..3bb7565` (3 commit)

---

## 1. Ringkasan

**FACT:** `apps/api/Dockerfile` sekarang membangun binary Go statis dalam stage builder dan menjalankannya pada runtime distroless sebagai user non-root. Runtime API-probe berukuran `14.6MB`, di bawah target `< 25MB`; target pembanding BusyBox berukuran `16.6MB` dan runtime `golang:alpine` berukuran `377MB`. Whitelist `apps/api/.dockerignore` menurunkan transfer build context yang dilaporkan BuildKit dari `942B` menjadi `621B` (`321B`, atau `34.1%`). Buildx juga membangun manifest OCI untuk `linux/amd64` dan `linux/arm64` tanpa push ke registry. Compose dan aplikasi web sengaja tidak dikerjakan dalam sesi ini sesuai scope yang ditetapkan.

---

## 2. Dikerjakan Agent

| # | Pekerjaan | File | Commit | Bukti |
|---|---|---|---|---|
| 1 | Menambahkan Dockerfile API multi-stage, static build, runtime distroless non-root, health check, dan target pembanding Alpine | `apps/api/Dockerfile` | `ad1bc11` | [ukuran](evidence/week-02/image-size.txt), [health check](evidence/week-02/healthcheck-distroless.txt), [multi-arch](evidence/week-02/buildx-multiarch.txt) |
| 2 | Membatasi Docker build context ke metadata modul dan source produksi | `apps/api/.dockerignore` | `8ac8e6b` | [sebelum](evidence/week-02/build-context-before.txt), [sesudah](evidence/week-02/build-context-after.txt) |
| 3 | Menambahkan mode `-healthcheck` dan mengubah runtime akhir agar tidak menyalin BusyBox | `apps/api/cmd/api/main.go`, `apps/api/internal/healthcheck/probe.go`, `apps/api/Dockerfile` | `3bb7565` | [test](evidence/week-02/go-test-healthcheck.txt), [runtime](evidence/week-02/api-probe-healthcheck.txt), [ukuran](evidence/week-02/busybox-vs-api-probe-size.txt) |

**Catatan implementasi:**

- **FACT:** `COPY go.mod` dan `go mod download` berada sebelum `COPY . .` pada `apps/api/Dockerfile`, sehingga perubahan source tidak membatalkan cache unduhan dependensi.
- **FACT:** runtime akhir menggunakan `HEALTHCHECK CMD ["/api", "-healthcheck"]`; target `busybox-runtime` dipertahankan hanya untuk perbandingan.
- **FACT:** mode `-healthcheck` melakukan HTTP GET lokal ke `/healthz`, dibatasi timeout dua detik, dan keluar gagal untuk status selain 200.
- **FACT:** `.dockerignore` memakai whitelist: hanya `go.mod`, `go.sum` bila ada, serta `cmd/` dan `internal/` yang masuk context; file test dikecualikan setelah whitelist diterapkan.

---

## 3. Dikerjakan Manusia

| # | Pekerjaan | Kenapa harus manual | Hasil |
|---|---|---|---|
| 1 | Konfirmasi tindakan manual selama sesi ini | Agent tidak boleh menebak tindakan manusia | Menunggu konfirmasi manusia |

---

## 4. Keputusan yang Diambil

| Keputusan | Alternatif yang ditolak | Alasan | ADR |
|---|---|---|---|
| Implementasi saat ini memakai probe internal `/api -healthcheck` | Target BusyBox tetap tersedia untuk pembanding | **FACT:** runtime lokal melaporkan `healthy`, proses API ber-UID `65532`, dan endpoint `/healthz` menjawab sukses | [ADR-0005](../adr/0005-api-healthcheck-probe.md) — Decision dan Consequences harus ditulis manusia |

**Alternatif untuk dicatat pada ADR-0005 (belum merupakan keputusan manusia):**

| Alternatif | Alasan tidak dipakai pada implementasi saat ini | Trade-off yang perlu diputuskan manusia |
|---|---|---|
| Shell-form `HEALTHCHECK CMD wget ...` di distroless | Distroless tidak menyediakan `/bin/sh`; command shell-form tidak dapat dimulai | Tidak ada binary probe tambahan, tetapi tidak memenuhi health check HTTP |
| `golang:1.23.6-alpine3.21` sebagai runtime produksi | **FACT:** hasil ukur lokal 377MB, dibanding 15.9MB untuk distroless | Shell dan toolchain memudahkan debugging, tetapi image lebih besar dan membawa lebih banyak tool runtime |
| Binary probe Go khusus | Diimplementasikan untuk dibuktikan; belum menjadi Decision ADR manusia | Mengurangi utility runtime, tetapi menambah kode dan test yang harus dipelihara |
| Tidak memasang `HEALTHCHECK` | Melanggar TASK BRIEF Minggu 2 | Image paling sederhana, tetapi Docker tidak dapat mendeteksi endpoint API yang gagal |

---

## 5. Angka & Bukti

| Metrik | Nilai | Cara diukur | File bukti |
|---|---|---|---|
| Ukuran image API-probe distroless | 14.6MB | `docker images campus-lms-api --format=...` | [image-size.txt](evidence/week-02/image-size.txt) |
| Ukuran image runtime BusyBox | 16.6MB | `docker images campus-lms-api --format=...` | [image-size.txt](evidence/week-02/image-size.txt) |
| Ukuran pembanding runtime `golang:alpine` | 377MB | `docker images campus-lms-api --format=...` | [image-size.txt](evidence/week-02/image-size.txt) |
| UID efektif container API | 65532 | `docker inspect` lalu `docker exec ... id -u` | [docker-inspect-user.txt](evidence/week-02/docker-inspect-user.txt) |
| Platform build | linux/amd64 dan linux/arm64 | `docker buildx build --platform ... --output type=oci` | [buildx-multiarch.txt](evidence/week-02/buildx-multiarch.txt) |
| Runtime health check API-probe | `healthy`; UID 65532; `/healthz` menjawab JSON | `docker inspect`, `docker top`, dan `curl` | [api-probe-healthcheck.txt](evidence/week-02/api-probe-healthcheck.txt) |
| Transfer build context sebelum `.dockerignore` | 942B | Buildx builder bersih, `load build context` | [build-context-before.txt](evidence/week-02/build-context-before.txt) |
| Transfer build context sesudah `.dockerignore` | 621B | Buildx builder bersih, `load build context` | [build-context-after.txt](evidence/week-02/build-context-after.txt) |

**Perbandingan sebelum/sesudah (runtime):**

| Metrik | Sebelum: `golang:alpine` runtime | Sesudah: distroless runtime | Perubahan | Bukti |
|---|---:|---:|---:|---|
| Ukuran image API | 377MB | 14.6MB | 362.4MB lebih kecil | [image-size.txt](evidence/week-02/image-size.txt) |
| Runtime health check | BusyBox: 16.6MB | API-probe: 14.6MB | 2.0MB lebih kecil dalam tampilan `docker images` | [image-size.txt](evidence/week-02/image-size.txt) |
| Transfer build context | 942B | 621B | 321B (34.1%) lebih kecil | [sebelum](evidence/week-02/build-context-before.txt), [sesudah](evidence/week-02/build-context-after.txt) |

---

## 6. Konsep yang Dipelajari

### Multi-stage build dan runtime distroless

- **Apa:** multi-stage build memisahkan alat kompilasi dari image yang menjalankan program.
- **Kenapa dipakai di sini:** builder membutuhkan Go, sedangkan stage akhir di `apps/api/Dockerfile` hanya menerima `/api`; BusyBox dipisahkan ke target pembanding.
- **Alternatif yang tidak dipilih:** runtime `golang:alpine` lebih mudah dieksplorasi karena memiliki shell dan toolchain, tetapi hasil ukur lokalnya 377MB.
- **Cara membuktikan sendiri:** `make docker-build && make docker-size`.
- **Pertanyaan interview terkait:** Mengapa image runtime tidak sebaiknya membawa compiler dan shell secara default?

### CGO dan cross-compilation

- **Apa:** `CGO_ENABLED=0` meminta Go menghasilkan binary tanpa dependensi pustaka C dinamis.
- **Kenapa dipakai di sini:** perintah build pada `apps/api/Dockerfile` menetapkan `CGO_ENABLED=0`, `GOOS`, dan `GOARCH` agar binary dapat dibawa ke runtime minimal dan dibangun untuk dua arsitektur.
- **Alternatif yang tidak dipilih:** binary dengan CGO dapat memerlukan libc yang cocok pada runtime; biaya debugging dan kompatibilitasnya lebih tinggi.
- **Cara membuktikan sendiri:** `docker buildx build --platform linux/amd64,linux/arm64 --output type=oci,dest=/tmp/campus-lms-api-multiarch.tar -f apps/api/Dockerfile apps/api`.
- **Pertanyaan interview terkait:** Mengapa `CGO_ENABLED=0` membantu saat memakai image scratch atau distroless?

### Urutan COPY dan cache layer

- **Apa:** Docker menggunakan hasil instruction sebelumnya sebagai cache layer berikutnya.
- **Kenapa dipakai di sini:** `go.mod` disalin dan dependensi diunduh sebelum seluruh source disalin pada `apps/api/Dockerfile`.
- **Alternatif yang tidak dipilih:** menyalin seluruh source sebelum `go mod download`; setiap perubahan source akan membuat layer unduhan dependensi dibangun ulang.
- **Cara membuktikan sendiri:** ubah satu file Go, lalu jalankan `make docker-build` dan amati apakah step `go mod download` berstatus `CACHED`.
- **Pertanyaan interview terkait:** Mengapa urutan `COPY` dapat lebih memengaruhi kecepatan rebuild daripada pilihan base image?

### Docker build context

- **Apa:** build context adalah kumpulan file yang dikirim Docker client ke builder sebelum instruction `COPY` dapat dijalankan.
- **Kenapa dipakai di sini:** `.dockerignore` di `apps/api/` mengirim hanya file yang dibutuhkan `go build`, sehingga test, script, migration placeholder, dan artefak lokal tidak ikut terkirim.
- **Alternatif yang tidak dipilih:** blacklist satu per satu; pendekatan itu mudah lupa diperbarui saat jenis artefak baru muncul.
- **Cara membuktikan sendiri:** bandingkan `load build context` pada [build-context-before.txt](evidence/week-02/build-context-before.txt) dan [build-context-after.txt](evidence/week-02/build-context-after.txt).
- **Pertanyaan interview terkait:** Mengapa `.dockerignore` bisa mempercepat build bahkan ketika image final tidak berubah?

---

## 7. Belum Terverifikasi

| Hal | Kenapa belum terverifikasi | Rencana verifikasi |
|---|---|---|
| `make up` dan networking antarlayanan | Compose berada di luar scope sesi ini | Sesi Compose Minggu 2 berikutnya |
| Push manifest multi-arch ke GHCR | Remote/registry belum dikonfirmasi manusia | Setelah manusia mengonfirmasi remote tersedia |
| Batas memori 256MB dan OOM-kill terkontrol | Compose dan limit runtime belum dikerjakan | Sesi Compose Minggu 2 berikutnya |
| Perilaku health check pada VM produksi | Baru diverifikasi pada Docker Engine lokal | Verifikasi saat deploy Minggu 4 |
| ADR pilihan distroless | Decision dan Consequences adalah milik manusia | Manusia menulis ADR setelah menjelaskan trade-off |

**Asumsi yang dipakai tapi belum dibuktikan:**

- Perilaku mode `-healthcheck` dengan dependency nyata pada Week 3; saat ini `/healthz` hanya menguji liveness proses.

---

## 8. Masalah & Cara Diselesaikan

### Masalah: target `make evidence` gagal saat format Docker diberi tanda kutip

- **Gejala:** `/bin/bash: line 1: test: too many arguments`, diikuti `usage: make evidence W=01 SLUG=image-size CMD="docker images"`.
- **Hipotesis yang salah:** Docker inspect atau container non-root yang gagal.
- **Akar masalah:** nilai `CMD` pada target Make diekspansi di dalam tanda kutip shell; tanda kutip internal pada format `docker inspect` memecah argumen `test`.
- **Solusi:** menjalankan format Go template tanpa tanda kutip internal: `--format={{.Config.User}}`.
- **Pencegahan:** gunakan `CMD` tanpa tanda kutip bersarang untuk target evidence yang ada; jangan ubah Makefile hanya agar perekaman bukti lolos.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: port host 18080 sudah dipakai saat container verifikasi dijalankan

- **Gejala (raw error):**

  ```text
  docker: Error response from daemon: failed to set up container networking: driver failed programming external connectivity on endpoint campus-lms-api-healthcheck-verify (3b9838eec6ea86a73ae00a5f91f22f05ba0a9680083487def4d926bd96db08d8): Bind for 0.0.0.0:18080 failed: port is already allocated

  Run 'docker run --help' for more information
  ```

- **Hipotesis yang salah:** Dockerfile distroless atau BusyBox health check gagal dibuat.
- **Akar masalah:** container verifikasi lama `campus-lms-api-verify` masih memetakan host port `18080` ke port container `8080`.
- **Solusi:** hapus kedua container verifikasi yang ditargetkan secara eksplisit, kemudian jalankan ulang container baru.
- **Pencegahan:** sebelum memakai port tetap untuk verifikasi, periksa `docker ps` dan hapus container verifikasi yang sudah selesai.
- **Waktu terbuang:** NOT MEASURED.

---

## 9. Status Definition of Done

| DoD dari roadmap | Usulan agent | Bukti | Dicentang manusia |
|---|---|---|---|
| `make up` membuat seluruh stack sehat dan web memakai nama service | ❌ belum — Compose/web di luar scope | Tidak ada | ☐ |
| Image API < 25MB dan UID non-root | ✅ kandidat terpenuhi secara lokal | [ukuran](evidence/week-02/docker-size-api-probe.txt), [runtime](evidence/week-02/api-probe-healthcheck.txt) | ☐ |
| Build amd64/arm64 sukses dan ter-push ke GHCR | ⚠️ sebagian — build lokal sukses, tidak ada push | [buildx-multiarch.txt](evidence/week-02/buildx-multiarch.txt) | ☐ |
| Container dibatasi 256MB dan OOM-kill dijelaskan | ❌ belum | Tidak ada | ☐ |
| ADR trade-off distroless tersedia | ⚠️ Context dan Options tersedia; Decision menunggu manusia | [ADR-0005](../adr/0005-api-healthcheck-probe.md) | ☐ |

---

## 10. Untuk Minggu Depan

- **Carry-over:** Docker Compose, networking nama service, limit memori 256MB, dan eksperimen OOM-kill.
- **Utang teknis yang sengaja diambil:** target `alpine-runtime` dipertahankan untuk perbandingan ukuran Minggu 2; hapus atau pindahkan setelah manusia memutuskan bentuk dokumentasi eksperimennya.
- **Persiapan yang perlu dilakukan manusia lebih dulu:** konfirmasi remote GHCR sebelum ada perintah push dan tulis ADR distroless setelah memilih trade-off-nya.

---

## 11. Verifikasi Manusia

- [ ] Saya sudah spot-check 3 file bukti secara acak dan isinya cocok dengan klaim
- [ ] Skor quiz: ____ / ____ (minimal 70% untuk lanjut)
- [ ] Explain-back 3 menit sudah direkam: `docs/progress/explain/week-02.<mp3|txt>`
- [ ] Saya bisa menjelaskan setiap keputusan di section 4 tanpa membuka catatan

**Self-audit agent (dari `agent/evidence-protocol.md` §8):** Bukti angka ditautkan; file yang dibuat dibaca ulang; build dan inspeksi yang diklaim dijalankan dalam sesi ini; section Belum Terverifikasi diisi; tidak ada checkbox DoD yang dicentang; label FACT/INFERENCE/RECOMMENDATION digunakan saat diperlukan. Quiz Minggu 2 masih template dan belum diisi.

**Ditandatangani:** ______________  **Tanggal:** __________
