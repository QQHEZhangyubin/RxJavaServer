package models

import (
	"../controllers/been"
	_"github.com/mattn/go-sqlite3"
	"github.com/astaxie/beego/orm"
	"github.com/Unknwon/com"
	"os"
	"path"
	"github.com/astaxie/beego"
	"errors"
	"time"
	"fmt"
	"../utils"
)

// 数据库建表总models
//
const (
	_DB_NAME = "data/beeblog.db"
	_SQLITE3_DRIVER = "sqlite3"
)

/**
用户
 */
type User struct {
	Id           int64                          //int, int32 - 设置 auto 或者名称为 Id 时 integer AUTO_INCREMENT
	Name         string `orm:"size(36)"`        //设置不为空并不能大小24个字符
	Mobile       string `orm:"unique;size(32)"` //手机
	Sex          bool   `orm:"default(false)"`  //性别
	Age          int    `orm:"default(0)"`
	Birth        string  `orm:"default(2016-10-10)"`
	Address      string `orm:"default(广州市天河区)"`
	Pwd          string `orm:"size(128)"`       //设置不为空并不能大小64个字符,数据库表默认为 NOT NULL，设置 null 代表 ALLOW NULL
	Pic          string  `orm:"size(512);default(/static/img/default.jpg)"`

	PostFavorts  []*PostFavort `orm:"reverse(many)"`
	PostComments []*PostComment `orm:"reverse(many)"`
	Posts        []*Post `orm:"reverse(many)"`  // 设置一对多的反向关系
}

type PostImage struct {
	Id     int64
	Url    string
	Size   string
	Name   string
	Belong *Post `orm:"rel(fk)"`
}

type PostFavort struct {
	Id     int64
	User   *User  `orm:"rel(fk)"`
	Belong *Post `orm:"rel(fk)"`
}

type PostComment struct {
	Id           int64
	Type         int `orm:"default(0)"`
	User         *User  `orm:"rel(fk)"`
	ToReplayUser *User  `orm:"rel(fk)"`
	Content      string `orm:"size(1024)"`
	Belong       *Post `orm:"rel(fk)"`
}

type Post struct {
	Id          int64
	Content     string   `orm:"size(1024)"`         //设置不为空并不能大小24个字符
	CreateTime  string
	Type        int  `orm:"default(1)"`
	LinkImg     string
	LinkTitle   string
	LinkUrl     string
	LinkDesc    string
	Author      *User  `orm:"rel(fk)"`              //设置一对多关系(外键)
	Images      []*PostImage  `orm:"reverse(many)"` // 设置一对多的反向关系
	Comments    []*PostComment `orm:"reverse(many)"`
	Favorts     []*PostFavort `orm:"reverse(many)"`
	VideoUrl    string
	VideoImgUrl string
}

func InitTable() {
	orm.RegisterModel(new(User), new(Post), new(PostComment), new(PostFavort), new(PostImage));
	//判断是否创建了sqlite3数据库
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 50)

}

func Login(u *been.UserFrom) (*User, bool) {
	o := orm.NewOrm()
	table := o.QueryTable("user")
	seter := table.Filter("Mobile", u.Mobile).Filter("Pwd", u.Pwd)
	user := &User{};
	seter.One(user)
	return user, seter.Exist()
}

//http://localhost:8080/register?mobile=18820792655&pwd=123456
func Register(u *been.UserFrom) (*User, error) {
	o := orm.NewOrm()
	conut, _ := o.QueryTable("user").Filter("Mobile", u.Mobile).Count()
	beego.Info("Register", u.Mobile)
	if conut > 0 {
		beego.Info("Register", "该手机号码已注册")
		return new(User),
			errors.New("该手机号码已注册")
	}

	beego.Info("Register", "该手机号码可使用")
	time := time.Now().Format("2006-01-02 15:04:05")
	user := &User{
		Mobile:u.Mobile,
		Pwd:u.Pwd,
		Name:u.Mobile,
		Age:0,
		Sex:false,
		Address:"广州市天河区",
		Pic:"/static/img/default.jpg",
		Birth:time}
	_, err := o.Insert(user)
	if err == nil {
		return user, err
	}
	return new(User), err
}

func QueryPost(pageIndex int, pageSize int) ([]*Post, error) {
	o := orm.NewOrm()
	var p []*Post
	_, err := o.QueryTable("post").OrderBy("-CreateTime").Limit(pageSize, pageIndex * pageSize).RelatedSel().All(&p)
	if err == nil {
		for k, v := range p {
			o.QueryTable("PostImage").Filter("Belong", v).All(&p[k].Images)
			o.QueryTable("PostComment").Filter("Belong", v).RelatedSel("User", "ToReplayUser").All(&p[k].Comments)
			o.QueryTable("PostFavort").Filter("Belong", v).OrderBy("-id").RelatedSel("User").All(&p[k].Favorts)
		}
	}
	return p, err

}

func AddVideoPost(userId int, content string, cType int, videpath string, imgurl string) (*Post, error) {
	o := orm.NewOrm()
	u := &User{Id:int64(userId)}
	err := o.Read(u)
	if err == nil {
		time := time.Now().Format("2006-01-02 15:04:05")
		post := &Post{Author:u, Content:content, Type:cType, VideoImgUrl:imgurl, VideoUrl:videpath, CreateTime:time}
		_, err = o.Insert(post)
		if err == nil {
			return post, nil
		}
	}

	return &Post{}, err
}

func AddUrlPost(userId int, content string, cType int, icon string, url string, title string, desc string) (*Post, error) {
	o := orm.NewOrm()
	u := &User{Id:int64(userId)}
	err := o.Read(u)
	if err == nil {
		time := time.Now().Format("2006-01-02 15:04:05")
		post := &Post{Author:u, Content:content, Type:cType, LinkImg:icon, LinkTitle:title, LinkUrl:url, LinkDesc:desc, CreateTime:time}
		_, err = o.Insert(post)
		if err == nil {
			return post, nil
		}
	}

	return &Post{}, err
}

func AddPost(userId int, content string, cType int, haveimg bool, imgs []string, sizes []string) (*Post, error) {

	o := orm.NewOrm()
	u := &User{Id:int64(userId)}
	err := o.Read(u)
	o.Begin()
	if err == nil {
		time := time.Now().Format("2006-01-02 15:04:05")
		post := new(Post)
		post.Author = u
		post.Content = content
		post.Type = cType
		post.CreateTime = time
		id, err := o.Insert(post)
		if err == nil {
			if haveimg {
				var size string = ""
				for i, v := range sizes {
					if i == 0 {
						size = fmt.Sprint("#", v)
					} else {
						size = fmt.Sprint(size, "#", v)
					}

				}

				for _, v := range imgs {
					p := &PostImage{Belong:post, Url:v, Size:size, Name:v}
					_, err = o.Insert(p)
				}

			}

			o.Commit()

			var P Post
			o.QueryTable("Post").Filter("Id", id).RelatedSel().One(&P)
			o.QueryTable("PostImage").Filter("Belong", P).All(&P.Images)
			return &P, nil
		}

	}

	return &Post{}, err

}

func DeletePost(postId int, userId int) error {
	o := orm.NewOrm()
	var p Post;
	seter := o.QueryTable("post").Filter("id", postId).Filter("Author__id", userId)
	err := seter.One(&p)
	beego.Info("delete Type:", p.Type)
	if err == nil {
		if p.Type == 1 {
			//图片类型
			o.QueryTable("PostImage").Filter("Belong", p).All(&p.Images)
			//删除图片
			if p.Images != nil&&len(p.Images) > 0 {
				for _, v := range p.Images {
					beego.Info("delete Image:", v.Url)
					if utils.StringIsNotEmpty(v.Url) {

						os.Remove(utils.GetSavePathBySize(v.Url, "l"))
						os.Remove(utils.GetSavePathBySize(v.Url, "m"))
						os.Remove(utils.GetSavePathBySize(v.Url, "s"))
						os.Remove(v.Url)
					}
				}
			}
		}
		if p.Type == 3 {

			//视频类型
			if utils.StringIsNotEmpty(p.VideoImgUrl) {
				os.Remove(p.VideoImgUrl)
				beego.Info("delete Video img:", p.VideoImgUrl)
			}
			if utils.StringIsNotEmpty(p.VideoUrl) {
				beego.Info("delete Video PATH:", p.VideoUrl)
				os.Remove(p.VideoUrl)
			}
		}

		_, err = seter.Delete()

	}
	return err

}

func AddFavort(postId int, userId int) (*PostFavort, error) {
	o := orm.NewOrm()
	p := &Post{Id: int64(postId)}
	u := &User{Id:int64(userId)}

	b := o.QueryTable("PostFavort").Filter("Belong", postId).Filter("User__id", userId).Exist()

	beego.Debug("判断用户是否点赞啦...", b)
	if b == false {
		o.Read(u)
		f := &PostFavort{User:u, Belong:p}
		_, err := o.Insert(f)
		return f, err
	}

	return &PostFavort{}, errors.New("已点赞")
}

func AddComment(content string, cType int, userId int, touserId int, postId int) (*PostComment, error) {
	o := orm.NewOrm()
	u := &User{Id:int64(userId)}
	tu := &User{Id:int64(touserId)}
	o.Read(u)
	o.Read(tu)
	p := &Post{Id: int64(postId)}
	pc := &PostComment{Type:cType, Content:content, User:u, ToReplayUser:tu, Belong:p}

	_, err := o.Insert(pc)

	return pc, err

}
func DeleteFavort(postId int, userId int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("PostFavort").Filter("Belong", postId).Filter("User__id", userId).Delete()
	return err
}

func DeleteComment(commentId int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("PostComment").Filter("id", commentId).Delete()
	return err
}
