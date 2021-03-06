// Package openssh provides a wrapper around the ssh command-line tool,
// allowing to run commands on remote machines.
package openssh

import (
	"errors"
	"strconv"
	"strings"

	"github.com/mlafeldt/chef-runner/exec"
)

// A Client is an OpenSSH client.
type Client struct {
	ConfigFile  string
	Host        string
	User        string
	Port        int
	PrivateKeys []string
	Options     []string
}

// NewClient creates a new Client from the given host string. The host string
// has the format [user@]hostname[:port]
func NewClient(host string) (*Client, error) {
	var user string
	a := strings.Split(host, "@")
	if len(a) > 1 {
		user = a[0]
		host = a[1]
	}

	var port int
	a = strings.Split(host, ":")
	if len(a) > 1 {
		host = a[0]
		var err error
		if port, err = strconv.Atoi(a[1]); err != nil {
			return nil, errors.New("invalid SSH port")
		}
	}

	return &Client{Host: host, User: user, Port: port}, nil
}

// Command returns the ssh command that will be executed by Copy.
func (c Client) Command(args []string) []string {
	cmd := []string{"ssh"}

	if c.ConfigFile != "" {
		cmd = append(cmd, "-F", c.ConfigFile)
	}

	if c.User != "" {
		cmd = append(cmd, "-l", c.User)
	}

	if c.Port != 0 {
		cmd = append(cmd, "-p", strconv.Itoa(c.Port))
	}

	for _, pk := range c.PrivateKeys {
		cmd = append(cmd, "-i", pk)
	}

	for _, o := range c.Options {
		cmd = append(cmd, "-o", o)
	}

	if c.Host != "" {
		cmd = append(cmd, c.Host)
	}

	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	return cmd
}

// RunCommand uses ssh to execute a command on a remote machine.
func (c Client) RunCommand(args []string) error {
	if len(args) == 0 {
		return errors.New("no command given")
	}
	if c.Host == "" {
		return errors.New("no host given")
	}
	return exec.RunCommand(c.Command(args))
}

// Shell returns a connection string that can be used by tools like rsync. Each
// argument is double-quoted to preserve spaces.
func (c Client) Shell() string {
	cmd := c.Command([]string{})
	var quoted []string
	for _, i := range cmd[:len(cmd)-1] {
		quoted = append(quoted, "'"+i+"'")
	}
	return strings.Join(quoted, " ")
}
