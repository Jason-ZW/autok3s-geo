package geo

import (
	"context"
	"fmt"

	"github.com/Jason-ZW/autok3s-geo/pkg/db"
	"github.com/Jason-ZW/autok3s-geo/pkg/util"

	"github.com/rancher/apiserver/pkg/apierror"
	"github.com/rancher/apiserver/pkg/store/empty"
	"github.com/rancher/apiserver/pkg/types"
	"github.com/rancher/wrangler/pkg/schemas/validation"
)

// Store holds provider's API state.
type Store struct {
	empty.Store
}

// Create creates geo by IP.
func (s *Store) Create(apiOp *types.APIRequest, schema *types.APISchema, data types.APIObject) (types.APIObject, error) {
	if !util.ValidRequest(apiOp.Request) {
		return types.APIObject{}, apierror.NewAPIError(validation.InvalidReference, "request is invalid")
	}

	ip, err := util.GetRealIPAddress(apiOp.Request)
	if err != nil {
		return types.APIObject{}, apierror.NewAPIError(validation.InvalidFormat, err.Error())
	}

	db, err := db.NewDB(context.Background())
	if err != nil {
		return types.APIObject{}, apierror.NewAPIError(validation.ServerError, err.Error())
	}
	defer db.Close()

	// filter duplicated IP address.
	existGeos, err := db.QueryByIPIn24Hours(ip)
	if len(existGeos) != 0 {
		return types.APIObject{}, apierror.NewAPIError(validation.NotUnique, fmt.Sprintf("request ip %s not unique", ip))
	}

	geo, err := FreeGeoXhr(ip)
	if err != nil {
		return types.APIObject{}, apierror.NewAPIError(validation.ServerError, err.Error())
	}

	// write to the db.
	geo, err = db.Write(geo)
	if err != nil {
		return types.APIObject{}, apierror.NewAPIError(validation.ServerError, err.Error())
	}

	return types.APIObject{
		Type:   schema.ID,
		ID:     ip,
		Object: geo,
	}, nil
}

// List returns geos as list.
func (s *Store) List(apiOp *types.APIRequest, schema *types.APISchema) (types.APIObjectList, error) {
	return types.APIObjectList{}, nil
}

// Delete deletes geo by IP.
func (s *Store) Delete(apiOp *types.APIRequest, schema *types.APISchema, id string) (types.APIObject, error) {
	return types.APIObject{}, nil
}
