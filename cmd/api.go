package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/spf13/cobra"
	"github.com/three-body/hertz-scaffold/biz/dal"
	"github.com/three-body/hertz-scaffold/biz/dal/query"
	"github.com/three-body/hertz-scaffold/biz/router"
	"github.com/three-body/hertz-scaffold/config"
	"github.com/three-body/hertz-scaffold/docs"
)

//	@title			Hertz Scaffold API
//	@version		0.1.0
//	@description	API Server for Hertz Scaffold
//	@termsOfService

//	@contact.name	threebody
//	@contact.url	https://github.com/three-body
//	@contact.email	shiyi@threebody.xyz

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/
//	@schemes	http

//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

//	@externalDocs.description	How to write the API docs
//	@externalDocs.url			https://github.com/swaggo/swag/blob/master/README.md

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		docs.SwaggerInfo.Host = config.GetConf().Swagger.Host
		docs.SwaggerInfo.BasePath = config.GetConf().Swagger.BasePath
		docs.SwaggerInfo.Schemes = strings.Split(config.GetConf().Swagger.Schemes, ",")
		docs.SwaggerInfo.Version = config.GetConf().Swagger.Version

		if err := dal.InitMySQL(); err != nil {
			fmt.Printf("dal init failed! %v\n", err)
			os.Exit(1)
		}
		query.SetDefault(dal.DB)

		h := server.New(
			server.WithExitWaitTime(time.Second*time.Duration(config.GetConf().Hertz.ExitWaitTime)),
			server.WithHostPorts(config.GetConf().Hertz.Address),
		)
		h.Name = config.ProjectName
		h.OnRun = append(h.OnRun, func(ctx context.Context) error {
			return nil
		})
		h.OnShutdown = append(h.OnShutdown, func(ctx context.Context) {
		})
		router.Register(h)
		h.Spin()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
