package dao

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
)

//UserInfo struct
type UserInfo struct {
	Userid    int
	Username  string
	Departmen string
}

// GetUsername get username
func GetUsername(ctx context.Context, user string) (us UserInfo, value bool) {

	span, _ := opentracing.StartSpanFromContext(ctx, "GetUsername")
	defer span.Finish()

	query := "SELECT * FROM userinfo WHERE username=?"

	rows, err := Connect.Query(query, user)
	if err != nil {
		span.LogFields(
			otlog.String("Infor: ", "Can not get information of user!"),
			otlog.String("Query", query),
			otlog.String("Value", user),
		)
		span.SetTag("error", true)
		span.SetTag("db.statement", query)
		span.SetTag("db.type", "sql")
	} else {
		span.LogFields(
			otlog.String("Infor: ", "Get all information of user!"),
			otlog.String("Query", query),
			otlog.String("Value", user),
		)
	}

	defer rows.Close()

	if rows.Next() == false {
		span.LogFields(
			otlog.String("Infor: ", user+" "+"does not exist!"),
		)
		span.SetTag("error", true)
		span.SetTag("db.type", "sql")

		value = false
		return us, value
	}
	rows.Scan(&us.Userid, &us.Username, &us.Departmen)
	value = true
	return us, value

}
