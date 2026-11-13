// Package forward: pola handler Portal yang meneruskan permintaan ke BE User dan
// mengembalikan jawabannya apa adanya (bentuk respons sama dengan sebelum BE User ada).
package forward

import (
	"portal_api/errors"
	request_helper "portal_api/helper/request"
	response_helper "portal_api/helper/response"
	"portal_api/pkg/authclient"

	"github.com/gin-gonic/gin"
)

// Call: fungsi service yang meneruskan body ke BE User.
type Call func(body []byte) (*authclient.Response, error)

// CallWithID: seperti Call, dengan id angka dari parameter URL.
type CallWithID func(id int, body []byte) (*authclient.Response, error)

// Handle membaca body apa adanya, memanggil service, lalu menulis jawaban BE User.
// Penolakan BE User dan gangguan layanan diteruskan lewat c.Error (ErrorHandlerMiddleware).
func Handle(c *gin.Context, call Call) {
	body, err := request_helper.ReadRawBody(c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: "Body permintaan tidak dapat dibaca"})
		return
	}
	resp, err := call(body)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WriteRaw(c, resp.StatusCode, resp.Body)
}

// HandleWithID: seperti Handle, tetapi parameter URL param wajib angka positif (mencegah
// path lain ikut terbawa ke URL BE User).
func HandleWithID(c *gin.Context, param string, call CallWithID) {
	id, err := request_helper.ConvertParamStringToInt(c, param)
	if err != nil || id <= 0 {
		c.Error(&errors.BadRequestError{Message: "ID tidak valid"})
		return
	}
	Handle(c, func(body []byte) (*authclient.Response, error) { return call(id, body) })
}
