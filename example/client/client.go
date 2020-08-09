// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"encoding/json"
	rk_context "github.com/rookie-ninja/rk-interceptor/context"
	"github.com/rookie-ninja/rk-interceptor/example/proto"
	rk_logging_zap "github.com/rookie-ninja/rk-interceptor/logging/zap"
	rk_retry "github.com/rookie-ninja/rk-interceptor/retry"
	rk_logger "github.com/rookie-ninja/rk-logger"
	rk_query "github.com/rookie-ninja/rk-query"
	"google.golang.org/grpc"
	"log"
	"time"
)

var (
	bytes = []byte(`{
     "level": "info",
     "encoding": "console",
     "outputPaths": ["stdout"],
     "errorOutputPaths": ["stderr"],
     "initialFields": {},
     "encoderConfig": {
       "messageKey": "msg",
       "levelKey": "",
       "nameKey": "",
       "timeKey": "",
       "callerKey": "",
       "stacktraceKey": "",
       "callstackKey": "",
       "errorKey": "",
       "timeEncoder": "iso8601",
       "fileKey": "",
       "levelEncoder": "capital",
       "durationEncoder": "second",
       "callerEncoder": "full",
       "nameEncoder": "full"
     },
    "maxsize": 1,
    "maxage": 7,
    "maxbackups": 3,
    "localtime": true,
    "compress": true
   }`)

	logger, _, _ = rk_logger.NewZapLoggerWithBytes(bytes, rk_logger.JSON)
)

func main() {
	// create event factory
	factory := rk_query.NewEventFactory(
		rk_query.WithAppName("app"),
		rk_query.WithLogger(logger),
		rk_query.WithFormat(rk_query.RK))

	// create client interceptor
	opt := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			rk_logging_zap.UnaryClientInterceptor(factory),
			rk_retry.UnaryClientInterceptor()),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

	// Set up a connection to the server.
	conn, err := grpc.DialContext(context.Background(), "localhost:8080", opt...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create grpc client
	c := proto.NewGreeterClient(conn)
	// create with rk context
	ctx, cancel := context.WithTimeout(rk_context.NewContext(), 5*time.Second)
	defer cancel()

	// add metadata
	rk_context.AddToOutgoingMD(ctx, "key", "1", "2")
	// add request id
	rk_context.AddRequestIdToOutgoingMD(ctx)

	// call server
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "name"})

	// print incoming metadata
	bytes, _ := json.Marshal(rk_context.GetIncomingMD(ctx))
	println(string(bytes))

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
