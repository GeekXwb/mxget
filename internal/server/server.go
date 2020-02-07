package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/winterssy/mxget/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type errorBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func customHTTPError(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	s, _ := status.FromError(err)
	st := runtime.HTTPStatusFromCode(s.Code())
	body := &errorBody{
		Code: st,
		Msg:  s.Message(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(st)
	_ = json.NewEncoder(w).Encode(body)
}

func init() {
	runtime.HTTPError = customHTTPError
}

func RunRPC(ctx context.Context, srv api.MusicServer, rpcPort int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", rpcPort))
	if err != nil {
		return err
	}

	rpcServer := grpc.NewServer()
	api.RegisterMusicServer(rpcServer, srv)
	return rpcServer.Serve(lis)
}

func RunRest(ctx context.Context, rpcPort int, restPort int) error {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, new(runtime.JSONBuiltin)),
	)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	endpoint := fmt.Sprintf("localhost:%d", rpcPort)
	if err := api.RegisterMusicHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", restPort),
		Handler: mux,
	}

	return srv.ListenAndServe()
}
