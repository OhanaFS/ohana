package selfsign_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/selfsign"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Create a new server with the TLS config given.
func NewServer(port string, tlsconfig *tls.Config) *http.Server {
	addr := fmt.Sprintf(":%s", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	return &http.Server{
		Addr:      addr,
		Handler:   mux,
		TLSConfig: tlsconfig,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func TestSelfSign(t *testing.T) {

	ogc := config.LoadFlagsConfig()

	caPath2 := "certs/wee"
	tp := true

	t.Run("Create CA", func(t *testing.T) {
		Assert := assert.New(t)

		c := &config.FlagsConfig{}
		*c = *ogc

		// setting new configs (need vars here because we need their pointers)
		c.GenCA = &tp

		err := selfsign.ProcessFlags(c)
		Assert.Nil(err)

		// Making them from another folder

		c = &config.FlagsConfig{}
		*c = *ogc

		c.GenCA = &tp
		c.GenCAPath = &caPath2

		err = selfsign.ProcessFlags(c)
		Assert.Nil(err)

	})

	t.Run("Create Certs", func(t *testing.T) {

		Assert := assert.New(t)

		c := &config.FlagsConfig{}
		*c = *ogc

		numCerts := 2

		// Creating a certhosts.yaml file

		fakeHosts := selfsign.Hosts{Hosts: []string{"localhost", "localhost2"}}

		hostFile, err := os.Create("certhosts.yaml")
		Assert.Nil(err)
		defer func(hostFile *os.File) {
			_ = hostFile.Close()
		}(hostFile)

		encoder := yaml.NewEncoder(hostFile)
		Assert.Nil(encoder.Encode(fakeHosts))

		// Default config with 2 num of certs

		// setting new configs (need vars here because we need their pointers)
		c.GenCerts = &tp
		c.NumOfCerts = &numCerts

		err = selfsign.ProcessFlags(c)
		Assert.Nil(err)

		// Making them from another folder

		c = &config.FlagsConfig{}
		*c = *ogc

		csrPath := caPath2 + "_csr.json"
		certPath := caPath2 + "_GLOBAL_CERTIFICATE.pem"
		pkPath := caPath2 + "_PRIVATE_KEY.pem"
		genCertsPath := caPath2

		c.GenCerts = &tp
		c.NumOfCerts = &numCerts
		c.CsrPath = &csrPath
		c.CertPath = &certPath
		c.PkPath = &pkPath
		c.GenCertsPath = &genCertsPath

		err = selfsign.ProcessFlags(c)
		Assert.Nil(err)

	})

	t.Run("Run Server and Check Certs", func(t *testing.T) {

		Assert := assert.New(t)

		// Reading default certificate and adding it to tlsconfig to run the server.

		caCert, err := ioutil.ReadFile(*ogc.CertPath)
		Assert.NoError(err)
		if err != nil {
			fmt.Println(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}

		srv := NewServer("5555", tlsConfig)
		go srv.ListenAndServeTLS(*ogc.GenCertsPath+"_cert.pem", *ogc.GenCertsPath+"_key.pem")
		defer srv.Close()

		time.Sleep(100 * time.Millisecond)

		// Client with valid cert

		clientCert, err := tls.LoadX509KeyPair(*ogc.GenCertsPath+"_cert_1.pem",
			*ogc.GenCertsPath+"_key_1.pem")
		Assert.NoError(err)

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: []tls.Certificate{clientCert},
				},
			},
		}

		res, err := client.Get("https://localhost:5555")
		Assert.NoError(err)

		defer res.Body.Close()

		Assert.Equal(http.StatusOK, res.StatusCode)

		body, err := ioutil.ReadAll(res.Body)
		Assert.NoError(err)

		expected := []byte("Hello World")
		Assert.Equal(expected, body)

		// Client 2 with bad cert (Should fail when using wrong cert)

		caCert2, _ := ioutil.ReadFile(caPath2 + "_GLOBAL_CERTIFICATE.pem")
		caCertPool2 := x509.NewCertPool()
		caCertPool2.AppendCertsFromPEM(caCert2)

		clientCert2, err := tls.LoadX509KeyPair(caPath2+"_cert.pem",
			caPath2+"_key.pem")
		Assert.NoError(err)

		client2 := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool2,
					Certificates: []tls.Certificate{clientCert2},
				},
			},
		}

		res, err = client2.Get("https://localhost:5555")
		Assert.IsType(err, &url.Error{}) // should be  x509.UnknownAuthorityError{}

	})

	t.Run("Cleanup", func(t *testing.T) {
		Assert := assert.New(t)

		path, err := os.Getwd()
		Assert.NoError(err)

		Assert.True(strings.Contains(path, "ohana/selfsign"))

		// ensure we are in the right folder before we overwrite the files

		if strings.Contains(path, "ohana/selfsign") {
			caPathFolder, _ := filepath.Split(*ogc.GenCAPath)
			err = os.RemoveAll(caPathFolder)
			Assert.Nil(err)

			caPathFolder2, _ := filepath.Split(caPath2)
			err = os.RemoveAll(caPathFolder2)
			Assert.Nil(err)

			err = os.RemoveAll("certhosts.yaml")
			Assert.Nil(err)
		}
	})

}
