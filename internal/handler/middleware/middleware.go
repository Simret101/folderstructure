package middleware

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/chacha20poly1305"
)

type EncryptedPayload struct {
	CipherText string `json:"cipherText"`
	IV         string `json:"iv"`
	Tag        string `json:"tag"`
}

type CryptoMiddleware struct {
	key []byte
}

func NewCryptoMiddleware(key string) (*CryptoMiddleware, error) {

	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 crypto key: %w", err)
	}

	if len(decoded) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("crypto key must be 32 bytes, got %d", len(decoded))
	}

	return &CryptoMiddleware{key: decoded}, nil
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (m *CryptoMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.GetHeader("enable_encryption") != "enabled" {
			c.Next()
			return
		}

		var req EncryptedPayload
		if err := c.ShouldBindJSON(&req); err != nil {
			if err != io.EOF {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "invalid encrypted payload",
				})
				return
			}
		}

		if req.CipherText != "" {

			nonce, err := base64.StdEncoding.DecodeString(req.IV)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			ct, err := base64.StdEncoding.DecodeString(req.CipherText)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			tag, err := base64.StdEncoding.DecodeString(req.Tag)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			aead, err := chacha20poly1305.New(m.key)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			plaintext, err := aead.Open(nil, nonce, append(ct, tag...), nil)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "decryption failed",
				})
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewReader(plaintext))
			c.Request.ContentLength = int64(len(plaintext))
		}

		writer := &bodyWriter{
			ResponseWriter: c.Writer,
			body:           new(bytes.Buffer),
		}
		c.Writer = writer

		c.Next()

		rawResponse := writer.body.Bytes()
		if len(rawResponse) == 0 {
			return
		}

		var resp interface{}
		if err := json.Unmarshal(rawResponse, &resp); err != nil {
			resp = string(rawResponse)
		}

		plain, _ := json.Marshal(resp)

		nonce := make([]byte, chacha20poly1305.NonceSize)
		if _, err := rand.Read(nonce); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		aead, err := chacha20poly1305.New(m.key)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		cipher := aead.Seal(nil, nonce, plain, nil)

		tag := cipher[len(cipher)-aead.Overhead():]
		ct := cipher[:len(cipher)-aead.Overhead()]

		c.JSON(http.StatusOK, EncryptedPayload{
			CipherText: base64.StdEncoding.EncodeToString(ct),
			IV:         base64.StdEncoding.EncodeToString(nonce),
			Tag:        base64.StdEncoding.EncodeToString(tag),
		})
	}
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}
