# Laporan Minggu <NN> — <Judul Minggu>

> Template. Salin ke `docs/progress/week-<NN>.md`. Disusun agent, ditandatangani manusia.
> Jangan hapus satu section pun. Section yang tidak relevan diisi "Tidak ada minggu ini".

- **Periode:** <YYYY-MM-DD> s/d <YYYY-MM-DD>
- **Fokus roadmap:** <fokus minggu ini>
- **Total jam:** <jam> (target 35)
- **Commit range:** `<sha-awal>..<sha-akhir>` (<jumlah> commit)

---

## 1. Ringkasan

<3–5 kalimat bahasa manusia. Apa yang berubah di sistem minggu ini, dan kenapa itu penting.
Tanpa bahasa marketing. Kalau minggu ini berantakan, tulis berantakan.>

---

## 2. Dikerjakan Agent

| # | Pekerjaan | File | Commit | Bukti |
|---|---|---|---|---|
| 1 | <apa> | `path/file.go` | `abc1234` | [link](evidence/week-NN/xxx.txt) |

**Catatan implementasi:**
- <keputusan teknis kecil yang diambil agent saat menulis kode, dan alasannya>

---

## 3. Dikerjakan Manusia

> Agent WAJIB menanyakan bagian ini, tidak boleh menebak.

| # | Pekerjaan | Kenapa harus manual | Hasil |
|---|---|---|---|
| 1 | <mis. verifikasi Azure di portal> | <tidak bisa diotomasi / kredensial / keputusan> | <hasil> |

---

## 4. Keputusan yang Diambil

| Keputusan | Alternatif yang ditolak | Alasan | ADR |
|---|---|---|---|
| <apa> | <alternatif> | <trade-off> | [ADR-000X](../adr/000X-....md) |

---

## 5. Angka & Bukti

> Setiap baris WAJIB punya file bukti. Tidak ada bukti → tulis `NOT MEASURED`.

| Metrik | Nilai | Cara diukur | File bukti |
|---|---|---|---|
| <mis. ukuran image API> | <21.4 MB> | `docker images ...` | [image-size.txt](evidence/week-NN/image-size.txt) |

**Perbandingan sebelum/sesudah (kalau ada optimasi):**

| Metrik | Sebelum | Sesudah | Perubahan | Bukti |
|---|---|---|---|---|
| | | | | |

---

## 6. Konsep yang Dipelajari

> Bagian terpenting untuk pembelajaranmu. Satu blok per konsep signifikan.

### <Nama konsep>

- **Apa:** <definisi satu kalimat, tanpa menumpuk jargon>
- **Kenapa dipakai di sini:** <alasan konkret di codebase INI, sebut file>
- **Alternatif yang tidak dipilih:** <dan biayanya>
- **Cara membuktikan sendiri:** `<perintah persis yang bisa kamu jalankan>`
- **Pertanyaan interview terkait:** <pertanyaan yang benar-benar ditanyakan hiring manager>

---

## 7. Belum Terverifikasi

> WAJIB DIISI. Section kosong = laporan ditolak.
> Ini bukan kelemahan — ini yang membedakan laporan jujur dari laporan karangan.

| Hal | Kenapa belum terverifikasi | Rencana verifikasi |
|---|---|---|
| <mis. perilaku di bawah 500 concurrent user> | <belum dijalankan load test> | <Minggu 6> |

**Asumsi yang dipakai tapi belum dibuktikan:**
- <asumsi>

---

## 8. Masalah & Cara Diselesaikan

> Termasuk jalan buntu. Bagian ini yang paling berguna saat interview.

### Masalah: <judul>
- **Gejala:** <apa yang terlihat, error persisnya>
- **Hipotesis yang salah:** <yang dicoba dan ternyata bukan penyebabnya>
- **Akar masalah:** <penyebab sebenarnya, dan bagaimana dibuktikan>
- **Solusi:** <apa yang diperbaiki>
- **Pencegahan:** <test/alert/tipe apa yang sekarang menangkap ini>
- **Waktu terbuang:** <jam>

---

## 9. Status Definition of Done

> Agent hanya MENGUSULKAN. Kolom "Dicentang manusia" diisi olehmu setelah melihat bukti.

| DoD dari roadmap | Usulan agent | Bukti | Dicentang manusia |
|---|---|---|---|
| <item DoD> | ✅ terpenuhi / ⚠️ sebagian / ❌ belum | [link](...) | ☐ |

---

## 10. Untuk Minggu Depan

- **Carry-over:** <yang belum selesai>
- **Utang teknis yang sengaja diambil:** <dan kapan harus dibayar>
- **Persiapan yang perlu dilakukan manusia lebih dulu:** <mis. daftar akun>

---

## 11. Verifikasi Manusia

- [ ] Saya sudah spot-check 3 file bukti secara acak dan isinya cocok dengan klaim
- [ ] Skor quiz: ____ / ____ (minimal 70% untuk lanjut)
- [ ] Explain-back 3 menit sudah direkam: `docs/progress/explain/week-NN.<mp3|txt>`
- [ ] Saya bisa menjelaskan setiap keputusan di section 4 tanpa membuka catatan

**Self-audit agent (dari `agent/evidence-protocol.md` §8):** <agent menyatakan sudah/belum dijalankan>

**Ditandatangani:** ______________  **Tanggal:** __________
