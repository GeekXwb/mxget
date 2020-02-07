package serve

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/winterssy/mxget/internal/server"
	"github.com/winterssy/mxget/pkg/service"
)

var (
	rpcPort  int
	restPort int
)

var CmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Run mxget as an API server",
}

func Run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	srv := new(service.MusicServerImpl)
	go func() {
		_ = server.RunRest(ctx, rpcPort, restPort)
	}()
	_ = server.RunRPC(ctx, srv, rpcPort)
}

func init() {
	CmdServe.Flags().IntVar(&rpcPort, "rpc-port", 8090, "rpc server listening port")
	CmdServe.Flags().IntVar(&restPort, "rest-port", 8080, "rest server listening port")
	CmdServe.Run = Run
}
