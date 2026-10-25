package middleware

import (
	"crypto/rand"

	"github.com/gin-gonic/gin"
)

// RequestIDHeader adalah nama header HTTP yang membawa request ID.
const RequestIDHeader = "X-Request-ID"

// requestIDKey adalah key untuk menyimpan request ID di gin.Context.
const requestIDKey = "request_id"

// RequestID memberi setiap request sebuah ID unik, menyimpannya di context,
// dan mengirimkannya kembali ke client lewat header X-Request-ID.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := rand.Text()

		c.Set(requestIDKey, id)
		c.Header(RequestIDHeader, id)

		c.Next()
	}
}

// GetRequestID mengambil request ID dari context.
// Hasilnya string kosong jika middleware RequestID belum dijalankan.
func GetRequestID(c *gin.Context) string {
	return c.GetString(requestIDKey)
}
