# ʕ•ᴥ•ʔ Kuma Blog

> _"Sebuah platform blog yang akhirnya rilis juga, setelah melewati fase 'it works on my machine' berkali-kali."_

[![Deployment](https://img.shields.io/badge/Status-Live_on_Koyeb-success?style=for-the-badge&logo=koyeb)](https://outdoor-hedi-casio-57ac17f4.koyeb.app/)

## 🌐 Coba Sekarang (Udah Online, Bro!)

Gak perlu _clone_ repo, gak perlu install Docker, gak perlu menuhin RAM laptop lo. Gue udah deploy ini ke awan (cloud).

👉 **Akses di sini:** [**https://outdoor-hedi-casio-57ac17f4.koyeb.app/**](https://outdoor-hedi-casio-57ac17f4.koyeb.app/)

_(Kalau loading pertamanya agak lama, maklum ya, server gratisan lagi bangun tidur. Tunggu 10 detik, refresh, nanti dia ngebut lagi kayak gue pas dikejar deadline)._

---

## 🤔 Apaan sih ini?

Ini adalah **Kuma Blog**. Simpelnya, ini kayak Medium.com tapi versi _indie_. Gue bikin ini karena gue pengen punya tempat nulis yang bersih, kenceng, dan gak banyak iklan obat peninggi badan.

Dibangun pakai **Clean Architecture**, jadi kalau lo liat kodingannya, itu rapi banget. Serapi kamar gue kalau lagi mau ada tamu doang.

## ✨ Fitur yang Bisa Lo Mainin

Pas lo buka link di atas, lo bisa ngapain aja?

1.  **Login via Google:** Gak perlu repot bikin password baru (gue tau lo pasti lupa password lo sendiri). Klik, login, beres.
2.  **Nulis Cerita:** Ada editor teks yang _distraction-free_. Fokus nulis aja, jangan fokus mikirin dia yang gak bales chat.
3.  **Komentar Ala Medium:** Coba buka salah satu postingan, terus klik ikon komentar. _Sidebar_-nya bakal muncul dari kanan. UX mahal nih, Bos.
4.  **Bookmark:** Simpan tulisan yang menarik buat dibaca nanti (wacana).
5.  **Hapus Tulisan:** Kalau lo nulis sesuatu pas lagi galau terus nyesel, tinggal hapus aja. Jejak digital aman.

## 🛠 Dapur Pacu (Tech Stack)

Biar kelihatan pinter dikit, ini teknologi yang gue pake di belakang layar:

- **Bahasa:** [Go (Golang)](https://go.dev/) — Biar performanya ngebut, gak kayak sinyal di gunung.
- **Framework:** [Fiber](https://gofiber.io/) — Ringan, kenceng, _expressive_.
- **Database:** PostgreSQL (via Supabase) — Tempat nyimpen semua curhatan user.
- **Frontend:** HTML/CSS Native (Glassmorphism Style) — Gak pake React biar _loading_-nya instan.
- **Deployment:** Docker Container di **Koyeb**.

## 💻 Cara Jalanin di Laptop Sendiri (Kalau Penasaran)

Kalau lo programmer dan pengen ngotak-ngatik isinya di laptop lo (Localhost), silakan:

1.  **Clone Repo ini:**
    ```bash
    git clone [https://github.com/username-lo/valeth-clean-blogPlatform.git](https://github.com/username-lo/valeth-clean-blogPlatform.git)
    ```
2.  **Setup Environment:**
    Copy `.env.example` jadi `.env` terus isi _credentials_-nya.
3.  **Jalanin Docker:**
    ```bash
    docker build -t kuma-blog .
    docker run -p 8080:8080 kuma-blog
    ```
4.  **Buka:** `http://localhost:8080`

---

_Dibuat dengan ❤️ dan sedikit ☕ di malam minggu._

```
valeth-clean-blogPlatform
├─ .DS_Store
├─ .dockerignore
├─ Dockerfile
├─ README.md
├─ cmd
│  └─ main.go
├─ config
│  ├─ database.go
│  └─ google_oauth.go
├─ go.mod
├─ go.sum
├─ internal
│  ├─ .DS_Store
│  ├─ delivery
│  │  ├─ .DS_Store
│  │  └─ http
│  │     ├─ auth_handler.go
│  │     ├─ comment_handler.go
│  │     ├─ post_handler.go
│  │     └─ profile_handler.go
│  ├─ domain
│  │  ├─ comment.go
│  │  ├─ post.go
│  │  └─ user.go
│  ├─ infrastructure
│  │  └─ supabase.go
│  ├─ middleware
│  │  └─ auth.go
│  ├─ repository
│  │  ├─ post_repository.go
│  │  ├─ postgres_comment_repository.go
│  │  └─ user_repository.go
│  ├─ usecase
│  │  ├─ comment_usecase.go
│  │  ├─ usecase.go
│  │  └─ user_usecase.go
│  └─ utils
│     └─ jwt.go
└─ web
   ├─ .DS_Store
   ├─ static
   │  └─ images
   │     ├─ Bookmark-2.svg
   │     ├─ Bookmark-3.svg
   │     ├─ Comment-2.svg
   │     ├─ Removal-858.png
   │     ├─ Screenshot_2026-01-19_at_19.14.37-removebg-preview-2.png
   │     ├─ hold-on-i-got-a-meme-for-this.gif
   │     ├─ ichikuma.png
   │     └─ image_processing20250206-832295-jowoon.gif
   └─ templates
    ├─ 404.html
      ├─ create.html
      ├─ index.html
      ├─ library.html
      ├─ login.html
      ├─ post.html
      ├─ profile.html
      └─ register.html

```
