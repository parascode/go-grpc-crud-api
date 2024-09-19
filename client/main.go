package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	pb "example.com/go-grpc-crud-api/proto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

type Company struct {
	ID			string 		`json:"id"`
	Name 		string 		`json:"Name"`
	Address 	string		`json:"Address"`
	Location 	string 	`json:"Location"`
}

type Pagination struct {
	PageSize	int64 `json:"pageSize"`
	PageToken	string	`json:"pageToken"`
}

func main() {
	flag.Parse()
	// conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err :=grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewCompanyServiceClient(conn)

	r := gin.Default()
	r.Use(cors.Default())
	
	r.POST("/all-companies", func(ctx *gin.Context) {
		body := json.NewDecoder(ctx.Request.Body)
		var t Pagination
		body.Decode(&t)

		fmt.Println("client start time" + time.Now().String())
		res, err := client.GetCompanies(ctx, &pb.ReadCompaniesRequest{
			// PageSize: strconv.Atoi((ctx.Params.ByName("pageSize"))),
			// PageToken: ctx.,
			PageSize: t.PageSize,
			PageToken: t.PageToken,
		})
		
		fmt.Println("client end time" + time.Now().String())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error" : err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"companies" : res.Companies,
			"nextPageToken" : res.NextPageToken,
		})
	})
	r.GET("/companies/:id", func(ctx *gin.Context){
		id := ctx.Param("id")
		res, err := client.GetCompany(ctx, &pb.ReadCompanyRequest{Id: id})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message" : err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"company": res.Company,
		})
	})
	r.POST("/companies", func(ctx *gin.Context) {
		var company pb.Company
  
		err := ctx.ShouldBind(&company)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		fmt.Println(company);
		data := &pb.Company{
			Name: company.Name,
			Address: company.Address,
			Location: company.Location,
		}
		res, err := client.CreateCompany(ctx, &pb.CreateCompanyRequest{
			Company: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"company": res.Company,
		})
	})
	r.PUT("/companies/:id", func(ctx *gin.Context) {
		var company pb.Company

		err :=ctx.ShouldBind(&company)
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":err,
			})
			return
		}

		data := &pb.Company{
			Id: ctx.Param("id"),
			Name: company.Name,
			Address: company.Address,
			Location: company.Location,
		}
		res, err := client.UpdateCompany(ctx, &pb.UpdateCompanyRequest{
			Company:data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error" : err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"company":res.Company,
		})
	})
	r.DELETE("/companies/:id", func(ctx *gin.Context) {
		
		id := ctx.Param("id")
		
		res, err :=client.DeleteCompany(ctx, &pb.DeleteCompanyRequest{
			Id: id,
		})
		if err!= nil{
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":err,
			})
			return
		}
		if res.Success {
			ctx.JSON(http.StatusOK, gin.H{
				"message":"company deleted successfully",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message":"error while deleting company",
			})
		}

	})
// code block for adding temp entries in the db
	// for i:=1; i<10000;i++{
	// client.CreateCompany(&gin.Context{}, &pb.CreateCompanyRequest{
	// 	Company: &pb.Company{
	// 		Name: "test" + string(i),
	// 		Address: "temp",
	// 		Location: "temp",
	// 	},
	// })
	// }

	r.Run(":5000")
}

// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/movie.proto
// command to compile the proto file in golang from cli