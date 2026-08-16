# Pre-Commit Checklist

> Untuk agent sebelum mengusulkan commit, dan untukmu sebelum merge.

## Agent

- [ ] `make lint` bersih
- [ ] `make test` hijau — output tersimpan sebagai bukti
- [ ] Tidak ada secret, token, atau data pribadi di diff
- [ ] Tidak ada `fmt.Println` / `print()` / `console.log` sisa debugging
- [ ] Semua error di-wrap dengan konteks, tidak ada yang ditelan
- [ ] Fungsi baru punya test kalau membawa perilaku baru
- [ ] Komentar menjelaskan **kenapa**, bukan **apa**
- [ ] Commit message format Conventional Commits, satu perubahan logis
- [ ] Bukti tersimpan di `docs/progress/evidence/week-NN/`

## Manusia sebelum merge

- [ ] Saya membaca seluruh diff, bukan hanya judul PR
- [ ] Saya bisa menjelaskan setiap file yang berubah
- [ ] Tidak ada baris yang saya tidak paham (kalau ada: tanya `/teach` dulu)
- [ ] CI hijau
- [ ] Perubahan ini muncul di laporan mingguan
