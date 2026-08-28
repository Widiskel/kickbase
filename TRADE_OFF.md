# Trade Off & Assumptions

Dokumen ini mencatat semua asumsi, keputusan desain, dan hal-hal yang **tidak disebutkan secara eksplisit** dalam dokumen ujian teknis (Soal 1).

---

## Asumsi Domain

| # | Asumsi | Alasan |
|---|--------|--------|
| 1 | Logo tim disimpan sebagai string URL, bukan file upload | Tidak ada requirement upload file; infrastruktur CDN/file storage di luar scope |
| 2 | Posisi pemain mengikuti sistem eFootball 15 Posisi (lihat tabel di bawah) | Dokumen hanya menyebut 4 posisi umum; eFootball lebih detail dan relevan untuk konteks sepakbola Indonesia |
| 3 | Pemain memiliki atribut `playstyle` yang menentukan gaya bermain eFootball | Menambah kedalaman data pemain; tidak ada di dokumen tapi relevan untuk konteks sepakbola |
| 4 | Nomor punggung unik per tim, bukan per pertandingan | Dokumen: "nomor punggung antar pemain dalam satu tim tidak boleh sama" (rentang 1-99) |
| 5 | 1 pertandingan hanya melibatkan 2 tim (home dan away) | Format standar sepakbola amatir |
| 6 | 1 pertandingan hanya bisa memiliki 1 laporan hasil | Tidak ada konsep leg atau replay |
| 7 | Skor akhir berupa angka integer (0, 1, 2, ...) | Standar sepakbola; tidak ada setengah gol |
| 8 | Waktu gol dicatat sebagai string (misal: "23'", "45+2'", "90'") | Fleksibel untuk injury time dan extra time |
| 9 | "Pencetak gol terbanyak" pada report mengacu ke pencetak gol terbanyak di pertandingan tersebut | Jika ada beberapa pemain dengan jumlah gol tertinggi yang sama, tampilkan semua |
| 10 | "Akumulasi total kemenangan" dihitung dari seluruh pertandingan selesai yang melibatkan tim tersebut | Agregasi historis head-to-head yang akurat |
| 11 | Tim home dan tim away dalam satu pertandingan boleh dari perusahaan yang sama (antar tim binaan) | Perusahaan XYZ menaungi beberapa tim |

---

## Posisi Pemain (eFootball)

Mengikuti sistem posisi eFootball oleh Konami. Dokumen ujian hanya menyebut 4 posisi umum (penyerang, gelandang, bertahan, penjaga gawang), tetapi implementasi menggunakan 15 posisi spesifik eFootball:

### Forwards (Penyerang)

| Posisi | Nama | Keterangan |
|--------|------|------------|
| `CF` | Centre Forward | Striker utama, pencetak gol |
| `SS` | Second Striker | Penyerang kedua, bermain di belakang CF |
| `LWF` | Left Wing Forward | Penyerang sayap kiri |
| `RWF` | Right Wing Forward | Penyerang sayap kanan |

### Midfielders (Gelandang)

| Posisi | Nama | Keterangan |
|--------|------|------------|
| `AMF` | Attacking Midfielder | Gelandang serang |
| `CMF` | Central Midfielder | Gelandang tengah |
| `DMF` | Defensive Midfielder | Gelandang bertahan |
| `LMF` | Left Midfielder | Gelandang sayap kiri |
| `RMF` | Right Midfielder | Gelandang sayap kanan |

### Defenders (Bertahan)

| Posisi | Nama | Keterangan |
|--------|------|------------|
| `CB` | Centre Back | Bek tengah |
| `LB` | Left Back | Bek kiri |
| `RB` | Right Back | Bek kanan |

### Goalkeeper (Penjaga Gawang)

| Posisi | Nama | Keterangan |
|--------|------|------------|
| `GK` | Goalkeeper | Penjaga gawang |

---

## Playstyle Pemain (eFootball)

Playstyle menentukan gaya bermain pemain di lapangan. Setiap playstyle kompatibel dengan posisi tertentu:

- **Forwards**: `Goal Poacher`, `Dummy Runner`, `Fox in the Box`, `Target Man`, `Deep-Lying Forward`, `Prolific Winger`, `Roaming Flank`, `Cross Specialist`.
- **Midfielders**: `Classic No. 10`, `Creative Playmaker`, `Hole Player`, `Box-to-Box`, `The Destroyer`, `The Orchestrator`, `Anchor Man`.
- **Defenders**: `Build Up`, `The Destroyer`, `Extra Frontman`, `Offensive Fullback`, `Defensive Fullback`, `Fullback Finisher`.
- **Goalkeepers**: `Offensive Goalkeeper`, `Defensive Goalkeeper`.

---

## Asumsi Teknis & Keamanan

| # | Asumsi | Alasan |
|---|--------|--------|
| 12 | Autentikasi menggunakan JWT Dual-Token (Access Token 1 Jam + Refresh Token Rotation 7 Hari di DB) | Memenuhi standar keamanan API modern dan requirement security assessment |
| 13 | Otorisasi berbasis Granular RBAC Permissions (`<domain>:<action>` dengan 24 permissions) | Memungkinkan pemisahan hak akses fleksibel antar role (`admin`, `staff`, `viewer`) |
| 14 | Endpoint pembacaan (GET teams, players, matches, reports) dapat diakses publik | Memudahkan integrasi dengan aplikasi mobile Android tanpa mewajibkan login untuk membaca data |
| 15 | Seluruh endpoint list mendukung Pagination (`page`, `limit`), Dynamic Filtering, dan Sorting (`sort_by`, `order`) | Mencegah bottleneck performa dan memudahkan query data spesifik |
| 16 | Database menggunakan PostgreSQL 16 dengan GORM | Relational data dengan constraint integritas ketat (FK, Unique Index, Cascading) cocok untuk RDBMS |
| 17 | Soft delete + Audit Trail via `_history` tables bertipe JSONB indexed | Setiap mutasi tercatat lengkap dengan kemampuan rollback / revert snapshot versi |
| 18 | Observabilitas terstandarisasi dengan Zerolog JSON logging, Prometheus `/metrics`, dan Grafana | Memudahkan monitoring latensi, request rate, dan metrik bisnis |
| 19 | Containerization multi-container via Docker/Podman Compose | Memastikan environment seragam bagi evaluator dalam satu perintah `docker compose up` |

---

## Keputusan Arsitektur & Aturan Bisnis

| # | Keputusan | Pilihan | Alternatif yang Ditolak |
|---|-----------|---------|------------------------|
| 20 | Format tanggal & waktu pertandingan | `YYYY-MM-DD` & `HH:MM:SS` | Format lokal (`DD/MM/YYYY`) — ISO 8601 lebih universal untuk API |
| 21 | Response envelope format | Envelope konsisten `{success, data, message, error}` | Bare response — envelope lebih terstruktur untuk integrasi client |
| 22 | Penghapusan tim yang memiliki pemain | Restrict (409 Conflict) | Cascade delete — berisiko menghapus data historis pemain tanpa sengaja |
| 23 | Penghapusan pemain yang memiliki catatan gol | Restrict (409 Conflict) | Cascade delete — skor pertandingan masa lalu menjadi tidak akurat |
| 24 | Transisi status pertandingan | `scheduled` ➔ `completed` / `cancelled` / `deferred` | Langsung delete jadwal — status lifecycle menjaga rekam jejak |
| 25 | Larangan tim bertanding melawan dirinya sendiri | Constraint `home_team_id != away_team_id` (400 Bad Request) | Membiarkan di DB — tidak masuk akal dalam aturan sepakbola |
| 26 | Validasi pencetak gol | Pemain wajib terdaftar di salah satu tim yang bertanding | Bebas input ID pemain — mencegah data anomali |

---

## Out of Scope (Tidak Masuk Scope Ujian)

| Item | Alasan |
|------|--------|
| Custom Web Frontend / Mobile App | Dokumen ujian fokus pada evaluasi Backend REST API JSON |
| Push notification ke device | Tidak ada requirement notifikasi real-time |
| Multi-tenant (multi perusahaan) | Hanya melayani entitas tim amatir Perusahaan XYZ |
| Sistem transfer bursa pemain | Di luar cakupan Soal 1 |
