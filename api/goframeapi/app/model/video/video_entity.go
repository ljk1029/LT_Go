package video

import (
	"database/sql"
	"github.com/gogf/gf/database/gdb"
)

// Entity is the golang structure for table user.
type Entity struct {
	Id               uint   `orm:"id,primary" json:"id"`                           // 模板id
	Title            string `orm:"title"   json:"title"`                           // 模板标题
	Content          string `orm:"content"   json:"content"`                       // 模板内容
	VideoTemplateUrl string `orm:"video_template_url"   json:"video_template_url"` // 模板视频地址
	Thumbnail        string `orm:"thumbnail"   json:"thumbnail"`                   // 模板缩略图
	FixedPicNum      string `orm:"fixed_pic_num"   json:"fixed_pic_num"`           // 固定图片数
	TemplateFunc     string `orm:"template_func"   json:"template_func"`           // 模板
	IsHead           string `orm:"is_head"   json:"is_head"`                       // 是否需要头像
	IsShow           string `orm:"is_show"   json:"is_show"`                       // 是否需要头像
	DefaultTitle     string `orm:"default_title"   json:"default_title"`           // 是否需要头像
}

// OmitEmpty sets OPTION_OMITEMPTY option for the model, which automatically filers
// the data and where attributes for empty values.
func (r *Entity) OmitEmpty() *arModel {
	return Model.Data(r).OmitEmpty()
}

// Inserts does "INSERT...INTO..." statement for inserting current object into table.
func (r *Entity) Insert() (result sql.Result, err error) {
	return Model.Data(r).Insert()
}

// InsertIgnore does "INSERT IGNORE INTO ..." statement for inserting current object into table.
func (r *Entity) InsertIgnore() (result sql.Result, err error) {
	return Model.Data(r).InsertIgnore()
}

// Replace does "REPLACE...INTO..." statement for inserting current object into table.
// If there's already another same record in the table (it checks using primary key or unique index),
// it deletes it and insert this one.
func (r *Entity) Replace() (result sql.Result, err error) {
	return Model.Data(r).Replace()
}

// Save does "INSERT...INTO..." statement for inserting/updating current object into table.
// It updates the record if there's already another same record in the table
// (it checks using primary key or unique index).
func (r *Entity) Save() (result sql.Result, err error) {
	return Model.Data(r).Save()
}

// Update does "UPDATE...WHERE..." statement for updating current object from table.
// It updates the record if there's already another same record in the table
// (it checks using primary key or unique index).
func (r *Entity) Update() (result sql.Result, err error) {
	return Model.Data(r).Where(gdb.GetWhereConditionOfStruct(r)).Update()
}

// Delete does "DELETE FROM...WHERE..." statement for deleting current object from table.
func (r *Entity) Delete() (result sql.Result, err error) {
	return Model.Where(gdb.GetWhereConditionOfStruct(r)).Delete()
}
