package comment

import (
	"context"
	"errors"
	"log"

	"github.com/Mickey327/rcsp-backend/internal/app/user"
	"github.com/labstack/echo/v4"
)

type Repository interface {
	Create(ctx context.Context, comment *Comment) (uint64, error)
	ReadByProductID(ctx context.Context, productID uint64) ([]*Comment, error)
	ReadByUserAndProductID(ctx context.Context, userID, productID uint64) (*Comment, error)
	Update(ctx context.Context, comment *Comment) (bool, error)
}

type CommentService struct {
	repository Repository
}

func NewService(repository Repository) *CommentService {
	return &CommentService{
		repository: repository,
	}
}

func (s *CommentService) WriteComment(c echo.Context, userID, productID uint64, message string) (uint64, error) {
	comment, err := s.repository.ReadByUserAndProductID(c.Request().Context(), userID, productID)
	var id uint64

	if errors.Is(err, CommentNotFoundErr) && comment == nil {
		log.Println("CREATE COMMENT")
		comment = &Comment{
			ProductID: productID,
			Message:   message,
			User: &user.User{
				ID: userID,
			},
		}
		id, err = s.repository.Create(c.Request().Context(), comment)
		if err != nil {
			return 0, err
		}
	} else {
		log.Println("UPDATE COMMENT")
		comment.Message = message
		comment.User = &user.User{ID: userID}
		comment.ProductID = productID
		log.Println("COMMENT MESSAGE:", message)
		_, err = s.repository.Update(c.Request().Context(), comment)
		id = comment.ID
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (s *CommentService) ReadByProductID(c echo.Context, productID uint64) ([]*DTO, error) {
	comments, err := s.repository.ReadByProductID(c.Request().Context(), productID)

	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return nil, CommentNotFoundErr
	}

	return ToDTOs(comments), nil
}

func (s *CommentService) ReadByUserAndProductID(c echo.Context, userID, productID uint64) (*DTO, error) {
	comment, err := s.repository.ReadByUserAndProductID(c.Request().Context(), userID, productID)

	if err != nil {
		return nil, err
	}

	return comment.ToDTO(), nil
}
