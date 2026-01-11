package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"github.com/junaid9001/lattrix-backend/internal/utils/jwtutil"
	"gorm.io/gorm"
)

var (
	ErrInvalidInvite   = errors.New("invalid invitation")
	ErrInviteExpired   = errors.New("invitation expired")
	ErrInviteProcessed = errors.New("invitation already processed")
)

type InvitationService struct {
	db               *gorm.DB
	inviteRepo       repository.InvitationRepository
	notificationRepo repository.NotificationRepository
	rbacRepo         repository.RBACrepository
}

func NewInvitationService(
	db *gorm.DB,
	inviteRepo repository.InvitationRepository,
	notificationRepo repository.NotificationRepository,
	rbacRepo repository.RBACrepository,
) *InvitationService {
	return &InvitationService{
		db:               db,
		inviteRepo:       inviteRepo,
		notificationRepo: notificationRepo,
		rbacRepo:         rbacRepo,
	}
}

func (s *InvitationService) CreateNewInvite(workspaceID, roleID uuid.UUID, email string, invitedBy uint) error {
	// Generate a simple token (uuid)
	token := uuid.New().String()

	invite := &models.WorkspaceInvitation{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		Email:       email,
		RoleID:      roleID,
		InvitedBy:   invitedBy,
		Token:       token,
		Status:      "pending",
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour), // 7 days expiry
	}

	return s.inviteRepo.Create(invite)
}

func (s *InvitationService) AcceptInvitation(
	userID uint,
	userEmail string,
	token string,
) (string, error) {

	var newAccessToken string

	err := s.db.Transaction(func(tx *gorm.DB) error {

		inviteRepo := s.inviteRepo.WithTx(tx)
		notifRepo := s.notificationRepo.WithTx(tx)
		rbacRepo := s.rbacRepo.WithTx(tx)

		invite, err := inviteRepo.FindByToken(token)
		if err != nil {
			return ErrInvalidInvite
		}

		if invite.Status != "pending" {
			return ErrInviteProcessed
		}

		if time.Now().After(invite.ExpiresAt) {
			return ErrInviteExpired
		}

		if invite.Email != userEmail {
			return ErrInvalidInvite
		}

		if err := rbacRepo.AssignRoleToUser(
			userID,
			invite.RoleID,
			invite.WorkspaceID,
		); err != nil {
			return err
		}

		if err := inviteRepo.UpdateStatus(invite.ID, "accepted"); err != nil {
			return err
		}

		if err := notifRepo.MarkAsReadByReference(
			invite.ID,
			"workspace_invite",
		); err != nil {
			return err
		}
		var role models.Role
		if err := s.db.
			Where("id = ? AND workspace_id = ?", invite.RoleID, invite.WorkspaceID).
			First(&role).Error; err != nil {
			return err
		}

		jwt, err := jwtutil.CreateAccessToken(int(userID), invite.WorkspaceID.String(), role.Name)
		if err != nil {
			return err
		}

		newAccessToken = jwt
		return nil
	})

	return newAccessToken, err
}
