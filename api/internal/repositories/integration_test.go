//go:build integration

package repositories

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mr-pixel-kg/shopshredder/api/internal/models"
	"github.com/mr-pixel-kg/shopshredder/api/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateAndFind(t *testing.T) {
	db := testutil.OpenIntegrationDB(t)
	testutil.ResetIntegrationDB(t, db)

	repo := NewUserRepository(db)
	user := &models.User{
		ID:           uuid.New(),
		Email:        fmt.Sprintf("user-%d@example.com", time.Now().UnixNano()),
		PasswordHash: "hashed-password",
	}

	require.NoError(t, repo.Create(user))

	byEmail, err := repo.FindByEmail(user.Email)
	require.NoError(t, err)
	assert.Equal(t, user.ID, byEmail.ID)

	byID, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.Email, byID.Email)
}
