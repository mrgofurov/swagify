package swagify

import (
	"encoding/json"
	"net/http"
	"reflect"
)

var (
	ctxPtrType = reflect.TypeOf((*Ctx)(nil))
	errorType  = reflect.TypeOf((*error)(nil)).Elem()
)

// scanResult holds the http.HandlerFunc produced from a typed handler
// together with the schema types extracted from its signature.
type scanResult struct {
	handler   http.HandlerFunc
	reqType   reflect.Type // body (POST/PUT/PATCH) or nil
	resType   reflect.Type
	queryType reflect.Type // query params (GET/DELETE) or nil
}

// scanHandler inspects h and returns a scanResult.
//
// Supported handler signatures:
//
//	func(http.ResponseWriter, *http.Request)   — plain net/http, no schema inference
//	func(*Ctx) error                           — access path/query/headers via ctx
//	func(*Ctx) (Res, error)                    — response schema inferred
//	func(*Ctx, Req) (Res, error)               — request + response schemas inferred
//
// For GET and DELETE, the Req type is treated as URL query parameters.
// For POST, PUT, PATCH, it is treated as a JSON request body.
func scanHandler(method string, h any) scanResult {
	// Fast path: plain net/http handlers
	if fn, ok := h.(http.HandlerFunc); ok {
		return scanResult{handler: fn}
	}
	if fn, ok := h.(func(http.ResponseWriter, *http.Request)); ok {
		return scanResult{handler: fn}
	}

	v := reflect.ValueOf(h)
	t := v.Type()

	if t.Kind() != reflect.Func {
		panic("swagify: handler must be a function")
	}
	if t.NumIn() == 0 || t.In(0) != ctxPtrType {
		panic("swagify: handler first parameter must be *swagify.Ctx")
	}
	if t.NumOut() == 0 || !t.Out(t.NumOut()-1).Implements(errorType) {
		panic("swagify: handler must return error as its last value")
	}

	var reqType, resType, queryType reflect.Type

	if t.NumIn() == 2 {
		rt := t.In(1)
		if isQueryMethod(method) {
			queryType = rt
		} else {
			reqType = rt
		}
	}
	if t.NumOut() == 2 {
		resType = t.Out(0)
	}

	return scanResult{
		handler:   buildWrapper(v, t, method),
		reqType:   reqType,
		resType:   resType,
		queryType: queryType,
	}
}

// buildWrapper wraps a typed handler function into an http.HandlerFunc.
func buildWrapper(fn reflect.Value, t reflect.Type, method string) http.HandlerFunc {
	hasReq := t.NumIn() == 2
	hasRes := t.NumOut() == 2
	useQuery := isQueryMethod(method)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Ctx{
			Context:  r.Context(),
			Request:  r,
			Response: w,
		}
		args := []reflect.Value{reflect.ValueOf(ctx)}

		if hasReq {
			reqType := t.In(1)
			isPtr := reqType.Kind() == reflect.Ptr
			elemType := reqType
			if isPtr {
				elemType = reqType.Elem()
			}

			req := reflect.New(elemType) // always allocate *T for parsing

			var parseErr error
			if useQuery {
				parseErr = parseQueryString(r, req.Interface())
			} else {
				parseErr = json.NewDecoder(r.Body).Decode(req.Interface())
			}
			if parseErr != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{
					"error": parseErr.Error(),
				})
				return
			}

			if isPtr {
				args = append(args, req)
			} else {
				args = append(args, req.Elem())
			}
		}

		results := fn.Call(args)
		errVal := results[len(results)-1]
		if !errVal.IsNil() {
			err := errVal.Interface().(error)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		if hasRes {
			writeJSON(w, successStatus(method), results[0].Interface())
		} else if method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

// isQueryMethod reports whether an HTTP method carries its parameters
// in the URL rather than a request body.
func isQueryMethod(method string) bool {
	return method == "GET" || method == "DELETE" || method == "HEAD"
}

// successStatus returns the default success HTTP status for a method.
func successStatus(method string) int {
	switch method {
	case "POST":
		return http.StatusCreated
	case "DELETE":
		return http.StatusNoContent
	default:
		return http.StatusOK
	}
}
