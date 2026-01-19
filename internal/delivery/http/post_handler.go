package http

import (
	"valeth-clean-blogPlatform/internal/domain"

	"github.com/gofiber/fiber/v2"
)

// 1. Bikin struct pelayana
type PostHandler struct {
	// jadi ini dimana variabel yang mana fungsinya bisa kita pake buat di method nanti.
	postUseCase domain.PostUseCase
}

// 2. Masukin pelayan
func NewPostHandler(app *fiber.App, u domain.PostUseCase) {
	handler := &PostHandler{
		postUseCase: u,
	}

	// 3. Daftar Menu
	// api ini variable grouping, jadi yang mau dibawah ini harus lewat grouping dahulu
	api := app.Group("/api")
	api.Get("/posts", handler.Fetch)
	api.Get("/posts/:id", handler.GetByID)
	api.Post("/posts", handler.Store)
	api.Delete("/posts/:id", handler.Delete)

	app.Get("/", handler.PageHome) // Open Home
	app.Get("/post", handler.PagePostDetail)
	app.Get("/login", handler.PageLogin)
	app.Get("/register", handler.PageRegister)
	app.Get("/create", handler.PageCreate)

}

func (h *PostHandler) Fetch(c *fiber.Ctx) error {
	// Misal: /posts?search=coding
	keyword := c.Query("search")

	// 2. Si pelaayan teriak Ke managaer sambil bawa nampan namanya fiber, dimana dia bisa menerima pesanan user
	//dan juga bawa nanti makanannya ke customer
	posts, err := h.postUseCase.Fetch(keyword)
	if err != nil {
		// Kalau dapur meledak, bilang ke tamu (Internal Server Error)
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// 3. Sajikan Makanan (200 OK)
	return c.Status(200).JSON(posts)
}

func (h *PostHandler) Store(c *fiber.Ctx) error {
	// 1. Siapin piring kosong
	var post domain.Post

	// 2. Dapet bahan bahan dari user.
	// Tamu ngirim JSON: {"title": "Halo", "content": "..."}
	// Kita tuang ke piring 'post'
	if err := c.BodyParser(&post); err != nil {
		// if tamunya ngirim sampah (JSON rusak), marahin (Bad Request)
		return c.Status(400).JSON(fiber.Map{
			"message": "Data lu ngaco bro: " + err.Error(),
		})
	}

	//3. Kasih ke Manajer ngasih tau
	if err := h.postUseCase.Store(&post); err != nil {
		// Manajer ngecheck kenapa ini kosong (misal judul kosong)
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// 4. Berhasil (201 Created)
	return c.Status(201).JSON(fiber.Map{
		"message": "Mantap, postingan udah tayang!",
	})
}

func (h *PostHandler) GetByID(c *fiber.Ctx) error {
	// 1. Ambil makanan dari nomor resi.
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "ID harus angka woy"})
	}

	// 2. PANGGIL MANAJER
	post, err := h.postUseCase.GetByID(id)
	if err != nil {
		// Asumsi simpel: Kalau error berarti gak ketemu (404)
		// (Nanti bisa didetailin lagi logic error-nya)
		return c.Status(404).Render("404", nil) // Kalo mau niat bikin file 404.html
	}

	return c.Status(200).JSON(post)
}

func (h *PostHandler) Delete(c *fiber.Ctx) error {
	// 1. AMBIL NOMOR MEJA (ID dari URL)
	// Contoh: DELETE /api/posts/5
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID-nya mana woy, harus angka ya!",
		})
	}

	// 2. TERIAK KE MANAJER (Usecase)
	// "Bos, hapus data nomor 5!"
	err = h.postUseCase.Delete(id)
	if err != nil {
		// Bisa jadi error karena ID gak ketemu, atau DB error
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal hapus bro: " + err.Error(),
		})
	}

	// 3. LAPORAN SUKSES
	// Biasanya kalau delete sukses, kita kasih pesan simple aja.
	return c.Status(200).JSON(fiber.Map{
		"message": "Mantap, postingan udah lenyap dari muka bumi.",
	})
}

// ... method JSON yang lama (Fetch, Store, dll) biarin aja di atas ... jadi ktia bakal pake yang diatas buat yang dibawah

// ==========================================
// BAGIAN FRONTEND (RENDER HTML) rendering html
// ==========================================

// 1. HALAMAN DEPAN (HOME)
func (h *PostHandler) PageHome(c *fiber.Ctx) error {
	// Panggil Manajer (Usecase) buat ambil data
	posts, err := h.postUseCase.Fetch("")
	if err != nil {
		return c.Status(500).SendString("Error database bro: " + err.Error())
	}

	// Render file 'index.html' dengan data posts
	return c.Render("index", posts)
}

// 2. HALAMAN BACA POSTINGAN (DETAIL)
func (h *PostHandler) PagePostDetail(c *fiber.Ctx) error {
	// Ambil ?id=1 dari URL
	id := c.QueryInt("id")

	post, err := h.postUseCase.GetByID(id)
	if err != nil {
		return c.Status(404).Render("404", nil) // Kalo mau niat bikin file 404.html
	}

	// Render file 'post.html'
	return c.Render("post", post)
}

// 3. HALAMAN LOGIN (Cuma nampilin doang)
func (h *PostHandler) PageLogin(c *fiber.Ctx) error {
	return c.Render("login", nil)
}

// 4. HALAMAN REGISTER
func (h *PostHandler) PageRegister(c *fiber.Ctx) error {
	return c.Render("register", nil)
}

// 5. HALAMAN CREATE
func (h *PostHandler) PageCreate(c *fiber.Ctx) error {
	return c.Render("create", nil) // Asumsi lu nanti bikin create.html
}
