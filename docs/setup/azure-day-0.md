# Checklist Minggu 0 — Persiapan Akun & Agent

> **100% tugas manusia.** Bukan karena agent tidak mampu, tapi karena portal
> cloud, verifikasi identitas, dan kredensial memang tidak boleh didelegasikan.
> Estimasi total: ±4 jam. Kerjakan sebelum Minggu 1.
>
> Versi naratif lengkap ada di roadmap bagian "MINGGU 0".

## 0.1 Kunci pengaman akun Azure (±45 menit)

- [x] **Subscriptions** → offer bernama **"Azure for Students"** (bukan "Free Trial" / "Pay-As-You-Go")
- [x] **Spending limit: ON** — selama ini aktif, kamu secara struktural tidak bisa ditagih
- [x] ⚠️ Jangan pernah terima tawaran *"Remove spending limit"* — itu satu-satunya jalan akun berubah jadi berbayar
- [x] **Cost Management → Budgets** → buat `campus-lms-guard`, **$10/bulan**, alert 50% / 80% / 100%
- [ ] Alert masuk ke email yang benar-benar kamu baca (tes: turunkan sementara ambangnya, pastikan email datang)
- [X] **Subscriptions → Usage + quotas** → region **Southeast Asia** → `Standard BS Family vCPUs` **= 0/4**
- [x] Catat konvensi di `docs/adr/0002c-azure-conventions.md`:
  - Region: **Southeast Asia** (Singapore — terdekat dari Bali)
  - Resource group: `rg-campuslms-prod`
  - Tag wajib: `project=campus-lms`, `env=prod`, `owner=<nama>`

**Belum dikerjakan sekarang:** VM (Minggu 4) · Neon (Minggu 4) · Cloudflare (Minggu 4) · AKS (jangan pernah — melahap kredit)

## 0.2 Klaim benefit Student Pack (±30 menit)

- [x] **Namecheap** — domain `.me` gratis 1 tahun. Klaim sekarang meski baru
      dipakai Minggu 4 (kupon/offer memiliki masa berlaku)
- [x] **GitHub Education / Student Developer Pack** — akun mahasiswa
      sudah terverifikasi. Benefit aktif sampai **13 Maret 2028**.
      GitHub Codespaces tersedia hingga **180 core hours/bulan** untuk
      akun personal.
- [ ] **JetBrains** — klaim student subscription untuk PyCharm dan IDE
      JetBrains lainnya. Benefit student saat ini berupa subscription gratis
      yang dapat diperpanjang setiap tahun.
- [X] **New Relic** — klaim benefit student untuk APM/observability gratis,
      untuk digunakan sebagai pembanding dengan observability self-host
      pada Minggu 6.

## 0.3 Kredensial & lingkungan lokal (±60 menit)

- [x] SSH key dibuat:
      `ssh-keygen -t ed25519 -C "campus-lms-azure" -f ~/.ssh/campus_lms_azure`
- [x] `chmod 600 ~/.ssh/campus_lms_azure` — private key **tidak pernah** masuk repo
- [x] Azure CLI: `curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash`
- [x] `az login --use-device-code` lalu `az account show --output table` menampilkan langganan student
- [x] Daftar API key LLM gratis **tanpa kartu** (baru dipakai Minggu 7, tapi daftar sekarang):
  - [x] Google AI Studio (Gemini)
  - [x] Groq
  - [x] Cerebras
  - [x] OpenRouter
- [x] `cp .env.example .env` dan isi yang sudah diketahui — **agent dilarang menyentuh file ini**
- [x] `.wslconfig` diatur kalau RAM 8 GB: `memory=5GB`, `processors=8`, `swap=8GB`

## 0.4 Siapkan agent (±60 menit)

- [x] Repo `campus-lms` sudah ada di laptop
- [ ] `git init` → commit pertama → push ke GitHub sebagai **repo publik** (Actions gratis unlimited)
- [ ] Ganti `CHANGE_ME` di `apps/api/go.mod` dan `Makefile` dengan username GitHub-mu
- [ ] **Baca `AGENTS.md` sampai selesai** — ini bacaan wajib, bukan formalitas
- [ ] **Baca `agent/policy.md`** — terutama tabel hard stop
- [ ] **Baca `agent/README.md`** — cara memakai sistem ini sehari-hari
- [ ] Uji agent dengan prompt: *"Baca AGENTS.md dan Makefile, lalu jalankan `make todo` dan ringkas pekerjaan Minggu 1. Jangan menambah task yang tidak ada di repo."*
  - Kalau agent mengarang task yang tidak ada → perbaiki setup konteks sebelum lanjut
- [ ] Salin `agent/templates/weekly-report.md` → `docs/progress/week-00.md` dan isi

## Definition of Done Minggu 0

- [ ] Spending limit ON dan budget alert terbukti terkirim
- [ ] Kuota vCPU dikonfirmasi tersedia di region pilihan
- [ ] Repo publik hidup di GitHub, `make todo` jalan di laptopmu
- [ ] Saya bisa menjelaskan beda `az vm stop` dan `az vm deallocate`
      *(petunjuk: hanya satu yang menghentikan tagihan compute)*
- [ ] `docs/progress/week-00.md` terisi dan ditandatangani

**Tanggal selesai:** __________  **Tanda tangan:** ______________
