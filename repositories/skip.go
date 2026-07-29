package repositories

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/secret"
	"time"
)

type skipKey struct {
	KeyID string `json:"keyId"`
	Key   string `json:"key"`
}

type skipResponse struct {
	KeyID string `json:"keyId"`
	Key   string `json:"key"`
}

type SKIPRepository struct {
	baseURL          string
	remoteSystemId   string
	maxRetries       int
	backoffBaseDelay time.Duration
	conn             *http.Client
	Managed          bool
}

func NewSKIPRepository(url string, remoteSystemID string, timeout time.Duration, maxRetries int, backoffBaseDelay time.Duration, auth *KMSAuth) *SKIPRepository {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		Proxy: http.ProxyFromEnvironment,
	}
	if auth.IsClientCertAuth() { //Check for files for cert auth
		clientCert, err := tls.LoadX509KeyPair(*auth.cert, *auth.key) //load client cert and  key
		if err != nil {
			log.Fatal(err)
		}
		tr.TLSClientConfig.Certificates = []tls.Certificate{clientCert} // Set client cert for TLS auth
		caCert, err := os.ReadFile(*auth.cacert)                        //load plain ca cert
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()        //Create blank pool of trusted certs
		caCertPool.AppendCertsFromPEM(caCert)   //Add our CA
		tr.TLSClientConfig.RootCAs = caCertPool //use this pool for server auth
	}
	return &SKIPRepository{
		baseURL:          url,
		remoteSystemId:   remoteSystemID,
		maxRetries:       maxRetries,
		backoffBaseDelay: backoffBaseDelay,
		conn: &http.Client{
			Timeout:   timeout,
			Transport: tr,
		},
		Managed: true,
	}
}

// This is main fuction to proceed skip request
func (r *SKIPRepository) skipRequest(requestUrl string) (string, []byte, error) {
	res, err := r.conn.Get(requestUrl)

	if err != nil {
		return "", nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("[ERROR] server returned status: %d", res.StatusCode)

	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil, fmt.Errorf("[ERROR] failed to read response: %w", err)
	}
	defer clear(body)

	var skipResp skipResponse

	if err := json.Unmarshal(body, &skipResp); err != nil {
		return "", nil, fmt.Errorf("[ERROR] failed to parse JSON: %w", err)
	}

	if skipResp.KeyID == "" || skipResp.Key == "" {
		return "", nil, fmt.Errorf("[ERROR] received empty key or keyID from server")
	}

	var rawKey []byte
	var decodeErr error

	secret.Do(func() {
		rawKey, decodeErr = hex.DecodeString(skipResp.Key)
	})

	if decodeErr != nil {
		return "", nil, fmt.Errorf("[ERROR] failed to decode hex key: %w", decodeErr)
	}

	return skipResp.KeyID, rawKey, nil
}

func (r *SKIPRepository) GetNewKey() (string, []byte, error) {
	requestUrl := r.baseURL + "/key?remoteSystemID=" + r.remoteSystemId
	KeyID, key, err := r.skipRequest(requestUrl)

	if err != nil {
		return "", nil, fmt.Errorf("[ERROR] failed to get new key: %w", err)
	}

	return KeyID, key, nil
}

func (r *SKIPRepository) GetKeyByID(keyID *string) ([]byte, error) {
	if keyID == nil || *keyID == "" {
		return nil, fmt.Errorf("[ERROR] keyID is nil or empty")
	}
	requestUrl := r.baseURL + "/key/" + *keyID + "?remoteSystemID=" + r.remoteSystemId
	_, key, err := r.skipRequest(requestUrl)

	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to get this key: %w", err)
	}

	return key, nil
}