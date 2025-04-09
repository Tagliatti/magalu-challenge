package notifications

import (
	"context"
	"database/sql"
	"github.com/Tagliatti/magalu-challenge/database"
	"github.com/Tagliatti/magalu-challenge/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PostgresRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *PostgresRepository
	db          *sql.DB
	ctx         context.Context
}

func (suite *PostgresRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := testhelpers.NewPostgresContainer(suite.ctx)
	require.Nil(suite.T(), err, "failed to start postgres container: %v", err)

	suite.pgContainer = pgContainer

	db, err := database.ConnectTest(pgContainer.ConnectionString)
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

		id, err := suite.repository.CreateNotification(createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		notification, err := suite.repository.FindNotificationByID(id)
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

		id, err := suite.repository.CreateNotification(createNotification)
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

		id, err := suite.repository.CreateNotification(createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		notificationStatus, err := suite.repository.FindNotificationStatusByID(id)
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

		id, err := suite.repository.CreateNotification(createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		notificationStatus, err := suite.repository.FindNotificationStatusByID(id + 1)
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

		id, err := suite.repository.CreateNotification(createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		updated, err := suite.repository.UpdateNotificationAsSent(id)
		require.Nilf(t, err, "failed to update notification as sent: %v", err)

		notificationStatus, err := suite.repository.FindNotificationStatusByID(id)
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

		id, err := suite.repository.CreateNotification(createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		deleted, err := suite.repository.DeleteNotificationByID(id)
		require.Nilf(t, err, "failed to delete notification: %v", err)

		notification, err := suite.repository.FindNotificationByID(id)
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

		id, err := suite.repository.CreateNotification(createNotification)
		require.Nilf(t, err, "failed to create notification: %v", err)

		deleted, err := suite.repository.DeleteNotificationByID(id + 1)
		require.Nilf(t, err, "failed to delete notification: %v", err)

		notification, err := suite.repository.FindNotificationByID(id)
		require.Nilf(t, err, "failed to find notification by ID: %v", err)

		assert.False(t, deleted)
		assert.NotNil(t, notification)
	})
}
