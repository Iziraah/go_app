package tests

import "os"

const mainURL = "http://localhost:8081/api/v1/"

func SetupEnv() {
	os.Setenv("USERS_DB_PORT", "5434")
	os.Setenv("USERS_DB_NAME", "final")
	os.Setenv("USERS_GRPC_ADDR", ":52001")
	os.Setenv("LINKS_DB_PORT", "27018")
	os.Setenv("LINKS_GRPC_ADDR", ":51001")
	os.Setenv("LINKS_AMQP_PORT", "5674")
	os.Setenv("LINKS_AMQP_QNAME", "final")
	os.Setenv("APIGW_ADDR", ":8081")
	os.Setenv("APIGW_USERS_CLIENT_ADDR", ":52001")
	os.Setenv("APIGW_LINKS_CLIENT_ADDR", ":51001")
}
