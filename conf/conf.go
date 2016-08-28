package conf

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gravitational/teleconsole/lib"
	"github.com/gravitational/teleport/lib/client"
	"github.com/gravitational/trace"
)

// Config stores the configuration of Teleconsole process
type Config struct {
	// APIEndpointURL is the API of the Teleconsole server API
	APIEndpointURL *url.URL

	// Verbosity defines the level of debugging output (greater means
	// more output)
	Verbosity int

	// when set, it means that instead of launching shell, another
	// command is launched
	RunCommand string

	// command line arguments
	Args []string

	// if 'true', the client will trust unknown SSL certificates
	// can be set via -insecure flag
	InsecureHTTPS bool

	// Ports to forward
	ForwardPorts []client.ForwardedPort

	// Forward-by-invite:
	ForwardPort *client.ForwardedPort
}

// Get() returns Teleconsole configuration: default values overwritten
// via config file
func Get() (c *Config, err error) {
	u, err := user.Current()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	// read ini-file ~/.teleconsolerc
	configFile := filepath.Join(u.HomeDir, DefaultConfigFileName)
	i, err := lib.ParseIniFile(configFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, trace.Wrap(err)
		}
	}

	c = &Config{}

	// apply ini-file vlaues to config:
	serverHostPort := i.GetOrDefault("", "server",
		net.JoinHostPort(DefaultServerHost, DefaultServerPort))
	err = c.SetServer(serverHostPort)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return c, nil
}

func (this *Config) SetServer(hostPort string) (err error) {
	var host, port string
	// missing port spec?
	if strings.LastIndex(hostPort, ":") < 0 {
		port = DefaultServerPort
		host = hostPort
	} else {
		host, port, err = net.SplitHostPort(hostPort)
		if err != nil {
			return trace.Wrap(err)
		}
	}
	// validate endpoint URL:
	this.APIEndpointURL, err = url.Parse(fmt.Sprintf("https://%s:%s", host, port))
	return trace.Wrap(err)
}
