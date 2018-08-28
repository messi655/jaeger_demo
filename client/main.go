package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jlog "github.com/opentracing/opentracing-go/log"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/http"
	"trustingsocial.com/jeager/jaeger_demo/tracing"
)

func selectdb(ctx context.Context, userremoteurl string, username string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "selectDB")
	defer span.Finish()

	v := url.Values{}
	v.Set("user", username)
	remoteurl := userremoteurl + v.Encode()
	req, err := http.NewRequest("GET", remoteurl, nil)
	if err != nil {
		panic(err.Error())
	}

	////////////

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, remoteurl)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	resp, err := xhttp.Do(req)
	if err != nil {
		span.LogFields(
			jlog.String("Infor", "Can not retrive remote url"),
			jlog.String("value", remoteurl),
		)
		span.SetTag("error", true)
	} else {
		span.LogFields(
			jlog.String("Infor", "Connected remote url"),
			jlog.String("value", remoteurl),
		)
	}

	helloStr := string(resp)

	return helloStr

}

func main() {

	//"http://localhost:3000/users?"

	tracer, closer := tracing.InitJaeger("Jaeger Demo Client")
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	userUrls := "http://localhost:3000/users?"

	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("Show User Infor", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		ctx := opentracing.ContextWithSpan(context.Background(), span)

		user := r.FormValue("client")
		helloTo := selectdb(ctx, userUrls, user)

		helloStr := fmt.Sprintf("Hello, %s!", helloTo)

		span.LogFields(
			jlog.String("url", userUrls),
			jlog.String("value", helloStr),
		)

		w.Write([]byte(helloStr))

	})
	log.Fatal(http.ListenAndServe(":4000", nil))

}
