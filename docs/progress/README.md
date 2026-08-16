# docs/progress — Rekam Jejak Pembelajaran

Folder ini adalah **bukti bahwa kamu benar-benar memahami sistem ini**, bukan
sekadar memilikinya. Ini juga yang akan kamu tunjukkan saat interview ketika
ditanya "ceritakan proses kamu membangun ini".

## Struktur

```
docs/progress/
  week-00.md .. week-12.md    laporan mingguan (disusun agent, ditandatangani kamu)
  evidence/week-NN/           output mentah perintah — sumber kebenaran semua angka
  quiz/week-NN.md             quiz verifikasi (soal dari agent, jawaban dari kamu)
  explain/week-NN.mp3|txt     rekaman explain-back 3 menit
  sessions/                   session log (opsional, berguna saat debugging panjang)
```

## Aturan gerbang antar-minggu

Minggu berikutnya **tidak dimulai** sebelum keempat hal ini beres:

1. Laporan mingguan ditandatangani
2. Skor quiz ≥ 70%
3. 3 file bukti sudah kamu spot-check secara acak
4. Explain-back 3 menit terekam

Checklist lengkapnya: `agent/checklists/human-verification.md`

## Kenapa serepot ini

Karena kamu memilih mode agent-first. Agent akan menulis sebagian besar kode.
Tanpa mekanisme ini, dalam 12 minggu kamu akan punya repo yang mengesankan dan
pemahaman yang rapuh — dan interview akan membongkarnya dalam 10 menit.

Folder ini adalah harga yang dibayar supaya agent-first tetap menghasilkan
engineer, bukan operator.

## Untuk agent

Baca `agent/evidence-protocol.md` sebelum menulis apa pun di sini.
Tiga aturan yang tidak bisa ditawar:

1. Setiap angka wajib punya file bukti berisi perintah + output mentah + timestamp
2. Section "Belum Terverifikasi" wajib terisi jujur
3. Dilarang mencentang DoD — itu hak manusia
