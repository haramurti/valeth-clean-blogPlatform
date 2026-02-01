package http

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
	"valeth-clean-blogPlatform/internal/domain"
	"valeth-clean-blogPlatform/internal/infrastructure"
	"valeth-clean-blogPlatform/internal/middleware"
	"valeth-clean-blogPlatform/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// 1. Bikin struct pelayan
type PostHandler struct {
	postUseCase domain.PostUseCase
	userUseCase domain.UserUseCase
	storage     *infrastructure.SupabaseStorage
}

// 2. Masukin pelayan
func NewPostHandler(app *fiber.App, u domain.PostUseCase, userU domain.UserUseCase, s *infrastructure.SupabaseStorage) {
	handler := &PostHandler{
		postUseCase: u,
		userUseCase: userU,
		storage:     s,
	}

	// 3. Daftar Menu
	api := app.Group("/api")

	// --- RUTE PUBLIC ---
	api.Get("/posts", handler.Fetch)
	api.Get("/posts/:id", handler.GetByID)

	// Page View Public
	app.Get("/", handler.PageHome)
	app.Get("/post", handler.PagePostDetail)
	app.Get("/login", handler.PageLogin)
	app.Get("/register", handler.PageRegister)

	// --- RUTE PRIVATE (Dijaga Satpam) ---
	protected := api.Group("/", middleware.AuthProtected)

	protected.Post("/posts", handler.Store)
	protected.Delete("/posts/:id", handler.Delete)
	protected.Post("/bookmarks/:id", handler.ToggleBookmark)

	// Page View Private
	app.Get("/library", middleware.AuthProtected, handler.PageLibrary)
	app.Get("/create", middleware.AuthProtected, handler.PageCreate)
	app.Get("/profile", middleware.AuthProtected, handler.PageProfile)
}

// ==========================================
// API JSON HANDLERS
// ==========================================

func (h *PostHandler) Fetch(c *fiber.Ctx) error {
	keyword := c.Query("search")
	posts, err := h.postUseCase.Fetch(keyword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(200).JSON(posts)
}

func (h *PostHandler) Store(c *fiber.Ctx) error {
	var post domain.Post

	if err := c.BodyParser(&post); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Data error: " + err.Error()})
	}

	// Logic Tags
	tagsRaw := c.FormValue("tags")
	if tagsRaw != "" {
		splitted := strings.Split(tagsRaw, ",")
		var cleanTags []string
		for _, t := range splitted {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				cleanTags = append(cleanTags, trimmed)
			}
		}
		post.Tags = cleanTags
	}

	// Logic Upload Gambar
	fileHeader, err := c.FormFile("image")
	if err == nil {
		file, _ := fileHeader.Open()
		defer file.Close()

		fileBytes, _ := io.ReadAll(file)
		userID := c.Locals("user_id").(int)
		ext := filepath.Ext(fileHeader.Filename)
		filename := fmt.Sprintf("post-%d-%d%s", userID, time.Now().Unix(), ext)

		imageURL, errUpload := h.storage.UploadFile(fileBytes, filename, fileHeader.Header.Get("Content-Type"), "posts")
		if errUpload != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Gagal upload gambar: " + errUpload.Error()})
		}
		post.Image = imageURL
	}

	userID := c.Locals("user_id").(int)
	post.UserID = uint(userID)

	if err := h.postUseCase.Store(&post); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Postingan berhasil dibuat!"})
}

func (h *PostHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "ID harus angka"})
	}

	post, err := h.postUseCase.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Post tidak ditemukan"})
	}

	return c.Status(200).JSON(post)
}

func (h *PostHandler) Delete(c *fiber.Ctx) error {
	postID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "ID invalid"})
	}

	requesterID := c.Locals("user_id")
	if requesterID == nil {
		return c.Status(401).JSON(fiber.Map{"message": "Login dulu"})
	}
	currentUserID := requesterID.(int)

	post, err := h.postUseCase.GetByID(postID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Post tidak ditemukan"})
	}

	if post.UserID != uint(currentUserID) {
		return c.Status(403).JSON(fiber.Map{"message": "Bukan tulisan lu!"})
	}

	if err := h.postUseCase.Delete(postID); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus: " + err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Post berhasil dihapus"})
}

// ==========================================
// PAGE HANDLERS (HTML)
// ==========================================

func (h *PostHandler) PageHome(c *fiber.Ctx) error {
	searchKeyword := c.Query("search")
	posts, err := h.postUseCase.Fetch(searchKeyword)
	if err != nil {
		return c.Status(500).SendString("Error database: " + err.Error())
	}

	myID := h.getMyUserID(c)
	isLoggedIn := myID != 0

	// LOGIC CEK BOOKMARK
	if isLoggedIn {
		bookmarkedPosts, err := h.userUseCase.GetBookmarks(myID)
		if err == nil {
			bookmarkMap := make(map[uint]bool)
			for _, b := range bookmarkedPosts {
				bookmarkMap[b.ID] = true
			}
			for i := range posts {
				if bookmarkMap[posts[i].ID] {
					posts[i].IsBookmarked = true
				}
			}
		}
	}

	return c.Render("index", fiber.Map{
		"Posts":       posts,
		"IsLoggedIn":  isLoggedIn,
		"UserAvatar":  c.Cookies("avatar"),
		"SearchQuery": searchKeyword,
	})
}

func (h *PostHandler) PagePostDetail(c *fiber.Ctx) error {
	id := c.QueryInt("id")
	post, err := h.postUseCase.GetByID(id)
	if err != nil {
		return c.Status(404).Render("404", nil)
	}

	myID := h.getMyUserID(c)

	// Cek bookmark di halaman detail
	if myID != 0 {
		bookmarkedPosts, _ := h.userUseCase.GetBookmarks(myID)
		for _, b := range bookmarkedPosts {
			if b.ID == post.ID {
				post.IsBookmarked = true
				break
			}
		}
	}

	isOwnPost := (post.UserID == uint(myID))

	return c.Render("post", fiber.Map{
		"Post":       post,
		"IsOwnPost":  isOwnPost,
		"IsLoggedIn": myID != 0,
		"UserAvatar": c.Cookies("avatar"),
	})
}

func (h *PostHandler) PageProfile(c *fiber.Ctx) error {
	targetID := c.QueryInt("id")
	myID := h.getMyUserID(c)

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
		profileUser.Name = "Member"
		profileUser.Avatar = "https://api.dicebear.com/9.x/micah/svg?seed=new"
	}

	isOwnProfile := (targetID == myID)

	return c.Render("profile", fiber.Map{
		"Posts":        posts,
		"ProfileUser":  profileUser,
		"IsOwnProfile": isOwnProfile,
		"IsLoggedIn":   myID != 0,
		"UserAvatar":   c.Cookies("avatar"),
	})
}

func (h *PostHandler) PageLogin(c *fiber.Ctx) error {
	return c.Render("login", nil)
}

func (h *PostHandler) PageRegister(c *fiber.Ctx) error {
	return c.Render("register", nil)
}

func (h *PostHandler) PageCreate(c *fiber.Ctx) error {
	return c.Render("create", nil)
}

// ==========================================
// BOOKMARK HANDLERS
// ==========================================

func (h *PostHandler) ToggleBookmark(c *fiber.Ctx) error {
	userID := h.getMyUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"message": "Login dulu bro!"})
	}

	postID, _ := c.ParamsInt("id")

	err := h.userUseCase.ToggleBookmark(userID, postID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Bookmark updated!"})
}

func (h *PostHandler) PageLibrary(c *fiber.Ctx) error {
	userID := h.getMyUserID(c)
	if userID == 0 {
		return c.Redirect("/login")
	}

	bookmarks, err := h.userUseCase.GetBookmarks(userID)
	if err != nil {
		return c.Status(500).SendString("Error: " + err.Error())
	}

	// Di library, semua post otomatis adalah bookmark (IsBookmarked = true)
	for i := range bookmarks {
		bookmarks[i].IsBookmarked = true
	}

	// Render file template "library"
	return c.Render("library", fiber.Map{
		"Posts":      bookmarks,
		"IsLoggedIn": true,
		"UserAvatar": c.Cookies("avatar"),
	})
}

// ==========================================
// HELPER
// ==========================================

func (h *PostHandler) getMyUserID(c *fiber.Ctx) int {
	tokenString := c.Cookies("jwt_token")
	if tokenString == "" {
		return 0
	}
	userID, err := utils.ParseToken(tokenString)
	if err != nil {
		return 0
	}
	return int(userID)
}
