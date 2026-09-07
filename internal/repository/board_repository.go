package repository

import (
	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/internal/model"
	"gorm.io/gorm"
)

type BoardRepository interface {
	Create(board *model.Board) error
	FindByPublicID(publicID uuid.UUID) (*model.Board, error)
	FindAllByOwner(ownerPublicID uuid.UUID, search string, offset, limit int) ([]model.Board, int64, error)
	FindUserBoards(ownerPublicID uuid.UUID, memberBoardIDs []int64, search string, offset, limit int) ([]model.Board, int64, error)
	FindAll(search string, offset, limit int) ([]model.Board, int64, error)
	Update(board *model.Board) error
	Delete(publicID uuid.UUID) error
}

type boardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{db: db}
}

func (r *boardRepository) Create(board *model.Board) error {
	return r.db.Create(board).Error
}

func (r *boardRepository) FindByPublicID(publicID uuid.UUID) (*model.Board, error) {
	var board model.Board
	err := r.db.Where("public_id = ?", publicID).First(&board).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *boardRepository) FindAllByOwner(ownerPublicID uuid.UUID, search string, offset, limit int) ([]model.Board, int64, error) {
	var boards []model.Board
	var total int64

	query := r.db.Model(&model.Board{}).Where("owner_public_id = ?", ownerPublicID)
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&boards).Error
	if err != nil {
		return nil, 0, err
	}

	return boards, total, nil
}

func (r *boardRepository) FindUserBoards(ownerPublicID uuid.UUID, memberBoardIDs []int64, search string, offset, limit int) ([]model.Board, int64, error) {
	var boards []model.Board
	var total int64

	query := r.db.Model(&model.Board{})
	if len(memberBoardIDs) > 0 {
		query = query.Where("owner_public_id = ? OR internal_id IN ?", ownerPublicID, memberBoardIDs)
	} else {
		query = query.Where("owner_public_id = ?", ownerPublicID)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&boards).Error
	if err != nil {
		return nil, 0, err
	}

	return boards, total, nil
}

func (r *boardRepository) FindAll(search string, offset, limit int) ([]model.Board, int64, error) {
	var boards []model.Board
	var total int64

	query := r.db.Model(&model.Board{})
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&boards).Error
	if err != nil {
		return nil, 0, err
	}

	return boards, total, nil
}

func (r *boardRepository) Update(board *model.Board) error {
	return r.db.Save(board).Error
}

func (r *boardRepository) Delete(publicID uuid.UUID) error {
	return r.db.Where("public_id = ?", publicID).Delete(&model.Board{}).Error
}
