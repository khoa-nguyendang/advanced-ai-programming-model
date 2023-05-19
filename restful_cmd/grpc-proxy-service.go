package main

import (
	"aapi/config"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	// "aapi/pkg/trace"

	gw "aapi/protos/gateway"
	"aapi/shared/constants"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// gorilla_mux "github.com/gorilla/mux"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

}

var (
	// command-line options:
	// gRPC server endpoint
	configfilepath     = flag.String("config-path", "", "The absolute path to configuration file")
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
	swaggerDir         = flag.String("swaggerDir", "/temp", "swagger api document folder")
)

func RunServer(address string, broker *Broker, opts ...runtime.ServeMuxOption) error {
	var configPath string
	if configfilepath != nil && *configfilepath != "" {
		configPath = *configfilepath
	} else {
		//load default
		configPath = config.GetConfigPath(os.Getenv("config"))
		log.Printf("Loadding default config path : %v \n", configPath)
	}
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := http.NewServeMux()

	handler := NewProxyHandler(dial(cfg.UserService), cfg)
	mux.HandleFunc("/swagger/", serveSwagger)
	mux.HandleFunc("/ping", handler.Ping)
	mux.HandleFunc("/server-publishing", broker.handleSSE)
	mux.HandleFunc("/enroll", handler.Enroll)
	mux.HandleFunc("/verify", handler.Verify)
	gw, err := newGateway(ctx, opts...)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)

	return http.ListenAndServe(address, allowCORS(mux))

}

func main() {
	log.Println("grpc Proxy Service is running now")
	flag.Parse()
	defer glog.Flush()
	broker := &Broker{
		Notifier:       make(chan []byte, 5),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()
	if err := RunServer(":8081", broker); err != nil {
		log.Fatal(err)
	}
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
		log.Printf("Swagger JSON not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	log.Printf("Serving %s", r.URL.Path)
	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join(*swaggerDir, p)
	http.ServeFile(w, r, p)
}

func newGateway(ctx context.Context, svropts ...runtime.ServeMuxOption) (http.Handler, error) {
	mux := runtime.NewServeMux(svropts...)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(constants.MAX_MESSAGE_LENGTH), grpc.MaxCallSendMsgSize(constants.MAX_MESSAGE_LENGTH))}
	var err error
	err = gw.RegisterLoggingHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Default().Println("RegisterLoggingHandlerFromEndpoint failed")
	}
	err = gw.RegisterUserHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Default().Println("RegisterUserHandlerFromEndpoint failed")
	}

	return mux, nil
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Request-Headers", "*")
	log.Printf("preflight request for %s", r.URL.Path)
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func dial(addr string) *grpc.ClientConn {
	log.Printf("addr: %#v", addr)

	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(constants.MAX_MESSAGE_LENGTH), grpc.MaxCallSendMsgSize(constants.MAX_MESSAGE_LENGTH))}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(fmt.Sprintf("ERROR: dial error: %v", err))
	}

	return conn
}
