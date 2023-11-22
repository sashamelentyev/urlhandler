package urlhandler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sashamelentyev/urlhandler/internal/urlhandler"
)

func TestURLHandler_CheckURL(t *testing.T) {
	noFile := "" // cuz test CheckURL only.
	reqTimeout := 5 * time.Second

	urlHandler := urlhandler.New(noFile, reqTimeout)

	ctx := context.Background()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test url handler")
	})
	s := httptest.NewServer(h)

	resp, err := urlHandler.CheckURL(ctx, s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.URL != s.URL {
		t.Fatalf("response url not equal %q", s.URL)
	}
	if resp.ContentLength != 17 {
		t.Fatal("content length not equal 17")
	}
	if resp.ResponseTime > reqTimeout {
		t.Fatal("response time greater that request timeout")
	}
}
