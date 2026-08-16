# Laporan Minggu 00 — Persiapan Akun & Agent

> **Status: BELUM DIISI.** Ini kerangka yang sudah disesuaikan untuk Minggu 0.
> Isi sambil/sesudah mengerjakan `docs/setup/azure-day-0.md`.
> Minggu 0 hampir seluruhnya tugas manusia, jadi section "Dikerjakan Agent"
> memang akan tipis — itu normal, bukan kekurangan.

- **Periode:** ______ s/d ______
- **Fokus roadmap:** Persiapan akun cloud, klaim benefit, setup agent
- **Total jam:** ______ (target ±4)
- **Commit range:** ______

---

## 1. Ringkasan

<3–5 kalimat. Apa yang sudah aman sekarang yang sebelumnya belum.>

---

## 2. Dikerjakan Agent

| # | Pekerjaan | File | Commit | Bukti |
|---|---|---|---|---|
| 1 | <mis. verifikasi `make todo` dan ringkasan task Minggu 1> | — | | |

> Minggu ini agent hanya membantu verifikasi setup. Kalau agent mengarang task
> yang tidak ada di repo saat diuji, catat di section 8 — itu temuan penting.

---

## 3. Dikerjakan Manusia

| # | Pekerjaan | Kenapa harus manual | Hasil |
|---|---|---|---|
| 1 | Verifikasi Azure for Students | Verifikasi identitas akademik | ✅ / ❌ |
| 2 | Cek spending limit ON | Portal cloud, tidak bisa diotomasi | |
| 3 | Buat budget alert $10/bulan | Portal cloud | |
| 4 | Aktifkan MFA | Keamanan identitas | |
| 5 | Cek kuota vCPU Southeast Asia | Portal cloud | tersedia: ___ vCPU |
| 6 | Klaim Namecheap / Codespaces / JetBrains / New Relic | Akun pribadi | |
| 7 | Buat SSH key ed25519 | Kredensial — agent dilarang menyentuh | |
| 8 | Install & login Azure CLI | Kredensial | |
| 9 | Daftar API key LLM (Gemini/Groq/Cerebras/OpenRouter) | Kredensial | |
| 10 | `git init` + push repo publik | Keputusan & akun | URL: |
| 11 | Baca AGENTS.md + agent/policy.md | Pemahaman, tidak bisa didelegasikan | |

---

## 4. Keputusan yang Diambil

| Keputusan | Alternatif yang ditolak | Alasan | ADR |
|---|---|---|---|
| Region Southeast Asia | Australia East, East Asia | <isi> | ADR-0002c |
| Nama resource group | | | ADR-0002c |

---

## 5. Angka & Bukti

| Metrik | Nilai | Cara diukur | File bukti |
|---|---|---|---|
| Kuota BS Family vCPU (Southeast Asia) | | Portal → Usage + quotas | screenshot/catatan |
| Kredit Azure tersisa | $ | Portal → Cost Management | |
| Proyeksi biaya bulanan | $ | Portal → Cost analysis | |

> Minggu 0 sebagian buktinya berupa screenshot portal, bukan output terminal.
> Itu sah — tapi tetap simpan di `docs/progress/evidence/week-00/` dan sebutkan
> tanggal pengambilannya.

---

## 6. Konsep yang Dipelajari

### Spending limit vs budget alert

- **Apa:** <isi sendiri setelah membacanya>
- **Kenapa dipakai di sini:**
- **Alternatif yang tidak dipilih:**
- **Cara membuktikan sendiri:**
- **Pertanyaan interview terkait:** *"Bagaimana kamu mencegah cost overrun di cloud?"*

### `az vm stop` vs `az vm deallocate`

- **Apa:**
- **Kenapa penting di sini:** kredit $100 harus bertahan 12 bulan
- **Cara membuktikan sendiri:** bandingkan status VM dan proyeksi biaya sebelum/sesudah
- **Pertanyaan interview terkait:** *"Bagaimana kamu mengelola biaya infrastruktur?"*

---

## 7. Belum Terverifikasi

| Hal | Kenapa belum terverifikasi | Rencana verifikasi |
|---|---|---|
| Apakah VM B1s benar-benar masuk kuota gratis | Belum membuat VM | Minggu 4 |
| Apakah budget alert benar-benar terkirim | Belum melewati ambang | Turunkan ambang sementara untuk tes |
| Apakah agent patuh pada AGENTS.md di task nyata | Baru diuji sekali | Minggu 1 |

---

## 8. Masalah & Cara Diselesaikan

### Masalah: <judul, kalau ada>
- **Gejala:**
- **Akar masalah:**
- **Solusi:**
- **Waktu terbuang:**

---

## 9. Status Definition of Done

| DoD Minggu 0 | Usulan | Bukti | Dicentang manusia |
|---|---|---|---|
| Spending limit ON + budget alert terkirim | | | ☐ |
| Kuota vCPU dikonfirmasi | | | ☐ |
| Repo publik hidup, `make todo` jalan | | | ☐ |
| Bisa jelaskan `stop` vs `deallocate` | | | ☐ |
| Laporan ini terisi & ditandatangani | | | ☐ |

---

## 10. Untuk Minggu Depan

- **Persiapan Minggu 1:** WSL2 siap, dotfiles, `.wslconfig` (kalau RAM 8 GB)
- **Carry-over:**

---

## 11. Verifikasi Manusia

- [ ] Spot-check bukti — *(Minggu 0 sebagian besar screenshot, cukup pastikan tanggalnya benar)*
- [ ] Skor quiz: ____ / ____
- [ ] Explain-back 3 menit direkam: `docs/progress/explain/week-00.___`
- [ ] Saya sudah benar-benar membaca `AGENTS.md` dan `agent/policy.md`, bukan sekadar membuka

**Ditandatangani:** ______________  **Tanggal:** __________
