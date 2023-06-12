//go:build !search && !mgmt && !userprefs

package extensions_test

import (
	"encoding/json"
	"os"
	"testing"

	distext "github.com/opencontainers/distribution-spec/specs-go/v1/extensions"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/resty.v1"

	"zotregistry.io/zot/pkg/api"
	"zotregistry.io/zot/pkg/api/config"
	"zotregistry.io/zot/pkg/api/constants"
	extconf "zotregistry.io/zot/pkg/extensions/config"
	"zotregistry.io/zot/pkg/test"
)

func TestGetExensionsDisabled(t *testing.T) {
	Convey("start zot server with extensions but no extensions built", t, func(c C) {
		conf := config.New()
		port := test.GetFreePort()
		baseURL := test.GetBaseURL(port)

		conf.HTTP.Port = port

		defaultVal := true

		searchConfig := &extconf.SearchConfig{
			BaseConfig: extconf.BaseConfig{Enable: &defaultVal},
		}

		mgmtConfg := &extconf.MgmtConfig{
			BaseConfig: extconf.BaseConfig{Enable: &defaultVal},
		}

		conf.Extensions = &extconf.ExtensionConfig{
			Search: searchConfig,
			Mgmt:   mgmtConfg,
		}

		logFile, err := os.CreateTemp("", "zot-log*.txt")
		So(err, ShouldBeNil)
		conf.Log.Output = logFile.Name()
		defer os.Remove(logFile.Name()) // clean up

		ctlr := makeController(conf, t.TempDir(), "")

		cm := test.NewControllerManager(ctlr)
		cm.StartAndWait(port)
		defer cm.StopServer()

		var extensionList distext.ExtensionList

		resp, err := resty.R().Get(baseURL + constants.RoutePrefix + constants.ExtOciDiscoverPrefix)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		err = json.Unmarshal(resp.Body(), &extensionList)
		So(err, ShouldBeNil)
		So(len(extensionList.Extensions), ShouldEqual, 0)
	})
}

func makeController(conf *config.Config, dir string, copyTestDataDest string) *api.Controller {
	ctlr := api.NewController(conf)

	if copyTestDataDest != "" {
		test.CopyTestFiles(copyTestDataDest, dir)
	}
	ctlr.Config.Storage.RootDirectory = dir

	return ctlr
}