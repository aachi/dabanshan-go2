package transport

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	// p_endpoint "github.com/laidingqing/dabanshan-go/svcs/product/endpoint"
	"github.com/laidingqing/dabanshan-go/svcs/product/model"
	"github.com/laidingqing/dabanshan-go/svcs/product/service"
)

var (
	// ErrUploadPartParams ...
	ErrUploadPartParams = errors.New("file part error.")
)

func decodeHTTPCreateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	a := model.CreateProductRequest{}
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func decodeHTTPGetProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	// err := json.NewDecoder(r.Body).Decode(&req)
	// todo convert params..
	a, _ := strconv.ParseInt(r.FormValue("userid"), 10, 64)
	b, _ := strconv.ParseInt(r.FormValue("size"), 10, 64)
	return model.GetProductsRequest{A: a, B: b}, nil
}

func decodeHTTPUploadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	file, handle, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		return ErrUploadPartParams, nil
	}
	buff := make([]byte, 512)
	if _, err = file.Read(buff); err != nil {
		return ErrUploadPartParams, nil
	}
	vmd5 := fmt.Sprintf("%x", md5.Sum(buff))

	return model.UploadProductRequest{
		Name: handle.Filename,
		Md5:  vmd5,
		Body: buff,
	}, nil
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

// encodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// encodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(model.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorWrapper struct {
	Error string `json:"error"`
}

func err2code(err error) int {
	switch err {
	case service.ErrTwoZeroes, service.ErrMaxSizeExceeded, service.ErrIntOverflow:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
