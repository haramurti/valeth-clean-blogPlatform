# ʕ•ᴥ•ʔ Kuma Blog

> *"Sebuah platform blog yang akhirnya rilis juga, setelah melewati fase 'it works on my machine' berkali-kali."*

[![Deployment](https://img.shields.io/badge/Status-Live_on_Koyeb-success?style=for-the-badge&logo=koyeb)](https://outdoor-hedi-casio-57ac17f4.koyeb.app/)

## 🌐 Coba Sekarang (Udah Online, Bro!)

Gak perlu *clone* repo, gak perlu install Docker, gak perlu menuhin RAM laptop lo. Gue udah deploy ini ke awan (cloud).

👉 **Akses di sini:** [**https://outdoor-hedi-casio-57ac17f4.koyeb.app/**](https://outdoor-hedi-casio-57ac17f4.koyeb.app/)

*(Kalau loading pertamanya agak lama, maklum ya, server gratisan lagi bangun tidur. Tunggu 10 detik, refresh, nanti dia ngebut lagi kayak gue pas dikejar deadline).*

---

## 🤔 Apaan sih ini?

Ini adalah **Kuma Blog**. Simpelnya, ini kayak Medium.com tapi versi *indie*. Gue bikin ini karena gue pengen punya tempat nulis yang bersih, kenceng, dan gak banyak iklan obat peninggi badan.

Dibangun pakai **Clean Architecture**, jadi kalau lo liat kodingannya, itu rapi banget. Serapi kamar gue kalau lagi mau ada tamu doang.

## ✨ Fitur yang Bisa Lo Mainin

Pas lo buka link di atas, lo bisa ngapain aja?

1.  **Login via Google:** Gak perlu repot bikin password baru (gue tau lo pasti lupa password lo sendiri). Klik, login, beres.
2.  **Nulis Cerita:** Ada editor teks yang *distraction-free*. Fokus nulis aja, jangan fokus mikirin dia yang gak bales chat.
3.  **Komentar Ala Medium:** Coba buka salah satu postingan, terus klik ikon komentar. *Sidebar*-nya bakal muncul dari kanan. UX mahal nih, Bos.
4.  **Bookmark:** Simpan tulisan yang menarik buat dibaca nanti (wacana).
5.  **Hapus Tulisan:** Kalau lo nulis sesuatu pas lagi galau terus nyesel, tinggal hapus aja. Jejak digital aman.

## 🛠 Dapur Pacu (Tech Stack)

Biar kelihatan pinter dikit, ini teknologi yang gue pake di belakang layar:

* **Bahasa:** [Go (Golang)](https://go.dev/) — Biar performanya ngebut, gak kayak sinyal di gunung.
* **Framework:** [Fiber](https://gofiber.io/) — Ringan, kenceng, *expressive*.
* **Database:** PostgreSQL (via Supabase) — Tempat nyimpen semua curhatan user.
* **Frontend:** HTML/CSS Native (Glassmorphism Style) — Gak pake React biar *loading*-nya instan.
* **Deployment:** Docker Container di **Koyeb**.

## 💻 Cara Jalanin di Laptop Sendiri (Kalau Penasaran)

Kalau lo programmer dan pengen ngotak-ngatik isinya di laptop lo (Localhost), silakan:

1.  **Clone Repo ini:**
    ```bash
    git clone [https://github.com/username-lo/valeth-clean-blogPlatform.git](https://github.com/username-lo/valeth-clean-blogPlatform.git)
    ```
2.  **Setup Environment:**
    Copy `.env.example` jadi `.env` terus isi *credentials*-nya.
3.  **Jalanin Docker:**
    ```bash
    docker build -t kuma-blog .
    docker run -p 8080:8080 kuma-blog
    ```
4.  **Buka:** `http://localhost:8080`

---

*Dibuat dengan ❤️ dan sedikit ☕ di malam minggu.*