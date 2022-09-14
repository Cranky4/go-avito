package integrationtests_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	StartsAt    string `json:"startsAt"`
	EndsAt      string `json:"endsAt"`
	NotifyAfter string `json:"notify,omitempty"`
}

var events = []string{
	`{
		"id": "28cd8858-9103-4c6a-9a83-1d58307f071a",
		"title": "first event",
		"startsAt": "2022-06-01T15:00:00+07:00",
		"endsAt": 	"2022-06-01T17:00:00+07:00",
		"notify": 	"2022-06-01T17:00:00+07:00"
	}`,
	`{
		"id": "28cd8858-9103-4c6a-9a83-1d58307f071b",
		"title": "second event",
		"startsAt": "2022-06-01T17:00:01+07:00",
		"endsAt": 	"2022-06-01T19:00:00+07:00",
		"notify": 	"2022-06-01T17:00:01+07:00"
	}`,
	`{
		"id": "28cd8858-9103-4c6a-9a83-1d58307f071c",
		"title": "third event",
		"startsAt": "2022-06-03T19:00:01+07:00",
		"endsAt": 	"2022-06-03T21:00:00+07:00",
		"notify": 	"2022-06-03T21:00:01+07:00"
	}`,
	`{
		"id": "28cd8858-9103-4c6a-9a83-1d58307f071d",
		"title": "fourth event",
		"startsAt": "2022-06-17T19:00:00+07:00",
		"endsAt": 	"2022-06-17T21:00:00+07:00",
		"notify": 	"2022-06-17T19:00:00+07:00"
	}`,
}

var _ = Describe("Get events via HTTP", Ordered, func() {
	BeforeAll(func() {
		for _, ev := range events {
			resp, err := http.Post(
				"http://localhost:8888/events",
				"application/json",
				bytes.NewReader([]byte(ev)),
			)
			if err != nil {
				Fail("error while do http request" + err.Error())
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				Fail("error while do http request" + err.Error())
			}

			if resp.StatusCode != http.StatusCreated {
				Fail("request is " + ev + "; error response" + string(body))
			}
		}
	})

	It("Get day events", func() {
		resp, err := http.Get("http://localhost:8888/events?day=2022-06-01&period=day")
		if err != nil {
			Fail("error while do http request" + err.Error())
		}
		defer resp.Body.Close()

		Expect(http.StatusOK).To(Equal(resp.StatusCode))

		body, err := io.ReadAll(resp.Body)

		Expect(err).To(BeNil())
		Expect("[{\"id\":\"28cd8858-9103-4c6a-9a83-1d58307f071a\",\"title\":\"first event\"," +
			"\"startsAt\":\"2022-06-01T15:00:00Z\",\"endsAt\":\"2022-06-01T17:00:00Z\",\"notify\":" +
			"\"2022-06-01T17:00:01Z\"},{\"id\":\"28cd8858-9103-4c6a-9a83-1d58307f071b\",\"title\":" +
			"\"second event\",\"startsAt\":\"2022-06-01T17:00:01Z\",\"endsAt\":\"2022-06-01T19:00:00Z\"," +
			"\"notify\":\"2022-06-01T17:00:01Z\"}]\n").To(Equal(string(body)))
	})

	It("Get week events", func() {
		resp, err := http.Get("http://localhost:8888/events?day=2022-06-01&period=week")
		if err != nil {
			Fail("error while do http request" + err.Error())
		}
		defer resp.Body.Close()

		Expect(http.StatusOK).To(Equal(resp.StatusCode))

		body, err := io.ReadAll(resp.Body)

		Expect(err).To(BeNil())
		Expect("[{\"id\":\"28cd8858-9103-4c6a-9a83-1d58307f071a\",\"title\":\"first event\",\"startsAt\":"+
			"\"2022-06-01T15:00:00Z\",\"endsAt\":\"2022-06-01T17:00:00Z\",\"notify\":\"2022-06-03T21:00:01Z\"}"+
			",{\"id\":\"28cd8858-9103-4c6a-9a83-1d58307f071b\",\"title\":\"second event\",\"startsAt\":"+
			"\"2022-06-01T17:00:01Z\",\"endsAt\":\"2022-06-01T19:00:00Z\",\"notify\":\"2022-06-03T21:00:01Z\"},"+
			"{\"id\":\"28cd8858-9103-4c6a-9a83-1d58307f071c\",\"title\":\"third event\",\"startsAt\""+
			":\"2022-06-03T19:00:01Z\",\"endsAt\":\"2022-06-03T21:00:00Z\",\"notify\":\"2022-06-03T21:00:01Z\"}"+
			"]\n").To(Equal(string(body)), string(body))
	})

	It("Get month events", func() {
		resp, err := http.Get("http://localhost:8888/events?day=2022-06-01&period=month")
		if err != nil {
			Fail("error while do http request" + err.Error())
		}
		defer resp.Body.Close()

		Expect(http.StatusOK).To(Equal(resp.StatusCode))

		body, err := io.ReadAll(resp.Body)

		Expect(err).To(BeNil())
		Expect("[{\"id\":\"28cd8858-9103-4c6a-9a83-1d58307f071a\",\"title\":\"first event\",\"startsAt\":"+
			"\"2022-06-01T15:00:00Z\",\"endsAt\":\"2022-06-01T17:00:00Z\",\"notify\":\"2022-06-17T19:00:00Z\"},"+
			"{\"id\":\"28cd8858-9103-4c6a-9a83-1d58307f071b\",\"title\":\"second event\",\"startsAt\":"+
			"\"2022-06-01T17:00:01Z\",\"endsAt\":\"2022-06-01T19:00:00Z\",\"notify\":\"2022-06-17T19:00:00Z\"},"+
			"{\"id\":\"28cd8858-9103-4c6a-9a83-1d58307f071c\",\"title\":\"third event\",\"startsAt\":"+
			"\"2022-06-03T19:00:01Z\",\"endsAt\":\"2022-06-03T21:00:00Z\",\"notify\":\"2022-06-17T19:00:00Z\"},"+
			"{\"id\":\"28cd8858-9103-4c6a-9a83-1d58307f071d\",\"title\":\"fourth event\",\"startsAt\":"+
			"\"2022-06-17T19:00:00Z\",\"endsAt\":\"2022-06-17T21:00:00Z\",\"notify\":\"2022-06-17T19:00:00Z\"}]"+
			"\n").To(Equal(string(body)), string(body))
	})

	AfterAll(func() {
		var event Event

		for _, ev := range events {
			err := json.Unmarshal([]byte(ev), &event)
			if err != nil {
				Fail("error while json unmarshal" + err.Error())
			}

			req, err := http.NewRequest(
				http.MethodDelete,
				"http://localhost:8888/events?id="+event.ID,
				nil,
			)
			if err != nil {
				Fail("error while building request" + err.Error())
			}
			req.Header.Add("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				Fail("error while do http request" + err.Error())
				return
			}
			defer resp.Body.Close()
		}
	})
})
