package api_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gin-gonic/gin"
	"github.com/kornypoet/lakitu/api"
)

var _ = Describe("Router", func() {

	var router *gin.Engine

	BeforeEach(func() {
		router = api.Router(false)
		api.AssetDir = GinkgoT().TempDir()
		api.Version = "test"
	})

	When("GET /v1/version", func() {
		It("returns 200", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/v1/version", nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(200))
			Expect(w.Body.String()).To(Equal("test"))
		})
	})

	When("GET /v1/manage_file", func() {
		It("returns 405", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/v1/manage_file", nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(405))
			Expect(w.Body.String()).To(Equal("405 method not allowed"))
		})
	})

	When("POST /v1/manage_file", func() {
		It("returns 400", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/v1/manage_file", nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(400))
			res := `{"status":"failure","err":"invalid request"}`
			Expect(w.Body.String()).To(MatchJSON(res))
		})
	})

	When(`POST /v1/manage_file {"invalid":"json"}`, func() {
		It("returns 400", func() {
			w := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(`{"invalid":"json"}`))
			req, _ := http.NewRequest("POST", "/v1/manage_file", body)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(400))
			res := `{"status":"failure","err":"Key: 'Payload.Action' Error:Field validation for 'Action' failed on the 'required' tag"}`
			Expect(w.Body.String()).To(MatchJSON(res))
		})
	})

	When(`POST /v1/manage_file {"action":"invalid"}`, func() {
		It("returns 400", func() {
			w := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(`{"action":"invalid"}`))
			req, _ := http.NewRequest("POST", "/v1/manage_file", body)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(400))
			res := `{"status":"failure","err":"Key: 'Payload.Action' Error:Field validation for 'Action' failed on the 'oneof' tag"}`
			Expect(w.Body.String()).To(MatchJSON(res))
		})
	})

	When(`POST /v1/manage_file {"action":"download"}`, func() {
		It("returns 200", func() {
			w := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(`{"action":"download"}`))
			req, _ := http.NewRequest("POST", "/v1/manage_file", body)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(200))
			res := `{"status":"success","action":"download"}`
			Expect(w.Body.String()).To(MatchJSON(res))
		})
	})

	When(`POST /v1/manage_file {"action":"download"}`, func() {
		It("returns 500", func() {
			_, _ = os.Create(api.AssetFile())
			w := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(`{"action":"download"}`))
			req, _ := http.NewRequest("POST", "/v1/manage_file", body)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(500))
			res := `{"status":"failure","err":"file already downloaded"}`
			Expect(w.Body.String()).To(MatchJSON(res))
		})
	})

	When(`POST /v1/manage_file {"action":"download"}`, func() {
		It("returns 500", func() {
			w := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(`{"action":"download"}`))
			req, _ := http.NewRequest("POST", "/v1/manage_file", body)
			api.BlockDownload <- true
			router.ServeHTTP(w, req)
			<-api.BlockDownload

			Expect(w.Code).To(Equal(429))
			res := `{"status":"failure","err":"file download in progress"}`
			Expect(w.Body.String()).To(MatchJSON(res))
		})
	})

	When(`POST /v1/manage_file {"action":"read"}`, func() {
		It("returns 500", func() {
			w := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(`{"action":"read"}`))
			req, _ := http.NewRequest("POST", "/v1/manage_file", body)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(500))
			res := `{"status":"failure","err":"file must be downloaded first"}`
			Expect(w.Body.String()).To(MatchJSON(res))
		})
	})

	When(`POST /v1/manage_file {"action":"read"}`, func() {
		It("returns 200", func() {
			contents := "lorum ipsum"
			out, _ := os.Create(api.AssetFile())
			defer out.Close()
			_, _ = out.WriteString(contents)
			w := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(`{"action":"read"}`))
			req, _ := http.NewRequest("POST", "/v1/manage_file", body)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(200))
			Expect(w.Body.String()).To(Equal(contents))
		})
	})
})
