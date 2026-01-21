package http

import (
	"fmt"
	"valeth-clean-blogPlatform/internal/domain"

	"valeth-clean-blogPlatform/internal/middleware"

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

	// --- RUTE PUBLIC (Bebas Masuk) ---
	api.Get("/posts", handler.Fetch)
	api.Get("/posts/:id", handler.GetByID)
	app.Get("/", handler.PageHome)
	app.Get("/post", handler.PagePostDetail)
	app.Get("/login", handler.PageLogin)
	app.Get("/register", handler.PageRegister)

	// --- RUTE PRIVATE (Dijaga Satpam) ---
	// Kita bikin grup baru khusus yang diproteksi
	protected := api.Group("/", middleware.AuthProtected)

	// Semua rute di bawah ini otomatis dicek login dulu
	protected.Post("/posts", handler.Store)
	protected.Delete("/posts/:id", handler.Delete)

	// Page Create juga harus dijaga satpam
	// Kalau belum login, jangan kasih buka halaman nulis
	app.Get("/create", middleware.AuthProtected, handler.PageCreate)
	app.Get("/profile", middleware.AuthProtected, handler.PageProfile)

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

	userID := c.Locals("user_id").(int) // Assert ke int

	post.UserID = uint(userID)

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
// Buka file: internal/delivery/http/post_handler.go

func (h *PostHandler) PageHome(c *fiber.Ctx) error {
	// 1. Ambil Data Postingan
	posts, err := h.postUseCase.Fetch("")
	if err != nil {
		return c.Status(500).SendString("Error database bro: " + err.Error())
	}

	// 2. Cek Cookie Login & Avatar
	cookieUserID := c.Cookies("user_id")

	// --- [BARIS BARU] AMBIL FOTO DARI SAKU ---
	cookieAvatar := c.Cookies("avatar")
	// -----------------------------------------

	isLoggedIn := cookieUserID != ""

	// 3. Kirim Data Lengkap ke HTML
	return c.Render("index", fiber.Map{
		"Posts":      posts,
		"IsLoggedIn": isLoggedIn,

		// --- [BARIS BARU] KIRIM KE FRONTEND ---
		"UserAvatar": cookieAvatar,
		// --------------------------------------
	})
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
	fmt.Println("📄 [HANDLER] Halaman Create dipanggil!")
	return c.Render("create", nil) // Asumsi lu nanti bikin create.html
}

// ... import dan kode lain ...

func (h *PostHandler) PageProfile(c *fiber.Ctx) error {
	// 1. Cek ID siapa yang mau dilihat?
	// Kalau ada ?id=5 di URL, pakai itu.
	targetID := c.QueryInt("id")

	// 2. Ambil ID kita sendiri dari Cookie (buat cek login & fallback)
	myCookieID := c.Cookies("user_id")
	var myID int
	fmt.Sscanf(myCookieID, "%d", &myID)

	// Kalau URL gak ada ?id=..., berarti user mau liat profil sendiri
	if targetID == 0 {
		if myID != 0 {
			targetID = myID
		} else {
			// Kalau gak ada ID target DAN belum login -> Tendang ke Login
			return c.Redirect("/login")
		}
	}

	// 3. Ambil Data Postingan milik Target ID
	posts, err := h.postUseCase.FetchByUserID(targetID)
	if err != nil {
		return c.Status(500).SendString("Gagal ambil post: " + err.Error())
	}

	// 4. Ambil Data User Target (Dari postingan pertama aja biar hemat query)
	// Trik: Kalau dia punya postingan, data user ada di posts[0].User
	// Kalau dia GAK punya postingan, kita harus query manual (PR nanti),
	// tapi sementara kita handle kalau posts ada isinya aja.

	var profileUser domain.User
	if len(posts) > 0 {
		profileUser = posts[0].User
	} else {
		// TODO: Kalau user belum pernah posting, idealnya kita fetch user by ID lewat UserUseCase.
		// Tapi biar cepet, kita kosongin dulu atau redirect home.
		// return c.Redirect("/")
		// Biar ga error, kita bikin dummy user (atau fetch manual via repo user kalau lu udah siapin)
		profileUser.Name = "New Member"
		profileUser.Avatar = "https://api.dicebear.com/9.x/micah/svg?seed=new"
	}

	// 5. Cek apakah ini profil kita sendiri? (Buat nampilin tombol Edit nanti)
	isOwnProfile := (targetID == myID)

	// 6. Data buat Header (Cookie)
	cookieAvatar := c.Cookies("avatar")

	return c.Render("profile", fiber.Map{
		"Posts":        posts,
		"ProfileUser":  profileUser, // Data pemilik profil
		"IsOwnProfile": isOwnProfile,
		"IsLoggedIn":   myID != 0,
		"UserAvatar":   cookieAvatar, // Data foto kita di pojok kanan atas
	})
}
