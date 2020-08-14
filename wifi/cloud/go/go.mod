module magma/wifi/cloud/go

replace (
	magma/gateway => ./../../../orc8r/gateway/go
	magma/orc8r/cloud/go => ./../../../orc8r/cloud/go
	magma/orc8r/lib/go => ./../../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ./../../../orc8r/lib/go/protos
)

require (
	github.com/Masterminds/squirrel v1.1.1-0.20190513200039-d13326f0be73
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.18.0
	github.com/go-openapi/validate v0.18.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.3
	github.com/google/uuid v1.1.1
	github.com/labstack/echo v0.0.0-20181123063414-c54d9e8eed6c
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	github.com/thoas/go-funk v0.7.0
	magma/orc8r/cloud/go v0.0.0
	magma/orc8r/lib/go v0.0.0
	magma/orc8r/lib/go/protos v0.0.0
)

go 1.12
