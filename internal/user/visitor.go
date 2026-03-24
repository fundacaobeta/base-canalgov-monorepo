package user

import (
	"fmt"
	"strings"

	"github.com/fundacaobeta/base-canalgov-monorepo/internal/user/models"
	"github.com/taion809/haikunator"
	"github.com/volatiletech/null/v9"
)

// CreateVisitor creates a new visitor user.
func (u *Manager) CreateVisitor(user *models.User) error {
	// Normalize email address.
	user.Email = null.NewString(strings.ToLower(user.Email.String), user.Email.Valid)

	if user.FirstName == "" && user.LastName == "" {
		h := haikunator.NewHaikunator()
		user.FirstName = h.Haikunate()
	}

	if err := u.q.InsertVisitor.Get(user, user.Email, user.FirstName, user.LastName, user.CustomAttributes); err != nil {
		u.lo.Error("error inserting contact", "error", err)
		return fmt.Errorf("insert contact: %w", err)
	}
	return nil
}

// GetVisitor retrieves a visitor user by ID
func (u *Manager) GetVisitor(id int) (models.User, error) {
	return u.Get(id, "", []string{models.UserTypeVisitor})
}

// IsOffline returns true if the user's availability status is offline.
func (u *Manager) IsOffline(id int) bool {
	user, err := u.Get(id, "", []string{})
	if err != nil {
		return true
	}
	return user.AvailabilityStatus == models.Offline
}

// MergeVisitorToContact reassigns all conversations from a visitor to a contact user,
// then deletes the visitor record.
func (u *Manager) MergeVisitorToContact(visitorID, contactID int) error {
	_, err := u.db.Exec(`
		UPDATE conversations SET contact_id = $1 WHERE contact_id = $2;
		DELETE FROM users WHERE id = $2 AND type = 'visitor'`,
		contactID, visitorID)
	if err != nil {
		u.lo.Error("error merging visitor to contact", "visitor_id", visitorID, "contact_id", contactID, "error", err)
		return fmt.Errorf("merging visitor to contact: %w", err)
	}
	return nil
}
