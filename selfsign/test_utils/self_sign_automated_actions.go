package selfsigntestutils

import (
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/selfsign"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type CertConfig struct {
	CaCertPath     string
	PublicCertPath string
	PrivateKeyPath string
}

// GenCertsTest generates certs for testing.
// It disregards flags given at compilation and thus can be used multiple times.
// It will also automatically generate certhosts.yaml file.
func GenCertsTest(tempDir string) (*CertConfig, error) {
	// Setting up certs for configs

	trueBool := true
	tempDirCA := filepath.Join(tempDir, "certificates/main")
	tempDirCerts := filepath.Join(tempDir, "certificates/output")
	tempCsrPath := filepath.Join(tempDir, "certificates/main_csr.json")
	tempCertPath := filepath.Join(tempDir, "certificates/main_GLOBAL_CERTIFICATE.pem")
	tempPkPath := filepath.Join(tempDir, "certificates/main_PRIVATE_KEY.pem")
	tempHostsFile := filepath.Join(tempDir, "certhosts.yaml")
	tempNumOfCerts := 1
	ogc := &config.FlagsConfig{
		GenCA:        &trueBool,
		GenCerts:     &trueBool,
		GenCAPath:    &tempDirCA,
		GenCertsPath: &tempDirCerts,
		CsrPath:      &tempCsrPath,
		CertPath:     &tempCertPath,
		PkPath:       &tempPkPath,
		AllHosts:     &tempHostsFile,
		NumOfCerts:   &tempNumOfCerts,
	}

	fakeHosts := selfsign.Hosts{Hosts: []string{"localhost", "localhost2"}}

	hostFile, err := os.Create(filepath.Join(tempDir, "certhosts.yaml"))
	if err != nil {
		return nil, err
	}
	defer hostFile.Close()

	encoder := yaml.NewEncoder(hostFile)
	if err := encoder.Encode(fakeHosts); err != nil {
		return nil, err
	}

	err = selfsign.ProcessFlags(ogc, true)

	if err != nil {
		return nil, err
	}

	return &CertConfig{
		CaCertPath:     filepath.Join(tempDir, "certificates/main_GLOBAL_CERTIFICATE.pem"),
		PublicCertPath: filepath.Join(tempDir, "certificates/output_cert.pem"),
		PrivateKeyPath: filepath.Join(tempDir, "certificates/output_key.pem"),
	}, nil

}
