package inc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go.uber.org/zap"
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

func NewInc(config *config.Config, db *gorm.DB, logger *zap.Logger) *Inc {

	// Setting up a router for it
	router := mux.NewRouter()

	// dont' think we need a mw since auth is purely certs based

	caCert, err := ioutil.ReadFile(config.Inc.CaCert)
	if err != nil {
		logger.Fatal("Error reading CA cert", zap.Error(err))
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	}

	incServer := &http.Server{
		Addr:      ":" + config.Inc.Port,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	// Loading client certificates
	clientCert, err := tls.LoadX509KeyPair(config.Inc.PublicCert, config.Inc.PrivateKey)
	if err != nil {
		logger.Fatal("Error loading client certificates", zap.Error(err))
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{clientCert},
			},
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, time.Second*5)
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
	router.HandleFunc("/api/v1/node/replace_shard", newInc.ReplaceShardRoute).Methods("POST")
	router.HandleFunc(FragmentHealthCheckPath, newInc.ShardHealthCheckRoute).Methods("GET")
	router.HandleFunc(FragmentPath, newInc.DeleteShardRoute).Methods("DELETE")
	router.HandleFunc(FragmentOrphanedPath, newInc.OrphanedShardsRoute).Methods("POST")
	router.HandleFunc(FragmentMissingPath, newInc.MissingShardsRoute).Methods("POST")
	router.HandleFunc(CurrentFilesHealthPath, newInc.CurrentFilesFragmentsHealthCheckRoute).Methods("POST")
	router.HandleFunc(AllFilesHealthPath, newInc.AllFilesFragmentsHealthCheckRoute).Methods("POST")
	router.HandleFunc(ReplaceShardPath, newInc.AllFilesFragmentsHealthCheckRoute).Methods("POST")

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
