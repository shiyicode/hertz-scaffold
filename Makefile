module = github.com/three-body/hertz-scaffold

router_dir = biz/router
handler_dir = biz/handler
model_dir = biz/hmodel
client_dir = client

gorm_dir = biz/dal/query

gen_api_init:
	hz new -mod $(module) -idl api/book.thrift --thrift-plugins=validator --handler_dir $(handler_dir) --model_dir $(model_dir) --client_dir $(client_dir) --router_dir $(router_dir)
	hz new -mod $(module) -idl api/user.thrift --thrift-plugins=validator --handler_dir $(handler_dir) --model_dir $(model_dir) --client_dir $(client_dir) --router_dir $(router_dir)

gen_api:
	hz update -idl api/book.thrift --thrift-plugins=validator --handler_dir $(handler_dir) --model_dir $(model_dir) --client_dir $(client_dir)
	hz update -idl api/user.thrift --thrift-plugins=validator --handler_dir $(handler_dir) --model_dir $(model_dir) --client_dir $(client_dir)

gen_swagger:
	swag init --parseDependency --parseInternal --parseDepth 5 --generalInfo cmd/api.go; swag fmt

gen_gorm:
	rm -rf $(gorm_dir)
	go run tools/gormgen/* -o $(gorm_dir)

