module github.com/omec-project/pcf

go 1.16

replace github.com/omec-project/nrf => ../nrf

replace github.com/omec-project/MongoDBLibrary => ../mongodblibrary

require (
	github.com/antihax/optional v1.0.0
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/cydev/zero v0.0.0-20160322155811-4a4535dd56e7
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.3
	github.com/google/uuid v1.3.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/omec-project/MongoDBLibrary v1.1.3
	github.com/omec-project/config5g v1.2.0
	github.com/omec-project/http2_util v1.1.0
	github.com/omec-project/http_wrapper v1.1.0
	github.com/omec-project/idgenerator v1.1.0
	github.com/omec-project/logger_conf v1.1.0
	github.com/omec-project/logger_util v1.1.0
	github.com/omec-project/nrf v1.0.1
	github.com/omec-project/openapi v1.1.0
	github.com/omec-project/path_util v1.1.0
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.0
	github.com/tim-ywliu/event v0.1.0
	github.com/urfave/cli v1.22.5
	gopkg.in/yaml.v2 v2.4.0
)
