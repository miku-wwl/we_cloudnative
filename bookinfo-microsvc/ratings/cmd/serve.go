/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"

	"bookinfo.com/ratings/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		p := ginprometheus.NewPrometheus("gin")
		// client := tools.NewClient(time.Second * 10)
		p.Use(r)
		r.Static("/static", "./static")
		dataMap := make(map[string]models.Rating)
		dataMap["1"] = models.Rating{
			Stars: 4,
			Color: "red",
		}
		dataMap["2"] = models.Rating{
			Stars: 5,
			Color: "red",
		}
		r.GET("/ratings", func(c *gin.Context) {
			id := c.Query("id")
			c.JSON(http.StatusOK, dataMap[id])
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
