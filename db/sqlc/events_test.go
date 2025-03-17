package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/musthafa-vakkayil/event_scheduler_v2/utils"
	"github.com/sqlc-dev/pqtype"
	"github.com/stretchr/testify/require"
)

func CreateRandomEvent(t *testing.T) Event {
	user := CreateRandomUser(t)

	args := CreateEventParams{
		Name:        utils.RandomString(8),
		Type:        "API",
		ApiEndPoint: "test-abc.com/test",
		ApiMethod: sql.NullString{
			String: "GET",
			Valid:  true,
		},
		ApiRequestBody: pqtype.NullRawMessage{
			RawMessage: []byte("{}"), // Empty JSON object
			Valid:      true,
		},
		CreatedBy: user.Username,
	}

	event, err := testQueries.CreateEvent(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, event)

	require.Equal(t, args.Name, event.Name)
	require.Equal(t, args.ApiMethod, event.ApiMethod)
	require.Equal(t, args.ApiEndPoint, event.ApiEndPoint)
	require.Equal(t, args.Type, event.Type)
	require.Equal(t, args.CreatedBy, event.CreatedBy)

	require.NotNil(t, event.CreatedAt)

	return event
}

func TestEvent(t *testing.T) {
	event := CreateRandomEvent(t)

	// test get user
	event1, err := testQueries.GetEvent(context.Background(), event.ID)
	require.NoError(t, err)
	require.NotEmpty(t, event1)
	require.Equal(t, event.Name, event1.Name)
	require.Equal(t, event.Type, event1.Type)
	require.Equal(t, event.ApiMethod, event1.ApiMethod)

	// test list users
	arg := ListEventsParams{
		Limit:  1,
		Offset: 0,
	}

	events, err := testQueries.ListEvents(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, events)

	// Test delete event
	err = testQueries.DeleteEvent(context.Background(), event.ID)
	require.NoError(t, err)

	err = testQueries.DeleteUser(context.Background(), event.CreatedBy)
	require.NoError(t, err)
}
