# Quiz Verifikasi — Minggu 00

> Template. Salin ke `docs/progress/quiz/week-<NN>.md`.
>
> **Aturan untuk agent:** 8–12 soal. Setiap soal WAJIB bisa dijawab dari kode
> yang benar-benar ada di repo ini — buka filenya dulu sebelum menulis soal.
> Dilarang membuat soal tentang kode yang kamu tulis tapi tidak kamu jalankan.
>
> **Aturan untuk manusia:** jawab tanpa membuka kode. Tanpa mencari. Tanpa
> bertanya ke AI. Skor < 70% berarti minggu ini diulang, bukan dilanjutkan.
> Ini bukan hukuman — ini alat ukur supaya kamu tidak menemukan lubangnya
> saat interview.

**Komposisi soal:** 40% penalaran ("kenapa X") · 30% failure mode ("apa yang terjadi kalau") · 20% orientasi ("di mana X, apa fungsinya") · 10% hafalan

---

## Soal

**1. [Penalaran]** Bagaimana kamu mencegah *cost overrun* di cloud, dan apa perbedaan fungsi spending limit dengan budget alert?

<br>

**2. [Failure mode]** Apa perbedaan menghentikan VM menggunakan `az vm stop` dan melakukan `az vm deallocate`, serta bagaimana keduanya memengaruhi compute billing?

<br>

**3. [Penalaran]** Bagaimana public-key authentication SSH bekerja dan kenapa private key tidak boleh masuk repository?

<br>

**4. [Failure mode]** Bagaimana kamu mengelola secret antara development dan production?

<br>

**5. [Orientasi]** Bagaimana kamu mengelola resource development environment pada mesin dengan RAM terbatas menggunakan WSL2?

<br>

**6. [Penalaran]** Bagaimana kamu menggunakan AI coding agent tanpa kehilangan pemahaman dan kontrol terhadap codebase?

<br>

**7. [Remedial — Failure mode]** Kamu memiliki `GROQ_API_KEY` untuk development di laptop dan aplikasi `campus-lms` nantinya berjalan di server production. Apakah production menggunakan credential development yang sama atau credential tersendiri? Di mana secret production seharusnya berada, bagaimana backend mendapatkannya, dan kenapa secret tersebut tidak boleh berada di Git?

---

## Jawabanmu

| No | Jawaban                                                                                                                                                                                                                                                                                                                                                                                                                 | Yakin? (1-5) |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------ |
| 1  | Dengan konfigurasi spending limit untuk membatasi penggunaan billing agar tagihan tidak melewati kredit yang tersedia, serta menggunakan budget alert pada threshold 50%, 80%, dan 100% yang dikirimkan ke email manusia. Spending limit bertindak sebagai guardrail, sedangkan budget alert berfungsi sebagai peringatan.                                                                                              | 4            |
| 2  | Menghentikan atau `stop` VM berarti hanya mematikan VM, tetapi alokasi compute masih dapat tetap diberikan kepada VM sehingga compute billing masih dapat berjalan. `deallocate` berarti mematikan VM sekaligus melepaskan alokasi compute dari VM sehingga compute billing tidak terus berjalan ketika tidak digunakan.                                                                                                | 5            |
| 3  | Public key ditempatkan di server, sedangkan private key tetap berada di komputer lokal. Client membuktikan bahwa ia memiliki private key yang sesuai dengan public key yang terdaftar sehingga hanya pihak yang memiliki private key yang valid yang dapat melakukan autentikasi. Private key tidak boleh masuk repository karena jika bocor dapat disalahgunakan pihak lain untuk menyamar sebagai pemilik credential. | 5            |
| 4  | Pada development, secret hanya boleh berada pada environment lokal dan tidak boleh dimasukkan ke Git atau di-hard-code karena dapat disalahgunakan. Pada production, secret juga harus dikelola secara terpisah dan tidak berada pada source code. Pemahaman awal saya mengenai mekanisme production masih mencampurkan secret management dengan provisioning, IaC, dan TLS sehingga bagian ini memerlukan remedial.    | 3            |
| 5  | RAM 8 GB merupakan kapasitas terbatas untuk development, sehingga WSL2 perlu dibatasi agar tidak mengambil seluruh resource komputer host dan membuat Windows atau workload lain kekurangan resource. Pembatasannya dilakukan melalui `.wslconfig` dan kemudian diverifikasi dari dalam WSL.                                                                                                                            | 4            |
| 6  | Dengan membaca dan memahami laporan serta perubahan kode/file yang dibuat agent, lalu tidak hanya percaya pada klaim agent tetapi juga menjalankan command verification dan memeriksa evidence untuk memastikan hasilnya benar-benar sesuai.                                                                                                                                                                            | 4            |
| 7  | Production sebaiknya membuat credential tersendiri yang berbeda dari development. Secret production disediakan kepada backend melalui environment atau mekanisme secret management saat runtime. Secret tidak boleh dimasukkan ke Git karena repository dapat terekspos dan credential tersebut dapat disalahgunakan oleh pihak yang tidak berwenang.                                                                   | 4            |

---

## Kunci Jawaban

> Jangan dibuka sebelum menjawab semua.

<details>
<summary>Buka kunci jawaban</summary>

**1.** Spending limit dan budget alert memiliki fungsi berbeda. Spending limit merupakan guardrail terhadap penggunaan kredit/billing pada subscription yang mendukungnya, sedangkan budget alert hanya memberi peringatan ketika penggunaan mencapai threshold tertentu. Budget alert tidak dengan sendirinya menghentikan penggunaan resource.

📁 Konfirmasi di: `docs/progress/week-00.md` Section 6

**2.** `az vm stop` menghentikan VM tetapi compute masih dapat tetap berada dalam keadaan allocated. `az vm deallocate` menghentikan VM sekaligus melepaskan alokasi compute sehingga compute billing berhenti. Deallocate bukan berarti menghapus VM dan resource persisten lain masih dapat tetap ada.

📁 Konfirmasi di: `docs/progress/week-00.md` Section 6

**3.** Public key boleh ditempatkan pada server, sedangkan private key tetap rahasia pada client. Private key tidak dikirim ke server; client membuktikan secara kriptografis bahwa ia memiliki private key yang sesuai dengan public key yang terdaftar. Private key tidak boleh masuk repository karena kebocorannya dapat memungkinkan pihak lain mencoba melakukan autentikasi sebagai pemilik key.

📁 Konfirmasi di: `docs/progress/week-00.md` Section 6

**4.** Secret tidak boleh di-hard-code atau dimasukkan ke source control. Development dan production sebaiknya menggunakan credential yang sesuai dengan masing-masing environment. Pada production, secret diberikan kepada backend melalui environment atau mekanisme secret management pada runtime dengan kontrol akses yang sesuai. Secret management berbeda dengan IaC dan TLS walaupun ketiganya dapat digunakan bersama dalam proses deployment.

📁 Konfirmasi di: `.env.example`, `AGENTS.md`, `agent/policy.md`, `docs/progress/week-00.md`

**5.** `.wslconfig` digunakan untuk membatasi resource WSL2 agar guest environment tidak mengambil seluruh resource host. Pada setup Minggu 0 digunakan konfigurasi memory, processor, dan swap, kemudian hasil runtime diverifikasi menggunakan `free -h` dan `nproc`.

📁 Konfirmasi di: `docs/progress/evidence/week-00/wslconfig.txt`, `docs/progress/evidence/week-00/wsl-resources.txt`

**6.** Output atau laporan agent adalah klaim yang harus diverifikasi. Manusia harus memahami perubahan repo, menjalankan command verification, memeriksa raw output/evidence, dan tidak menerima pernyataan seperti "seharusnya bekerja" tanpa bukti.

📁 Konfirmasi di: `AGENTS.md`, `agent/policy.md`, `docs/progress/evidence/week-00/make-todo.txt`

**7.** Production sebaiknya memiliki credential tersendiri yang terpisah dari development. Secret production tidak disimpan di source code atau Git, tetapi tersedia bagi backend melalui environment atau secret-management mechanism pada runtime. Backend menggunakan secret tersebut untuk mengakses layanan pihak ketiga tanpa mengeksposnya ke frontend atau repository.

📁 Konfirmasi di: `.env.example`, `AGENTS.md`, `agent/policy.md`

</details>

---

## Hasil

* **Skor:** 86 / 100 = 86%
* **Lulus (≥70%)?** ya
* **Soal yang salah:** No. 4 pada jawaban awal memerlukan remedial; remedial No. 7 kemudian lulus
* **Konsep yang perlu diulang:** Pengelolaan secret production, khususnya perbedaan antara secret management, environment/runtime configuration, IaC, dan TLS
* **Rencana perbaikan:** Menguji kembali pemahaman secret management ketika credential production benar-benar diterapkan pada tahap deployment, serta tetap menggunakan pola explain-back mandiri pada quiz minggu berikutnya