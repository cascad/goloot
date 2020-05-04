package erlang

import (
	"github.com/zegl/goriak/v3"
)

func Connect(host string) (*goriak.Session, error) {
	session, err := goriak.Connect(goriak.ConnectOpts{
		Address: host,
		//Port:    8087,
	})

	return session, err
}
