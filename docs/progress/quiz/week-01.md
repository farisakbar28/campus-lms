# Quiz Verifikasi — Minggu 01

> Template. Salin ke docs/progress/quiz/week-<NN>.md.
>
> **Aturan untuk agent:** 8–12 soal. Setiap soal WAJIB bisa dijawab dari kode
> yang benar-benar ada di repo ini — buka filenya dulu sebelum menulis soal.
> Dilarang membuat soal tentang kode yang kamu tulis tapi tidak kamu jalankan.
>
> **Aturan untuk manusia:** jawab tanpa membuka kode. Tanpa mencari. Tanpa
> bertanya ke AI. Skor < 70% berarti minggu ini diulang, bukan dilanjutkan.
> Ini bukan hukuman — ini alat ukur supaya kamu tidak menemukan lubangnya
> saat interview.

**Komposisi soal:** 40% penalaran (1, 3, 8, 10) · 30% failure mode (2, 4, 9) · 20% orientasi (5, 6) · 10% hafalan (7)

---

## Soal

**1. [Penalaran]** Mengapa main membuat context dengan signal.NotifyContext untuk SIGINT dan SIGTERM, lalu memanggil server.Shutdown, alih-alih membiarkan signal menghentikan proses segera?

<br>

**2. [Failure mode]** Apa yang terjadi bila APP_SHUTDOWN_TIMEOUT tidak tersedia, tidak dapat diparse sebagai duration, atau bernilai nol/negatif?

<br>

**3. [Penalaran]** Mengapa /healthz dan /readyz tetap terpisah walaupun sekarang mengembalikan status yang sama?

<br>

**4. [Failure mode]** Apa yang dilakukan main bila server.Shutdown(shutdownCtx) mengembalikan error, misalnya karena timeout shutdown habis?


<br>

**5. [Orientasi]** Di file dan fungsi mana route GET /healthz dan GET /readyz didaftarkan, dan fungsi bersama apa yang menulis HTTP 200 serta body JSON?

<br>

**6. [Orientasi]** Di mana logger JSON utama dibuat, dan tiga field konstan apa yang otomatis ditambahkan ke setiap log?

<br>

**7. [Hafalan]** Sebutkan empat timeout HTTP yang diatur server dan nilainya.

<br>

**8. [Penalaran]** Mengapa host key SSH di /etc/ssh/ssh_host_* tidak boleh disamakan dengan user key ~/.ssh/localhost_key? Jelaskan siapa yang dibuktikan oleh masing-masing key.


<br>

**9. [Failure mode]** Mengapa grep isi sshd_config tidak cukup untuk membuktikan password login telah mati, dan apa risiko me-restart konfigurasi SSH yang belum tervalidasi pada server remote?

<br>

**10. [Penalaran]** Mengapa target build API menetapkan CGO_ENABLED=0, tetapi target test tetap menjalankan go test -race dan membutuhkan compiler C pada lingkungan ini?

---

## Jawabanmu

| No | Jawaban | Yakin? (1-5) |
|---|---|---|
| 1 | alasan kenapa main membuat context untuk SIGINT dan SIGTERM adalah karena untuk mengirimkan sinyal kapan harus dishutdown, dan .shutdown menggunakan sinyal untuk memberhentikan request yang akan datang dan membiarkan request yang sudah masuk selesai (graceful).
 | 5 |
| 2 | jika timeout tidak tersedia, aplikasi akan menggunakan nilai default, ketika nilai default tidak dapat diparse sebagai duration atau bernilai nol/negatif, konfigurasi seharusnya valid dan aplikasi harusnya gagal startup daripada dengan timeout shutdown yang tidak aman. | 4.5 |
| 3 | karena keduanya memiliki fungsi dan kegunaan yang berbeda. /healthz digunakan untuk mendapatkan status apakah backend/API tersedia, sementara /readyz digunakan untuk mendapatkan status apakah backend/API sudah aman/ready untuk menerima traffic atau tidak. kenapa itu dipisah, karena jika hanya menjadikannya sebagai satu kesatuan (misal /health saja), itu bisa memiliki konsekuensi kehilangan informasi penting ketika ada gangguan. | 5 |
| 4 | kalau server.Shutdown(shutdownCtx) gagal, misalnya karena cfg.ShutdownTimeout habis sebelum semua request selesai, program akan mencatat log "graceful shutdown failed" beserta error-nya lalu menjalankan os.Exit(1), jadi proses dianggap berhenti tidak normal karena graceful shutdown tidak berhasil diselesaikan. | 4.5 |
| 5 | keduanya berjalan di file server.go dan fungsi yang mendaftarkan GET /healthz dan GET /readyz adalah fungsi "func NewServer". dan fungsi bersama yang menulis HTTP 200 serta body JSON adalah fungsi "func writeStatus". | 4 |
| 6 | logger JSON utama dibuat pada file config.go terutama fungsi "func Load" yang menampilkan 3 field konstan ke setiap log, yaitu field port yang digunakan, log level, dan shutdown timeout. | 4 |
| 7 | ReadHeaderTimeout, ReadTimeout, WriteTimeout, IdleTimeout | 5 |
| 8 | host key di /etc/ssh/ssh_host_* dipakai server untuk membuktikan “saya benar-benar server/host yang kamu tuju”, sedangkan user key seperti ~/.ssh/localhost_key dipakai client untuk membuktikan “saya benar-benar user yang berhak login.” Jadi host key mengautentikasi identitas server ke client, sementara user key mengautentikasi identitas user ke server. | 3 |
| 9 | grep isi sshd konfig tidak cukup karena hanya untuk digunakan untuk membuktikan teks apa yang ditulis pada file, bukan akhir yang benar benar dipakai sshd, karena masih bisa ada default, include, atau hal semacamnya. untuk itu sshd -T digunakan untuk melihat konfigurasi efektif. pada server remote, me-restar ssh sebelum konfigurasi divalidasi sshd dapat membuat gagal start atau menolak login, yang dapat membuat bisa terkunci di server yang pemulihannya harus diakses melalui console/provider. | 3.5 |
| 10 | karena build production ingin menghasilkan binary Go yang lebih portable dan tidak bergantung pada library C, sehingga memakai CGO_ENABLED=0; sedangkan go test -race menggunakan Go race detector yang pada environment ini memerlukan CGO dan compiler C seperti gcc, jadi target test sengaja tidak mematikan CGO agar pengecekan race condition tetap bisa dijalankan. | 3 |

---

## Kunci Jawaban

> Jangan dibuka sebelum menjawab semua.

<details>
<summary>Buka kunci jawaban</summary>

**1.** Signal dijadikan cancellation context agar main masuk ke jalur graceful shutdown. Jalur itu membuat context bertimeout dan memanggil Shutdown, bukan membiarkan default signal menghentikan proses segera.
📁 Konfirmasi di: apps/api/cmd/api/main.go:81-105

**2.** config.Load mengembalikan error. main mencatat error melalui bootstrap logger dan keluar dengan status 1 sebelum server dibuat atau ListenAndServe berjalan.
📁 Konfirmasi di: apps/api/internal/config/config.go:49-59; apps/api/cmd/api/main.go:71-76

**3.** Healthz adalah liveness, sedangkan readyz disiapkan untuk dependency check. Pada Minggu 1 dependency check belum ada, sehingga proses hidup dianggap siap; pemisahan route menghindari perubahan kontrak probe nanti.
📁 Konfirmasi di: apps/api/internal/http/server.go:24-25; apps/api/internal/http/server.go:37-49

**4.** main mencatat "graceful shutdown failed" beserta error, lalu memanggil os.Exit(1). Shutdown yang gagal tidak ditandai sebagai exit normal.
📁 Konfirmasi di: apps/api/cmd/api/main.go:99-105

**5.** Keduanya didaftarkan di NewServer. Handler healthz dan readyz sama-sama memanggil writeStatus; fungsi ini mengatur Content-Type, menulis HTTP 200, lalu meng-encode status ok.
📁 Konfirmasi di: apps/api/internal/http/server.go:22-25; apps/api/internal/http/server.go:37-58

**6.** Logger utama dibuat oleh newLogger menggunakan slog.NewJSONHandler. Field yang ditambahkan adalah service=campus-api, env dari konfigurasi, dan version dari buildVersion().
📁 Konfirmasi di: apps/api/cmd/api/main.go:109-115

**7.** ReadHeaderTimeout 5 detik, ReadTimeout 10 detik, WriteTimeout 15 detik, dan IdleTimeout 60 detik.
📁 Konfirmasi di: apps/api/internal/http/server.go:11-16; apps/api/internal/http/server.go:27-34

**8.** Host key milik sshd dan dipakai client untuk memeriksa identitas host/server; konfigurasi HostKey menunjuk key di /etc/ssh. User key milik user/client dan dipakai server untuk memeriksa identitas user yang login. Evidence membuktikan key authentication localhost berhasil dan password ditolak.
📁 Konfirmasi di: /etc/ssh/sshd_config:40-42; docs/progress/evidence/week-01/ssh-keyonly.txt:2-16

**9.** Isi file dapat berupa komentar, default, Include, atau override; nilai yang dipakai daemon adalah hasil parsing efektif dari sshd -T. Bukti menunjukkan passwordauthentication no dan kbdinteractiveauthentication no. Restart SSH yang salah pada remote server dapat mengunci administrator di luar server, sehingga sshd -t dijalankan sebelum restart.
📁 Konfirmasi di: /etc/ssh/sshd_config:7-24; docs/progress/evidence/week-01/ssh-keyonly.txt:3,13-16; docs/progress/week-01.md:138-143

**10.** Build target sengaja membuat binary statis dengan CGO_ENABLED=0. Test target menjalankan race detector melalui -race; bukti Minggu 1 mencatat gcc sebelumnya tidak tersedia lalu test lulus setelah build-essential dipasang. Kebutuhan toolchain test tidak sama dengan runtime binary yang dibangun.
📁 Konfirmasi di: Makefile:43-49; apps/api/Dockerfile:14-20; docs/progress/week-01.md:114-120; docs/progress/evidence/week-01/go-test.txt:2-13

</details>

---

## Hasil

- **Skor:** 70 / 100 = 70%
- **Lulus (≥70%)?** ya
- **Soal yang salah:** 2, 6, 7
- **Konsep yang perlu diulang:** perilaku config.Load() terhadap APP_SHUTDOWN_TIMEOUT yang missing/invalid/non-positive; lokasi pembuatan logger JSON dan field konstan service, env, version; nilai empat HTTP timeout (ReadHeaderTimeout=5s, ReadTimeout=10s, WriteTimeout=15s, IdleTimeout=60s).
- **Rencana perbaikan:** baca ulang apps/api/internal/config/config.go, apps/api/cmd/api/main.go, dan apps/api/internal/http/server.go, lalu jelaskan kembali soal 2, 6, dan 7 tanpa melihat catatan sampai bisa menyebut perilaku konfigurasi, field logger, dan keempat nilai timeout dengan tepat.

