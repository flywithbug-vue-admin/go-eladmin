package model_app

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-eladmin/common"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var (
	appPlatformMap = map[string]string{"IOS": "iOS", "ANDROID": "Android", "H5": "H5", "SERVER": "Server"}
)

const (
	appVersionCollection = mongo_index.CollectionAppVersion
)

const (
	appStatusTypeUnDetermined typeStatus = iota //待定
	appStatusTypePrepare                        //准备中 待开发
	appStatusTypeDeveloping                     //开发中 待灰度
	appStatusTypeGray                           //灰度  待发布
	appStatusTypeRelease                        //已发布  已发布不能再更改
	appStatusTypeWorkDone                       //工作全部结束

)

type AppVersion struct {
	Id            int64      `json:"id,omitempty" bson:"_id,omitempty"`
	AppId         int64      `json:"app_id,omitempty" bson:"app_id,omitempty"` //所属App DB Id
	Version       string     `json:"version,omitempty" bson:"version,omitempty"`
	ParentVersion string     `json:"parent_version,omitempty" bson:"parent_version,omitempty"`
	VersionNum    int        `json:"version_num,omitempty" bson:"version_num,omitempty"`
	Platform      []string   `json:"platform,omitempty" bson:"platform,omitempty"`           //(iOS,Android,H5,Server)["iOS","Android","H5","Server"]
	Status        typeStatus `json:"status,omitempty" bson:"status,omitempty"`               //状态    1(准备中) 2(开发中) 3(灰度) 4(已发布)
	ApprovalTime  int64      `json:"approval_time,omitempty" bson:"approval_time,omitempty"` //立项时间
	LockTime      int64      `json:"lock_time,omitempty" bson:"lock_time,omitempty"`         //锁版时间
	GrayTime      int64      `json:"gray_time,omitempty" bson:"gray_time,omitempty"`         //灰度时间
	CreateTime    int64      `json:"create_time,omitempty" bson:"create_time,omitempty"`     //添加时间
	AppStatus     string     `json:"app_status,omitempty" bson:"app_status,omitempty"`       //app状态
	ReleaseTime   int64      `json:"release_time,omitempty" bson:"release_time,omitempty"`
}

func versionTime(a, b int64) int64 {
	if a > 0 {
		return a
	}
	return b
}

func (app AppVersion) checkAppVersionTime() error {
	if app.Id == 0 {
		if app.ApprovalTime > app.LockTime {
			return fmt.Errorf("立项时间不能晚于锁版时间")
		}
		if app.LockTime > app.GrayTime {
			return fmt.Errorf("锁版时间不能晚于灰度时间")
		}
		if app.ReleaseTime > 0 && app.GrayTime > app.ReleaseTime {
			return fmt.Errorf("灰度时间不能晚于发布时间")
		}
		return nil
	}
	a, err := app.FindOne()
	if err != nil {
		return err
	}
	if app.ApprovalTime > 0 {
		if app.ApprovalTime > versionTime(app.LockTime, a.LockTime) {
			return fmt.Errorf("立项时间不能晚于锁版时间")
		}
	}
	if app.LockTime > 0 {
		if app.LockTime > versionTime(app.GrayTime, a.GrayTime) {
			return fmt.Errorf("锁版时间时间不能早于灰度时间")
		}
	}
	if app.GrayTime > 0 {
		if app.GrayTime < versionTime(app.LockTime, a.LockTime) {
			return fmt.Errorf("灰度时间不能早于锁版时间")
		}
	}
	if app.ReleaseTime > 0 {
		if app.ReleaseTime < versionTime(app.GrayTime, a.GrayTime) {
			return fmt.Errorf("发布时间不能早于灰度时间")
		}
	}
	return nil
}

func (app AppVersion) ToJson() string {
	js, _ := json.Marshal(app)
	return string(js)
}

func (app AppVersion) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), appVersionCollection, docs...)
}

func (app AppVersion) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), appVersionCollection, query)
}

func (app AppVersion) findOne(query, selector interface{}) (AppVersion, error) {
	ap := AppVersion{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), appVersionCollection, query, selector, &ap)
	return ap, err
}

func (app AppVersion) findAll(query, selector interface{}) (results []AppVersion, err error) {
	results = []AppVersion{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), appVersionCollection, query, selector, &results)
	return results, err
}

func (app AppVersion) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), appVersionCollection, selector, update, true)
}

func (app AppVersion) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), appVersionCollection, selector)
}

func (app AppVersion) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), appVersionCollection, selector)
}

func (app AppVersion) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), appVersionCollection, query, selector)
}

func (app AppVersion) findPage(page, limit int, query, selector interface{}, fields ...string) (results []AppVersion, err error) {
	results = []AppVersion{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), appVersionCollection, page, limit, query, selector, &results, fields...)
	return
}

func (app AppVersion) FindOne() (AppVersion, error) {
	app, err := app.findOne(bson.M{"_id": app.Id}, nil)
	return app, err
}

func (app AppVersion) FindAll(query, selector interface{}) ([]AppVersion, error) {
	results, err := app.findAll(query, selector)
	return results, err
}

func (app *AppVersion) Insert() error {
	err := app.checkTimeValid()
	if err != nil {
		return err
	}
	var application = Application{}
	if !application.isExist(bson.M{"_id": app.AppId}) {
		return fmt.Errorf("appID:%d not found", app.AppId)
	}
	if app.isExist(bson.M{"version": app.Version, "app_id": app.AppId}) {
		return fmt.Errorf("version exist")
	}

	if len(app.ParentVersion) > 0 && !strings.EqualFold(app.ParentVersion, "-") {
		if !app.isExist(bson.M{"version": app.ParentVersion, "app_id": app.AppId}) {
			return errors.New("parent_version not exist")
		}
	}
	if len(app.Platform) == 0 {
		return errors.New("platform must choose")
	}
	for _, platform := range app.Platform {
		_, ok := appPlatformMap[strings.ToUpper(platform)]
		if !ok {
			return fmt.Errorf("platform must like (iOS,Android,H5,Server) ")
		}
	}
	err = app.checkAppVersionTime()
	if err != nil {
		return err
	}
	app.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), appVersionCollection)
	app.CreateTime = time.Now().Unix()
	app.Status = appStatusTypePrepare
	app.AppStatus = makeStatusString(appStatusTypePrepare)
	compareState, err := common.VersionCompare(app.Version, app.ParentVersion)
	if err != nil {
		return err
	}
	if compareState != common.CompareVersionStateGreater {
		return errors.New("new Version must bigger than ParentVersion")
	}
	if len(app.ParentVersion) == 0 {
		app.ParentVersion = "-"
	}

	return app.insert(app)
}

func (app AppVersion) checkTimeValid() error {
	if app.ApprovalTime > app.LockTime {
		return errors.New("approval time must early than lock time")
	}

	if app.LockTime > app.GrayTime {
		return errors.New("lock time must early than gray time")
	}
	if app.Status == appStatusTypeRelease {
		if app.GrayTime > app.ReleaseTime {
			return errors.New("gray time must early than release time")
		}
	}
	return nil
}

func (app AppVersion) Remove() error {
	selector := bson.M{"_id": app.Id}
	app, err := app.findOne(selector, nil)
	if err != nil {
		return err
	}
	if app.Status != appStatusTypePrepare {
		return errors.New("版本已锁定，不能删除")
	}
	return app.remove(selector)
}

func (app AppVersion) Update() error {
	if app.ParentVersion == "-" {
		app.ParentVersion = ""
	}
	if err := app.checkTimeValid(); err != nil {
		return err
	}
	err := app.checkAppVersionTime()
	if err != nil {
		return err
	}
	appOld, _ := app.findOne(bson.M{"_id": app.Id}, nil)
	if app.Status > appStatusTypeWorkDone || app.Status < appStatusTypeUnDetermined {
		return errors.New("status not right")
	}
	if app.Status < appOld.Status {
		app.Status = appOld.Status
	}

	selector := bson.M{"_id": app.Id}
	if app.Status > 1 {
		//状态大于1时，可以更新锁版时间，灰度时间，状态，和发布时间
		app.AppStatus = makeStatusString(app.Status)
		if app.Status == appStatusTypeWorkDone {
			if app.ReleaseTime == 0 {
				app.ReleaseTime = time.Now().Unix()
			}
		}
		if app.Status == appStatusTypeRelease {
			if app.ReleaseTime == 0 {
				app.ReleaseTime = time.Now().Unix()
			}
		}
		if app.Status > appStatusTypePrepare {
			app.ApprovalTime = 0
		}
		if app.Status > appStatusTypeDeveloping {
			app.LockTime = 0
		}
		if app.Status > appStatusTypeGray {
			app.GrayTime = 0
		}
		app.AppId = 0
		return app.update(selector, app)
	} else {
		app.ReleaseTime = 0
		//判断非当前version id的版本号是否存在
		if app.isExist(bson.M{"version": app.Version, "app_id": app.AppId, "_id": bson.M{"$ne": app.Id}}) {
			return fmt.Errorf("version exist")
		}
		if len(app.ParentVersion) > 0 {
			if !app.isExist(bson.M{"version": app.ParentVersion, "app_id": app.AppId}) {
				return errors.New("parent_version not exist")
			}
		}
		if len(app.Platform) == 0 {
			return errors.New("platform must choose")
		}
		for _, platform := range app.Platform {
			_, ok := appPlatformMap[strings.ToUpper(platform)]
			if !ok {
				return fmt.Errorf("platform must like (iOS,Android,H5,Server) ")
			}
		}
		compareState, err := common.VersionCompare(app.Version, app.ParentVersion)
		if err != nil {
			return err
		}
		if compareState != common.CompareVersionStateGreater {
			return errors.New("new Version must bigger than ParentVersion")
		}
		if len(app.ParentVersion) == 0 {
			app.ParentVersion = "-"
		}
	}
	app.AppId = 0
	return app.update(selector, app)
}

func makeStatusString(status typeStatus) string {
	statusString := "待定"
	switch status {
	case appStatusTypePrepare:
		statusString = "准备中"
		break
	case appStatusTypeDeveloping:
		statusString = "开发中"
		break
	case appStatusTypeGray:
		statusString = "灰度"
		break
	case appStatusTypeRelease:
		statusString = "已发布"
		break
	case appStatusTypeWorkDone:
		statusString = "Done"
		break
	case appStatusTypeUnDetermined:
		statusString = "待定"
		break
	default:
		statusString = "未定义状态"
	}
	return statusString
}

func (app AppVersion) TotalCount(query, selector interface{}) (int, error) {
	return app.totalCount(query, selector)
}
func (app AppVersion) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) (apps []AppVersion, err error) {
	return app.findPage(page, limit, query, selector, fields...)
}
