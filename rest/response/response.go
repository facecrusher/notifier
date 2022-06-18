package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Response interface {
	Bytes()
	Decode(interface{})
}

type ReqResponse struct {
	*http.Response
	Err      error
	ByteBody []byte
}

// Bytes return the Response Body as bytes.
func (r *ReqResponse) Bytes() []byte {
	return r.ByteBody
}

func (r *ReqResponse) Decode(decode interface{}) error {
	ctypeJSON := "application/json"

	ctype := strings.ToLower(r.Header.Get("Content-Type"))

	for i := 0; i < 2; i++ {

		switch {
		case strings.Contains(ctype, ctypeJSON):
			return json.Unmarshal(r.ByteBody, decode)
		case i == 0:
			ctype = http.DetectContentType(r.ByteBody)
		}

	}
	return errors.New("Response format neither JSON nor XML")
}
