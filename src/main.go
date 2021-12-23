package main

import (
	"fmt"
	"log"
	"net"
	"path"
	"path/filepath"
	"runtime"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/centerpb"
	"swd_project/src/pbs/userpb"
	"swd_project/src/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func main() {
	fmt.Println("Server is running...")

	postgresdb.PostgresInit()
	if err := postgresdb.DB.AutoMigrate(&model.User{}, &model.Consultant{}, &model.Report{}, &model.Schudule{}, &model.Reserve{}, &model.Center{}); err != nil {
		log.Fatalf("error while AutoMigrate : %v", err)
	}

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("error while listening to host: %v\n", err)
	}

	s := grpc.NewServer()

	reflection.Register(s)

	serviceRegistry(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error while serving the server : %v", err)
	}
}

func serviceRegistry(s *grpc.Server) {
	userpb.RegisterUserServiceServer(s, &server.UserServer{})
	centerpb.RegisterCenterServiceServer(s, &server.CenterServer{})
}
