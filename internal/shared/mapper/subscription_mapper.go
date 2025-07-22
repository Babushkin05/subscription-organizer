package mapper

import (
	"time"

	"github.com/Babushkin05/subscription-organizer/internal/domain/model"
	"github.com/Babushkin05/subscription-organizer/internal/shared/dto"
	"github.com/google/uuid"
)

func ToSubscriptionModel(dto dto.CreateSubscriptionRequest) (*model.Subscription, error) {
	startDate, err := time.Parse("01-2006", dto.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if dto.EndDate != "" {
		t, err := time.Parse("01-2006", dto.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &t
	}

	return &model.Subscription{
		ID:          uuid.New(),
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      uuid.MustParse(dto.UserID),
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}
