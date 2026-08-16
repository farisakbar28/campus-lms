# Checklist Verifikasi Manusia (akhir tiap minggu)

> Untuk KAMU, bukan agent. ±45 menit. Ini gerbang antar-minggu.

## 1. Spot-check bukti (10 menit)

Pilih **3 file bukti secara acak** dari `docs/progress/evidence/week-NN/`:

- [ ] Header lengkap (CLAIM, COMMAND, RUN AT, COMMIT, EXIT)
- [ ] Perintahnya **saya jalankan ulang sendiri** — hasilnya mirip (tidak harus identik, tapi ordo besarnya sama)
- [ ] Angka di laporan cocok dengan angka di file bukti
- [ ] Commit SHA-nya benar-benar ada (`git show <sha>`)

**Kalau ada satu saja yang tidak cocok:** hentikan. Minta agent menjelaskan. Ini sinyal paling awal bahwa sistemnya bocor.

## 2. Baca laporan dengan curiga (10 menit)

- [ ] Section "Belum Terverifikasi" **terisi dan masuk akal** (kosong = curigai)
- [ ] Section "Dikerjakan Manusia" sesuai dengan yang benar-benar saya kerjakan
- [ ] Ada minimal satu jalan buntu / kesalahan yang dicatat (minggu tanpa kesalahan itu tidak realistis)
- [ ] Setiap angka punya link bukti yang bisa diklik
- [ ] Tidak ada klaim DoD yang dicentang agent

## 3. Quiz (15 menit)

- [ ] Dijawab **tanpa** membuka kode, tanpa mencari, tanpa bertanya ke AI
- [ ] Skor dicatat. **≥ 70% untuk lanjut**
- [ ] Soal yang salah → konsepnya dicatat di "perlu diulang"

## 4. Explain-back (5 menit)

Rekam suara 3 menit menjelaskan **satu** keputusan teknis minggu ini, seolah menjawab interviewer.

- [ ] Rekaman tersimpan di `docs/progress/explain/week-NN.<mp3|txt>`
- [ ] Saya tidak tersendat lebih dari 5 detik
- [ ] Saya menyebutkan trade-off, bukan hanya nama teknologi

**Kalau tersendat:** itu konsep yang belum nyangkut. Kembali ke situ sebelum lanjut.

## 5. Baca kode dengan mata sendiri (5 menit)

- [ ] `git log --oneline` minggu ini — saya mengenali setiap commit
- [ ] `git diff` pada 2 file terpenting — saya paham setiap perubahan
- [ ] Tidak ada kode yang membuat saya berpikir "kok bisa begini ya?" tanpa penjelasan

## 6. Sign-off

- [ ] Semua di atas beres → tanda tangani laporan, lanjut minggu berikutnya
- [ ] Ada yang gagal → **ulangi bagian itu**. Roadmap 14 minggu dengan pemahaman utuh > 12 minggu dengan repo yang tidak bisa dipertahankan

---

## Tanda bahaya yang harus langsung dihentikan

| Tanda | Artinya |
|---|---|
| Bukti ada tapi perintahnya tidak bisa saya jalankan ulang | Bukti kemungkinan dikarang |
| Angka terlalu bulat dan indah (tepat 50%, tepat 100 ms) | Curigai fabrikasi |
| Agent menyebut file/fungsi yang tidak ada | Konteks kotor — mulai sesi baru |
| Saya tidak paham > 30% kode minggu ini | Terlalu cepat. Pelankan, minta `/teach` |
| Quiz terasa terlalu mudah | Agent membuat soal dangkal — minta soal failure-mode |
