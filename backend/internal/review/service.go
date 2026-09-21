package review

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type ReviewRepository interface {
	Create(ctx context.Context, review *domain.Review) error
	GetAll(ctx context.Context) ([]*domain.Review, error)
	GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Review, error)
	GetByRestaurantAll(ctx context.Context, restID uuid.UUID) ([]*domain.Review, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Review, error)
	GetByOrder(ctx context.Context, orderID uuid.UUID) (*domain.Review, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error)
	ExistsByOrder(ctx context.Context, orderID uuid.UUID) (bool, error)
	GetPhotosByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.ReviewPhoto, error)
	SetReply(ctx context.Context, review *domain.Review) error
	SetVisibility(ctx context.Context, id uuid.UUID, isHidden bool) error
}

type OrderRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	SetReviewed(ctx context.Context, orderID uuid.UUID) error
}

type RestaurantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error)
	IncrementRating(ctx context.Context, restID uuid.UUID, rating int8) error
}

type Service struct {
	tm         TransactionManager
	reviewRepo ReviewRepository
	orderRepo  OrderRepository
	restRepo   RestaurantRepository
}

func NewService(tm TransactionManager, reviewRepo ReviewRepository, orderRepo OrderRepository, restRepo RestaurantRepository) *Service {
	return &Service{
		tm:         tm,
		reviewRepo: reviewRepo,
		orderRepo:  orderRepo,
		restRepo:   restRepo,
	}
}

func (s *Service) Create(ctx context.Context, input domain.CreateReviewInput) (*domain.Review, error) {
	const op = "review.Service.Create"

	order, err := s.orderRepo.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if order.UserID != input.UserID {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	// Check if the order is completed
	if order.Status != domain.OrderStatusCollected {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrOrderNotCompleted)
	}

	// Check if the review already exists
	exists, err := s.reviewRepo.ExistsByOrder(ctx, input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// if the review already exists, return an error
	if exists {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
	}

	input.RestaurantID = order.RestaurantID
	review := domain.NewReview(input)

	err = s.tm.Do(ctx, func(txCtx context.Context) error {
		if err := s.reviewRepo.Create(txCtx, review); err != nil {
			return fmt.Errorf("failed to create review: %w", err)
		}

		if err := s.restRepo.IncrementRating(txCtx, input.RestaurantID, input.Rating); err != nil {
			return fmt.Errorf("failed to increment restaurant rating: %w", err)
		}

		if err := s.orderRepo.SetReviewed(txCtx, input.OrderID); err != nil {
			return fmt.Errorf("failed to set order reviewed: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return review, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*domain.Review, error) {
	const op = "review.Service.GetAll"

	reviews, err := s.reviewRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return reviews, nil
}

func (s *Service) GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.Review, error) {
	const op = "review.Service.GetByRestaurant"

	reviews, err := s.reviewRepo.GetByRestaurant(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return reviews, nil
}

func (s *Service) GetByRestaurantAll(ctx context.Context, restID uuid.UUID) ([]*domain.Review, error) {
	const op = "review.Service.GetByRestaurantAll"

	reviews, err := s.reviewRepo.GetByRestaurantAll(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return reviews, nil
}

func (s *Service) GetByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Review, error) {
	const op = "review.Service.GetByUser"

	reviews, err := s.reviewRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return reviews, nil
}

func (s *Service) GetByOrder(ctx context.Context, actorID uuid.UUID, orderID uuid.UUID) (*domain.Review, error) {
	const op = "review.Service.GetByOrder"

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Check if the actor is the user who placed the order
	if order.UserID != actorID {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	review, err := s.reviewRepo.GetByOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return review, nil
}

func (s *Service) GetPhotosByRestaurant(ctx context.Context, restID uuid.UUID) ([]*domain.ReviewPhoto, error) {
	const op = "review.Service.GetPhotosByRestaurant"

	photos, err := s.reviewRepo.GetPhotosByRestaurant(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return photos, nil
}

func (s *Service) AddReply(ctx context.Context, actor *domain.RestaurantStaff, reviewID uuid.UUID, reply string) (*domain.Review, error) {
	const op = "review.Service.AddReply"

	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if review.RestaurantID != actor.RestaurantID {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrPermissionDenied)
	}

	review.AddReply(reply)

	if err := s.reviewRepo.SetReply(ctx, review); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return review, nil
}

func (s *Service) Hide(ctx context.Context, reviewID uuid.UUID) error {
	const op = "review.Service.Hide"

	if err := s.reviewRepo.SetVisibility(ctx, reviewID, true); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Unhide(ctx context.Context, reviewID uuid.UUID) error {
	const op = "review.Service.Unhide"

	if err := s.reviewRepo.SetVisibility(ctx, reviewID, false); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
