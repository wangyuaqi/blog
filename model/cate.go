package model

// Cate 分类
type Cate struct {
	Id    int    `xorm:"pk autoincr INT(11)" json:"id" form:"id"`
	Name  string `xorm:"unique VARCHAR(64)" json:"name" form:"name"`
	Intro string `xorm:"VARCHAR(64)" json:"intro,omitempty" form:"intro"`
}

// CateIds 通过id返回新闻类别信息集合
func cateIds(ids []int) map[int]*Cate {
	mods := make([]Cate, 0, 6)
	DB.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapSet := make(map[int]*Cate, len(mods))
		for idx := range mods {
			mods[idx].Intro = ""
			mapSet[mods[idx].Id] = &mods[idx]
		}
		return mapSet
	}
	return nil
}

//CateGet 一个分类
func CateGet(id int) (*Cate, bool) {
	mod := &Cate{
		Id: id,
	}
	has, _ := DB.Get(mod)
	return mod, has
}

// CateName 通过name 查询分类
func CateName(nam string) (*Cate, bool) {
	mod := &Cate{
		Name: nam,
	}
	has, _ := DB.Get(mod)
	return mod, has
}

// CateAll 所有分类
func CateAll() ([]Cate, error) {
	mods := make([]Cate, 0, 4)
	err := DB.Asc("id").Find(&mods)
	return mods, err
}

// CateAdd 添加分类
func CateAdd(mod *Cate) bool {
	sess := DB.NewSession()
	defer sess.Close()
	sess.Begin()
	affect, _ := sess.InsertOne(mod)
	if affect != 1 {
		sess.Rollback()
		return false
	}
	sess.Commit()
	return true
}

// CateEdit 修改分类
func CateEdit(mod *Cate) bool {
	sess := DB.NewSession()
	defer sess.Close()
	sess.Begin()
	affect, err := sess.ID(mod.Id).Cols("Name", "Intro").Update(mod)
	if affect >= 0 && err == nil {
		sess.Commit()
		return true
	}
	sess.Rollback()
	return false
}

// CateDel 删除分类
func CateDel(id int) bool {
	sess := DB.NewSession()
	defer sess.Close()
	sess.Begin()
	if affect, err := sess.ID(id).Delete(&Cate{}); affect > 0 && err == nil {
		sess.Commit()
		DB.ClearCacheBean(&Cate{}, string(id))
		return true
	}
	sess.Rollback()
	return false
}

// CatePostCount 通过标签查询文章分页总数
// lmt 是否前台限制
func CatePostCount(cid int, lmt bool) int {
	sess := DB.NewSession()
	defer sess.Close()
	if cid > 0 {
		sess.Where("Cate_id = ?  ", cid)
	}
	if lmt {
		sess.Where("Is_Public = 1 and Status = 3 ")
	}
	sess.Where("Type = 0")
	mod := &Post{}
	count, _ := sess.Count(mod)
	return int(count)
}

// CatePostList 通过分类查询文章分页
// lmt 是否前台限制
func CatePostList(cid, pi, ps int, lmt bool) ([]Post, error) {
	mods := make([]Post, 0, ps)
	sess := DB.NewSession()
	defer sess.Close()
	if cid > 0 {
		sess.Where("Cate_id = ? ", cid)
	}
	if lmt {
		sess.Where("Is_Public = 1 and Status = 3 ")
	}
	sess.Where("Type = 0")
	err := sess.Cols("id", "title", "path", "create_time", "summary", "comment_num", "options").Desc("create_time").Limit(ps, (pi-1)*ps).Find(&mods)
	if len(mods) > 0 {
		for idx := range mods {
			if !lmt {
				mods[idx].Summary = ""
			}
			mods[idx].Options = ""
			mods[idx].Content = ""
			mods[idx].MarkdownContent = ""
		}
	}
	return mods, err
}
