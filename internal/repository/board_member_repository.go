package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/internal/dto"
	"github.com/irpanzy/Task-Forge/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BoardMemberRepository interface {
	AddMember(boardID, userID int64) error
	AddMembers(boardID int64, userIDs []int64) error
	RemoveMember(boardID, userID int64) error
	IsMember(boardID, userID int64) (bool, error)
	GetMembers(boardID int64) ([]dto.MemberResponse, error)
	FindMemberBoardIDs(userID int64) ([]int64, error)
}

type boardMemberRepository struct {
	db *gorm.DB
}

func NewBoardMemberRepository(db *gorm.DB) BoardMemberRepository {
	return &boardMemberRepository{db: db}
}

func (r *boardMemberRepository) AddMember(boardID, userID int64) error {
	member := model.BoardMember{
		BoardID:  boardID,
		UserID:   userID,
		JoinedAt: time.Now(),
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&member).Error
}

func (r *boardMemberRepository) AddMembers(boardID int64, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}

	var members []model.BoardMember
	now := time.Now()
	for _, uid := range userIDs {
		members = append(members, model.BoardMember{
			BoardID:  boardID,
			UserID:   uid,
			JoinedAt: now,
		})
	}

	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&members).Error
}

func (r *boardMemberRepository) RemoveMember(boardID, userID int64) error {
	return r.db.Where("board_id = ? AND user_id = ?", boardID, userID).Delete(&model.BoardMember{}).Error
}

func (r *boardMemberRepository) IsMember(boardID, userID int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.BoardMember{}).
		Where("board_id = ? AND user_id = ?", boardID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

type memberScanResult struct {
	PublicID uuid.UUID
	Name     string
	Email    string
	Role     string
	JoinedAt time.Time
}

func (r *boardMemberRepository) GetMembers(boardID int64) ([]dto.MemberResponse, error) {
	var scanResults []memberScanResult

	err := r.db.Table("board_members").
		Select("users.public_id, users.name, users.email, users.role, board_members.joined_at").
		Joins("JOIN users ON users.internal_id = board_members.user_id").
		Where("board_members.board_id = ? AND users.deleted_at IS NULL", boardID).
		Order("board_members.joined_at ASC").
		Scan(&scanResults).Error

	if err != nil {
		return nil, err
	}

	var members []dto.MemberResponse
	for _, res := range scanResults {
		members = append(members, dto.MemberResponse{
			PublicID: res.PublicID,
			Name:     res.Name,
			Email:    res.Email,
			Role:     res.Role,
			JoinedAt: res.JoinedAt,
		})
	}

	return members, nil
}

func (r *boardMemberRepository) FindMemberBoardIDs(userID int64) ([]int64, error) {
	var boardIDs []int64
	err := r.db.Model(&model.BoardMember{}).
		Where("user_id = ?", userID).
		Pluck("board_id", &boardIDs).Error
	if err != nil {
		return nil, err
	}
	return boardIDs, nil
}
