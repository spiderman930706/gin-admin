package gin_admin

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spiderman930706/gin_admin/config"
	"github.com/spiderman930706/gin_admin/core"
	"github.com/spiderman930706/gin_admin/global"
	"github.com/spiderman930706/gin_admin/middleware"
	"github.com/spiderman930706/gin_admin/models"
	"github.com/spiderman930706/gin_admin/routers"
	"gorm.io/gorm/schema"
)

func RegisterConfigAndRouter(config config.Config, Router *gin.RouterGroup) error {
	global.Config = config
	db, err := core.MysqlInit(global.Config.Mysql)
	if err != nil {
		return err
	}
	global.DB = db
	//Router.Use(middleware.JWTAuth())
	Router.Use(middleware.Cors()).Use(middleware.Recovery())
	routers.InitRouter(Router)
	return nil
}

func RegisterTables(migrate bool, dst ...models.AdminOperation) error {
	if migrate {
		var res []interface{}
		for _, n := range dst {
			res = append(res, n)
		}
		if err := core.MigrateMysqlTables(global.DB, res...); err != nil {
			return err
		}
	}
	global.Tables = make(map[string]*global.Table)
	for _, n := range dst {
		if err := ParseSchema(n); err != nil {
			return err
		}
	}
	return nil
}

//获取表信息
func ParseSchema(n models.AdminOperation) error {
	model := global.DB.Model(n)
	if err := model.Statement.Parse(n); err != nil {
		log.Printf("MySQL异常 %s", err)
		return err
	}
	tableName := model.Statement.Table
	//先用表名和字段名来做数据的增删改，但在这里要取出所有定义了admin的参数，比如admin:"list:id;type:int"
	schemaField := model.Statement.Schema.FieldsByName
	newFields := make(map[string]*global.Field)
	for _, v := range schemaField {
		field := ParseTag(v)
		newFields[v.DBName] = field
	}
	global.Tables[tableName] = &global.Table{
		Field:     newFields,
		Source:    n,
		CanDelete: n.CanDelete(),
		CanModify: n.CanModify(),
		CanAdd:    n.CanAdd(),
	}
	return nil
}

//获取struct tag中admin的信息
func ParseTag(v *schema.Field) (m *global.Field) {
	tag := v.Tag
	admin := tag.Get("admin")
	arr := strings.Split(admin, ";")
	m = &global.Field{}
	m.Schema = v
	for _, n := range arr {
		rr := strings.Split(n, ":")
		switch rr[0] {
		case "list":
			m.List = true
		case "type":
			m.Type = rr[1]
		case "name":
			m.Name = rr[1]
		}
	}
	return
}
