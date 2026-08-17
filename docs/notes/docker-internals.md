# Docker Internals — Process, Namespace, Cgroup, Filesystem, Network, dan Perbedaan Container vs VM

## Tujuan

Catatan ini mendokumentasikan eksperimen langsung untuk memahami cara kerja Docker container di lingkungan Linux/WSL.

Eksperimen difokuskan pada beberapa pertanyaan dasar:

* Apakah container sebenarnya merupakan process Linux biasa?
* Bagaimana container diisolasi dari host?
* Apa fungsi Linux namespace?
* Apa fungsi cgroup?
* Bagaimana Docker membatasi memory dan CPU?
* Bagaimana filesystem container terlihat dari dalam container?
* Apakah container mempunyai network interface sendiri?
* Apa perbedaan utama antara container dan virtual machine?

Mental model utama yang ingin dibuktikan adalah:

```text
Container
│
├── menjalankan process Linux biasa
├── memakai namespace untuk isolasi
├── memakai cgroup untuk pembatasan resource
├── memiliki filesystem view sendiri
└── memiliki network view sendiri
```

---

# 1. Membuat Container Percobaan

Container percobaan dibuat menggunakan image Alpine:

```bash
docker run -d --name probe alpine sleep 300
```

Output:

```text
e580555780bb8f06b3c80dd9c3e39d1556dc8793de59daba0d3d8bdda78d7269
```

Perintah tersebut berarti:

* `docker run` membuat dan menjalankan container baru;
* `-d` menjalankan container di background;
* `--name probe` memberi nama container `probe`;
* `alpine` adalah image yang digunakan;
* `sleep 300` adalah command yang dijalankan di dalam container.

Command `sleep 300` sengaja digunakan agar container tetap hidup selama 300 detik sehingga process-nya dapat diperiksa dari host.

Secara sederhana:

```text
Alpine image
     │
     ↓
Docker membuat container "probe"
     │
     ↓
menjalankan process
     │
     ↓
sleep 300
```

---

# 2. Mengambil PID Process Utama Container

PID process utama container diambil dengan:

```bash
PID=$(docker inspect -f '{{.State.Pid}}' probe)
```

Command substitution:

```bash
$(...)
```

menjalankan command di dalamnya terlebih dahulu.

Sedangkan:

```bash
docker inspect -f '{{.State.Pid}}' probe
```

meminta Docker menampilkan PID process utama container `probe` dari sudut pandang host.

Nilai tersebut kemudian disimpan ke shell variable:

```text
PID
```

Sehingga command berikutnya dapat menggunakan:

```bash
$PID
```

daripada berulang kali menjalankan `docker inspect`.

---

# 3. Container Tetap Merupakan Process Linux Biasa

Untuk melihat process tersebut dari host digunakan:

```bash
ps -p $PID -o pid,ppid,comm
```

Hasil:

```text
  PID  PPID COMMAND
11945 11921 sleep
```

Informasi tersebut menunjukkan:

```text
PID     = 11945
PPID    = 11921
COMMAND = sleep
```

## PID

`PID` adalah **Process ID**, yaitu nomor identitas process di Linux.

Dalam eksperimen ini:

```text
PID 11945
```

adalah process `sleep` milik container dari sudut pandang host.

## PPID

`PPID` adalah **Parent Process ID**.

Hasil:

```text
PPID 11921
```

berarti process `sleep` dengan PID `11945` mempunyai parent process dengan PID `11921`.

Hal ini memperlihatkan bahwa process container tetap berada di process tree Linux host.

Dengan kata lain, container bukan komputer baru yang sepenuhnya terpisah.

Dari sisi host, container tetap mempunyai process nyata:

```text
Linux / WSL environment
│
├── process lain
├── process lain
│
└── PID 11921
    │
    └── PID 11945 → sleep
                      ↑
                      process container "probe"
```

---

# 4. Memeriksa Jumlah Process `sleep`

Eksperimen berikut dijalankan:

```bash
ps aux | grep -c sleep
```

Hasil:

```text
2
```

Artinya pada saat command tersebut dijalankan, output `ps aux` memiliki dua baris yang mengandung teks `sleep`.

Eksperimen ini membantu menunjukkan bahwa process `sleep` milik container dapat terlihat dari sisi host sebagai process Linux biasa.

Namun angka `2` sendiri tidak digunakan sebagai bukti bahwa kedua process tersebut pasti berasal dari Docker. Untuk mengetahui asal masing-masing process secara pasti, PID dan parent process-nya perlu diperiksa secara individual.

Poin pentingnya adalah:

> Process di dalam container tetap mempunyai representasi process nyata pada Linux host.

---

# 5. Melihat Namespace Container

Namespace process utama container diperiksa menggunakan:

```bash
sudo ls -l /proc/$PID/ns/
```

Hasil:

```text
total 0
lrwxrwxrwx 1 root root 0 Aug 17 23:56 cgroup -> 'cgroup:[4026532234]'
lrwxrwxrwx 1 root root 0 Aug 17 23:56 ipc -> 'ipc:[4026532223]'
lrwxrwxrwx 1 root root 0 Aug 17 23:56 mnt -> 'mnt:[4026532221]'
lrwxrwxrwx 1 root root 0 Aug 17 23:56 net -> 'net:[4026532235]'
lrwxrwxrwx 1 root root 0 Aug 17 23:56 pid -> 'pid:[4026532233]'
lrwxrwxrwx 1 root root 0 Aug 17 23:56 pid_for_children -> 'pid:[4026532233]'
lrwxrwxrwx 1 root root 0 Aug 17 23:56 time -> 'time:[4026532315]'
lrwxrwxrwx 1 root root 0 Aug 17 23:56 time_for_children -> 'time:[4026532315]'
lrwxrwxrwx 1 root root 0 Aug 17 23:56 user -> 'user:[4026531837]'
lrwxrwxrwx 1 root root 0 Aug 17 23:56 uts -> 'uts:[4026532222]'
```

Direktori:

```text
/proc/<PID>/ns/
```

menampilkan namespace yang digunakan oleh sebuah process.

Namespace merupakan salah satu fitur Linux kernel yang dipakai Docker untuk memberikan **isolasi pandangan sistem** kepada container.

Mental model sederhananya:

> Namespace menentukan bagian sistem apa yang terlihat oleh sebuah process.

---

# 6. Membandingkan Namespace Container dengan Process PID 1

Untuk membuktikan bahwa process container berada pada namespace yang berbeda dari process utama environment Linux, namespace PID `1` juga diperiksa:

```bash
sudo ls -l /proc/1/ns/
```

Hasil:

```text
total 0
lrwxrwxrwx 1 root root 0 Aug 17 23:57 cgroup -> 'cgroup:[4026531835]'
lrwxrwxrwx 1 root root 0 Aug 17 23:57 ipc -> 'ipc:[4026532206]'
lrwxrwxrwx 1 root root 0 Aug 17 23:57 mnt -> 'mnt:[4026532217]'
lrwxrwxrwx 1 root root 0 Aug 17 23:57 net -> 'net:[4026531840]'
lrwxrwxrwx 1 root root 0 Aug 17 23:57 pid -> 'pid:[4026532219]'
lrwxrwxrwx 1 root root 0 Aug 17 23:57 pid_for_children -> 'pid:[4026532219]'
lrwxrwxrwx 1 root root 0 Aug 17 23:57 time -> 'time:[4026531834]'
lrwxrwxrwx 1 root root 0 Aug 17 23:57 time_for_children -> 'time:[4026531834]'
lrwxrwxrwx 1 root root 0 Aug 17 23:04 user -> 'user:[4026531837]'
lrwxrwxrwx 1 root root 0 Aug 17 23:57 uts -> 'uts:[4026532218]'
```

Jika dibandingkan, banyak namespace identifier-nya berbeda.

| Namespace | Container `probe` |        PID 1 |
| --------- | ----------------: | -----------: |
| cgroup    |      `4026532234` | `4026531835` |
| ipc       |      `4026532223` | `4026532206` |
| mnt       |      `4026532221` | `4026532217` |
| net       |      `4026532235` | `4026531840` |
| pid       |      `4026532233` | `4026532219` |
| time      |      `4026532315` | `4026531834` |
| user      |      `4026531837` | `4026531837` |
| uts       |      `4026532222` | `4026532218` |

Perbedaan identifier menunjukkan bahwa process container dan process PID 1 berada pada namespace yang berbeda untuk banyak kategori.

Contohnya:

```text
Container net namespace
net:[4026532235]

PID 1 net namespace
net:[4026531840]
```

Karena identifier berbeda, keduanya berada pada network namespace yang berbeda.

Hal serupa terlihat pada:

* PID namespace;
* mount namespace;
* IPC namespace;
* UTS namespace;
* cgroup namespace;
* time namespace.

Pada eksperimen ini, `user` namespace mempunyai identifier yang sama:

```text
user:[4026531837]
```

Artinya, dalam konfigurasi yang sedang digunakan, process container dan PID 1 berada pada user namespace yang sama.

Jadi tidak benar jika diasumsikan bahwa **setiap jenis namespace pasti selalu berbeda** untuk setiap container. Hal tersebut bergantung pada konfigurasi runtime dan host.

---

# 7. Jenis-Jenis Namespace yang Terlihat

## PID Namespace

Entry:

```text
pid
```

berhubungan dengan isolasi process ID.

PID namespace membuat container dapat mempunyai pandangan process sendiri.

Sebuah process yang mempunyai PID tertentu dari sisi host dapat mempunyai PID yang berbeda ketika dilihat dari dalam namespace container.

Konsep sederhananya:

```text
Host melihat:
PID 11945 → sleep

Container:
mempunyai process view sendiri
```

Jadi container tidak harus melihat seluruh process milik host.

---

## Mount Namespace

Entry:

```text
mnt
```

mengisolasi filesystem mount yang terlihat oleh process.

Dengan mount namespace, container dapat mempunyai pandangan filesystem sendiri tanpa harus melihat seluruh filesystem host.

---

## Network Namespace

Entry:

```text
net
```

mengisolasi networking.

Container dapat mempunyai:

* network interface sendiri;
* IP address sendiri;
* routing table sendiri;
* loopback interface sendiri.

Eksperimen network akan dibahas pada bagian berikutnya.

---

## UTS Namespace

Entry:

```text
uts
```

mengisolasi hostname dan domain name.

Dengan UTS namespace, sebuah container dapat mempunyai hostname sendiri tanpa mengubah hostname host.

---

## IPC Namespace

Entry:

```text
ipc
```

mengisolasi beberapa mekanisme **Inter-Process Communication**, misalnya:

* shared memory;
* message queues;
* semaphore tertentu.

---

## User Namespace

Entry:

```text
user
```

berkaitan dengan isolasi dan mapping:

* user ID;
* group ID.

Dalam eksperimen ini, user namespace process container dan PID 1 memiliki identifier yang sama.

---

## Cgroup Namespace

Entry:

```text
cgroup
```

mengisolasi pandangan process terhadap hierarchy cgroup.

Perlu dibedakan:

```text
cgroup namespace
```

dan:

```text
cgroup resource control
```

bukan konsep yang persis sama.

Cgroup namespace mengatur **pandangan process terhadap hierarchy cgroup**, sedangkan cgroup sendiri digunakan Linux untuk mengelompokkan dan mengontrol resource process.

---

## Time Namespace

Entry:

```text
time
```

berhubungan dengan isolasi terhadap clock tertentu yang disediakan Linux.

---

# 8. Mental Model Namespace

Cara sederhana memahami namespace adalah:

```text
Process container
       │
       ↓
Linux namespace
       │
       ├── "process apa yang bisa saya lihat?"
       ├── "filesystem apa yang bisa saya lihat?"
       ├── "network apa yang bisa saya lihat?"
       └── "hostname apa yang saya lihat?"
```

Jadi:

> **Namespace terutama memberikan isolasi terhadap pandangan sebuah process terhadap sistem.**

---

# 9. Membuktikan Pembatasan Memory dengan Cgroup

Container sementara dijalankan dengan batas memory:

```bash
docker run --rm --memory=64m alpine cat /sys/fs/cgroup/memory.max
```

Hasil:

```text
67108864
```

Docker diminta memberikan batas memory:

```text
64 MiB
```

Dalam byte:

```text
64 × 1024 × 1024
= 67,108,864
```

Sehingga:

```text
67108864 byte = 64 MiB
```

Nilai:

```text
/sys/fs/cgroup/memory.max
```

menunjukkan batas maksimum memory yang diterapkan oleh cgroup pada container tersebut.

Eksperimen ini membuktikan bahwa opsi Docker:

```bash
--memory=64m
```

diterjemahkan menjadi pembatasan resource pada Linux cgroup.

---

# 10. Membuktikan Pembatasan CPU dengan Cgroup

Eksperimen berikut dijalankan:

```bash
docker run --rm --cpus=0.5 alpine cat /sys/fs/cgroup/cpu.max
```

Hasil:

```text
50000 100000
```

Pada cgroup v2, `cpu.max` memiliki bentuk dasar:

```text
quota period
```

Dalam eksperimen:

```text
quota  = 50000
period = 100000
```

Perbandingannya:

```text
50000 / 100000
= 0.5
```

Hal tersebut sesuai dengan konfigurasi Docker:

```bash
--cpus=0.5
```

Artinya container diberi quota CPU setara dengan sekitar setengah dari satu CPU secara rata-rata dalam periode scheduler tersebut.

Mental model sederhananya:

```text
Host memiliki CPU
       │
       ↓
Docker --cpus=0.5
       │
       ↓
cgroup cpu.max
       │
       ↓
50000 100000
```

Jadi cgroup tidak hanya dapat membatasi memory, tetapi juga CPU.

---

# 11. Namespace dan Cgroup Memiliki Fungsi Berbeda

Dua konsep ini sangat penting dan tidak boleh tertukar.

## Namespace

Menjawab pertanyaan:

> Apa yang boleh dilihat oleh process?

Contohnya:

* process;
* network;
* filesystem mount;
* hostname.

## Cgroup

Menjawab pertanyaan:

> Berapa banyak resource yang boleh digunakan oleh process?

Contohnya:

* memory;
* CPU;
* I/O;
* jumlah process.

Mental model:

```text
Docker Container
│
├── Namespace
│   └── isolasi pandangan sistem
│
└── Cgroup
    └── pembatasan dan pengelolaan resource
```

Atau lebih singkat:

```text
Namespace → isolation
Cgroup    → resource control
```

---

# 12. Melihat Filesystem dari Dalam Container

Eksperimen berikut menjalankan container Alpine sementara:

```bash
docker run --rm alpine ls /
```

Hasil:

```text
bin
dev
etc
home
lib
media
mnt
opt
proc
root
run
sbin
srv
sys
tmp
usr
var
```

Dari dalam container, root directory `/` terlihat seperti filesystem Linux tersendiri:

```text
/
├── bin
├── dev
├── etc
├── home
├── lib
├── proc
├── root
├── sys
├── tmp
├── usr
└── var
```

Namun ini tidak berarti container mempunyai disk fisik atau operating system kernel sendiri.

Docker memberikan **filesystem view** yang sesuai dengan image dan mount container.

Dalam kasus ini, filesystem dasarnya berasal dari image:

```text
alpine
```

Mount namespace membantu membuat process di dalam container melihat filesystem tersebut sebagai dunianya sendiri.

Mental model:

```text
Host filesystem
│
│    Docker + mount namespace
│              ↓
│      Container filesystem view
│              │
│              ├── /bin
│              ├── /etc
│              ├── /proc
│              ├── /usr
│              └── /var
```

---

# 13. Melihat Network dari Dalam Container

Eksperimen berikut dijalankan:

```bash
docker run --rm alpine ip addr
```

Hasil:

```text
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever

2: eth0@if29: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP
    link/ether aa:99:6e:e9:bc:05 brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.3/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
```

Container memiliki setidaknya dua interface yang terlihat:

```text
lo
eth0
```

## Loopback Interface

```text
lo
```

mempunyai alamat:

```text
127.0.0.1
```

Ini adalah loopback interface dari sudut pandang network namespace container.

## Ethernet Interface

Container juga melihat:

```text
eth0
```

dengan alamat IPv4:

```text
172.17.0.3/16
```

Hal ini menunjukkan bahwa container mempunyai network view sendiri.

Secara sederhana:

```text
Container
│
├── lo
│   └── 127.0.0.1
│
└── eth0
    └── 172.17.0.3
```

Network interface tersebut tidak berarti container mempunyai kartu jaringan fisik sendiri.

Docker dan Linux menyediakan virtual networking sehingga process di dalam container dapat melihat interface jaringan seolah-olah merupakan lingkungan networking tersendiri.

Ini merupakan salah satu bukti nyata fungsi network namespace.

---

# 14. Image dan Container Bukan Hal yang Sama

Dalam eksperimen:

```text
alpine
```

adalah **image**.

Sedangkan:

```text
probe
```

adalah **container**.

Image dapat dianggap sebagai template read-only yang digunakan sebagai dasar membuat container.

Mental model:

```text
Alpine image
     │
     ├── Container A
     ├── Container B
     └── Container probe
```

Ketika menjalankan:

```bash
docker rm -f probe
```

yang dihapus adalah container:

```text
probe
```

bukan image:

```text
alpine
```

Karena itu image Alpine masih dapat digunakan untuk membuat container baru.

---

# 15. Menghapus Container Percobaan

Setelah eksperimen selesai:

```bash
docker rm -f probe
```

Output:

```text
probe
```

Opsi:

```text
-f
```

berarti Docker akan memaksa container berhenti apabila masih berjalan, kemudian menghapus container tersebut.

Setelah itu:

```bash
docker inspect probe
```

tidak akan menemukan object tersebut lagi.

Hal ini juga menjelaskan mengapa PID yang sebelumnya diperoleh tidak seharusnya dianggap sebagai identitas permanen container.

PID berlaku untuk process yang sedang hidup pada saat itu.

Jika container dihentikan, dihapus, dan dibuat kembali, process baru dapat memperoleh PID yang berbeda.

---

# 16. Peran Linux Kernel

Untuk memahami Docker, penting memahami istilah **kernel**.

Kernel adalah bagian inti operating system yang mengelola hal seperti:

* process;
* CPU;
* memory;
* filesystem;
* networking;
* device;
* permission.

Gambaran sederhananya:

```text
Application
    │
    ↓
Linux Kernel
    │
    ├── CPU
    ├── RAM
    ├── Disk
    └── Network
```

Docker menggunakan kemampuan yang sudah disediakan Linux kernel.

Contohnya:

```text
Linux Kernel
│
├── namespaces
├── cgroups
├── process management
├── networking
└── filesystem mechanisms
```

Docker kemudian menggabungkan kemampuan-kemampuan tersebut menjadi pengalaman yang disebut:

```text
container
```

---

# 17. Container Tidak Membawa Kernel Linux Sendiri

Salah satu kesalahpahaman yang umum adalah menganggap setiap Docker container merupakan mini virtual machine.

Container tidak bekerja seperti itu.

Dalam container Linux:

```text
Container A ─┐
Container B ─┼──→ Linux kernel yang sama
Container C ─┘
```

Container menjalankan process yang menggunakan kernel host.

Container mendapatkan lingkungan yang tampak terpisah karena mekanisme seperti:

* namespace;
* cgroup;
* filesystem layers;
* virtual networking.

Karena itu sebuah container dapat terlihat seperti environment Linux tersendiri tanpa harus menjalankan kernel Linux baru di dalam setiap container.

---

# 18. Bagaimana Posisi WSL dalam Eksperimen Ini?

Eksperimen dilakukan melalui environment WSL.

Secara konseptual, lapisannya dapat dipahami sebagai:

```text
Windows
│
└── WSL2
    │
    └── Linux environment
        │
        ├── Linux kernel
        │
        └── Docker
            │
            └── Container
                │
                └── application process
```

Jadi istilah berikut tidak sama:

### Windows

Operating system utama pada mesin.

### WSL

Fitur Windows yang menyediakan environment Linux.

### Linux Kernel

Kernel yang menyediakan process, memory, networking, namespace, cgroup, dan fungsi sistem lainnya.

### Ubuntu

Salah satu distribusi Linux yang dapat digunakan sebagai user-space environment pada WSL.

### Docker

Software dan runtime yang mengelola container.

### Container

Environment terisolasi untuk menjalankan satu atau beberapa process.

### Process

Program yang sedang dijalankan oleh kernel.

---

# 19. Perbedaan Container dan Virtual Machine

Perbedaan utama container dan VM terletak pada tingkat virtualisasi dan penggunaan kernel.

## Container

Container menggunakan kernel host yang sama.

```text
Host Linux
│
└── Linux Kernel
    │
    ├── Container A
    │   └── application
    │
    ├── Container B
    │   └── application
    │
    └── Container C
        └── application
```

Container memperoleh isolasi melalui fitur kernel seperti namespace dan cgroup.

Container tidak memerlukan guest kernel sendiri.

---

## Virtual Machine

Virtual machine menjalankan guest operating system dan guest kernel sendiri.

```text
Physical Machine
│
├── Host / Hypervisor
│
├── VM A
│   ├── Guest OS
│   ├── Guest Kernel
│   └── Application
│
└── VM B
    ├── Guest OS
    ├── Guest Kernel
    └── Application
```

VM memvirtualisasikan sebuah mesin pada tingkat yang lebih rendah dibanding container.

Karena setiap VM biasanya menjalankan kernel dan operating system environment sendiri, kebutuhan resource-nya umumnya lebih besar dibanding container.

---

# 20. Ringkasan Container vs Virtual Machine

| Aspek             | Container                           | Virtual Machine                 |
| ----------------- | ----------------------------------- | ------------------------------- |
| Unit utama        | Process yang diisolasi              | Mesin virtual                   |
| Kernel            | Berbagi kernel host                 | Memiliki guest kernel sendiri   |
| Guest OS lengkap  | Tidak diperlukan                    | Umumnya ada                     |
| Isolasi           | Fitur kernel seperti namespace      | Virtualisasi machine/hardware   |
| Resource control  | Cgroup dan mekanisme kernel         | Resource virtual machine        |
| Startup           | Umumnya cepat                       | Umumnya lebih lambat            |
| Resource overhead | Umumnya lebih kecil                 | Umumnya lebih besar             |
| Filesystem        | Filesystem view/container layers    | Disk/filesystem milik VM        |
| Networking        | Virtual network namespace/interface | Virtual network device milik VM |

---

# 21. Bukti yang Diperoleh dari Eksperimen

Eksperimen menghasilkan beberapa bukti nyata.

## 1. Container mempunyai process nyata di host

Command:

```bash
ps -p $PID -o pid,ppid,comm
```

menghasilkan:

```text
PID   PPID  COMMAND
11945 11921 sleep
```

Berarti process container dapat dilihat dari Linux host.

---

## 2. Container menggunakan namespace yang berbeda

Namespace process container antara lain:

```text
net -> net:[4026532235]
pid -> pid:[4026532233]
mnt -> mnt:[4026532221]
uts -> uts:[4026532222]
```

Sedangkan PID 1 mempunyai:

```text
net -> net:[4026531840]
pid -> pid:[4026532219]
mnt -> mnt:[4026532217]
uts -> uts:[4026532218]
```

Identifier yang berbeda menunjukkan penggunaan namespace yang berbeda.

---

## 3. Memory dapat dibatasi melalui cgroup

Command:

```bash
docker run --rm --memory=64m alpine cat /sys/fs/cgroup/memory.max
```

menghasilkan:

```text
67108864
```

yang sama dengan:

```text
64 MiB
```

---

## 4. CPU dapat dibatasi melalui cgroup

Command:

```bash
docker run --rm --cpus=0.5 alpine cat /sys/fs/cgroup/cpu.max
```

menghasilkan:

```text
50000 100000
```

Dengan rasio:

```text
50000 / 100000 = 0.5
```

sesuai batas:

```text
--cpus=0.5
```

---

## 5. Container mempunyai filesystem view sendiri

Command:

```bash
docker run --rm alpine ls /
```

menampilkan root filesystem Alpine:

```text
bin
dev
etc
home
lib
media
mnt
opt
proc
root
run
sbin
srv
sys
tmp
usr
var
```

---

## 6. Container mempunyai network view sendiri

Command:

```bash
docker run --rm alpine ip addr
```

menunjukkan interface:

```text
lo
eth0
```

dan alamat container:

```text
172.17.0.3/16
```

Ini merupakan bukti praktis bahwa container mendapatkan lingkungan networking tersendiri.

---

# 22. Mental Model Akhir Docker Container

Seluruh eksperimen dapat dirangkum menjadi:

```text
Windows
│
└── WSL / Linux environment
    │
    └── Linux Kernel
        │
        ├── Namespace
        │   ├── PID isolation
        │   ├── network isolation
        │   ├── mount isolation
        │   └── hostname/IPC/etc.
        │
        ├── Cgroup
        │   ├── memory limit
        │   └── CPU limit
        │
        └── Docker
            │
            └── Container
                │
                └── Linux process
                    └── sleep
```

Hal terpenting adalah:

> Docker container bukan sebuah komputer kecil yang mempunyai kernel sendiri.

Container pada dasarnya menjalankan process Linux dengan environment yang diisolasi dan dikontrol menggunakan fasilitas Linux kernel.

---

# 23. Jawaban: Apa Bedanya Container dan VM?

Jika ditanya:

> **Apa bedanya container dan VM?**

Jawaban yang dapat diberikan adalah:

Container menjalankan process yang terisolasi tetapi tetap berbagi kernel host. Docker menggunakan fitur Linux seperti namespace untuk mengisolasi pandangan process terhadap process lain, filesystem, networking, dan hostname, serta cgroup untuk membatasi resource seperti memory dan CPU.

Virtual machine berbeda karena menjalankan guest operating system dan guest kernel sendiri di atas lapisan virtualisasi.

Karena container tidak membutuhkan guest kernel sendiri untuk setiap instance, container umumnya mempunyai overhead lebih kecil dan dapat dibuat serta dijalankan lebih cepat dibanding virtual machine.

Versi paling singkat:

```text
Container
= process terisolasi yang berbagi kernel host

Virtual Machine
= mesin virtual yang menjalankan guest OS dan kernel sendiri
```

---

# 24. Kesimpulan

Eksperimen ini menunjukkan bahwa Docker container dibangun di atas mekanisme Linux, bukan dengan membuat mesin virtual baru untuk setiap container.

Container `probe` menjalankan process nyata:

```text
PID     11945
PPID    11921
COMMAND sleep
```

Process tersebut mempunyai namespace yang berbeda dari process PID 1 untuk beberapa aspek penting seperti:

* PID;
* network;
* mount;
* IPC;
* UTS;
* cgroup;
* time.

Eksperimen cgroup menunjukkan bahwa Docker dapat menerapkan batas resource:

```text
Memory:
--memory=64m
→ memory.max = 67108864
→ 64 MiB
```

dan:

```text
CPU:
--cpus=0.5
→ cpu.max = 50000 100000
→ 0.5 CPU
```

Eksperimen filesystem menunjukkan bahwa container melihat root filesystem Alpine sendiri, sedangkan eksperimen network menunjukkan bahwa container memiliki interface `lo` dan `eth0` beserta alamat IP sendiri dalam network namespace-nya.

Mental model utama yang perlu disimpan adalah:

```text
Container
│
├── adalah process Linux
├── berbagi kernel host
├── Namespace → mengisolasi apa yang terlihat
├── Cgroup    → membatasi resource yang digunakan
├── Filesystem view sendiri
└── Network view sendiri
```

Sedangkan:

```text
Virtual Machine
│
├── merupakan mesin yang divirtualisasikan
├── mempunyai guest operating system
└── mempunyai guest kernel sendiri
```

Dengan demikian, perbedaan fundamentalnya adalah:

> **Container mengisolasi process sambil berbagi kernel host, sedangkan VM memvirtualisasikan mesin dan menjalankan kernel sendiri.**
