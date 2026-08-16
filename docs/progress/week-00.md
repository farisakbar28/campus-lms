# Laporan Minggu 00 — Persiapan Akun & Agent

> Minggu 0 berfokus pada pekerjaan persiapan yang sebagian besar hanya dapat dilakukan manusia: pengamanan akun cloud, klaim benefit mahasiswa, pengelolaan kredensial, persiapan environment lokal, dan validasi awal AI agent.
>
> Agent hanya digunakan untuk membantu membaca aturan repo, memverifikasi workflow, dan menguji apakah instruksi di `AGENTS.md` dipatuhi.

* **Periode:** 16 Agustus 2026 s/d 16 Agustus 2026
* **Fokus roadmap:** Persiapan akun cloud, klaim benefit mahasiswa, kredensial, environment lokal, repository, dan setup agent
* **Total jam:** 3 jam (estimasi roadmap ±4 jam)
* **Commit range:** `initial commit` (1)

---

## 1. Ringkasan

Minggu 0 digunakan untuk menyiapkan fondasi operasional `campus-lms` sebelum implementasi aplikasi dimulai. Subscription Azure for Students telah diverifikasi, spending limit dipastikan aktif, budget guard dibuat, kuota VM family di Southeast Asia diperiksa, serta konvensi resource Azure didokumentasikan. Environment lokal juga telah disiapkan melalui SSH key khusus Azure, Azure CLI, kredensial empat provider LLM, `.env`, dan konfigurasi WSL2 yang disesuaikan dengan laptop RAM 8 GB.

Repository publik dan aturan kerja AI agent juga telah disiapkan dan dibaca oleh manusia. Acceptance test agent sempat gagal karena GNU Make belum tersedia di WSL, kemudian environment diperbaiki oleh manusia dan pengujian diulang hingga agent berhasil menjalankan `make todo` serta merangkum pekerjaan Minggu 1 tanpa menambahkan task di luar repo. Satu guardrail yang masih belum terverifikasi adalah pengiriman email budget alert Azure, sehingga Minggu 0 belum boleh dianggap selesai sepenuhnya sebelum bukti tersebut tersedia.

---

## 2. Dikerjakan Agent

| # | Pekerjaan                                                                                                                                    | File                                          | Commit | Bukti                           |
| - | -------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------- | ------ | ------------------------------- |
| 1 | Membaca `AGENTS.md` dan `Makefile`, kemudian mencoba menjalankan `make todo` untuk acceptance test awal                                      | `AGENTS.md`, `Makefile`                       | —      | Belum disimpan sebagai evidence |
| 2 | Melaporkan kegagalan eksekusi awal karena `make` belum tersedia di environment WSL                                                           | —                                             | —      | Belum disimpan sebagai evidence |
| 3 | Mengulang acceptance test setelah environment diperbaiki manusia, menjalankan `make todo`, dan merangkum pekerjaan Minggu 1 berdasarkan repo | `AGENTS.md`, `Makefile` dan file TODO terkait | —      | Belum disimpan sebagai evidence |

**Catatan implementasi:**

* Agent tidak melakukan perubahan source code pada acceptance test Minggu 0.
* Pada pengujian pertama, agent tidak memasang GNU Make sendiri karena instalasi system package termasuk tindakan yang memerlukan persetujuan/manusia menurut `agent/policy.md`.
* Setelah GNU Make tersedia, acceptance test berhasil dijalankan melalui WSL.
* Acceptance test ini hanya membuktikan perilaku agent pada satu skenario awal; kepatuhan pada task engineering nyata masih harus terus diverifikasi pada Minggu 1 dan seterusnya.

---

## 3. Dikerjakan Manusia

| #  | Pekerjaan                                                                             | Kenapa harus manual                                              | Hasil                                      |
| -- | ------------------------------------------------------------------------------------- | ---------------------------------------------------------------- | ------------------------------------------ |
| 1  | Verifikasi subscription **Azure for Students**                                        | Akun dan verifikasi identitas pribadi                            | ✅ Subscription aktif                       |
| 2  | Memastikan **Spending limit: ON**                                                     | Pengaturan billing Azure Portal                                  | ✅ Aktif                                    |
| 3  | Membuat budget `campus-lms-guard` sebesar **$10/bulan** dengan alert 50% / 80% / 100% | Pengaturan Cost Management Azure                                 | ✅ Dibuat                                   |
| 4  | Menguji pengiriman email budget alert                                                 | Membutuhkan perubahan threshold dan akses email pribadi          | ⚠️ Belum terverifikasi                     |
| 5  | Memeriksa kuota `Standard BS Family vCPUs` di **Southeast Asia**                      | Azure Portal                                                     | ✅ Kuota family tersedia                    |
| 6  | Mendokumentasikan region, resource group, dan mandatory tags Azure                    | Keputusan konvensi cloud milik manusia                           | ✅ `docs/adr/0002c-azure-conventions.md`    |
| 7  | Klaim domain `.me` melalui Namecheap Student benefit                                  | Akun pribadi / Student Pack                                      | ✅ Berhasil                                 |
| 8  | Aktivasi GitHub Education / Student Developer Pack                                    | Akun dan verifikasi mahasiswa                                    | ✅ Aktif                                    |
| 9  | Klaim JetBrains Student Subscription                                                  | Akun pribadi                                                     | ⚠️ Tidak dikerjakan pada Minggu 0          |
| 10 | Klaim New Relic student benefit                                                       | Akun pribadi / GitHub Education                                  | ✅ Berhasil                                 |
| 11 | Membuat SSH key ED25519 khusus Azure                                                  | Private key hanya boleh dikelola manusia                         | ✅ `~/.ssh/campus_lms_azure`                |
| 12 | Memastikan permission private key `600`                                               | Keamanan credential lokal                                        | ✅ Terverifikasi                            |
| 13 | Menginstal Azure CLI                                                                  | System package / perubahan environment lokal                     | ✅ Terpasang                                |
| 14 | Login Azure CLI dengan device code                                                    | Kredensial akun Azure                                            | ✅ Berhasil                                 |
| 15 | Memverifikasi subscription Azure dari CLI                                             | Kredensial akun Azure                                            | ✅ `Azure for Students`, `Enabled`, default |
| 16 | Membuat API key Google AI Studio / Gemini                                             | Credential pribadi                                               | ✅ Berhasil                                 |
| 17 | Membuat API key Groq                                                                  | Credential pribadi                                               | ✅ Berhasil                                 |
| 18 | Membuat API key Cerebras                                                              | Credential pribadi                                               | ✅ Berhasil                                 |
| 19 | Membuat API key OpenRouter                                                            | Credential pribadi                                               | ✅ Berhasil                                 |
| 20 | Menyalin `.env.example` menjadi `.env` dan mengisi credential yang sudah tersedia     | `.env` dan secret adalah human-only                              | ✅ Selesai                                  |
| 21 | Membuat random `JWT_SECRET` lokal                                                     | Secret adalah human-only                                         | ✅ Selesai                                  |
| 22 | Membuat `.wslconfig` untuk laptop RAM 8 GB                                            | Konfigurasi OS/WSL berada di luar scope agent                    | ✅ `memory=5GB`, `processors=8`, `swap=8GB` |
| 23 | Restart dan memverifikasi konfigurasi WSL2                                            | Perubahan environment host                                       | ✅ Berhasil                                 |
| 24 | Menyiapkan repo `campus-lms` di laptop                                                | Environment dan akun GitHub manusia                              | ✅ Selesai                                  |
| 25 | `git init`, initial commit, dan push repo publik                                      | Akun GitHub dan keputusan publikasi                              | ✅ Selesai                                  |
| 26 | Mengganti placeholder `CHANGE_ME` dengan username GitHub                              | Identitas repository                                             | ✅ Selesai                                  |
| 27 | Membaca `AGENTS.md` sampai selesai                                                    | Pemahaman tidak dapat didelegasikan                              | ✅ Selesai                                  |
| 28 | Membaca `agent/policy.md` sampai selesai                                              | Pemahaman hard stop dan permission                               | ✅ Selesai                                  |
| 29 | Membaca `agent/README.md`                                                             | Pemahaman workflow agent                                         | ✅ Selesai                                  |
| 30 | Menginstal GNU Make setelah acceptance test menemukan `make` belum tersedia           | System package hanya boleh dipasang manusia / dengan persetujuan | ✅ Selesai                                  |
| 31 | Mengulang acceptance test agent                                                       | Verifikasi manusia terhadap perilaku agent                       | ✅ Lulus untuk skenario Minggu 0            |

---

## 4. Keputusan yang Diambil

| Keputusan                                                                         | Alternatif yang ditolak                                         | Alasan                                                                                                                       | ADR                                            |
| --------------------------------------------------------------------------------- | --------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| Menggunakan region Azure **Southeast Asia**                                       | Australia East, East Asia, atau region lain                     | Dipilih sebagai region utama project dan digunakan secara konsisten untuk resource production                                | `ADR-0002c`                                    |
| Menggunakan resource group `rg-campuslms-prod`                                    | Penamaan resource group ad-hoc                                  | Memberikan konvensi resource yang konsisten dan mudah diaudit                                                                | `ADR-0002c`                                    |
| Mewajibkan tags `project`, `env`, dan `owner` pada resource Azure                 | Resource tanpa tagging standar                                  | Mempermudah identifikasi, audit, dan cost tracking                                                                           | `ADR-0002c`                                    |
| Tidak mengunci arsitektur production pada SKU `Standard_B1s` sebelum provisioning | Menganggap B1s pasti tersedia sejak Minggu 0                    | Kuota VM family sudah tersedia, tetapi SKU final dan kapasitas aktual baru dapat diverifikasi ketika VM dibuat pada Minggu 4 | Belum dituangkan sebagai Decision ADR terpisah |
| Menjaga spending limit Azure for Students tetap aktif                             | Menghapus spending limit untuk memperoleh fleksibilitas billing | Project dirancang agar tidak membuka jalur tagihan cloud tak terbatas                                                        | `ADR-0002c` / konvensi Azure                   |

---

## 5. Angka & Bukti

> Setiap nilai pada bagian ini berasal dari verifikasi aktual Minggu 0 dan harus dapat ditelusuri ke raw output atau screenshot di `docs/progress/evidence/week-00/`. Informasi sensitif seperti API key, private key, Subscription ID, Tenant ID, dan isi `.env` tidak disimpan sebagai evidence.

| Metrik | Nilai | Cara diukur | File bukti |
| ------------------------------------------------- | ------------------------------------------ | ------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------- |
| Subscription Azure aktif | `Azure for Students` — `Enabled` — default | `az account show --query "{Name:name,State:state,IsDefault:isDefault,Environment:environmentName}" --output table` | [azure-account-show.txt](evidence/week-00/azure-account-show.txt) |
| Spending limit Azure for Students | ON | Azure Portal → Subscription | [azure-spending-limit.png](evidence/week-00/azure-spending-limit.png) Tidak Ada screenshoot |
| Budget guard | $10/bulan | Azure Portal → Cost Management → Budgets | [azure-budget.png](evidence/week-00/azure-budget.png) Tidak Ada screenshoot |
| Threshold budget alert | 50% / 80% / 100% | Azure Portal → Cost Management → Budgets | [azure-budget.png](evidence/week-00/azure-budget.png) Tidak Ada screenshoot |
| Kuota `Standard BS Family vCPUs` — Southeast Asia | 0/4 — 0 digunakan dari kuota 4 vCPU | Azure Portal → Subscriptions → Usage + quotas | [azure-vcpu-quota.png](evidence/week-00/azure-vcpu-quota.png) Tidak Ada screenshoot |
| Limit RAM WSL2 yang dikonfigurasi | 5 GB | `.wslconfig` | [wslconfig.txt](evidence/week-00/wslconfig.txt) |
| Processor WSL2 yang dikonfigurasi | 8 | `.wslconfig` | [wslconfig.txt](evidence/week-00/wslconfig.txt) |
| Swap WSL2 yang dikonfigurasi | 8 GB | `.wslconfig` | [wslconfig.txt](evidence/week-00/wslconfig.txt) |
| RAM WSL2 yang terdeteksi setelah restart | 4.8 GiB | `free -h` | [wsl-resources.txt](evidence/week-00/wsl-resources.txt) |
| Swap WSL2 yang terdeteksi | 8.0 GiB | `free -h` | [wsl-resources.txt](evidence/week-00/wsl-resources.txt) |
| Logical processor yang tersedia di WSL2 | 8 | `nproc` | [wsl-resources.txt](evidence/week-00/wsl-resources.txt) |
| Permission SSH private key | `600` / `-rw-------` | `ls -l ~/.ssh/campus_lms_azure ~/.ssh/campus_lms_azure.pub` | [ssh-key-permissions.txt](evidence/week-00/ssh-key-permissions.txt) |
| `.env` diabaikan oleh Git | Terverifikasi | `git check-ignore -v .env` | [env-gitignore.txt](evidence/week-00/env-gitignore.txt) |
| Target `make todo` | Berhasil dijalankan | `make todo` | [make-todo.txt](evidence/week-00/make-todo.txt) |
| Pengiriman budget alert ke email | NOT MEASURED | Azure Cost Management + email penerima | Belum diverifikasi |
| Kredit Azure tersisa | NOT MEASURED | Azure Portal → Cost Management | Belum diukur |
| Proyeksi biaya bulanan | NOT MEASURED | Azure Portal → Cost Management / Cost analysis | Belum diukur |

**Perbandingan sebelum/sesudah:**

Tidak ada optimasi performa yang dilakukan pada Minggu 0. Perubahan `.wslconfig` merupakan konfigurasi resource development environment, bukan optimasi performa aplikasi, sehingga tidak dibuat klaim peningkatan performa dari perubahan tersebut.

## 6. Konsep yang Dipelajari

### Spending limit vs budget alert

* **Apa:** Spending limit adalah guardrail billing pada subscription yang membatasi penggunaan kredit, sedangkan budget alert adalah mekanisme pemantauan yang memberi peringatan ketika pengeluaran mencapai threshold tertentu.
* **Kenapa dipakai di sini:** `campus-lms` menggunakan Azure for Students dan project sengaja dirancang dengan disiplin biaya yang ketat. Spending limit menjadi perlindungan utama, sementara budget alert memberikan peringatan lebih awal sebelum kredit digunakan terlalu banyak.
* **Alternatif yang tidak dipilih:** Menghapus spending limit dan hanya mengandalkan budget alert. Alternatif ini memberikan fleksibilitas lebih besar tetapi membuka kemungkinan biaya melewati kredit yang tersedia.
* **Cara membuktikan sendiri:** Buka Azure Portal → Subscription dan verifikasi spending limit, kemudian buka Cost Management → Budgets dan periksa `campus-lms-guard`.
* **Pertanyaan interview terkait:** *"Bagaimana kamu mencegah cost overrun di cloud?"*. **Jawaban Manusia:** *"Dengan konfigurasi spending limit untuk membatasi spending billing agar tagihan tidak bengkak melebihi kredit gratis yang diberikan (saat ini kredit gratis azure sebesar $100 selama 12 bulan) dan penggunaan Alert berdasarkan konfigurasi budget ketika sudah 50%, 80%, dan 100% yang dikirimkan ke email manusia"*.

### `az vm stop` vs `az vm deallocate`

* **Apa:** `az vm stop` menghentikan sistem operasi VM, tetapi resource compute dapat tetap dialokasikan. `az vm deallocate` menghentikan VM sekaligus melepaskan alokasi compute sehingga billing compute berhenti; resource persisten lain seperti disk tetap dapat memiliki biaya.
* **Kenapa penting di sini:** Project menggunakan kredit Azure for Students yang harus dijaga agar tidak terbuang karena VM tetap berada pada status allocated ketika tidak digunakan.
* **Alternatif yang tidak dipilih:** Membiarkan VM terus hidup selama development. Ini lebih praktis tetapi menggunakan resource cloud terus-menerus.
* **Cara membuktikan sendiri:** Setelah VM tersedia pada Minggu 4, bandingkan state VM menggunakan `az vm get-instance-view` setelah `stop` dan setelah `deallocate`, kemudian cocokkan dengan Cost Management.
* **Pertanyaan interview terkait:** *"Apa perbedaan menghentikan dan mendeallocate VM di cloud, dan bagaimana hal itu memengaruhi biaya?"*. **Jawaban Manusia:** *"Menghentikan/stop VM di cloud itu berarti kita hanya mematikkan VM, tetapi beban compute billing masih dialokasikan ke VM kita. sementara kalo untuk deallocate, itu berarti kita mematikan VM sekaligus membuang alokasi compute dari VM kita, sehingga compute billing tidak terus berjalan ketika tidak digunakan"*.

### Private key vs public key SSH

* **Apa:** Private key adalah credential rahasia yang membuktikan identitas pemilik, sedangkan public key adalah bagian yang boleh diberikan kepada server untuk memverifikasi autentikasi.
* **Kenapa dipakai di sini:** VM Azure pada Minggu 4 akan menggunakan SSH key khusus `campus_lms_azure` agar autentikasi tidak bergantung pada password.
* **Alternatif yang tidak dipilih:** Login SSH berbasis password, karena meningkatkan risiko credential reuse dan brute-force.
* **Cara membuktikan sendiri:** `ls -l ~/.ssh/campus_lms_azure*`
* **Pertanyaan interview terkait:** *"Bagaimana public-key authentication SSH bekerja dan kenapa private key tidak boleh masuk repository?"*. **Jawaban Manusia:** *"public key auth SSH bekerja dengan menempatkannya di server dan mencocokkannya dengan private key di komputer local agar hanya pihak yang valid dengan private key yang cocok yang bisa login, private key tidak boleh masuk ke repository karena berisiko disalahgunakan oleh pihak yang menyamar menggunakan identitas kita"*.

### Environment variable dan secret management

* **Apa:** Konfigurasi aplikasi disediakan melalui environment variable, sedangkan secret asli disimpan di `.env` lokal dan tidak dimasukkan ke source control.
* **Kenapa dipakai di sini:** API key LLM, JWT secret, dan credential lain tidak boleh berada di source code maupun file yang di-commit.
* **Alternatif yang tidak dipilih:** Hard-code credential di source code atau config yang masuk Git.
* **Cara membuktikan sendiri:** `git check-ignore -v .env`
* **Pertanyaan interview terkait:** *"Bagaimana kamu mengelola secret antara development dan production?"*. **Jawaban Manusia:** *"pada mode development, key secret harus boleh hanya berada pada komputer lokal dan tidak boleh diakses sembarangan orang, itu kenapa tidak boleh sampai masuk ke git ataupun dengan hard-code, ini bisa disalahgunakan seperti mengakses api llm yang beban billingnya kena kepada kita, atau hal semacam itu. kalo untuk production, kita menggunakan mekanisme provisioning untuk menerapkan IaC ketika membuat server untuk dibuatkan sertifikat khusus demi keamanan komunikasi dalam jaringan"*.

### WSL resource limiting

* **Apa:** `.wslconfig` mengatur batas resource global WSL2 dari sisi Windows.
* **Kenapa dipakai di sini:** Laptop development hanya memiliki RAM 8 GB, sehingga WSL tidak boleh mengambil seluruh memory dan membuat Windows atau Docker kehabisan resource.
* **Alternatif yang tidak dipilih:** Membiarkan WSL menggunakan konfigurasi resource default tanpa batas eksplisit.
* **Cara membuktikan sendiri:** `free -h && nproc`
* **Pertanyaan interview terkait:** *"Bagaimana kamu mengelola resource development environment pada mesin dengan RAM terbatas?"*. **Jawaban Manusia:** *"ram 8gb adalah kapasitas yang terbatas untuk menjalankan komputasi berat, maka WSL sebagai virtual machine wajib dibatasi agar tidak membebani window/host secara brutal. maka pada wsl harus di konfigurasi agar tidak menggunakan seluruh kapasitas ram yang tersedia di komputer lokal."*.

### Agent-first, evidence-verified

* **Apa:** AI agent boleh melakukan sebagian besar pekerjaan yang dapat diotomatisasi secara lokal, tetapi setiap hasil harus diverifikasi dengan command dan evidence yang dapat diperiksa manusia.
* **Kenapa dipakai di sini:** Tujuan project bukan hanya menghasilkan software, tetapi memastikan pemilik mampu memahami dan mempertahankan setiap keputusan teknis.
* **Alternatif yang tidak dipilih:** Memberikan agent kebebasan penuh dan menerima hasil berdasarkan klaim bahwa implementasi “seharusnya bekerja”.
* **Cara membuktikan sendiri:** Jalankan acceptance test `make todo`, lalu cocokkan ringkasan agent dengan output dan isi repository.
* **Pertanyaan interview terkait:** *"Bagaimana kamu menggunakan AI coding agent tanpa kehilangan pemahaman dan kontrol terhadap codebase?"*. **Jawaban Manusia:** *"dengan membaca dan memahami laporan hasil yang sudah ditambahkan detail oleh agent dan juga melalui pemahaman dari kode apa saja yang berubah, file apa saja yang ditambah/dihapus, termasuk sekaligus menjalankan command verification dan evidence, dan masih banyak lagi"*.

---

## 7. Belum Terverifikasi

| Hal                                                                       | Kenapa belum terverifikasi                                                                          | Rencana verifikasi                                  |
| ------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | --------------------------------------------------- |
| Budget alert benar-benar masuk ke email yang dibaca manusia               | Threshold belum diuji sampai menghasilkan notifikasi nyata                                          | Turunkan threshold sementara dan simpan bukti email |
| SKU VM production final tersedia di Southeast Asia untuk subscription ini | VM belum diprovision dan availability SKU dapat berbeda ketika provisioning dilakukan               | Verifikasi ulang saat Minggu 4 sebelum membuat VM   |
| Resource envelope VM production cukup untuk workload nyata                | Aplikasi production belum dibuat dan belum dilakukan profiling                                      | Benchmark dan observability pada minggu berikutnya  |
| Agent konsisten mematuhi `AGENTS.md` pada engineering task nyata          | Acceptance test baru mencakup workflow awal `make todo`                                             | Verifikasi pada seluruh task Minggu 1               |
| Seluruh angka setup Week 0 memiliki evidence file di repository           | Sebagian verifikasi dilakukan interaktif tetapi belum disimpan ke `docs/progress/evidence/week-00/` | Simpan screenshot/raw output sebelum sign-off       |
| Behavior aplikasi di production                                           | Aplikasi belum diimplementasikan dan belum dideploy                                                 | Minggu 4 dan seterusnya                             |
| JetBrains Student Subscription                                            | Benefit sengaja tidak diklaim pada sesi Week 0 ini                                                  | Klaim nanti jika memang diperlukan                  |

**Asumsi yang dipakai tapi belum dibuktikan:**

* Target workload production dapat dijalankan pada Azure B-series dengan resource envelope kecil.
* Free tier layanan eksternal yang direncanakan masih mencukupi ketika workload nyata tersedia.
* Pemisahan komponen berat dari VM Azure akan cukup menjaga konsumsi RAM production; hal ini baru dapat dibuktikan setelah sistem dapat dijalankan dan diukur.

---

## 8. Masalah & Cara Diselesaikan

### Masalah: `make todo` tidak dapat dijalankan pada acceptance test pertama

* **Gejala:** Agent melaporkan bahwa `make todo` tidak dapat dijalankan secara literal. Verifikasi manual di WSL menghasilkan `Command 'make' not found`.
* **Hipotesis yang salah:** Pada awalnya terdapat kemungkinan bahwa masalah berasal dari akses shell/WSL agent karena percobaan eksekusi juga sempat menghasilkan masalah akses.
* **Akar masalah:** GNU Make memang belum terpasang pada environment WSL yang digunakan untuk development.
* **Solusi:** Manusia memasang GNU Make sebagai system package, memverifikasi ketersediaannya, kemudian menjalankan ulang acceptance test agent.
* **Pencegahan:** Prerequisite development environment harus diverifikasi sebelum task yang bergantung pada command tersebut dijalankan.
* **Waktu terbuang:** sebentar saja karena ini hanya masalah kecil, sekitar 1-2 menit.

### Masalah: Agent mengganti `make todo` dengan pencarian TODO pada percobaan pertama

* **Gejala:** Karena command asli tidak tersedia, agent mencoba memperoleh hasil ekuivalen dengan membaca pola TODO dari repository.
* **Akar masalah:** Environment belum memenuhi prerequisite untuk menjalankan command yang secara eksplisit diminta.
* **Solusi:** Hasil percobaan pertama tidak dianggap sebagai acceptance test yang lulus. Environment diperbaiki dan prompt yang sama dijalankan ulang.
* **Pencegahan:** Jika command verifikasi wajib tidak dapat dijalankan, agent harus melakukan hard stop atau menyatakan task belum terverifikasi, bukan mengganti verification step dengan pendekatan lain.
* **Waktu terbuang:** sebentar saja karena ini hanya masalah kecil, sekitar 3-5 menit.

---

## 9. Status Definition of Done

> Status di kolom **Usulan agent** hanya merupakan evaluasi berdasarkan pekerjaan yang sudah dilakukan. Kolom terakhir hanya boleh dicentang manusia setelah memeriksa bukti.

| DoD Minggu 0                                                   | Usulan agent                                                                   | Bukti                                                        | Dicentang manusia |
| -------------------------------------------------------------- | ------------------------------------------------------------------------------ | ------------------------------------------------------------ | ----------------- |
| Spending limit ON dan budget alert terbukti terkirim           | ⚠️ sebagian — spending limit sudah diverifikasi, pengiriman alert belum        | Evidence spending limit dan email alert masih perlu disimpan | ☐                 |
| Kuota vCPU dikonfirmasi tersedia di region pilihan             | ✅ terpenuhi berdasarkan verifikasi manusia; evidence repo masih perlu disimpan | Screenshot Usage + quotas perlu disimpan (catatan manusia: screenshoot lupa dilakukan, tapi sudah dipastikan tersedia)                     | ✅                 |
| Repo publik hidup di GitHub, `make todo` jalan di laptop       | ✅ terpenuhi berdasarkan acceptance test                                        | Raw output `make todo` perlu disimpan sebagai evidence       | ✅                 |
| Saya bisa menjelaskan beda `az vm stop` dan `az vm deallocate` | ✅ sudah dikonfirmasi manusia melalui explain-back                            | Explain-back tersedia                                  | ✅                 |
| `docs/progress/week-00.md` terisi dan ditandatangani           | ✅ laporan terisi; sudah ditandatangani                                        | File ini                                                     | ✅                 |

---

## 10. Untuk Minggu Depan

* **Carry-over:**

  * Verifikasi budget alert benar-benar masuk ke email.
  * Simpan evidence Week 0 yang masih hanya tersedia dari verifikasi interaktif.
  * Selesaikan human sign-off dan tandatangani laporan.
  * JetBrains Student Subscription tetap opsional/carry-over jika nantinya diperlukan.

* **Utang teknis yang sengaja diambil:**

  * SKU VM production belum dikunci; keputusan final menunggu verifikasi availability pada Minggu 4.
  * Belum ada benchmark atau profiling karena aplikasi belum diimplementasikan.
  * Sebagian besar evidence Week 0 masih harus diformalisasi menjadi file di `docs/progress/evidence/week-00/`.

* **Persiapan Minggu 1:**

  * WSL2 sudah dikonfigurasi.
  * Azure CLI sudah tersedia.
  * GNU Make sudah tersedia.
  * Repository dan agent governance sudah siap.
  * Jalankan `make todo` sebagai sumber pekerjaan Minggu 1.
  * Gunakan task loop `ORIENT → PLAN → CONFIRM → IMPLEMENT → VERIFY → EVIDENCE → REPORT → TEACH` pada setiap pekerjaan agent.

---

## 11. Verifikasi Manusia

* [x] Saya sudah spot-check minimal 3 file bukti secara acak dan isinya cocok dengan klaim
* [x] Skor quiz: 86 / 100 (minimal 70% untuk lanjut)
* [x] Explain-back 3 menit sudah direkam: `docs/progress/explain/week-00.<mp3|txt>`
* [x] Saya sudah benar-benar membaca `AGENTS.md` dan `agent/policy.md`, bukan sekadar membuka
* [x] Saya dapat menjelaskan fungsi spending limit dan budget alert tanpa membuka catatan
* [x] Saya dapat menjelaskan perbedaan `az vm stop` dan `az vm deallocate` tanpa membuka catatan
* [x] Saya sudah memastikan credential, API key, private key, dan `.env` tidak masuk repository

**Self-audit agent (dari `agent/evidence-protocol.md` §8):** Belum dapat dinyatakan selesai sampai evidence Week 0 diformalisasi dan self-audit dijalankan sesuai protokol.

**Ditandatangani:** FfFfFfFf  **Tanggal:** 16 Agustus 2026
