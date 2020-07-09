package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	r, authMiddleware := CreateRouter(viper.Get("log.gin").(string),
		PayloadFuncHandler,
		AuthenticatorHandler,
		IdentityHandler,
		AuthorizatorHandler)
	auth := r.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/listenaddr", ListListenaddrHandler)
	}
	//{
	//	// 监听地址
	//	auth.GET("/listenaddr", ListListenaddrHandler)
	//	auth.PUT("/listenaddr", vpdn.ModifyListenaddrHandler)
	//	auth.POST("/listenaddr", vpdn.AddListenaddrHandler)
	//	auth.DELETE("/listenaddr", vpdn.DelListenaddrHandler)
	//
	//	// 回连映射
	//	auth.GET("/proxymap", vpdn.ListProxyMapHandler)
	//	auth.PUT("/proxymap", vpdn.ModifyProxyMapHandler)
	//	auth.POST("/proxymap", vpdn.AddPorxyMapHandler)
	//	auth.DELETE("/proxymap", vpdn.DelProxyMapHandler)
	//
	//	// IP白名单
	//	auth.GET("/whitelist", vpdn.ListWhiteListHandler)
	//	auth.PUT("/whitelist", vpdn.ModifyWhiteListHandler)
	//	auth.POST("/whitelist", vpdn.AddWhiteListHandler)
	//	auth.DELETE("/whitelist", vpdn.DelWhiteListHandler)
	//	auth.PATCH("/whitelist", vpdn.ChangeWhiteStatusHandler)
	//
	//	// 连接管理
	//	auth.GET("/connection", vpdn.ListOnlineConnectionsHandler)
	//	auth.POST("/connection", vpdn.DisConnectionHandler)
	//
	//	// 日志
	//	auth.GET("/logs/connection/authorized", vpdn.ListConnectionsLogsHandler)
	//	auth.GET("/logs/connection/unauthorized", vpdn.ListAtkConnectionsHandler)
	//
	//	// 版本信息
	//	auth.GET("/version", vpdn.GetVersion)
	//}

	fmt.Println(authMiddleware)

	fmt.Println(r)
}
// 获取监听配置列表
func ListListenaddrHandler(c *gin.Context) {
	type Param struct {
		Page       int    `form:"page" json:"page" binding:"required,numeric,gt=0"`
		PageSize   int    `form:"pagesize" json:"pagesize" binding:"required,numeric,gt=0,lte=100"`
		Creator    string `form:"creator" json:"creator" binding:"omitempty"`
		ListenPort int    `form:"listenport" json:"listenport" binding:"omitempty,numeric,gt=0,lt=65536"`
		IP         string `form:"ip" json:"ip" binding:"omitempty,ipv4"`
		Port       int    `form:"port" json:"port" binding:"omitempty,numeric,gt=0,lt=65536"`
		Sdate      int    `form:"sdate" json:"sdate" binding:"omitempty,numeric,ltfield=Edate"`
		Edate      int    `form:"edate" json:"edate" binding:"omitempty,numeric,gtfield=Sdate"`
	}
	var param Param
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status_code": 400, "reason": fmt.Sprint(err.Error())})
		return // exit on first error
	}

	response := listenAddrResponse{}
	response.Page = param.Page
	response.Pagesize = param.PageSize
	response.Total = listenaddrs.Count(param.IP, param.Creator, param.Sdate, param.Edate)
	response.Data = listenaddrs.Find(param.IP, param.Creator, param.Sdate, param.Edate, param.Page, param.PageSize)
	c.IndentedJSON(http.StatusOK, response)
}
// User 用户信息
type ListenAddr struct {
	ID         int64  `json:"id"`
	IP         string `json:"ip"`
	Creator    string `json:"creator"`
	CreateTime uint32 `json:"timestamp"`
}
type listenAddrResponse struct {
	Page     int          `json:"page"`
	Pagesize int          `json:"pagesize"`
	Total    int          `json:"total"`
	Data     []ListenAddr `json:"data"`
}

var listenaddrs *ListenAddrs

// Users 用户管理
type ListenAddrs struct {
	Db             *sql.DB
	ListenAddrInfo map[string](*ListenAddr)
}

func (la *ListenAddrs) Count(ip string, creator string, sdate int, edate int) int {
	ipQueryCondition := ""
	andstr := ""
	if ip != "" {
		ipQueryCondition = fmt.Sprintf(" ip='%s'", ip)
		andstr = " AND"
	}
	creatorQueryCondition := ""
	if creator != "" {
		creatorQueryCondition = fmt.Sprintf("%s creator='%s'", andstr, creator)
		andstr = " AND"
	}
	dateQueryCondition := ""
	if sdate != edate {
		dateQueryCondition = fmt.Sprintf("%s create_time BETWEEN '%d' AND '%d'", andstr, sdate, edate)
	}
	conditions := ""
	if ip != "" || creator != "" || sdate != edate {
		conditions = fmt.Sprintf(" where%s%s%s", ipQueryCondition, creatorQueryCondition, dateQueryCondition)
	}
	query := fmt.Sprintf("select count(*) from listenip%s;", conditions)
	rows, err := la.Db.Query(query)
	if err != nil {
		log.Printf(query, err.Error())
		return 0
	}

	defer rows.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Printf(err.Error())
			return 0
		}
	}

	return count
}


// 查询条目
func (la *ListenAddrs) Find(ip string, creator string, sdate int, edate int, page int, pagesize int) []ListenAddr {
	ipQueryCondition := ""
	andstr := ""
	if ip != "" {
		ipQueryCondition = fmt.Sprintf(" ip='%s'", ip)
		andstr = " AND"
	}
	creatorQueryCondition := ""
	if creator != "" {
		creatorQueryCondition = fmt.Sprintf("%s creator=='%s'", andstr, creator)
		andstr = " AND"
	}
	dateQueryCondition := ""
	if sdate != edate {
		dateQueryCondition = fmt.Sprintf("%s create_time BETWEEN '%d' AND '%d'", andstr, sdate, edate)
	}
	conditions := ""
	if ip != "" || creator != "" || sdate != edate {
		conditions = fmt.Sprintf(" where%s%s%s", ipQueryCondition, creatorQueryCondition, dateQueryCondition)
	}
	query := fmt.Sprintf("select * from listenip%s limit %d offset %d;",
		conditions,
		pagesize, (page-1)*pagesize)
	rows, err := la.Db.Query(query)
	if err != nil {
		log.Printf(query, err.Error())
		return nil
	}

	defer rows.Close()

	listenaddrs := make([]ListenAddr, 0, pagesize)
	listenaddr := ListenAddr{}
	for rows.Next() {
		err = rows.Scan(&listenaddr.ID, &listenaddr.IP, &listenaddr.Creator, &listenaddr.CreateTime)
		if err != nil {
			log.Printf(err.Error())
			return nil
		}
		listenaddrs = append(listenaddrs, listenaddr)
	}

	return listenaddrs
}