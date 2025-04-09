package grpctool

import (
	"google.golang.org/grpc"
	"testing"
)

func TestRun(t *testing.T) {
	err := Run(
		ServerConfig{Port: 8080}, func(s *grpc.Server) {

		},
	)
	if err != nil {
		t.Fatal(err)
		return
	}
}
