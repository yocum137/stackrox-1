package cas

import (
	"context"
	"os"
	"path/filepath"
)

const (
	caCert      = "/Users/nchander/go/src/github.com/stackrox/stackrox/local/database-restore/full/current/keys/ca.pem"
	caKey       = "/Users/nchander/go/src/github.com/stackrox/stackrox/local/database-restore/full/current/keys/ca-key.pem"
	jwtKeyInDer = "/Users/nchander/go/src/github.com/stackrox/stackrox/local/database-restore/full/current/keys/jwt-key.der"
	jwtKeyInPem = "/Users/nchander/go/src/github.com/stackrox/stackrox/local/database-restore/full/current/keys/jwt-key.pem"
)

// NewCertsBackup returns a generator of certificate backups.
func NewCertsBackup() *CertsBackup {
	// Include jwt key in either der or in pem format, preferable in der format.
	jwtKey := jwtKeyInDer
	if _, err := os.Stat(jwtKeyInDer); os.IsNotExist(err) {
		jwtKey = jwtKeyInPem
	}
	return &CertsBackup{
		certFiles: []string{caCert, caKey, jwtKey},
	}
}

// CertsBackup is an implementation of a PathMapGenerator which generate the layout of cert files to backup.
type CertsBackup struct {
	certFiles []string
}

// GeneratePathMap generates the map from the path within backup to its source certificate.
func (c *CertsBackup) GeneratePathMap(_ context.Context) (map[string]string, error) {
	certMap := make(map[string]string)
	// Put all the certs under the same root directory.
	for _, p := range c.certFiles {
		certMap[filepath.Base(p)] = p
	}
	return certMap, nil
}
