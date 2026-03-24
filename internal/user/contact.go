package user

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/fundacaobeta/base-canalgov-monorepo/internal/dbutil"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/envelope"
	"github.com/fundacaobeta/base-canalgov-monorepo/internal/user/models"
	"github.com/volatiletech/null/v9"
)

// CreateContact creates a new contact user.
func (u *Manager) CreateContact(user *models.User) error {
	password, err := u.generatePassword()
	if err != nil {
		u.lo.Error("generating password", "error", err)
		return fmt.Errorf("generating password: %w", err)
	}

	// Normalize email address.
	user.Email = null.NewString(strings.ToLower(user.Email.String), user.Email.Valid)

	if err := u.q.InsertContact.QueryRow(user.Email, user.FirstName, user.LastName, password, user.AvatarURL, user.InboxID, user.SourceChannelID).Scan(&user.ID, &user.ContactChannelID); err != nil {
		u.lo.Error("error inserting contact", "error", err)
		return fmt.Errorf("insert contact: %w", err)
	}
	return nil
}

// UpdateContact updates a contact in the database.
func (u *Manager) UpdateContact(id int, user models.User) error {
	if _, err := u.q.UpdateContact.Exec(id, user.FirstName, user.LastName, user.Email, user.AvatarURL, user.PhoneNumber, user.PhoneNumberCountryCode); err != nil {
		u.lo.Error("error updating user", "error", err)
		return envelope.NewError(envelope.GeneralError, u.i18n.Ts("globals.messages.errorUpdating", "name", "{globals.terms.contact}"), nil)
	}
	return nil
}

// GetContact retrieves a contact by ID.
func (u *Manager) GetContact(id int, email string) (models.User, error) {
	return u.Get(id, email, []string{models.UserTypeContact})
}

// GetAllContacts returns a list of all contacts.
func (u *Manager) GetContacts(page, pageSize int, order, orderBy string, filters string) ([]models.User, error) {
	if pageSize > maxListPageSize {
		return nil, envelope.NewError(envelope.InputError, u.i18n.Ts("globals.messages.pageTooLarge", "max", fmt.Sprintf("%d", maxListPageSize)), nil)
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// If filters contains a segment_id, we need to load the segment filters first.
	if strings.Contains(filters, "segment_id") {
		var f []dbutil.Filter
		if err := json.Unmarshal([]byte(filters), &f); err == nil {
			for _, filter := range f {
				if filter.Field == "segment_id" && filter.Operator == "eq" {
					segmentID, _ := strconv.Atoi(fmt.Sprintf("%v", filter.Value))
					if segmentID > 0 {
						segment, err := u.GetContactSegment(segmentID)
						if err == nil {
							// Merge segment filters with existing filters.
							var segmentFilters []dbutil.Filter
							if err := json.Unmarshal(segment.Filters, &segmentFilters); err == nil {
								// Remove the segment_id filter itself.
								var newFilters []dbutil.Filter
								for _, existingFilter := range f {
									if existingFilter.Field != "segment_id" {
										newFilters = append(newFilters, existingFilter)
									}
								}
								newFilters = append(newFilters, segmentFilters...)
								b, _ := json.Marshal(newFilters)
								filters = string(b)
							}
						}
					}
					break
				}
			}
		}
	}

	query, qArgs, err := dbutil.BuildPaginatedQuery(`
		SELECT COUNT(*) OVER() as total, id, created_at, updated_at, email, first_name, last_name, avatar_url, phone_number, phone_number_country_code, enabled, custom_attributes
		FROM users
		WHERE type = 'contact' AND deleted_at IS NULL`,
		[]interface{}{}, dbutil.PaginationOptions{
			Order:    order,
			OrderBy:  orderBy,
			Page:     page,
			PageSize: pageSize,
		}, filters, dbutil.AllowedFields{
			"users": {"id", "email", "first_name", "last_name", "phone_number", "created_at", "updated_at", "enabled", "custom_attributes"},
		})

	if err != nil {
		u.lo.Error("error building contacts query", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, u.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.user}"), nil)
	}

	var users []models.User
	if err := u.db.Select(&users, query, qArgs...); err != nil {
		u.lo.Error("error fetching contacts from DB", "error", err)
		return nil, envelope.NewError(envelope.GeneralError, u.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.user}"), nil)
	}

	return users, nil
}
