package dao

import (
	"context"
	"database/sql"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
)

// Connect export for all
var Connect *sql.DB
var err error

// ConnectDB connect to database
func ConnectDB(ctx context.Context, username string, password string, address string, dbport string, dbname string) {

	span, _ := opentracing.StartSpanFromContext(ctx, "Connect Database")
	defer span.Finish()

	datasource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, address, dbport, dbname)
	Connect, err = sql.Open("mysql", datasource)

	if err != nil {
		span.LogFields(
			otlog.String("Infor", "Can not connect to DB"),
			otlog.String("value", datasource),
		)
		span.SetTag("error", true)
		span.SetTag("db.instance", dbname)
	}
	err = Connect.Ping()
	if err != nil {
		span.LogFields(
			otlog.String("Infor", "error! Something wrong in your DB connection information!"),
			otlog.String("value", datasource),
		)
		span.SetTag("error", true)
		span.SetTag("db.instance", dbname)
	} else {
		span.LogFields(
			otlog.String("Infor", "Connected to Database!"),
			otlog.String("value", datasource),
		)
	}

}
