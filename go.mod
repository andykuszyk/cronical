module github.com/andykuszyk/cronical

go 1.16

require (
	github.com/arran4/golang-ical v0.0.0-20210601225245-48fd351b08e7
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/google/uuid v1.2.0
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/arran4/golang-ical => github.com/andykuszyk/golang-ical v1.0.0
