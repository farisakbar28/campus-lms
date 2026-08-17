# ADR-0005: API binary health check versus BusyBox probe

- **Date:** 2026-08-18
- **Status:** Proposed
- **Roadmap week:** 2

## Context

The Week 2 Dockerfile requires both a distroless runtime and an HTTP
`HEALTHCHECK`. Distroless intentionally contains no shell or HTTP client.
The initial implementation copied BusyBox into the final runtime and invoked
its `wget` applet. That made the check work, but also introduced a multi-call
binary with utilities that distroless intentionally omits.

The implementation now provides an API binary mode, `/api -healthcheck`, that
performs a bounded HTTP GET to `http://127.0.0.1:$APP_PORT/healthz` and exits
zero only for HTTP 200. The BusyBox implementation remains available as the
`busybox-runtime` Docker target solely for a measured comparison.

## Options considered

1. **Copy BusyBox and run `wget` in the final distroless image**
   - Provides an HTTP client without changing application code.
   - `docker images` reports a final size of 16.6MB.
   - The copied `/busybox` file is 1,214,736 bytes. Docker history reports its
     COPY layer as 1.22MB.
   - The final image reports 4,054,678 bytes through `docker image inspect`.
   - It adds a multi-call binary, including utilities not needed by the API.

2. **Use an API binary health-check mode**
   - Keeps the final runtime to the API binary and the distroless base.
   - `docker images` reports a final size of 14.6MB.
   - The final image reports 3,285,268 bytes through `docker image inspect`.
   - This is 769,410 bytes smaller than the BusyBox runtime by that Docker
     image-size measurement.
   - It adds a small HTTP client code path that requires unit tests and must
     stay compatible with the health endpoint.

3. **Use `golang:alpine` as the production runtime**
   - Makes interactive debugging easier because a shell and utilities exist.
   - A previous local comparison measured 377MB, which is outside the API
     image-size target.

4. **Omit Docker HEALTHCHECK**
   - Avoids adding any probe utility or code path.
   - Does not meet the Week 2 task brief.

## Evidence pending human decision

- [BusyBox versus API-probe exact image sizes](../progress/evidence/week-02/busybox-vs-api-probe-size.txt)
- [BusyBox and API-probe layer histories](../progress/evidence/week-02/busybox-layer-history.txt)
- [API-probe runtime health check](../progress/evidence/week-02/api-probe-healthcheck.txt)
- [Earlier Alpine versus distroless comparison](../progress/evidence/week-02/image-size.txt)

This ADR intentionally has no Decision or Consequences section. The human
owner chooses and records those sections after reviewing the evidence.

## Decision

* **Keputusan Menggunakan API Binary Health-check Mode:** Keputusan final adalah menggunakan mode health check langsung melalui binary API dengan flag `-healthcheck`, bukan dengan menambahkan BusyBox kedalam image distroless. keputusan ini didasari dari hasil pengukuran yang menunjukkan bahwa penggunaan distroless dapat mengurangi ukuran runtime dari alpine yang sebelumnya mencapai 377 MB menjadi sekitar 15,9 MB atau berkurang sekitar 95,8%. ketika BusyBox ditambahkan hanya untuk menyediakan `wget` sebagai HTTP probe, terdapat tambahan layer sebesar 1,22 MB, padahal fungsi yang benar-benar dibutuhkan hanya melakukan request ke endpoint `/healthz`. dengan menambahkan mode `/api -healthcheck` langsung pada binary aplikasi, Docker `HEALTHCHECK` tetap dapat melakukan pengecekan HTTP tanpa membutuhkan shell ataupun binary tambahan, dan solusi tersebut tidak menambah ukuran image secara terukur atau +0 MB dibandingkan runtime distroless final. keputusan ini juga mempertahankan tujuan utama penggunaan distroless, yaitu membuat runtime sekecil mungkin dan hanya membawa komponen yang memang dibutuhkan oleh aplikasi.

## Consequences

konsekuensi utama dari keputusan menggunakan API binary sebagai health check adalah aplikasi sekarang memiliki satu code path tambahan yang khusus digunakan untuk melakukan HTTP GET ke endpoint `/healthz`. code path tersebut harus tetap diuji dan dijaga agar sesuai dengan perilaku health endpoint, karena jika implementasi health check atau endpoint berubah tanpa diperbarui bersama, Docker dapat menganggap container tidak sehat walaupun proses API masih berjalan. tetapi konsekuensi tersebut lebih kecil dibandingkan menambahkan BusyBox kedalam runtime hanya untuk mendapatkan `wget`, karena BusyBox membawa banyak utility lain yang sebenarnya tidak dibutuhkan oleh API dan menambahkan sekitar 1,22 MB kedalam image.

konsekuensi lainnya berasal dari penggunaan distroless itu sendiri. image distroless tidak menyediakan shell dan utility debugging seperti yang biasanya tersedia pada alpine atau image Linux umum lainnya. akibatnya developer tidak dapat masuk kedalam container menggunakan cara debugging tradisional seperti menjalankan `/bin/sh` lalu memeriksa filesystem atau menjalankan command secara langsung. jika membutuhkan debugging pada container distroless, developer harus menggunakan mekanisme seperti `docker debug` atau menggunakan ephemeral container yang memiliki debugging tools. ini adalah harga yang diterima untuk mendapatkan runtime yang jauh lebih kecil, yaitu sekitar 15,9 MB dibandingkan alpine runtime yang sebelumnya mencapai 377 MB, sekaligus mengurangi jumlah executable dan utility yang tidak diperlukan sehingga permukaan serangan pada production runtime juga menjadi lebih kecil.

## Catatan

Keputusan ini dibuat berdasarkan hasil pengukuran image dan layer yang sudah dilakukan pada Week 2, bukan hanya berdasarkan asumsi bahwa distroless selalu lebih baik. BusyBox sebenarnya dapat menyelesaikan masalah HTTP `HEALTHCHECK`, tetapi setelah dibandingkan, menambahkan utility tersebut hanya untuk satu fungsi probe dianggap tidak sebanding dengan tambahan komponen didalam runtime. karena binary Go yang sama sudah dapat melakukan health check sendiri melalui flag `-healthcheck`, solusi tersebut dipilih sebagai implementasi final.