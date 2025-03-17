package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomLog(t *testing.T) (Log, Event) {
	event := CreateRandomEvent(t)

	args := CreateLogParams{
		EventID:    event.ID,
		ExecutedOn: time.Now(),
		Status:     "SUCCESS",
	}

	log, err := testQueries.CreateLog(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, log)

	require.Equal(t, args.EventID, log.EventID)
	require.Equal(t, args.Status, log.Status)

	return log, event
}

func TestLog(t *testing.T) {
	log, event := CreateRandomLog(t)

	// test get user
	log1, err := testQueries.GetLog(context.Background(), log.ID)
	require.NoError(t, err)
	require.NotEmpty(t, log1)
	require.Equal(t, log.ID, log1.ID)
	require.Equal(t, log.EventID, log1.EventID)
	require.Equal(t, log.ExecutedOn, log1.ExecutedOn)

	// test list users
	arg := ListLogsParams{
		Limit:  1,
		Offset: 0,
	}

	logs, err := testQueries.ListLogs(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, logs)

	// Test delete log
	err = testQueries.DeleteLog(context.Background(), log.ID)
	require.NoError(t, err)

	err = testQueries.DeleteEvent(context.Background(), log.EventID)
	require.NoError(t, err)

	err = testQueries.DeleteUser(context.Background(), event.CreatedBy)
	require.NoError(t, err)
}
