package integrationtests_test

import (
	"bytes"
	"io"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delete event via HTTP", func() {
	baseURL := os.Getenv("CALENDAR_API_BASE_URL")

	BeforeEach(func() {
		resp, err := http.Post(
			baseURL+"/events",
			"application/json",
			bytes.NewReader([]byte(`{
					"id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
					"title": "first event",
					"startsAt": "2022-09-23T15:04:05+07:00",
					"endsAt": "2022-09-23T15:04:05+07:00",
					"notify": "2022-09-23T15:04:05+07:00"
				}`)),
		)
		if err != nil {
			Fail("error while do http request" + err.Error())
		}
		defer resp.Body.Close()
	})

	AfterEach(func() {
		req, err := http.NewRequest(
			http.MethodDelete,
			baseURL+"/events?id=48cd8858-9103-4c6a-9a83-1d58307f071b",
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

	Context("event not found", func() {
		req, err := http.NewRequest(
			http.MethodDelete,
			baseURL+"/events?id=48cd8858-9103-4c6a-9a83-1d58307f071z",
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

		It("no errors appears while body reading", func() {
			Expect(err).To(BeNil())
		})

		body, err := io.ReadAll(resp.Body)

		It("is status code 422", func() {
			Expect(http.StatusUnprocessableEntity).To(Equal(resp.StatusCode))
		})

		It("is error body", func() {
			Expect("{\"Code\":422,\"Message\":\"invalid UUID format\",\"Data\":[\"id\"]}\n").To(Equal(string(body)))
		})
	})

	Context("event deleted", func() {
		req, err := http.NewRequest(
			http.MethodDelete,
			baseURL+"/events?id=48cd8858-9103-4c6a-9a83-1d58307f071b",
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

		It("no errors appears while body reading", func() {
			Expect(err).To(BeNil())
		})

		body, err := io.ReadAll(resp.Body)

		It("is status code 204", func() {
			Expect(http.StatusNoContent).To(Equal(resp.StatusCode))
		})

		It("is error body", func() {
			Expect("").To(Equal(string(body)))
		})
	})
})
