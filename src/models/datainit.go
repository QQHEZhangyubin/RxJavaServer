package models

import (
	"math/rand"
	"github.com/astaxie/beego/orm"
	"time"
	"strconv"
	"github.com/astaxie/beego"
	"fmt"
)

var IMGAGES []string
var CONTENTS = [...]string{"",
	"beego 是一个快速开发 Go 应用的 HTTP 框架，他可以用来快速开发 API、Web 及后端服务等各种应用，是一个 RESTful 的框架，主要设计灵感来源于 tornado、sinatra 和 flask 这三个框架，但是结合了 Go 本身的一些特性（interface、struct 嵌入等）而设计的一个框架。",
	"哈哈，18123456789,ChinaAr  http://www.ChinaAr.com;一个不错的VR网站。哈哈，ChinaAr  http://www.ChinaAr.com;一个不错的VR网站。哈哈，ChinaAr  http://www.ChinaAr.com;一个不错的VR网站。哈哈，ChinaAr  http://www.ChinaAr.com;一个不错的VR网站。",
	"既然 beego 是基于这些模块构建的，那么它的执行逻辑是怎么样的呢？beego 是一个典型的 MVC 架构，它的执行逻辑如下图所示：",
	"beego 的架构",
	"使用beego实现的"}

func createUser() []*User {

	users := make([]*User, 9)
	users[0] = &User{Name:"FORME", Mobile:"13800138000", Sex:false, Age:20, Birth:"2017-3-3", Address: "广东天河区", Pwd:"123456", Pic:IMGAGES[0]}
	users[1] = &User{Name:"NOKIA", Mobile:"13800138001", Sex:false, Age:20, Birth:"2017-3-3", Address:"广东天河区", Pwd:"123456", Pic:IMGAGES[0]}
	users[2] = &User{Name:"TCL", Mobile:"13800138002", Sex:false, Age:20, Birth:"2017-3-3", Address:"广东天河区", Pwd:"123456", Pic:IMGAGES[1]}
	users[3] = &User{Name:"HUAWEI", Mobile:"13800138003", Sex:false, Age:20, Birth:"2017-3-3", Address: "广东天河区", Pwd:"123456", Pic:IMGAGES[2]}
	users[4] = &User{Name:"APPLE", Mobile:"13800138004", Sex:false, Age:20, Birth:"2017-3-3", Address: "广东天河区", Pwd:"123456", Pic: IMGAGES[3]}
	users[5] = &User{Name:"LG", Mobile:"13800138005", Sex:false, Age:20, Birth:"2017-3-3", Address: "广东天河区", Pwd:"123456", Pic:IMGAGES[4]}
	users[6] = &User{Name:"MEIZU", Mobile:"13800138006", Sex:false, Age:20, Birth:"2017-3-3", Address:"广东天河区", Pwd:"123456", Pic:IMGAGES[5]}
	users[7] = &User{Name:"XIAOMI", Mobile:"13800138007", Sex:false, Age:20, Birth:"2017-3-3", Address: "广东天河区", Pwd:"123456", Pic: IMGAGES[6]}
	users[8] = &User{Name:"这个名字是不是很长，哈哈！因为我是用来测试换行的", Mobile:"13800138008", Sex:false, Age:20, Birth:"2017-3-3", Address: "广东天河区", Pwd:"123456", Pic: IMGAGES[6]}
	return users
}

func createImage() []*PostImage {

	ilen := rand.Intn(9)

	var posts = make([]*PostImage, ilen)

	IMGLEN := len(IMGAGES)
	for i := 0; i < ilen; i++ {
		intn := rand.Intn(IMGLEN)
		posts[i] = &PostImage{Name:"img", Url:IMGAGES[intn], Size:""}
	}
	return posts
}

func getContent() string {
	return CONTENTS[rand.Intn(len(CONTENTS))]

}

func createFavort(userarr []*User, p *Post) {
	o := orm.NewOrm()
	o.Begin()
	l := len(userarr)
	left := rand.Intn(l)
	right := rand.Intn(l)
	if right < left {
		temp := right
		right = left
		left = temp

	}

	for i := left; i < right; i++ {
		f := &PostFavort{User:userarr[i], Belong:p}
		o.Insert(f)
	}
	o.Commit()

}

func createComment(userarr []*User, p *Post) {
	o := orm.NewOrm()
	o.Begin()
	//time := time.Now().Format("2006-01-02 15:04:05")
	l := len(userarr)
	a := rand.Intn(l)
	for i := 0; i < a; i++ {
		var s string = p.CreateTime + "评价测试:" + strconv.Itoa(i)
		c := &PostComment{Content:s, Belong:p}
		user := userarr[rand.Intn(l)]
		c.User = user
		c.Type = 0
		if rand.Intn(100) % 2 == 0 {
			c.Type = 1
			c.ToReplayUser = userarr[rand.Intn(l)]
		} else {
			c.ToReplayUser = p.Author
		}
		cou, err := o.Insert(c)
		beego.Info("插入评价=", cou, "error:", err)

	}
	o.Commit()

}

func createPost() {

	o := orm.NewOrm()
	//All 的参数支持 *[]Type 和 *[]*Type 两种形式的 slice
	var users []*User
	o.QueryTable("user").All(&users)

	o.Begin()
	ulen := len(users)
	time := time.Now().Format("2006-01-02 15:04:05")
	for i := 0; i < 15; i++ {
		randnum := rand.Intn(ulen)
		user := users[randnum]
		p := new(Post)
		p.Author = user
		p.Content = getContent()
		p.CreateTime = time
		p.Type = 1
		o.Insert(p)
		images := createImage()
		for k, _ := range images {
			images[k].Belong = p
			o.Insert(images[k])
		}
	}
	o.Commit()

	var posts []*Post
	o.QueryTable("post").RelatedSel().All(&posts)
	for _, v := range posts {
		createFavort(users, v)
		createComment(users, v)
	}

}
func AutoCreateData() {
	IMGAGES = make([]string, 20)
	for i := 0; i < 20; i++ {
		file := fmt.Sprint("static/img/", i + 1, ".jpg")
		IMGAGES[i] = file
	}

	o := orm.NewOrm()
	c, _ := o.QueryTable("user").Count()
	if c == 0 {
		o.Begin()
		userArr := createUser()
		for _, v := range userArr {
			o.Insert(v)
		}
		o.Commit()
		createPost()
	}

}
