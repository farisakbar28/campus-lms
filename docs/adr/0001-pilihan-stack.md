# ADR-0001: Pilihan Stack untuk Campus LMS

- **Tanggal:** 2026-08-17
- **Status:** Proposed
- **Minggu roadmap:** 1

## Konteks

Campus LMS adalah proyek SaaS Learning Management System multi-tenant yang
dikerjakan oleh satu mahasiswa TI semester 7 dalam kurikulum 12 minggu. Proyek
ini harus menghasilkan aplikasi yang dapat dideploy sekaligus bukti pemahaman
yang dapat dijelaskan saat interview backend/fullstack dengan kekuatan AI.
Target pasar awal adalah Indonesia, lalu pekerjaan remote internasional.

Waktu implementasi sekitar 12 minggu. Anggaran infrastruktur adalah $0 dan
tidak boleh bergantung pada kartu kredit atau debit. Lingkungan produksi utama
adalah VM dengan 1 GB RAM; karena itu komponen yang membutuhkan memori besar
tidak dapat diasumsikan selalu berjalan di VM. Mesin pengembangan adalah laptop
Ryzen 5 dengan RAM 8--16 GB, sehingga stack lokal juga perlu dapat dijalankan
bertahap tanpa membuat laptop terus-menerus kehabisan memori.

Proyek perlu melatih operasi yang biasanya tersembunyi oleh platform managed:
migrasi, query dan index PostgreSQL, backup/restore, container, observability,
serta deployment. Pada saat yang sama, fitur AI membutuhkan ekosistem Python
yang matang untuk RAG, evaluasi, dan orkestrasi agent. ADR ini membandingkan
bahasa API, pemisahan layanan AI, pilihan database, dan tingkat pemecahan
arsitektur dengan batasan tersebut.

## Opsi yang dipertimbangkan

### 1. Bahasa API

#### Go

- **Kelebihan:** bahasa yang sudah paling dikuasai pemilik; standard library
  memadai untuk HTTP; binary dan runtime yang relatif ringan cocok untuk VM
  1 GB; concurrency dan cancellation dapat dipelajari langsung untuk layanan
  backend.
- **Biaya nyata:** perlu waktu belajar idiom Go yang disiplin, terutama
  `context.Context`, error handling, goroutine, dan testing. Ekosistem AI
  tingkat aplikasi lebih terbatas daripada Python, sehingga integrasi AI yang
  kompleks akan lebih sulit bila dipaksakan tetap di Go. Menjalankan API Go
  tidak menambah runtime besar, tetapi tetap berbagi RAM VM dengan proxy,
  Redis, dan proses lain.
- **Yang dipersulit kemudian:** jika seluruh fitur AI ditulis di Go, penggunaan
  library RAG, evaluasi, dan agent yang umumnya Python dapat membutuhkan
  wrapper sendiri atau integrasi API tambahan.

#### Node.js / TypeScript

- **Kelebihan:** satu bahasa dengan frontend; TypeScript memberi type checking
  dan ekosistem web yang luas; waktu awal untuk membangun endpoint CRUD dapat
  lebih singkat bagi developer frontend/fullstack.
- **Biaya nyata:** runtime Node.js dan dependency tree menambah penggunaan RAM
  dan kompleksitas dependency pada VM 1 GB. Waktu belajar dapat bergeser ke
  pengelolaan package, perbedaan module system, serta pola async JavaScript
  yang benar. Mengoperasikan frontend dan API pada ekosistem JavaScript juga
  meningkatkan risiko dependency churn dalam jadwal 12 minggu.
- **Yang dipersulit kemudian:** observasi memory leak dan pengendalian
  dependency production dapat menjadi pekerjaan tambahan; pembelajaran Go
  sebagai kekuatan backend/infra tidak terjadi melalui API utama.

#### Python

- **Kelebihan:** satu ekosistem untuk API dan AI; FastAPI, tooling data, RAG,
  evaluasi, dan agent dapat digunakan tanpa batas layanan antarbahasa.
- **Biaya nyata:** runtime dan dependency AI Python cenderung membutuhkan RAM
  lebih besar daripada API Go sederhana; VM 1 GB membatasi ruang untuk API,
  dependency AI, dan proses pendukung berjalan bersama. Waktu belajar terbagi
  antara pola API produksi Python dan tooling AI. Menginstal package AI juga
  dapat menambah waktu setup dan tekanan storage laptop.
- **Yang dipersulit kemudian:** memisahkan beban AI dari API setelah semuanya
  menyatu membutuhkan kontrak layanan, deployment, dan observability baru di
  tengah proyek.

### 2. Layanan AI

#### Python sebagai layanan terpisah

- **Kelebihan:** API inti dapat tetap ringan, sedangkan Python digunakan hanya
  untuk RAG, embedding, evaluasi, dan orkestrasi agent. Batas HTTP antarproses
  membuat kegagalan dan kebutuhan resource AI dapat diisolasi.
- **Biaya nyata:** ada dua runtime, dua set dependency, kontrak API internal,
  logging/tracing lintas layanan, dan deployment yang harus dipelajari. Python
  service tetap mengonsumsi RAM VM ketika aktif, sehingga pekerjaan batch berat
  perlu dijadwalkan atau dipindahkan ke laptop/layanan eksternal.
- **Yang dipersulit kemudian:** perubahan kontrak antara API dan AI service
  perlu versioning dan integration test; debugging sebuah request dapat
  melewati dua proses.

#### Semua fitur di satu bahasa dan satu proses API

- **Kelebihan:** satu runtime, satu deployment, dan tidak ada network boundary
  internal; alur awal mungkin lebih cepat dipahami.
- **Biaya nyata:** bahasa yang dipilih harus menanggung kompromi: Go/Node akan
  membuat akses tooling AI lebih sulit, sedangkan Python akan membawa beban
  runtime AI ke jalur API utama. Kebocoran memori, dependency besar, atau task
  AI lambat lebih mudah mengganggu request HTTP biasa pada RAM 1 GB.
- **Yang dipersulit kemudian:** pemisahan API dan worker/AI service di masa
  depan mengharuskan ekstraksi kode, desain ulang kontrak, dan perubahan
  deployment ketika fitur sudah bergantung padanya.

### 3. Database

#### PostgreSQL self-managed

- **Kelebihan:** memberi pengalaman langsung tentang migrasi, koneksi, index,
  `EXPLAIN`, backup, restore, serta kegagalan operasi database--materi yang
  relevan untuk target backend.
- **Biaya nyata:** PostgreSQL membutuhkan RAM, storage, backup discipline, dan
  waktu troubleshooting. Menjalankannya bersama API dan layanan pendukung pada
  VM 1 GB membatasi headroom; untuk lokal, Docker database ikut bersaing dengan
  IDE, browser, dan service lain pada laptop 8--16 GB.
- **Yang dipersulit kemudian:** availability, patching, backup off-site, dan
  pemulihan insiden menjadi tanggung jawab proyek bila database production
  tetap self-managed.

#### Supabase

- **Kelebihan:** Postgres managed dengan auth, storage, dan dashboard dapat
  mempercepat fitur aplikasi; pemilik sudah memiliki pengalaman dengannya.
- **Biaya nyata:** sebagian operasi database disembunyikan, sehingga waktu
  belajar migrasi, tuning, backup, dan koneksi production berkurang. Ada
  batasan free tier serta ketergantungan pada produk dan cara operasi vendor;
  penggunaan RAM VM rendah karena database berada di luar VM.
- **Yang dipersulit kemudian:** menjelaskan dan mengoperasikan Postgres tanpa
  abstraksi Supabase akan tetap menjadi gap; migrasi ke provider lain dapat
  melibatkan penggantian layanan di luar database seperti auth atau storage.

#### Neon

- **Kelebihan:** PostgreSQL managed memindahkan beban database dari VM 1 GB
  dan tetap mempertahankan SQL/Postgres sebagai antarmuka utama. Ini memberi
  ruang RAM production untuk API dan proses pendukung, sambil memungkinkan
  PostgreSQL self-managed dipakai di lokal sebagai lingkungan belajar.
- **Biaya nyata:** ada dependency jaringan dan layanan vendor; waktu belajar
  tetap diperlukan untuk memahami connection string, pooling, migrasi, dan
  perbedaan lingkungan lokal-production. Batasan free tier perlu diverifikasi
  sebelum menjadi asumsi production.
- **Yang dipersulit kemudian:** outage atau perubahan batas layanan vendor
  berada di luar kontrol proyek; perpindahan provider memerlukan rencana
  export/import dan validasi koneksi aplikasi.

### 4. Arsitektur aplikasi

#### Modular monolith

- **Kelebihan:** satu codebase dan satu deployment untuk domain LMS inti,
  dengan batas modul yang dapat dipelajari dan diuji tanpa biaya jaringan antar
  layanan. Cocok untuk solo developer dan jadwal 12 minggu. Layanan AI dapat
  tetap menjadi proses terpisah bila memang membutuhkan Python.
- **Biaya nyata:** disiplin batas modul harus dijaga secara manual; satu proses
  API masih berbagi resource dengan fitur domain lain. Waktu belajar berfokus
  pada desain modul, testing, dan deployment yang sederhana, bukan pada
  orkestrasi banyak layanan. Penggunaan RAM lebih rendah daripada banyak
  service untuk domain yang sama, tetapi tidak menghilangkan kebutuhan RAM
  layanan AI atau database lokal.
- **Yang dipersulit kemudian:** bila domain tumbuh jauh lebih besar, ekstraksi
  modul menjadi service terpisah membutuhkan kontrak dan migrasi ownership
  data yang belum ada sejak awal.

#### Microservices

- **Kelebihan:** setiap domain dapat dideploy, diskalakan, dan diisolasi secara
  mandiri; batas antar layanan dipaksa melalui API.
- **Biaya nyata:** banyak proses berarti RAM baseline lebih tinggi--berisiko
  pada VM 1 GB--serta membutuhkan service discovery, kontrak API, autentikasi
  antarservice, observability terdistribusi, retry, dan deployment per layanan.
  Waktu belajar akan habis pada plumbing infrastruktur sebelum fitur LMS dan
  AI memiliki nilai yang dapat didemonstrasikan.
- **Yang dipersulit kemudian:** debugging, pengujian end-to-end, transaksi
  lintas domain, dan perubahan skema menjadi lebih mahal karena setiap
  perubahan dapat melibatkan beberapa repository atau deployment.

## Keputusan
- **Keputusan Menggunakan GO sebagai Bahasa API:** Go adalah bahasa backend yang dikenal efisien, memiliki performa tinggi dengan native binary yang tidak dimiliki bahasa interpreted seperti Python/Node.js, Go memiliki tingkat concurrency yang sangat tinggi ketika melayani banyak request dalam 1 waktu bersamaan, hal-hal tersebut sangat cocok untuk project LMS SaaS yang sedang dikembangkan. bayangkan jika ada 1000 tenant campus yang sedang melakukan request absensi di jam 08:00 pagi secara bersamaan, melalui penggunaan concurrency goroutine dan channel, backend API mampu memecah request yang masuk tanpa membuat request lainnya menunggu.
- **Keputusan Pemisahan Antara API dan AI Service:** Go memang unggul jika digunakan untuk API service, tapi tidak dengan AI service. keputusan pemisahan antara keduanya didasari pada proses komputasi yang bisa sangat berat jika keduanya digabungkan, dan keduanya memang memiliki perannya masing-masing. jika API service ada kendala, AI service tidak terkena dampaknya, begitupun sebaliknya. maka, proyek ini dibangung dengan Python sebagai bahasa yang digunakan untuk AI service. Python memang dikenal umum sebagai bahasa paling cocok ketika kita membicarakan AI, hal ini sangat berkaitan dengan luasnya dukungan framework/dependency Python untuk mengembangkan model AI atau semacamnya.
- **Keputusan Menggunakan Postgresql Self-managed:** Keputusan inti dari ini adalah karena melalui penggunaan Postgresql Self-managed dapat memberikan kuasa yang jauh lebih dalam untuk proyek, terutama jika proyek sudah jauh lebih besar. dengan ini, developer atau tim data analyst dapat lebih leluasa dalam pengelolaan database atau migrasi. kenapa tidak supabase yang sudah langsung menyiapkan dan memberikan UI yang mudah? supabase itu sendiri berbasis Postgresql, dan supabase tidak memberikan akses yang leluasa untuk developer/tim mendalami lebih jauh kedalam mekanisme database itu sendiri, lalu juga akan lebih sulit ketika tim ingin melakukan migrasi karena beberapa hal yang diatur supabase mungkin tidak relevan dengan tools lain. disinilah Neon hadir, Neon mampu memberikan akses penuh terhadap database proyek yang dikembangkan dan dapat self-managed, hal ini dapat mengurangi beban komputer tim/developer dengan spesifikasi rendah (saat ini developer utama menggunakan VM 1GB).
- **Keputusan Menggunakan Arsitektur Modular Monolith:** Penggunaan arsitektur ini dipilih karena mengutumakan skala kebutuhan developer untuk proses pengembangan project latihan. arsitektur lebih mudah diimplementasikan karena lebih mudah dipelajari dan tidak ada biaya antar jaringan layanan seperti pada Microservices. Microservices memang unggul jika untuk project dengan skala besar dan pengguna besar, karena setiap layanan berjalan masing masing dan tidak berbagi resource satu sama lain, tapi untuk skala pembelajaran, arsitektur modular monolith jauh lebih mudah diimplementasikan.

## Konsekuensi
konsenkuensi dari keputusan diatas murni pada 1 hal utama, yaitu skala project untuk latihan mandiri developer. seperti memisahkan API dan AI yang jika keduanya disatukan, hal tersebut bisa justru akan menganggu satu sama lain jika terdapat masalah disalahsatunya, karena kedua layanan ini fokus yang berbeda dan tidak harus satu sama lainnya mengganggu. konsekuensinya adalah integrasi kedua layanan dapat jadi lebih kompleks karena keduanya dibangun menggunakan 2 bahasa yang berbeda. kalau untuk konsekuensi dari keputusan Postgresql Self-manage, ada pada tingkat pemahaman yang ingin dicapai oleh developer, melalui penggunaan ini developer dapat lebih leluasa dalam mengelola ekosistem database secara penuh, konsekuensinya mungkin akan jadi lebih kompleks karena database di manage sendiri oleh developer. kalau untuk pemilihan arsitektur, konsekuensinya adaah jika skala project ini domainnya menjadi besar, akan lebih sulit untuk dikelola jika dibandingkan arsitektur modern seperti microservices jauh dapat lebih fleksibel untuk pengelolaan per layanan tanpa mengganggu layanan lainnya.

## Catatan
Pemahaman diatas adalah pemahaman developer yang mungkin masih terlalu abstak atau masih banyak kurangnya penjelasan yang detail, hal ini murni karena developer masih sedang mempelajari dan mendalami hal-hal tersebut.