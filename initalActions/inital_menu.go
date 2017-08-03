package initalActions

import (
	"encoding/xml"
	md "golangERP/models"
	service "golangERP/services"
	"golangERP/utils"
	"io/ioutil"
	"os"

	"github.com/astaxie/beego/orm"
)

// InitMenu  菜单数据解析
type InitMenu struct {
	Name      string `xml:"name"`
	XMLID     string `xml:"id,attr"`
	Path      string `xml:"path"`
	Icon      string `xml:"icon"`
	Component string `xml:"component"`
	Sequence  int64  `xml:"sequence"`
	ParentID  string `xml:"parent_id,attr"`
}

// InitMenus 菜单数据列表
type InitMenus struct {
	XMLName xml.Name   `xml:"Menus"`
	Menus   []InitMenu `xml:"menu"`
}

// InitMenus2DB 菜单初始化，数据库创建记录
func InitMenus2DB(split string) {
	if xmlDir, err := os.Getwd(); err == nil {
		xmlBase := utils.StringsJoin(xmlDir, split, "inital_data", split, "xml", split, "menu")
		if dirList, err := ioutil.ReadDir(xmlBase); err == nil {
			for _, dir := range dirList {
				if dir.IsDir() {
					continue
				}
				if file, err := os.Open(utils.StringsJoin(xmlBase, split, dir.Name())); err == nil {
					defer file.Close()
					if data, err := ioutil.ReadAll(file); err == nil {
						var initMenus InitMenus
						var moduleName = "BaseMenu"
						if xml.Unmarshal(data, &initMenus) == nil {
							for _, menuXML := range initMenus.Menus {
								ormObj := orm.NewOrm()
								var xmlid = utils.StringsJoin(moduleName, ".", menuXML.XMLID)
								// 检查在系统中是否已经存在
								if _, err = md.GetModuleDataByXMLID(xmlid, ormObj); err != nil {
									var menu md.BaseMenu
									var parent md.BaseMenu
									menu.Name = menuXML.Name
									menu.Path = menuXML.Path
									menu.Component = menuXML.Component
									menu.Icon = menuXML.Icon
									menu.Sequence = menuXML.Sequence
									parentIDStr := menuXML.ParentID
									if parentIDStr != "" {
										if mobuleData, err := md.GetModuleDataByXMLID(utils.StringsJoin(moduleName, ".", parentIDStr), ormObj); err == nil {
											parent.ID = mobuleData.InsertID
											menu.Parent = &parent
										}
									}
									if insertID, err := service.ServiceCreateBaseMenu(&menu); err == nil {
										var moduleData md.ModuleData
										moduleData.InsertID = insertID
										moduleData.XMLID = xmlid
										moduleData.Descrition = menu.Name
										moduleData.ModuleName = moduleName
										md.AddModuleData(&moduleData, ormObj)
									}
								}
							}
						}
					}

				}
			}
		}
	}
}