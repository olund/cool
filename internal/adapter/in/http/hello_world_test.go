package http

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
)

func TestHelloWorld(t *testing.T) {
	apitest.New().
		HandlerFunc(GetHelloWorld()).
		Get("/").
		Expect(t).
		Body(`Hello World`).
		Status(http.StatusOK).
		End()
}
