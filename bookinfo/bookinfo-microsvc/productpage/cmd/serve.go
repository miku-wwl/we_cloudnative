/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"bookinfo.com/productpage/models"
	"bookinfo.com/productpage/tools"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// 获取details服务的数据
func getDetails(client *tools.Client, header http.Header, parentSpan opentracing.Span) models.DetailsWithStatus {
	url := fmt.Sprintf("%s%s", globalCfg.ServerMap["details"], "/details")
	dataBytes, err := client.Get(url, header, "details", parentSpan)
	detailsWithStatus := models.DetailsWithStatus{
		Base: models.Base{
			Status: 200,
		},
	}
	if err != nil {
		detailsWithStatus.Base = models.Base{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}
	} else {
		details := models.Details{}
		err := json.Unmarshal(dataBytes, &details)
		if err != nil {
			detailsWithStatus.Base = models.Base{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}
		} else {
			detailsWithStatus.Details = details
		}

	}
	return detailsWithStatus
}

func getReviews(client *tools.Client, header http.Header, parentSpan opentracing.Span) models.ReviewsWithStatus {
	url := fmt.Sprintf("%s%s", globalCfg.ServerMap["reviews"], "/reviews")
	dataBytes, err := client.Get(url, header, "reviews", parentSpan)
	reviewsWithStatus := models.ReviewsWithStatus{
		Base: models.Base{
			Status: 200,
		},
	}
	if err != nil {
		reviewsWithStatus.Base = models.Base{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}
	} else {
		reviewsBody := models.ReviewsBody{}
		err := json.Unmarshal(dataBytes, &reviewsBody)
		if err != nil {
			reviewsWithStatus.Base = models.Base{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}
		} else {
			reviewsWithStatus.ReviewsBody = reviewsBody
		}
	}
	return reviewsWithStatus
}

var service_name = "productpage"

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
		r.GET("/productpage", func(c *gin.Context) {
			wireContext, err := opentracing.GlobalTracer().Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(c.Request.Header),
			)
			if err != nil {
				log.Printf("read span context err %v", err)
			}
			span := opentracing.StartSpan("hander_"+globalCfg.ServiceName, opentracing.ChildOf(wireContext))
			defer span.Finish()
			data, _ := os.ReadFile("templates/productpage.html")
			funcs := template.FuncMap{
				"html_format": func(s string) template.HTML {
					return template.HTML(s)
				},
				"inRange": func(stars int) []int {
					return make([]int, stars)
				},
				"sub": func(a, b int) int {
					return a - b
				},
			}
			tmpl, err := template.New("productpage").Funcs(funcs).Parse(string(data))
			if err != nil {
				c.Data(http.StatusInternalServerError, "text/html,charset=utf-8", []byte(err.Error()))
				return
			}
			product_page := models.ProductPage{
				User: "admin",
				Product: models.Product{
					DescriptionHtml: `三国演义》又名《三国志演义》、《三国志通俗演义》，是我国小说史上最著名最杰出的长篇章回体历史小说。 《三国演义》的作者是元末明初人罗贯中，由毛纶，毛宗岗父子批改。在其成书前，“三国故事”已经历了数百年的历史发展过程。在唐代，三国故事已广为流传，连儿童都很熟悉。随着市民文艺的发展，宋代的“说话”艺人，已有专门说三国故事的，当时称为“说三分”。元代出现的《三国志平话》，实际上是从说书人使用的本子，虽较简略粗糙，但已初肯《三国演义》的规模。罗贯中在群众传说和民间艺人创作的基础上，又依据陈寿《三国志》及裴松之注中所征引的资料（还包括《世说新语》及注中的资料），经过巨大的创作劳动，写在了规模宏伟的巨著——《三国志通俗演义》全书24卷，240回。后来经过毛纶，毛宗岗父子批改，形成我们现在所见的《三国演义》120回版
	
						由于人文出版社选取的毛批版本和明嘉靖版的《三国志通俗演义》颇有出入，将拥刘抑曹的思想发展到一个几乎病态的角度，且大量地删除了原文中赞扬曹操一方人物的诗词和评论，所以给世人造成一种三国演义打击曹操、歌颂刘备的错误印象，且该批改版将刘备一方无限神化，甚至将诸葛亮准备将魏延烧死在上方谷这样的情节统统删除，将诸葛说成是一个完美近乎神人的形象，这是大家阅读时候值得注意的地方。
						
						由于该版本选用的是毛批改版，而且没有将毛家父子的批语选入，就会给人造成一种错觉，认为该版本和毛批本是不相同的，其实本书和毛批本完全一致，不过是删除了批语部分，加入了编者的注脚。`,
					Title: "三国演义",
				},
			}
			product_page.Details = getDetails(client, c.Request.Header, span)
			product_page.Reviews = getReviews(client, c.Request.Header, span)
			buf := bytes.Buffer{}
			err = tmpl.Execute(&buf, product_page)
			if err != nil {
				c.Data(http.StatusInternalServerError, "text/html,charset=utf-8", []byte(err.Error()))
				return
			}
			c.Data(http.StatusOK, "text/html,charset=utf-8", buf.Bytes())
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
