package main

import (
	"github.com/TelephoneTan/GoHTTPServer/net/http/server"
	"net/http"
)

func main() {
	server.NewContainer(func() []server.Service {
		return []server.Service{
			{Network: "tcp4", Address: "127.0.0.1:10086", UseTLS: false, UseGzip: true},
			{Network: "tcp6", Address: "[::1]:10086", UseTLS: false, UseGzip: true},
		}
	}, func() server.HandleFunc {
		s := server.NewServer(nil, nil).Use(manager)
		return func(writer http.ResponseWriter, request *http.Request) {
			if !s.Handle(writer, request) {
				writer.WriteHeader(http.StatusNotFound)
			}
		}
	}, nil, func(container server.Container) {
		container.ShouldListenOnDefaultPorts = func() bool {
			return false
		}
	}).Boot()
}
