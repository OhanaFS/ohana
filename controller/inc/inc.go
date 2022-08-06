package inc

import (
	"github.com/OhanaFS/ohana/config"
	"gorm.io/gorm"
)

type Inc struct {
	ShardsLocation string
	ServerName     string
	HostName       string
	Port           string
	CaCert         string
	PublicCert     string
	PrivateKey     string
	BindIp         string
	Db             *gorm.DB
	Shutdown       chan bool
}

func NewInc(config *config.Config, db *gorm.DB) *Inc {

	newInc := &Inc{
		ShardsLocation: config.Stitch.ShardsLocation,
		ServerName:     config.Inc.ServerName,
		HostName:       config.Inc.HostName,
		Port:           config.Inc.Port,
		CaCert:         config.Inc.CaCert,
		PublicCert:     config.Inc.PublicCert,
		PrivateKey:     config.Inc.PrivateKey,
		BindIp:         config.Inc.BindIp,
		Db:             db,
		Shutdown:       make(chan bool),
	}

	return newInc
}
