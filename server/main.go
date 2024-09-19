package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	// "time"

	// "net/url"

	// "os"

	database "example.com/go-grpc-crud-api/database"
	// "example.com/go-grpc-crud-api/database"
	pb "example.com/go-grpc-crud-api/proto"

	// "github.com/google/uuid"
	"github.com/oklog/ulid/v2"

	// "github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	// "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
// var err error

// const (
	// port = ":50051"
	// dbURL = "postgres://user:password@localhost:5432/companydb"
//   )

func init() {
	godotenv.Load()
	// DatabaseConnection()
	DB = database.DatabaseConnection()
}

type Company struct {
	ID			string 		`gorm:"primarykey"`
	Name		string 
	Address		string 
	Location	string
}



// creating grpc server
var (
	port = flag.Int("port", 50051, "gRPC server port")
)

type server struct {
	pb.UnimplementedCompanyServiceServer
	// db *pgx.Conn
}

// implpement rpc methods

func (*server) CreateCompany(ctx context.Context, req *pb.CreateCompanyRequest) (*pb.CreateCompanyResponse, error) {
	fmt.Println("Create Company")
	company := req.GetCompany()
	// company.Id = uuid.New().String()
	company.Id = ulid.Make().String()
	
	

	data := Company{
		ID: company.GetId(),
		Name: company.GetName(),
		Address: company.GetAddress(),
		Location: company.GetLocation(),
	}

	res := DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("company creation unsuccessful")
	}
	return &pb.CreateCompanyResponse{
		Company: &pb.Company{
			Id: company.GetId(),
			Name: company.GetName(),
			Address: company.GetAddress(),
			Location: company.GetLocation(),

		},
	}, nil

}

func (*server) GetCompany(ctx context.Context, req *pb.ReadCompanyRequest) (*pb.ReadCompanyResponse, error) {
	fmt.Println("Read Company", req.GetId())
	var company Company
	res := DB.Find(&company, "id = ?", req.GetId())
	if res.RowsAffected == 0 {
		return nil, errors.New("company not found")
	}
	return &pb.ReadCompanyResponse{
		Company: &pb.Company{
			Id: company.ID,
			Name: company.Name,
			Address: company.Address,
			Location: company.Location,
		},
	}, nil
}

const (
	defaultPageSize = 5
	maxPageSize = 20
)

func (*server) GetCompanies(ctx context.Context, req *pb.ReadCompaniesRequest) (*pb.ReadCompaniesResponse, error) {
	fmt.Println("Read all companies")
	companies := []*pb.Company{}
	
	// fmt.Println("start time server" + time.Now().String())
	// res := DB.Find(&companies)
	// fmt.Println("end time server" + time.Now().String())
	req.PageSize = defaultPageSize

	fmt.Println(req.GetPageToken())
	query := DB.Table("companies").Limit(int(req.GetPageSize())).Order("id ASC")

	if len(req.GetPageToken()) != 0 {
		query = query.Where("id > ?", req.GetPageToken())
	}

	query.Scan(&companies)

	// if res.RowsAffected == 0
	if len(companies) == 0 {
		return nil, errors.New("no company found")
	}

	lastIndex := len(companies) - 1
	nextPageToken := companies[lastIndex].Id

	if len(companies) < int(req.GetPageSize()) {
		nextPageToken = ""
	}

	
	return &pb.ReadCompaniesResponse{
		Companies: companies,
		NextPageToken: nextPageToken,
	}, nil
	
}

func (*server) UpdateCompany(ctx context.Context, req *pb.UpdateCompanyRequest) (*pb.UpdateCompanyResponse, error) {
	fmt.Println("Update Company")
	var company Company
	reqCompany := req.GetCompany()

	res := DB.Model(&company).Where("id=?", reqCompany.Id).Updates(Company{Name : reqCompany.Name, Address: reqCompany.Address, Location: reqCompany.Location})

	if res.RowsAffected == 0 {
		return nil, errors.New("company not found")
	}

	return &pb.UpdateCompanyResponse{
		Company: &pb.Company{
			Id: company.ID,
			Name: company.Name,
			Address: company.Address,
			Location: company.Location,
		},
	}, nil
}

func (*server) DeleteCompany(ctx context.Context, req *pb.DeleteCompanyRequest) (*pb.DeleteCompanyResponse, error){
	fmt.Println("Delete company")
	var company Company
	res := DB.Where("id = ?", req.GetId()).Delete(&company)

	if res.RowsAffected == 0 {
		return nil, errors.New("company not found")
	}

	return &pb.DeleteCompanyResponse{
		Success: true,
	}, nil
}

// starting the grpc server

func main() {
	fmt.Println("gRPC server running...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	// docker start
	// db, err := pgx.Connect(pgx.ConnConfig{Host: "localhost", Port: 5432, User: "postgres", Password: "Paras@435"})
	// if err != nil {
	//   log.Fatalf("failed to connect to database: %v", err)
	// }
	// defer db.Close()
	// docker end

	s := grpc.NewServer()

	pb.RegisterCompanyServiceServer(s, &server{})
	// pb.RegisterCompanyServiceServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}

	// for i:=0; i < 10000; i++{
	// 	// call create company function
	// 	// with name, address, location 
	// 	CreateCompany(gin.Context, pb.CreateCompanyRequest{
	// 		Company: &pb.Company{
	// 			Name: "test" + string(i),
	// 			Address: "temp",
	// 			Location: "temp",
	// 		},
	// 	});
		
		
	// }
	
	// run this function once and load the data in the db



}