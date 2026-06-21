package main

import (
	"context"
	"log"
	"time"

	"example.com/pz2-grpc/gen/studentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.NewClient(
		"127.0.0.1:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := studentpb.NewStudentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Ping
	pingResp, err := client.Ping(ctx, &studentpb.PingRequest{Message: "hello grpc"})
	if err != nil {
		log.Fatal("Ping error:", err)
	}
	log.Println("Ping response:", pingResp.GetMessage())

	// GetStudentByID — существующий
	studentResp, err := client.GetStudentByID(ctx, &studentpb.GetStudentRequest{Id: 1})
	if err != nil {
		log.Fatal("GetStudentByID error:", err)
	}
	st := studentResp.GetStudent()
	log.Printf("Student: id=%d, full_name=%s, group=%s, email=%s",
		st.GetId(), st.GetFullName(), st.GetGroup(), st.GetEmail())

	// GetStudentByID — несуществующий
	_, err = client.GetStudentByID(ctx, &studentpb.GetStudentRequest{Id: 999})
	if err != nil {
		log.Println("GetStudentByID(999) error:", err)
	}

	// ListStudents
	listResp, err := client.ListStudents(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatal("ListStudents error:", err)
	}
	log.Printf("ListStudents: total=%d", len(listResp.GetStudents()))
	for _, s := range listResp.GetStudents() {
		log.Printf("  - id=%d, full_name=%s, group=%s", s.GetId(), s.GetFullName(), s.GetGroup())
	}
}
