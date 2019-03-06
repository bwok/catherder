package main

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPageNewHandler(t *testing.T) {
	request := httptest.NewRequest("GET", "https://localhost/new.html", nil)
	w := httptest.NewRecorder()
	pageEditHandler(w, request)

	response := w.Result()
	body, _ := ioutil.ReadAll(response.Body)

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

	file, err := os.Open("templates/edit.html")
	if err != nil {
		t.Fatal(err)
	}
	htmlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(body, htmlBytes) != 0 {
		t.Error("returned bytes for new.html differed to the one found on the filesystem.")
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
