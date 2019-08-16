package controllers

import (
	"fmt"
	"net/http"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// QRHandler creates and serves QR Codes for notes
func QRHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/temp/qr/notes/")
	if !isUUID(code) {
		return
	}
	code = fmt.Sprintf("%sdropcode?voucher=%s", site, code)
	qr, err := qrcode.Encode(code, qrcode.Medium, 256)
	Handle(err)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/png")
	w.Write(qr)
	return
}
