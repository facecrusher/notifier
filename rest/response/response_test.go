package response

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Decode(t *testing.T) {
	var decode map[string]string
	mockHeaders := make(http.Header)
	mockHeaders.Add("Content-Type", "application/json")

	response := ReqResponse{
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Header:     mockHeaders,
		},
		ByteBody: []byte(`{}`),
	}

	err := response.Decode(&decode)
	assert.Nil(t, err)

	mockHeaders.Set("Content-Type", "unknown")
	assert.Error(t, response.Decode(response))
}
