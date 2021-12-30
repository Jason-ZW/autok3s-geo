package geo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Jason-ZW/autok3s-geo/pkg/types"
)

func FreeGeoXhr(ip string) (*types.GeoIP, error) {
	geo := &types.GeoIP{}
	r, err := http.Get("https://freegeoip.live/json/" + ip)
	if err != nil {
		return geo, err
	}
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return geo, err
	}

	err = json.Unmarshal(body, geo)
	return geo, err
}
