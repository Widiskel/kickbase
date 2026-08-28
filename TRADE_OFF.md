# Trade Off & Assumptions

Dokumen ini mencatat semua asumsi, keputusan desain, dan hal-hal yang **tidak disebutkan secara eksplisit** dalam dokumen ujian teknis (Soal 1).

---

## Asumsi Domain

| # | Asumsi | Alasan |
|---|--------|--------|
| 1 | Logo tim disimpan sebagai string URL, bukan file upload | Tidak ada requirement upload file; infrastruktur CDN/file storage di luar scope |
| 2 | Posisi pemain mengikuti sistem eFootball (lihat tabel di bawah) | Dokumen hanya menyebut 4 posisi umum; eFootball lebih detail dan relevan untuk konteks sepakbola Indonesia |
| 3 | Pemain memiliki atribut `playstyle` yang menentukan gaya bermain (lihat tabel di bawah) | Menambah kedalaman data pemain; tidak ada di dokumen tapi relevan untuk konteks sepakbola |
| 4 | Nomor punggung unik per tim, bukan per pertandingan | Dokumen: "nomor punggung antar pemain dalam satu tim tidak boleh sama" |
| 5 | 1 pertandingan hanya melibatkan 2 tim (home dan away) | Format standar sepakbola amatir |
| 6 | 1 pertandingan hanya bisa memiliki 1 laporan hasil | Tidak ada konsep leg atau replay |
| 7 | Skor akhir berupa angka integer (0, 1, 2, ...) | Standar sepakbola; tidak ada setengah gol |
| 8 | Waktu gol dicatat sebagai string (misal: "23'", "45+2'", "90'") | Fleksibel untuk injury time dan extra time |
| 9 | "Pencetak gol terbanyak" pada report mengacu ke pencetak gol terbanyak di pertandingan tersebut | Jika ada beberapa pemain dengan jumlah gol sama, tampilkan semua |
| 10 | "Akumulasi total kemenangan" dihitung dari semua pertandingan yang sudah ada hasilnya, bukan hanya yang dijadwalkan | Hanya pertandingan dengan laporan hasil yang dihitung |
| 11 | Tim home dan tim away dalam satu pertandingan boleh dari perusahaan yang sama (antar tim binaan) | Perusahaan XYZ menaungi beberapa tim |

## Posisi Pemain (eFootball)

Mengikuti sistem posisi eFootball oleh Konami. Dokumen ujian hanya menyebut 4 posisi umum (penyerang, gelandang, bertahan, penjaga gawang), tetapi implementasi menggunakan posisi spesifik eFootball:

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

## Playstyle Pemain (eFootball)

Playstyle menentukan gaya bermain pemain di lapangan. Setiap playstyle kompatibel dengan posisi tertentu.

### Forwards Playstyles

| Playstyle | Posisi Kompatibel | Keterangan |
|-----------|-------------------|------------|
| `goal_poacher` | CF | Striker yang selalu bergerak di belakang bek terakhir, siap menyambut peluang |
| `dummy_runner` | CF, SS, AMF | Penyerang yang menarik perhatian bek, membuka ruang untuk rekan |
| `fox_in_the_box` | CF | Striker yang aktif di kotak penalti, menunggu umpan untuk diselesaikan |
| `target_man` | CF | Striker bertubuh besar yang menggunakan fisik untuk menahan bola |
| `deep_lying_forward` | CF, SS | Penyerang yang turun ke tengah untuk menerima bola dan membangun serangan |
| `prolific_winger` | LWF, RWF | Winger yang bertahan di sayap untuk memberikan umpan silang |
| `roaming_flank` | LWF, RWF, LMF, RMF | Winger yang sering memotong ke tengah untuk menciptakan peluang |
| `cross_specialist` | LWF, RWF, LMF, RMF | Pemain yang fokus memberikan umpan silang akurat |

### Midfielders Playstyles

| Playstyle | Posisi Kompatibel | Keterangan |
|-----------|-------------------|------------|
| `classic_no_10` | SS, AMF | Playmaker tradisional yang beroperasi di dekat kotak penalti |
| `creative_playmaker` | SS, RWF, LWF, AMF, LMF, RMF | Pemain yang mengeksploitasi celah pertahanan dengan umpan akurat |
| `hole_player` | SS, AMF, CMF, LMF, RMF | Gelandang yang membuat lari menusuk ke kotak penalti saat tim menguasai bola |
| `box_to_box` | CMF, DMF, LMF, RMF | Gelandang pekerja keras yang menutupi seluruh lapangan selama 90 menit |
| `destroyer` | CMF, DMF | Pemain yang agresif menekan dan merebut bola dengan tekel keras |
| `orchestrator` | CMF, DMF | Gelandang yang mengatur tempo permainan dari posisi dalam |
| `anchor_man` | DMF | Gelandang yang duduk di depan bek untuk melindungi lini pertahanan |

### Defenders Playstyles

| Playstyle | Posisi Kompatibel | Keterangan |
|-----------|-------------------|------------|
| `build_up` | CB | Bek yang turun mengambil bola dan membangun serangan dengan umpan panjang |
| `destroyer` | CB, CMF, DMF | Bek yang agresif menekan dan merebut bola |
| `extra_frontman` | CB | Bek yang sering naik membantu serangan |
| `offensive_fullback` | LB, RB | Bek sayap yang aktif naik membantu serangan |
| `defensive_fullback` | LB, RB | Bek sayap yang fokus bertahan dan tetap di posisi |
| `fullback_finisher` | LB, RB | Bek sayap yang masuk ke area sentral untuk mencetak gol |

### Goalkeeper Playstyles

| Playstyle | Posisi Kompatibel | Keterangan |
|-----------|-------------------|------------|
| `attacking_goalkeeper` | GK | Kiper yang sering keluar dari garis untuk mengantisipasi |
| `defensive_goalkeeper` | GK | Kiper yang tetap dekat garis gawang, fokus pada penyelamatan |

## Asumsi Teknis

| # | Asumsi | Alasan |
|---|--------|--------|
| 12 | Autentikasi & Otorisasi menggunakan JWT Bearer Token dengan Role 'admin' | Dokumen menyebut "Admin perusahaan dapat...", mutasi data dilindungi role admin |
| 13 | Endpoint pembacaan (GET teams, players, matches, reports) dapat diakses publik | Aplikasi Android dapat menampilkan info jadwal & report tanpa login wajib |
| 14 | Database menggunakan PostgreSQL | Relational data dengan constraint ketat cocok untuk RDBMS |
| 15 | Soft delete + audit trail via `_history` tables | Setiap entity punya tabel history untuk tracking perubahan dan revert |
| 16 | Tidak ada pagination pada list pertandingan untuk report | Report bersifat agregasi, bukan list biasa |
| 17 | Tidak ada filter/pencarian pada endpoint list | Dokumen tidak menyebutkan fitur filter |
| 18 | Tidak ada upload file (logo tim) | Logo disimpan sebagai URL string |
| 19 | Tidak ada rate limiting | Single-admin internal tool |
| 20 | Tidak ada caching | Data volume rendah (tim amatir) |
| 21 | Tidak ada WebSocket/real-time | Tidak ada requirement notifikasi real-time |

## Keputusan yang Tidak Dijelaskan di Dokumen

| # | Keputusan | Pilihan | Alternatif yang Ditolak |
|---|-----------|---------|------------------------|
| 22 | Format tanggal pertandingan | `YYYY-MM-DD` (ISO 8601) | Format lokal Indonesia (`DD/MM/YYYY`) — ISO lebih universal untuk API |
| 23 | Format waktu pertandingan | `HH:MM:SS` (24 jam) | 12-hour format — 24 jam lebih unambiguous |
| 24 | Response error format | Envelope `{success, error}` | Langsung return error string — envelope lebih konsisten |
| 25 | Cara menanganihapus tim yang punya pemain | Restrict (gagal hapus jika masih ada pemain) | Cascade delete — berisiko kehilangan data historis |
| 26 | Cara menanganihapus pemain yang punya catatan gol | Restrict (gagal hapus jika punya gol) | Cascade delete — skor pertandingan jadi tidak akurat |
| 27 | Status pertandingan sebelum hasil dilaporkan | `scheduled` | Tidak ada status — tapi perlu distinguish jadwal vs selesai |
| 28 | Apakah tim boleh bertanding melawan dirinya sendiri | Tidak boleh (constraint `home ≠ away`) | — |
| 29 | Apakah pemain dari tim A boleh mencetak gol di pertandingan tim B | Tidak boleh (validasi di service layer) | — |
| 30 | Posisi pemain menggunakan enum eFootball, bukan 4 posisi umum | Lebih spesifik dan relevan | 4 posisi umum terlalu abstrak |
| 31 | Playstyle bersifat opsional (nullable) | Tidak semua pemain punya playstyle yang jelas | Wajib — bisa memaksa data yang tidak akurat |

## Out of Scope (Tidak Akan Dibangun)

| Item | Alasan |
|------|--------|
| Admin dashboard / web UI | Dokumen hanya minta backend API |
| Push notification | Tidak ada requirement |
| Multi-tenant (multi perusahaan) | Hanya 1 perusahaan XYZ |
| Liga/kompetisi management | Di luar scope soal |
| Transfer pemain antar tim | Tidak disebutkan |
| Statistik pemain (assist, kartu, dll) | Hanya gol yang disebutkan |
| Export PDF/Excel untuk report | Report via API JSON saja |
| API versioning | Single Android client, tidak perlu backward compat |
| Docker / containerization | Simplifikasi submission |
| CI/CD pipeline | Technical test deliverable |

---

*Dokumen ini wajid diupdate setiap kali ada keputusan desain baru yang tidak eksplisit di dokumen ujian.*
