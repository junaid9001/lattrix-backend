package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/http/dto"
	"github.com/junaid9001/lattrix-backend/internal/services"
	"github.com/junaid9001/lattrix-backend/internal/utils/jwtutil"
)

type AuthHandler struct {
	authSevice  *services.AuthService
	rbacService *services.RbacService
}

func NewAuthHandler(authService *services.AuthService, rbacService *services.RbacService) *AuthHandler {
	return &AuthHandler{
		authSevice:  authService,
		rbacService: rbacService,
	}
}

// signup handler
func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var req dto.RegisterRequet

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "validation failed", "error": err.Error()})
	}

	err := h.authSevice.SignUP(req.Username, req.Email, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "signup success"})
}

// login handler
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid request",
		})
	}
	access, refresh, err := h.authSevice.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    access,
		Path:     "/",
		SameSite: "Lax",
		Secure:   false, //true in prod
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/",
		SameSite: "Lax",
		Secure:   false, //true in prod
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{"success": true})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"success": false, "error": "missing refresh token"})
	}

	claims, err := jwtutil.ValidateRefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"success": false, "error": "invalid refresh token"})
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"success": false, "error": "invalid subject"})
	}

	accessToken, err := h.authSevice.RefreshAccessToken(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   false, // true in prod
	})

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	workspaceID := c.Locals("workspaceID").(uuid.UUID)
	role := c.Locals("role")
	permissions, err := h.rbacService.RbacRepo.UserPermissions(uint(userID), workspaceID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user_id":      userID,
			"workspace_id": workspaceID,
			"role":         role,
			"permissions":  permissions,
		},
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false, // true in prod
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false, // true in prod
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}
