package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	apperrors "music/internal/common/errors"
)

const maxJSONBodySize = 1 << 20 // 1 MB

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return apperrors.BadRequest("invalid JSON syntax", map[string]interface{}{
				"offset": syntaxError.Offset,
			})

		case errors.Is(err, io.ErrUnexpectedEOF):
			return apperrors.BadRequest("invalid JSON", "unexpected end of JSON")

		case errors.As(err, &unmarshalTypeError):
			return apperrors.BadRequest("invalid JSON field type", map[string]interface{}{
				"field":  unmarshalTypeError.Field,
				"type":   unmarshalTypeError.Type.String(),
				"offset": unmarshalTypeError.Offset,
			})

		case errors.Is(err, io.EOF):
			return apperrors.BadRequest("request body is required", nil)

		default:
			return apperrors.BadRequest("invalid request body", err.Error())
		}
	}

	if decoder.Decode(&struct{}{}) != io.EOF {
		return apperrors.BadRequest("request body must contain only one JSON object", nil)
	}

	return nil
}

func DecodePathID(value string, name string) error {
	if value == "" {
		return apperrors.BadRequest(fmt.Sprintf("%s is required", name), nil)
	}

	return nil
}
