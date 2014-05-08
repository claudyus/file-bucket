package main

import (
  "log"
  "testing"
  "net/http"
  "github.com/stretchr/testify/assert"
)


func TestSomething(t *testing.T) {
    res, err := http.NewRequest("POST", "http://localhost:1234/failtoken", nil)
    if assert.NotNil(t, err) {
        log.Printf("ds %d", err)
        assert.Equal(t, res.StatusCode, 403)
    }

}
