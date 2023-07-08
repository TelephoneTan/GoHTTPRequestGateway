package main

import (
	"github.com/TelephoneTan/GoHTTPRequest/net/http"
	"github.com/TelephoneTan/GoHTTPServer/net/http/method"
	"github.com/TelephoneTan/GoHTTPServer/net/http/server"
	"github.com/TelephoneTan/GoHTTPServer/types"
	"github.com/TelephoneTan/GoPromise/async/promise"
	"mime"
	httpGo "net/http"
)

type packet struct {
	bad             bool
	incomingRequest http.Request
}

var manager = server.NewResourceManager(func() types.WordList {
	return types.WordList{{"gateway"}, {"网关"}}
}, nil, func(r server.ResourceManager[packet]) {
	r.Guide = map[method.Method]server.ResourceRequestHandler[packet]{
		method.GET: {
			Peek: func(r *httpGo.Request, paths server.PathPack) (pack packet, hijacked bool) {
				goto start
			hijack:
				return pack, true
			bad:
				pack.bad = true
				goto hijack
			start:
				reqJSON := r.Header.Get("HTTP-Request")
				if reqJSON == "" {
					goto bad
				}
				reqJSON, err := new(mime.WordDecoder).DecodeHeader(reqJSON)
				if err != nil {
					goto bad
				}
				pack.incomingRequest = http.NewRequest().Deserialize(reqJSON)
				goto hijack
			},
			Reply: func(w httpGo.ResponseWriter, pack packet) {
				if pack.bad {
					w.WriteHeader(httpGo.StatusBadRequest)
					return
				}
				success := promise.Then(pack.incomingRequest.ByteSlice(), promise.FulfilledListener[http.Result[[]byte], any]{
					OnFulfilled: func(value http.Result[[]byte]) any {
						_, _ = w.Write([]byte(value.Request.Serialize()))
						return nil
					},
				})
				promise.Catch(success, promise.RejectedListener[any]{
					OnRejected: func(reason error) any {
						w.Header().Set("HTTP-Request-Error", mime.QEncoding.Encode("utf-8", reason.Error()))
						w.WriteHeader(httpGo.StatusOK)
						return nil
					},
				}).Await()
			},
		},
	}
})
