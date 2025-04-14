package notifications

import (
	"context"
	"github.com/Tagliatti/magalu-challenge/database"
	"github.com/Tagliatti/magalu-challenge/testhelpers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PostgresRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *PostgresRepository
	db          *pgxpool.Pool
	ctx         context.Context
}

func (suite *PostgresRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := testhelpers.NewPostgresContainer(suite.ctx)
	require.Nil(suite.T(), err, "failed to start postgres container: %v", err)

	suite.pgContainer = pgContainer

	db, err := database.ConnectTest(suite.ctx, pgContainer.ConnectionString)
	require.Nil(suite.T(), err, "failed to connect to database: %v", err)

	suite.db = db
	suite.repository = NewPostgresRepository(db)
}

func (suite *PostgresRepositoryTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.T().Fatalf("failed to terminate pgContainer: %s", err)
	}
}

func TestPostgresRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepositoryTestSuite))
}

func (suite *PostgresRepositoryTestSuite) TestSuccessCreateNotification() {
	t := suite.T()

	t.Run("Should create and find notification successfully", func(t *testing.T) {
		err := testhelpers.TruncateAllTables(suite.ctx, suite.db)
		require.Nilf(t, err, "failed to truncate tables: %v", err)

		createNotification := &CreateNotification{
			Type:      "email",
			Recipient: "test@example.com",
		}

		id, err := suite.repository.CreateNotification(suite.ctx, createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		notification, err := suite.repository.FindNotificationByID(suite.ctx, id)
		require.Nilf(t, err, "failed to find notification by ID: %v", err)

		assert.NotNil(t, notification, "notification should not be nil")
		assert.Equal(t, id, notification.Id)
		assert.Equal(t, createNotification.Type, notification.Type)
		assert.Equal(t, createNotification.Recipient, notification.Recipient)
		assert.NotEmpty(t, notification.CreatedAt)
		assert.False(t, notification.Sent)
		assert.Nil(t, notification.SentAt)
	})
}

func (suite *PostgresRepositoryTestSuite) TestErrorOnCreateNotificationWithInvalidType() {
	t := suite.T()

	t.Run("Should return error when creating notification with invalid type", func(t *testing.T) {
		err := testhelpers.TruncateAllTables(suite.ctx, suite.db)
		require.Nilf(t, err, "failed to truncate tables: %v", err)

		createNotification := &CreateNotification{
			Type:      "invalid",
			Recipient: "test@example.com",
		}

		id, err := suite.repository.CreateNotification(suite.ctx, createNotification)
		assert.NotNil(t, err)
		assert.Zero(t, id)
	})
}

func (suite *PostgresRepositoryTestSuite) TestSuccessFindNotificationStatus() {
	t := suite.T()

	t.Run("Should find notification status successfully", func(t *testing.T) {
		err := testhelpers.TruncateAllTables(suite.ctx, suite.db)
		require.Nilf(t, err, "failed to truncate tables: %v", err)

		createNotification := &CreateNotification{
			Type:      "email",
			Recipient: "test@example.com",
		}

		id, err := suite.repository.CreateNotification(suite.ctx, createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		notificationStatus, err := suite.repository.FindNotificationStatusByID(suite.ctx, id)
		require.Nilf(t, err, "failed to find notification by ID: %v", err)

		assert.NotNil(t, notificationStatus, "notification should not be nil")
		assert.False(t, notificationStatus.Sent)
		assert.Nil(t, notificationStatus.SentAt)
	})
}

func (suite *PostgresRepositoryTestSuite) TestNotFoundOnFindNotificationStatus() {
	t := suite.T()

	t.Run("Should not find notification status", func(t *testing.T) {
		err := testhelpers.TruncateAllTables(suite.ctx, suite.db)
		require.Nilf(t, err, "failed to truncate tables: %v", err)

		createNotification := &CreateNotification{
			Type:      "email",
			Recipient: "test@example.com",
		}

		id, err := suite.repository.CreateNotification(suite.ctx, createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		notificationStatus, err := suite.repository.FindNotificationStatusByID(suite.ctx, id+1)
		require.Nilf(t, err, "failed to find notification by ID: %v", err)

		assert.Nil(t, notificationStatus)
	})
}

func (suite *PostgresRepositoryTestSuite) TestSuccessFindSentNotificationStatus() {
	t := suite.T()

	t.Run("Should find sent notification status successfully", func(t *testing.T) {
		err := testhelpers.TruncateAllTables(suite.ctx, suite.db)
		require.Nilf(t, err, "failed to truncate tables: %v", err)

		createNotification := &CreateNotification{
			Type:      "email",
			Recipient: "test@example.com",
		}

		id, err := suite.repository.CreateNotification(suite.ctx, createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		updated, err := suite.repository.UpdateNotificationAsSent(suite.ctx, id)
		require.Nilf(t, err, "failed to update notification as sent: %v", err)

		notificationStatus, err := suite.repository.FindNotificationStatusByID(suite.ctx, id)
		require.Nilf(t, err, "failed to find notification by ID: %v", err)
		require.NotNil(t, notificationStatus, "notification should not be nil")

		assert.True(t, updated)
		assert.True(t, notificationStatus.Sent)
		assert.NotNil(t, notificationStatus.SentAt)
	})
}

func (suite *PostgresRepositoryTestSuite) TestSuccessDeleteNotification() {
	t := suite.T()

	t.Run("Should delete notification successfully", func(t *testing.T) {
		err := testhelpers.TruncateAllTables(suite.ctx, suite.db)
		require.Nilf(t, err, "failed to truncate tables: %v", err)

		createNotification := &CreateNotification{
			Type:      "email",
			Recipient: "test@example.com",
		}

		id, err := suite.repository.CreateNotification(suite.ctx, createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		deleted, err := suite.repository.DeleteNotificationByID(suite.ctx, id)
		require.Nilf(t, err, "failed to delete notification: %v", err)

		notification, err := suite.repository.FindNotificationByID(suite.ctx, id)
		require.Nilf(t, err, "failed to find notification by ID: %v", err)

		assert.True(t, deleted)
		assert.Nil(t, notification)
	})
}

func (suite *PostgresRepositoryTestSuite) TestNotFoundOnDeleteNotification() {
	t := suite.T()

	t.Run("Should not delete notification", func(t *testing.T) {
		err := testhelpers.TruncateAllTables(suite.ctx, suite.db)
		require.Nilf(t, err, "failed to truncate tables: %v", err)

		createNotification := &CreateNotification{
			Type:      "email",
			Recipient: "test@example.com",
		}

		id, err := suite.repository.CreateNotification(suite.ctx, createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		deleted, err := suite.repository.DeleteNotificationByID(suite.ctx, id+1)
		require.Nilf(t, err, "failed to delete notification: %v", err)

		notification, err := suite.repository.FindNotificationByID(suite.ctx, id)
		require.Nilf(t, err, "failed to find notification by ID: %v", err)

		assert.False(t, deleted)
		assert.NotNil(t, notification)
	})
}

func (suite *PostgresRepositoryTestSuite) TestSuccessListPendingNotifications() {
	t := suite.T()

	t.Run("Should list only pending notifications", func(t *testing.T) {
		err := testhelpers.TruncateAllTables(suite.ctx, suite.db)
		require.Nilf(t, err, "failed to truncate tables: %v", err)

		pendingNotification := &CreateNotification{
			Type:      "email",
			Recipient: "pending@example.com",
		}
		_, err = suite.repository.CreateNotification(suite.ctx, pendingNotification)
		require.Nilf(t, err, "failed to create pending notification: %v", err)

		sentNotification := &CreateNotification{
			Type:      "email",
			Recipient: "sent@example.com",
		}
		sentID, err := suite.repository.CreateNotification(suite.ctx, sentNotification)
		require.Nilf(t, err, "failed to create sent notification: %v", err)

		_, err = suite.repository.UpdateNotificationAsSent(suite.ctx, sentID)
		require.Nilf(t, err, "failed to update notification as sent: %v", err)

		pendingNotifications, err := suite.repository.ListPendingNotifications(suite.ctx)
		require.Nilf(t, err, "failed to list pending notifications: %v", err)

		assert.Len(t, pendingNotifications, 1)
		assert.Equal(t, pendingNotification.Recipient, pendingNotifications[0].Recipient)
		assert.Equal(t, pendingNotification.Type, pendingNotifications[0].Type)
		assert.False(t, pendingNotifications[0].Sent)
		assert.Nil(t, pendingNotifications[0].SentAt)
	})
}
