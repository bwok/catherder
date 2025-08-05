package main

import (
	"net/http/httptest"
	"testing"
)

func TestPageNewHandler(t *testing.T) {
	request := httptest.NewRequest("GET", "https://localhost/edit", nil)
	w := httptest.NewRecorder()
	pageEditHandler(w, request)

	response := w.Result()

	if response.StatusCode != 200 {
		t.Error("http status code was not 200")
	}
	if response.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Error("Content-Type header was wrong")
	}
	if response.Header.Get("Strict-Transport-Security") != "max-age=31536000; includeSubDomains" {
		t.Error("Strict-Transport-Security header was wrong")
	}
	if response.Header.Get("Content-Security-Policy") != "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';" {
		t.Error("Content-Security-Policy header was wrong")
	}

}

func TestAjaxCreateHandler(t *testing.T) {
	// TODO test
}

func TestPageViewHandler(t *testing.T) {
	// TODO test
}

func TestAddUserHandler(t *testing.T) {
	// TODO test
}

func TestPageAdminHandler(t *testing.T) {
	// TODO test
}

func TestAjaxAdminSaveHandler(t *testing.T) {
	// TODO test
}

func TestAjaxAdminDeleteHandler(t *testing.T) {
	// TODO test
}
