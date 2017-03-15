package main

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/astaxie/beego/orm"
	"../models"
)

type ReturnMsg struct {
	Msg  string
	Is   int
	Data interface{}
}

type Address struct {
	Id       int
	CityCode string
	Add      string
}
type Login struct {
	Name string
	Age  int
	Add  []*Address
}

func test_json() {
	a := new(Address)
	a.Add = "`AA"
	a.CityCode = "020"
	a.Id = 1

	b := &Address{2, "020", "BB"}
	c := Address{3, "050", "CC"}

	l := &Login{"ADB", 20, make([]*Address, 10)}
	l.Add [0] = a
	l.Add [1] = b
	l.Add [2] = &c

	r := &ReturnMsg{"ok", 1, l}
	js, _ := json.Marshal(r)
	fmt.Println("JSON format: %s", js)

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

}

func init() {
	models.InitTable()
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("default", false, true)
	models.AutoCreateData()
}
func main() {
	o := orm.NewOrm()
	var p []models.Post
	o.QueryTable("post").All(&p)
	for k, v := range p {
		o.QueryTable("PostImage").Filter("Belong", v).All(&p[k].Images)
		o.QueryTable("PostComment").Filter("Belong", v).All(&p[k].Comments)
		o.QueryTable("PostFavort").Filter("Belong", v).All(&p[k].Favorts)
	}
	for _, v := range p {
		fmt.Println(v, "\n", &v.Comments, "\n", &v.Images, "\n", &v.Favorts)
	}

}
