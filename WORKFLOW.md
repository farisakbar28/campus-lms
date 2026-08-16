# WORKFLOW — Cara Kerja Repo Ini

> Kartu referensi. Cetak atau buka di layar kedua selama 12 minggu.
> Penjelasan mendalam: `agent/README.md`. Kontrak agent: `AGENTS.md`.

---

## Model mental: tiga lingkaran

```
┌─ LINGKARAN MINGGU (Senin → Minggu) ─────────────────────────────┐
│                                                                  │
│  ┌─ LINGKARAN HARI (±7 jam) ────────────────────────────────┐   │
│  │                                                            │   │
│  │  ┌─ LINGKARAN TASK (30-90 menit) ──────────────────┐      │   │
│  │  │  PLAN → DEV → VERIFY → EVIDENCE → COMMIT        │      │   │
│  │  └──────────────────────────────────────────────────┘      │   │
│  │            ↑ diulang 3-5x per hari                         │   │
│  │  ditutup dengan: make journal                              │   │
│  └────────────────────────────────────────────────────────────┘   │
│            ↑ diulang 5-6x per minggu                             │
│  ditutup dengan: LAPORAN → QUIZ → EXPLAIN-BACK → GATE            │
└──────────────────────────────────────────────────────────────────┘
```

**Aturan tunggal yang menopang semuanya:** tidak ada yang dianggap selesai
tanpa bukti eksekusi. Tidak ada minggu baru tanpa lolos gate.

---

## Siapa mengerjakan apa

| | Agent | Kamu |
|---|---|---|
| Menulis kode, config, test, skrip | ✅ | — |
| Menjalankan & mengukur | ✅ | — |
| Menyusun draf laporan & quiz | ✅ | — |
| Portal Azure, DNS, kredensial | ❌ | ✅ |
| Keputusan & ADR (Decision) | ❌ | ✅ |
| Centang DoD | ❌ | ✅ |
| Jawab quiz, explain-back | ❌ | ✅ |
| Review sebelum merge | draf | ✅ final |

Detail lengkap: `agent/policy.md`

---

## Lingkaran TASK (inti harian)

```bash
make agent-context      # ingatkan agent apa yang wajib dibaca
```

**1. PLAN** — pakai `agent/prompts/plan.md`
Agent membaca aturan, menyusun rencana, menyebut file yang akan disentuh dan
DoD yang ditarget. **Kamu baca rencananya.** Kalau tidak masuk akal, perbaiki
sekarang — bukan setelah 200 baris kode ditulis.

**2. DEV** — pakai `agent/prompts/dev.md`
Agent implementasi.

**3. VERIFY + EVIDENCE**
```bash
CLAIM="deskripsi klaim" make evidence W=01 SLUG=nama-bukti CMD="perintah"
```
Output mentah tersimpan dengan command, timestamp, commit SHA, exit code.
Perintah gagal **tetap disimpan** — kegagalan adalah bukti.

**4. REVIEW + COMMIT** — pakai `agent/prompts/review.md`
Cek `agent/checklists/pre-commit.md`. Kamu baca `git diff`, bukan hanya judul.

> **Macet paham?** `agent/prompts/teach.md`
> **Ada yang rusak?** `agent/prompts/debug.md`
> **Bingung pilih pendekatan?** `agent/prompts/spark.md`

---

## Lingkaran HARI

| Waktu | Aktivitas |
|---|---|
| Awal | `make journal` → tulis 3 target hari ini |
| Isi | 3–5 lingkaran task |
| Akhir | Isi jurnal: yang berhasil, yang macet, satu hal yang dipelajari |

---

## Lingkaran MINGGU

**Senin pagi:**
```bash
make week-init W=01     # buat laporan + quiz + folder bukti
```
Baca bagian "MINGGU 1" di `docs/roadmap.md`. Catat DoD-nya.

**Jumat/Sabtu — tutup minggu:**

1. Agent menyusun laporan → `agent/prompts/heavy.md`
2. Agent membuat quiz dari kode yang **benar-benar ada**
3. Kamu jalankan gate:

```bash
make gate
```

| # | Gate | Waktu |
|---|---|---|
| 1 | Spot-check 3 file bukti acak — jalankan ulang perintahnya | 10 mnt |
| 2 | Baca laporan dengan curiga (cek "Belum Terverifikasi") | 10 mnt |
| 3 | Jawab quiz tanpa membuka kode — **min. 70%** | 15 mnt |
| 4 | Rekam explain-back 3 menit | 5 mnt |
| 5 | Baca `git log` + `git diff` 2 file terpenting | 5 mnt |
| 6 | Tanda tangan laporan | — |

Checklist lengkap: `agent/checklists/human-verification.md`

**Skor < 70% = ulangi bagian itu, jangan lanjut.** Bukan hukuman — alat ukur.

---

## Perintah yang dipakai sehari-hari

```bash
make help                # semua perintah
make todo                # pekerjaan tersisa, berlabel minggu
make agent-context       # bacaan wajib agent
make journal             # jurnal hari ini
make week-init W=NN      # buka minggu baru
make evidence W=NN SLUG=x CMD="..."   # rekam bukti
make gate                # checklist tutup minggu

make up / down / logs / ps            # stack (mulai Minggu 2)
make test / lint / test-cover         # kualitas kode
make docker-size                      # cek target < 25MB (Minggu 2)
make health                           # cek /healthz & /readyz
make prune                            # bersihkan Docker (tiap Jumat)
```

---

## Tanda bahaya (hentikan, jangan lanjut)

| Tanda | Artinya | Tindakan |
|---|---|---|
| Angka tanpa file bukti | Kemungkinan fabrikasi | Hapus angkanya, minta ukur ulang |
| "Belum Terverifikasi" kosong | Laporan tidak jujur | Tolak, minta tinjau ulang |
| Agent sebut file yang tidak ada | Konteks kotor | Mulai sesi baru, muat ulang konteks |
| Saya tidak paham >30% kode minggu ini | Terlalu cepat | Pelankan, pakai `/teach` |
| Quiz terasa terlalu mudah | Soal dangkal | Minta soal failure-mode |
| Bukti tidak bisa saya reproduksi | Bukti dikarang | Investigasi serius |

---

## Kalau tertinggal

Boleh dipotong (urutan dari yang paling boleh):
Terraform → Kubernetes → Grading assistant → E2E Playwright

**Tidak boleh dipotong:** Docker · Postgres/RLS · deploy · CI/CD ·
observability dasar · RAG · **eval harness** · guardrails · **gate mingguan**
