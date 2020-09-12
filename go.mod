module ginweb

go 1.13

require (
	github.com/Shopify/sarama v1.27.0
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
	github.com/google/uuid v1.1.2 // indirect
	github.com/hpcloud/tail v1.0.0
	github.com/jinzhu/gorm v1.9.16
	github.com/spf13/viper v1.7.1
	go.etcd.io/etcd v3.3.25+incompatible
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	gopkg.in/fsnotify.v1 v1.0.0-00010101000000-000000000000 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace gopkg.in/fsnotify.v1 => gopkg.in/fsnotify.v1 v1.4.7
