package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PentestGPT-Project/htb-cli/config"
	"go.uber.org/zap"
)

func prepareRequestTest(t *testing.T) {
	t.Helper()
	t.Setenv("HTB_TOKEN", "header.payload.signature")
	config.GlobalConfig.Logger = zap.NewNop()
	config.GlobalConfig.ProxyParam = ""
}

func TestHtbRequestRejectsHTTPFailure(t *testing.T) {
	prepareRequestTest(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = io.WriteString(w, `{"message":"upstream failed"}`)
	}))
	defer server.Close()

	if _, err := HtbRequest(http.MethodGet, server.URL, nil); err == nil {
		t.Fatal("HtbRequest accepted an HTTP 502 response")
	}
}

func TestHtbRequestRejectsHTMLInsteadOfAPIJSON(t *testing.T) {
	prepareRequestTest(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, "<html>proxy error</html>")
	}))
	defer server.Close()

	if _, err := HtbRequest(http.MethodGet, server.URL, nil); err == nil {
		t.Fatal("HtbRequest accepted HTML where API JSON was expected")
	}
}

func TestHtbRequestVerifiesTLSCertificates(t *testing.T) {
	prepareRequestTest(t)

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	if _, err := HtbRequest(http.MethodGet, server.URL, nil); err == nil {
		t.Fatal("HtbRequest accepted an untrusted TLS certificate")
	}
}

func TestHtbRequestReturnsSuccessfulJSONBody(t *testing.T) {
	prepareRequestTest(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer server.Close()

	resp, err := HtbRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("HtbRequest returned an error: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response: %v", err)
	}
	if strings.TrimSpace(string(body)) != `{"ok":true}` {
		t.Fatalf("body = %q", body)
	}
}
