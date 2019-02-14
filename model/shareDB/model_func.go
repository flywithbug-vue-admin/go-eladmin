package shareDB

var (
	docManager = "doc_manager"
	monitorDB  = "monitor"
)

func SetMonitorDBName(dbName string) {
	monitorDB = dbName
}

func SetDocMangerDBName(dbName string) {
	docManager = dbName
}

func DocManagerDBName() string {
	return docManager
}
func MonitorDBName() string {
	return monitorDB
}

type OperationModel interface {
	ToJson() string

	isExist(query interface{}) bool
	insert(docs ...interface{}) error
	update(selector, update interface{}) error
	findOne(query, selector interface{}) (interface{}, error)
	findAll(query, selector interface{}) (results []interface{}, err error)
	remove(selector interface{}) error
	removeAll(selector interface{}) error
	totalCount(query, selector interface{}) (int, error)
	findPage(page, limit int, query, selector interface{}, fields ...string) (results []interface{}, err error)

	TotalCount(query, selector interface{}) (int, error)
	FindPage(page, limit int, query, selector interface{}, fields ...string) (results []interface{}, err error)
	Insert() (int64, error)
	Update() error
	FindOne(query, selector interface{}) (interface{}, error)
	Remove() error
	RemoveAll() error
}
