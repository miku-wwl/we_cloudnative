/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"html/template"
	"net/http"
	_ "net/http/pprof"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type Base struct {
	Status int
	Error  string
}
type Rating struct {
	Base
	Stars int
	Color string
}
type Review struct {
	Rating   Rating
	Text     string
	Reviewer string
}
type Reviews struct {
	Base
	Reviews     []Review
	Podname     string
	Clustername string
}
type Details struct {
	Base
	ISBN10    string
	Publisher string
	Pages     int
	Type      string
	Language  string
}
type Product struct {
	Title           string
	DescriptionHtml string
}
type ProductPage struct {
	User    string
	Product Product
	Details Details
	Reviews Reviews
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "server",
	Long:  `start bookinfo server.`,
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		p := ginprometheus.NewPrometheus("gin")
		p.Use(r)
		r.Static("/static", "./static")
		r.GET("/productpage", func(c *gin.Context) {
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
			product_page := ProductPage{
				User: "admin",
				Product: Product{
					DescriptionHtml: `三国演义》又名《三国志演义》、《三国志通俗演义》，是我国小说史上最著名最杰出的长篇章回体历史小说。 《三国演义》的作者是元末明初人罗贯中，由毛纶，毛宗岗父子批改。在其成书前，“三国故事”已经历了数百年的历史发展过程。在唐代，三国故事已广为流传，连儿童都很熟悉。随着市民文艺的发展，宋代的“说话”艺人，已有专门说三国故事的，当时称为“说三分”。元代出现的《三国志平话》，实际上是从说书人使用的本子，虽较简略粗糙，但已初肯《三国演义》的规模。罗贯中在群众传说和民间艺人创作的基础上，又依据陈寿《三国志》及裴松之注中所征引的资料（还包括《世说新语》及注中的资料），经过巨大的创作劳动，写在了规模宏伟的巨著——《三国志通俗演义》全书24卷，240回。后来经过毛纶，毛宗岗父子批改，形成我们现在所见的《三国演义》120回版

					由于人文出版社选取的毛批版本和明嘉靖版的《三国志通俗演义》颇有出入，将拥刘抑曹的思想发展到一个几乎病态的角度，且大量地删除了原文中赞扬曹操一方人物的诗词和评论，所以给世人造成一种三国演义打击曹操、歌颂刘备的错误印象，且该批改版将刘备一方无限神化，甚至将诸葛亮准备将魏延烧死在上方谷这样的情节统统删除，将诸葛说成是一个完美近乎神人的形象，这是大家阅读时候值得注意的地方。
					
					由于该版本选用的是毛批改版，而且没有将毛家父子的批语选入，就会给人造成一种错觉，认为该版本和毛批本是不相同的，其实本书和毛批本完全一致，不过是删除了批语部分，加入了编者的注脚。`,
					Title: "三国演义",
				},
				Details: Details{
					Base: Base{
						Status: 200,
					},
					ISBN10:    "9787020008728",
					Publisher: "人民文学出版社",
					Pages:     990,
					Type:      "历史",
					Language:  "中文",
				},
				Reviews: Reviews{
					Base: Base{
						Status: 200,
					},
					Podname:     "pod-001",
					Clustername: "cluster-001",
					Reviews: []Review{
						{
							Rating: Rating{
								Base: Base{
									Status: 200,
								},
								Stars: 4,
								Color: "red",
							},
							Text:     "“流星的光芒虽短促，但天上还有什么星能比它更灿烂，辉煌？ 当流星出现的时候，就算是永恒不变的星座，也夺不去它的光芒。” ——三国，周郎，赤壁。",
							Reviewer: "陶者无缰",
						},
						{
							Rating: Rating{
								Base: Base{
									Status: 200,
								},
								Stars: 5,
								Color: "red",
							},
							Text:     "鲁迅说此书状刘备仁慈而近伪，状诸葛多谋而近妖，真是一点不假。又第106回跨越11年，第120回跨越15年，也够蛋疼",
							Reviewer: "不要焦虑浣熊妹",
						},
					},
				},
			}
			buf := bytes.Buffer{}
			err = tmpl.Execute(&buf, product_page)
			if err != nil {
				c.Data(http.StatusInternalServerError, "text/html,charset=utf-8", []byte(err.Error()))
				return
			}
			c.Data(http.StatusOK, "text/html,charset=utf-8", buf.Bytes())
		})
		r.GET("/", func(c *gin.Context) {
			data, _ := os.ReadFile("templates/index.html")
			c.Data(http.StatusOK, "text/html,charset=utf-8", data)
		})
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			http.ListenAndServe(":6060", nil)
		}()
		r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
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
