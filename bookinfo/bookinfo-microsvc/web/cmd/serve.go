/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"bookinfo.com/web/tools"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	"github.com/uber/jaeger-client-go"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var service_name = "web"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		p := ginprometheus.NewPrometheus("gin")
		client := tools.NewClient(time.Second * 10)
		err := tools.InitJaegerClient(globalCfg.ServiceName, globalCfg.JaegerHost)
		if err != nil {
			panic(err)
		}
		p.Use(r)
		r.Static("/static", "./static")
		r.GET("/", func(c *gin.Context) {
			data, _ := os.ReadFile("templates/index.html")
			c.Data(http.StatusOK, "text/html,charset=utf-8", data)
		})
		r.GET("/productpage", func(c *gin.Context) {
			span := opentracing.StartSpan("hander_" + globalCfg.ServiceName)
			defer span.Finish()
			userType := c.Query("u")
			if userType == "" {
				userType = "normal"
			}
			c.Request.Header.Set("u", userType)
			url := fmt.Sprintf("%s%s", globalCfg.ServerMap["productpage"], "/productpage")
			data, err := client.Get(url, c.Request.Header, "productpage", span)
			if jaegerContext, ok := span.Context().(jaeger.SpanContext); ok {
				tranceid := jaegerContext.TraceID().String()
				c.Header("X-Trance-ID", tranceid)
			}
			if err != nil {
				c.Data(http.StatusInternalServerError, "text/html,charset=utf-8", []byte(err.Error()))
			} else {
				c.Data(http.StatusOK, "text/html,charset=utf-8", data)
			}
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
