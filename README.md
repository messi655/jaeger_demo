# What is Jaeger: 
- https://www.jaegertracing.io/

# Setup Jaeger backend for development environment

```
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest
```


You can then navigate to http://localhost:16686 to access the Jaeger UI.

# How to run this demo
## Prepare database
- Make sure you have MySQL on your machine
- Create database name: `jeager`
- Create table `userinfo` with columns:

```
uid int(10) PK
username varchar(64) 
departname varchar(64) 
```

## How to run 
- Make sure Go 1.9 or higher installed on your machine.
- Clone the source code into $GOPATH/src/trustingsocial.com/jeager
- run:
```
go run selectdb/main.go
go run client/main.go
```
- Open browse and access to:
```
http://ip_add:3000/users/user=tinhuynh
```
```
http://ip_add:4000/clients/client=tinhuynh
``` 

- Access to Jaeger UI and select service `Jaeger Demo Backend API`

- Access to Jaeger UI and select service `Jaeger Demo Client`

# How to send tracing data to Jaeger backend 

## Create an instance of Jaeger Tracer 
```
func InitJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: "Jeager Demo",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 4,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "IP_Address_Jaeger_Agent:6831",
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
```

## Send tracing data to Span (Apply this into logical functions)
### Declare span context
```
    span, _ := opentracing.StartSpanFromContext(ctx, "Name of logical function")
	defer span.Finish()
```
#### Send tracing data to span, this data will show on Jaeger UI
```
    span.LogFields(
		log.String("Infor", "Can not connect to DB"),
		log.String("value", datasource),
	)
	span.SetTag("error", true)
```

## Init Jaeger (Apply this in the main function)
```
    tracer, closer := initJaeger("Name Of API or Services")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
```

