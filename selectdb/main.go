package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/config"
	"trustingsocial.com/jeager/jaeger_demo/dao"
	"trustingsocial.com/jeager/jaeger_demo/tracing"
)

func userInfor(ctx context.Context, username string) (resultstring string, value bool) {

	//var resultstring string

	uss, value := dao.GetUsername(ctx, username)
	span, _ := opentracing.StartSpanFromContext(ctx, "Return User Information")
	defer span.Finish()

	if value == false {
		span.LogFields(
			otlog.String("Infor: ", username+" "+"does not exist!"),
		)
	} else {
		resultstring = strconv.Itoa(uss.Userid) + "," + uss.Username + "," + uss.Departmen
		span.LogFields(
			otlog.String("value: ", username),
			otlog.String("User Information: ", resultstring),
		)
	}
	en, _ := config.FromEnv()
	fmt.Printf("%v\n", en.Reporter.LocalAgentHostPort)

	return resultstring, value
}

func main() {
	tracer, closer := tracing.InitJaeger("Jaeger Demo Backend API")
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("Show User Infor", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		ctx := opentracing.ContextWithSpan(context.Background(), span)
		dao.ConnectDB(ctx, "root", "123abc", "localhost", "3306", "jeager")
		user := r.FormValue("user")
		helloTo, result := userInfor(ctx, user)

		var helloStr string

		if result == false {
			helloStr = fmt.Sprintf("%s does not exist!", user)
		} else {
			helloStr = fmt.Sprintf("Hello, %s!", helloTo)
			span.LogFields(
				otlog.String("event", "User Information!"),
				otlog.String("value", helloStr),
			)
		}

		w.Write([]byte(helloStr))

	})
	log.Fatal(http.ListenAndServe(":3000", nil))
	///=====================
	/////////////

}
