# Roadmap 12 Minggu: Dari Mahasiswa TI Semester 7 → Backend/Fullstack Engineer yang Kuat AI

**Disusun:** 16 Agustus 2026
**Diperbarui:** 20 Agustus 2026
**Sinkronisasi domain:** `docs/domain.md` §37–38 dan `docs/domain-ai.md` §5 menjadi sumber kebenaran untuk urutan implementasi entity. Jika daftar deliverable lama bertentangan dengan tier implementasi, tier di `domain.md` yang berlaku. Tabel Tier 3 tidak boleh muncul di migrasi tanpa keputusan tertulis untuk menaikkan tier.
**Sinkronisasi progres:** status Minggu 0–2 direkonsiliasi dengan laporan manusia yang ditandatangani. `✅` berarti sudah dikonfirmasi manusia pada laporan/evidence; `☐` berarti belum selesai atau belum diverifikasi. Perubahan scope tidak boleh dipakai untuk mengubah kegagalan menjadi keberhasilan secara retroaktif.
**Untuk:** Mahasiswa TI semester 7, pengalaman frontend + backend + Supabase, terbiasa memakai AI coding agent via CLI (opencode + 9router/freebuff dengan combo `heavy/plan/dev/review/spark/quick`)
**Target peran:** Backend/Fullstack Engineer dengan spesialisasi AI (AI Product Engineer)
**Pasar:** Indonesia dulu (Bali/Jakarta) → remote internasional
**Intensitas:** ±35 jam/minggu × 12 minggu ≈ **420 jam**
**Budget infrastruktur:** Rp 0 — **dan tanpa kartu kredit/debit sama sekali** (semua jalur pembayaran: gratis, verifikasi status mahasiswa, atau QRIS/GoPay/OVO)
**Project pengikat:** **SaaS Learning Management System multi-tenant** (satu project, dibangun bertahap, setiap minggu menambah satu lapisan skill industri)

---

## Daftar Isi

1. [Cara Membaca Dokumen Ini](#1-cara-membaca-dokumen-ini)
2. [Analisis Gap: Yang Kamu Punya vs Standar Industri 2026](#2-analisis-gap-yang-kamu-punya-vs-standar-industri-2026)
3. [Kendala Nyata: Hardware & Budget](#3-kendala-nyata-hardware--budget)
4. [Keputusan Stack (dan Alasannya)](#4-keputusan-stack-dan-alasannya)
5. [Arsitektur Target Project LMS](#5-arsitektur-target-project-lms)
6. [Peta Infrastruktur Gratis 2026 (Tanpa Kartu)](#6-peta-infrastruktur-gratis-2026-tanpa-kartu)
7. [Mode Kerja Agent-First & Aturan Anti-Halusinasi](#7-naik-level-cara-pakai-ai-dari-konsumen--pembangun)
8. [Roadmap Mingguan (Minggu 0–12)](#8-roadmap-mingguan)
9. [Ritme Harian & Aturan Main](#9-ritme-harian--aturan-main)
10. [Portofolio, CV, dan Strategi Melamar](#10-portofolio-cv-dan-strategi-melamar)
11. [Checklist Skill Akhir](#11-checklist-skill-akhir)
12. [Sumber & Catatan Verifikasi](#12-sumber--catatan-verifikasi)

---

## 1. Cara Membaca Dokumen Ini

Setiap minggu punya format yang sama:

| Bagian | Maksudnya |
|---|---|
| **Tujuan** | Satu kalimat: kemampuan apa yang kamu miliki di akhir minggu |
| **Konsep wajib** | Teori minimum. Bukan untuk dihafal, untuk dipakai hari itu juga |
| **Deliverable project** | Kode/infra nyata yang masuk ke repo LMS |
| **Definition of Done (DoD)** | Kriteria objektif. Kalau belum lolos, jangan lanjut minggu berikutnya |
| **Alokasi jam** | Pembagian 35 jam |
| **Sinyal CV** | Kalimat konkret yang boleh kamu tulis di CV setelah minggu itu selesai |

**Aturan besi:** DoD yang belum terpenuhi tetap berstatus terbuka dan tidak boleh dicentang. Jika manusia secara sadar mengubah scope atau memindahkan item ke minggu lain, perubahan itu harus tertulis di laporan/roadmap sebagai **carry-over atau scope revision**, bukan dianggap selesai secara retroaktif. Lebih baik roadmap selesai lebih lambat dengan status yang jujur daripada semua kotak tampak hijau tanpa bukti.

---

## 2. Analisis Gap: Yang Kamu Punya vs Standar Industri 2026

### 2.1 Yang sudah kamu kuasai (aset)

- Frontend + backend + database (Supabase) → kamu bisa **membangun fitur**.
- Go, Python, Node.js — Go paling kuat → bahasa backend yang serius dan bernilai tinggi.
- Terbiasa dengan AI coding agent (IDE extension, built-in agent, CLI) dan sudah mengatur routing model per-tugas lewat `opencode.json` → ini **di atas rata-rata fresh grad**. Kemampuan verifikasi & orkestrasi AI memang salah satu ekspektasi baru untuk junior 2026. ([nucamp, Jan 2026](https://www.nucamp.co/blog/the-junior-developer-hiring-crisis-in-2026-how-to-get-your-first-backend-job))

### 2.2 Gap yang nyata

Yang sering terjadi pada profil sepertimu: **kamu bisa membuat aplikasi, tapi belum bisa mengoperasikan sistem.** Supabase menyembunyikan seluruh lapisan yang justru ditanyakan saat interview: koneksi pooling, index, migrasi, backup, TLS, proses yang mati jam 3 pagi, siapa yang tahu kalau error naik 5%.

Ekspektasi pasar untuk "junior" 2026 (dari posting kerja & panduan hiring):

| Area | Ekspektasi minimum | Statusmu | Prioritas |
|---|---|---|---|
| Satu bahasa backend mendalam + REST idiomatik | Wajib | Ada (Go) | Perkuat |
| SQL: join, index, transaction, desain skema | Wajib | Parsial (lewat Supabase) | **Tinggi** |
| Docker + docker compose | Wajib | Belum | **Kritis** |
| CI/CD (GitHub Actions) | Wajib | Belum | **Kritis** |
| Deploy & operasikan di Linux VPS / cloud | Wajib | Belum | **Kritis** |
| Testing otomatis (unit + integration) | Wajib | Belum jelas | **Tinggi** |
| Observability: log terstruktur, metrics, tracing | Nilai plus kuat, makin jadi baseline | Belum | Tinggi |
| Keamanan API dasar (OWASP API Top 10, JWT/RBAC, rate limit, secrets) | Wajib disebut di banyak posting | Parsial | Tinggi |
| Kubernetes | *Helpful* di level junior, wajib di mid | Belum | Sedang |
| Terraform / IaC | Preferred, bukan blocker junior | Belum | Sedang-rendah |

Sumber ekspektasi di atas: ringkasan posting kerja backend & panduan hiring 2026 yang secara konsisten menempatkan Docker + CI/CD + cloud + SQL sebagai non-negotiable, dan K8s/Terraform sebagai "helpful di junior, wajib di mid". ([KORE1, Jul 2026](https://www.kore1.com/hire-backend-developer/), [nucamp](https://www.nucamp.co/blog/the-junior-developer-hiring-crisis-in-2026-how-to-get-your-first-backend-job))

### 2.3 Gap di sisi AI Engineer

Ini penting dan sering disalahpahami. Kemampuanmu sekarang = **pengguna canggih AI coding agent**. Itu meningkatkan produktivitas, **tapi bukan skill AI Engineer**. Tidak ada perusahaan yang membayarmu karena kamu bisa `/plan` di opencode.

AI Engineer 2026 = orang yang **mengintegrasikan LLM ke dalam produk** dan bertanggung jawab atas kualitas, biaya, dan keandalannya. Yang dicari:

| Kemampuan | Status di pasar 2026 |
|---|---|
| Panggil LLM API + prompt engineering + structured output | **Baseline** (dianggap sudah bisa, tidak bikin kamu menonjol) |
| RAG sederhana (chunk → embed → retrieve → jawab) | **Baseline** |
| Vector DB / pgvector, embedding, chunking strategy | Baseline-menengah |
| **Evaluation harness** (ragas / promptfoo / DeepEval, LLM-as-judge, regression test di CI) | **Pembeda utama** — disebut sebagai "universal screen" di interview 2026 |
| **Agentic RAG** (retrieve iteratif, agen menilai kecukupan informasi, validasi sumber) | Pembeda |
| **Tool calling + MCP** (Model Context Protocol) | Pembeda kuat |
| **Observability AI**: tracing, token/cost tracking, latency p95 (Langfuse/LangSmith) | Pembeda |
| Guardrails & red-teaming (prompt injection, PII) | Naik cepat |
| Cost optimization: model routing, caching, prompt compression | Pembeda |

Sumber: analisis skill AI engineer 2026 yang konsisten menyebut agent orchestration, MCP, evaluation design, agentic RAG, dan production observability sebagai skill paling terdiferensiasi — sementara "basic LLM API integration + simple RAG" sudah turun jadi baseline. ([gnxt, Jun 2026](https://gnxtsystems.com/the-ai-skills-that-mattered-in-2025-are-already-obsolete-heres-what-2026-demands/), [technovids](https://technovids.com/ai-engineer-skills), [doit.software](https://doit.software/blog/ai-developer-skills))

Satu kutipan yang harus kamu pegang: **"Projects must be deployed and publicly accessible — not just in Colab notebooks."** ([technovids](https://technovids.com/ai-engineer-skills)) Itulah kenapa roadmap ini menaruh VPS/Docker/CI di depan, baru AI.

### 2.4 Soal Hermes / OpenClaw — jawaban jujur

Kamu menyebut belum punya pengalaman "hermes agent, openclaw, dan lain sebagainya". Klarifikasi supaya kamu tidak salah alokasi waktu:

- **OpenClaw** itu nyata: framework AI agent self-hosted open-source (rilis November 2025 sebagai "Clawdbot", ganti nama Januari 2026, dibuat Peter Steinberger), berjalan sebagai proses Node.js headless, terhubung ke WhatsApp/Telegram/Discord/Slack, arsitektur berbasis "skills". Populer sekali secara komunitas. ([petronellatech](https://petronellatech.com/blog/openclaw-ai-agent-guide/), [Medium/Kanerika](https://medium.com/@kanerika/openclaw-how-a-self-hosted-ai-agent-changed-automation-in-2026-6ba728345d53))
- **Tetapi**: OpenClaw adalah *personal assistant automation*, bukan skill yang diminta di job description AI Engineer. Yang diminta di JD adalah MCP, orkestrasi agent (LangGraph/OpenAI Agents SDK/CrewAI), evaluasi, dan observability.
- Untuk "hermes agent" saya **tidak menemukan rujukan yang jelas dan konsisten** sebagai standar industri. Bisa jadi maksudmu model *Hermes* dari Nous Research (keluarga model open-weight yang kuat di function calling), atau tool internal komunitas. **Saya tidak akan menebak** — ini bagian dari komitmen tanpa halusinasi.

**Keputusan roadmap:** OpenClaw & sejenisnya masuk kategori *optional weekend fun* (ada slot di Minggu 12), **bukan** jalur utama. Waktumu 420 jam terlalu berharga untuk dipakai mengejar tool yang tidak muncul di JD.

---

## 3. Kendala Nyata: Hardware & Budget

### 3.1 Laptop: Axioo Hype 5 AMD X5-2

Spesifikasi resmi varian umum:

| Komponen | Spesifikasi |
|---|---|
| CPU | AMD Ryzen 5 7430U — 6 core / 12 thread, 2.3 GHz turbo 4.3 GHz (Zen 3, 15W) |
| GPU | AMD Radeon Vega 7 (iGPU, **tanpa VRAM dedicated**) |
| RAM | **8 GB DDR4-3200** pada varian dasar (2 slot SODIMM, upgradeable hingga 64 GB); ada varian 16 GB |
| Storage | 256 GB NVMe PCIe Gen3 (ada varian 512 GB), upgradeable |
| Layar | 14" FHD IPS 60Hz |

([liputan6](https://www.liputan6.com/tekno/read/6280876/harga-axioo-hype-5-amd-x5-2-laptop-lokal-stylish-dengan-performa-tangguh-dan-baterai-tahan-lama), [ponselesa](https://www.ponselesa.com/2024/12/spesifikasi-dan-harga-axioo-hype-5-amd-x5-2.html), [agres.id — varian 16 GB](https://www.agres.id/products/axioo-hype-5-amd-x5-2-ryzen-5-7430-16gb-256gb-w11-140fhd-ips-blit-hdmi-gry))

**Konsekuensi teknis yang harus kamu terima:**

1. **CPU-mu cukup.** 6C/12T Zen 3 sangat memadai untuk Docker, Postgres, Go build, dan Python.
2. **RAM adalah bottleneck utama.** Kalau kamu di varian 8 GB: WSL2 + Docker + Postgres + Redis + Grafana stack + browser + IDE akan swap-thrashing.
   → **Rekomendasi #1 dan satu-satunya belanja yang saya sarankan:** tambah 1 keping SODIMM DDR4-3200 8 GB (slot kedua kosong) supaya jadi 16 GB dual-channel. Harganya jauh lebih murah daripada nilai waktumu yang hilang karena laptop tersendat, dan bonusnya iGPU Vega 7 ikut lebih cepat karena dual-channel. Kalau tidak memungkinkan, roadmap ini tetap jalan — lihat mitigasi di bawah.
3. **Tidak ada GPU NVIDIA.** Artinya:
   - **Jangan** rencanakan fine-tuning, training, atau menjalankan model besar lokal.
   - Menjalankan LLM lokal via Ollama hanya realistis untuk model kecil (≤4B, kuantisasi Q4) dan itu pun lambat. **Bukan bagian wajib roadmap ini.**
   - Semua kebutuhan LLM diarahkan ke **free tier API cloud** (lihat §6). Ini justru sesuai praktik industri: AI Engineer memakai API, bukan melatih model.
4. **Storage 256 GB akan sesak.** Docker images + WSL2 + node_modules dapat menghabiskan ruang dengan cepat. Untuk pembersihan rutin gunakan `make prune` / `docker system prune -af` **tanpa `--volumes`**. Penghapusan volume dipisahkan ke `make prune-hard` dan wajib memakai konfirmasi eksplisit, karena named volume Postgres adalah data, bukan sampah build.

**Mitigasi kalau tetap di 8 GB RAM:**

| Masalah | Mitigasi |
|---|---|
| Compose stack terlalu berat | Jalankan per-profile: `docker compose --profile core up` (app+db) vs `--profile obs` (Grafana/Prometheus/Loki). Jangan semua sekaligus |
| Grafana stack rakus | Jalankan observability stack hanya saat dibutuhkan (profile `obs`), atau pindahkan ke **GitHub Codespaces** (gratis, tanpa kartu) dan akses lewat browser |
| Kubernetes lokal | Pakai **k3d** (k3s dalam Docker, 1 server + 1 agent), bukan Minikube/kind multi-node. Batasi memori |
| WSL2 memakan RAM | Buat `%UserProfile%\.wslconfig`: `memory=5GB`, `processors=8`, `swap=8GB` |
| IDE berat | Pakai VS Code + remote WSL, matikan extension yang tidak dipakai. Editor terminal (helix/neovim) untuk sesi panjang |

### 3.2 Konsekuensi budget Rp 0 **tanpa kartu**

Kartu debit Mastercard BCA-mu selalu *declined* di platform internasional. Ini masalah yang sangat umum di Indonesia dan **bukan kesalahanmu** — penjelasan singkatnya ada di §6.5. Yang penting: **roadmap ini sudah saya susun ulang supaya tidak satu pun langkahnya memerlukan kartu.**

Tiga jalur pembayaran yang kita pakai:

| Jalur | Contoh | Butuh kartu? |
|---|---|---|
| **A. Gratis murni** | Cloudflare, Neon, Supabase, GitHub Actions, semua LLM API free tier | Tidak |
| **B. Verifikasi status mahasiswa** | Azure for Students ($100, 12 bulan), GitHub Student Developer Pack | Tidak — cukup email `.ac.id` atau foto KTM |
| **C. QRIS / GoPay / OVO** (opsional, hanya kalau butuh) | VPS lokal Rp 25–87rb/bln, domain `.my.id` Rp ~15–30rb/thn | Tidak |

Yang **dihapus** dari rencana sebelumnya karena mewajibkan kartu: **Oracle Cloud Always Free** (verifikasi kartu saat signup), AWS, Google Cloud, Fly.io, Railway, DigitalOcean. Semua sudah diganti dengan padanan yang setara atau lebih baik.

Konsekuensi lain yang harus kamu terima:
- Domain: pakai **Namecheap `.me` gratis dari GitHub Student Pack**, atau `.my.id` lewat registrar lokal (bayar QRIS), atau gratis total via **Cloudflare Tunnel** (`*.trycloudflare.com`) / DuckDNS.
- Kuota LLM terbatas per hari → kamu **dipaksa** belajar caching, model routing, dan token budgeting. Kebetulan itu justru salah satu skill pembeda 2026.
- Kredit Azure $100 punya batas waktu 12 bulan → kamu **dipaksa** belajar cost engineering: matikan VM saat tidak dipakai, scale-up hanya saat load test. Ini bahan cerita interview yang bagus, bukan kekurangan.

---

## 4. Keputusan Stack (dan Alasannya)

| Lapisan | Pilihan | Alasan |
|---|---|---|
| API utama | **Go** (net/http + `chi` atau Echo) | Bahasa terkuatmu. Go = sinyal kuat untuk backend/infra roles, binary kecil (ideal untuk VPS 12 GB), concurrency bawaan, build image Docker ~15 MB |
| AI service | **Python + FastAPI** | Seluruh ekosistem AI (LangGraph, ragas, LlamaIndex, sentence-transformers) hidup di Python. FastAPI disebut eksplisit di daftar skill AI Engineer 2026 |
| Pola arsitektur | **Modular monolith Go + 1 sidecar Python** | Realistis, bukan microservice-cosplay. Polyglot dua-service = persis pola industri untuk produk ber-AI |
| Frontend | **Next.js (App Router) + TypeScript + Tailwind** | Kamu sudah bisa frontend; ini stack paling umum di pasar lokal & remote |
| Database | **PostgreSQL 16 self-managed di Docker (lokal/dev)** + **Neon (produksi)**, keduanya dengan **pgvector** | Kamu harus lepas dari ketergantungan Supabase: di lokal kamu yang mengelola Postgres sendiri (migrasi, index, tuning, backup, restore) — di situlah pembelajarannya. Produksi dipindah ke Neon karena VM gratis hanya 1 GB RAM. **Wajib** kamu dokumentasikan sebagai keputusan arsitektur sadar, bukan kemalasan |
| Cache/queue | **Redis** (cache, rate limit) + **River** atau **asynq** (job queue Go) | Background job wajib untuk transcoding/ingest dokumen/eval batch |
| Object storage | **MinIO** lokal (belajar S3 API) → **Supabase Storage** di produksi | Upload materi kuliah. Keduanya gratis tanpa kartu. *(Cloudflare R2 sengaja dihindari: mengaktifkannya biasanya meminta metode pembayaran)* |
| Reverse proxy | **Caddy** | TLS otomatis (Let's Encrypt), config 5 baris. Nginx dipelajari sebagai pembanding |
| CI/CD | **GitHub Actions** + **GHCR** | Gratis unlimited untuk repo publik. Registry gratis |
| Observability | **OpenTelemetry SDK** → **Prometheus + Grafana + Loki + Tempo** (self-host) | Nama-nama ini muncul persis di JD. Semua open source |
| AI observability | **Langfuse (self-host via Docker)** | Tracing LLM, cost tracking, dataset & eval. Alternatif gratis dari LangSmith |
| Eval | **promptfoo** (CI regression) + **ragas** (kualitas RAG) | Dua nama yang paling sering muncul di daftar eval stack |
| Orkestrasi agent | **LangGraph** (Python) | Framework agent stateful yang paling banyak disebut di JD 2026 |
| Protokol tool | **MCP** (Model Context Protocol) | Standar 2026 untuk menghubungkan agent ke sistem eksternal |
| Container orchestration | **Docker Compose** (produksi utama) → **k3d/k3s** (belajar K8s) | Compose cukup dan jujur untuk skala portofolio; K8s dipelajari agar lolos screening |
| IaC | **Terraform** (provider Azure + Cloudflare) | Cukup "bisa baca & tulis dasar", tidak perlu dalam di level junior |
| Supabase | **Tetap dipakai — tapi sebagai perbandingan** | Di Minggu 3 kamu akan menulis dokumen "Supabase vs self-managed Postgres: apa yang disembunyikan". Ini bahan cerita interview yang bagus |

---

## 5. Arsitektur Target Project LMS

**Nama project:** `campus-lms` (SaaS LMS multi-tenant untuk universitas)

### 5.1 Kenapa LMS adalah pilihan yang tepat

LMS punya semua bentuk masalah yang membuktikan skill industri sekaligus:
- **Multi-tenancy** (banyak universitas dalam satu instance) → RLS, isolasi data, composite foreign key, sesuatu yang jarang dipahami fresh grad
- **Otorisasi berlapis** → tenant role melalui `memberships` + `membership_roles`, course-scoped role melalui `course_staff`, dan hak mahasiswa melalui active `enrollments`; role Super Admin tetap platform-scoped dan bukan hak akademik otomatis
- **Pemisahan `courses` dan `course_offerings`** → master mata kuliah tidak dicampur dengan pelaksanaan kelas pada semester tertentu
- **File besar** (PDF, slide, video) → object storage, background job, streaming
- **Data relasional kompleks** (course offering, module, lesson, enrollment, submission, grading) → SQL beneran
- **Beban baca tinggi** (mahasiswa membuka materi) → caching, index, N+1
- **Deadline & notifikasi** → scheduler, queue, idempotency
- **Muatan AI yang MASUK AKAL, bukan tempelan** → inilah kuncinya (§5.3)

### 5.2 Diagram arsitektur akhir (target Minggu 12)

Arsitektur ini **sengaja terdistribusi**: target awalnya memakai VM Azure B-series kecil, sehingga komponen berat digeser ke free tier lain. **SKU production belum dikunci pada Minggu 0**; `Standard_B1s` tetap kandidat utama dan harus diverifikasi kembali saat provisioning Minggu 4. Ini keputusan sadar agar roadmap tidak menganggap availability SKU sebagai fakta sebelum resource benar-benar dapat dibuat.

```
                            Internet
                               |
              [ Cloudflare — DNS / TLS / WAF / Tunnel ]
                    |                          |
          [ Cloudflare Pages ]        [ Caddy :443 ]
             Next.js frontend                |
             (gratis, unlimited BW)          |
                                    ==================================
                                    | Azure VM B-series (B1s kandidat)  |
                                    | SKU final diverifikasi Minggu 4 |
                                    ==================================
                                     /         |            \
                            [ api-go ]    [ ai-svc ]    [ worker-go ]
                            Go/chi :8080  FastAPI :8000  job queue
                            - auth JWT    - RAG pipeline  - ingest
                            - authz+tenant - LangGraph     - notifikasi
                            - LMS-owned   - MCP server    - eval batch
                            - enqueue     - eval runner
                                     \         |            /
                                      \        |           /
                                       [ redis ] (kecil, cache+queue)
                                              |
              ------------------ data plane (eksternal, gratis) ------------------
              [ Neon Postgres + pgvector ]        [ Supabase Storage ]
                 serverless, scale-to-zero            file materi kuliah
                 (DB utama, lepas dari VM)            (1 GB, S3-like)

   ---------------------- Observability plane ----------------------
   OTel Collector (di VM) -> Prometheus | Loki | Tempo | Grafana
                             (self-host di LAPTOP, atau di VM saat scale-up)
   Langfuse self-host (LLM traces, cost, eval datasets)

   ---------------------- Delivery plane ----------------------
   GitHub -> Actions (lint, test, gosec, trivy, promptfoo eval)
          -> GHCR (image multi-arch: amd64 + arm64)
          -> deploy via SSH: pull, health check, rollback
```

**Catatan penting soal RAM 1 GB:** stack di atas muat karena Go (~30 MB), FastAPI (~250 MB), Redis (~50 MB), Caddy (~20 MB) — dengan Postgres dan frontend sudah dipindah keluar. Tambahkan swap 2 GB. Saat butuh tenaga lebih (load test Minggu 6, batch embedding Minggu 8), **scale-up sementara** ke VM 4 GB pakai kredit $100, lalu turunkan lagi.

### 5.3 Fitur AI di LMS (yang punya alasan bisnis, bukan gimmick)

| Fitur | Teknik yang dilatih | Minggu |
|---|---|---|
| **Auto-summary & learning objectives** dari material yang diupload dosen | LLM API, structured output, prompt versioning, cost tracking; hasil disimpan sebagai `ai_artifacts`, bukan sebagai `learning_outcomes` | 7 |
| **Feature flag AI per tenant** | `tenant_settings`, deny-by-default, kemampuan mematikan seluruh AI per tenant | 7, melalui promosi Tier 2 yang dicatat tertulis |
| **"Tanya Materi"** — chat berbasis material yang telah dipublikasikan dengan sitasi halaman | `material_chunks`, pgvector, hybrid search, reranking, citation grounding, RLS | 8 |
| **Draft kuis dari materi** | Structured output + sitasi; disimpan sebagai `ai_artifacts.artifact_type = quiz_draft` dan wajib human review. Tidak membuat `quizzes`, `questions`, atau `question_banks` pada rilis pertama karena entity tersebut Tier 3 | 8–9 |
| **Eval harness + CI gate** | `ai_feedback`, golden dataset, ragas, promptfoo, regression test | 9 |
| **Study-plan agent** | LangGraph, tool calling, agentic RAG, checkpointing, guardrails; memakai data Tier 1 yang sudah tersedia | 10 |
| **MCP server `campus-lms`** | Mengekspos course offering, assignment/deadline, submission status, published grade, dan pencarian material sebagai tools dengan auth per-tenant | 10 |
| **Audit aksi agent** | `agent_actions`, human approval untuk tool berdampak, jejak immutable | 10, melalui promosi Tier 2 yang diprioritaskan |
| **Cost & quality dashboard** | `ai_interactions`, Langfuse + Grafana, token/cost/latency/fallback/cache metrics | 9–10 |
| **Kuota AI per tenant** | `ai_quotas`, request/token budget, hard limit atau degradasi terkontrol | 12, melalui promosi Tier 2 |
| **Grading assistant berbasis rubric** | Tetap bagian desain `domain-ai.md`, tetapi bukan deliverable wajib 12 minggu karena tabel `rubrics` berada di Tier 3. Tidak boleh memaksa pembuatan tabel Tier 3 hanya untuk mengejar demo | Tidak dijadwalkan sebagai implementasi wajib |

Perhatikan polanya: setiap fitur AI **mengonsumsi infrastruktur dan entity yang sudah tersedia lebih dulu**. AI tidak pernah menjadi source of truth akademik. Artefak AI yang berumur panjang wajib melalui status review; mahasiswa tidak boleh melihat konten AI sebelum `approved`. Konten AI seperti learning objectives atau draft kuis tidak boleh diam-diam membuat entity domain Tier 3.


---

## 6. Peta Infrastruktur Gratis 2026 (Tanpa Kartu)

> ⚠️ **Peringatan akurasi:** angka free tier berubah cepat. Semua nilai di bawah punya sumber bertanggal, tapi **kamu wajib cek halaman resmi vendor sebelum bergantung padanya.** Kolom "Kartu?" adalah kolom terpenting untukmu.

### 6.1 Compute / VPS — jalur utama: Azure for Students

**Ini pengganti Oracle, dan sebenarnya lebih baik untuk kasusmu.** Microsoft memberi mahasiswa terverifikasi **$100 kredit Azure berlaku 12 bulan, tanpa kartu kredit**, plus akses ke 25+ layanan always-free. Verifikasi memakai email kampus; kalau email `.ac.id`-mu tidak dikenali, Microsoft mengizinkan verifikasi lewat **upload kartu mahasiswa** atau lewat **GitHub Student Developer Pack** yang sudah terverifikasi.

| Opsi | Yang didapat | Kartu? | Catatan |
|---|---|---|---|
| **Azure for Students** ⭐ jalur utama | $100 kredit / 12 bulan + layanan populer gratis 12 bulan, termasuk **VM Linux B1s 750 jam/bulan** (1 vCPU, 1 GB RAM), 5 GB blob storage, Azure Functions 1 juta request/bulan | **Tidak** — email kampus / KTM / GitHub Student Pack | Akun ini punya **spending limit keras**: saat kredit habis, layanan berhenti — **tidak ada tagihan kejutan** karena tidak ada kartu tersimpan. Bisa diperbarui tiap tahun selama masih mahasiswa |
| **Cloudflare Pages** | Hosting frontend, bandwidth unlimited, 500 build/bulan | Tidak | Frontend Next.js pindah ke sini → beban VM jauh berkurang |
| **Cloudflare Workers** | 100.000 request/hari | Tidak | Edge function, webhook receiver |
| **Cloudflare Tunnel** (`cloudflared`) | Expose service dari laptop/VM ke internet dengan TLS, **tanpa membuka port apa pun** | Tidak | Penyelamat: demo publik bisa jalan bahkan dari laptop di kos |
| **GitHub Codespaces** | Kuota core-hours gratis per bulan (lebih besar dengan GitHub Pro dari Student Pack) | Tidak | Dev environment cloud — **solusi untuk laptop 8 GB RAM-mu**. Bisa jalankan Docker Compose penuh di sana |
| **GitHub Actions** | Gratis unlimited untuk **repo publik** | Tidak | Buat repo LMS publik — sekalian jadi portofolio |
| ~~Oracle Cloud Always Free~~ | ~~2 OCPU / 12 GB ARM~~ | **Ya — meminta kartu saat signup** | **Dicoret dari roadmap.** Kalau suatu saat kamu punya kartu yang jalan, ini masih opsi bagus (limit dipangkas 4/24 → 2/12 per 15 Juni 2026) |
| ~~AWS / GCP / Fly.io / Railway / DigitalOcean~~ | — | **Ya** | Dicoret. DigitalOcean juga sudah menghentikan program GitHub Student Pack (redemption tutup 31 Juli 2026) |

**Strategi memakai kredit $100 dengan cerdas** (ini sekaligus latihan cost engineering yang bisa kamu ceritakan):

| Periode | Ukuran VM | Alasan |
|---|---|---|
| Minggu 4–12 (default, 24/7) | **B1s kandidat** — gunakan hanya jika SKU benar-benar tersedia saat provisioning; jika tidak, pilih SKU B-series yang tersedia dan masih sesuai budget | Demo portofolio hidup dengan resource envelope sekecil mungkin; keputusan final dicatat dari hasil provisioning |
| Minggu 6 (load test k6) | Scale-up sementara ke **B2s** (2 vCPU / 4 GB) selama 2–3 hari | Butuh headroom untuk 500 VU. Lalu **turunkan lagi** |
| Minggu 8 (batch embedding) | Jalankan embedding **di laptop**, bukan di VM | CPU Ryzen 5-mu lebih kuat dari B1s. Hasilnya (vektor) yang di-push ke Neon |
| Kapan pun idle > 1 hari | `az vm deallocate` | Kredit berhenti terpakai (disk tetap kecil biayanya) |

Konfigurasi aktual Minggu 0 memakai budget **`campus-lms-guard` $10/bulan** dengan alert **50% / 80% / 100%**. Pengiriman email alert masih harus dibuktikan sebelum mekanisme ini boleh diklaim end-to-end. Ini tetap menjadi latihan cost engineering bersama spending limit dan right-sizing.

**Rencana B — kalau verifikasi Azure gagal:** jalankan seluruh stack di **laptop + Cloudflare Tunnel**. Demo tetap bisa diakses publik dengan TLS (selama laptop nyala). Ini sah untuk portofolio; sebutkan apa adanya di README.

**Rencana C — VPS lokal, bayar QRIS/GoPay/OVO** (kalau kamu bersedia keluar Rp 25–90rb/bulan):

| Provider | Harga awal | Metode bayar |
|---|---|---|
| Biznet Gio NEO Lite | ~Rp 59rb/bln | Transfer, VA, QRIS |
| IDCloudHost Cloud VPS | ~Rp 87rb/bln (top-up saldo, pay-as-you-go per jam) | **QRIS, GoPay, OVO, ShopeePay, VA, Alfamart/Indomaret** |
| CloudKilat / Dihostingin | mulai ~Rp 25rb/bln (spek kecil) | QRIS, GoPay, DANA, VA |

Keunggulan jalur ini: ditagih dalam rupiah, **tidak ada kartu yang kena flag transaksi luar negeri**, dan support berbahasa Indonesia. Kekurangannya: uptime dan kualitas jaringan bervariasi, dan sebagian tidak punya API/Terraform provider yang matang.

### 6.2 Database & storage

| Opsi | Free tier | Kartu? | Peran di roadmap |
|---|---|---|---|
| **Postgres self-host (Docker, di laptop)** ⭐ | Gratis | Tidak | **Tempat kamu belajar** — migrasi, index, EXPLAIN, backup, restore drill |
| **Neon** ⭐ | ~0.5 GB/project, hingga 10 project, 10 branch, scale-to-zero | Tidak | **Database produksi** (karena VM 1 GB tidak muat Postgres) + **database preview per-PR** di CI lewat branching copy-on-write — trik yang mengesankan di interview |
| **Supabase** | 500 MB DB, 1 GB storage, 5 GB bandwidth, 2 project, pause setelah 7 hari idle | Tidak | **Object storage** untuk materi kuliah + bahan pembanding di Minggu 3 |
| MinIO (Docker, lokal) | Gratis | Tidak | Belajar S3 API tanpa layanan eksternal |

### 6.3 LLM API gratis (strategi "stacking") — semuanya tanpa kartu

Ini kunci menjalankan seluruh fase AI dengan Rp 0 dan tanpa kartu. Karena tiap provider punya batas harian, kamu bangun **router sendiri** — dan kebetulan kamu sudah punya intuisinya dari setup combo 9router.

| Provider | Kuota free (laporan Jun 2026) | Kartu? | Pakai untuk |
|---|---|---|---|
| **Google Gemini API (AI Studio)** | Flash tier: ~1.500 req/hari, 10–15 RPM; konteks 1M | Tidak | Default. Long-context (analisis silabus panjang), summarization |
| **Groq** | ~30 RPM, ~1.000–14.400 req/hari (bervariasi antar sumber) | Tidak | Jalur latensi rendah: chat "Tanya Materi", agent step cepat |
| **Cerebras** | ~1 juta token/hari | Tidak | Batch: ingest dokumen, generate eval dataset |
| **OpenRouter** | ~20 RPM, 50 req/hari | Tidak | Safety net + eksperimen banyak model dengan satu API key. *(Upgrade ke 1.000 req/hari butuh top-up $10 — lewati saja, tidak perlu)* |
| **GitHub Models** | 10–15 RPM, 50–150 req/hari | Tidak | Akses model frontier untuk perbandingan kualitas |
| **Cloudflare Workers AI** | ~10K neuron/hari | Tidak | Embedding & task kecil di edge |
| **Mistral La Plateforme** | ~1 miliar token/bulan (2 RPM) | Bervariasi menurut sumber | Cadangan. Kalau diminta kartu, lewati |

**Embedding tanpa API sama sekali:** jalankan `sentence-transformers` (`multilingual-e5-small` atau `bge-m3`) di CPU laptopmu. Ryzen 5 7430U sanggup meng-embed ribuan chunk — lambat, tapi ini batch job. Ini menghilangkan ketergantungan kuota embedding **dan** menghemat kredit Azure.

### 6.4 Observability & pendukung

| Kebutuhan | Pilihan gratis | Kartu? |
|---|---|---|
| Metrics/logs/traces/dashboard | Prometheus + Loki + Tempo + Grafana OSS, self-host via Compose (di laptop, atau di VM saat scale-up) | Tidak |
| LLM tracing, cost, eval | **Langfuse** self-host (Docker) | Tidak |
| Container registry | GHCR (gratis, repo publik) | Tidak |
| Uptime monitoring | Uptime Kuma (self-host) atau UptimeRobot free | Tidak |
| Error tracking | GlitchTip (self-host, Sentry-compatible) | Tidak |
| APM tambahan | **New Relic gratis lewat GitHub Student Pack** | Tidak |
| Secrets | GitHub Actions Secrets + SOPS + age | Tidak |
| Email transaksional | Brevo free tier (300 email/hari) | Tidak |
| Notifikasi alert | Bot Telegram (gratis, instan) | Tidak |
| Domain | Namecheap `.me` gratis 1 tahun (Student Pack) · `.my.id` ~Rp 15–30rb/thn via registrar lokal (QRIS) · DuckDNS / Cloudflare Tunnel (gratis) | Tidak |
| IDE | JetBrains GoLand/PyCharm gratis (Student Pack) | Tidak |

### 6.5 Kenapa kartumu ditolak — dan apa yang bisa kamu lakukan

Ini bukan bagian dari roadmap teknis, tapi kamu berhak tahu penyebabnya supaya tidak membuang waktu.

**Penyebab paling umum kartu debit BCA Mastercard ditolak di merchant internasional:**

1. **Toggle transaksi internasional belum aktif.** Cek: **BCA Mobile → Akun Saya → Kontrol → Transaksi Internasional** (di sebagian versi aplikasi: `m-Admin → Atur Kartu Debit`). Sebagian besar orang berhenti di sini dan berhasil. Kamu bilang sudah aktif — lanjut ke poin berikutnya.
2. **Merchant tidak memakai 3D Secure (3DS/OTP).** Ini biang keladi yang sebenarnya. Kartu **debit** BCA umumnya hanya lolos di merchant yang menjalankan 3DS dengan OTP. Banyak platform cloud/SaaS internasional memakai alur *card-on-file* tanpa 3DS → otomatis ditolak, sebanyak apa pun saldomu. **Ini menjelaskan kenapa kamu ditolak "di mana pun".**
3. **Verifikasi pra-otorisasi $0 / $1.** Banyak penyedia cloud melakukan *authorization hold* kecil. Kartu debit sering gagal di langkah ini.
4. Kartu berlogo **GPN** hanya untuk domestik (tidak relevan kalau kartumu sudah Mastercard).

**Kalau suatu saat kamu butuh kartu yang benar-benar jalan** (bukan untuk roadmap ini — semua sudah bebas kartu — tapi mungkin nanti untuk beli domain atau API berbayar), yang secara resmi mendukung transaksi internasional online adalah kartu debit **bank digital** dengan toggle eksplisit. Contoh terdokumentasi: **blu by BCA Digital** menyediakan `bluDebit Card` Mastercard dengan kontrol transaksi terpisah — *Transaksi Domestik Online*, **Transaksi Internasional Online**, *Transaksi Internasional Offline*, *Contactless* — dan menyatakan kartunya dapat dipakai di e-commerce yang menerima Mastercard, domestik maupun internasional. Bank Jago dan Jenius juga sering disebut komunitas developer Indonesia sebagai yang paling lancar untuk keperluan ini.

> Catatan jujur: pernyataan tentang blu bersumber dari FAQ resminya; sebutan Jago/Jenius bersumber dari laporan komunitas (Reddit r/indotech), **bukan jaminan**. Keberhasilan tetap bergantung pada merchant dan alur 3DS-nya. Untuk Superbank/OVO/GoPay yang kamu miliki, **jangan berasumsi** bisa dipakai untuk merchant kartu internasional — verifikasi langsung di aplikasi masing-masing sebelum mengandalkannya. Karena itu, roadmap ini **tidak menggantungkan satu langkah pun** pada hal tersebut.


## 7. Naik Level Cara Pakai AI: dari Konsumen → Pembangun

Kamu sudah mahir **memakai** agent. Sekarang kamu harus bisa **membangunnya**. Tiga pergeseran pola pikir yang harus terjadi selama 12 minggu:

**Pergeseran 1 — Dari "prompt agent" ke "desain sistem yang memanggil LLM."**
Di opencode kamu yang jadi orkestrator (kamu memilih combo `plan` vs `heavy`). Di produk, **kode yang harus mengorkestrasi** — routing model, retry, fallback, timeout, budget token, cache — tanpa manusia di tengah. Setup combo-mu adalah *prototype mental* dari router yang akan kamu tulis di Minggu 7. Manfaatkan itu: kamu sudah tahu instingnya, tinggal jadikan kode.

**Pergeseran 2 — Dari "kelihatannya bagus" ke "terukur."**
Ini pembeda terbesar. Sekarang kamu menilai output agent dengan mata. Di industri, kamu harus punya angka: *faithfulness 0.87, context precision 0.79, p95 latency 2.3 s, $0.004/request, hallucination rate turun dari 12% ke 4%*. Eval harness (Minggu 9) adalah minggu paling berharga dalam roadmap ini untuk nilai jualmu.

**Pergeseran 3 — Dari "AI membantu saya coding" ke "AI adalah fitur produk yang saya operasikan."**
Recruiter tidak terkesan kamu memakai Claude Code. Mereka terkesan kalau kamu bisa menjawab: *"Bagaimana kamu mencegah prompt injection lewat PDF yang diupload dosen?"* — kamu akan bisa menjawabnya di Minggu 10.

### 7.1 Mode kerja yang dipilih: **agent-first dengan verifikasi berbasis bukti**

Kamu memilih agent mengerjakan semua yang bisa diotomasi secara lokal, dan manusia hanya menangani yang mustahil diotomasi. Roadmap ini mengikuti keputusan itu — dengan satu penyesuaian yang tidak bisa ditawar.

**Risiko yang harus kamu sadari (saya sampaikan sekali, lalu tidak lagi):**
Kalau agent menulis kode dan kamu hanya menerimanya, dalam 12 minggu kamu akan punya repo yang mengesankan dan pemahaman yang rapuh. Interview akan membongkarnya dalam 10 menit — pertanyaan seperti *"kenapa RLS dan bukan filter di aplikasi?"* atau *"komponen mana yang gagal duluan saat 10× beban?"* tidak bisa dijawab oleh orang yang tidak pernah bergulat dengan sistemnya.

**Mitigasinya — dan ini yang membuat mode agent-first tetap menghasilkan engineer:**
Beban belajarmu dipindahkan dari *mengetik kode* ke *memverifikasi, mengukur, dan menjelaskan*. Tiga mekanisme wajib:

1. **Bukti eksekusi, bukan klaim.** Agent dilarang menulis "sudah berjalan" tanpa melampirkan output perintah asli yang tersimpan di `docs/progress/evidence/`. Protokol lengkapnya ada di `agent/evidence-protocol.md`.
2. **Quiz verifikasi mingguan.** Tiap akhir minggu, agent menyusun 8–12 pertanyaan dari kode yang **benar-benar ada di repo**, dan kamu menjawabnya tanpa membuka kode. Skor < 70% = minggu itu diulang, bukan dilanjutkan.
3. **Explain-back 3 menit.** Rekam suara menjelaskan satu keputusan teknis minggu itu. Kalau tersendat, kamu belum paham — dan itu ketahuan sekarang, bukan saat interview.

**Pembagian kerja yang berlaku:**

| Kategori | Pelaku | Contoh |
|---|---|---|
| Implementasi lokal | **Agent** | Kode Go/Python, Dockerfile, Compose, YAML CI, skrip, migrasi, test, refactor, dokumentasi teknis |
| Eksekusi & pengukuran | **Agent** | Menjalankan test, load test, `EXPLAIN ANALYZE`, eval harness, mengumpulkan bukti |
| Penyusunan laporan | **Agent** | Draf `docs/progress/week-XX.md` lengkap dengan bukti dan daftar hal yang belum terverifikasi |
| **Klik portal & identitas** | **Manusia** | Azure Portal, verifikasi status mahasiswa, klaim Student Pack, DNS registrar, 2FA |
| **Keputusan & trade-off** | **Manusia** | Isi ADR, strategi multi-tenancy, prioritas fitur, definisi SLO |
| **Verifikasi pemahaman** | **Manusia** | Quiz mingguan, explain-back, review PR, sign-off laporan |
| **Kredensial** | **Manusia** | Semua API key & secret. Agent dilarang menyentuh `.env` asli |

Aturan operasional lengkap untuk agent ada di `campus-lms/AGENTS.md` dan folder `campus-lms/agent/`.

### 7.2 Tiga aturan anti-halusinasi yang berlaku di repo ini

1. **Nomor tanpa bukti = pelanggaran.** Setiap angka (latensi, coverage, skor eval, ukuran image) wajib punya file bukti berisi perintah + output mentah + timestamp.
2. **"Tidak tahu" adalah jawaban yang sah.** Agent wajib menulis bagian *Belum terverifikasi* di tiap laporan, bukan mengarang agar terlihat lengkap.
3. **DoD hanya boleh dicentang oleh manusia.** Agent boleh mengusulkan, kamu yang mencentang setelah melihat buktinya.

Sebutkan cara kerja ini di interview — kemampuan memverifikasi & men-debug kode hasil AI kini secara eksplisit jadi ekspektasi junior, dan kamu punya sistem terdokumentasi untuk itu. ([nucamp](https://www.nucamp.co/blog/the-junior-developer-hiring-crisis-in-2026-how-to-get-your-first-backend-job))

---

## 8. Roadmap Mingguan

### 8.0 Guardrail implementasi domain

Urutan tabel di bawah **mengikat roadmap**. Ini adalah rekonsiliasi langsung antara `docs/domain.md` §37–38 dan `docs/domain-ai.md` §5.

| Minggu | Entity yang boleh/direncanakan dibangun | Status |
|---|---|---|
| 3 | `tenants`, `users`, `auth_identities`, `memberships`, `membership_roles`, `audit_logs`, `academic_terms`, `courses`, `course_offerings`, `course_staff`, `enrollments` | Tier 1 |
| 4 | `auth_sessions` | Tier 1 |
| 5 | `modules`, `lessons`, `materials`, `files` | Tier 1 |
| 6 | `learning_activities`, `assignments`, `submissions`, `submission_versions`, `submission_files`, `grade_items`, `grades` | Tier 1 |
| 7 | `ai_interactions`, `ai_artifacts` | AI roadmap |
| 7 | `tenant_settings` | Tier 2; perlu promosi tertulis sebelum AI boleh diaktifkan per tenant, karena feature flag AI adalah requirement non-negotiable di `domain-ai.md` |
| 8 | `material_chunks` | AI roadmap |
| 9 | `ai_feedback` | AI roadmap |
| 10 | `agent_actions` | Tier 2; diprioritaskan oleh `domain.md` dan dipetakan ke Minggu 10 oleh `domain-ai.md`; promosi dicatat tertulis |
| 12 | `ai_quotas` | Tier 2; dipetakan ke Minggu 12 oleh `domain-ai.md`; promosi dicatat tertulis |

Entity Tier 2 lain (`attendance_sessions`, `attendance_credentials`, `attendance_records`, `activity_completions`, `announcements`, `grade_categories`) **tidak otomatis menjadi deliverable**. Entity tersebut hanya boleh ditambahkan bila waktu memungkinkan dan keputusan naik tier ditulis di laporan mingguan.

Entity Tier 3 — termasuk `quizzes`, `question_banks`, `rubrics`, `final_grades`, `integrations`, `academic_program_refs`, dan entity Tier 3 lain yang tercantum di `domain.md` — **tidak boleh muncul di migrasi rilis pertama**. Fitur roadmap harus disesuaikan dengan entity yang tersedia, bukan sebaliknya.

### Ikhtisar fase

| Fase | Minggu | Fokus |
|---|---|---|
| **0. Persiapan** | 0 (±4 jam) | Pengamanan akun cloud, klaim benefit, setup agent — **tugas manusia** |
| **A. Fondasi Operasional** | 1–3 | Linux/jaringan/Git, Docker, Postgres mendalam |
| **B. Produksi Nyata** | 4–6 | VPS + TLS + hardening, CI/CD, observability & performa |
| **C. AI Engineering** | 7–10 | LLM di produk, RAG, evaluation, agent + MCP |
| **D. Skala & Kesiapan Kerja** | 11–12 | Kubernetes/IaC, security, dokumentasi, portofolio |

---

### Deliverable yang berlaku di SETIAP minggu (tanpa kecuali)

Selain deliverable spesifik tiap minggu, tiga artefak berikut **wajib ada setiap minggu**. Ini yang mengubah mode agent-first dari "menumpuk kode" menjadi "membangun pemahaman":

| Artefak | Lokasi | Penyusun | Verifikator |
|---|---|---|---|
| **Laporan mingguan** — apa yang dikerjakan agent, apa yang dikerjakan manusia, keputusan yang diambil, konsep yang dipelajari, angka yang diukur, dan **daftar hal yang belum terverifikasi** | `docs/progress/week-XX.md` | Agent (draf) | Kamu (sign-off) |
| **Bukti eksekusi** — output mentah perintah, timestamp, commit SHA untuk setiap klaim terukur | `docs/progress/evidence/week-XX/` | Agent | Kamu (spot-check minimal 3) |
| **Quiz verifikasi** — 8–12 pertanyaan dari kode yang benar-benar ada di repo, dijawab tanpa membuka kode | `docs/progress/quiz/week-XX.md` | Agent (soal) | Kamu (jawab, target ≥ 70%) |

**Aturan gerbang:** laporan harus ditandatangani dan skor quiz harus ≥ 70%. DoD yang belum terpenuhi tetap dicatat sebagai gap terbuka; jika manusia memutuskan lanjut karena ada scope revision/carry-over yang eksplisit, status gap tersebut **tidak berubah menjadi selesai**. Template dan protokolnya ada di `campus-lms/agent/` dan `campus-lms/docs/progress/`.

---

### MINGGU 0 — Persiapan Akun & Agent (±4 jam, kerjakan sebelum Minggu 1)

**Tujuan:** Semua akun aman, semua benefit terklaim, agent siap bekerja — sebelum satu baris kode ditulis.

> **Ini satu-satunya blok kerja yang hampir seluruhnya MANUAL.** Bukan karena agent tidak mampu, tapi karena portal cloud, verifikasi identitas, dan kredensial memang tidak boleh dan tidak bisa didelegasikan. Checklist yang bisa dicentang ada di `campus-lms/docs/setup/azure-day-0.md`.

#### 0.1 Kunci pengaman akun Azure (manusia, ±45 menit)

| # | Langkah | Kenapa penting |
|---|---|---|
| 1 | **Subscriptions** → pastikan offer bernama **"Azure for Students"** dan **Spending limit: ON** | Selama ON, kamu **secara struktural tidak bisa ditagih**: kredit habis → layanan berhenti, dan tidak ada kartu tersimpan untuk ditagih. Jangan pernah terima tawaran *"Remove spending limit"* |
| 2 | **Cost Management → Budgets** → buat budget `campus-lms-guard`, **$10/bulan**, alert di 50/80/100% | Target operasionalmu mendekati $0. B1s hanya kandidat sampai provisioning Minggu 4 membuktikan availability; ambang budget ketat tetap dipakai sebagai deteksi dini |
| 3 | Nyalakan **MFA** di akun Microsoft | Akun ini memegang kredit + kredensial produksi. Perlakukan seperti akun kerja |
| 4 | **Subscriptions → Usage + quotas** → filter region **Southeast Asia**, cari `Standard BS Family vCPUs`, pastikan ≥ 2 | Langganan student kadang berkuota 0 di region tertentu. Ketahuan sekarang > kecewa di Minggu 4. Alternatif: Australia East / East Asia |
| 5 | Catat konvensi: region **Southeast Asia** (Singapore, terdekat dari Bali), resource group `rg-campuslms-prod`, tag wajib `project=campus-lms`, `env=prod`, `owner=<nama>` | Tanpa tag, laporan biaya tidak terbaca. Satu RG = satu perintah untuk reset total |

**Status aktual Minggu 0:** `✅` subscription Azure for Students aktif; `✅` spending limit ON; `✅` budget `$10/bulan` dengan threshold 50/80/100 dibuat; `✅` kuota BS Family dan konvensi Azure dicatat; `☐` pengiriman email budget alert belum terbukti; `☐` MFA tidak tercatat pada laporan Minggu 0 sehingga tidak boleh diasumsikan sudah aktif.

**Belum dikerjakan sekarang:** membuat VM (Minggu 4), Neon (Minggu 4), Cloudflare (Minggu 4), AKS (jangan pernah — akan melahap kredit).

#### 0.2 Klaim benefit Student Pack (manusia, ±30 menit)

- ✅ **Namecheap** — domain `.me` gratis 1 tahun sudah diklaim pada Minggu 0
- ☐ **GitHub Codespaces** — Student Developer Pack sudah aktif, tetapi aktivasi/kuota Codespaces tidak dibuktikan secara terpisah pada laporan Minggu 0 (Tidak Digunakan) 
- ☐ **JetBrains** GoLand + PyCharm — belum diklaim pada Minggu 0; tetap opsional bila diperlukan (Tidak Digunakan)
- ✅ **New Relic** — student benefit sudah berhasil diklaim pada Minggu 0

#### 0.3 Kredensial & lingkungan lokal (manusia, ±60 menit)

- ✅ SSH key ED25519 khusus Azure sudah dibuat; private key berpermission `600` dan tidak masuk repo
- ✅ Azure CLI sudah terpasang di WSL dan login device-code berhasil
- ✅ API key **Gemini AI Studio**, **Groq**, **Cerebras**, dan **OpenRouter** sudah dibuat
- ✅ `.env.example` sudah disalin ke `.env`, credential lokal diisi manusia, dan `.env` terbukti di-ignore Git
- ✅ `.wslconfig` sudah diatur dan diverifikasi: `memory=5GB`, `processors=8`, `swap=8GB`

#### 0.4 Siapkan agent (campuran, ±60 menit)

- ✅ Repo `campus-lms` sudah disiapkan di laptop/WSL
- ✅ `git init`, initial commit, dan push repo publik sudah selesai
- ✅ Placeholder `CHANGE_ME` sudah diganti dengan username GitHub
- ✅ `AGENTS.md` dan `agent/policy.md` sudah dibaca penuh oleh manusia
- ✅ Acceptance test agent sudah diulang setelah GNU Make dipasang; `make todo` berjalan dan ringkasan Minggu 1 tidak menambah task di luar repo

**Definition of Done Minggu 0**
- ☐ Spending limit sudah **ON** dan budget `$10/bulan` sudah dibuat, tetapi email budget alert **belum terverifikasi**
- ✅ Kuota `Standard BS Family vCPUs` di Southeast Asia dikonfirmasi tersedia (`0/4` pada laporan Minggu 0)
- ✅ Repo publik hidup di GitHub dan `make todo` berjalan di laptop/WSL
- ✅ Perbedaan `az vm stop` dan `az vm deallocate` sudah dijelaskan tanpa membuka catatan
- ✅ `docs/progress/week-00.md` terisi dan ditandatangani manusia

**Carry-over aktual dari laporan Minggu 0**
- ☐ Verifikasi budget alert benar-benar masuk ke email yang dibaca manusia.
- ☐ Formalisasi evidence interaktif yang belum tersimpan, terutama screenshot/artefak Azure yang laporan nyatakan tidak tersedia.
- ☐ Verifikasi/aktifkan MFA akun Microsoft; laporan Minggu 0 tidak menyebut hasil langkah ini.
- ☐ JetBrains Student Subscription tetap opsional; tidak perlu menjadi blocker bila memang tidak digunakan.
- ✅ SKU production **tidak dikunci** pada B1s di Minggu 0; availability final diverifikasi saat provisioning Minggu 4.

**Sinyal CV:** belum ada — ini persiapan. Tapi bagian "mengelola infrastruktur produksi dalam anggaran $0 dengan budget alerting" yang kamu klaim di Minggu 12 dimulai dari sini.

---

### MINGGU 1 — Linux, Jaringan, Git Profesional, & Fondasi Repo

**Tujuan:** Kamu nyaman di terminal Linux, paham apa yang sebenarnya terjadi saat request HTTP masuk, dan repo LMS berdiri dengan struktur profesional.

**Konsep wajib**
- Linux: filesystem hierarchy, permission (`chmod`/`chown`, kenapa 644 vs 755), user & group, **systemd** (`systemctl`, unit file, `journalctl`), proses & signal (SIGTERM vs SIGKILL — ini penting untuk graceful shutdown container), cron.
- Jaringan: model OSI seperlunya, TCP handshake, DNS (A/CNAME/TXT, TTL, propagasi), HTTP/1.1 vs HTTP/2, TLS handshake & rantai sertifikat, port & socket, NAT, firewall (`ufw`, iptables dasar), SSH (key pair, `~/.ssh/config`, agent forwarding, port forwarding).
- Tools yang harus jadi refleks: `ss -tulpn`, `curl -v`, `dig`, `traceroute`, `htop`, `journalctl -u`, `tail -f`, `grep`/`rg`, `jq`.
- Git profesional: trunk-based development, conventional commits, PR kecil, rebase vs merge, `git bisect`, tag semver, CODEOWNERS.

**Deliverable project**
1. WSL2 (Ubuntu 24.04) tertata; `.wslconfig` diatur; dotfiles disimpan di repo `dotfiles`.
2. Repo publik `campus-lms` (monorepo) dengan struktur:
   ```
   campus-lms/
   ├── apps/
   │   ├── api/          # Go
   │   ├── web/          # Next.js
   │   └── ai/           # Python FastAPI (kosong dulu)
   ├── deploy/           # compose, caddy, terraform nanti
   ├── docs/
   │   ├── adr/          # Architecture Decision Records
   │   └── runbook/
   ├── .github/workflows/
   ├── Makefile
   └── README.md
   ```
3. API Go minimal: `GET /healthz`, `GET /readyz`, log terstruktur (`log/slog` JSON), graceful shutdown (tangkap SIGTERM, drain 10 detik), config dari env (12-factor).
4. `docs/adr/0001-pilihan-stack.md` — ADR pertama: kenapa Go + Python + Postgres self-managed.
5. Domain model LMS v1 di `docs/domain.md`: entitas, relasi, batasan multi-tenant.

**Definition of Done**
- ✅ Explain-back perjalanan URL → DNS/TCP/TLS/HTTP sudah dikonfirmasi manusia pada laporan Minggu 1.
- ✅ `/healthz` mengembalikan HTTP 200 dan log JSON sesuai evidence Minggu 1.
- ✅ SIGTERM menghasilkan graceful shutdown dan exit code 0 sesuai evidence Minggu 1.
- ☐ Target **15 conventional commits** belum terpenuhi pada Minggu 1; laporan mencatat **7 commit**.
- ✅ SSH localhost key-only terverifikasi; password authentication dan keyboard-interactive authentication dimatikan.

**Carry-over aktual dari laporan Minggu 1**
- ☐ Target 15 conventional commits tetap terbuka; jangan menandainya selesai hanya karena minggu berikutnya sudah dimulai.
- ☐ Checkbox verifikasi quiz pada laporan Minggu 1 masih kosong walaupun skor tertulis **70/100**, tepat pada ambang gerbang. Jika itu hanya kelalaian pencatatan, manusia yang harus merekonsiliasinya di laporan; roadmap tidak mengubah sign-off secara otomatis.
- ✅ Explain-back Minggu 1 tercatat sudah direkam pada bagian verifikasi manusia.

**Alokasi 35 jam:** Linux 8 · Jaringan 8 · Git/workflow 4 · Kode Go + repo 12 · Dokumentasi/ADR 3

**Sinyal CV:** *"Structured logging, graceful shutdown, dan 12-factor config di service Go sejak commit pertama."*

---

### MINGGU 2 — Docker & Compose sampai ke Tulang

**Tujuan:** Kamu bisa menjelaskan apa itu container **tanpa memakai kata "VM ringan"**, dan menjalankan seluruh dev environment LMS dengan satu perintah.

**Konsep wajib**
- Isi container sebenarnya: namespaces (pid, net, mnt, uts, ipc, user), cgroups v2 (limit CPU/memori), union filesystem & layer, kenapa container ≠ VM.
- Image: Dockerfile instruction & layer caching, **multi-stage build**, distroless/`scratch` untuk Go, non-root user, `.dockerignore`, `HEALTHCHECK`, label OCI (Open Container Initiative), **multi-arch build (`buildx`, amd64 + arm64)** — tetap wajib meski VM Azure-mu x86, karena ini skill yang ditanyakan dan memberimu kebebasan pindah ke server ARM kapan saja.
- Runtime: `docker run` flags penting, bind mount vs named volume, port publishing, restart policy, resource limit (`--memory`, `--cpus`), `docker logs`/`exec`/`inspect`/`stats`.
- Networking: bridge default vs user-defined network, **DNS antar-service pakai nama service**, kenapa `localhost` di dalam container bukan host, `host.docker.internal`.
- Compose: services, depends_on + `condition: service_healthy`, profiles, env_file, override file (`docker-compose.override.yml` untuk dev), secrets.
- Kebersihan: `prune`, ukuran image, keamanan dasar (jangan jalankan root, jangan taruh secret di image layer).

**Deliverable project**
1. `apps/api/Dockerfile` — multi-stage, base `golang:1.2x-alpine` → runtime `gcr.io/distroless/static` atau `scratch`. **Target ukuran < 25 MB.** Non-root, HEALTHCHECK.
2. **Scope web dipindahkan ke Minggu 4** berdasarkan keputusan yang tercatat pada laporan Minggu 2. `apps/web/Dockerfile`, service web pada Compose, dan verifikasi frontend → API tidak lagi menjadi blocker Minggu 2.
3. `deploy/compose/docker-compose.yml` untuk core Minggu 2:
   - `core`: api, postgres, redis;
   - ketiganya mempunyai healthcheck/dependency yang dapat diverifikasi;
   - `storage` dipindahkan ke Minggu 5 saat flow `files/materials` dibangun;
   - `obs` dipindahkan ke Minggu 6 saat observability benar-benar diimplementasikan.
4. `docker-compose.override.yml` untuk development API: source mount + hot reload dengan Air. Scope `next dev` mengikuti pemindahan web ke Minggu 4.
5. `Makefile`: `make up`, `make down`, `make logs`, `make test`, build multi-arch, serta pembersihan aman (`make prune` tanpa volume; `make prune-hard` dengan konfirmasi eksplisit).
6. Buktikan pemahaman: `docs/notes/docker-internals.md` — eksperimen namespace/cgroup, resource limit/OOM, layer/image, dan build context yang mengukur objek yang benar.

**Definition of Done**
- ✅ `make up` → API, Postgres, dan Redis pada core sehat; DNS service `postgres` dan `redis` terverifikasi. Web tidak lagi menjadi blocker Minggu 2 karena dipindahkan eksplisit ke Minggu 4.
- ✅ Image API **14.6 MB** dan berjalan sebagai UID **65532** non-root berdasarkan evidence Minggu 2.
- ✅ Build multi-arch `linux/amd64` + `linux/arm64` berhasil dan manifest keduanya ter-push ke GHCR. Runtime arm64 masih belum diuji, tetapi itu bukan syarat DoD build/push ini.
- ☐ DoD **256 MB** belum terpenuhi persis: Compose tercatat 128MiB/512MiB/128MiB dan demonstrasi OOM dilakukan pada 32MiB. Jangan centang sampai target ini dipenuhi atau diubah melalui keputusan manusia yang eksplisit.
- ✅ ADR-0005 sudah memuat keputusan distroless + `/api -healthcheck` dan trade-off debugging tanpa shell.

**Carry-over aktual dari laporan Minggu 2**
- ☐ Jalankan ulang `make test` pada terminal manusia di SHA terkini; rerun agent tidak dapat menjadi bukti lulus karena sandbox memblokir listener `httptest`.
- ☐ Rekam explain-back 3 menit Minggu 2; checkbox ini masih kosong pada laporan.
- ☐ Selesaikan atau revisi secara eksplisit DoD demonstrasi **256 MiB**; bukti 32MiB tidak boleh dipakai sebagai pengganti diam-diam.
- ℹ️ Runtime `linux/arm64` masih belum diverifikasi. Ini gap non-blocking untuk DoD build/push multi-arch, tetapi tetap harus disebut sebagai belum terverifikasi bila dibahas.
- ✅ Scope web tidak dianggap hilang: dipindahkan ke Minggu 4. `storage` dipindahkan ke Minggu 5 dan `obs` ke Minggu 6.

**Alokasi:** Teori container 6 · Dockerfile/optimasi 10 · Compose & networking 10 · Eksperimen + tulisan 6 · Buffer 3

**Sinyal CV:** *"Multi-stage & multi-arch (amd64/arm64) Docker builds; image API Go 14.6 MB, distroless, non-root, dengan HEALTHCHECK tanpa shell."*

---

### MINGGU 3 — PostgreSQL Mendalam & Multi-Tenancy

**Tujuan:** Kamu lepas dari "Supabase sebagai kotak ajaib". Kamu mendesain skema Tier 1 yang benar, memisahkan course master dari course offering, menulis migrasi, membaca `EXPLAIN`, dan mengamankan isolasi antar-tenant sampai level database.

**Konsep wajib**
- Desain skema: normalisasi 3NF & kapan sengaja denormalisasi, tipe data (`uuid`, `timestamptz`, `numeric`, `jsonb`), constraint, soft delete vs hard delete.
- **`courses` ≠ `course_offerings`**: `courses` adalah master mata kuliah; `course_offerings` adalah pelaksanaan course pada academic term dan kelas tertentu. Data identitas akademik berasal dari SIAKAD dan tidak menjadi CRUD bebas LMS.
- **Identity & multi-role**: `users` global; hubungan user ke tenant melalui `memberships`; role tenant berada di `membership_roles`. Satu membership dapat memiliki lebih dari satu role aktif. Teaching Assistant tetap course-scoped melalui `course_staff`.
- **Strategi multi-tenancy**: shared schema + `tenant_id` + RLS vs schema-per-tenant vs database-per-tenant. Untuk project ini: shared schema + RLS.
- **RLS classification**: `tenants`, `users`, `auth_identities` adalah global dan tidak memakai tenant RLS biasa. `memberships`, `membership_roles`, `audit_logs`, `academic_terms`, `courses`, `course_offerings`, `course_staff`, dan `enrollments` tenant-scoped dan wajib RLS.
- **Composite foreign key**: setiap relasi antar-entity tenant wajib membawa `tenant_id`, sehingga referensi lintas tenant ditolak secara struktural oleh database.
- **Row Level Security**: `ENABLE ROW LEVEL SECURITY`, `USING`/`WITH CHECK`, `current_setting('app.tenant_id')`, dan `SET LOCAL` di dalam transaksi agar aman terhadap connection pooling.
- Index: B-tree, komposit & leftmost prefix, partial index, GIN, covering index, serta biaya index pada write.
- **`EXPLAIN (ANALYZE, BUFFERS)`**: seq scan vs index scan, nested loop vs hash join, rows estimate, `ANALYZE`.
- Transaksi: ACID, read committed, serializable, deadlock, `SELECT ... FOR UPDATE`, optimistic locking.
- Operasional: `pgxpool`, migrasi versioned, `pg_dump`/restore drill, `pg_stat_statements`, dan deteksi N+1.

**Deliverable project**
1. Migrasi Minggu 3 **hanya** untuk 11 tabel Tier 1:
   `tenants`, `users`, `auth_identities`, `memberships`, `membership_roles`, `audit_logs`, `academic_terms`, `courses`, `course_offerings`, `course_staff`, `enrollments`.
2. `memberships` tidak memiliki kolom role. Role aktif disimpan di `membership_roles` dengan histori pencabutan melalui `revoked_at`; aturan multi-role mengikuti amendemen `domain.md` §38-A1.
3. RLS aktif pada seluruh tabel tenant-scoped Minggu 3 + middleware Go yang menyetel `app.tenant_id` per-request di dalam transaksi.
4. Composite FK dan unique constraint tenant-aware diterapkan pada relasi Tier 1, termasuk:
   - `UNIQUE (tenant_id, user_id)` pada `memberships`;
   - `UNIQUE (course_offering_id, user_id)` pada `course_staff`;
   - `UNIQUE (course_offering_id, student_user_id)` pada `enrollments`;
   - FK `(tenant_id, student_user_id)` dari `enrollments` ke membership tenant yang sama;
   - seluruh child tenant-scoped tidak dapat menunjuk parent tenant lain.
5. Seeder dev/test yang merepresentasikan data sinkronisasi SIAKAD: **3 tenant, 50 dosen, 2.000 mahasiswa, 200 courses, sekitar 400 course_offerings, dan 20.000 enrollments**, plus supporting rows yang diperlukan oleh invariant pada `auth_identities`, `memberships`, `membership_roles`, `academic_terms`, dan `course_staff`. **Tidak ada 50.000 submission** pada Minggu 3 karena `submissions` baru dibangun Minggu 6.
6. Repository layer Go (`pgx`), tanpa ORM ajaib — SQL eksplisit (boleh `sqlc`). Endpoint/query membaca course offering dan peserta berdasarkan tenant context; tidak menyediakan jalur normal untuk mengubah fakta akademik SIAKAD-authoritative.
7. `docs/notes/query-tuning.md`: minimal 5 query yang dianalisis dengan `EXPLAIN (ANALYZE, BUFFERS)`, lengkap dengan before/after dan perubahan index/query. Angka yang ditulis harus berasal dari evidence aktual.
8. `docs/adr/0002-multi-tenancy.md` + `docs/notes/supabase-vs-self-managed.md`.
9. Skrip backup + **restore drill yang benar-benar dijalankan** (`deploy/scripts/backup.sh`), termasuk waktu restore aktual.

> Catatan: semua ini dikerjakan di Postgres yang kamu jalankan sendiri di Docker. Seeder adalah fixture pengembangan untuk mensimulasikan data akademik yang secara domain dimiliki SIAKAD; keberadaan seeder tidak mengubah ownership data.

**Definition of Done**
- [ ] Query sebagai Tenant A tidak dapat membaca atau memodifikasi row Tenant B, dibuktikan integration test RLS.
- [ ] Composite FK benar-benar menolak upaya referensi silang tenant.
- [ ] Satu membership dapat memegang dua role aktif yang berbeda; duplikasi role aktif yang sama ditolak; pencabutan menggunakan `revoked_at`, bukan `DELETE`.
- [ ] Tidak ada tabel Minggu 5/6 atau Tier 3 yang muncul prematur pada migrasi Minggu 3.
- [ ] Minimal satu query yang terbukti lambat/inefisien pada dataset seeder diperbaiki dan mempunyai `EXPLAIN` before/after; jangan mengarang threshold performa yang tidak muncul dari pengukuran.
- [ ] Endpoint daftar **course offering + peserta** bebas N+1, dibuktikan dengan hitungan query di test.
- [ ] `make db-restore` berhasil memulihkan backup ke database kosong dengan data lengkap.
- [ ] Kamu bisa menjelaskan read committed vs serializable memakai contoh dari skema LMS-mu sendiri.
- [ ] Kamu bisa menjelaskan kenapa `courses` tanpa `course_offerings` akan merusak histori saat masuk semester berikutnya.

**Alokasi:** Teori & desain skema 8 · RLS + composite FK + multi-role 8 · Index & EXPLAIN 8 · Repository + integration test 8 · Backup/dokumen 3

**Sinyal CV:** *"Mendesain Postgres multi-tenant dengan RLS + composite foreign key, memisahkan course master dari course offering, dan men-tuning query dengan bukti `EXPLAIN ANALYZE` before/after."*
---

### MINGGU 4 — Cloud Produksi: Deploy, TLS, Hardening, & Auth Session

**Tujuan:** LMS-mu hidup di internet, di server yang kamu amankan sendiri, dengan HTTPS dan sesi autentikasi yang dapat dirotasi/dicabut — semuanya **tanpa kartu**.

> **Langkah 0 (kerjakan hari Sabtu sebelum minggu ini dimulai, karena butuh waktu tunggu):**
> 1. Daftar **GitHub Student Developer Pack** di `education.github.com/pack` — verifikasi dengan email `.ac.id` atau foto KTM + bukti keaktifan.
> 2. Aktifkan **Azure for Students** — login dengan email kampus atau jalur verifikasi Student Pack.
> 3. Kalau keduanya gagal, jalankan Rencana B (laptop + Cloudflare Tunnel) atau Rencana C (VPS lokal bayar QRIS) dari §6.1.

**Konsep wajib**
- **Auth session lifecycle**: access token berumur pendek, refresh token disimpan sebagai hash, refresh rotation, revocation, reuse detection, dan hubungan `rotated_from`.
- `auth_sessions` adalah tabel **global**, bukan tenant-scoped; tenant authorization tetap diperiksa setelah user memilih/mengakses tenant.
- Provisioning Azure: Resource Group, VNet + Subnet, NSG, Public IP, VM Linux B-series, `az` CLI. `Standard_B1s` adalah kandidat, bukan SKU yang sudah dikunci; availability diverifikasi saat provisioning.
- Firewall dobel: NSG dan `ufw`/iptables.
- Cost management: budget alert, `deallocate` vs `stop`, disk tetap memiliki biaya.
- Hidup dengan resource VM kecil: bila SKU final B1s maka envelope-nya 1 GB RAM; sesuaikan swap, memory limit container, dan observasi OOM berdasarkan SKU yang benar-benar diprovision.
- Hardening: user non-root + sudo, SSH key-only, `PermitRootLogin no`, `PasswordAuthentication no`, `ufw`, `fail2ban`, unattended-upgrades, NTP.
- Caddy/TLS, DNS, deployment dari GHCR, rollback, secret management, runbook.

**Deliverable project**
1. Migrasi Tier 1 `auth_sessions` sesuai `domain.md` §38-A8:
   `id`, `user_id`, `refresh_token_hash`, `issued_at`, `expires_at`, `rotated_from`, `revoked_at`, `revoked_reason`, `user_agent`, `ip_address`, `last_seen_at`.
2. Implementasi refresh-token rotation + revocation. Pemakaian ulang token yang sudah dirotasi diperlakukan sebagai indikasi pencurian dan mencabut seluruh sesi user sesuai aturan domain.
3. Provision VM Azure B-series yang benar-benar tersedia di subscription/region; gunakan B1s bila tersedia dan sesuai budget. Catat SKU final sebagai keputusan aktual, lalu harden VM dan buat swap sesuai resource envelope.
4. Database produksi pindah ke Neon; jalankan migrasi yang sama dan aktifkan pgvector. Postgres Docker tetap untuk dev/test.
5. **Selesaikan carry-over web dari Minggu 2:** bangun/verifikasi frontend Next.js, `apps/web/Dockerfile` standalone non-root untuk jalur container lokal yang semula direncanakan, verifikasi frontend → API tanpa hard-coded container IP, lalu deploy frontend ke Cloudflare Pages; VM hanya melayani API/service yang diperlukan.
6. `campus-lms` live di HTTPS dengan sertifikat valid.
7. `deploy/compose/docker-compose.prod.yml` — image pin digest, restart policy, memory limit, log rotation.
8. `deploy/caddy/Caddyfile` — TLS, security headers, rate limit dasar, kompresi.
9. `deploy/scripts/deploy.sh` — pull image, health check, rollback otomatis jika readiness gagal.
10. Backup terjadwal: `pg_dump` dari Neon → storage tujuan yang telah dipilih, retensi 7 hari, notifikasi bila gagal.
11. `docs/runbook/incident.md` + `docs/adr/0002b-arsitektur-hemat-biaya.md`.

**Definition of Done**
- [ ] Refresh token tersimpan sebagai hash, bukan plaintext.
- [ ] Rotation menghasilkan session baru yang menunjuk `rotated_from`; token lama tidak dapat dipakai kembali.
- [ ] Test reuse token membuktikan seluruh sesi user dicabut sesuai rule.
- [ ] Situs dapat diakses publik lewat HTTPS dari jaringan di luar laptopmu.
- [ ] `nmap` dari luar hanya menampilkan port yang memang diizinkan; konfigurasi NSG dan firewall OS sama-sama terbukti.
- [ ] Kamu sengaja mematikan service/menyebabkan failure yang aman lalu memulihkannya mengikuti runbook dan mencatat MTTR aktual.
- [ ] Deploy versi baru tidak menghasilkan request gagal pada loop pengujian yang digunakan.
- [ ] Rollback dan backup/restore diuji; angka waktu yang dicantumkan di laporan berasal dari evidence.
- [ ] Budget alert Azure aktif dan proyeksi biaya dicatat dari portal/CLI, bukan asumsi.

**Alokasi:** Auth session 5 · Azure/provisioning 5 · Hardening + tuning 1 GB 7 · Neon + Pages 5 · Caddy/TLS/DNS 4 · Deploy/rollback 4 · Backup/runbook/ADR 5

**Sinyal CV:** *"Mengoperasikan LMS multi-tenant di cloud dengan refresh-token rotation/revocation, TLS, hardening jaringan, health-checked deployment, backup/restore, dan rollback yang terukur."*
---

### MINGGU 5 — CI/CD, Testing, Supply Chain, & Content Core

**Tujuan:** Setiap perubahan otomatis teruji dan terkirim dengan aman, sambil menambahkan struktur konten Tier 1 (`modules`, `lessons`, `materials`, `files`) tanpa merusak boundary tenant.

**Konsep wajib**
- Relasi `course_offering -> modules -> lessons -> materials`; satu module tidak dipaksa sama dengan satu minggu.
- `files` menyimpan metadata file dan tenant ownership; binary object berada di object storage, bukan row domain.
- `materials` mendukung text/file/link/video/audio/embed/learning package/external tool secara desain, tetapi implementasi pertama tidak perlu membuat entity Tier 3.
- Composite FK + RLS pada seluruh tabel tenant-scoped baru; file access wajib memvalidasi tenant.
- N+1 pada tree module/lesson/material dan cara memperbaikinya.
- Piramida test, Go testing, testcontainers, contract/e2e secukupnya.
- GitHub Actions, lint/type-check, supply-chain scanning, SBOM, action pinning, expand-contract migration.

**Deliverable project**
1. Migrasi Tier 1 Minggu 5: `modules`, `lessons`, `materials`, `files`.
2. RLS + composite FK tenant-aware untuk keempat tabel; `files.tenant_id` wajib eksplisit.
3. Repository/API minimal untuk:
   - course staff mengelola module/lesson/material pada `course_offering` yang berhak diakses;
   - mahasiswa hanya membaca konten yang tersedia/published pada offering tempat dirinya active enrolled;
   - file access menolak cross-tenant reference.
4. Object-storage flow: row `files` menyimpan metadata (`storage_key`, filename, MIME, size, checksum, scan status); binary tidak disimpan di Postgres. Aktifkan Compose profile `storage` dengan MinIO pada minggu ini sebagai realisasi scope yang ditunda dari Minggu 2.
5. `.github/workflows/ci.yml`: lint → unit → integration (testcontainers) → build multi-arch → Trivy scan → push GHCR.
6. `.github/workflows/cd.yml`: deploy setelah gate yang dipilih, jalankan migrasi, smoke test, rollback otomatis jika gagal.
7. Coverage bermakna di package domain/service; angka coverage dicatat dari output nyata.
8. Branch protection + required CI + review manual.
9. `docs/notes/testing-strategy.md` dan demo expand-contract pada satu perubahan aman.

**Definition of Done**
- [ ] Migrations Minggu 5 hanya menambahkan `modules`, `lessons`, `materials`, `files`; tidak ada tabel assessment/quiz/rubric yang muncul prematur.
- [ ] Integration test membuktikan Tenant A tidak dapat membaca material/file Tenant B meskipun mengetahui ID-nya.
- [ ] Endpoint tree course offering → modules → lessons → materials bebas N+1 berdasarkan query-count evidence.
- [ ] Upload/download metadata file menjaga tenant ownership dan tidak menyimpan binary di row domain.
- [ ] Push ke branch menjalankan CI; merge/deploy mengikuti gate yang ditetapkan tanpa langkah manual yang tidak terdokumentasi.
- [ ] Integration test memakai Postgres asli via testcontainers.
- [ ] Trivy/gosec/govulncheck dijalankan dan setiap pengecualian HIGH/CRITICAL, bila ada, dijustifikasi tertulis.
- [ ] PR yang sengaja merusak test ditahan oleh branch protection.

**Alokasi:** Domain content + object storage 8 · Test 8 · Actions/CI 8 · CD & migrasi aman 5 · Supply-chain security 4 · Dokumentasi 2

**Sinyal CV:** *"Membangun content core tenant-scoped (`modules`, `lessons`, `materials`, `files`) dengan RLS/composite FK dan object storage, lalu mengamankannya lewat CI/CD, testcontainers, dan supply-chain scanning."*
---

### MINGGU 6 — Observability, Performa, & Learning/Assessment Core

**Tujuan:** Kamu menambahkan alur assignment → submission version → grade Tier 1, lalu mengukur dan mengoptimalkan sistemnya dengan observability yang nyata.

**Konsep wajib**
- `learning_activities` sebagai **supertype** agar subtype lain dapat ditambahkan nanti tanpa membuat tabel Tier 3 sekarang.
- `assignments` sebagai subtype yang diimplementasikan pada rilis pertama.
- Submission logical record vs `submission_versions` immutable; resubmit membuat version baru, bukan overwrite.
- `submission_files` tenant-scoped dan wajib `tenant_id` eksplisit.
- `grade_items` dan `grades`: raw/max/normalized score, status, actor/timestamp, unique `(grade_item_id, enrollment_id)`.
- Amendemen `grade_items` §38-A5: gunakan typed FK ketika source type benar-benar tersedia. Pada Minggu 6, implementasikan source assignment dan manual tanpa membuat placeholder table `quizzes`/attendance. Kolom FK untuk source baru ditambahkan ketika tier terkait benar-benar dipromosikan.
- Tiga pilar observability + korelasi log/metric/trace, OpenTelemetry, Prometheus, Grafana, Loki, Tempo.
- SLI/SLO/error budget, profiling `pprof`, k6, Redis cache-aside, stampede protection, keyset pagination.

**Deliverable project**
1. Migrasi Tier 1 Minggu 6:
   `learning_activities`, `assignments`, `submissions`, `submission_versions`, `submission_files`, `grade_items`, `grades`.
2. RLS + composite FK pada seluruh tabel baru. `submission_files` menyimpan `tenant_id` eksplisit.
3. Alur assignment/submission:
   - activity dan assignment terkait tenant/offering yang sama;
   - submission dimiliki enrollment yang sah;
   - resubmission menambah `submission_versions`;
   - version lama immutable;
   - file submission menunjuk metadata `files` pada tenant yang sama.
4. Gradebook minimum:
   - `grade_items` mendukung assignment source yang sudah ada dan manual item;
   - `grades` mempunyai unique `(grade_item_id, enrollment_id)`;
   - TA hanya dapat menghasilkan/mengubah draft score sesuai authorization service;
   - **tidak membuat `final_grades`, `rubrics`, `quizzes`, atau `grade_categories`** pada Tier 1 Minggu 6.
5. Instrumentasi OTel penuh di `api-go`: trace web/api/postgres, `trace_id` pada log.
6. Compose profile `obs`: OTel Collector + Prometheus + Loki + Tempo + Grafana.
7. Dashboard RED, Postgres/Redis, resource host; alert error rate/latency/disk/backup.
8. `docs/slo.md`.
9. Load test k6 pada katalog course offering, material, dan alur submission; laporan before/after di `docs/notes/load-test-1.md`.
10. Caching Redis untuk katalog **course offering** + keyset pagination pada daftar yang volumenya cukup besar.

**Definition of Done**
- [ ] Resubmit membuat `submission_versions` baru dan test membuktikan version lama tidak berubah.
- [ ] Cross-tenant submission/file/grade reference ditolak oleh RLS/composite FK.
- [ ] Tidak ada tabel Tier 3 yang dibuat untuk “melengkapi” gradebook.
- [ ] Dari request lambat di Grafana kamu dapat menelusuri trace dan log terkait.
- [ ] Minimal satu bottleneck asli ditemukan lewat trace/pprof/EXPLAIN, diperbaiki, dan memiliki evidence before/after.
- [ ] Alert benar-benar diterima saat kondisi test dipicu.
- [ ] Load test dan optimasi menghasilkan angka yang dicatat apa adanya; klaim CV baru boleh memakai angka tersebut setelah evidence ada.
- [ ] Kamu dapat menjelaskan kenapa `learning_activities` dipakai sebagai supertype dan kenapa `submission_versions` immutable.

**Alokasi:** Domain assessment + test 10 · OTel/instrumentasi 7 · Observability stack 6 · Load test/profiling 7 · Optimasi/cache/pagination 5

> 💡 Jika perlu scale-up sementara untuk load test, catat biaya aktual dan turunkan kembali setelah eksperimen selesai.

**Sinyal CV:** *"Membangun assignment/submission versioning dan gradebook minimum tenant-scoped, lalu menginstrumentasi OpenTelemetry dan mengoptimalkan bottleneck berdasarkan trace/load-test evidence."*

> 🎯 **Checkpoint tengah jalan.** Setelah Minggu 6, baseline backend yang dibuktikan mencakup Docker, cloud deployment, SQL/RLS/composite FK, auth session, testing/CI/CD, content core, submission versioning, dan observability.
---

### MINGGU 7 — LLM dalam Produk: Fondasi AI Engineering

**Tujuan:** LMS punya fitur AI pertama di produksi dengan structured output, routing/fallback, cost tracking, dan human review — tanpa menjadikan AI source of truth.

> **Gate domain sebelum AI diaktifkan:** `domain-ai.md` mewajibkan fitur AI dapat dimatikan per tenant melalui `tenant_settings`. Karena `tenant_settings` berada di Tier 2, roadmap ini hanya boleh mempromosikannya setelah keputusan tertulis di laporan Minggu 7 sesuai `domain.md` §37.4. Jika promosi tidak dilakukan, fitur AI tetap disabled dan DoD “AI di produksi” belum boleh dicentang.

**Konsep wajib**
- Tokenisasi/context window, non-determinisme, halusinasi.
- Structured output + Pydantic/JSON Schema; prompt versioning.
- Timeout, retry, jitter, circuit breaker, idempotency, fallback, SSE.
- Cost engineering: token count, cache, routing, budget.
- Async/job queue.
- Privacy: PII redaction; data mahasiswa tidak dikirim tanpa alasan terdokumentasi.
- `ai_artifacts` lifecycle `draft -> approved -> published` atau `rejected`; published hanya boleh jika approved.
- `ai_interactions` menyimpan metadata biaya/latensi/model, **bukan isi prompt/jawaban**.
- Dokumen/material yang diproses model adalah data tak tepercaya, bukan instruksi.

**Deliverable project**
1. `apps/ai/` FastAPI service: schema Pydantic, Dockerfile multi-arch, `/healthz`, OTel, log terstruktur.
2. Migrasi AI: `ai_interactions`, `ai_artifacts`, keduanya tenant-scoped, RLS, dan mengikuti composite FK pattern. Terapkan invariant:
   - `UNIQUE (tenant_id, subject_type, subject_id, artifact_type, prompt_version)` pada `ai_artifacts`;
   - `CHECK (published_at IS NULL OR review_status = 'approved')`;
   - token dan `cost_estimate` pada `ai_interactions` tidak boleh negatif.
3. Promosi tertulis `tenant_settings` dari Tier 2 lalu implementasi feature flag AI per tenant. Bila tidak dipromosikan, AI tidak boleh diaktifkan.
4. `llm-router` multi-provider dengan retry/fallback/circuit breaker, token/cost accounting, dan cache.
5. **Auto-Summary + Learning Objectives Material**:
   - input berasal dari `materials`/`files` tenant yang sah;
   - model menghasilkan envelope tervalidasi `{summary, learning_objectives[], keywords[], estimated_read_minutes}`;
   - persist sebagai **dua** `ai_artifacts` untuk subject material yang sama: `artifact_type = summary` dan `artifact_type = learning_objectives`;
   - keduanya menyimpan `model`, `provider`, `prompt_version`, `input_hash`, `review_status = draft`;
   - tidak membuat `learning_outcomes` Tier 3.
6. UI review course staff: approve/reject. `published_at` hanya dapat terisi jika `review_status = approved`. Mahasiswa tidak dapat melihat draft/rejected artifact. Jika `input_hash` tidak lagi cocok dengan source material, artifact diperlakukan stale dan tidak boleh ditampilkan sebagai representasi materi terkini.
7. Setiap model call menulis satu row `ai_interactions` dengan feature, provider/model, token, cost estimate, latency, cache status, fallback, trace.
8. Prompt versioned di `apps/ai/prompts/*.md`.
9. `docs/adr/0003-llm-routing-dan-budget.md`.

**Definition of Done**
- [ ] Feature flag tenant dapat menonaktifkan AI; ketika off, endpoint AI menolak operasi tanpa memanggil provider.
- [ ] Auto-summary berjalan pada environment produksi hanya setelah `tenant_settings` promotion dicatat.
- [ ] Mahasiswa tidak dapat membaca `ai_artifacts` sebelum approved.
- [ ] Summary dan learning objectives tersimpan sebagai artifact type yang berbeda; tidak ada `learning_outcomes` Tier 3 yang dibuat.
- [ ] Database constraint menolak `published_at` jika review status bukan `approved`.
- [ ] Perubahan source material membuat `input_hash` lama terdeteksi stale sehingga artifact lama tidak ditampilkan sebagai representasi materi terkini.
- [ ] Matikan satu provider → fallback bekerja dan metadata fallback tercatat di `ai_interactions`.
- [ ] Setiap model call mempunyai token in/out, cost estimate, latency, status, trace metadata.
- [ ] Isi prompt/jawaban tidak disimpan di `ai_interactions`.
- [ ] Output JSON tervalidasi; kegagalan schema ditangani tanpa menyimpan artifact seolah-olah valid.

**Alokasi:** FastAPI + schema/RLS 8 · Router/reliability 9 · Auto-summary + artifact review 9 · Tenant feature flag 3 · Streaming/UI 3 · Dokumentasi 3

**Sinyal CV:** *"Membangun LLM gateway multi-provider dengan fallback/cost tracking dan lifecycle `ai_artifacts` human-reviewed; AI dapat dimatikan per tenant dan tidak pernah menjadi source of truth akademik."*
---

### MINGGU 8 — RAG Produksi: "Tanya Materi"

**Tujuan:** Mahasiswa bertanya hanya pada material yang memang boleh diakses dari course offering tempat dirinya enrolled, lalu mendapat jawaban bersitasi tanpa retrieval lintas tenant.

**Konsep wajib**
- Ingest → parse → chunk → embed → index → retrieve → rerank → augment → generate → cite.
- Parsing PDF/PPTX/DOCX, struktur dan metadata halaman/section.
- Chunking strategy dan content hash.
- Embedding multilingual, cosine similarity, dimensi/model.
- pgvector HNSW vs IVFFlat, parameter index, metadata filter, RLS.
- Hybrid search BM25/full-text + vector + RRF, reranking.
- Grounding: sitasi `material_id` + page; refusal bila konteks tidak memadai.
- Incremental re-embed; material unpublish/delete tidak boleh tetap retrievable.
- Retrieval authorization: tenant sama, course offering yang boleh diakses user, material sudah published.

**Deliverable project**
1. Migrasi `material_chunks` sesuai `domain-ai.md`: `tenant_id`, `material_id`, `course_offering_id`, `chunk_index`, `content`, `token_count`, `page`, `section`, `content_hash`, `embedding_model`, `embedding`, `created_at`; unique `(material_id, chunk_index, embedding_model)`.
2. RLS + composite FK tenant-aware untuk `material_chunks`; HNSW index dan parameter eksperimennya dicatat di ADR/notes.
3. Pipeline ingest background: material/file → parse → chunk → embed → `material_chunks`, idempotent dan incremental berdasarkan `content_hash`.
4. Retrieval hybrid + rerank hanya atas chunks yang satu tenant, satu offering yang berhak diakses user, dan berasal dari material published.
5. Endpoint `POST /api/course-offerings/{id}/ask` (streaming) + UI sitasi yang menunjuk `material_id` dan page.
6. Guardrail refusal jika retrieval di bawah threshold yang dipilih berdasarkan eksperimen.
7. **Draft kuis** dari material disimpan sebagai `ai_artifacts.artifact_type = quiz_draft`, tetap `draft` sampai course staff review. **Tidak membuat `question_banks`, `questions`, `quizzes`, atau tabel quiz Tier 3.**
8. `docs/notes/rag-experiments.md`: 3 strategi chunking × 2 embedding model × vector-only/hybrid × rerank/no-rerank pada test set yang sama.

**Definition of Done**
- [ ] `material_chunks` mempunyai tenant ownership eksplisit dan retrieval cross-tenant ditolak oleh test.
- [ ] User yang tidak enrolled/tidak menjadi course staff tidak dapat memakai endpoint ask untuk offering tersebut.
- [ ] Material yang di-unpublish tidak lagi muncul di retrieval.
- [ ] Setiap jawaban yang diberikan mempunyai sitasi yang dapat diverifikasi ke material/page; pertanyaan di luar materi ditolak.
- [ ] Eksperimen chunking/embedding/retrieval menghasilkan angka dari run nyata, bukan angka contoh.
- [ ] Draft kuis hanya tersimpan sebagai AI artifact dan tidak menyebabkan tabel Tier 3 muncul di migrasi.
- [ ] Latensi dan kualitas dicatat dari evidence aktual; threshold yang diklaim di CV hanya boleh berasal dari hasil tersebut.

**Alokasi:** Parsing & ingest 9 · pgvector/RLS 8 · Hybrid + rerank 7 · Authz/UI/sitasi 5 · Draft kuis artifact 2 · Eksperimen/dokumentasi 4

**Sinyal CV:** *"Membangun RAG multi-tenant di pgvector dengan hybrid retrieval, reranking, authorization berbasis course offering, dan sitasi material/page; konfigurasi dipilih dari eksperimen terukur."*
---

### MINGGU 9 — Evaluation, LLM Observability, & Quality Gates

> **Ini minggu pembeda:** fitur AI yang tidak diukur tidak boleh dianggap selesai.

**Tujuan:** Kualitas AI terukur, feedback manusia masuk ke loop evaluasi, dan regresi diblokir otomatis oleh CI.

**Konsep wajib**
- Eval vs assertion deterministic.
- Offline/online eval, reference-based/reference-free, LLM-as-judge.
- RAG metrics: faithfulness, answer relevancy, context precision/recall; retrieval hit rate/MRR/nDCG.
- Golden dataset yang diverifikasi manusia.
- Bias judge dan kalibrasi terhadap label manusia.
- promptfoo quality gate.
- Langfuse tracing.
- `ai_feedback` sebagai feedback user terhadap `ai_interactions` dan sumber kandidat golden cases.
- Metrik operasional AI: cost/request, cost/tenant, p95 latency, cache hit, refusal, fallback/error.

**Deliverable project**
1. Migrasi `ai_feedback` (`tenant_id`, `user_id`, `ai_interaction_id`, `rating`, `comment`, `created_at`) + RLS/composite tenant consistency.
2. UI/API feedback `helpful | not_helpful | wrong | unsafe`; kasus buruk dapat diekspor/dipilih sebagai kandidat golden dataset, tetapi verifikasi manusia tetap wajib.
3. `evals/` — golden dataset ≥ 80 kasus untuk summary, tanya-materi, dan **draft kuis AI artifact**; JSONL versioned di git.
4. Eval harness Python: ragas + judge kustom. `make eval`.
5. Langfuse self-hosted/terpasang pada environment yang mampu menampungnya; seluruh model call dapat dikorelasikan dengan `trace_id` dan `ai_interactions`.
6. Workflow `ai-eval.yml` untuk PR yang menyentuh AI/prompt; threshold gate ditetapkan dan dicatat dari baseline eval yang benar-benar dijalankan.
7. Dashboard AI Ops: token, cost estimate, latency, error/fallback, cache.
8. `docs/notes/eval-report-v1.md` + eksperimen perubahan yang dibandingkan berdasarkan angka.

**Definition of Done**
- [ ] Feedback pada interaction Tenant A tidak dapat direferensikan dari Tenant B.
- [ ] Kasus `wrong`/`unsafe` dapat ditelusuri ke interaction/trace dan masuk kandidat golden dataset.
- [ ] `make eval` menghasilkan laporan metrik dari dataset versioned.
- [ ] PR yang sengaja memperburuk kasus kritis ditahan oleh quality gate yang telah dibaseline.
- [ ] Langfuse/trace memperlihatkan alur model call beserta token/cost/latency metadata.
- [ ] Minimal satu perubahan AI dipilih berdasarkan perbandingan metrik before/after.
- [ ] Kamu dapat menjelaskan kelemahan LLM-as-judge dan proses kalibrasinya.

**Alokasi:** `ai_feedback` + integration 4 · Golden dataset 7 · Eval harness 8 · Langfuse/dashboard 6 · CI gate 5 · Eksperimen/laporan 5

**Sinyal CV:** *"Membangun eval harness dan feedback loop (`ai_feedback`) untuk RAG/LLM, dengan golden dataset versioned, trace/cost observability, dan quality gate CI berbasis baseline terukur."*
---

### MINGGU 10 — Agent, Tool Calling, MCP, & Guardrails

**Tujuan:** Kamu membangun agent multi-langkah yang hanya memakai data dan tool yang memang tersedia di scope rilis pertama, dengan audit `agent_actions`, approval, dan proteksi tenant.

> `agent_actions` adalah Tier 2 yang secara eksplisit diprioritaskan di `domain.md` dan dipetakan ke Minggu 10 di `domain-ai.md`. Implementasinya harus didahului keputusan promosi tertulis di laporan Minggu 10.

**Konsep wajib**
- Agent loop, state, stop condition, step/token/time budget.
- Tool design yang ketat, idempotency, error semantics.
- LangGraph state/edge/checkpoint/human interrupt.
- Agentic RAG.
- Memory scoped per user/course offering.
- MCP host/client/server, tools/resources/prompts, transport, authz.
- Prompt injection dari material sebagai data tak tepercaya.
- Allow-list tool per role, least privilege, PII redaction, output filtering.
- `agent_actions` sebagai audit immutable; action berdampak tidak boleh `executed` tanpa approval.

**Deliverable project**
1. Promosi tertulis + migrasi `agent_actions` sesuai `domain-ai.md`: tenant/user/session/step/tool/arguments-redacted/result-summary/status/requires-approval/approval metadata/token-cost/trace/timestamps; RLS + composite tenant consistency; immutable.
2. **Study-Plan Agent** memakai entity yang sudah tersedia:
   - `get_enrollment_context`
   - `list_upcoming_assignments`
   - `list_submission_status`
   - `list_published_grades`
   - `search_course_materials`
   - `estimate_study_time`
   - `draft_study_plan`
   - `schedule_reminder` dengan konfirmasi/approval yang sesuai.
   Progress tidak boleh mengasumsikan `activity_completions` tersedia karena tabel itu masih Tier 2 opsional.
3. MCP server `campus-lms` dengan token/auth per-tenant. Demo nyata menggunakan operasi yang ada, misalnya membaca offering aktif, assignment/deadline, status submission, published grade, dan mencari material.
4. Guardrail: material tidak dapat mengubah system instruction; PII redaction; allow-list tool per role; setiap tool call menulis `agent_actions`.
5. Approval flow: action berdampak disimpan `proposed`, baru dapat `executed` setelah approval yang sah. Constraint database mengikuti `domain-ai.md`.
6. Red-team suite minimal 15 skenario: prompt injection, cross-tenant extraction, privilege escalation, tool abuse, jailbreak.
7. Eval agent: success rate, steps, token/cost, tool error.
8. `docs/notes/agent-design.md` + `docs/notes/threat-model-ai.md`.

**Bukan deliverable wajib Minggu 10:** grading assistant berbasis tabel `rubrics`. `rubrics` berada di Tier 3, sehingga tidak boleh dibuat hanya untuk mengejar fitur tersebut. Capability itu tetap terdokumentasi di `domain-ai.md` untuk iterasi sesudah scope pertama.

**Definition of Done**
- [ ] Setiap tool call memiliki row `agent_actions` dengan tenant/user/session/step/tool/status/trace dan arguments yang sudah direduksi PII-nya.
- [ ] RLS mencegah agent membaca action/data tenant lain.
- [ ] Constraint menolak status `executed` untuk action yang `requires_approval = true` tanpa `approved_by`.
- [ ] MCP server terbukti digunakan dari client MCP pada operasi yang benar-benar tersedia.
- [ ] Agent berhenti pada hard limit langkah/token/waktu yang dikonfigurasi dan memberi error yang eksplisit.
- [ ] Red-team suite dijalankan; hasil aktual dicatat. Jangan menulis “100% tertahan” sebelum evidence benar-benar menunjukkan itu.
- [ ] Tidak ada `rubrics`, `final_grades`, `announcements`, atau tabel lain di luar scope Minggu 10 yang dibuat diam-diam untuk mendukung demo agent.

**Alokasi:** `agent_actions` + approval 5 · LangGraph/study-plan 9 · MCP 6 · Guardrails/red-team 9 · Eval 3 · Dokumentasi 3

**Sinyal CV:** *"Merancang LangGraph agent + MCP server dengan RLS, immutable tool audit (`agent_actions`), human approval, hard budget, dan red-team suite yang hasilnya dapat ditelusuri ke evidence."*
---

### MINGGU 11 — Kubernetes, IaC, & Kematangan Platform

**Tujuan:** Kamu bisa berbicara Kubernetes dengan percaya diri di interview, dan infrastrukturmu terdefinisi sebagai kode.

> **Sikap yang jujur:** untuk portofolio skala ini, Docker Compose sudah tepat dan lebih hemat. Kubernetes dipelajari karena **muncul di screening**, dan supaya kamu tahu kapan *tidak* memakainya. Itu justru jawaban yang mengesankan.

**Konsep wajib**
- Arsitektur K8s: control plane (api-server, etcd, scheduler, controller-manager) vs node (kubelet, kube-proxy, container runtime).
- Objek inti: Pod, ReplicaSet, Deployment (rolling update, maxSurge/maxUnavailable), Service (ClusterIP/NodePort/LoadBalancer), Ingress, ConfigMap, Secret, PVC/StorageClass, StatefulSet (untuk Postgres — dan diskusi kenapa DB di K8s itu perdebatan), Job/CronJob, Namespace.
- Operasional: resource requests/limits, QoS class, liveness/readiness/startup probe, HPA, PodDisruptionBudget, node affinity/taint-toleration, RBAC, ServiceAccount, NetworkPolicy.
- Debugging: `kubectl describe`, `logs -p`, `exec`, `port-forward`, `events`, membaca CrashLoopBackOff / ImagePullBackOff / OOMKilled / Pending.
- Packaging: Helm (chart, values, template) atau Kustomize (base + overlay). Konsep GitOps (ArgoCD) — cukup paham alurnya.
- **IaC dengan Terraform**: provider, resource, state (dan bahaya state lokal), variable, output, module, `plan` vs `apply`, drift. Targetkan hal yang gratis: DNS/Cloudflare, dan resource Azure (Resource Group, VNet, NSG, VM) — `azurerm` provider adalah salah satu yang paling banyak dipakai di pasar kerja, jadi ini pilihan yang menguntungkan.

**Deliverable project**
1. Cluster **k3d** lokal (1 server + 1 agent, memori dibatasi) menjalankan seluruh LMS: manifest lengkap di `deploy/k8s/`.
2. Helm chart `campus-lms` (atau Kustomize base + overlay dev/prod), dengan probe, resource limit, HPA pada `api`, PDB, NetworkPolicy yang membatasi akses ke DB.
3. Skenario latihan yang **kamu jalankan dan dokumentasikan**: rolling update tanpa downtime, pod dibunuh saat trafik jalan (buktikan self-healing), OOMKill dan perbaikan limit, scale 1→5 replika di bawah load k6.
4. Terraform untuk: DNS Cloudflare + Resource Group/VNet/NSG/VM Azure, dengan `plan` bersih dan state ter-backup (state di Azure Blob Storage — masuk kuota gratis).
5. Job queue & CronJob: reminder deadline, batch ingest nightly, eval nightly.
6. `docs/adr/0004-compose-vs-kubernetes.md` — **kenapa produksi tetap di Compose** dan pada kondisi apa kamu akan pindah ke K8s. (Ini jawaban interview yang matang: engineer yang tahu kapan tidak memakai teknologi.)
7. `docs/notes/k8s-debugging-lab.md` — 6 kegagalan yang kamu buat sendiri dan cara kamu mendiagnosisnya.

**Definition of Done**
- [ ] `kubectl get pods` menampilkan seluruh stack LMS Running di k3d, aplikasi bisa diakses lewat Ingress.
- [ ] Rolling update berjalan tanpa satu pun request gagal (loop k6 selama update).
- [ ] Kamu bisa mendiagnosis CrashLoopBackOff, ImagePullBackOff, Pending, dan OOMKilled **tanpa googling**.
- [ ] `terraform plan` bersih (no drift) dan infrastruktur bisa dibangun ulang dari nol.
- [ ] ADR compose-vs-k8s ditulis dengan argumen biaya & kompleksitas yang konkret.

**Alokasi:** Konsep K8s 8 · Manifest & Helm 10 · Debugging lab & chaos 8 · Terraform 6 · ADR/dokumentasi 3

**Sinyal CV:** *"Menjalankan LMS di Kubernetes (k3s) dengan Helm, HPA, PDB, dan NetworkPolicy; Terraform untuk DNS & jaringan cloud; mendokumentasikan keputusan arsitektur Compose vs K8s berdasarkan biaya operasional."*

---

### MINGGU 12 — Security Hardening, AI Quota, Dokumentasi, Portofolio, & Kesiapan Interview

**Tujuan:** Menutup rilis pertama dengan security audit, kuota AI tenant-aware, dokumentasi yang jujur, dan bukti yang siap dipakai saat melamar.

> `ai_quotas` adalah Tier 2 yang dipetakan ke Minggu 12 oleh `domain-ai.md`. Implementasinya diawali keputusan promosi tertulis sesuai aturan tier.

**Konsep wajib**
- OWASP API Security Top 10 dan OWASP Top 10 for LLM Applications.
- Auth: access token pendek, refresh rotation/revocation melalui `auth_sessions`, password hashing, MFA bila diimplementasikan, tenant/course authorization.
- Rate limiting berlapis: IP, user, tenant.
- `ai_quotas`: period, token/request budget, hard limit vs degradasi.
- CORS, security headers, input validation, upload safety, SSRF prevention, secret rotation.
- Kesadaran UU PDP sebagai desain/operasional; jangan klaim compliance formal tanpa audit yang sah.
- Dokumentasi: README, ADR, runbook, OpenAPI, evidence.

**Deliverable project**
1. Promosi tertulis + migrasi `ai_quotas`: `tenant_id`, period start/end, token/request budget & usage, `hard_limit`, `updated_at`; RLS; unique `(tenant_id, period_start)`.
2. Enforcement quota + rate limit per tenant **dan** per user. Saat quota habis: tolak secara jelas atau degradasi sesuai policy; jangan menggantung.
3. Audit keamanan mandiri terhadap endpoint/core flow yang benar-benar ada. Semua finding nyata diperbaiki dan diberi regression test. **Jangan mengarang kerentanan hanya untuk memenuhi target angka.** Jika finding nyata kurang dari lima, lengkapi latihan dengan minimal lima negative/security test scenario yang sengaja dibuat, dan bedakan jelas antara “finding” dan “test scenario”.
4. Security CI: cross-tenant IDOR/BOLA tests, role escalation tests, agent approval bypass tests, quota bypass tests, dan scanner yang dipilih.
5. README kelas satu: arsitektur, domain ownership, implementation tiers, trade-off, angka performa/eval/cost yang hanya diambil dari evidence, cara run lokal, demo, dashboard/laporan.
6. Demo video 5 menit: alur LMS → AI artifact review/RAG → observability/eval → CI/CD → satu keputusan arsitektur.
7. Tiga tulisan teknis berbasis pekerjaan yang benar-benar dilakukan:
   - multi-tenant Postgres dengan RLS + composite FK;
   - RAG + eval harness;
   - prompt injection/tool audit pada agent LMS.
8. CV 1 halaman + LinkedIn. Setiap angka menggunakan output nyata dari `docs/progress/evidence/`; gunakan placeholder sampai bukti ada.
9. Persiapan interview: 20 pertanyaan dari desain dan kode aktual, termasuk course vs course_offering, RLS vs app filter, composite FK, immutable submission version, AI human review, fallback provider, RAG grounding, tiering, dan agent approval.
10. Opsional: eksplorasi OpenClaw hanya bila seluruh DoD inti selesai.

**Definition of Done**
- [ ] `ai_quotas` tenant-scoped, RLS aktif, unique period constraint bekerja, dan cross-tenant quota reference ditolak.
- [ ] Quota/rate-limit enforcement diuji pada tenant dan user boundary.
- [ ] Seluruh finding keamanan yang benar-benar ditemukan sudah memiliki fix + regression test; finding dan skenario latihan tidak dicampur.
- [ ] Orang lain dapat menjalankan project dari README dengan langkah yang terdokumentasi.
- [ ] Demo/video/blog/CV hanya memuat fitur dan angka yang sudah ada evidencenya.
- [ ] Tidak ada klaim implementasi Tier 3 yang sebenarnya hanya didokumentasikan.
- [ ] Melamar ke minimal 15 posisi setelah materi portofolio yang diperlukan siap.

**Alokasi:** Security audit/test/fix 10 · `ai_quotas` + rate limit 5 · README/dokumentasi 6 · Video/blog 7 · CV/interview/lamaran 5 · Buffer 2

**Sinyal CV:** tulis **setelah** evidence tersedia, dengan pola *aksi + teknologi + hasil terukur*. Jangan mengisi angka performa/eval/security dari target roadmap.
---

## 9. Ritme Harian & Aturan Main

### Template hari kerja (7 jam × 5 hari + 1 hari ringan)

| Blok | Durasi | Isi |
|---|---|---|
| 08.00–08.30 | 30 m | Review catatan kemarin, tulis 3 target hari ini di `docs/journal/YYYY-MM-DD.md` |
| 08.30–10.30 | 2 j | **Deep work 1** — konsep baru, tanpa AI autocomplete, tanpa notifikasi |
| 10.30–10.45 | 15 m | Istirahat (jauh dari layar) |
| 10.45–12.45 | 2 j | **Deep work 2** — implementasi ke project |
| 12.45–13.45 | 1 j | Makan & istirahat |
| 13.45–15.45 | 2 j | **Deep work 3** — implementasi/debugging (AI agent boleh dipakai penuh) |
| 15.45–16.45 | 1 j | Test, commit, PR, update dokumentasi/ADR |
| 16.45–17.00 | 15 m | **Jurnal**: apa yang berhasil, apa yang macet, satu hal yang dipelajari |

**Hari ke-6 (ringan, ~3–4 jam):** kejar ketertinggalan, tulis blog, review DoD minggu ini.
**Hari ke-7:** libur total. Otak butuh konsolidasi; ini bukan kemalasan, ini bagian dari metode.

### Aturan yang menentukan keberhasilan

1. **Semua masuk repo.** Catatan, eksperimen gagal, ADR — semua di git. Jejak berpikir ini yang membuat interview lancar.
2. **Angka, selalu angka.** Setiap klaim terukur harus menunjuk evidence. Jangan menyalin angka contoh dari roadmap menjadi klaim CV; gunakan hanya hasil run milikmu.
3. **Rusakkan dengan sengaja, tetapi terkendali.** Setiap minggu, picu satu failure yang aman di environment latihan (mis. matikan DB/service atau cabut satu provider LLM), lalu pulihkan dan simpan evidence. Jangan merusak data yang tidak dapat dipulihkan.
4. **Jelaskan keras-keras.** Rekam penjelasan 3 menit tiap akhir minggu. Kalau tersendat, kamu belum paham.
5. **Deploy sesering mungkin.** Fitur yang tidak di-deploy tidak dihitung.
6. **Jangan gold-plating.** UI LMS-mu cukup rapi & fungsional. Nilai jualmu ada di backend, infra, dan AI — bukan di animasi.
7. **Jangan sendirian.** Post progres mingguan di LinkedIn/X. Ini membangun jaringan sekaligus menciptakan akuntabilitas publik.

### Rencana kalau tertinggal

Jika di Minggu 6 kamu ketinggalan > 1 minggu, **potong dengan urutan ini** (paling boleh dipotong lebih dulu):
1. Terraform (Minggu 11) — cukup baca, tidak wajib implementasi penuh.
2. Kubernetes (Minggu 11) — turunkan jadi lab debugging minimum.
3. Eksperimen AI tambahan di luar konfigurasi minimum yang dibutuhkan untuk memilih baseline.
4. E2E Playwright — cukup integration test bila waktu sempit.
5. Tier 2 tambahan yang tidak dipetakan eksplisit oleh `domain-ai.md` — **jangan** menaikkan tier hanya untuk menambah jumlah fitur.

**Yang TIDAK BOLEH dipotong:** Docker, Postgres/RLS + composite FK, `course_offerings`, auth session, migrasi Tier 1, CI/CD, observability dasar, immutable submission versioning, RAG tenant-safe, **eval harness**, human review AI, dan guardrails agent.

---

## 10. Portofolio, CV, dan Strategi Melamar

### 10.1 Satu project besar > lima project kecil

Repo `campus-lms` adalah keseluruhan portofoliomu. Sistem yang dalam dan nyata membuktikan lebih banyak dibanding lima tutorial. Rapikan repo agar bercerita sendiri: README, `docs/adr/`, `docs/runbook/`, dashboard publik, laporan eval, demo live.

### 10.2 Kapan mulai melamar

**Minggu 7.** Bukan Minggu 12. Alasannya: setelah Minggu 6 kamu sudah memenuhi baseline backend, dan proses rekrutmen makan waktu 3–8 minggu. Lamaran di Minggu 7 akan sampai tahap interview saat kamu sedang mengerjakan fase AI — dan kamu bisa bercerita tentang pekerjaan yang sedang berlangsung, yang justru terdengar hidup.

### 10.3 Positioning per pasar

**Indonesia (Bali/Jakarta):** posisikan diri sebagai *"Backend engineer yang bisa mengoperasikan sendiri sistemnya dan sudah mengirim fitur AI ke produksi."* Sebagian besar startup lokal sedang mencari orang yang bisa menambahkan AI ke produk mereka tanpa perlu tim ML. Target: startup edtech, agency produk digital, software house, dan kampus (produkmu literally LMS — pitch ke kampusmu sendiri; pilot user nyata adalah keunggulan besar).

**Remote internasional:** naikkan penekanan pada bahasa Inggris tertulis (README, blog, commit message semua dalam Inggris), timezone overlap, dan bukti kerja asinkron (PR yang deskriptif, ADR, dokumentasi). Sertakan angka biaya & performa — tim remote sangat peduli pada engineer yang paham cost.

### 10.4 Format baris CV yang berhasil

Pola: **[Aksi] + [teknologi spesifik] + [hasil terukur]**

Contoh yang lemah:
> *Membuat aplikasi LMS dengan Go dan AI.*

Contoh yang kuat:
> *Membangun SaaS LMS multi-tenant (Go, Next.js, Postgres+pgvector, FastAPI) dengan RLS + composite FK, `course_offerings`, CI/CD GitHub Actions, dan auto-rollback — lead time serta p95 dicantumkan dari evidence aktual.*
> *Mengirim fitur AI human-reviewed (auto-summary), RAG tanya-materi bersitasi, dan study-plan agent dengan LLM gateway multi-provider, eval harness, dan tracing — faithfulness, latensi, serta cost/request hanya dicantumkan setelah diukur.*

### 10.5 Kanal melamar

- **Lokal:** LinkedIn, Dealls, Glints, Kalibrr, Techinasia Jobs, Discord/Telegram komunitas dev Indonesia, dan **jaringan kampus** (dosen sering punya kanal ke industri; LMS-mu adalah pembuka pintu yang sempurna).
- **Remote:** Wellfound, RemoteOK, WeWorkRemotely, Otta, HN "Who is hiring", proyek open source (kontribusi kecil ke repo tool yang kamu pakai — Langfuse, promptfoo, dan sejenisnya menerima kontribusi pemula, dan ini jalur referral yang nyata).

---

## 11. Checklist Skill Akhir

Centang saat kamu bisa melakukannya **tanpa tutorial**.

### Linux, jaringan, operasional
- [ ] Navigasi & administrasi Linux, permission, systemd unit, journalctl
- [ ] Diagnosis jaringan: `dig`, `ss`, `curl -v`, `tcpdump` dasar
- [ ] Hardening server: SSH key-only, firewall, fail2ban, auto-update
- [ ] Reverse proxy + TLS otomatis (Caddy), paham setara-nya di Nginx
- [ ] Provisioning & operasi VM cloud (Azure): VNet, NSG, public IP, deallocate vs stop
- [ ] **Cost engineering**: budget alert, right-sizing, scale-up sementara, menjalankan produksi dalam anggaran $0
- [ ] Menulis dan mengikuti runbook incident

### Container & orkestrasi
- [ ] Menjelaskan namespaces & cgroups tanpa analogi VM
- ✅ Multi-stage, multi-arch build/push, non-root, image < 25 MB — dibuktikan pada Minggu 2 (runtime arm64 tetap belum diuji)
- [ ] Compose dengan profiles, healthcheck, resource limit
- [ ] Deploy & debug workload di Kubernetes (k3s), Helm dasar
- [ ] Menjelaskan **kapan tidak** memakai Kubernetes

### Backend & data
- [ ] Desain API REST idiomatik + OpenAPI
- [ ] Desain skema Postgres, migrasi versioned, expand-contract
- [ ] Menjelaskan `courses` vs `course_offerings` dan menjaga histori antarsemester
- [ ] Multi-role tenant melalui `membership_roles` dan course-scoped authorization melalui `course_staff`
- [ ] Multi-tenancy dengan RLS **dan composite foreign key**, beserta integration test cross-tenant
- [ ] Immutable `submission_versions` dan tenant-aware file/submission relation
- [ ] Membaca `EXPLAIN ANALYZE` dan memperbaiki query lambat berdasarkan evidence
- [ ] Transaksi, isolation level, mencegah race condition
- [ ] Caching Redis, background job queue, idempotency
- [ ] Backup & **restore yang benar-benar diuji**

### Delivery & kualitas
- [ ] Unit + integration test (testcontainers) + sedikit e2e
- [ ] Pipeline CI/CD lengkap dengan security scanning
- [ ] Rilis otomatis dengan health check & rollback
- [ ] Trunk-based development, PR kecil, conventional commits

### Observability
- [ ] Instrumentasi OpenTelemetry (traces, metrics, logs berkorelasi)
- [ ] PromQL dasar, dashboard Grafana yang berguna, alert bermakna
- [ ] Mendefinisikan SLI/SLO & error budget
- [ ] Load test k6 dan profiling pprof

### AI Engineering
- [ ] Prompt engineering produksi + structured output tervalidasi
- [ ] `ai_artifacts` human review: draft/approved/rejected/published dan stale detection via `input_hash`
- [ ] Feature flag AI per tenant melalui `tenant_settings`
- [ ] LLM gateway: routing, fallback, retry, caching, token/cost tracking
- [ ] RAG: `material_chunks`, parsing, chunking, embedding, pgvector, hybrid search, reranking, sitasi
- [ ] **Eval harness** + `ai_feedback`: golden dataset, ragas, LLM-as-judge, gate di CI
- [ ] LLM observability: `ai_interactions`, tracing, cost per tenant, latency p95
- [ ] Agent: LangGraph, tool design, agentic RAG, checkpointing, human-in-the-loop
- [ ] `agent_actions`: immutable tool audit + approval constraint
- [ ] **MCP**: membangun server, mengekspos tools dengan aman
- [ ] `ai_quotas`: token/request budget per tenant
- [ ] Guardrails: prompt injection, PII, least privilege, RLS, red-teaming

### Keamanan
- [ ] OWASP API Top 10 & OWASP LLM Top 10 — bisa mengaudit sistem sendiri
- [ ] AuthN/AuthZ: JWT rotation, Argon2id, RBAC per-tenant
- [ ] Rate limiting berlapis, upload aman, manajemen secret
- [ ] Kesadaran UU PDP untuk data pendidikan

### Profesional
- [ ] Menulis ADR & dokumentasi teknis yang dibaca orang
- [ ] Menjelaskan trade-off, bukan hanya menyebut teknologi
- [ ] Memverifikasi & men-debug kode hasil AI, dan mengartikulasikan prosesnya
- [ ] Punya angka untuk setiap klaim

---

## 12. Sumber & Catatan Verifikasi

### Bagaimana dokumen ini disusun

Klaim tentang pasar kerja, free tier, spesifikasi hardware, dan metode pembayaran **bersumber dari pencarian web bertanggal 2026**, bukan dari asumsi. Yang tidak bisa saya verifikasi, saya sebutkan sebagai tidak terverifikasi (mis. "hermes agent" di §2.4, dan status kartu Superbank di §6.5). Angka free tier saya beri peringatan eksplisit karena berubah cepat — Oracle memangkas alokasinya di pertengahan 2026 tanpa pengumuman publik, dan DigitalOcean menutup program Student Pack-nya pada Juli 2026. Keduanya membuktikan kenapa kolom "Kartu?" dan tanggal sumber itu penting.

**Satu hal yang sengaja TIDAK saya sarankan:** membuat email `.edu` palsu lewat generator untuk mengakali verifikasi Azure/Student Pack. Situs semacam itu muncul di hasil pencarian, tapi itu pemalsuan identitas akademik — melanggar ToS Microsoft/GitHub, akun bisa dihapus beserta seluruh datamu, dan tidak ada gunanya karena **kamu memang mahasiswa sungguhan**. Pakai email `.ac.id` kampusmu atau upload KTM.

### Rujukan utama

**Pasar kerja & skill**
- [Junior Developer Hiring Crisis 2026 — nucamp](https://www.nucamp.co/blog/the-junior-developer-hiring-crisis-in-2026-how-to-get-your-first-backend-job) — ekspektasi "junior" 2026: satu bahasa + SQL + Docker + CI/CD + cloud + verifikasi kode AI
- [How to Hire a Backend Developer in 2026 — KORE1](https://www.kore1.com/hire-backend-developer/) — matriks skill entry/mid/senior; Docker wajib di mid, K8s "helpful" di junior
- [AI Engineer Skills 2026 — technovids](https://technovids.com/ai-engineer-skills) — peta skill AI engineer; "projects must be deployed and publicly accessible"
- [The AI Skills That Mattered in 2025 Are Already Obsolete — gnxt](https://gnxtsystems.com/the-ai-skills-that-mattered-in-2025-are-already-obsolete-heres-what-2026-demands/) — agent orchestration, MCP, evaluation design, agentic RAG, observability sebagai pembeda
- [13 Essential AI Developer Skills 2026 — doit.software](https://doit.software/blog/ai-developer-skills) — RAG hybrid + reranking, tool calling, MCP
- [AI Engineer Resume Examples 2026](https://resumeoptimizerpro.com/blog/ai-engineer-resume-examples) — apa yang dicari parser CV: framework orkestrasi + tool eval + angka produksi

**Infrastruktur gratis tanpa kartu**
- [Azure for Students — GitHub Education Pack](https://education.github.com/pack) — halaman resmi, menyatakan eksplisit *"no credit card required"* + $100 kredit Azure
- [Azure for Students: eligibility & verifikasi — studentdealslab, Jul 2026](https://studentdealslab.com/tools/azure-for-students-discount/) — $100 / 12 bulan, tanpa kartu, verifikasi via email kampus / upload KTM / GitHub Student Pack; satu kredit per orang, bisa diperbarui tiap tahun
- [Azure for Students: layanan always-free — lifeatlas, Mar 2026](https://www.lifeatlas.site/news/azure-student-credits) — termasuk 750 jam VM B1s, 5 GB blob storage, Functions 1 juta request/bulan
- [GitHub Student Developer Pack — daftar penawaran](https://education.github.com/pack) — JetBrains, Codespaces, Namecheap `.me` gratis, New Relic, Clerk Pro
- [DigitalOcean menghentikan program Student Pack — DO Community, Jun 2026](https://www.digitalocean.com/community/questions/does-the-github-students-pack-removed-from-digitalocean) — redemption tutup 31 Juli 2026, jangan dikejar lagi
- [Cloud Free Tiers 2026 — leanopstech](https://leanopstech.com/blog/cloud-free-tier-comparison-2026/) — Cloudflare Workers/Pages sebagai free tier paling predictable
- [Neon vs Supabase — vela.run](https://vela.run/neon-vs-supabase/) dan [Supabase Pricing 2026](https://designrevision.com/blog/supabase-pricing)
- *Dicoret karena mewajibkan kartu:* [Oracle memangkas Always Free 4/24 → 2/12 — InfoQ, Jul 2026](https://www.infoq.com/news/2026/07/oracle-cloud-free-tier-limits/) (tetap saya cantumkan sebagai arsip keputusan)

**VPS lokal bayar QRIS/GoPay/OVO (Rencana C)**
- [IDCloudHost Cloud VPS](https://idcloudhost.com/en/cloud-vps/) — FAQ resmi mencantumkan GOPAY, OVO, ShopeePay, VA, Alfamart/Indomaret; sistem top-up saldo, tagihan per jam
- [Review Cloud VPS IDCloudHost — dwiay.com](https://dwiay.com/2025/10/02/review-cloud-vps-dari-idcloudhost/) — pengalaman nyata top-up tanpa kartu kredit via QRIS/VA
- [Biznet Gio vs IDCloudHost, Jul 2026 — matthewswong.com](https://www.matthewswong.com/id/blog/biznet-gio-vs-idcloudhost-indonesia-cloud/) — NEO Lite mulai ~Rp59rb, IDCloudHost ~Rp87rb, ditagih rupiah sehingga tidak kena flag transaksi luar negeri
- [Diskusi VPS murah pembayaran lokal — r/indotech](https://www.reddit.com/r/indotech/comments/1l7pn27/saran_vps_buat_pemula_murah_good_tech_support/) — pengalaman komunitas (CloudKilat, Dihostingin, QRIS)

**Kartu ditolak & alternatif bank digital**
- [FAQ resmi bluDebit Card — blu by BCA Digital](https://blubybcadigital.id/info/faq/blu-debit-card) — toggle terpisah *Transaksi Internasional Online*; kartu Mastercard dapat dipakai domestik & internasional; OTP hanya berlaku di merchant 3D Secure
- [Solusi kartu debit BCA Mastercard tidak bisa transaksi online](https://bca.emingko.com/2023/03/solusi-kartu-debit-bca-mastercard-tidak.html) — syarat: fitur debit online aktif, merchant memakai 3D Secure, nomor HP aktif untuk OTP
- [Laporan pengguna: BCA Mastercard declined di merchant internasional — r/indotech, Jul 2026](https://www.reddit.com/r/indotech/comments/1utltck/masalah_gagal_transaksi_pakai_kartu_bca/) — kasus identik dengan kasusmu, termasuk saran toggle *BCA Mobile → Akun Saya → Kontrol → Transaksi Internasional* dan pengalaman komunitas dengan Jago/Jenius
- [Cara mengatasi kartu BCA gagal tertaut ke PayPal](https://viapaypal.id/blog/cara-mengatasi-kartu-bca-gagal-tertaut-ke-paypal/) — penjelasan verifikasi pra-otorisasi ~$1,95 yang sering menggagalkan kartu debit

**LLM API gratis**
- [Best Free LLM 2026 — klymentiev](https://klymentiev.com/blog/best-free-llm-2026)
- [Free LLM APIs Compared — OpenRouter, Jun 2026](https://openrouter.ai/blog/tutorials/free-llm-apis-compared/)
- [Every Free AI API in 2026 — awesomeagents](https://awesomeagents.ai/tools/free-ai-inference-providers-2026/)
> Angka antar sumber tidak selalu cocok. Selalu cek dashboard resmi sebelum bergantung.

**Hardware**
- [Spesifikasi Axioo Hype 5 AMD X5-2 — liputan6](https://www.liputan6.com/tekno/read/6280876/harga-axioo-hype-5-amd-x5-2-laptop-lokal-stylish-dengan-performa-tangguh-dan-baterai-tahan-lama)
- [Axioo Hype 5 AMD X5-2 — ponselesa](https://www.ponselesa.com/2024/12/spesifikasi-dan-harga-axioo-hype-5-amd-x5-2.html)
- [Varian 16 GB — agres.id](https://www.agres.id/products/axioo-hype-5-amd-x5-2-ryzen-5-7430-16gb-256gb-w11-140fhd-ips-blit-hdmi-gry)

**OpenClaw (untuk konteks §2.4)**
- [OpenClaw Tutorial — petronellatech](https://petronellatech.com/blog/openclaw-ai-agent-guide/)
- [OpenClaw: How a Self-Hosted AI Agent Changed Automation in 2026 — Medium/Kanerika](https://medium.com/@kanerika/openclaw-how-a-self-hosted-ai-agent-changed-automation-in-2026-6ba728345d53)

### Dokumentasi primer yang harus jadi rujukan harianmu

Utamakan dokumentasi resmi di atas kursus berbayar. Semuanya gratis:
Docker docs · PostgreSQL manual (bab Indexes, MVCC, Performance Tips) · pgvector README · Go by Example + Effective Go · Kubernetes docs (Concepts) · OpenTelemetry docs · Prometheus docs · Caddy docs · GitHub Actions docs · Terraform Registry · Model Context Protocol spec (`modelcontextprotocol.io`) · LangGraph docs · ragas docs · promptfoo docs · Langfuse docs · OWASP API Security Top 10 & OWASP Top 10 for LLM Applications · The Twelve-Factor App · Google SRE Book (bab SLO — gratis online).

---

## Penutup

Tiga hal yang membuat roadmap ini berbeda dari daftar "belajar Docker, belajar K8s" yang biasa:

1. **Satu project mengikat semuanya.** Setiap skill baru menyelesaikan masalah nyata di LMS yang sama, jadi kamu tidak pernah belajar sesuatu "untuk nanti".
2. **Setiap minggu menghasilkan angka.** Angka inilah yang berubah menjadi baris CV dan jawaban interview — dan ini yang paling langka pada fresh graduate.
3. **Fase AI dibangun di atas fondasi operasional, bukan sebaliknya.** Banyak orang belajar RAG lalu bingung cara men-deploy-nya. Kamu akan melakukan yang kebalikannya, dan itu urutan yang benar.

Keunggulan yang sudah kamu miliki — orkestrasi model per-tugas lewat combo di opencode — bukan hal sepele. Sebagian besar pelamar belum berpikir sejauh itu tentang routing dan efisiensi token. Tugasmu 12 minggu ke depan adalah mengubah intuisi itu menjadi **kode produksi yang bisa kamu tunjukkan**.

Mulai Senin. Minggu 1, jam 08.00. Buat repo-nya hari ini.
