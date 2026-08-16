# Runbook Insiden

**MINGGU 4 — TULIS SENDIRI, lalu LATIH.**

Runbook yang tidak pernah dilatih hanyalah dokumen. Untuk tiap skenario,
kamu harus benar-benar menyebabkannya dengan sengaja, memperbaikinya
sambil mengikuti runbook, lalu mencatat MTTR-nya.

## Format tiap skenario

```
### <Gejala yang terlihat>
Dampak      : siapa yang terganggu, seberapa parah
Deteksi     : dari mana kamu tahu (alert apa? dashboard mana?)
Diagnosis   : perintah yang dijalankan, urutannya
Mitigasi    : langkah cepat menghentikan pendarahan
Perbaikan   : solusi permanen
Pencegahan  : apa yang diubah supaya tidak terulang
MTTR        : <diisi setelah latihan sungguhan>
```

## Skenario wajib (minimal 6)

- [ ] Situs mengembalikan 502 Bad Gateway
- [ ] Disk VM penuh
- [ ] Container dibunuh OOM-killer (sangat mungkin di VM 1 GB)
- [ ] Database tidak bisa diakses / connection pool habis
- [ ] Sertifikat TLS gagal diperbarui
- [ ] Deploy gagal di tengah jalan
- [ ] Kredit Azure menipis / VM berhenti
- [ ] Provider LLM down atau kuota harian habis (Minggu 7)
- [ ] Agent melakukan aksi tidak diinginkan (Minggu 10)
