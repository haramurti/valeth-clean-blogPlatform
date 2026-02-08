# ʕ•ᴥ•ʔ Kuma Blog

> *"Sebuah platform blog yang dibikin bukan karena butuh tempat curhat, tapi karena butuh validasi kalau gue bisa ngoding Go."*

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat&logo=go)
![Fiber](https://img.shields.io/badge/Fiber-v2-black?style=flat)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supabase-316192?style=flat&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)

## 🤔 Apaan nih?

Ini adalah **Kuma Blog**. Simpelnya, ini kayak Medium.com tapi versi *lite*, gratisan, dan yang bikin belum dapet pendanaan seri A.

Gue bikin ini pakai **Clean Architecture**. Kenapa? Biar kalau ada fitur yang error, gue tau harus nyalahin file yang mana tanpa perlu nangis di pojok kamar. UI-nya pakai konsep *Glassmorphism* (kaca-kaca buram gitu), biar kelihatan futuristik dan estetik, meskipun isinya mungkin cuma tulisan "Hello World".

## ✨ Fitur Unggulan (Yang Lumayan Bikin Pusing)

* **Authentication:** Bisa login pakai **Google OAuth** (karena gue tau lo pasti males ngafalin password baru) atau email biasa.
* **Clean Writing Space:** Editor teks simpel. Gak banyak tombol aneh-aneh. Fokus nulis, jangan fokus nyari *font* Comic Sans.
* **Sidebar Comments:** Komentar muncul dari samping ala Medium. UX mahal, modal *JavaScript* native.
* **Bookmark System:** Simpan tulisan yang mau dibaca nanti (wacana doang biasanya).
* **Responsive UI:** Jalan mulus di Laptop, HP, maupun kalkulator (kalau ada browsernya).

## 🛠 Teknologi di Balik Layar

Dibangun dengan keringat, air mata, dan kopi sachet:

* **Backend:** [Go (Golang)](https://go.dev/) + [Fiber](https://gofiber.io/) — Karena hidup udah lambat, web jangan ikutan lambat.
* **Database:** PostgreSQL (via Supabase) — Nyimpen data mantan... eh, data user.
* **ORM:** GORM — Biar gak perlu nulis raw SQL panjang-panjang.
* **Frontend:** Go `html/template` + CSS Native — Gak pake React/Vue. Laki itu ngoding CSS manual.
* **Deployment:** Docker — Biar gak ada alasan *"It works on my machine"* pas demo.

## 🚀 Cara Jalanin (Buat Lo yang Mau Coba)

Pastikan di laptop lo udah ada **Go** dan **Docker**. Kalau belum, install dulu. Jangan males.

1.  **Clone Repo ini:**
    ```bash
    git clone [https://github.com/username-lo/kuma-blog.git](https://github.com/username-lo/kuma-blog.git)
    cd kuma-blog
    ```

2.  **Siapin File Rahasia (.env):**
    Copy file `.env.example` jadi `.env`, terus isi data database & Google Auth lo.
    ```bash
    cp .env.example .env
    ```

3.  **Jalanin via Docker (Jalur VIP):**
    ```bash
    docker-compose up --build
    ```

4.  **Buka Browser:**
    Akses di `http://localhost:8080`. Selamat, webnya jalan. Kalau error, coba restart laptop, biasanya sembuh.

## 🤝 Kontribusi

Kalau nemu *bug*, tolong *issue*-nya ditulis yang sopan. Kalau mau benerin sendiri (*Pull Request*), gue doain masuk surga.

---
*Dibuat oleh [Nama Lo]. Masih jomblo (opsional, hapus kalau udah laku).*