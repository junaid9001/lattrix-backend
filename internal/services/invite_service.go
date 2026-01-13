package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"github.com/junaid9001/lattrix-backend/internal/utils/jwtutil"
	"gorm.io/datatypes"
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
	userRepo         repository.UserRepository
}

func NewInvitationService(
	db *gorm.DB,
	inviteRepo repository.InvitationRepository,
	notificationRepo repository.NotificationRepository,
	rbacRepo repository.RBACrepository,
	userRepo repository.UserRepository,
) *InvitationService {
	return &InvitationService{
		db:               db,
		inviteRepo:       inviteRepo,
		notificationRepo: notificationRepo,
		rbacRepo:         rbacRepo,
		userRepo:         userRepo,
	}
}

func (s *InvitationService) CreateNewInvite(workspaceID, roleID uuid.UUID, email string, invitedBy uint) error {
	existingUser, _ := s.userRepo.FindByEmail(email)

	token := uuid.New().String()

	return s.db.Transaction(func(tx *gorm.DB) error {
		inviteRepo := s.inviteRepo.WithTx(tx)
		notifRepo := s.notificationRepo.WithTx(tx)

		invite := &models.WorkspaceInvitation{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			Email:       email,
			RoleID:      roleID,
			InvitedBy:   invitedBy,
			Token:       token,
			Status:      "pending",
			ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
		}

		if err := inviteRepo.Create(invite); err != nil {
			return err
		}

		if existingUser != nil && existingUser.ID != 0 {

			dataMap := map[string]string{"token": token}
			dataBytes, _ := json.Marshal(dataMap)

			notif := &models.Notification{
				ID:          uuid.New(),
				UserID:      existingUser.ID,
				Type:        "invitation",
				Title:       "Workspace Invitation",
				Message:     "You have been invited to join a workspace.",
				ReferenceID: invite.ID,
				Data:        datatypes.JSON(dataBytes),
				IsRead:      false,
			}

			if err := notifRepo.Create(notif); err != nil {
				return err
			}
		} else {

		}

		return nil
	})
}

func (s *InvitationService) AcceptInvitation(
	userID uint,
	token string,
) (string, error) {

	var newAccessToken string

	user, err := s.userRepo.FindByID(int(userID))
	if err != nil {
		return "", err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {

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

		if invite.Email != user.Email {
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

		if err := tx.Model(&models.User{}).
			Where("id = ?", userID).
			Updates(map[string]interface{}{
				"workspace_id": invite.WorkspaceID,
				"role":         role.Name,
			}).Error; err != nil {
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
