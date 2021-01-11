# gin_admin

目的是做一个类似于django框架自带的admin后台管理，将注册的表在后台都可以方便管理，这是后端接口部分，前端有大佬感兴趣的话可以联系我

还在开发中，只实现了预想中的小部分功能

## 简单使用
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spiderman930706/gin_admin"
	"github.com/spiderman930706/gin_admin/config"
	"github.com/spiderman930706/gin_admin/models"
)

type User struct {
	models.User
	Phone string `gorm:"column:phone_num;type:varchar(15);unique_index" json:"phone"`
}

func main() {
	con := config.Config{
		Mysql: config.Mysql{
			DbName:   "blog2",
			User:     "root",
			Password: "930706",
			Host:     "127.0.0.1",
		},
		JWT: config.JWT{
			SigningKey:   "example-key",
			ExpireSecond: 7 * 24 * 3600,
		},
	}
	r := gin.New()

	//注册gin_admin
	group := r.Group("admin")
	if err := gin_admin.RegisterConfigAndRouter(con, group); err != nil {
		log.Println(err)
		os.Exit(0)
	}
	if err := gin_admin.RegisterTables(true, &User{}); err != nil {
		log.Println(err)
		os.Exit(0)
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8888),
		Handler:        r,
		ReadTimeout:    60,
		WriteTimeout:   60,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe() //.Error()
}

```
