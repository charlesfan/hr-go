package app

import (
	"context"
	"fmt"
	"net"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/charlesfan/hr-go/config"
	"github.com/charlesfan/hr-go/repository/cache"
	"github.com/charlesfan/hr-go/repository/db/daos"
	"github.com/charlesfan/hr-go/service"
)

var server = newServer()

type Server struct {
	router *Router
}

func (s *Server) Run(c config.Config) error {
	cache.Init(c)
	state, err := cache.NewRedis().Ping(context.Background())
	if err != nil {
		fmt.Println(fmt.Errorf("cache init error: %+v", err))
	}
	fmt.Printf("cache state: %s\n", state)
	if err := daos.Init(c); err != nil {
		return err
	}
	service.Init(daos.NewDBRepoFactory())
	s.router = NewRouter(net.JoinHostPort(c.Server.Host, c.Server.Port))
	s.router.Config(c)
	s.router.Run()
	return nil
}

func newServer() *Server {
	return &Server{}
}

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "server",
		Long:         `Here is hr-go backend server.`,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			config.Init()
			c := config.NewConfig()

			server.Run(c)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "", "config file (default is $HOME/.cobra.yaml)")

	viper.BindPFlags(cmd.Flags())

	return cmd
}
