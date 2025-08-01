package omaha

import (
	"bytes"
	"encoding/xml"
	"os"
	"reflect"
	"strconv"
	"testing"

	omahaSpec "github.com/flatcar/go-omaha/omaha"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/flatcar/nebraska/backend/pkg/api"
)

const (
	defaultTestDbURL string = "postgres://postgres:nebraska@127.0.0.1:5432/nebraska_tests?sslmode=disable&connect_timeout=10"

	flatcarAppID = "e96281a6-d1af-4bde-9a0a-97b76e56dc57"

	reqVersion  string = "3"
	reqPlatform string = "coreos"
	reqSp       string = "linux"
	reqArch     string = "x64"
)

func checkDB(t *testing.T) {
	t.Helper()

	if _, ok := os.LookupEnv("NEBRASKA_DB_URL"); !ok {
		t.Logf("NEBRASKA_DB_URL not set, setting to default %q\n", defaultTestDbURL)
		_ = os.Setenv("NEBRASKA_DB_URL", defaultTestDbURL)
	}

	a, err := api.New(api.OptionInitDB)
	if err != nil {
		t.Error("These tests require PostgreSQL running and a tests database created, please adjust NEBRASKA_DB_URL as needed.")
		t.Fatalf("Failed to init DB: %v\n", err)
	}
	a.Close()
}

func newForTest(t *testing.T) *api.API {
	checkDB(t)

	a, err := api.NewForTest(api.OptionInitDB, api.OptionDisableUpdatesOnFailedRollout)

	require.NoError(t, err)
	require.NotNil(t, a)

	return a
}

func TestMain(m *testing.M) {
	if os.Getenv("NEBRASKA_SKIP_TESTS") != "" {
		return
	}

	os.Exit(m.Run())
}

func TestInvalidRequests(t *testing.T) {
	a := newForTest(t)
	defer a.Close()
	h := NewHandler(a)

	tTeam, _ := a.AddTeam(&api.Team{Name: "test_team"})
	tApp, _ := a.AddApp(&api.Application{Name: "test_app", Description: "Test app", TeamID: tTeam.ID})
	tPkg, _ := a.AddPackage(&api.Package{Type: api.PkgTypeFlatcar, URL: "http://sample.url/pkg", Version: "640.0.0", ApplicationID: tApp.ID, Arch: api.ArchAMD64})
	tChannel, _ := a.AddChannel(&api.Channel{Name: "test_channel", Color: "blue", ApplicationID: tApp.ID, PackageID: null.StringFrom(tPkg.ID), Arch: api.ArchAMD64})
	tGroup, _ := a.AddGroup(&api.Group{Name: "test_group", ApplicationID: tApp.ID, ChannelID: null.StringFrom(tChannel.ID), PolicyUpdatesEnabled: true, PolicySafeMode: true, PolicyPeriodInterval: "15 minutes", PolicyMaxUpdatesPerPeriod: 2, PolicyUpdateTimeout: "60 minutes"})

	validUnregisteredIP := "127.0.0.1"
	validUnregisteredMachineID := "some-id"
	validUnverifiedAppVersion := "100.0.1"
	addPing := false
	updateCheck := true
	noEventInfo := (*eventInfo)(nil)

	omahaResp := doOmahaRequest(t, h, tApp.ID, validUnverifiedAppVersion, validUnregisteredMachineID, "invalid-track", validUnregisteredIP, addPing, updateCheck, noEventInfo)
	checkOmahaResponse(t, omahaResp, tApp.ID, omahaSpec.AppStatus("error-failedToRetrieveUpdatePackageInfo"))

	omahaResp = doOmahaRequest(t, h, tApp.ID, validUnverifiedAppVersion, validUnregisteredMachineID, tGroup.ID, "invalid-ip", addPing, updateCheck, noEventInfo)
	checkOmahaResponse(t, omahaResp, tApp.ID, omahaSpec.AppStatus("error-instanceRegistrationFailed"))

	omahaResp = doOmahaRequest(t, h, "invalid-app-uuid", validUnverifiedAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, noEventInfo)
	checkOmahaResponse(t, omahaResp, "invalid-app-uuid", omahaSpec.AppStatus("error-failedToRetrieveUpdatePackageInfo"))

	omahaResp = doOmahaRequest(t, h, tApp.ID, "", validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, noEventInfo)
	checkOmahaResponse(t, omahaResp, tApp.ID, omahaSpec.AppStatus("error-instanceRegistrationFailed"))
}

func TestAppNoUpdateForAppWithChannelAndPackageName(t *testing.T) {
	a := newForTest(t)
	defer a.Close()
	h := NewHandler(a)

	tAppFlatcar, _ := a.GetApp(flatcarAppID)
	tPkgFlatcar640, _ := a.AddPackage(&api.Package{Type: api.PkgTypeFlatcar, URL: "http://sample.url/pkg", Version: "640.0.0", ApplicationID: tAppFlatcar.ID, Arch: api.ArchAMD64})
	tChannel, _ := a.AddChannel(&api.Channel{Name: "mychannel", Color: "white", ApplicationID: tAppFlatcar.ID, PackageID: null.StringFrom(tPkgFlatcar640.ID), Arch: api.ArchAMD64})
	tGroup, _ := a.AddGroup(&api.Group{Name: "Production", ApplicationID: tAppFlatcar.ID, ChannelID: null.StringFrom(tChannel.ID), PolicyUpdatesEnabled: true, PolicySafeMode: true, PolicyPeriodInterval: "15 minutes", PolicyMaxUpdatesPerPeriod: 2, PolicyUpdateTimeout: "60 minutes"})

	validUnregisteredIP := "127.0.0.1"
	validUnregisteredMachineID := "65e1266d-6f54-4b87-9080-23b99ca9c12f"
	expectedAppVersion := "640.0.0"
	updateCheck := true
	addPing := true

	// Now with an error event tag, no updatecheck tag
	omahaResp := doOmahaRequest(t, h, tAppFlatcar.ID, expectedAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, !addPing, !updateCheck, ei(omahaSpec.EventTypeUpdateComplete, omahaSpec.EventResultError, "268437959"))
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaEventResponse(t, omahaResp, tAppFlatcar.ID, 1)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, !addPing)
	checkOmahaNoUpdateResponse(t, omahaResp)

	// Now updatetag, successful event, no previous version
	omahaResp = doOmahaRequest(t, h, tAppFlatcar.ID, expectedAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, !addPing, updateCheck, ei(omahaSpec.EventTypeUpdateComplete, omahaSpec.EventResultSuccessReboot, "0.0.0.0"))
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaEventResponse(t, omahaResp, tAppFlatcar.ID, 1)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, !addPing)
	checkOmahaUpdateResponse(t, omahaResp, expectedAppVersion, "", "", omahaSpec.NoUpdate)

	// Now updatetag, successful event, no previous version
	omahaResp = doOmahaRequest(t, h, tAppFlatcar.ID, expectedAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, ei(omahaSpec.EventTypeUpdateComplete, omahaSpec.EventResultSuccessReboot, ""))
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaEventResponse(t, omahaResp, tAppFlatcar.ID, 1)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, addPing)
	checkOmahaUpdateResponse(t, omahaResp, expectedAppVersion, "", "", omahaSpec.NoUpdate)

	// Now updatetag, successful event, with previous version
	omahaResp = doOmahaRequest(t, h, tAppFlatcar.ID, expectedAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, ei(omahaSpec.EventTypeUpdateComplete, omahaSpec.EventResultSuccessReboot, "614.0.0"))
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaEventResponse(t, omahaResp, tAppFlatcar.ID, 1)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, addPing)
	checkOmahaUpdateResponse(t, omahaResp, expectedAppVersion, "", "", omahaSpec.NoUpdate)

	// Now updatetag, successful event, with previous version, greater than current active version
	omahaResp = doOmahaRequest(t, h, tAppFlatcar.ID, "666.0.0", validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, ei(omahaSpec.EventTypeUpdateComplete, omahaSpec.EventResultSuccessReboot, "614.0.0"))
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaEventResponse(t, omahaResp, tAppFlatcar.ID, 1)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, addPing)
	checkOmahaUpdateResponse(t, omahaResp, expectedAppVersion, "", "", omahaSpec.NoUpdate)
}

func TestAppRegistrationForAppWithChannelAndPackageName(t *testing.T) {
	a := newForTest(t)
	defer a.Close()
	h := NewHandler(a)

	tAppFlatcar, _ := a.GetApp(flatcarAppID)
	tPkgFlatcar640, _ := a.AddPackage(&api.Package{Type: api.PkgTypeFlatcar, URL: "http://sample.url/pkg", Version: "640.0.0", ApplicationID: tAppFlatcar.ID, Arch: api.ArchAMD64})
	tChannel, _ := a.AddChannel(&api.Channel{Name: "mychannel", Color: "white", ApplicationID: tAppFlatcar.ID, PackageID: null.StringFrom(tPkgFlatcar640.ID), Arch: api.ArchAMD64})
	tGroup, _ := a.AddGroup(&api.Group{Name: "Production", ApplicationID: tAppFlatcar.ID, ChannelID: null.StringFrom(tChannel.ID), PolicyUpdatesEnabled: true, PolicySafeMode: true, PolicyPeriodInterval: "15 minutes", PolicyMaxUpdatesPerPeriod: 2, PolicyUpdateTimeout: "60 minutes"})

	validUnregisteredIP := "127.0.0.1"
	validUnregisteredMachineID := "65e1266d-6f54-4b87-9080-23b99ca9c12f"
	expectedAppVersion := "640.0.0"
	updateCheck := true
	eventPreviousVersion := ""
	addPing := true
	noEventInfo := (*eventInfo)(nil)

	omahaResp := doOmahaRequest(t, h, tAppFlatcar.ID, expectedAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, noEventInfo)
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, addPing)
	checkOmahaUpdateResponse(t, omahaResp, expectedAppVersion, "", "", omahaSpec.NoUpdate)

	omahaResp = doOmahaRequest(t, h, tAppFlatcar.ID, expectedAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, !updateCheck, ei(omahaSpec.EventTypeUpdateComplete, omahaSpec.EventResultSuccess, eventPreviousVersion))
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, addPing)
}

func TestAppUpdateForAppWithChannelAndPackageName(t *testing.T) {
	a := newForTest(t)
	defer a.Close()
	h := NewHandler(a)

	tAppFlatcar, _ := a.GetApp(flatcarAppID)
	tFilenameFlatcar := "flatcarupdate.tgz"
	tPkgFlatcar640, _ := a.AddPackage(&api.Package{Type: api.PkgTypeFlatcar, URL: "http://sample.url/pkg", Filename: null.StringFrom(tFilenameFlatcar), Version: "99640.0.0", ApplicationID: tAppFlatcar.ID, Arch: api.ArchAMD64})
	tChannel, _ := a.AddChannel(&api.Channel{Name: "mychannel", Color: "white", ApplicationID: tAppFlatcar.ID, PackageID: null.StringFrom(tPkgFlatcar640.ID), Arch: api.ArchAMD64})
	tGroup, _ := a.AddGroup(&api.Group{Name: "Production", ApplicationID: tAppFlatcar.ID, ChannelID: null.StringFrom(tChannel.ID), PolicyUpdatesEnabled: true, PolicySafeMode: true, PolicyPeriodInterval: "15 minutes", PolicyMaxUpdatesPerPeriod: 2, PolicyUpdateTimeout: "60 minutes"})
	flatcarAction, _ := a.AddFlatcarAction(&api.FlatcarAction{Event: "postinstall", Sha256: "fsdkjjfghsdakjfgaksdjfasd", PackageID: tPkgFlatcar640.ID})

	validUnregisteredIP := "127.0.0.1"
	validUnregisteredMachineID := "65e1266d-6f54-4b87-9080-23b99ca9c12f"
	oldAppVersion := "610.0.0"
	updateCheck := true
	addPing := true

	omahaResp := doOmahaRequest(t, h, tAppFlatcar.ID, oldAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, nil)
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaUpdateResponse(t, omahaResp, tPkgFlatcar640.Version, tFilenameFlatcar, tPkgFlatcar640.URL, omahaSpec.UpdateOK)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, addPing)
	checkOmahaFlatcarAction(t, flatcarAction, omahaResp.Apps[0].UpdateCheck.Manifest.Actions[0])

	// Send download started
	omahaResp = doOmahaRequest(t, h, tAppFlatcar.ID, oldAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, !updateCheck, ei(omahaSpec.EventTypeUpdateDownloadStarted, omahaSpec.EventResultSuccess, ""))
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, addPing)
	checkOmahaNoUpdateResponse(t, omahaResp)

	// Send download finished
	omahaResp = doOmahaRequest(t, h, tAppFlatcar.ID, oldAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, !updateCheck, ei(omahaSpec.EventTypeUpdateDownloadFinished, omahaSpec.EventResultSuccess, ""))
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, addPing)
	checkOmahaNoUpdateResponse(t, omahaResp)

	// Send complete
	omahaResp = doOmahaRequest(t, h, tAppFlatcar.ID, oldAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, !updateCheck, ei(omahaSpec.EventTypeUpdateComplete, omahaSpec.EventResultSuccess, ""))
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, addPing)
	checkOmahaNoUpdateResponse(t, omahaResp)

	// Send rebooted
	omahaResp = doOmahaRequest(t, h, tAppFlatcar.ID, tPkgFlatcar640.Version, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, ei(omahaSpec.EventTypeUpdateComplete, omahaSpec.EventResultSuccessReboot, oldAppVersion))
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, addPing)
	checkOmahaUpdateResponse(t, omahaResp, tPkgFlatcar640.Version, "", "", omahaSpec.NoUpdate)

	// Expect no update
	omahaResp = doOmahaRequest(t, h, tAppFlatcar.ID, tPkgFlatcar640.Version, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, nil)
	checkOmahaResponse(t, omahaResp, tAppFlatcar.ID, omahaSpec.AppOK)
	checkOmahaPingResponse(t, omahaResp, tAppFlatcar.ID, addPing)
	checkOmahaUpdateResponse(t, omahaResp, tPkgFlatcar640.Version, "", "", omahaSpec.NoUpdate)
}

func TestFlatcarGroupNamesConversionToIds(t *testing.T) {
	a := newForTest(t)
	defer a.Close()
	h := NewHandler(a)

	flatcarAppIDWithCurlyBraces := "{" + flatcarAppID + "}"
	machineID := "65e1266d-6f54-4b87-9080-23b99ca9c12f"
	machineIP := "10.0.0.1"

	omahaResp := doOmahaRequest(t, h, flatcarAppID, "2000.0.0", machineID, "invalid-group", machineIP, false, true, nil)
	checkOmahaResponse(t, omahaResp, flatcarAppID, omahaSpec.AppStatus("error-failedToRetrieveUpdatePackageInfo"))

	omahaResp = doOmahaRequest(t, h, flatcarAppID, "2000.0.0", machineID, "alpha", machineIP, false, true, nil)
	checkOmahaResponse(t, omahaResp, flatcarAppID, omahaSpec.AppOK)

	omahaResp = doOmahaRequest(t, h, flatcarAppIDWithCurlyBraces, "2000.0.0", machineID, "alpha", machineIP, false, true, nil)
	checkOmahaResponse(t, omahaResp, flatcarAppIDWithCurlyBraces, omahaSpec.AppOK)
}

func TestProductIDBasedRequest(t *testing.T) {
	a := newForTest(t)
	defer a.Close()
	h := NewHandler(a)

	tTeam, _ := a.AddTeam(&api.Team{Name: "test_team"})
	tApp, _ := a.AddApp(&api.Application{Name: "test_app", Description: "Test app", TeamID: tTeam.ID, ProductID: null.StringFrom("io.flatcar.MyApp")})
	tPkg, _ := a.AddPackage(&api.Package{Type: api.PkgTypeFlatcar, URL: "http://sample.url/pkg", Version: "640.0.0", ApplicationID: tApp.ID, Arch: api.ArchAMD64})
	tChannel, _ := a.AddChannel(&api.Channel{Name: "test_channel", Color: "blue", ApplicationID: tApp.ID, PackageID: null.StringFrom(tPkg.ID), Arch: api.ArchAMD64})
	tGroup, _ := a.AddGroup(&api.Group{Name: "test_group", ApplicationID: tApp.ID, ChannelID: null.StringFrom(tChannel.ID), PolicyUpdatesEnabled: true, PolicySafeMode: false, PolicyPeriodInterval: "15 minutes", PolicyMaxUpdatesPerPeriod: 4, PolicyUpdateTimeout: "60 minutes"})

	validAppProductID := "io.flatcar.MyApp"
	validUnregisteredIP := "127.0.0.1"
	validUnregisteredMachineID := "some-id"
	validUnverifiedAppVersion := "100.0.1"
	addPing := false
	updateCheck := true
	noEventInfo := (*eventInfo)(nil)

	omahaResp := doOmahaRequest(t, h, tApp.ID, validUnverifiedAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, noEventInfo)
	checkOmahaResponse(t, omahaResp, tApp.ID, omahaSpec.AppOK)

	omahaResp = doOmahaRequest(t, h, validAppProductID, validUnverifiedAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, noEventInfo)
	checkOmahaResponse(t, omahaResp, *tApp.ProductID.Ptr(), omahaSpec.AppOK)
}

func TestMultiPackageResponse(t *testing.T) {
	a := newForTest(t)
	defer a.Close()
	h := NewHandler(a)

	extraFiles := []api.File{
		{
			Name:    null.StringFrom("myfile1.txt"),
			Size:    null.StringFrom("1234"),
			Hash:    null.StringFrom("abcd"),
			Hash256: null.StringFrom("xyz"),
		},
		{
			Name: null.StringFrom("myfile2.txt"),
			Size: null.StringFrom("12345"),
			Hash: null.StringFrom("abcde"),
		},
	}

	tTeam, _ := a.AddTeam(&api.Team{Name: "test_team"})
	tApp, _ := a.AddApp(&api.Application{Name: "test_app", Description: "Test app", TeamID: tTeam.ID, ProductID: null.StringFrom("io.flatcar.MyApp")})
	tPkg, _ := a.AddPackage(&api.Package{Type: api.PkgTypeFlatcar, URL: "http://sample.url/pkg", Version: "640.0.0", ApplicationID: tApp.ID, Arch: api.ArchAMD64, ExtraFiles: extraFiles})
	tChannel, _ := a.AddChannel(&api.Channel{Name: "test_channel", Color: "blue", ApplicationID: tApp.ID, PackageID: null.StringFrom(tPkg.ID), Arch: api.ArchAMD64})
	tGroup, _ := a.AddGroup(&api.Group{Name: "test_group", ApplicationID: tApp.ID, ChannelID: null.StringFrom(tChannel.ID), PolicyUpdatesEnabled: true, PolicySafeMode: false, PolicyPeriodInterval: "15 minutes", PolicyMaxUpdatesPerPeriod: 100, PolicyUpdateTimeout: "60 minutes"})
	_, err := a.AddFlatcarAction(&api.FlatcarAction{Event: "postinstall", Sha256: "fsdkjjfghsdakjfgaksdjfasd", PackageID: tPkg.ID})

	assert.NoError(t, err)

	validUnregisteredIP := "127.0.0.1"
	validUnregisteredMachineID := "65e1266d-6f54-4b87-9080-23b99ca9c12f"
	oldAppVersion := "610.0.0"
	updateCheck := true
	addPing := true

	omahaResp := doOmahaRequest(t, h, tApp.ID, oldAppVersion, validUnregisteredMachineID, tGroup.ID, validUnregisteredIP, addPing, updateCheck, nil)
	checkOmahaResponse(t, omahaResp, tApp.ID, omahaSpec.AppOK)

	appResp := omahaResp.Apps[0]
	assert.NotNil(t, appResp)

	updateResp := appResp.UpdateCheck

	assert.NotNil(t, updateResp)
	assert.NotNil(t, updateResp.Manifest)

	// The packages shipped in the response are the main one + the extras
	packages := updateResp.Manifest.Packages
	assert.Equal(t, len(extraFiles)+1, len(packages))

	extraFile1 := extraFiles[0]
	shippedExtraFile1 := packages[1]

	assert.Equal(t, extraFile1.Name.String, shippedExtraFile1.Name)

	extraFile1Size, _ := strconv.ParseUint(extraFile1.Size.String, 10, 64)
	assert.Equal(t, extraFile1Size, shippedExtraFile1.Size)

	assert.Equal(t, extraFile1.Hash.String, shippedExtraFile1.SHA1)
	assert.Equal(t, extraFile1.Hash256.String, shippedExtraFile1.SHA256)

	extraFile2 := extraFiles[1]
	shippedExtraFile2 := packages[2]

	assert.Equal(t, extraFile2.Name.String, shippedExtraFile2.Name)

	mainFile := packages[0]
	assert.Equal(t, tPkg.Filename.String, mainFile.Name)
}

type eventInfo struct {
	Type            omahaSpec.EventType
	Result          omahaSpec.EventResult
	PreviousVersion string
}

func ei(t omahaSpec.EventType, r omahaSpec.EventResult, pv string) *eventInfo {
	return &eventInfo{
		Type:            t,
		Result:          r,
		PreviousVersion: pv,
	}
}

func doOmahaRequest(t *testing.T, h *Handler, appID, appVersion, appMachineID, appTrack, ip string, addPing, updateCheck bool, eventInfo *eventInfo) *omahaSpec.Response {
	omahaReq := omahaSpec.NewRequest()
	omahaReq.OS.Version = reqVersion
	omahaReq.OS.Platform = reqPlatform
	omahaReq.OS.ServicePack = reqSp
	omahaReq.OS.Arch = reqArch
	appReq := omahaReq.AddApp(appID, appVersion)
	appReq.MachineID = appMachineID
	appReq.Track = appTrack
	if updateCheck {
		appReq.AddUpdateCheck()
	}
	if eventInfo != nil {
		eReq := appReq.AddEvent()
		eReq.Type = eventInfo.Type
		eReq.Result = eventInfo.Result
		eReq.PreviousVersion = eventInfo.PreviousVersion
	}
	if addPing {
		appReq.AddPing()
	}

	omahaReqXML, err := xml.Marshal(omahaReq)
	assert.NoError(t, err)

	omahaRespXML := new(bytes.Buffer)
	err = h.Handle(bytes.NewReader(omahaReqXML), omahaRespXML, ip)
	assert.NoError(t, err)

	var omahaResp *omahaSpec.Response
	err = xml.NewDecoder(omahaRespXML).Decode(&omahaResp)
	assert.NoError(t, err)

	return omahaResp
}

func checkOmahaResponse(t *testing.T, omahaResp *omahaSpec.Response, expectedAppID string, expectedError omahaSpec.AppStatus) {
	appResp := omahaResp.Apps[0]

	assert.Equal(t, expectedError, appResp.Status)
	assert.Equal(t, expectedAppID, appResp.ID)
}

func checkOmahaNoUpdateResponse(t *testing.T, omahaResp *omahaSpec.Response) {
	appResp := omahaResp.Apps[0]

	assert.Nil(t, appResp.UpdateCheck)
}

func checkOmahaUpdateResponse(t *testing.T, omahaResp *omahaSpec.Response, expectedVersion, expectedPackageName, expectedUpdateURL string, expectedError omahaSpec.UpdateStatus) {
	appResp := omahaResp.Apps[0]

	assert.NotNil(t, appResp.UpdateCheck)
	assert.Equal(t, expectedError, appResp.UpdateCheck.Status)

	if appResp.UpdateCheck.Manifest != nil {
		assert.True(t, appResp.UpdateCheck.Manifest.Version >= expectedVersion)
		assert.Equal(t, expectedPackageName, appResp.UpdateCheck.Manifest.Packages[0].Name)
	}

	if appResp.UpdateCheck.URLs != nil {
		assert.Equal(t, 1, len(appResp.UpdateCheck.URLs))
		assert.Equal(t, expectedUpdateURL, appResp.UpdateCheck.URLs[0].CodeBase)
	}
}

func checkOmahaEventResponse(t *testing.T, omahaResp *omahaSpec.Response, expectedAppID string, expectedEventCount int) {
	appResp := omahaResp.Apps[0]

	assert.Equal(t, expectedAppID, appResp.ID)
	assert.Equal(t, expectedEventCount, len(appResp.Events))
	for i := 0; i < expectedEventCount; i++ {
		assert.Equal(t, "ok", appResp.Events[i].Status)
	}
}

func checkOmahaPingResponse(t *testing.T, omahaResp *omahaSpec.Response, expectedAppID string, expectedPingResponse bool) {
	appResp := omahaResp.Apps[0]

	assert.Equal(t, expectedAppID, appResp.ID)
	if expectedPingResponse {
		assert.Equal(t, "ok", appResp.Ping.Status)
		assert.NotNil(t, appResp.Ping)
	} else {
		assert.Nil(t, appResp.Ping)
	}
}

func checkOmahaFlatcarAction(t *testing.T, c *api.FlatcarAction, r *omahaSpec.Action) {
	assert.Equal(t, c.Event, r.Event)
	assert.Equal(t, c.Sha256, r.SHA256)
	assert.Equal(t, c.IsDelta, r.IsDeltaPayload)
	assert.Equal(t, c.Deadline, r.Deadline)
	assert.Equal(t, c.DisablePayloadBackoff, r.DisablePayloadBackoff)
	assert.Equal(t, c.ChromeOSVersion, r.DisplayVersion)
	assert.Equal(t, c.MetadataSize, r.MetadataSize)
	assert.Equal(t, c.NeedsAdmin, r.NeedsAdmin)
	assert.Equal(t, c.MetadataSignatureRsa, r.MetadataSignatureRsa)
}

func Test_getArch(t *testing.T) {
	type args struct {
		os     *omahaSpec.OS
		appReq *omahaSpec.AppRequest
	}
	tests := []struct {
		name string
		args args
		want api.Arch
	}{
		{
			name: "all",
			args: args{
				os: &omahaSpec.OS{
					Arch: "all",
				},
				appReq: nil,
			},
			want: api.ArchAll,
		},
		{
			name: "empty",
			args: args{
				os: &omahaSpec.OS{
					Arch: "",
				},
				appReq: &omahaSpec.AppRequest{},
			},
			want: api.ArchAMD64,
		},
		{
			name: "coreOS amd64",
			args: args{
				os: nil,
				appReq: &omahaSpec.AppRequest{
					Board: "amd64-usr",
				},
			},
			want: api.ArchAMD64,
		},
		{
			name: "coreOS arm64",
			args: args{
				os: &omahaSpec.OS{
					Arch: "",
				},
				appReq: &omahaSpec.AppRequest{
					Board: "arm64-usr",
				},
			},
			want: api.ArchAArch64,
		},
		{
			name: "flatpak amd64",
			args: args{
				os: &omahaSpec.OS{
					Arch: "x64",
				},
				appReq: &omahaSpec.AppRequest{},
			},
			want: api.ArchAMD64,
		},
		{
			name: "flatpak arm64",
			args: args{
				os: &omahaSpec.OS{
					Arch: "arm",
				},
				appReq: &omahaSpec.AppRequest{},
			},
			want: api.ArchAArch64,
		},
		{
			name: "flatpak x86",
			args: args{
				os: &omahaSpec.OS{
					Arch: "x86",
				},
				appReq: &omahaSpec.AppRequest{},
			},
			want: api.ArchX86,
		},
		{
			name: "nil",
			args: args{
				os:     nil,
				appReq: nil,
			},
			want: api.ArchAMD64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getArch(tt.args.os, tt.args.appReq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getArch() = %v, want %v", got, tt.want)
			}
		})
	}
}
