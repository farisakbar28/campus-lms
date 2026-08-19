# Quiz Verifikasi — Minggu 02

> Template. Salin ke `docs/progress/quiz/week-02.md`.
>
> **Aturan untuk agent:** 8–12 soal. Setiap soal WAJIB bisa dijawab dari kode
> yang benar-benar ada di repo ini — buka filenya dulu sebelum menulis soal.
> Dilarang membuat soal tentang kode yang kamu tulis tapi tidak kamu jalankan.
>
> **Aturan untuk manusia:** jawab tanpa membuka kode. Tanpa mencari. Tanpa
> bertanya ke AI. Skor < 70% berarti minggu ini diulang, bukan dilanjutkan.
> Ini bukan hukuman — ini alat ukur supaya kamu tidak menemukan lubangnya
> saat interview.
>
> Soal disusun dari file yang telah dibaca pada 2026-08-18.

**Komposisi soal:** 40% penalaran (1, 3, 6, 9) · 30% failure mode (2, 4, 7) · 20% orientasi (5, 8) · 10% hafalan (10)

---

## Soal

**1. [Penalaran]** Mengapa `apps/api/Dockerfile` menyalin `go.mod` dan `go.sum*` lalu menjalankan `go mod download` sebelum `COPY . ./`?
Jawaban: Supaya dependency download berada pada layer yang masih dapat dipakai ulang ketika hanya source application berubah. `COPY . ./` setelahnya baru membatalkan cache saat source berubah

<br>

**2. [Failure mode]** Apa yang terjadi bila Dockerfile memakai `COPY go.mod go.sum ./` ketika modul belum mempunyai dependency dan karena itu belum memiliki `go.sum`?
Jawaban: Build gagal karena `go.sum` diminta sebagai source `COPY` tetapi file itu tidak ada. Suffix wildcard pada `go.sum*` membuat file opsional tersebut tidak wajib ada

<br>

**3. [Penalaran]** Mengapa runtime distroless final memilih `HEALTHCHECK CMD ["/api", "-healthcheck"]` dibanding shell-form command atau `wget` biasa?
Jawaban: Distroless tidak membawa shell atau HTTP client. Binary API sudah ada di runtime dan memiliki mode probe sendiri, sehingga Docker dapat menjalankan command exec-form tanpa menambah BusyBox.


<br>

**4. [Failure mode]** Dalam `healthcheck.Probe`, apa hasilnya jika `/healthz` mengembalikan HTTP 503 atau request melebihi timeout?
Jawaban: Probe mengembalikan error untuk status selain 200. Request yang melewati timeout juga mengembalikan error dari HTTP client, sehingga mode healthcheck keluar dengan status gagal.

<br>

**5. [Orientasi]** Sebutkan file dan function yang membuat request HTTP probe lokal, serta endpoint yang dipanggil oleh `main` ketika flag `-healthcheck` aktif.
Jawaban: `healthcheck.Probe` pada `apps/api/internal/healthcheck/probe.go` membuat request. Saat flag aktif, `main` membangun `http://127.0.0.1:<port>/healthz` lalu memanggil function tersebut.

<br>

**6. [Penalaran]** Mengapa `postgres-data` adalah named volume, sedangkan source API pada override development adalah bind mount?
Jawaban: Named volume menjaga data Postgres terpisah dari container agar tetap ada setelah container dihentikan. Bind mount source membuat Air melihat perubahan file host dan dapat rebuild saat development.

<br>

**7. [Failure mode]** Apa risiko `make prune` lama yang menjalankan `docker system prune -af --volumes`, dan bagaimana target Makefile sekarang membatasi risiko itu?
Jawaban: Opsi `--volumes` dapat menghapus volume database. `prune` sekarang tidak memakai opsi itu; tindakan yang juga menghapus volume dipindah ke `prune-hard` dan meminta konfirmasi `DESTROY`.

<br>

**8. [Orientasi]** Di mana limit memori API, Postgres, dan Redis ditetapkan, dan berapa nilai yang ditulis untuk masing-masing service?
Jawaban: Semua berada di `deploy/compose/docker-compose.yml`: API `128m`, Postgres `512m`, Redis `128m`.

<br>

**9. [Penalaran]** Mengapa `depends_on` API memakai `condition: service_healthy` untuk Postgres dan Redis, bukan sekadar menunggu container dibuat?
Jawaban: Container yang baru dibuat belum tentu siap menerima koneksi. Condition tersebut menahan API sampai healthcheck Postgres dan Redis menyatakan sehat.

<br>

**10. [Hafalan]** Versi Air apa yang dipin pada override development, dan command apa yang dipakai Air untuk membangun binary sementara?
Jawaban: Air dipin pada `v1.61.7`; command build-nya adalah `go build -o /tmp/air/api ./cmd/api`.

---

## Jawabanmu

| No | Jawaban | Yakin? (1-5) |
|---|---|---|
| 1 | Supaya dependency download berada pada layer yang masih dapat dipakai ulang ketika hanya source application berubah. `COPY . ./` setelahnya baru membatalkan cache saat source berubah
 | 5 |
| 2 | Build gagal karena `go.sum` diminta sebagai source `COPY` tetapi file itu tidak ada. Suffix wildcard pada `go.sum*` membuat file opsional tersebut tidak wajib ada | 5 |
| 3 | Distroless tidak membawa shell atau HTTP client. Binary API sudah ada di runtime dan memiliki mode probe sendiri, sehingga Docker dapat menjalankan command exec-form tanpa menambah BusyBox. | 5 |
| 4 | Probe mengembalikan error untuk status selain 200. Request yang melewati timeout juga mengembalikan error dari HTTP client, sehingga mode healthcheck keluar dengan status gagal. | |
| 5 | `healthcheck.Probe` pada `apps/api/internal/healthcheck/probe.go` membuat request. Saat flag aktif, `main` membangun `http://127.0.0.1:<port>/healthz` lalu memanggil function tersebut. | |
| 6 | Named volume menjaga data Postgres terpisah dari container agar tetap ada setelah container dihentikan. Bind mount source membuat Air melihat perubahan file host dan dapat rebuild saat development. | 5 |
| 7 | Opsi `--volumes` dapat menghapus volume database. `prune` sekarang tidak memakai opsi itu; tindakan yang juga menghapus volume dipindah ke `prune-hard` dan meminta konfirmasi `DESTROY`. | |
| 8 | Jawaban: Semua berada di `deploy/compose/docker-compose.yml`: API `128m`, Postgres `512m`, Redis `128m`. | 5 |
| 9 | Container yang baru dibuat belum tentu siap menerima koneksi. Condition tersebut menahan API sampai healthcheck Postgres dan Redis menyatakan sehat. | 5 |
| 10 | Air dipin pada `v1.61.7`; command build-nya adalah `go build -o /tmp/air/api ./cmd/api`. | 5 |

---

## Kunci Jawaban

> Jangan dibuka sebelum menjawab semua.

<details>
<summary>Buka kunci jawaban</summary>

**1.** Supaya dependency download berada pada layer yang masih dapat dipakai ulang ketika hanya source application berubah. `COPY . ./` setelahnya baru membatalkan cache saat source berubah.
📁 Konfirmasi di: `apps/api/Dockerfile`

**2.** Build gagal karena `go.sum` diminta sebagai source `COPY` tetapi file itu tidak ada. Suffix wildcard pada `go.sum*` membuat file opsional tersebut tidak wajib ada.
📁 Konfirmasi di: `apps/api/Dockerfile`

**3.** Distroless tidak membawa shell atau HTTP client. Binary API sudah ada di runtime dan memiliki mode probe sendiri, sehingga Docker dapat menjalankan command exec-form tanpa menambah BusyBox.
📁 Konfirmasi di: `apps/api/Dockerfile`; `apps/api/cmd/api/main.go`

**4.** Probe mengembalikan error untuk status selain 200. Request yang melewati timeout juga mengembalikan error dari HTTP client, sehingga mode healthcheck keluar dengan status gagal.
📁 Konfirmasi di: `apps/api/internal/healthcheck/probe.go`; `apps/api/cmd/api/main.go`

**5.** `healthcheck.Probe` pada `apps/api/internal/healthcheck/probe.go` membuat request. Saat flag aktif, `main` membangun `http://127.0.0.1:<port>/healthz` lalu memanggil function tersebut.
📁 Konfirmasi di: `apps/api/cmd/api/main.go`; `apps/api/internal/healthcheck/probe.go`

**6.** Named volume menjaga data Postgres terpisah dari container agar tetap ada setelah container dihentikan. Bind mount source membuat Air melihat perubahan file host dan dapat rebuild saat development.
📁 Konfirmasi di: `deploy/compose/docker-compose.yml`; `deploy/compose/docker-compose.override.yml`

**7.** Opsi `--volumes` dapat menghapus volume database. `prune` sekarang tidak memakai opsi itu; tindakan yang juga menghapus volume dipindah ke `prune-hard` dan meminta konfirmasi `DESTROY`.
📁 Konfirmasi di: `Makefile`

**8.** Semua berada di `deploy/compose/docker-compose.yml`: API `128m`, Postgres `512m`, Redis `128m`.
📁 Konfirmasi di: `deploy/compose/docker-compose.yml`

**9.** Container yang baru dibuat belum tentu siap menerima koneksi. Condition tersebut menahan API sampai healthcheck Postgres dan Redis menyatakan sehat.
📁 Konfirmasi di: `deploy/compose/docker-compose.yml`

**10.** Air dipin pada `v1.61.7`; command build-nya adalah `go build -o /tmp/air/api ./cmd/api`.
📁 Konfirmasi di: `deploy/compose/docker-compose.override.yml`

</details>

---

## Hasil

- **Skor:** 10 / 10 = 100%
- **Lulus (≥70%)?** ya
- **Soal yang salah:** tidak ada
- **Konsep yang perlu diulang:** tidak ada
- **Rencana perbaikan:** tidak ada
