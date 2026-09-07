package service

import (
	"errors"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/irpanzy/Task-Forge/internal/dto"
	"github.com/irpanzy/Task-Forge/internal/model"
	"github.com/irpanzy/Task-Forge/internal/repository"
	"gorm.io/gorm"
)

type BoardService interface {
	CreateBoard(ownerPublicID uuid.UUID, req *dto.CreateBoardRequest) (*dto.BoardResponse, error)
	GetUserBoards(userPublicID uuid.UUID, search string, page, limit int) (*dto.PaginatedBoardsResponse, error)
	GetBoardDetail(boardPublicID, userPublicID uuid.UUID, userRole string) (*dto.BoardResponse, error)
	UpdateBoard(boardPublicID, userPublicID uuid.UUID, userRole string, req *dto.UpdateBoardRequest) (*dto.BoardResponse, error)
	DeleteBoard(boardPublicID, userPublicID uuid.UUID, userRole string) error

	AddMembers(boardPublicID, userPublicID uuid.UUID, userRole string, memberPublicIDs []string) error
	GetMembers(boardPublicID, userPublicID uuid.UUID, userRole string) ([]dto.MemberResponse, error)
	RemoveMember(boardPublicID, userPublicID, targetMemberPublicID uuid.UUID, userRole string) error
}

type boardService struct {
	boardRepo       repository.BoardRepository
	boardMemberRepo repository.BoardMemberRepository
	userRepo        repository.UserRepository
}

func NewBoardService(
	boardRepo repository.BoardRepository,
	boardMemberRepo repository.BoardMemberRepository,
	userRepo repository.UserRepository,
) BoardService {
	return &boardService{
		boardRepo:       boardRepo,
		boardMemberRepo: boardMemberRepo,
		userRepo:        userRepo,
	}
}

func (s *boardService) CreateBoard(ownerPublicID uuid.UUID, req *dto.CreateBoardRequest) (*dto.BoardResponse, error) {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return nil, errors.New("board title is required")
	}

	user, err := s.userRepo.FindByPublicID(ownerPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("owner user not found")
		}
		return nil, err
	}

	newBoard := model.Board{
		OwnerID:       user.InternalID,
		OwnerPublicID: ownerPublicID,
		Title:         req.Title,
		Description:   strings.TrimSpace(req.Description),
		DueDate:       req.DueDate,
	}

	if err := s.boardRepo.Create(&newBoard); err != nil {
		return nil, err
	}

	res := dto.ToBoardResponse(&newBoard)
	res.IsOwner = true
	return &res, nil
}

func (s *boardService) GetUserBoards(userPublicID uuid.UUID, search string, page, limit int) (*dto.PaginatedBoardsResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return nil, err
	}

	memberBoardIDs, err := s.boardMemberRepo.FindMemberBoardIDs(user.InternalID)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * limit

	boards, totalData, err := s.boardRepo.FindUserBoards(userPublicID, memberBoardIDs, search, offset, limit)
	if err != nil {
		return nil, err
	}

	var boardResponses []dto.BoardResponse
	for _, b := range boards {
		res := dto.ToBoardResponse(&b)
		res.IsOwner = (b.OwnerPublicID == userPublicID)
		boardResponses = append(boardResponses, res)
	}

	totalPages := int(math.Ceil(float64(totalData) / float64(limit)))

	return &dto.PaginatedBoardsResponse{
		Boards:      boardResponses,
		TotalData:   totalData,
		CurrentPage: page,
		TotalPages:  totalPages,
		Limit:       limit,
	}, nil
}

func (s *boardService) GetBoardDetail(boardPublicID, userPublicID uuid.UUID, userRole string) (*dto.BoardResponse, error) {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("board not found")
		}
		return nil, err
	}

	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return nil, err
	}

	isOwner := (board.OwnerPublicID == userPublicID)
	isAdmin := strings.EqualFold(userRole, "admin")

	isMember := false
	if !isOwner && !isAdmin {
		isMember, err = s.boardMemberRepo.IsMember(board.InternalID, user.InternalID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, errors.New("access denied: you do not have permission to view this board")
		}
	}

	members, err := s.boardMemberRepo.GetMembers(board.InternalID)
	if err != nil {
		return nil, err
	}

	res := dto.ToBoardResponse(board)
	res.IsOwner = isOwner
	res.Members = members

	return &res, nil
}

func (s *boardService) UpdateBoard(boardPublicID, userPublicID uuid.UUID, userRole string, req *dto.UpdateBoardRequest) (*dto.BoardResponse, error) {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("board not found")
		}
		return nil, err
	}

	// Authorization check: non-admin can only update boards they own
	if !strings.EqualFold(userRole, "admin") && board.OwnerPublicID != userPublicID {
		return nil, errors.New("access denied: you do not have permission to modify this board")
	}

	if req.Title != "" {
		board.Title = strings.TrimSpace(req.Title)
	}
	if req.Description != "" {
		board.Description = strings.TrimSpace(req.Description)
	}
	if req.DueDate != nil {
		board.DueDate = req.DueDate
	}

	if err := s.boardRepo.Update(board); err != nil {
		return nil, err
	}

	res := dto.ToBoardResponse(board)
	res.IsOwner = (board.OwnerPublicID == userPublicID)
	return &res, nil
}

func (s *boardService) DeleteBoard(boardPublicID, userPublicID uuid.UUID, userRole string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("board not found")
		}
		return err
	}

	// Authorization check: non-admin can only delete boards they own
	if !strings.EqualFold(userRole, "admin") && board.OwnerPublicID != userPublicID {
		return errors.New("access denied: you do not have permission to delete this board")
	}

	return s.boardRepo.Delete(boardPublicID)
}

func (s *boardService) AddMembers(boardPublicID, userPublicID uuid.UUID, userRole string, memberPublicIDs []string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("board not found")
		}
		return err
	}

	// Only owner or admin can add members
	if !strings.EqualFold(userRole, "admin") && board.OwnerPublicID != userPublicID {
		return errors.New("access denied: only the board owner can add members")
	}

	var userInternalIDs []int64
	for _, idStr := range memberPublicIDs {
		targetUUID, err := uuid.Parse(strings.TrimSpace(idStr))
		if err != nil {
			continue // skip invalid UUID strings
		}

		// Skip if target is the owner
		if targetUUID == board.OwnerPublicID {
			continue
		}

		targetUser, err := s.userRepo.FindByPublicID(targetUUID)
		if err != nil {
			continue // skip if user doesn't exist
		}

		userInternalIDs = append(userInternalIDs, targetUser.InternalID)
	}

	if len(userInternalIDs) == 0 {
		return errors.New("no valid members to add")
	}

	return s.boardMemberRepo.AddMembers(board.InternalID, userInternalIDs)
}

func (s *boardService) GetMembers(boardPublicID, userPublicID uuid.UUID, userRole string) ([]dto.MemberResponse, error) {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("board not found")
		}
		return nil, err
	}

	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return nil, err
	}

	isOwner := (board.OwnerPublicID == userPublicID)
	isAdmin := strings.EqualFold(userRole, "admin")

	if !isOwner && !isAdmin {
		isMember, err := s.boardMemberRepo.IsMember(board.InternalID, user.InternalID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, errors.New("access denied: you do not have permission to view members of this board")
		}
	}

	return s.boardMemberRepo.GetMembers(board.InternalID)
}

func (s *boardService) RemoveMember(boardPublicID, userPublicID, targetMemberPublicID uuid.UUID, userRole string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("board not found")
		}
		return err
	}

	// Owner cannot be removed as a member
	if targetMemberPublicID == board.OwnerPublicID {
		return errors.New("cannot remove the board owner")
	}

	isOwner := (board.OwnerPublicID == userPublicID)
	isAdmin := strings.EqualFold(userRole, "admin")
	isSelf := (userPublicID == targetMemberPublicID)

	// Only owner, admin, or the member themselves (leaving) can remove a member
	if !isOwner && !isAdmin && !isSelf {
		return errors.New("access denied: you do not have permission to remove this member")
	}

	targetUser, err := s.userRepo.FindByPublicID(targetMemberPublicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("member user not found")
		}
		return err
	}

	return s.boardMemberRepo.RemoveMember(board.InternalID, targetUser.InternalID)
}
