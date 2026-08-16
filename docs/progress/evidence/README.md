# Evidence — Bukti Eksekusi

Setiap file di sini adalah output mentah sebuah perintah, tidak diedit.

Format wajib (lihat `agent/evidence-protocol.md`):

```
=== EVIDENCE ===
CLAIM:    <klaim satu kalimat>
COMMAND:  <perintah persis>
CWD:      <working directory>
RUN AT:   <ISO-8601 timestamp>
COMMIT:   <git sha>
EXIT:     <exit code>
=== RAW OUTPUT ===
<output apa adanya — tidak dipotong, tidak dirapikan>
=== END ===
```

**Kegagalan juga bukti.** Perintah yang gagal disimpan dengan `EXIT: 1`,
tidak disembunyikan. Justru itu yang paling berguna saat kamu menelusuri
kenapa sesuatu rusak dua minggu kemudian.

Simpan satu klaim per file. Nama file deskriptif:
`image-size.txt`, `explain-analyze-course-list-before.txt`, `k6-500vu-run1.txt`.
