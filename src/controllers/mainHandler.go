package controllers

import (
	"strconv"
	"os"
	"io"
	"../utils"
	"github.com/hunterhug/go_image"
	"../models"
	"github.com/astaxie/beego"
	"./been"
)
//获取帖子列表
//http://localhost:8080//postList?pageIndex=0&pageSize=5
func (this *MainController)GetPostList() {
	pageIndex, _ := strconv.Atoi(this.Input().Get("pageIndex"))
	pageSize, _ := strconv.Atoi(this.Input().Get("pageSize"))
	posts, err := models.QueryPost(pageIndex, pageSize)
	if err == nil {
		ReturnSuccess(&(this.Controller), posts)
	} else {
		ReturnError(&(this.Controller), err.Error(), FAIL)
	}
}

func (this *MainController)DeletePost() {
	postId, _ := strconv.Atoi(this.Input().Get("postId"))
	userId, _ := strconv.Atoi(this.Input().Get("userId"))
	err := models.DeletePost(postId, userId)
	if err == nil {
		ReturnSuccess(&(this.Controller), "")
	} else {
		ReturnError(&(this.Controller), err.Error(), FAIL)
	}
}

func (this *MainController)DeleteComment() {
	commentId, _ := strconv.Atoi(this.Input().Get("commentId"))
	err := models.DeleteComment(commentId)
	if err == nil {
		ReturnSuccess(&(this.Controller), commentId)
	} else {
		ReturnError(&(this.Controller), err.Error(), FAIL)
	}
}

func (this *MainController)AddComment() {
	optype, _ := strconv.Atoi(this.Input().Get("cType"))
	content := this.Input().Get("content")
	userId, _ := strconv.Atoi(this.Input().Get("userId"))
	touserId, _ := strconv.Atoi(this.Input().Get("touserId"))
	postId, _ := strconv.Atoi(this.Input().Get("postId"))
	pc, err := models.AddComment(content, optype, userId, touserId, postId)
	if err == nil {
		ReturnSuccess(&(this.Controller), pc)
	} else {
		ReturnError(&(this.Controller), err.Error(), FAIL)
	}
}

func (this *MainController)AddFavort() {
	postId, _ := strconv.Atoi(this.Input().Get("postId"))
	userId, _ := strconv.Atoi(this.Input().Get("userId"))
	f, err := models.AddFavort(postId, userId)
	if err == nil {
		ReturnSuccess(&(this.Controller), f)
	} else {
		ReturnError(&(this.Controller), err.Error(), FAIL)
	}
}

func (this *MainController)DeleteFavort() {
	postId, _ := strconv.Atoi(this.Input().Get("postId"))
	userId, _ := strconv.Atoi(this.Input().Get("userId"))
	err := models.DeleteFavort(postId, userId)
	if err == nil {
		ReturnSuccess(&(this.Controller), userId)
	} else {
		ReturnError(&(this.Controller), err.Error(), FAIL)
	}
}

/**
 点赞操作
 */
func (this *MainController)FavortOp() {

	op := this.Input().Get("op")
	postId, _ := strconv.Atoi(this.Input().Get("postId"))
	userId, _ := strconv.Atoi(this.Input().Get("userId"))
	switch op {
	case "add":
		f, err := models.AddFavort(postId, userId)
		if err == nil {
			ReturnSuccess(&(this.Controller), f)
		} else {
			ReturnError(&(this.Controller), err.Error(), FAIL)
		}
		break

	case "delete":
		err := models.DeleteFavort(postId, userId)
		if err == nil {
			ReturnSuccess(&(this.Controller), userId)
		} else {
			ReturnError(&(this.Controller), err.Error(), FAIL)
		}
		break

	default:
		ReturnError(&(this.Controller), "操作失败", ERROR)
		break


	}

}

func (this *MainController)AddUrlPost() {

	userId, _ := strconv.Atoi(this.Input().Get("userId"))
	content := this.Input().Get("content")
	contenttype, _ := strconv.Atoi(this.Input().Get("type"))

	shareIcon := this.Input().Get("shareIcon")
	shareTitle := this.Input().Get("shareTitle")
	shareDesc := this.Input().Get("shareDesc")
	shareUrl := this.Input().Get("shareUrl")

	post, err := models.AddUrlPost(userId, content, contenttype, shareIcon, shareUrl, shareTitle, shareDesc)
	if err == nil {
		ReturnSuccess(&(this.Controller), post)
	} else {
		ReturnError(&(this.Controller), err.Error(), FAIL)
	}

}

func (this *MainController)AddVideoPost() {
	userId, _ := strconv.Atoi(this.Input().Get("userId"))
	content := this.Input().Get("content")
	contenttype, _ := strconv.Atoi(this.Input().Get("type"))
	//VIDEO_SCREENSHOT
	video, videoHead, err := this.GetFile("video")
	defer video.Close()
	if err == nil {
		//2MB
		buf := make([]byte, 2 * 1024 * 1024)
		var l int
		l, err = video.Read(buf)
		if err == nil {
			beego.Info("上传视频文件：", videoHead.Filename)
			//限制上传文件大小
			if l < (2 * 1024 * 1024) {
				videoPath := utils.GetSavePath(videoHead.Filename, VideoPath)
				err = this.SaveToFile("video", videoPath)
				if err == nil {
					img, imgHead, _ := this.GetFile("videoImg")
					defer img.Close()
					imgPath := utils.GetSavePath(imgHead.Filename, VideoImgPath)
					this.SaveToFile("videoImg", imgPath)

					post, err := models.AddVideoPost(userId, content, contenttype, videoPath, imgPath)

					if err == nil {
						ReturnSuccess(&(this.Controller), post)
					} else {
						ReturnError(&(this.Controller), err.Error(), FAIL)
					}

				} else {
					ReturnError(&(this.Controller), "获取文件失败", FAIL)
				}
			} else {
				ReturnError(&(this.Controller), "上传文件过大", FAIL)
			}

		}
	}

}

var ImageSize = []string{"l", "m", "s"}

func (this *MainController)AddPost() {

	userId, _ := strconv.Atoi(this.Input().Get("userId"))
	content := this.Input().Get("content")
	contenttype, _ := strconv.Atoi(this.Input().Get("type"))
	haveimg := this.Input().Get("haveimg")

	beego.Info("userId", userId, content, contenttype, haveimg)
	if haveimg == "have" {
		imgs := make([]string, 0)
		files, err := this.GetFiles("photos")
		if err == nil {
			for i, f := range files {
				beego.Error("获取文件：", f.Filename)
				//for each fileheader, get a handle to the actual file
				file, err := files[i].Open()
				defer file.Close()
				if err == nil {
					//path := utils.GetSavePath(f.Filename, IamgePath)
					original := utils.GetSavePath(f.Filename, IamgePath)

					//create destination file making sure the path is writeable.
					//return OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
					//    os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
					dst, err := os.Create(original)
					defer dst.Close()
					if err == nil {
						if _, err = io.Copy(dst, file); err != nil {
							beego.Error("---- save upload  filename:", f.Filename, "  user Name:", userId)
							ReturnError(&(this.Controller), err.Error(), FAIL)
							return
						} else {

							imgs = append(imgs, original)
							ch1 := make(chan bool, len(ImageSize))
							defer close(ch1)
							go GenerateMultImage(ch1, original)
							<-ch1
						}
					}
				} else {
					beego.Error("---- open upload filename :", err.Error())
					ReturnError(&(this.Controller), err.Error(), FAIL)
					return
				}
			}

			post, err := models.AddPost(userId, content, contenttype, true, imgs, ImageSize)
			if err == nil {
				ReturnSuccess(&(this.Controller), post)
			} else {
				ReturnError(&(this.Controller), err.Error(), FAIL)
			}
			return

		} else {
			beego.Error("---- open upload filename :", err.Error())
			ReturnError(&(this.Controller), err.Error(), FAIL)
			return
		}

	}

	post, err := models.AddPost(userId, content, contenttype, false, nil, nil)
	if err == nil {
		ReturnSuccess(&(this.Controller), post)
	} else {
		ReturnError(&(this.Controller), err.Error(), FAIL)
	}

}

/**
生成多种尺寸图片
 */
func GenerateMultImage(c chan bool, file string) {

	for _, v := range ImageSize {

		newFileName := utils.GetSavePathBySize(file, v)
		if v == "l" {
			//生成大图
			go_image.ScaleF2F(file, newFileName, 720)
		} else if v == "m" {
			//中图
			go_image.ScaleF2F(file, newFileName, 480)
		} else if v == "s" {
			//小图
			go_image.ScaleF2F(file, newFileName, 240)
		}

	}

	c <- true
}

func ReturnSuccess(this *beego.Controller, data interface{}) {

	beego.Info("返回Json:", data)
	msg := &been.ReturnMsg{"success", SUCCESS, data}
	this.Data["json"] = msg
	this.ServeJSON()
}

func ReturnError(this *beego.Controller, errmsg string, errcode int) {
	if len(errmsg) == 0 {
		errmsg = "操作失败";
	}
	msg := &been.ReturnMsg{errmsg, errcode, nil}
	this.Data["json"] = msg
	this.ServeJSON()
}
