# Folder `agent/` — Panduan untuk KAMU (manusia)

> Bahasa Indonesia karena ini materi belajarmu. File di dalam `rules/` dan
> `prompts/` berbahasa Inggris karena itu dikonsumsi oleh model.

## Kenapa folder ini ada

Kamu memilih mode **agent-first**: agent mengerjakan semua yang bisa
diotomasi lokal, kamu menangani yang mustahil diotomasi. Itu keputusan yang
sah dan realistis untuk 2026 — tapi ada satu bahaya yang harus dijinakkan:

**Kalau agent menulis kode dan kamu hanya menerimanya, kamu akan punya repo
yang mengesankan dan pemahaman yang rapuh.**

Folder ini adalah sistem penjinaknya. Isinya tiga mekanisme:

1. **Kontrak kerja** (`AGENTS.md`, `policy.md`) — batas tegas apa yang boleh
   dan tidak boleh dikerjakan agent.
2. **Protokol bukti** (`evidence-protocol.md`) — agent dilarang mengklaim
   tanpa melampirkan output perintah asli.
3. **Mekanisme belajar** (`templates/`, `checklists/`) — laporan mingguan,
   quiz verifikasi, dan explain-back yang memaksa pemahaman terbentuk.

## Isi folder

| File | Untuk siapa | Isi |
|---|---|---|
| `policy.md` | Agent | Izin detail, hard stop, prosedur eskalasi |
| `evidence-protocol.md` | Agent | Format bukti & laporan, aturan anti-halusinasi |
| `rules/00-global.md` | Agent | Aturan lintas area |
| `rules/10-go-api.md` | Agent | Konvensi Go, error handling, struktur |
| `rules/20-database.md` | Agent | Migrasi, RLS, index, query |
| `rules/30-docker-deploy.md` | Agent | Dockerfile, Compose, deploy, batasan RAM |
| `rules/40-ci-cd.md` | Agent | GitHub Actions, testing, supply chain |
| `rules/50-python-ai.md` | Agent | FastAPI, LLM, RAG, eval, agent |
| `rules/60-security.md` | Agent | OWASP API & LLM, secrets, multi-tenancy |
| `rules/70-docs.md` | Agent | Cara menulis laporan, ADR, catatan |
| `prompts/*.md` | Kamu → Agent | Template prompt untuk tiap combo opencode-mu |
| `templates/*.md` | Agent | Kerangka laporan mingguan, quiz, session log |
| `checklists/*.md` | Kamu | Checklist verifikasi sebelum sign-off |

## Cara memakainya sehari-hari

### Saat memulai sesi kerja

Arahkan agent ke konteksnya. Contoh dengan opencode:

```
/plan Baca AGENTS.md, agent/rules/10-go-api.md, dan bagian MINGGU 1 di
roadmap. Buat rencana untuk task: implementasi main.go sesuai spesifikasi.
Sebutkan file yang akan disentuh dan DoD yang ditargetkan.
```

Setelah rencananya masuk akal:

```
/dev Jalankan rencana tadi. Setelah selesai, jalankan verifikasi dan simpan
buktinya sesuai agent/evidence-protocol.md.
```

Template lengkap untuk tiap combo ada di `prompts/`.

### Saat menutup minggu

```
/heavy Susun docs/progress/week-01.md dari pekerjaan minggu ini.
Ikuti agent/templates/weekly-report.md persis. Wajib ada bagian
"Belum Terverifikasi". Setiap angka wajib menunjuk file bukti.
Lalu buat quiz di docs/progress/quiz/week-01.md sesuai
agent/templates/quiz.md — soal HANYA dari kode yang benar-benar ada di repo.
```

Lalu kamu: jawab quiz **tanpa membuka kode**, spot-check 3 bukti secara acak,
rekam explain-back 3 menit, dan tanda tangani laporan.

## Aturan gerbang (jangan dilanggar, ini inti sistemnya)

**Minggu berikutnya tidak dimulai sebelum:**

- [ ] Laporan mingguan ditandatangani
- [ ] Skor quiz ≥ 70%
- [ ] 3 bukti acak sudah kamu spot-check dan cocok
- [ ] Explain-back 3 menit terekam

Skor di bawah 70% **bukan kegagalan** — itu sinyal bagian mana yang belum
nyangkut. Ulangi bagian itu, jangan lanjut. Roadmap 14 minggu dengan
pemahaman utuh jauh lebih berharga daripada 12 minggu dengan repo yang tidak
bisa kamu pertahankan di interview.

## Cara mendeteksi agent yang berhalusinasi

Tanda bahaya yang harus langsung kamu curigai:

| Tanda | Yang harus kamu lakukan |
|---|---|
| Angka bulat mencurigakan ("p95 turun 50%") tanpa file bukti | Minta perintah + output mentah. Kalau tidak ada, angka itu dihapus |
| Bagian "Belum Terverifikasi" kosong | Hampir pasti tidak jujur. Minta agent meninjau ulang |
| Klaim "test lulus" tapi tidak ada output test | Minta jalankan ulang di depanmu |
| Menyebut file/fungsi yang ternyata tidak ada | Hentikan sesi, mulai ulang dengan konteks bersih |
| Jawaban terlalu mulus untuk masalah yang sulit | Tanya "apa yang bisa membuat ini gagal?" — jawaban jujur pasti ada |

Spot-check acak 3 bukti per minggu itu **wajib**, bukan opsional. Agent yang
tahu buktinya diperiksa akan jauh lebih hati-hati.
