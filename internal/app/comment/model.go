package comment

import (
	"time"

	"github.com/Mickey327/rcsp-backend/internal/app/user"
)

type Comment struct {
	ID        uint64     `db:"id"`
	ProductID uint64     `db:"product_id"`
	Message   string     `db:"message"`
	User      *user.User `scan:"notate"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}

func (c *Comment) ToDTO() *DTO {
	dto := &DTO{
		Message:   c.Message,
		UpdatedAt: c.UpdatedAt.Format("02 Jan 2006 15:04"),
	}
	if c.User != nil {
		dto.UserEmail = c.User.Email
	}
	return dto
}

func ToDTOs(comments []*Comment) []*DTO {
	var commentDTOs []*DTO

	for _, comment := range comments {
		commentDTOs = append(commentDTOs, comment.ToDTO())
	}

	return commentDTOs
}
