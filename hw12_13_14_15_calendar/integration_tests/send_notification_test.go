package integrationtests_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get events via HTTP", Ordered, func() {
	baseURL := os.Getenv("CALENDAR_API_BASE_URL")
	BeforeAll(func() {
		resp, err := http.Post(
			baseURL+"/events",
			"application/json",
			bytes.NewReader([]byte(`{
				"id": "38cd8858-9103-4c6a-9a83-1d58307f071a",
				"title": "first event",
				"startsAt": "2022-06-01T15:00:00+07:00",
				"endsAt": 	"2022-06-01T17:00:00+07:00",
				"notify": 	"2022-06-01T17:00:00+07:00"
			}`)),
		)
		if err != nil {
			Fail("error while do http request" + err.Error())
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			Fail("request failed")
		}
	})

	It("check email log", func() {
		tries := 20
		timout := 1 * time.Second

		try := 1

		for {
			f, err := os.Open("/logs/email.log")
			if try > tries {
				f.Close()
				Fail("Maximum tries to read email log is reached")
			}

			Expect(err).To(BeNil())

			r := bufio.NewReader(f)
			line, err := r.ReadString('\n')

			if errors.Is(err, io.EOF) {
				f.Close()
				time.Sleep(timout)
				try++
				continue
			}

			Expect(err).To(BeNil())

			Expect(line).To(Equal("[NOTIFICATION SENT] {\"ID\":\"38cd8858-9103-4c6a-9a83-1d58307f071a\",\"Title\"" +
				":\"first event\",\"StartedAt\":\"2022-06-01T15:00:00Z\",\"UserID\":\"\"} \n"))

			f.Close()
			break
		}
	})

	AfterAll(func() {
		req, err := http.NewRequest(
			http.MethodDelete,
			baseURL+"/events?id=38cd8858-9103-4c6a-9a83-1d58307f071a",
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
	})
})
