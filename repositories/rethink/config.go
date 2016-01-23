package rethink

import (
	"time"

	"github.com/dancannon/gorethink"
)

// Config for repo
type Config struct {
	Address      string        `yaml:"address,omitempty"`
	Addresses    []string      `yaml:"addresses,omitempty"`
	Database     string        `yaml:"database,omitempty"`
	AuthKey      string        `yaml:"authkey,omitempty"`
	Timeout      time.Duration `yaml:"timeout,omitempty"`
	WriteTimeout time.Duration `yaml:"write_timeout,omitempty"`
	ReadTimeout  time.Duration `yaml:"read_timeout,omitempty"`

	//TODO:
	//TLSConfig    *tls.Config   `yaml:"tlsconfig,omitempty"`

	MaxIdle int `yaml:"max_idle,omitempty"`
	// By default a maximum of 2 connections are opened per host.
	MaxOpen int `yaml:"max_open,omitempty"`

	// DiscoverHosts is used to enable host discovery, when true the driver
	// will attempt to discover any new nodes added to the cluster and then
	// start sending queries to these new nodes.
	DiscoverHosts bool `yaml:"discover_hosts,omitempty"`
	// NodeRefreshInterval is used to determine how often the driver should
	// refresh the status of a node.
	//
	// Deprecated: This function is no longer used due to changes in the
	// way hosts are selected.
	NodeRefreshInterval time.Duration `yaml:"node_refresh_interval,omitempty"`
	// HostDecayDuration is used by the go-hostpool package to calculate a weighted
	// score when selecting a host. By default a value of 5 minutes is used.
	HostDecayDuration time.Duration `yaml:"node_decay_duration,omitempty"`

	// Indicates whether the cursors running in this session should use json.Number instead of float64 while
	// unmarshaling documents with interface{}. The default is `false`.
	UseJSONNumber bool `yaml:"use_json_number,omitempty"`
}

func (c *Config) getConnOps() gorethink.ConnectOpts {
	return gorethink.ConnectOpts{
		Address:             c.Address,
		Addresses:           c.Addresses,
		Database:            c.Database,
		AuthKey:             c.AuthKey,
		Timeout:             c.Timeout,
		WriteTimeout:        c.WriteTimeout,
		ReadTimeout:         c.ReadTimeout,
		MaxIdle:             c.MaxIdle,
		MaxOpen:             c.MaxOpen,
		DiscoverHosts:       c.DiscoverHosts,
		NodeRefreshInterval: c.NodeRefreshInterval,
		HostDecayDuration:   c.HostDecayDuration,
		UseJSONNumber:       c.UseJSONNumber,
	}
}

// Session gets rethink session
func (c *Config) Session() (*gorethink.Session, error) {
	return gorethink.Connect(c.getConnOps())
}
