/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"bookinfo.com/details/models"
	"bookinfo.com/details/tools"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var service_name = "details"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		p := ginprometheus.NewPrometheus("gin")
		// client := tools.NewClient(time.Second * 10)
		err := tools.InitJaegerClient(globalCfg.ServiceName, globalCfg.JaegerHost)
		if err != nil {
			panic(err)
		}
		p.Use(r)
		r.GET("/details", func(c *gin.Context) {
			wireContext, err := opentracing.GlobalTracer().Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(c.Request.Header),
			)
			if err != nil {
				log.Printf("read span context err %v", err)
			}
			span := opentracing.StartSpan("hander_"+globalCfg.ServiceName, opentracing.ChildOf(wireContext))
			defer span.Finish()
			details := models.Details{
				ISBN10:    "9787020008728",
				Publisher: "人民文学出版社",
				Pages:     990,
				Type:      "历史",
				Language:  "中文",
			}
			c.JSON(http.StatusOK, details)
		})

		r.Run(fmt.Sprintf(":%d", globalCfg.Port))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
