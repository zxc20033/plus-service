package main

import (
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/zxc20033/plus-service/pb"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/plus", func(c *gin.Context) {
		value1, err := strconv.Atoi(c.Query("value1"))
		value2, err := strconv.Atoi(c.Query("value2"))

		if err != nil {
			c.JSON(403, gin.H{
				"err": "Null value.",
			})
			panic("Null value.")
		}
		// 連線到遠端 gRPC 伺服器。
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			c.JSON(403, gin.H{
				"err": err,
			})
			panic("Server no response.")
		}
		defer conn.Close()

		// 建立新的 Calculator 客戶端，所以等一下就能夠使用 Calculator 的所有方法。
		cal := pb.NewCalculatorClient(conn)

		// 傳送新請求到遠端 gRPC 伺服器 Calculator 中，並呼叫 Plus 函式，讓兩個數字相加。
		r, err := cal.Plus(context.Background(), &pb.CalcRequest{NumberA: int32(value1), NumberB: int32(value2)})
		if err != nil {
			c.JSON(403, gin.H{
				"err": err,
			})
			panic(err)
		}
		c.JSON(200, gin.H{
			"result": r.Result,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
