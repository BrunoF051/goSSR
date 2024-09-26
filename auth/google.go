package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

var (
	GoogleOAuthConfig *oauth2.Config
	oauthStateString  = "random-string" // In production, use a proper random string
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func InitializeOAuthConfig() {
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func (h *Handler) GoogleLoginHandler(c *fiber.Ctx) error {
	url := GoogleOAuthConfig.AuthCodeURL(oauthStateString)
	return c.Redirect(url)
}

func (h *Handler) GoogleCallbackHandler(c *fiber.Ctx) error {
	if c.Query("state") != oauthStateString {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid OAuth state")
	}

	token, err := GoogleOAuthConfig.Exchange(c.Context(), c.Query("code"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Code exchange failed")
	}

	user, err := getUserInfo(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
	}

	fmt.Printf("User: %+v\n", user)
	fmt.Printf("Token: %+v\n", token)

	// Here you can use h.DB to interact with the database
	// For example:
	// dbUser, err := h.createOrUpdateUser(user)
	// if err != nil {
	//     return c.Status(fiber.StatusInternalServerError).SendString("Failed to process user")
	// }

	// Set up user session, etc.

	//
	sess, _ := c.Locals("session").(*session.Session)
	sess.Set("user_id", user.ID)
	sess.Set("user_email", user.Email)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
	}

	// For now, just redirect to home page
	return c.Redirect("/")
}

func getUserInfo(token *oauth2.Token) (*GoogleUser, error) {
	client := GoogleOAuthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %v", err)
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %v", err)
	}

	var user GoogleUser
	if err = json.Unmarshal(contents, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	return &user, nil
}

// Middleware to require authentication
func RequireAuth(c *fiber.Ctx) error {
	sess, ok := c.Locals("session").(*session.Session)
	if !ok {
		return c.Redirect("/auth/google")
	}

	userID := sess.Get("user_id")
	if userID == nil {
		return c.Redirect("/auth/google")
	}

	return c.Next()
}
