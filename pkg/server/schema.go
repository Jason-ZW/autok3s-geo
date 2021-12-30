package server

import (
	"net/http"

	geostore "github.com/Jason-ZW/autok3s-geo/pkg/store/geo"
	geotypes "github.com/Jason-ZW/autok3s-geo/pkg/types"

	"github.com/rancher/apiserver/pkg/types"
)

func initGeoIP(s *types.APISchemas) {
	s.MustImportAndCustomize(geotypes.GeoIP{}, func(schema *types.APISchema) {
		schema.Store = &geostore.Store{}
		schema.CollectionMethods = []string{http.MethodGet, http.MethodPost}
		schema.ResourceMethods = []string{http.MethodGet, http.MethodDelete}
	})
}
