# Architecture Decision Records

ADR adalah catatan singkat: **keputusan apa yang diambil, kenapa, dan apa
konsekuensinya.** Satu file per keputusan, tidak pernah dihapus — kalau
berubah pikiran, tulis ADR baru yang menggantikan (`Supersedes ADR-000X`).

## Kenapa ini ada di roadmap

Dua alasan yang sangat praktis:

1. **Interview.** Pertanyaan "kenapa kamu pilih X?" adalah pertanyaan
   paling sering di interview mid-level. Kandidat yang menyebut nama
   teknologi kalah dari kandidat yang menjelaskan trade-off. ADR memaksamu
   berlatih itu tiap minggu.
2. **Recruiter membaca repo.** Folder ADR yang terisi rapi adalah sinyal
   kematangan yang jarang ditemukan di portofolio fresh graduate.

## Daftar

| No | Judul | Minggu | Status |
|---|---|---|---|
| 0001 | Pilihan stack | 1 | ✅ |
| 0002 | Strategi multi-tenancy | 3 | ⬜ TODO |
| 0002b | Arsitektur hemat biaya (Azure + Neon + Pages) | 4 | ⬜ TODO |
| 0002c | Konvensi Azure (region, naming, tagging) | Day-0 | ✅ |
| 0003 | LLM routing & token budget | 7 | ⬜ TODO |
| 0004 | Compose vs Kubernetes untuk produksi | 11 | ⬜ TODO |
| 0005 | API health-check probe versus BusyBox | 2 | ✅ |

Pakai `template.md` sebagai kerangka.
