package mongo_index

import (
	"go-eladmin/model/shareDB"

	"gopkg.in/mgo.v2"
)

// monitor
const (
	CollectionLog          = "log"
	CollectionVisitUId     = "visit_uid"
	CollectionVisitApi     = "visit_api"
	CollectionMonitorCount = "monitor_count"
)

func monitorIndex() []Index {
	var MonitorIndexes = []Index{
		{
			DBName:     shareDB.MonitorDBName(),
			Collection: CollectionLog,
			Index: mgo.Index{
				Key:        []string{"request_id"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_request_id_f_index",
			},
		},
		{
			DBName:     shareDB.MonitorDBName(),
			Collection: CollectionVisitUId,
			Index: mgo.Index{
				Key:        []string{"client_ip", "uuid", "time_date"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_client_ip_f_uuid_time_date_index",
			},
		},
		{
			DBName:     shareDB.MonitorDBName(),
			Collection: CollectionVisitApi,
			Index: mgo.Index{
				Key:        []string{"path", "method", "time_date"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_path_f_method_time_date_index",
			},
		},
		{
			DBName:     shareDB.MonitorDBName(),
			Collection: CollectionMonitorCount,
			Index: mgo.Index{
				Key:        []string{"monitor", "time_date"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_monitor_f_time_date_index",
			},
		},
	}
	return MonitorIndexes
}
