# Checklist Verifikasi Manusia (akhir tiap minggu)

> Untuk KAMU, bukan agent. ±45 menit. Ini gerbang antar-minggu.

## 1. Spot-check bukti (10 menit)

Pilih **3 file bukti secara acak** dari `docs/progress/evidence/week-01/`:

- [x] Header lengkap (CLAIM, COMMAND, RUN AT, COMMIT, EXIT)
- [x] Perintahnya **saya jalankan ulang sendiri** — hasilnya mirip (tidak harus identik, tapi ordo besarnya sama)
- [x] Angka di laporan cocok dengan angka di file bukti
- [x] Commit SHA-nya benar-benar ada (`git show <sha>`)

**Kalau ada satu saja yang tidak cocok:** hentikan. Minta agent menjelaskan. Ini sinyal paling awal bahwa sistemnya bocor.

## 2. Baca laporan dengan curiga (10 menit)

- [x] Section "Belum Terverifikasi" **terisi dan masuk akal** (kosong = curigai)
- [x] Section "Dikerjakan Manusia" sesuai dengan yang benar-benar saya kerjakan
- [x] Ada minimal satu jalan buntu / kesalahan yang dicatat (minggu tanpa kesalahan itu tidak realistis)
- [x] Setiap angka punya link bukti yang bisa diklik
- [x] Tidak ada klaim DoD yang dicentang agent

## 3. Quiz (15 menit)

- [ ] Dijawab **tanpa** membuka kode, tanpa mencari, tanpa bertanya ke AI
- [x] Skor dicatat. **≥ 70% untuk lanjut**
- [x] Soal yang salah → konsepnya dicatat di "perlu diulang"

## 4. Explain-back (5 menit)

Rekam suara 3 menit menjelaskan **satu** keputusan teknis minggu ini, seolah menjawab interviewer.

- [x] Rekaman tersimpan di `docs/progress/explain/week-NN.<mp3|txt>`
- [x] Saya tidak tersendat lebih dari 5 detik
- [x] Saya menyebutkan trade-off, bukan hanya nama teknologi

**Kalau tersendat:** itu konsep yang belum nyangkut. Kembali ke situ sebelum lanjut.

## 5. Baca kode dengan mata sendiri (5 menit)

- [x] `git log --oneline` minggu ini — saya mengenali setiap commit
- [x] `git diff` pada 2 file terpenting — saya paham setiap perubahan
- [x] Tidak ada kode yang membuat saya berpikir "kok bisa begini ya?" tanpa penjelasan

## 6. Sign-off

- [x] Semua di atas beres → tanda tangani laporan, lanjut minggu berikutnya
- [x] Ada yang gagal → **ulangi bagian itu**. Roadmap 14 minggu dengan pemahaman utuh > 12 minggu dengan repo yang tidak bisa dipertahankan

---

## Tanda bahaya yang harus langsung dihentikan

| Tanda | Artinya |
|---|---|
| Bukti ada tapi perintahnya tidak bisa saya jalankan ulang | Bukti kemungkinan dikarang |
| Angka terlalu bulat dan indah (tepat 50%, tepat 100 ms) | Curigai fabrikasi |
| Agent menyebut file/fungsi yang tidak ada | Konteks kotor — mulai sesi baru |
| Saya tidak paham > 30% kode minggu ini | Terlalu cepat. Pelankan, minta `/teach` |
| Quiz terasa terlalu mudah | Agent membuat soal dangkal — minta soal failure-mode |
