package internalhttp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/app"
	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestEventApiHandlerCreateErrors(t *testing.T) {
	t.Run("empty request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events", nil)
		w := httptest.NewRecorder()

		logg := logger.New("error", 0)
		calendar := app.New(logg, memorystorage.New())
		handler := NewEventAPIHandler(calendar)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "{\"Code\":422,\"Message\":\"EOF\",\"Data\":null}\n", string(data))
	})

	t.Run("invalid id", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{
			"id": "ID",
			"title": "zxc",
			"startsAt": "2022-08-23T15:04:05+07:00",
			"endsAt": "2022-08-23T15:04:05+07:00"
		}`))
		req := httptest.NewRequest(http.MethodPost, "/events", body)
		w := httptest.NewRecorder()

		logg := logger.New("error", 0)
		calendar := app.New(logg, memorystorage.New())
		handler := NewEventAPIHandler(calendar)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "{\"Code\":422,\"Message\":\"invalid UUID length: 2\",\"Data\":[\"id\"]}\n", string(data))
	})

	t.Run("invalid datetime", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
			"title": "zxc",
			"startsAt": "2022-08-23",
			"endsAt": "2022-08-25"
		}`))
		req := httptest.NewRequest(http.MethodPost, "/events", body)
		w := httptest.NewRecorder()

		logg := logger.New("error", 0)
		calendar := app.New(logg, memorystorage.New())
		handler := NewEventAPIHandler(calendar)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "{\"Code\":422,\"Message\":\"parsing time \\\"2022-08-23\\\" as \\\"2006-01-02T15:04"+
			":05Z07:00\\\": cannot parse \\\"\\\" as \\\"T\\\"\",\"Data\":[\"startsAt\"]}\n", string(data))
	})

	t.Run("date busy", func(t *testing.T) {
		// create valid
		body := bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
			"title": "zxc",
			"startsAt": "2022-08-23T15:04:05+07:00",
			"endsAt": "2022-08-25T15:04:05+07:00"
		}`))
		req := httptest.NewRequest(http.MethodPost, "/events", body)
		w := httptest.NewRecorder()
		logg := logger.New("error", 0)
		calendar := app.New(logg, memorystorage.New())
		handler := NewEventAPIHandler(calendar)
		handler.ServeHTTP(w, req)

		// create another one
		body = bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071c",
			"title": "zxc",
			"startsAt": "2022-08-24T15:04:05+07:00",
			"endsAt": "2022-08-26T15:04:05+07:00"
		}`))
		req = httptest.NewRequest(http.MethodPost, "/events", body)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "{\"Code\":422,\"Message\":\"date is busy\",\"Data\":null}\n", string(data))
	})
}

func TestEventApiHandlerUpdaterErrors(t *testing.T) {
	t.Run("empty request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/events", nil)
		w := httptest.NewRecorder()

		logg := logger.New("error", 0)
		calendar := app.New(logg, memorystorage.New())
		handler := NewEventAPIHandler(calendar)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "{\"Code\":422,\"Message\":\"EOF\",\"Data\":null}\n", string(data))
	})

	t.Run("invalid id", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{
			"id": "ID",
			"title": "zxc",
			"startsAt": "2022-08-23T15:04:05+07:00",
			"endsAt": "2022-08-23T15:04:05+07:00"
		}`))
		req := httptest.NewRequest(http.MethodPut, "/events", body)
		w := httptest.NewRecorder()

		logg := logger.New("error", 0)
		calendar := app.New(logg, memorystorage.New())
		handler := NewEventAPIHandler(calendar)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "{\"Code\":422,\"Message\":\"invalid UUID length: 2\",\"Data\":[\"id\"]}\n", string(data))
	})

	t.Run("invalid datetime", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
			"title": "zxc",
			"startsAt": "2022-08-23",
			"endsAt": "2022-08-25"
		}`))
		req := httptest.NewRequest(http.MethodPut, "/events", body)
		w := httptest.NewRecorder()

		logg := logger.New("error", 0)
		calendar := app.New(logg, memorystorage.New())
		handler := NewEventAPIHandler(calendar)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "{\"Code\":422,\"Message\":\"parsing time \\\"2022-08-23\\\" as \\\"2006-01-02T15:04"+
			":05Z07:00\\\": cannot parse \\\"\\\" as \\\"T\\\"\",\"Data\":[\"startsAt\"]}\n", string(data))
	})

	t.Run("date busy", func(t *testing.T) {
		// create  2 valid
		body := bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
			"title": "zxc",
			"startsAt": "2022-08-20T15:04:05+07:00",
			"endsAt": "2022-08-23T15:04:05+07:00"
		}`))
		req := httptest.NewRequest(http.MethodPost, "/events", body)
		w := httptest.NewRecorder()
		logg := logger.New("error", 0)
		calendar := app.New(logg, memorystorage.New())
		handler := NewEventAPIHandler(calendar)
		handler.ServeHTTP(w, req)

		body = bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071c",
			"title": "zxc",
			"startsAt": "2022-08-24T15:04:05+07:00",
			"endsAt": "2022-08-26T15:04:05+07:00"
		}`))
		req = httptest.NewRequest(http.MethodPost, "/events", body)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		// update first
		body = bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071c",
			"title": "zxc",
			"startsAt": "2022-08-22T15:04:05+07:00",
			"endsAt": "2022-08-26T15:04:05+07:00"
		}`))
		req = httptest.NewRequest(http.MethodPut, "/events", body)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "{\"Code\":422,\"Message\":\"date is busy\",\"Data\":null}\n", string(data))

		// update not found
		body = bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071d",
			"title": "zxc",
			"startsAt": "2022-08-22T15:04:05+07:00",
			"endsAt": "2022-08-26T15:04:05+07:00"
		}`))
		req = httptest.NewRequest(http.MethodPut, "/events", body)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		res = w.Result()
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "{\"Code\":422,\"Message\":\"event not found\",\"Data\":null}\n", string(data))
	})
}

func TestEventApiHandlerDeleteErrors(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		// create
		body := bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
			"title": "zxc",
			"startsAt": "2022-08-23T15:04:05+07:00",
			"endsAt": "2022-08-23T15:04:05+07:00"
		}`))
		req := httptest.NewRequest(http.MethodPost, "/events", body)
		w := httptest.NewRecorder()
		logg := logger.New("error", 0)
		calendar := app.New(logg, memorystorage.New())
		handler := NewEventAPIHandler(calendar)
		handler.ServeHTTP(w, req)

		body = bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071d",
			"title": "zxc",
			"startsAt": "2022-08-22T15:04:05+07:00",
			"endsAt": "2022-08-26T15:04:05+07:00"
		}`))
		req = httptest.NewRequest(http.MethodPut, "/events", body)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "{\"Code\":422,\"Message\":\"event not found\",\"Data\":null}\n", string(data))
	})
}

func TestEventApiHandlerSuccess(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// create
		body := bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
			"title": "zxc",
			"startsAt": "2022-08-23T15:04:05+07:00",
			"endsAt": "2022-08-23T15:04:05+07:00",
			"notify": "2022-08-23T13:04:05+07:00"
		}`))
		req := httptest.NewRequest(http.MethodPost, "/events", body)
		w := httptest.NewRecorder()
		logg := logger.New("error", 0)
		calendar := app.New(logg, memorystorage.New())
		handler := NewEventAPIHandler(calendar)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, w.Code)
		require.Empty(t, data)

		// get
		req = httptest.NewRequest(http.MethodGet, "/events?day=2022-08-23&period=month", nil)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		res = w.Result()
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.NotEmpty(t, data)
		require.Equal(
			t,
			"[{\"id\":\"48cd8858-9103-4c6a-9a83-1d58307f071b\",\"title\":\"zxc\",\"startsAt\":"+
				"\"2022-08-23T15:04:05+07:00\",\"endsAt\":\"2022-08-23T15:04:05+07:00\",\"notify\":"+
				"\"2022-08-23T13:04:05+07:00\"}]\n",
			string(data),
		)

		// update
		body = bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
			"title": "NEW TITLE!",
			"startsAt": "2022-08-25T15:04:05+07:00",
			"endsAt": "2022-08-27T15:04:05+07:00"
		}`))
		req = httptest.NewRequest(http.MethodPut, "/events", body)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		res = w.Result()
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		require.Nil(t, err)
		fmt.Printf("%#v", string(data))
		require.Equal(t, http.StatusOK, w.Code)
		require.Empty(t, data)

		// get
		req = httptest.NewRequest(http.MethodGet, "/events?day=2022-08-23&period=month", nil)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		res = w.Result()
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.NotEmpty(t, data)
		require.Equal(
			t,
			"[{\"id\":\"48cd8858-9103-4c6a-9a83-1d58307f071b\",\"title\":\"NEW TITLE!\",\"startsAt\":"+
				"\"2022-08-25T15:04:05+07:00\",\"endsAt\":\"2022-08-27T15:04:05+07:00\"}]\n",
			string(data),
		)

		// add another one
		body = bytes.NewReader([]byte(`{
			"id": "12cd8858-9103-4c6a-9a83-1d58307f071c",
			"title": "another one",
			"startsAt": "2022-08-28T15:04:05+07:00",
			"endsAt": "2022-08-29T15:04:05+07:00"
		}`))

		req = httptest.NewRequest(http.MethodPost, "/events", body)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		res = w.Result()
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, w.Code)
		require.Empty(t, data)

		// get
		req = httptest.NewRequest(http.MethodGet, "/events?day=2022-08-23&period=month", nil)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		res = w.Result()
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.NotEmpty(t, data)
		require.Equal(
			t,
			"[{\"id\":\"48cd8858-9103-4c6a-9a83-1d58307f071b\",\"title\":\"NEW TITLE!\",\"startsAt\":\""+
				"2022-08-25T15:04:05+07:00\",\"endsAt\":\"2022-08-27T15:04:05+07:00\"},{\"id\":\"12cd8858-9103"+
				"-4c6a-9a83-1d58307f071c\",\"title\":\"another one\",\"startsAt\":\"2022-08-28T15:04:05+07:00"+
				"\",\"endsAt\":\"2022-08-29T15:04:05+07:00\"}]\n",
			string(data),
		)

		// delete first
		req = httptest.NewRequest(http.MethodDelete, "/events?id=48cd8858-9103-4c6a-9a83-1d58307f071b", nil)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		res = w.Result()
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusNoContent, w.Code)
		require.Empty(t, data)

		// get
		req = httptest.NewRequest(http.MethodGet, "/events?day=2022-08-23&period=month", nil)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		res = w.Result()
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.NotEmpty(t, data)
		require.Equal(
			t,
			"[{\"id\":\"12cd8858-9103-4c6a-9a83-1d58307f071c\",\"title\":\"another one\","+
				"\"startsAt\":\"2022-08-28T15:04:05+07:00\",\"endsAt\":\"2022-08-29T15:04:05+07:00\"}]\n",
			string(data),
		)
	})
}
