package inc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/OhanaFS/ohana/config"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Inc struct {
	ShardsLocation    string
	ServerName        string
	HostName          string
	Port              string
	CaCert            string
	PublicCert        string
	PrivateKey        string
	CaCertPool        *x509.CertPool
	ClientCertificate tls.Certificate
	HttpClient        *http.Client
	HttpServer        *http.Server
	BindIp            string
	Db                *gorm.DB
	Shutdown          chan bool
}

func NewInc(config *config.Config, db *gorm.DB) *Inc {

	// TODO ? Put a zap logger here

	// Setting up a router for it
	router := mux.NewRouter()

	// dont' think we need a mw since auth is purely certs based

	caCert, err := ioutil.ReadFile(config.Inc.CaCert)
	if err != nil {
		panic(err) // TODO put logger here
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	incServer := &http.Server{
		Addr:      config.Inc.HostName + ":" + config.Inc.Port,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	// Loading client certificates
	clientCert, err := tls.LoadX509KeyPair(config.Inc.PublicCert, config.Inc.PrivateKey)
	if err != nil {
		panic(err) // TODO put logger here
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{clientCert},
			},
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, time.Millisecond*100)
			},
		},
	}

	// register routes

	newInc := &Inc{
		ShardsLocation:    config.Stitch.ShardsLocation,
		ServerName:        config.Inc.ServerName,
		HostName:          config.Inc.HostName,
		Port:              config.Inc.Port,
		CaCert:            config.Inc.CaCert,
		PublicCert:        config.Inc.PublicCert,
		PrivateKey:        config.Inc.PrivateKey,
		BindIp:            config.Inc.BindIp,
		CaCertPool:        caCertPool,
		ClientCertificate: clientCert,
		HttpClient:        client,
		HttpServer:        incServer,
		Db:                db,
		Shutdown:          make(chan bool),
	}

	router.HandleFunc("/api/v1/node/ping", Pong)
	router.HandleFunc("/api/v1/node/details", newInc.ReturnServerDetails)
	router.HandleFunc("/api/v1/node/shard/{shardId}", newInc.handleShardStream)

	// start server
	go func() {
		err := incServer.ListenAndServeTLS(config.Inc.PublicCert, config.Inc.PrivateKey)
		if err != nil {
			if err != http.ErrServerClosed {
				fmt.Println("Server closed")
			}
		}
	}()

	return newInc
}
