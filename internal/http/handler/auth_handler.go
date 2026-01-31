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
	authSevice     *services.AuthService
	rbacService    *services.RbacService
	profileService *services.ProfileService
}

func NewAuthHandler(authService *services.AuthService, rbacService *services.RbacService, profileService *services.ProfileService) *AuthHandler {
	return &AuthHandler{
		authSevice:     authService,
		rbacService:    rbacService,
		profileService: profileService,
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
	refresh, workspaces, err := h.authSevice.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/",
		SameSite: "Lax",
		Secure:   false, //true in prod
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{"success": true, "workspaces": workspaces})
}

// select a workspace
func (h *AuthHandler) SelectWorkspace(c *fiber.Ctx) error {
	var req dto.SelectWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "invalid request"})
	}

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "missing refresh token"})
	}
	claims, err := jwtutil.ValidateRefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "invalid refresh token"})
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "invalid token subject"})
	}

	accessToken, err := h.authSevice.GenerateAccessTokenForWorkspace(userID, req.WorkspaceID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false,
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
	user, err := h.profileService.GetUserProfile(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "user not found"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user_id":        userID,
			"workspace_id":   workspaceID,
			"role":           role,
			"permissions":    permissions,
			"plan":           user.Plan,
			"is_super_admin": user.IsSuperAdmin,
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

//create workspace handler

func (h *AuthHandler) CreateWorkspace(c *fiber.Ctx) error {

	var req struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return err
	}
	userIdVal := c.Locals("userID")
	if userIdVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	userId := userIdVal.(int)

	wsID, err := h.authSevice.CreateWorkspace(uint(userId), req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"workspace_id": wsID},
	})
}

//list user workspace handler

func (h *AuthHandler) GetUserWorkspaces(c *fiber.Ctx) error {
	userIDVal := c.Locals("userID")
	if userIDVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	userID := userIDVal.(int)

	workspaces, err := h.authSevice.UserWorkspaces(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    workspaces,
	})
}
