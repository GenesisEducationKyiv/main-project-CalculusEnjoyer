package ctrl

import "net/http"

type ErrorTransformer interface {
	TransformToHTTPErr(err error, w http.ResponseWriter)
}
