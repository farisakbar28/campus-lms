# Laporan Minggu 02 — Docker & Compose sampai ke Tulang

> Template. Salin ke `docs/progress/week-02.md`. Disusun agent, ditandatangani manusia.
> Jangan hapus satu section pun. Section yang tidak relevan diisi "Tidak ada minggu ini".

- **Periode:** 2026-08-18 s/d 2026-08-18
- **Fokus roadmap:** Dockerfile API multi-stage, Compose core/dev, networking, resource limit, dan kebersihan runtime.
- **Total jam:** NOT MEASURED.
- **Commit range:** `a8c99bf..2fbb3b4` ([14 commit](evidence/week-02/report-source-git-log.txt)).

---

## 1. Ringkasan

**FACT:** API mempunyai Dockerfile multi-stage dengan runtime distroless, user non-root, dan `HEALTHCHECK` yang memakai mode `-healthcheck` pada binary sendiri. **FACT:** Compose core sekarang menjalankan API, Postgres, dan Redis dalam keadaan sehat pada verifikasi manusia; override development menjalankan Air untuk hot reload. **FACT:** image API yang diukur terakhir berukuran [14.6MB](evidence/week-02/image-size.txt), di bawah target roadmap. Scope web, storage, dan observability belum selesai, sehingga Minggu 02 belum dapat dianggap selesai sepenuhnya.

## 2. Dikerjakan Agent

| # | Pekerjaan | File | Commit | Bukti |
|---|---|---|---|---|
| 1 | Menambahkan modul awal sebagai dasar build Go tanpa dependency eksternal | `apps/api/go.mod` | `a8c99bf` | [build runtime](evidence/week-02/api-probe-runtime-build.txt) |
| 2 | Menambahkan Dockerfile multi-stage: static build, distroless, non-root, dan target pembanding | `apps/api/Dockerfile` | `ad1bc11` | [build multi-arch](evidence/week-02/buildx-multiarch.txt), [UID](evidence/week-02/docker-inspect-user.txt) |
| 3 | Membatasi context build dengan whitelist `.dockerignore` | `apps/api/.dockerignore` | `8ac8e6b` | [sebelum](evidence/week-02/build-context-before.txt), [sesudah](evidence/week-02/build-context-after.txt) |
| 4 | Mengganti probe BusyBox pada image akhir dengan `/api -healthcheck`, beserta unit test-nya | `apps/api/cmd/api/main.go`, `apps/api/internal/healthcheck/`, `apps/api/Dockerfile` | `3bb7565` | [test historis](evidence/week-02/go-test-healthcheck.txt), [runtime probe](evidence/week-02/api-probe-healthcheck.txt) |
| 5 | Menulis opsi probe untuk keputusan manusia dan mendokumentasikan eksperimen internals container | `docs/adr/0005-api-healthcheck-probe.md`, `docs/notes/docker-internals.md` | `db71c5c`, `4543ced` | [history layer](evidence/week-02/busybox-layer-history.txt) |
| 6 | Menetapkan penundaan web dari profile core ke Minggu 04 dan membuat draf awal laporan/bukti | `deploy/compose/docker-compose.yml`, `docs/progress/` | `ae32c04`, `4543ced` | [log commit](evidence/week-02/report-source-git-log.txt) |

**Catatan implementasi:**

- **FACT:** `COPY go.mod go.sum* ./` dan `go mod download` mendahului `COPY . ./` dalam `apps/api/Dockerfile`; wildcard membuat build tetap berjalan ketika `go.sum` belum ada karena modul belum memiliki dependency.
- **FACT:** `healthcheck.Probe` menggunakan `http.NewRequestWithContext`, timeout dua detik, dan hanya menerima HTTP 200 (`apps/api/internal/healthcheck/probe.go`).
- **INFERENCE:** memisahkan target `busybox-runtime` dari runtime akhir menjaga pembanding tetap dapat direproduksi tanpa memasukkan BusyBox ke image produksi.

## 3. Dikerjakan Manusia

| # | Pekerjaan | Kenapa harus manual | Hasil |
|---|---|---|---|
| 1 | Menghasilkan seluruh perubahan pada rentang `3f559f56e922131dd8dc5daf4719384462186b56` sampai `2fbb3b467f2f464f7b901832813f538bf76823f6`, baik tindakan langsung maupun lewat command Bash | Klasifikasi pelaku dan tindakan manual ditetapkan pemilik | Tujuh commit tercatat dalam [log commit](evidence/week-02/report-source-git-log.txt) |
| 2 | Memigrasikan lingkungan dari Docker Desktop ke Docker Engine | Keputusan tooling host berada pada pemilik; alasan yang dinyatakan adalah batas RAM laptop 8GB | Docker Engine dipilih; penghematan RAM belum terukur dari WSL |
| 3 | Menjalankan verifikasi Compose, DNS, volume, limit memori, hot reload, OOM-kill, dan inspeksi GHCR | Sandbox agent tidak diberi akses ke `/var/run/docker.sock` | Bukti verifikasi manusia tersimpan pada file Compose, hot reload, OOM, dan GHCR di section 5 |
| 4 | Memilih dan merekam keputusan final probe API pada ADR | Decision dan Consequences ADR adalah milik manusia | `ADR-0005` memuat keputusan menggunakan `-healthcheck` |

## 4. Keputusan yang Diambil

| Keputusan | Alternatif yang ditolak | Alasan | ADR |
|---|---|---|---|
| Runtime produksi memakai distroless dengan `/api -healthcheck` | Menyalin BusyBox ke image akhir | **FACT:** image akhir terakhir diukur [14.6MB](evidence/week-02/image-size.txt); BusyBox menambah layer [1.22MB](evidence/week-02/busybox-layer-history.txt). **INFERENCE:** probe dalam binary cukup untuk kebutuhan HTTP dan mempertahankan runtime lebih minimal. | [ADR-0005](../adr/0005-api-healthcheck-probe.md) |
| Docker Desktop diganti Docker Engine | Tetap memakai Docker Desktop | **FACT (pernyataan manusia):** keputusan dibuat karena laptop memiliki RAM 8GB. Penghematan RAM tidak diukur dari WSL. | Tidak ada minggu ini |
| `make prune` aman untuk data volume; penghapusan volume dipisah ke `make prune-hard` dengan konfirmasi | `docker system prune -af --volumes` pada target harian | **FACT:** target aman tidak memakai `--volumes`; target keras meminta teks `DESTROY` sebelum menjalankan perintah berbahaya. | Tidak ada minggu ini |
| Override development memakai Air `v1.61.7` | Air terbaru | **FACT (pernyataan manusia):** versi terbaru memerlukan Go 1.25, sedangkan lingkungan ini memakai Go 1.23.6; versi yang dipilih terlihat pada [log hot reload](evidence/week-02/hot-reload-demo.txt). | Tidak ada minggu ini |

## 5. Angka & Bukti

| Metrik | Nilai | Cara diukur | File bukti |
|---|---|---|---|
| Commit dalam rentang laporan | 14 | `git log` dan `git rev-list --count` | [report-source-git-log.txt](evidence/week-02/report-source-git-log.txt) |
| Image API terakhir | 14.6MB | `docker images` | [image-size.txt](evidence/week-02/image-size.txt) |
| UID efektif API | 65532 | `docker inspect` dan `docker exec id -u` | [docker-inspect-user.txt](evidence/week-02/docker-inspect-user.txt) |
| Runtime API tanpa BusyBox | `healthy`; user `65532` | `docker inspect`, `docker top`, dan `curl /healthz` | [api-probe-healthcheck.txt](evidence/week-02/api-probe-healthcheck.txt) |
| Runtime distroless dengan probe BusyBox (pembanding) | `healthy`; UID 65532 | `docker inspect`, `docker exec`, dan `curl /healthz` | [healthcheck-distroless.txt](evidence/week-02/healthcheck-distroless.txt) |
| Build runtime akhir tanpa BusyBox | berhasil | `docker build --no-cache` | [api-probe-runtime-build.txt](evidence/week-02/api-probe-runtime-build.txt) |
| Build pembanding BusyBox | berhasil | `docker build --target busybox-runtime --no-cache` | [busybox-runtime-build.txt](evidence/week-02/busybox-runtime-build.txt) |
| Build Dockerfile dengan context terbatasi | berhasil | `docker build` | [docker-build-current.txt](evidence/week-02/docker-build-current.txt) |
| Output `make docker-size` | 14.6MB | `make docker-size` | [docker-size-api-probe.txt](evidence/week-02/docker-size-api-probe.txt) |
| Platform yang berhasil dibangun secara lokal | `linux/amd64` dan `linux/arm64` | `docker buildx build --platform ...` | [buildx-multiarch.txt](evidence/week-02/buildx-multiarch.txt) |
| Manifest GHCR | `linux/amd64` dan `linux/arm64` | `docker buildx imagetools inspect` | [ghcr-push.txt](evidence/week-02/ghcr-push.txt) |
| Context build sebelum/ sesudah whitelist | 942B / 621B | Buildx builder baru | [sebelum](evidence/week-02/build-context-before.txt), [sesudah](evidence/week-02/build-context-after.txt) |
| Layer BusyBox pada history | 1.22MB | `docker history --no-trunc` | [busybox-layer-history.txt](evidence/week-02/busybox-layer-history.txt) |
| Ukuran exact BusyBox/API-probe dan file BusyBox | 4,054,678B / 3,285,268B / 1,214,736B | `docker image inspect` dan `stat` | [busybox-vs-api-probe-size.txt](evidence/week-02/busybox-vs-api-probe-size.txt) |
| API, Postgres, Redis pada Compose core | ketiganya `healthy` | `docker compose ps` | [compose-ps-healthy.txt](evidence/week-02/compose-ps-healthy.txt) |
| Jalur Compose runtime produksi | API `healthy`; probe `CMD /api -healthcheck` | `docker compose ps` dan `docker inspect` | [compose-prod-path.txt](evidence/week-02/compose-prod-path.txt) |
| DNS service dalam network Compose | `postgres` dan `redis` resolve | `getent hosts` dalam network | [compose-dns.txt](evidence/week-02/compose-dns.txt) |
| Limit memori API / Postgres / Redis | 128MiB / 512MiB / 128MiB | `docker inspect .HostConfig.Memory` | [compose-memlimit.txt](evidence/week-02/compose-memlimit.txt) |
| Postgres named volume setelah `make down` | baris `id=42` masih ada | `psql SELECT` | [compose-volume-persistence.txt](evidence/week-02/compose-volume-persistence.txt) |
| Dev hot reload | Air mendeteksi perubahan lalu build dan run ulang | `docker compose logs` | [hot-reload-demo.txt](evidence/week-02/hot-reload-demo.txt) |
| Demonstrasi batas memori | `OOMKilled=true`, exit 137 | container Alpine dengan `--memory=32m` | [oom-kill-demo.txt](evidence/week-02/oom-kill-demo.txt) |
| Unit test race detector pada commit implementasi | lulus secara historis | `go test -race -count=1 ./...` | [go-test-healthcheck.txt](evidence/week-02/go-test-healthcheck.txt) |

**Perbandingan sebelum/sesudah (kalau ada optimasi):**

| Metrik | Sebelum | Sesudah | Perubahan | Bukti |
|---|---|---|---|---|
| Runtime sebelum replacement probe vs API-probe | 15.9MB | 14.6MB | turun sekitar 1.3MB pada tampilan `docker images`; konsisten setelah pembulatan dengan layer BusyBox 1.22MB | [ukuran historis](evidence/week-02/healthcheck-size-comparison-historical.txt), [history](evidence/week-02/busybox-layer-history.txt) |
| Transfer build context | 942B | 621B | turun 321B (34.1%) | [sebelum](evidence/week-02/build-context-before.txt), [sesudah](evidence/week-02/build-context-after.txt) |

## 6. Konsep yang Dipelajari

### Multi-stage build dan distroless

- **Apa:** multi-stage memisahkan toolchain build dari runtime agar image akhir hanya membawa artefak yang dibutuhkan.
- **Kenapa dipakai di sini:** builder Go ada di `apps/api/Dockerfile`, sedangkan stage akhir memakai `gcr.io/distroless/static-debian12:nonroot` dan `/api`.
- **Alternatif yang tidak dipilih:** runtime `golang:alpine` lebih mudah diinspeksi, tetapi bukti historis mencatatnya 377MB ([ukuran historis](evidence/week-02/healthcheck-size-comparison-historical.txt)).
- **Cara membuktikan sendiri:** `make docker-build && make docker-size`.
- **Pertanyaan interview terkait:** Mengapa `CGO_ENABLED=0` penting ketika menjalankan binary Go pada image distroless atau scratch?

### HEALTHCHECK tanpa shell

- **Apa:** Docker HEALTHCHECK menjalankan command berkala dan menandai container unhealthy bila command tersebut gagal.
- **Kenapa dipakai di sini:** distroless tidak menyediakan shell maupun HTTP client; `apps/api/cmd/api/main.go` menerima flag `-healthcheck` dan `healthcheck.Probe` memanggil `/healthz`.
- **Alternatif yang tidak dipilih:** BusyBox menyediakan `wget`, tetapi membawa layer tambahan [1.22MB](evidence/week-02/busybox-layer-history.txt).
- **Cara membuktikan sendiri:** `docker inspect campus-lms-api:dev --format '{{.Config.Healthcheck.Test}}'`.
- **Pertanyaan interview terkait:** Apa risiko memakai liveness check yang gagal saat dependency eksternal sedang down?

### Compose: service DNS, health dependency, dan volume

- **Apa:** Compose membuat network internal yang memberi DNS berdasarkan nama service dan mengelola volume bernama terpisah dari lifecycle container.
- **Kenapa dipakai di sini:** API bergantung pada Postgres dan Redis yang harus healthy; `postgres-data` menyimpan data Postgres dalam `deploy/compose/docker-compose.yml`.
- **Alternatif yang tidak dipilih:** memakai IP statis atau bind mount untuk data database; IP rapuh dan bind mount mengikat data ke filesystem host.
- **Cara membuktikan sendiri:** `docker run --rm --network campus-lms_default alpine getent hosts postgres redis`.
- **Pertanyaan interview terkait:** Mengapa `depends_on` tanpa `condition: service_healthy` tidak cukup untuk API yang perlu database siap?

### Resource limit dan OOM-kill

- **Apa:** cgroup membatasi resource container; kernel dapat menghentikan process yang melampaui limit memori.
- **Kenapa dipakai di sini:** setiap service core mempunyai `mem_limit` untuk menjaga laptop dengan RAM terbatas.
- **Alternatif yang tidak dipilih:** tanpa limit membiarkan satu container memakai memori host tanpa batas yang disepakati.
- **Cara membuktikan sendiri:** jalankan command pada [oom-kill-demo.txt](evidence/week-02/oom-kill-demo.txt).
- **Pertanyaan interview terkait:** Apa arti `OOMKilled=true` dan mengapa exit code 137 muncul?

## 7. Belum Terverifikasi

| Hal | Kenapa belum terverifikasi | Rencana verifikasi |
|---|---|---|
| Penghematan RAM setelah Docker Desktop ke Docker Engine | WSL tidak memberi pengukuran host yang membandingkan kedua runtime; tidak ada angka yang sah untuk diklaim | Ukur dari host Windows dengan metode yang sama sebelum/sesudah bila environment pembanding tersedia |
| Runtime `linux/arm64` | Bukti hanya menunjukkan image arm64 dibangun dan dipublikasikan; laptop ini amd64 dan tidak menjalankan image arm64 | Jalankan pada host arm64 atau emulator yang diizinkan, lalu simpan health check dan arsitektur runtime |
| Verifikasi Compose oleh agent | Sandbox agent tidak dapat mengakses `/var/run/docker.sock`; bukti Compose dilakukan manusia | Jalankan ulang di terminal manusia dan simpan output bila konfigurasi berubah |
| Test suite pada sesi penyusunan laporan | Pengulangan agent gagal karena sandbox melarang `httptest` membuka listener IPv6 lokal, bukan karena assertion test | Jalankan `make test` pada terminal manusia; hasil gagal sandbox ada di [report-test-rerun.txt](evidence/week-02/report-test-rerun.txt) |
| Web memanggil API lewat nama service | Profile core sekarang tidak memuat web; bukti DNS hanya mencakup Postgres dan Redis | Tambahkan dan uji service web saat scope web disetujui |
| Profile `storage` dan `obs` | Belum memiliki implementasi service | Implementasi pada minggu yang direncanakan dan verifikasi per-profile |

**Asumsi yang dipakai tapi belum dibuktikan:**

- **INFERENCE:** `-healthcheck` akan berperilaku sama di host arm64; bukti runtime arm64 dapat membantahnya.
- **INFERENCE:** Docker Engine mengurangi tekanan RAM host dibanding Docker Desktop pada laptop ini; metrik host sebelum/sesudah dapat membantahnya.

## 8. Masalah & Cara Diselesaikan

### Masalah: Migrasi Docker Desktop ke Docker Engine
- **Gejala:** **FACT (pernyataan manusia):** Docker Desktop diganti karena keterbatasan RAM laptop 8GB.
- **Hipotesis yang salah:** Menganggap WSL `free -h` dapat mengukur penghematan RAM Docker Desktop versus Docker Engine.
- **Akar masalah:** `free -h` mengamati memori distro WSL, bukan seluruh pemakaian memori host Windows oleh kedua runtime.
- **Solusi:** **Keputusan manusia:** memakai Docker Engine. Tidak ada angka penghematan yang diklaim.
- **Pencegahan:** ukur pada lapisan host yang sama untuk semua kandidat runtime.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: HEALTHCHECK HTTP pada distroless
- **Gejala:** runtime distroless tidak memiliki shell atau HTTP client untuk menjalankan probe `wget` biasa.
- **Hipotesis yang salah:** shell-form HEALTHCHECK dapat berjalan otomatis pada distroless.
- **Akar masalah:** command HEALTHCHECK tetap membutuhkan executable yang tersedia pada filesystem runtime.
- **Solusi:** solusi awal menyalin BusyBox; solusi akhir memakai flag `/api -healthcheck` yang melakukan GET `/healthz` dengan timeout. Runtime akhir terukur 14.6MB dan health-nya `healthy` ([bukti](evidence/week-02/api-probe-healthcheck.txt)).
- **Pencegahan:** uji HEALTHCHECK pada image runtime aktual, bukan hanya build stage. Target pembanding BusyBox tetap tersedia untuk membandingkan efeknya.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: biaya BusyBox harus dibuktikan, bukan diasumsikan
- **Gejala:** BusyBox awalnya menyelesaikan probe, tetapi menambah executable yang tidak dipakai API.
- **Hipotesis yang salah:** penambahan utility kecil pasti tidak bermakna terhadap ukuran image.
- **Akar masalah:** ukuran layer dan ukuran image harus diukur pada target yang sebanding.
- **Solusi:** ukuran historis menunjukkan 15.9MB menjadi 14.6MB (penurunan sekitar 1.3MB), sementara `docker history` menunjukkan COPY BusyBox 1.22MB; dua pengukuran ini konsisten setelah pembulatan tampilan image ([ukuran](evidence/week-02/healthcheck-size-comparison-historical.txt), [history](evidence/week-02/busybox-layer-history.txt)).
- **Pencegahan:** pakai pengukuran image dan layer secara bersamaan ketika mengklaim optimasi.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: `go.sum` tidak ada pada modul tanpa dependency
- **Gejala:** Dockerfile awal perlu menyalin metadata Go, tetapi `go.sum` belum ada karena dependency masih nol.
- **Hipotesis yang salah:** `COPY go.mod go.sum ./` akan selalu tersedia pada semua modul Go.
- **Akar masalah:** Go hanya membuat `go.sum` bila ada checksum dependency yang perlu direkam.
- **Solusi:** gunakan `COPY go.mod go.sum* ./`; build akhir membuktikan wildcard bekerja dan `go mod download` melaporkan tidak ada dependency ([build](evidence/week-02/api-probe-runtime-build.txt)).
- **Pencegahan:** perlakukan file metadata yang opsional sebagai opsional pada build context, tetapi jangan menyembunyikan kegagalan file source wajib.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: `make prune` dapat menghapus database
- **Gejala:** target lama memakai `docker system prune -af --volumes`, sehingga named volume Postgres ikut menjadi target hapus.
- **Hipotesis yang salah:** `prune` rutin aman tanpa membedakan image/container sampah dari volume data.
- **Akar masalah:** opsi `--volumes` memperluas scope ke volume, termasuk `campus-lms_postgres-data`.
- **Solusi:** ditemukan sebelum ada data hilang; `make prune` kini tanpa `--volumes`, sedangkan `make prune-hard` membutuhkan konfirmasi `DESTROY`.
- **Pencegahan:** command destruktif diberi nama eksplisit, scope dipisah, dan membutuhkan konfirmasi manusia.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: sandbox agent memblokir socket Docker
- **Gejala:** agent tidak dapat melakukan verifikasi Compose melalui `/var/run/docker.sock`; bukti Compose karena itu berasal dari eksekusi manusia.
- **Hipotesis yang salah:** keberadaan Docker CLI berarti daemon pasti dapat diakses oleh proses agent.
- **Akar masalah:** izin sandbox tidak mencakup socket daemon Docker.
- **Solusi:** manusia menjalankan verifikasi Compose dan menyimpan raw output; agent hanya membaca serta melaporkan batasnya.
- **Pencegahan:** cek akses daemon di awal sesi sebelum menjanjikan verifikasi container oleh agent.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: versi Air terbaru tidak kompatibel dengan Go yang tersedia
- **Gejala:** **FACT (pernyataan manusia):** Air terbaru memerlukan Go 1.25, sedangkan environment menggunakan Go 1.23.6.
- **Hipotesis yang salah:** versi tool terbaru selalu kompatibel dengan compiler proyek yang dipin.
- **Akar masalah:** requirement Go dari tool development tidak sama dengan versi Go yang tersedia di builder.
- **Solusi:** pilih `github.com/air-verse/air@v1.61.7`; log runtime mengonfirmasi Air v1.61.7 berjalan dan me-reload source ([hot reload](evidence/week-02/hot-reload-demo.txt)).
- **Pencegahan:** verifikasi minimum Go version dependency/tool sebelum memasangnya atau mem-pinnya.
- **Waktu terbuang:** NOT MEASURED.

### Pola berulang: alat ukur mengukur hal yang salah

Empat kali alat ukur menghasilkan sinyal yang tidak menjawab pertanyaan sebenarnya: grep `fmt.Println` mencocoki komentar, bukan hanya kode executable; `free -h` mengukur distro WSL yang salah untuk membandingkan Docker Desktop dengan Docker Engine; ukuran image tidak membuktikan isi atau efektivitas `.dockerignore` tanpa melihat transfer build context; dan `$(...)` dimakan Make sebelum mencapai shell bila tanda `$` tidak di-escape. **RECOMMENDATION:** setiap command verifikasi harus menyatakan objek ukur, layer eksekusi, dan output yang diharapkan sebelum hasilnya dipakai sebagai klaim.

## 9. Status Definition of Done

> Agent hanya MENGUSULKAN. Kolom "Dicentang manusia" diisi olehmu setelah melihat bukti.

| DoD dari roadmap | Usulan agent | Bukti | Dicentang manusia |
|---|---|---|---|
| `make up` menghasilkan seluruh stack sehat dan web memanggil API lewat service name | ⚠️ sebagian: API, Postgres, Redis healthy; web belum ada pada core | [Compose sehat](evidence/week-02/compose-ps-healthy.txt), [DNS](evidence/week-02/compose-dns.txt) | ✅ |
| Image API di bawah target dan berjalan non-root | ✅ kandidat terpenuhi secara lokal, menunggu verifikasi manusia | [size](evidence/week-02/image-size.txt), [UID](evidence/week-02/docker-inspect-user.txt) | ✅ |
| Build amd64/arm64 dan push GHCR | ⚠️ sebagian: kedua platform dibangun dan manifest GHCR ada; runtime arm64 belum diuji | [buildx](evidence/week-02/buildx-multiarch.txt), [GHCR](evidence/week-02/ghcr-push.txt) | ✅ |
| Container dibatasi 256MB dan OOM-kill dapat dijelaskan | ⚠️ sebagian: OOM-kill dibuktikan pada 32MiB; limit Compose tercatat 128MiB/512MiB/128MiB, bukan demonstrasi 256MiB | [OOM](evidence/week-02/oom-kill-demo.txt), [limit](evidence/week-02/compose-memlimit.txt) | ☐ |
| ADR menjelaskan pilihan distroless dan trade-off debugging | ✅ kandidat terpenuhi; Decision/Consequences ditulis manusia | [ADR-0005](../adr/0005-api-healthcheck-probe.md), [history](evidence/week-02/busybox-layer-history.txt) | ✅ |

## 10. Untuk Minggu Depan

- **Carry-over:** putuskan scope web terhadap roadmap, tambahkan service web bila tetap menjadi DoD Minggu 02, lalu buktikan web → API melalui DNS service.
- **Utang teknis yang sengaja diambil:** `storage` dan `obs` masih belum berisi service; kerjakan sesuai minggu yang direncanakan.
- **Persiapan yang perlu dilakukan manusia lebih dulu:** jalankan `make test` di terminal manusia dan uji runtime arm64 pada host/emulator yang diizinkan; ukur penggunaan RAM host jika ingin mengklaim dampak migrasi Docker.

## 11. Verifikasi Manusia

- [x] Saya sudah spot-check 3 file bukti secara acak dan isinya cocok dengan klaim
- [x] Skor quiz: 100% / 100% (minimal 70% untuk lanjut)
- [x] Explain-back 3 menit sudah direkam: `docs/progress/explain/week-02.<mp3|txt>`
- [x] Saya bisa menjelaskan setiap keputusan di section 4 tanpa membuka catatan

**Self-audit agent (dari `agent/evidence-protocol.md` §8):** Sudah dijalankan. Semua angka pengukuran di laporan memiliki link bukti; seluruh file bukti Week 02, sumber code quiz, git log, dan file yang diklaim telah dibaca; `make test` dicoba ulang tetapi tidak dapat dinyatakan lulus karena sandbox menolak socket `httptest` ([raw output](evidence/week-02/report-test-rerun.txt)); Belum Terverifikasi berisi gap nyata; agent tidak mencentang DoD; FACT/INFERENCE/RECOMMENDATION diberi label; dan 10 pertanyaan quiz diverifikasi terhadap file sumber yang ada. Karena test pengulangan tidak selesai di sandbox, self-audit ini **tidak sepenuhnya lulus** dan memerlukan verifikasi manusia.

**Ditandatangani:** FfFfFf  **Tanggal:** 19 Agustus 2026
