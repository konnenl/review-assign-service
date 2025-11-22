package handler

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/konnen/review-assign-service/internal/dto"
	"github.com/konnen/review-assign-service/internal/errs"
	"github.com/konnen/review-assign-service/internal/model"
	"github.com/konnen/review-assign-service/internal/validator"
)

type mockTeamService struct {
	mock.Mock
}

func (m *mockTeamService) AddTeamWithMembers(ctx context.Context, team model.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *mockTeamService) GetTeamWithMembers(ctx context.Context, teamName string) (model.Team, error) {
	args := m.Called(ctx, teamName)
	return args.Get(0).(model.Team), args.Error(1)
}

func TestTeamHandler_GetTeam(t *testing.T) {
	lgr := slog.New(slog.NewTextHandler(io.Discard, nil))
	v := validator.NewValidator()
	teamServiceMock := new(mockTeamService)

	h := &teamHandler{
		lgr:         lgr,
		validator:   v,
		teamService: teamServiceMock,
	}
	active := true
	tests := []struct {
		name           string
		teamName       string
		mockTeam       model.Team
		mockErr        error
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "validation error missing team_name",
			teamName:       "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   dto.CodeValidationError,
		},
		{
			name:     "team found",
			teamName: "backend",
			mockTeam: model.Team{
				Name: "backend",
				Members: []model.User{{
					ID:       "u1",
					Name:     "Alice",
					IsActive: &active}, {
					ID:       "u2",
					Name:     "Bob",
					IsActive: &active},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "team not found",
			teamName:       "frontend",
			mockErr:        errs.ErrTeamNotFound,
			expectedStatus: http.StatusNotFound,
			expectedCode:   dto.CodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teamServiceMock.ExpectedCalls = nil
			if tt.mockErr != nil {
				teamServiceMock.On("GetTeamWithMembers", mock.Anything, tt.teamName).Return(model.Team{}, tt.mockErr)
			} else if tt.teamName != "" {
				teamServiceMock.On("GetTeamWithMembers", mock.Anything, tt.teamName).Return(tt.mockTeam, nil)
			}

			req := httptest.NewRequest(http.MethodGet, "/team/get?team_name="+tt.teamName, nil)
			w := httptest.NewRecorder()

			h.getTeam(w, req)

			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedStatus != http.StatusOK {
				var errResp dto.ErrorResp
				_ = json.NewDecoder(res.Body).Decode(&errResp)
				require.Equal(t, tt.expectedCode, errResp.Error.Code)
			} else {
				var teamResp dto.TeamDTO
				_ = json.NewDecoder(res.Body).Decode(&teamResp)
				require.Equal(t, tt.mockTeam.Name, teamResp.Name)
				require.Len(t, teamResp.Members, len(tt.mockTeam.Members))
			}

			teamServiceMock.AssertExpectations(t)
		})
	}
}
