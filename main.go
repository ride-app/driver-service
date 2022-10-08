package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	pb "github.com/ride-app/entity-service/api/gen/ride/entity/v1alpha1"
	"github.com/ride-app/entity-service/config"
	"github.com/ride-app/entity-service/di"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	log.SetReportCaller(true)
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		log.SetFormatter(&log.TextFormatter{
			DisableLevelTruncation: true,
			PadLevelText:           true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				dir, err := os.Getwd()
				if err != nil {
					dir = ""
				} else {
					dir = dir + "/"
				}

				filename := strings.Replace(f.File, dir, "", -1)

				return fmt.Sprintf("(%s)", path.Base(f.Function)), fmt.Sprintf(" %s:%d", filename, f.Line)

			},
		})
	}

	err := cleanenv.ReadConfig(".env", &config.Config)

	if err != nil {
		log.Warnf("Could not load config: %v", err)
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Config.Port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Infof("Listening to port: %d", config.Config.Port)

	var opts []grpc.ServerOption

	// if config.Config.Debug == false {
	// 	creds, err := credentials.NewServerTLSFromFile()
	// 	if err != nil {
	// 		log.Fatalf("Failed to generate credentials %v", err)
	// 	}
	// 	opts = []grpc.ServerOption{grpc.Creds(creds)}
	// }

	service, err := di.InitializeService()

	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}

	log.Info("Service Initialized")

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterEntityServiceServer(grpcServer, service)
	panic(grpcServer.Serve(lis))
}
