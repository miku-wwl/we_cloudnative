/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"bookinfo.com/reviews/models"
	"bookinfo.com/reviews/tools"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func getRating(client *tools.Client, id string, header http.Header, parentSpan opentracing.Span) models.RatingWithStatus {
	url := fmt.Sprintf("%s%s?id=%s", globalCfg.ServerMap["ratings"], "/ratings", id)
	dataBytes, err := client.Get(url, header, "ratings", parentSpan)
	ratingWithStatus := models.RatingWithStatus{
		Base: models.Base{
			Status: 200,
		},
	}
	if err != nil {
		ratingWithStatus.Base = models.Base{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}
	} else {
		rating := models.Rating{}
		err := json.Unmarshal(dataBytes, &rating)
		if err != nil {
			ratingWithStatus.Base = models.Base{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}
		} else {
			ratingWithStatus.Rating = rating
		}
	}
	return ratingWithStatus
}

var service_name = "reviews"

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
		r.GET("/reviews", func(c *gin.Context) {
			wireContext, err := opentracing.GlobalTracer().Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(c.Request.Header),
			)
			if err != nil {
				log.Printf("read span context err %v", err)
			}
			span := opentracing.StartSpan("hander_"+globalCfg.ServiceName, opentracing.ChildOf(wireContext))
			defer span.Finish()
			hostname, ok := os.LookupEnv("HOSTNAME")
			if !ok {
				hostname = "unknown"
			}
			reviewsBody := models.ReviewsBody{
				Podname:     hostname,
				Clustername: "cluster-001",
				Reviews: []models.Review{
					{
						Text:     "“流星的光芒虽短促，但天上还有什么星能比它更灿烂，辉煌？ 当流星出现的时候，就算是永恒不变的星座，也夺不去它的光芒。” ——三国，周郎，赤壁。",
						Reviewer: "陶者无缰",
					},
					{
						Text:     "鲁迅说此书状刘备仁慈而近伪，状诸葛多谋而近妖，真是一点不假。又第106回跨越11年，第120回跨越15年，也够蛋疼",
						Reviewer: "不要焦虑浣熊妹",
					},
				},
			}
			for i := 0; i < len(reviewsBody.Reviews); i++ {
				id := strconv.Itoa(i + 1)
				reviewsBody.Reviews[i].Rating = getRating(client, id, c.Request.Header, span)
			}
			c.JSON(http.StatusOK, reviewsBody)
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
