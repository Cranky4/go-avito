package memorystorage

import (
	"errors"
	"testing"
	"time"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("add event", func(t *testing.T) {
		st := New()
		newEvent := storage.Event{
			ID:       storage.NewEventID(),
			Title:    "new event",
			StartsAt: time.Now(),
			EndsAt:   time.Now().Add(time.Hour),
		}
		st.CreateEvent(newEvent)

		event, err := st.GetEvent(newEvent.ID)
		require.NotNil(t, event)
		require.NoError(t, err)

		newEvent2 := storage.Event{
			ID:       storage.NewEventID(),
			Title:    "new event",
			StartsAt: time.Now(),
			EndsAt:   time.Now().Add(time.Hour),
		}
		err = st.CreateEvent(newEvent2)
		require.Error(t, err)
		require.True(t, errors.Is(storage.ErrDateBusy, err))

		event, err = st.GetEvent(newEvent2.ID)
		require.NotNil(t, event)
		require.Error(t, err)
		require.True(t, errors.Is(storage.ErrEventNotFound, err))
	})

	t.Run("update event", func(t *testing.T) {
		st := New()
		newEvent1 := storage.Event{
			ID:       storage.NewEventID(),
			Title:    "new event",
			StartsAt: time.Now(),
			EndsAt:   time.Now().Add(time.Hour),
		}
		err := st.CreateEvent(newEvent1)
		require.NoError(t, err)

		newEvent2 := storage.Event{
			ID:       newEvent1.ID,
			Title:    "new title",
			StartsAt: time.Now().Add(time.Minute),
			EndsAt:   time.Now().Add(time.Hour),
		}
		err = st.UpdateEvent(newEvent1.ID, newEvent2)
		require.NoError(t, err)

		event, err := st.GetEvent(newEvent1.ID)

		require.NoError(t, err)
		require.Equal(t, "new title", event.Title)

		err = st.UpdateEvent(storage.NewEventID(), newEvent2)
		require.Error(t, err)
		require.True(t, errors.Is(storage.ErrEventNotFound, err))

		newEvent3 := storage.Event{
			ID:       storage.NewEventID(),
			Title:    "new event 3",
			StartsAt: time.Now().Add(2 * time.Hour),
			EndsAt:   time.Now().Add(3 * time.Hour),
		}
		err = st.CreateEvent(newEvent3)
		require.NoError(t, err)

		err = st.UpdateEvent(newEvent3.ID, storage.Event{
			ID:       newEvent1.ID,
			Title:    "new event 3",
			StartsAt: time.Now(),
			EndsAt:   time.Now().Add(3 * time.Hour),
		})

		require.Error(t, err)
		require.True(t, errors.Is(storage.ErrDateBusy, err))
	})

	t.Run("delete event", func(t *testing.T) {
		st := New()
		newEvent1 := storage.Event{
			ID:       storage.NewEventID(),
			Title:    "new event",
			StartsAt: time.Now(),
			EndsAt:   time.Now().Add(time.Hour),
		}
		err := st.CreateEvent(newEvent1)
		require.NoError(t, err)
		event, err := st.GetEvent(newEvent1.ID)

		require.NoError(t, err)
		require.NotNil(t, event)

		err = st.DeleteEvent(newEvent1.ID)
		require.NoError(t, err)

		_, err = st.GetEvent(newEvent1.ID)
		require.Error(t, err)
		require.True(t, errors.Is(storage.ErrEventNotFound, err))

		err = st.DeleteEvent(newEvent1.ID)
		require.Error(t, err)
		require.True(t, errors.Is(storage.ErrEventNotFound, err))
	})

	t.Run("get events", func(t *testing.T) {
		st := New()

		events, err := st.GetEvents(time.Now(), time.Now().Add(time.Hour))
		require.Empty(t, events)
		require.NoError(t, err)

		newEvent1 := storage.Event{
			ID:       storage.NewEventID(),
			Title:    "new event",
			StartsAt: time.Now().Add(time.Second),
			EndsAt:   time.Now().Add(10 * time.Minute),
		}
		err = st.CreateEvent(newEvent1)
		require.NoError(t, err)

		newEvent2 := storage.Event{
			ID:       storage.NewEventID(),
			Title:    "new event 2",
			StartsAt: time.Now().Add(11 * time.Minute),
			EndsAt:   time.Now().Add(30 * time.Minute),
		}
		err = st.CreateEvent(newEvent2)
		require.NoError(t, err)

		newEvent3 := storage.Event{
			ID:       storage.NewEventID(),
			Title:    "new event 3",
			StartsAt: time.Now().Add(2 * time.Hour),
			EndsAt:   time.Now().Add(3 * time.Hour),
		}
		err = st.CreateEvent(newEvent3)
		require.NoError(t, err)

		events, err = st.GetEvents(time.Now(), time.Now().Add(60*time.Minute))
		require.Len(t, events, 2)
		require.NoError(t, err)
	})
}
