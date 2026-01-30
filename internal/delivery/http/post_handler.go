package http

import (
	"fmt"
	"io"
	"path/filepath"
	"time"
	"valeth-clean-blogPlatform/internal/domain"
	"valeth-clean-blogPlatform/internal/infrastructure"
	"valeth-clean-blogPlatform/internal/utils"

	"valeth-clean-blogPlatform/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// 1. Bikin struct pelayana
type PostHandler struct {
	// jadi ini dimana variabel yang mana fungsinya bisa kita pake buat di method nanti.
	postUseCase domain.PostUseCase
	storage     *infrastructure.SupabaseStorage // 👈 1. Tambah ini
}

// 2. Masukin pelayan
func NewPostHandler(app *fiber.App, u domain.PostUseCase, s *infrastructure.SupabaseStorage) {
	handler := &PostHandler{
		postUseCase: u,
		storage:     s,
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

	fileHeader, err := c.FormFile("image")

	// Kalau tidak error (artinya ada file), kita upload
	if err == nil {
		file, _ := fileHeader.Open()
		defer file.Close()

		fileBytes, _ := io.ReadAll(file) // Baca file jadi bytes

		// Bikin nama unik: post-USERID-TIMESTAMP.jpg
		userID := c.Locals("user_id").(int)
		ext := filepath.Ext(fileHeader.Filename)
		filename := fmt.Sprintf("post-%d-%d%s", userID, time.Now().Unix(), ext)

		// Upload ke Supabase
		imageURL, errUpload := h.storage.UploadFile(
			fileBytes,
			filename,
			fileHeader.Header.Get("Content-Type"),
			"posts", // <--- Kirim ke bucket posts
		)
		if errUpload != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Gagal upload gambar: " + errUpload.Error()})
		}

		// Simpan URL ke struct Post
		post.Image = imageURL
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
	// 1. AMBIL ID POST DARI URL (Target yang mau dihapus)
	postID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID-nya mana woy, harus angka ya!",
		})
	}

	// 2. AMBIL ID USER DARI SATPAM/MIDDLEWARE (Siapa yang request?)
	// (Pastikan middleware AuthProtected sudah jalan)
	requesterID := c.Locals("user_id")
	if requesterID == nil {
		return c.Status(401).JSON(fiber.Map{"message": "Login dulu bro!"})
	}
	// Konversi ke int (sesuai tipe data user ID lu)
	currentUserID := requesterID.(int)

	// 3. CEK KEPEMILIKAN (Database Check)
	// Kita panggil GetByID dulu buat liat siapa pemilik aslinya
	post, err := h.postUseCase.GetByID(postID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Postingan gak ketemu, mungkin udah dihapus duluan.",
		})
	}

	// --- 🛡️ BENTENG PERTAHANAN 🛡️ ---
	// Bandingkan: "ID Pemilik Postingan" vs "ID Orang yang Request"
	if post.UserID != uint(currentUserID) {
		// Kalau beda, TENDANG!
		return c.Status(403).JSON(fiber.Map{
			"message": "HEH! JANGAN MALING! Ini bukan tulisan lu, gaboleh dihapus!",
		})
	}
	// -------------------------------

	// 4. EKSEKUSI HAPUS (Kalau lolos pengecekan di atas)
	err = h.postUseCase.Delete(postID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal hapus bro: " + err.Error(),
		})
	}

	// 5. LAPORAN SUKSES
	return c.Status(200).JSON(fiber.Map{
		"message": "Aman, postingan milikmu sudah dihapus.",
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

	searchKeyword := c.Query("search")
	posts, err := h.postUseCase.Fetch(searchKeyword)
	if err != nil {
		return c.Status(500).SendString("Error database bro: " + err.Error())
	}

	//baru jumat 23 jan 2026
	// 2. CEK LOGIN PAKE HELPER BARU
	myID := h.getMyUserID(c) // <--- Panggil fungsi bantuan tadi
	isLoggedIn := myID != 0  // Kalau ID bukan 0, berarti Login

	// 3. Render
	return c.Render("index", fiber.Map{
		"Posts":       posts,
		"IsLoggedIn":  isLoggedIn,
		"UserAvatar":  c.Cookies("avatars"), // Avatar gapapa baca cookie langsung, gak bahaya
		"SearchQuery": searchKeyword,
	})
}

// 2. HALAMAN BACA POSTINGAN (DETAIL)
// internal/delivery/http/post_handler.go

func (h *PostHandler) PagePostDetail(c *fiber.Ctx) error {
	// 1. Ambil Data Postingan
	id := c.QueryInt("id")
	post, err := h.postUseCase.GetByID(id)
	if err != nil {
		return c.Status(404).Render("404", nil)
	}

	//baru 23 jumat 2026
	myID := h.getMyUserID(c)

	// 2. Bandingkan ID Login vs Pemilik Post
	isOwnPost := (post.UserID == uint(myID))

	return c.Render("post", fiber.Map{
		"Post":       post,
		"IsOwnPost":  isOwnPost,
		"IsLoggedIn": myID != 0,
		"UserAvatar": c.Cookies("avatar"),
	})
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

// kode baaru jumat 23 2026
func (h *PostHandler) PageProfile(c *fiber.Ctx) error {
	targetID := c.QueryInt("id")

	// 1. CEK LOGIN PAKE HELPER BARU
	myID := h.getMyUserID(c)

	// Logika: Kalau gak ada targetID di URL, berarti mau liat profil sendiri
	if targetID == 0 {
		if myID != 0 {
			targetID = myID
		} else {
			return c.Redirect("/login")
		}
	}

	posts, err := h.postUseCase.FetchByUserID(targetID)
	if err != nil {
		return c.Status(500).SendString("Error: " + err.Error())
	}

	var profileUser domain.User
	if len(posts) > 0 {
		profileUser = posts[0].User
	} else {
		// (Optional) Fetch user data manual kalo postingan kosong
		profileUser.Name = "Member"
		profileUser.Avatar = "https://api.dicebear.com/9.x/micah/svg?seed=new"
	}

	// 2. Bandingkan ID Login vs Target Profil
	isOwnProfile := (targetID == myID)

	return c.Render("profile", fiber.Map{
		"Posts":        posts,
		"ProfileUser":  profileUser,
		"IsOwnProfile": isOwnProfile,
		"IsLoggedIn":   myID != 0,
		"UserAvatar":   c.Cookies("avatars"),
	})
}

// --- HELPER BUAT POST HANDLER ---
// Tugas: Cek "Siapa sih yang lagi login?" dengan baca JWT
func (h *PostHandler) getMyUserID(c *fiber.Ctx) int {
	tokenString := c.Cookies("jwt_token")
	if tokenString == "" {
		return 0 // Gak ada token = Belum login
	}

	// Buka segel token
	userID, err := utils.ParseToken(tokenString)
	if err != nil {
		return 0 // Token rusak/palsu = Belum login
	}

	return int(userID)
}
