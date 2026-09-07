package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/internal/model"
)

type CreateBoardRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type UpdateBoardRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type MemberResponse struct {
	PublicID uuid.UUID `json:"public_id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type BoardResponse struct {
	PublicID      uuid.UUID        `json:"public_id"`
	OwnerPublicID uuid.UUID        `json:"owner_public_id"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	DueDate       *time.Time       `json:"due_date,omitempty"`
	IsOwner       bool             `json:"is_owner"`
	Members       []MemberResponse `json:"members,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

type PaginatedBoardsResponse struct {
	Boards      []BoardResponse `json:"boards"`
	TotalData   int64           `json:"total_data"`
	CurrentPage int             `json:"current_page"`
	TotalPages  int             `json:"total_pages"`
	Limit       int             `json:"limit"`
}

func ToBoardResponse(b *model.Board) BoardResponse {
	return BoardResponse{
		PublicID:      b.PublicID,
		OwnerPublicID: b.OwnerPublicID,
		Title:         b.Title,
		Description:   b.Description,
		DueDate:       b.DueDate,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
}
