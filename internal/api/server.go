// Package server provides a server object that represents the GoCert backend
package server

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/canonical/gocert/internal/certdb"
)

type Environment struct {
	DB                      *certdb.CertificateRequestsRepository
	SendPebbleNotifications bool
	JWTSecret               []byte
}

func SendPebbleNotification(key, request_id string) error {
	cmd := exec.Command("pebble", "notify", key, fmt.Sprintf("request_id=%s", request_id))
	if err := cmd.Run(); err != nil {
		return errors.Join(errors.New("couldn't execute a pebble notify: "), err)
	}
	return nil
}

func generateJWTSecret() ([]byte, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return bytes, fmt.Errorf("failed to generate JWT secret: %w", err)
	}
	return bytes, nil
}

// NewServer creates an environment and an http server with handlers that Go can start listening to
func NewServer(port int, cert []byte, key []byte, dbPath string, pebbleNotificationsEnabled bool) (*http.Server, error) {
	serverCerts, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	db, err := certdb.NewCertificateRequestsRepository(dbPath, "CertificateRequests")
	if err != nil {
		log.Fatalf("Couldn't connect to database: %s", err)
	}

	jwtSecret, err := generateJWTSecret()
	if err != nil {
		return nil, err
	}
	env := &Environment{}
	env.DB = db
	env.SendPebbleNotifications = pebbleNotificationsEnabled
	env.JWTSecret = jwtSecret
	router := NewGoCertRouter(env)

	s := &http.Server{
		Addr: fmt.Sprintf(":%d", port),

		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{serverCerts},
		},
	}

	return s, nil
}
