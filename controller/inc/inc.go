package inc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"
	"io/ioutil"
	"net/http"
	"strconv"
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
		},
		// Quick timeout for quick failover
		Timeout: time.Second,
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

func MOCKDumpLogs(db *gorm.DB, i *Inc) {

	const (
		fatalCount   = 5
		errorCount   = 10
		warningCount = 15
		infoCount    = 20
		debugCount   = 25
		traceCount   = 30
	)

	logsToCreate := map[int8]int{
		dbfs.LogServerFatal:   fatalCount,
		dbfs.LogServerError:   errorCount,
		dbfs.LogServerWarning: warningCount,
		dbfs.LogServerInfo:    infoCount,
		dbfs.LogServerDebug:   debugCount,
		dbfs.LogServerTrace:   traceCount,
	}

	dbfsLogger := dbfs.NewLogger(db, i.ServerName)

	for logType, count := range logsToCreate {
		for i := 0; i < count; i++ {
			switch logType {
			case dbfs.LogServerFatal:
				dbfsLogger.LogFatal("fatal log " + strconv.Itoa(i))
			case dbfs.LogServerError:
				dbfsLogger.LogError("error log " + strconv.Itoa(i))
			case dbfs.LogServerWarning:
				dbfsLogger.LogWarning("warning log " + strconv.Itoa(i))
			case dbfs.LogServerInfo:
				dbfsLogger.LogInfo("info log " + strconv.Itoa(i))
			case dbfs.LogServerDebug:
				dbfsLogger.LogDebug("debug log " + strconv.Itoa(i))
			case dbfs.LogServerTrace:
				dbfsLogger.LogTrace("trace log	" + strconv.Itoa(i))
			}
		}
	}

}
