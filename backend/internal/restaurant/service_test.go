package restaurant_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/restaurant"
	"gitlab.com/2foood/backend/internal/restaurant/mocks"
)

func setupTest() (*mocks.UserRepositoryMock, *mocks.RestaurantRepositoryMock, *mocks.StaffRepositoryMock, *restaurant.Service) {
	tm := &mocks.TransactionManagerFake{}
	userRepo := &mocks.UserRepositoryMock{}
	restRepo := &mocks.RestaurantRepositoryMock{}
	staffRepo := &mocks.StaffRepositoryMock{}
	svc := restaurant.NewService(tm, userRepo, restRepo, staffRepo)
	return userRepo, restRepo, staffRepo, svc
}

func TestRestaurantService_Create(t *testing.T) {
	ownerID := uuid.New()
	description := "Good food"
	input := domain.CreateRestaurantInput{
		Name:        "Test Rest",
		Description: &description,
	}

	tests := []struct {
		name    string
		mockFn  func(userRepo *mocks.UserRepositoryMock, restRepo *mocks.RestaurantRepositoryMock, staffRepo *mocks.StaffRepositoryMock)
		wantErr bool
	}{
		{
			name: "failed to create restaurant in db",
			mockFn: func(u *mocks.UserRepositoryMock, r *mocks.RestaurantRepositoryMock, s *mocks.StaffRepositoryMock) {
				r.On("Create", mock.Anything, mock.AnythingOfType("*domain.Restaurant")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success (owner is admin, role not updated)",
			mockFn: func(u *mocks.UserRepositoryMock, r *mocks.RestaurantRepositoryMock, s *mocks.StaffRepositoryMock) {
				r.On("Create", mock.Anything, mock.AnythingOfType("*domain.Restaurant")).Return(nil)
				s.On("Add", mock.Anything, mock.AnythingOfType("*domain.RestaurantStaff")).Return(nil)

				u.On("GetByID", mock.Anything, ownerID).Return(&domain.User{Role: domain.UserRoleAdmin}, nil)
			},
			wantErr: false,
		},
		{
			name: "success (owner is customer, upgraded to partner)",
			mockFn: func(u *mocks.UserRepositoryMock, r *mocks.RestaurantRepositoryMock, s *mocks.StaffRepositoryMock) {
				r.On("Create", mock.Anything, mock.AnythingOfType("*domain.Restaurant")).Return(nil)
				s.On("Add", mock.Anything, mock.AnythingOfType("*domain.RestaurantStaff")).Return(nil)

				u.On("GetByID", mock.Anything, ownerID).Return(&domain.User{Role: domain.UserRoleCustomer}, nil)
				u.On("SetRole", mock.Anything, ownerID, domain.UserRolePartner).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, restRepo, staffRepo, svc := setupTest()
			tt.mockFn(userRepo, restRepo, staffRepo)

			rest, err := svc.Create(context.Background(), ownerID, input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, rest)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rest)
				assert.Equal(t, input.Name, rest.Name)
			}
			userRepo.AssertExpectations(t)
			restRepo.AssertExpectations(t)
			staffRepo.AssertExpectations(t)
		})
	}
}

func TestRestaurantService_RestGetters(t *testing.T) {
	restID := uuid.New()
	userID := uuid.New()

	t.Run("GetAll_Success", func(t *testing.T) {
		t.Parallel()
		_, restRepo, _, svc := setupTest()

		restRepo.On("GetAll", mock.Anything).Return([]*domain.Restaurant{{ID: restID}}, nil)
		res, err := svc.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("GetByPartner_Success", func(t *testing.T) {
		t.Parallel()
		_, restRepo, _, svc := setupTest()

		restRepo.On("GetByPartner", mock.Anything, userID).Return([]*domain.Restaurant{{ID: restID}}, nil)
		res, err := svc.GetByPartner(context.Background(), userID)

		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("GetByID_Success", func(t *testing.T) {
		t.Parallel()
		_, restRepo, _, svc := setupTest()

		restRepo.On("GetByID", mock.Anything, restID).Return(&domain.Restaurant{ID: restID}, nil)
		res, err := svc.GetByID(context.Background(), restID)

		assert.NoError(t, err)
		assert.Equal(t, restID, res.ID)
	})
}

func TestRestaurantService_UpdateInfo(t *testing.T) {
	restID := uuid.New()
	input := domain.UpdateRestaurantInput{}
	actorOwner := &domain.RestaurantStaff{Role: domain.StaffRoleOwner}
	actorManager := &domain.RestaurantStaff{Role: domain.StaffRoleManager}

	tests := []struct {
		name    string
		actor   *domain.RestaurantStaff
		mockFn  func(restRepo *mocks.RestaurantRepositoryMock)
		wantErr bool
	}{
		{
			name:    "permission denied (not owner)",
			actor:   actorManager,
			mockFn:  func(r *mocks.RestaurantRepositoryMock) {},
			wantErr: true,
		},
		{
			name:  "restaurant not found",
			actor: actorOwner,
			mockFn: func(r *mocks.RestaurantRepositoryMock) {
				r.On("GetByID", mock.Anything, restID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name:  "success",
			actor: actorOwner,
			mockFn: func(r *mocks.RestaurantRepositoryMock) {
				r.On("GetByID", mock.Anything, restID).Return(&domain.Restaurant{ID: restID}, nil)
				r.On("UpdateInfo", mock.Anything, mock.AnythingOfType("*domain.Restaurant")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, restRepo, _, svc := setupTest()
			tt.mockFn(restRepo)

			rest, err := svc.UpdateInfo(context.Background(), tt.actor, restID, input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, rest)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rest)
			}
			restRepo.AssertExpectations(t)
		})
	}
}

func TestRestaurantService_UpdateCuisines(t *testing.T) {
    restID := uuid.New()
    actorOwner := &domain.RestaurantStaff{Role: domain.StaffRoleOwner}
    actorManager := &domain.RestaurantStaff{Role: domain.StaffRoleManager}
    cuisineIDs := []int{1, 2, 3}

    tests := []struct {
        name    string
        actor   *domain.RestaurantStaff
        mockFn  func(restRepo *mocks.RestaurantRepositoryMock)
        wantErr bool
    }{
        {
            name:    "permission denied (not owner)",
            actor:   actorManager,
            mockFn:  func(r *mocks.RestaurantRepositoryMock) {},
            wantErr: true,
        },
        {
            name:  "db error",
            actor: actorOwner,
            mockFn: func(r *mocks.RestaurantRepositoryMock) {
                r.On("UpdateCuisines", mock.Anything, restID, cuisineIDs).Return(errors.New("db error"))
            },
            wantErr: true,
        },
        {
            name:  "success",
            actor: actorOwner,
            mockFn: func(r *mocks.RestaurantRepositoryMock) {
                r.On("UpdateCuisines", mock.Anything, restID, cuisineIDs).Return(nil)
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            _, restRepo, _, svc := setupTest()
            tt.mockFn(restRepo)

            err := svc.UpdateCuisines(context.Background(), tt.actor, restID, cuisineIDs)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
            restRepo.AssertExpectations(t)
        })
    }
}

func TestRestaurantService_Delete(t *testing.T) {
	restID := uuid.New()

	tests := []struct {
		name    string
		mockFn  func(restRepo *mocks.RestaurantRepositoryMock)
		wantErr bool
	}{
		{
			name: "db error",
			mockFn: func(r *mocks.RestaurantRepositoryMock) {
				r.On("Delete", mock.Anything, restID).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(r *mocks.RestaurantRepositoryMock) {
				r.On("Delete", mock.Anything, restID).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, restRepo, _, svc := setupTest()
			tt.mockFn(restRepo)

			err := svc.Delete(context.Background(), restID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			restRepo.AssertExpectations(t)
		})
	}
}

func TestRestaurantService_AddStaff(t *testing.T) {
	restID := uuid.New()
	targetID := uuid.New()
	actorOwner := &domain.RestaurantStaff{Role: domain.StaffRoleOwner}
	actorManager := &domain.RestaurantStaff{Role: domain.StaffRoleManager}

	tests := []struct {
		name    string
		actor   *domain.RestaurantStaff
		role    domain.StaffRole
		mockFn  func(userRepo *mocks.UserRepositoryMock, staffRepo *mocks.StaffRepositoryMock)
		wantErr bool
	}{
		{
			name:    "permission denied (actor is not owner)",
			actor:   actorManager,
			role:    domain.StaffRoleManager,
			mockFn:  func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {},
			wantErr: true,
		},
		{
			name:  "target already exists",
			actor: actorOwner,
			role:  domain.StaffRoleManager,
			mockFn: func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				s.On("GetMember", mock.Anything, restID, targetID).
					Return(&domain.RestaurantStaff{}, nil)
			},
			wantErr: true,
		},
		{
			name:  "cannot add owner directly",
			actor: actorOwner,
			role:  domain.StaffRoleOwner,
			mockFn: func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				s.On("GetMember", mock.Anything, restID, targetID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name:  "success (target is admin, role not updated)",
			actor: actorOwner,
			role:  domain.StaffRoleManager,
			mockFn: func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				s.On("GetMember", mock.Anything, restID, targetID).Return(nil, domain.ErrNotFound)
				s.On("Add", mock.Anything, mock.AnythingOfType("*domain.RestaurantStaff")).Return(nil)

				u.On("GetByID", mock.Anything, targetID).
					Return(&domain.User{Role: domain.UserRoleAdmin}, nil)
			},
			wantErr: false,
		},
		{
			name:  "success (target is customer, becomes partner)",
			actor: actorOwner,
			role:  domain.StaffRoleManager,
			mockFn: func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				s.On("GetMember", mock.Anything, restID, targetID).Return(nil, domain.ErrNotFound)
				s.On("Add", mock.Anything, mock.AnythingOfType("*domain.RestaurantStaff")).Return(nil)

				u.On("GetByID", mock.Anything, targetID).
					Return(&domain.User{Role: domain.UserRoleCustomer}, nil)
				u.On("SetRole", mock.Anything, targetID, domain.UserRolePartner).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, _, staffRepo, svc := setupTest()
			tt.mockFn(userRepo, staffRepo)

			err := svc.AddStaff(context.Background(), tt.actor, restID, targetID, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			userRepo.AssertExpectations(t)
			staffRepo.AssertExpectations(t)
		})
	}
}

func TestRestaurantService_UpdateStaffRole(t *testing.T) {
	restID := uuid.New()
	targetID := uuid.New()
	actorOwner := &domain.RestaurantStaff{Role: domain.StaffRoleOwner}

	tests := []struct {
		name    string
		role    domain.StaffRole
		mockFn  func(staffRepo *mocks.StaffRepositoryMock)
		wantErr bool
	}{
		{
			name: "target not found",
			role: domain.StaffRoleManager,
			mockFn: func(s *mocks.StaffRepositoryMock) {
				s.On("GetMember", mock.Anything, restID, targetID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "cannot update owner",
			role: domain.StaffRoleManager,
			mockFn: func(s *mocks.StaffRepositoryMock) {
				s.On("GetMember", mock.Anything, restID, targetID).
					Return(&domain.RestaurantStaff{Role: domain.StaffRoleOwner}, nil)
			},
			wantErr: true,
		},
		{
			name: "cannot set to owner",
			role: domain.StaffRoleOwner,
			mockFn: func(s *mocks.StaffRepositoryMock) {
				s.On("GetMember", mock.Anything, restID, targetID).
					Return(&domain.RestaurantStaff{Role: domain.StaffRoleManager}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, _, staffRepo, svc := setupTest()
			tt.mockFn(staffRepo)

			err := svc.UpdateStaffRole(context.Background(), actorOwner, restID, targetID, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			staffRepo.AssertExpectations(t)
		})
	}
}

func TestRestaurantService_TransferOwnership(t *testing.T) {
	restID := uuid.New()
	actorID := uuid.New()
	targetID := uuid.New()
	actorOwner := &domain.RestaurantStaff{UserID: actorID, Role: domain.StaffRoleOwner}

	actorPartnerUser := &domain.User{Role: domain.UserRolePartner}

	tests := []struct {
		name    string
		actor   *domain.RestaurantStaff
		mockFn  func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock)
		wantErr bool
	}{
		{
			name:  "target is already staff (update role)",
			actor: actorOwner,
			mockFn: func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				u.On("GetByID", mock.Anything, actorID).Return(actorPartnerUser, nil)

				s.On("SetRole", mock.Anything, restID, actorID, domain.StaffRoleManager).Return(nil)
				s.On("SetRole", mock.Anything, restID, targetID, domain.StaffRoleOwner).Return(nil)

				u.On("GetByID", mock.Anything, targetID).Return(&domain.User{Role: domain.UserRoleCustomer}, nil)
				u.On("SetRole", mock.Anything, targetID, domain.UserRolePartner).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "target is NOT staff (add new owner)",
			actor: actorOwner,
			mockFn: func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				u.On("GetByID", mock.Anything, actorID).Return(actorPartnerUser, nil)

				s.On("SetRole", mock.Anything, restID, actorID, domain.StaffRoleManager).Return(nil)
				s.On("SetRole", mock.Anything, restID, targetID, domain.StaffRoleOwner).Return(domain.ErrNotFound)
				s.On("Add", mock.Anything, mock.AnythingOfType("*domain.RestaurantStaff")).Return(nil)

				u.On("GetByID", mock.Anything, targetID).Return(&domain.User{Role: domain.UserRoleCustomer}, nil)
				u.On("SetRole", mock.Anything, targetID, domain.UserRolePartner).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "admin platform user cannot transfer via this method",
			actor: actorOwner,
			mockFn: func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				u.On("GetByID", mock.Anything, actorID).Return(&domain.User{Role: domain.UserRoleAdmin}, nil)
			},
			wantErr: true,
		},
		{
			name:  "actor is manager not owner",
			actor: &domain.RestaurantStaff{UserID: actorID, Role: domain.StaffRoleManager},
			mockFn: func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				u.On("GetByID", mock.Anything, actorID).Return(actorPartnerUser, nil)
			},
			wantErr: true,
		},
		{
			name:  "transfer to self",
			actor: actorOwner,
			mockFn: func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				u.On("GetByID", mock.Anything, actorID).Return(actorPartnerUser, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, _, staffRepo, svc := setupTest()
			tid := targetID
			if tt.name == "transfer to self" {
				tid = actorID
			}
			tt.mockFn(userRepo, staffRepo)

			err := svc.TransferOwnership(context.Background(), tt.actor, restID, tid)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			userRepo.AssertExpectations(t)
			staffRepo.AssertExpectations(t)
		})
	}
}

func TestRestaurantService_RemoveStaff(t *testing.T) {
	restID := uuid.New()
	targetID := uuid.New()
	actorOwner := &domain.RestaurantStaff{Role: domain.StaffRoleOwner}

	tests := []struct {
		name    string
		mockFn  func(userRepo *mocks.UserRepositoryMock, staffRepo *mocks.StaffRepositoryMock)
		wantErr bool
	}{
		{
			name: "success (remove manager, downgrade to customer)",
			mockFn: func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				s.On("GetMember", mock.Anything, restID, targetID).
					Return(&domain.RestaurantStaff{Role: domain.StaffRoleManager}, nil)

				s.On("Remove", mock.Anything, restID, targetID).Return(nil)

				s.On("CountRestaurants", mock.Anything, targetID).Return(0, nil)

				u.On("GetByID", mock.Anything, targetID).
					Return(&domain.User{Role: domain.UserRolePartner}, nil)

				u.On("SetRole", mock.Anything, targetID, domain.UserRoleCustomer).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success (remove manager, keep partner role if still in other restaurants)",
			mockFn: func(u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				s.On("GetMember", mock.Anything, restID, targetID).
					Return(&domain.RestaurantStaff{Role: domain.StaffRoleManager}, nil)

				s.On("Remove", mock.Anything, restID, targetID).Return(nil)

				s.On("CountRestaurants", mock.Anything, targetID).Return(1, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, _, staffRepo, svc := setupTest()
			tt.mockFn(userRepo, staffRepo)

			err := svc.RemoveStaff(context.Background(), actorOwner, restID, targetID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			userRepo.AssertExpectations(t)
			staffRepo.AssertExpectations(t)
		})
	}
}

func TestRestaurantService_StaffGetters(t *testing.T) {
	restID := uuid.New()
	userID := uuid.New()

	t.Run("GetStaff_Success", func(t *testing.T) {
		t.Parallel()
		_, _, staffRepo, svc := setupTest()

		staffRepo.On("GetByRestaurant", mock.Anything, restID).Return([]*domain.RestaurantStaff{{}}, nil)
		res, err := svc.GetStaff(context.Background(), restID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("GetStaffMember_Success", func(t *testing.T) {
		t.Parallel()
		_, _, staffRepo, svc := setupTest()

		staffRepo.On("GetMember", mock.Anything, restID, userID).Return(&domain.RestaurantStaff{UserID: userID}, nil)
		res, err := svc.GetStaffMember(context.Background(), restID, userID)
		assert.NoError(t, err)
		assert.Equal(t, userID, res.UserID)
	})
}

func TestRestaurantService_AssignOwner(t *testing.T) {
	restID := uuid.New()
	ownerUserID := uuid.New()
	targetID := uuid.New()

	tests := []struct {
		name    string
		mockFn  func(restRepo *mocks.RestaurantRepositoryMock, userRepo *mocks.UserRepositoryMock, staffRepo *mocks.StaffRepositoryMock)
		wantErr bool
	}{
		{
			name: "restaurant not found",
			mockFn: func(r *mocks.RestaurantRepositoryMock, u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				r.On("GetByID", mock.Anything, restID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "get owner fails (non-NotFound)",
			mockFn: func(r *mocks.RestaurantRepositoryMock, u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				r.On("GetByID", mock.Anything, restID).Return(&domain.Restaurant{ID: restID}, nil)
				s.On("GetOwner", mock.Anything, restID).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "no current owner, target not staff (add owner + set partner)",
			mockFn: func(r *mocks.RestaurantRepositoryMock, u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				r.On("GetByID", mock.Anything, restID).Return(&domain.Restaurant{ID: restID}, nil)
				s.On("GetOwner", mock.Anything, restID).Return(nil, domain.ErrNotFound)
				s.On("GetMember", mock.Anything, restID, targetID).Return(nil, domain.ErrNotFound)
				s.On("Add", mock.Anything, mock.AnythingOfType("*domain.RestaurantStaff")).Return(nil)
				u.On("GetByID", mock.Anything, targetID).Return(&domain.User{Role: domain.UserRoleCustomer}, nil)
				u.On("SetRole", mock.Anything, targetID, domain.UserRolePartner).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "demote current owner, target already staff (set owner + partner)",
			mockFn: func(r *mocks.RestaurantRepositoryMock, u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				r.On("GetByID", mock.Anything, restID).Return(&domain.Restaurant{ID: restID}, nil)
				s.On("GetOwner", mock.Anything, restID).Return(&domain.RestaurantStaff{UserID: ownerUserID, Role: domain.StaffRoleOwner}, nil)
				s.On("SetRole", mock.Anything, restID, ownerUserID, domain.StaffRoleManager).Return(nil)
				s.On("GetMember", mock.Anything, restID, targetID).Return(&domain.RestaurantStaff{UserID: targetID}, nil)
				s.On("SetRole", mock.Anything, restID, targetID, domain.StaffRoleOwner).Return(nil)
				u.On("GetByID", mock.Anything, targetID).Return(&domain.User{Role: domain.UserRoleCustomer}, nil)
				u.On("SetRole", mock.Anything, targetID, domain.UserRolePartner).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "target is admin (skip user SetRole)",
			mockFn: func(r *mocks.RestaurantRepositoryMock, u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				r.On("GetByID", mock.Anything, restID).Return(&domain.Restaurant{ID: restID}, nil)
				s.On("GetOwner", mock.Anything, restID).Return(nil, domain.ErrNotFound)
				s.On("GetMember", mock.Anything, restID, targetID).Return(nil, domain.ErrNotFound)
				s.On("Add", mock.Anything, mock.AnythingOfType("*domain.RestaurantStaff")).Return(nil)
				u.On("GetByID", mock.Anything, targetID).Return(&domain.User{Role: domain.UserRoleAdmin}, nil)
			},
			wantErr: false,
		},
		{
			name: "target already partner (skip user SetRole)",
			mockFn: func(r *mocks.RestaurantRepositoryMock, u *mocks.UserRepositoryMock, s *mocks.StaffRepositoryMock) {
				r.On("GetByID", mock.Anything, restID).Return(&domain.Restaurant{ID: restID}, nil)
				s.On("GetOwner", mock.Anything, restID).Return(nil, domain.ErrNotFound)
				s.On("GetMember", mock.Anything, restID, targetID).Return(&domain.RestaurantStaff{UserID: targetID}, nil)
				s.On("SetRole", mock.Anything, restID, targetID, domain.StaffRoleOwner).Return(nil)
				u.On("GetByID", mock.Anything, targetID).Return(&domain.User{Role: domain.UserRolePartner}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, restRepo, staffRepo, svc := setupTest()
			tt.mockFn(restRepo, userRepo, staffRepo)

			err := svc.AssignOwner(context.Background(), restID, targetID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			userRepo.AssertExpectations(t)
			restRepo.AssertExpectations(t)
			staffRepo.AssertExpectations(t)
		})
	}
}
