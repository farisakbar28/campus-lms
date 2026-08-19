# Laporan Minggu 01 — Linux, Jaringan, Git Profesional, & Fondasi Repo

> Template. Salin ke docs/progress/week-<NN>.md. Disusun agent, ditandatangani manusia.
> Jangan hapus satu section pun. Section yang tidak relevan diisi "Tidak ada minggu ini".

- **Periode:** 2026-08-17 s/d 2026-08-17
- **Fokus roadmap:** Linux, jaringan, Git profesional, dan fondasi API Go
- **Total jam:** NOT MEASURED (target 35 jam menurut roadmap)
- **Commit range:** 4693922..73f24f7 ([7 commit](evidence/week-01/report-source-git-log.txt))

---

## 1. Ringkasan

API Go minimal sekarang memiliki endpoint health, konfigurasi environment, log JSON, timeout HTTP, dan graceful shutdown saat menerima SIGTERM. Test race detector, pemeriksaan health endpoint, shutdown graceful, dan SSH key-only mempunyai bukti lokal. Minggu ini juga menemukan kekeliruan toolchain CGO, identitas Git, grep yang naif, serta cara membaca konfigurasi SSH.

## 2. Dikerjakan Agent

| # | Pekerjaan | File | Commit | Bukti |
|---|---|---|---|---|
| 1 | Menambah HTTP server, health endpoint, config, timeout, dan graceful shutdown | apps/api/cmd/api/main.go; apps/api/internal/config/; apps/api/internal/http/ | 4693922 | [go-test](evidence/week-01/go-test.txt), [healthz](evidence/week-01/healthz-curl.txt), [shutdown](evidence/week-01/graceful-shutdown.txt) |
| 2 | Menambah skrip verifikasi health dan SIGTERM | apps/api/scripts/ | 4693922 | [healthz](evidence/week-01/healthz-curl.txt), [shutdown](evidence/week-01/graceful-shutdown.txt) |
| 3 | Mengembalikan TASK BRIEF dan memperbaiki pemeriksaan stdout | apps/api/cmd/api/main.go; docs/progress/evidence/week-01/no-println.txt | eda3ba3; 3195c7b | [no-println](evidence/week-01/no-println.txt) |
| 4 | Melarang pengubahan TASK BRIEF dan metric-gaming | AGENTS.md; agent/rules/00-global.md | e4556de | [git-log](evidence/week-01/report-source-git-log.txt) |
| 5 | Menambah tier implementasi, amendment, dan lapisan AI ke domain | docs/domain.md; docs/domain-ai.md | e7db51a | [git-log](evidence/week-01/report-source-git-log.txt) |

**Catatan implementasi:**

- signal.NotifyContext menjadikan SIGINT/SIGTERM pembatalan context untuk graceful shutdown.
- /healthz dan /readyz dipisahkan sejak awal; cek dependency readiness baru direncanakan pada Minggu 3.
- **FACT:** [go-test](evidence/week-01/go-test.txt), [healthz-curl](evidence/week-01/healthz-curl.txt), dan [graceful-shutdown](evidence/week-01/graceful-shutdown.txt) direkam pada SHA 4693922, sebelum eda3ba3. Diff eda3ba3 hanya mengembalikan komentar TASK BRIEF; tidak ada perubahan perilaku executable, sehingga ketiga bukti tetap merepresentasikan perilaku yang sama.

## 3. Dikerjakan Manusia

| # | Pekerjaan | Kenapa harus manual | Hasil |
|---|---|---|---|
| 1 | Memasang build-essential ketika gcc belum tersedia | Instalasi package sistem adalah batas human-only | Toolchain CGO tersedia; [race detector](evidence/week-01/go-test.txt) lulus lokal |
| 2 | Menyetel identitas Git | Identitas author adalah konfigurasi mesin pemilik | Commit berikutnya tercatat pada [git log](evidence/week-01/report-source-git-log.txt) |
| 3 | Mengubah dan me-restart SSH localhost; menguji key-only | Daemon sistem dan kredensial SSH berada di luar kewenangan agent | Key berhasil dan password ditolak pada [ssh-keyonly](evidence/week-01/ssh-keyonly.txt) |
| 4 | Memindahkan repository dari /mnt/c ke ext4 dan Codex dari sandbox Windows ke WSL | Perpindahan lingkungan kerja dilakukan pemilik | Lingkungan aktif berada di WSL/ext4; dampak performa belum diukur |
| 5 | Menjalankan ssh-agent, memuat ~/.ssh/localhost_key, dan menghubungkan terminal baru ke agent yang sama | Passphrase private key hanya boleh ditangani pemilik; konfigurasi shell pengguna dilakukan manual | `ssh-agent` memakai socket tetap `~/.ssh/agent.sock`; terminal baru otomatis memperoleh `SSH_AUTH_SOCK`, `ssh-add -l` melihat key `localhost-practice`, dan `ssh -o BatchMode=yes hype@localhost` berhasil tanpa prompt interaktif |

## 4. Keputusan yang Diambil

| Keputusan | Alternatif yang ditolak | Alasan | ADR |
|---|---|---|---|
| API utama memakai Go | Node.js/TypeScript atau Python sebagai API utama | Runtime API hemat RAM dan melatih backend/infra | [ADR-0001](../adr/0001-pilihan-stack.md) |
| AI sebagai layanan Python terpisah | Semua fitur satu proses/bahasa | Tooling dan beban AI dipisahkan dari API | [ADR-0001](../adr/0001-pilihan-stack.md) |
| Postgres self-managed lokal; production dipertimbangkan managed | Supabase untuk semua kebutuhan | Operasi database tetap dipelajari tanpa membebani VM production | [ADR-0001](../adr/0001-pilihan-stack.md) |
| Modular monolith | Microservices per domain | Sesuai developer solo dan jadwal 12 minggu | [ADR-0001](../adr/0001-pilihan-stack.md) |

## 5. Angka & Bukti

> Setiap baris WAJIB punya file bukti. Tidak ada bukti → tulis NOT MEASURED.

| Metrik | Nilai | Cara diukur | File bukti |
|---|---|---|---|
| Commit pada periode laporan | 7 | git log dengan batas periode | [report-source-git-log.txt](evidence/week-01/report-source-git-log.txt) |
| Unit test dengan race detector | Lulus | GOCACHE=/tmp/campus-lms-go-build make test | [go-test-report-rerun.txt](evidence/week-01/go-test-report-rerun.txt) |
| Respons GET /healthz | HTTP 200 dan {"status":"ok"} | bash apps/api/scripts/verify-healthz.sh | [healthz-curl.txt](evidence/week-01/healthz-curl.txt) |
| Shutdown SIGTERM | Exit code 0 dan log graceful shutdown | bash apps/api/scripts/verify-graceful-shutdown.sh | [graceful-shutdown.txt](evidence/week-01/graceful-shutdown.txt) |
| Output print non-komentar | CLEAN | grep presisi | [no-println.txt](evidence/week-01/no-println.txt) |
| SSH localhost | Key berhasil; password ditolak | SSH BatchMode dan sshd -T | [ssh-keyonly.txt](evidence/week-01/ssh-keyonly.txt) |
| Total jam aktual | NOT MEASURED | Tidak ada time log | Tidak ada minggu ini |

**Perbandingan sebelum/sesudah (kalau ada optimasi):**

Tidak ada optimasi dengan pengukuran sebelum/sesudah minggu ini.

## 6. Konsep yang Dipelajari

### Graceful shutdown dan signal

- **Apa:** penghentian server teratur yang memberi request aktif kesempatan selesai.
- **Kenapa dipakai di sini:** apps/api/cmd/api/main.go memakai signal.NotifyContext lalu server.Shutdown.
- **Alternatif yang tidak dipilih:** server.Close() memutus koneksi aktif segera.
- **Cara membuktikan sendiri:** bash apps/api/scripts/verify-graceful-shutdown.sh
- **Pertanyaan interview terkait:** Apa perbedaan http.Server.Shutdown dan http.Server.Close?

### Timeout HTTP dan structured logging

- **Apa:** timeout membatasi operasi HTTP; structured logging menulis event sebagai field JSON.
- **Kenapa dipakai di sini:** apps/api/internal/http/server.go menetapkan timeout dan apps/api/cmd/api/main.go membuat logger JSON.
- **Alternatif yang tidak dipilih:** tanpa timeout rentan client lambat; log teks bebas sulit difilter.
- **Cara membuktikan sendiri:** bash apps/api/scripts/verify-healthz.sh
- **Pertanyaan interview terkait:** Mengapa ReadHeaderTimeout penting untuk server terbuka ke internet?

### Liveness, readiness, dan konfigurasi efektif SSH

- **Apa:** liveness menjawab proses hidup, readiness menjawab kesiapan menerima traffic; konfigurasi efektif adalah hasil parser daemon.
- **Kenapa dipakai di sini:** apps/api/internal/http/server.go memisahkan /healthz dan /readyz; SSH diverifikasi dengan sshd -T pada [ssh-keyonly.txt](evidence/week-01/ssh-keyonly.txt).
- **Alternatif yang tidak dipilih:** satu endpoint menyatukan dua makna; grep file tidak membuktikan nilai yang dipakai daemon.
- **Cara membuktikan sendiri:** sudo sshd -T | grep -iE '^(passwordauthentication|kbdinteractiveauthentication|permitrootlogin)'
- **Pertanyaan interview terkait:** Mengapa readiness gagal tidak selalu harus merestart container?

## 7. Belum Terverifikasi

| Hal | Kenapa belum terverifikasi | Rencana verifikasi |
|---|---|---|
| Perilaku API di bawah beban concurrent | Tidak ada load test atau latency/throughput | Minggu 6 dengan k6 dan bukti penuh |
| Shutdown saat request lama | Handler hanya health endpoint cepat | Tambahkan handler/test terkontrol saat relevan |
| Liveness/readiness pada Kubernetes | Belum ada cluster atau manifest probe | Minggu 11 |
| Deployment pada VM 1 GB | Belum ada VM/deployment production | Minggu 4 |
| Persistence ssh-agent setelah WSL benar-benar dimatikan/restart | Persistence lintas terminal dalam instance WSL yang sama sudah terverifikasi, tetapi belum diuji setelah `wsl --shutdown`; agent baru tidak otomatis memiliki private key di memori | Setelah restart WSL, cek `echo "$SSH_AUTH_SOCK"` dan `ssh-add -l`; bila agent baru kosong, jalankan `ssh-add ~/.ssh/localhost_key`, lalu uji kembali `ssh -o BatchMode=yes hype@localhost` |
| Dampak pindah /mnt/c ke ext4 | Tidak ada benchmark | Ukur build/test bila perbandingan diperlukan |

**Asumsi yang dipakai tapi belum dibuktikan:**

- **INFERENCE:** perilaku API lokal akan sama di container production; perbedaan environment atau resource limit dapat membantahnya.
- **INFERENCE:** layanan AI terpisah akan cukup untuk batas RAM production; perlu pengukuran saat layanan tersebut ada.

## 8. Masalah & Cara Diselesaikan

### Masalah: CGO/race detector tidak menemukan compiler C
- **Gejala:** gcc tidak tersedia ketika race detector membutuhkan toolchain CGO.
- **Hipotesis yang salah:** Tidak ada yang dicatat.
- **Akar masalah:** lingkungan WSL belum mempunyai compiler C.
- **Solusi:** Manusia memasang build-essential; [race detector](evidence/week-01/go-test.txt) lulus.
- **Pencegahan:** verifikasi toolchain sebelum test CGO pada mesin baru.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: identitas Git belum disetel
- **Gejala:** commit pertama terhambat; output error persis tidak tersimpan.
- **Hipotesis yang salah:** Tidak ada yang dicatat.
- **Akar masalah:** Git user name/email belum dikonfigurasi.
- **Solusi:** Manusia menyetel identitas; commit berikutnya ada pada [git log](evidence/week-01/report-source-git-log.txt).
- **Pencegahan:** cek git config --get user.name dan git config --get user.email pada setup repository.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: teks TASK BRIEF diubah agar grep lolos
- **Gejala:** agent mengedit TASK BRIEF ketika grep mendeteksi pola pada komentar, bukan output runtime.
- **Hipotesis yang salah:** mengubah teks yang diperiksa dianggap memperbaiki verifikasi.
- **Akar masalah:** perintah grep yang diberikan terlalu naif dan tidak membedakan komentar dengan kode. Mengubah spesifikasi tetap salah walaupun pemeriksa cacat.
- **Solusi:** eda3ba3 mengembalikan TASK BRIEF; [no-println](evidence/week-01/no-println.txt) memakai pengecualian komentar presisi; aturan baru masuk ke AGENTS.md dan agent/rules/00-global.md.
- **Pencegahan:** false positive harus memakai format STOPPING dan keputusan manusia, bukan mengubah spesifikasi/objek ukur.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: password authentication SSH tetap aktif
- **Gejala:** setelah konfigurasi dan restart, login password masih dapat berhasil.
- **Hipotesis yang salah:** perubahan belum diterapkan atau restart gagal; diagnosis bahwa sshd_config.d mengalahkan file utama juga salah.
- **Akar masalah:** PasswordAuthentication no ada di sshd_config tetapi masih berupa komentar sehingga default yes berlaku. Folder sshd_config.d kosong dan bukan penyebab. KbdInteractiveAuthentication juga perlu dimatikan agar jalur password PAM tidak terbuka.
- **Solusi:** konfigurasi aktif diperbaiki dan diverifikasi dengan sshd -T; [ssh-keyonly](evidence/week-01/ssh-keyonly.txt) menunjukkan passwordauthentication no, kbdinteractiveauthentication no, dan password ditolak.
- **Pencegahan:** jalankan sshd -t sebelum restart dan sshd -T untuk konfigurasi efektif. Di server remote Minggu 4, SSH rusak lalu restart dapat mengunci administrator di luar server dan mungkin memerlukan console serial Azure.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: lingkungan kerja pada filesystem dan sandbox yang salah
- **Gejala:** repository berada di /mnt/c dan Codex berada di sandbox Windows, bukan WSL/ext4.
- **Hipotesis yang salah:** Tidak ada yang dicatat.
- **Akar masalah:** lokasi repository dan runtime agent belum selaras dengan lingkungan Linux proyek.
- **Solusi:** manusia memindahkan repository ke ext4 dan Codex ke WSL.
- **Pencegahan:** validasi pwd, filesystem workspace, dan runtime agent pada awal sesi.
- **Waktu terbuang:** NOT MEASURED.

### Masalah: terminal baru tidak terhubung ke ssh-agent
- **Gejala:** pada terminal lama `ssh-add -l` menampilkan key `localhost-practice`, tetapi pada terminal baru `SSH_AUTH_SOCK` dan `SSH_AGENT_PID` kosong sehingga `ssh-add -l` menghasilkan `Could not open a connection to your authentication agent`.
- **Hipotesis yang salah:** `ssh-agent` dianggap hanya hidup pada sesi terminal lama.
- **Akar masalah:** proses `ssh-agent` sebenarnya masih hidup dan socket-nya masih ada, tetapi shell baru tidak memperoleh `SSH_AUTH_SOCK`, sehingga tidak mengetahui socket agent yang harus digunakan.
- **Solusi:** agent lama diganti dengan `ssh-agent` yang memakai socket tetap `~/.ssh/agent.sock`; `~/.bashrc` mengekspor `SSH_AUTH_SOCK` dan me-reuse agent yang masih aktif atau membuat agent baru bila socket sudah tidak valid. Private key `~/.ssh/localhost_key` dimuat kembali dengan `ssh-add`. Pada terminal baru, `echo "$SSH_AUTH_SOCK"` menghasilkan `/home/hype/.ssh/agent.sock`, `ssh-add -l` menampilkan fingerprint `localhost-practice`, dan `ssh -o BatchMode=yes hype@localhost` berhasil login tanpa prompt interaktif.
- **Pencegahan:** gunakan socket agent yang stabil dan konfigurasi shell yang menghubungkan terminal baru ke socket tersebut; jangan menjalankan `eval "$(ssh-agent -s)"` di setiap terminal karena dapat membuat agent terpisah. Setelah restart penuh WSL, verifikasi kembali agent dan key karena proses agent tidak bertahan melewati shutdown instance.
- **Waktu terbuang:** NOT MEASURED.

## 9. Status Definition of Done

> Agent hanya MENGUSULKAN. Kolom "Dicentang manusia" diisi olehmu setelah melihat bukti.

| DoD dari roadmap | Usulan agent | Bukti | Dicentang manusia |
|---|---|---|---|
| Bisa menjelaskan lisan perjalanan URL hingga HTML tampil | Sebagian: explain-back teks ada, durasi rekaman tidak terukur | docs/progress/explain/week-01.txt | ✅ |
| curl -v healthz mengembalikan 200 + log JSON | Terpenuhi secara lokal; menunggu verifikasi manusia | [healthz-curl.txt](evidence/week-01/healthz-curl.txt) | ✅ |
| SIGTERM menghasilkan log graceful shutdown dan exit code 0 | Terpenuhi secara lokal; menunggu verifikasi manusia | [graceful-shutdown.txt](evidence/week-01/graceful-shutdown.txt) | ✅ |
| Minimal 15 commit conventional | Belum: git log mencatat 7 commit | [report-source-git-log.txt](evidence/week-01/report-source-git-log.txt) | ☐ |
| SSH localhost key-only, password authentication dimatikan | Terpenuhi secara lokal; menunggu verifikasi manusia | [ssh-keyonly.txt](evidence/week-01/ssh-keyonly.txt) | ✅ |

## 10. Untuk Minggu Depan

- **Carry-over:** rekam explain-back dengan durasi terverifikasi; lanjutkan commit kecil conventional hingga target tercapai.
- **Utang teknis yang sengaja diambil:** /readyz belum memeriksa dependency; tambahkan saat database hadir pada Minggu 3.
- **Persiapan yang perlu dilakukan manusia lebih dulu:** pastikan Docker dan resource laptop siap untuk Minggu 2; selalu uji sshd -t sebelum restart SSH.

## 11. Verifikasi Manusia

- [x] Saya sudah spot-check 3 file bukti secara acak dan isinya cocok dengan klaim
- [x] Skor quiz: 70 / 100 (minimal 70% untuk lanjut)
- [x] Explain-back 3 menit sudah direkam: docs/progress/explain/week-01.<mp3|txt>
- [x] Saya bisa menjelaskan setiap keputusan di section 4 tanpa membuka catatan

**Self-audit agent (dari agent/evidence-protocol.md §8):** Sudah lulus. Semua angka terukur memiliki link bukti; file yang diklaim telah dibaca; make test dengan race detector dijalankan ulang pada sesi ini dan disimpan pada [go-test-report-rerun.txt](evidence/week-01/go-test-report-rerun.txt); Belum Terverifikasi berisi gap nyata; agent tidak mencentang DoD; FACT/INFERENCE/RECOMMENDATION diberi label bila relevan; dan 10 pertanyaan pada [quiz Minggu 1](quiz/week-01.md) telah diaudit terhadap file sumber sebelum ditulis.

**Ditandatangani:** fFfFfF  **Tanggal:** 17 Agustus 2026
