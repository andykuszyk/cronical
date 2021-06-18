module github.com/andykuszyk/cronical

go 1.16

require (
	github.com/arran4/golang-ical v0.0.0-20210601225245-48fd351b08e7
	github.com/google/uuid v1.2.0
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20200116001909-b77594299b42 // indirect
)

replace github.com/arran4/golang-ical => github.com/andykuszyk/golang-ical v1.0.0
