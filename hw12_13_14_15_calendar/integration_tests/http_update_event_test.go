package integrationtests_test

import (
	"bytes"
	"io"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update event via HTTP", func() {
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

	Context("with empty request", func() {
		req, err := http.NewRequest(http.MethodPut, baseURL+"/events", nil)
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

		It("status code is 400", func() {
			Expect(400).To(Equal(resp.StatusCode))
		})

		body, err := io.ReadAll(resp.Body)

		It("no erros appears while body reading", func() {
			Expect(err).To(BeNil())
		})
		It("error body", func() {
			Expect("{\"Code\":422,\"Message\":\"EOF\",\"Data\":null}\n").To(Equal(string(body)))
		})
	})

	Context("with invalid id", func() {
		req, err := http.NewRequest(http.MethodPut, baseURL+"/events", bytes.NewReader([]byte(`{
			"id": "ID",
			"title": "first event",
			"startsAt": "2022-08-23T15:04:05+07:00",
			"notify": "2022-08-23T15:04:05+07:00",
			"endsAt": "2022-08-23T15:04:05+07:00"
		}`)))
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

		body, err := io.ReadAll(resp.Body)

		It("no erros appears while body reading", func() {
			Expect(err).To(BeNil())
		})
		It("error body", func() {
			Expect("{\"Code\":422,\"Message\":\"invalid UUID length: 2\",\"Data\":[\"id\"]}\n").To(Equal(string(body)))
		})
	})

	Context("with invalid dates", func() {
		req, err := http.NewRequest(http.MethodPut, baseURL+"/events", bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
			"title": "first event",
			"startsAt": "2022-08-23",
			"endsAt": "2022-08-25"
		}`)))
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

		body, err := io.ReadAll(resp.Body)

		It("no erros appears while body reading", func() {
			Expect(err).To(BeNil())
		})
		It("error body", func() {
			Expect("{\"Code\":422,\"Message\":\"parsing time \\\"2022-08-23\\\" as \\\"2006-01-02T15:04" +
				":05Z07:00\\\": cannot parse \\\"\\\" as \\\"T\\\"\",\"Data\":[\"startsAt\"]}\n").To(Equal(string(body)))
		})
	})

	Context("with busy dates", func() {
		Describe("create second event", func() {
			resp, err := http.Post(
				baseURL+"/events",
				"application/json",
				bytes.NewReader([]byte(`{
						"id": "48cd8858-9103-4c6a-9a83-1d58307f071c",
						"title": "second event",
						"startsAt": "2022-10-23T15:04:05+07:00",
						"endsAt": "2022-10-23T15:04:05+07:00",
						"notify": "2022-10-23T15:04:05+07:00"
					}`)),
			)
			if err != nil {
				Fail("error while do http request" + err.Error())
			}
			defer resp.Body.Close()
		})

		Describe("try to update fist event with date of second event", func() {
			req, err := http.NewRequest(http.MethodPut, baseURL+"/events", bytes.NewReader([]byte(`{
				"id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
				"title": "first event",
				"startsAt": "2022-10-23T14:04:05+07:00",
				"endsAt": "2022-10-23T16:04:05+07:00",
				"notify": "2022-10-23T15:04:05+07:00"
			}`)))
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
				Expect("{\"Code\":422,\"Message\":\"date is busy\",\"Data\":null}\n").To(Equal(string(body)))
			})
		})

		Describe("delete second event", func() {
			req, err := http.NewRequest(
				http.MethodDelete,
				baseURL+"/events?id=48cd8858-9103-4c6a-9a83-1d58307f071c",
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

	Context("event not found", func() {
		req, err := http.NewRequest(http.MethodPut, baseURL+"/events", bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071z",
			"title": "first event new name",
			"startsAt": "2022-10-23T14:04:05+07:00",
			"endsAt": "2022-10-23T16:04:05+07:00",
			"notify": "2022-10-23T15:04:05+07:00"
		}`)))
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

	Context("with valid parameters", func() {
		req, err := http.NewRequest(http.MethodPut, baseURL+"/events", bytes.NewReader([]byte(`{
			"id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
			"title": "first event new name",
			"startsAt": "2022-10-23T14:04:05+07:00",
			"endsAt": "2022-10-23T16:04:05+07:00",
			"notify": "2022-10-23T15:04:05+07:00"
		}`)))
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

		It("is status code 200", func() {
			Expect(http.StatusOK).To(Equal(resp.StatusCode))
		})

		It("is error body", func() {
			Expect("").To(Equal(string(body)))
		})
	})
})
