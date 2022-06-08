package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Post(t *testing.T) {

	type mockResponse struct {
		About   string `json:"about"`
		EventID string `json:"event_id"`
	}

	url := "https://eouss1txxyn5t7x.m.pipedream.net"
	headers := make(map[string]string)
	client := NewNotifierRestClient(url, headers)

	body := "`{\"test\":\"event\"}`"
	decode := &mockResponse{}

	response := client.Post(body, decode)

	assert.Nil(t, response)
	assert.NotNil(t, decode)
	assert.NotNil(t, decode.About)
	assert.NotNil(t, decode.EventID)
}
