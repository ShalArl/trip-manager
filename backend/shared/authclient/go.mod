module github.com/ShalArl/trip-manager/backend/shared/authclient

go 1.25.8

require (
	"github.com/ShalArl/trip-manager/backend/shared/tenantdb" v0.0.0
)

replace (
	"github.com/ShalArl/trip-manager/backend/shared/tenantdb" => ../tenantdb
)