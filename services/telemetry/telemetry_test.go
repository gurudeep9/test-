// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package telemetry

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/mattermost/mattermost-server/v5/services/httpservice"
	"github.com/mattermost/mattermost-server/v5/services/searchengine"
	"github.com/mattermost/mattermost-server/v5/services/telemetry/mocks"
	storeMocks "github.com/mattermost/mattermost-server/v5/store/storetest/mocks"
)

type FakeConfigService struct {
	cfg *model.Config
}

func (fcs *FakeConfigService) Config() *model.Config                                       { return fcs.cfg }
func (fcs *FakeConfigService) AddConfigListener(f func(old, current *model.Config)) string { return "" }
func (fcs *FakeConfigService) RemoveConfigListener(key string)                             {}
func (fcs *FakeConfigService) AsymmetricSigningKey() *ecdsa.PrivateKey                     { return nil }

func initializeMocks(cfg *model.Config) (*mocks.ServerIface, *storeMocks.Store, func(t *testing.T), func()) {
	serverIfaceMock := &mocks.ServerIface{}

	configService := &FakeConfigService{cfg}
	serverIfaceMock.On("Config").Return(cfg)
	serverIfaceMock.On("IsLeader").Return(true)

	pluginDir, _ := ioutil.TempDir("", "")
	webappPluginDir, _ := ioutil.TempDir("", "")
	cleanUp := func() {
		os.RemoveAll(pluginDir)
		os.RemoveAll(webappPluginDir)
	}
	pluginsAPIMock := &plugintest.API{}
	pluginEnv, _ := plugin.NewEnvironment(func(m *model.Manifest) plugin.API { return pluginsAPIMock }, pluginDir, webappPluginDir, mlog.NewLogger(&mlog.LoggerConfiguration{}), nil)
	serverIfaceMock.On("GetPluginsEnvironment").Return(pluginEnv, nil)

	serverIfaceMock.On("License").Return(model.NewTestLicense(), nil)
	serverIfaceMock.On("GetRoleByName", "system_admin").Return(&model.Role{Permissions: []string{"sa-test1", "sa-test2"}}, nil)
	serverIfaceMock.On("GetRoleByName", "system_user").Return(&model.Role{Permissions: []string{"su-test1", "su-test2"}}, nil)
	serverIfaceMock.On("GetRoleByName", "system_user_manager").Return(&model.Role{Permissions: []string{"sum-test1", "sum-test2"}}, nil)
	serverIfaceMock.On("GetRoleByName", "system_manager").Return(&model.Role{Permissions: []string{"sm-test1", "sm-test2"}}, nil)
	serverIfaceMock.On("GetRoleByName", "system_read_only_admin").Return(&model.Role{Permissions: []string{"sra-test1", "sra-test2"}}, nil)
	serverIfaceMock.On("GetRoleByName", "team_admin").Return(&model.Role{Permissions: []string{"ta-test1", "ta-test2"}}, nil)
	serverIfaceMock.On("GetRoleByName", "team_user").Return(&model.Role{Permissions: []string{"tu-test1", "tu-test2"}}, nil)
	serverIfaceMock.On("GetRoleByName", "team_guest").Return(&model.Role{Permissions: []string{"tg-test1", "tg-test2"}}, nil)
	serverIfaceMock.On("GetRoleByName", "channel_admin").Return(&model.Role{Permissions: []string{"ca-test1", "ca-test2"}}, nil)
	serverIfaceMock.On("GetRoleByName", "channel_user").Return(&model.Role{Permissions: []string{"cu-test1", "cu-test2"}}, nil)
	serverIfaceMock.On("GetRoleByName", "channel_guest").Return(&model.Role{Permissions: []string{"cg-test1", "cg-test2"}}, nil)
	serverIfaceMock.On("GetSchemes", "team", 0, 100).Return([]*model.Scheme{}, nil)
	serverIfaceMock.On("HttpService").Return(httpservice.MakeHTTPService(configService))

	storeMock := &storeMocks.Store{}
	storeMock.On("GetDbVersion").Return("5.24.0", nil)

	systemStore := storeMocks.SystemStore{}
	props := model.StringMap{}
	props[model.SYSTEM_TELEMETRY_ID] = "test"
	systemStore.On("Get").Return(props, nil)
	systemStore.On("GetByName", model.ADVANCED_PERMISSIONS_MIGRATION_KEY).Return(nil, nil)
	systemStore.On("GetByName", model.MIGRATION_KEY_ADVANCED_PERMISSIONS_PHASE_2).Return(nil, nil)

	userStore := storeMocks.UserStore{}
	userStore.On("Count", model.UserCountOptions{IncludeBotAccounts: false, IncludeDeleted: true, ExcludeRegularUsers: false, TeamId: "", ViewRestrictions: nil}).Return(int64(10), nil)
	userStore.On("Count", model.UserCountOptions{IncludeBotAccounts: true, IncludeDeleted: false, ExcludeRegularUsers: true, TeamId: "", ViewRestrictions: nil}).Return(int64(100), nil)
	userStore.On("Count", model.UserCountOptions{Roles: []string{model.SYSTEM_MANAGER_ROLE_ID}}).Return(int64(5), nil)
	userStore.On("Count", model.UserCountOptions{Roles: []string{model.SYSTEM_USER_MANAGER_ROLE_ID}}).Return(int64(10), nil)
	userStore.On("Count", model.UserCountOptions{Roles: []string{model.SYSTEM_READ_ONLY_ADMIN_ROLE_ID}}).Return(int64(15), nil)
	userStore.On("AnalyticsGetGuestCount").Return(int64(11), nil)
	userStore.On("AnalyticsActiveCount", mock.Anything, model.UserCountOptions{IncludeBotAccounts: false, IncludeDeleted: false, ExcludeRegularUsers: false, TeamId: "", ViewRestrictions: nil}).Return(int64(5), nil)
	userStore.On("AnalyticsGetInactiveUsersCount").Return(int64(8), nil)
	userStore.On("AnalyticsGetSystemAdminCount").Return(int64(9), nil)

	teamStore := storeMocks.TeamStore{}
	teamStore.On("AnalyticsTeamCount", false).Return(int64(3), nil)
	teamStore.On("GroupSyncedTeamCount").Return(int64(16), nil)

	channelStore := storeMocks.ChannelStore{}
	channelStore.On("AnalyticsTypeCount", "", "O").Return(int64(25), nil)
	channelStore.On("AnalyticsTypeCount", "", "P").Return(int64(26), nil)
	channelStore.On("AnalyticsTypeCount", "", "D").Return(int64(27), nil)
	channelStore.On("AnalyticsDeletedTypeCount", "", "O").Return(int64(22), nil)
	channelStore.On("AnalyticsDeletedTypeCount", "", "P").Return(int64(23), nil)
	channelStore.On("GroupSyncedChannelCount").Return(int64(17), nil)

	postStore := storeMocks.PostStore{}
	postStore.On("AnalyticsPostCount", "", false, false).Return(int64(1000), nil)
	postStore.On("AnalyticsPostCountsByDay", &model.AnalyticsPostCountsOptions{TeamId: "", BotsOnly: false, YesterdayOnly: true}).Return(model.AnalyticsRows{}, nil)
	postStore.On("AnalyticsPostCountsByDay", &model.AnalyticsPostCountsOptions{TeamId: "", BotsOnly: true, YesterdayOnly: true}).Return(model.AnalyticsRows{}, nil)

	commandStore := storeMocks.CommandStore{}
	commandStore.On("AnalyticsCommandCount", "").Return(int64(15), nil)

	webhookStore := storeMocks.WebhookStore{}
	webhookStore.On("AnalyticsIncomingCount", "").Return(int64(16), nil)
	webhookStore.On("AnalyticsOutgoingCount", "").Return(int64(17), nil)

	groupStore := storeMocks.GroupStore{}
	groupStore.On("GroupCount").Return(int64(25), nil)
	groupStore.On("GroupTeamCount").Return(int64(26), nil)
	groupStore.On("GroupChannelCount").Return(int64(27), nil)
	groupStore.On("GroupMemberCount").Return(int64(32), nil)
	groupStore.On("DistinctGroupMemberCount").Return(int64(22), nil)
	groupStore.On("GroupCountWithAllowReference").Return(int64(13), nil)

	schemeStore := storeMocks.SchemeStore{}
	schemeStore.On("CountByScope", "channel").Return(int64(8), nil)
	schemeStore.On("CountByScope", "team").Return(int64(7), nil)
	schemeStore.On("CountWithoutPermission", "channel", "create_post", model.RoleScopeChannel, model.RoleTypeUser).Return(int64(6), nil)
	schemeStore.On("CountWithoutPermission", "channel", "create_post", model.RoleScopeChannel, model.RoleTypeGuest).Return(int64(7), nil)
	schemeStore.On("CountWithoutPermission", "channel", "add_reaction", model.RoleScopeChannel, model.RoleTypeUser).Return(int64(8), nil)
	schemeStore.On("CountWithoutPermission", "channel", "add_reaction", model.RoleScopeChannel, model.RoleTypeGuest).Return(int64(9), nil)
	schemeStore.On("CountWithoutPermission", "channel", "manage_public_channel_members", model.RoleScopeChannel, model.RoleTypeUser).Return(int64(10), nil)
	schemeStore.On("CountWithoutPermission", "channel", "use_channel_mentions", model.RoleScopeChannel, model.RoleTypeUser).Return(int64(11), nil)
	schemeStore.On("CountWithoutPermission", "channel", "use_channel_mentions", model.RoleScopeChannel, model.RoleTypeGuest).Return(int64(12), nil)

	storeMock.On("System").Return(&systemStore)
	storeMock.On("User").Return(&userStore)
	storeMock.On("Team").Return(&teamStore)
	storeMock.On("Channel").Return(&channelStore)
	storeMock.On("Post").Return(&postStore)
	storeMock.On("Command").Return(&commandStore)
	storeMock.On("Webhook").Return(&webhookStore)
	storeMock.On("Group").Return(&groupStore)
	storeMock.On("Scheme").Return(&schemeStore)

	return serverIfaceMock, storeMock, func(t *testing.T) {
		serverIfaceMock.AssertExpectations(t)
		storeMock.AssertExpectations(t)
		systemStore.AssertExpectations(t)
		pluginsAPIMock.AssertExpectations(t)
	}, cleanUp
}

func TestPluginSetting(t *testing.T) {
	settings := &model.PluginSettings{
		Plugins: map[string]map[string]interface{}{
			"test": {
				"foo": "bar",
			},
		},
	}
	assert.Equal(t, "bar", pluginSetting(settings, "test", "foo", "asd"))
	assert.Equal(t, "asd", pluginSetting(settings, "test", "qwe", "asd"))
}

func TestPluginActivated(t *testing.T) {
	states := map[string]*model.PluginState{
		"foo": {
			Enable: true,
		},
		"bar": {
			Enable: false,
		},
	}
	assert.True(t, pluginActivated(states, "foo"))
	assert.False(t, pluginActivated(states, "bar"))
	assert.False(t, pluginActivated(states, "none"))
}

func TestPluginVersion(t *testing.T) {
	plugins := []*model.BundleInfo{
		{
			Manifest: &model.Manifest{
				Id:      "test.plugin",
				Version: "1.2.3",
			},
		},
		{
			Manifest: &model.Manifest{
				Id:      "test.plugin2",
				Version: "4.5.6",
			},
		},
	}
	assert.Equal(t, "1.2.3", pluginVersion(plugins, "test.plugin"))
	assert.Equal(t, "4.5.6", pluginVersion(plugins, "test.plugin2"))
	assert.Empty(t, pluginVersion(plugins, "unknown.plugin"))
}

func TestRudderTelemetry(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	type batch struct {
		MessageId  string
		UserId     string
		Event      string
		Timestamp  time.Time
		Properties map[string]interface{}
	}

	type payload struct {
		MessageId string
		SentAt    time.Time
		Batch     []struct {
			MessageId  string
			UserId     string
			Event      string
			Timestamp  time.Time
			Properties map[string]interface{}
		}
		Context struct {
			Library struct {
				Name    string
				Version string
			}
		}
	}

	data := make(chan payload, 100)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var p payload
		err = json.Unmarshal(body, &p)
		require.NoError(t, err)

		data <- p
	}))
	defer server.Close()

	marketplaceServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		json, err := json.Marshal([]*model.MarketplacePlugin{{
			BaseMarketplacePlugin: &model.BaseMarketplacePlugin{
				Manifest: &model.Manifest{
					Id: "testplugin",
				},
			},
		}})
		require.NoError(t, err)
		res.Write(json)
	}))

	defer func() { marketplaceServer.Close() }()

	telemetryID := "test-telemetry-id-12345"

	cfg := &model.Config{}
	cfg.SetDefaults()
	serverIfaceMock, storeMock, deferredAssertions, cleanUp := initializeMocks(cfg)
	defer cleanUp()
	defer deferredAssertions(t)

	telemetryService := New(serverIfaceMock, storeMock, searchengine.NewBroker(cfg, nil), mlog.NewLogger(&mlog.LoggerConfiguration{}))
	telemetryService.TelemetryID = telemetryID
	telemetryService.rudderClient = nil
	telemetryService.initRudder(server.URL, "")

	assertPayload := func(t *testing.T, actual payload, event string, properties map[string]interface{}) {
		t.Helper()
		assert.NotEmpty(t, actual.MessageId)
		assert.False(t, actual.SentAt.IsZero())
		if assert.Len(t, actual.Batch, 1) {
			assert.NotEmpty(t, actual.Batch[0].MessageId, "message id should not be empty")
			assert.Equal(t, telemetryID, actual.Batch[0].UserId)
			if event != "" {
				assert.Equal(t, event, actual.Batch[0].Event)
			}
			assert.False(t, actual.Batch[0].Timestamp.IsZero(), "batch timestamp should not be the zero value")
			if properties != nil {
				assert.Equal(t, properties, actual.Batch[0].Properties)
			}
		}
		assert.Equal(t, "analytics-go", actual.Context.Library.Name)
		assert.Equal(t, "3.0.0", actual.Context.Library.Version)
	}

	collectInfo := func(info *[]string) {
		t.Helper()
		for {
			select {
			case result := <-data:
				assertPayload(t, result, "", nil)
				*info = append(*info, result.Batch[0].Event)
			case <-time.After(time.Second * 1):
				return
			}
		}
	}

	collectBatches := func(info *[]batch) {
		t.Helper()
		for {
			select {
			case result := <-data:
				assertPayload(t, result, "", nil)
				*info = append(*info, result.Batch[0])
			case <-time.After(time.Second * 1):
				return
			}
		}
	}

	// Should send a client identify message
	select {
	case identifyMessage := <-data:
		assertPayload(t, identifyMessage, "", nil)
	case <-time.After(time.Second * 1):
		require.Fail(t, "Did not receive ID message")
	}

	t.Run("Send", func(t *testing.T) {
		testValue := "test-send-value-6789"
		telemetryService.sendTelemetry("Testing Telemetry", map[string]interface{}{
			"hey": testValue,
		})
		select {
		case result := <-data:
			assertPayload(t, result, "Testing Telemetry", map[string]interface{}{
				"hey": testValue,
			})
		case <-time.After(time.Second * 1):
			require.Fail(t, "Did not receive telemetry")
		}
	})

	// Plugins remain disabled at this point
	t.Run("SendDailyTelemetryPluginsDisabled", func(t *testing.T) {
		telemetryService.sendDailyTelemetry(true)

		var info []string
		// Collect the info sent.
		collectInfo(&info)

		for _, item := range []string{
			TRACK_CONFIG_SERVICE,
			TRACK_CONFIG_TEAM,
			TRACK_CONFIG_SQL,
			TRACK_CONFIG_LOG,
			TRACK_CONFIG_NOTIFICATION_LOG,
			TRACK_CONFIG_FILE,
			TRACK_CONFIG_RATE,
			TRACK_CONFIG_EMAIL,
			TRACK_CONFIG_PRIVACY,
			TRACK_CONFIG_OAUTH,
			TRACK_CONFIG_LDAP,
			TRACK_CONFIG_COMPLIANCE,
			TRACK_CONFIG_LOCALIZATION,
			TRACK_CONFIG_SAML,
			TRACK_CONFIG_PASSWORD,
			TRACK_CONFIG_CLUSTER,
			TRACK_CONFIG_METRICS,
			TRACK_CONFIG_SUPPORT,
			TRACK_CONFIG_NATIVEAPP,
			TRACK_CONFIG_EXPERIMENTAL,
			TRACK_CONFIG_ANALYTICS,
			TRACK_CONFIG_PLUGIN,
			TRACK_ACTIVITY,
			TRACK_SERVER,
			TRACK_CONFIG_MESSAGE_EXPORT,
			TRACK_PLUGINS,
		} {
			require.Contains(t, info, item)
		}
	})

	// Enable plugins for the remainder of the tests.
	// th.Server.UpdateConfig(func(cfg *model.Config) { *cfg.PluginSettings.Enable = true })

	t.Run("SendDailyTelemetry", func(t *testing.T) {
		telemetryService.sendDailyTelemetry(true)

		var info []string
		// Collect the info sent.
		collectInfo(&info)

		for _, item := range []string{
			TRACK_CONFIG_SERVICE,
			TRACK_CONFIG_TEAM,
			TRACK_CONFIG_SQL,
			TRACK_CONFIG_LOG,
			TRACK_CONFIG_NOTIFICATION_LOG,
			TRACK_CONFIG_FILE,
			TRACK_CONFIG_RATE,
			TRACK_CONFIG_EMAIL,
			TRACK_CONFIG_PRIVACY,
			TRACK_CONFIG_OAUTH,
			TRACK_CONFIG_LDAP,
			TRACK_CONFIG_COMPLIANCE,
			TRACK_CONFIG_LOCALIZATION,
			TRACK_CONFIG_SAML,
			TRACK_CONFIG_PASSWORD,
			TRACK_CONFIG_CLUSTER,
			TRACK_CONFIG_METRICS,
			TRACK_CONFIG_SUPPORT,
			TRACK_CONFIG_NATIVEAPP,
			TRACK_CONFIG_EXPERIMENTAL,
			TRACK_CONFIG_ANALYTICS,
			TRACK_CONFIG_PLUGIN,
			TRACK_ACTIVITY,
			TRACK_SERVER,
			TRACK_CONFIG_MESSAGE_EXPORT,
			TRACK_PLUGINS,
		} {
			require.Contains(t, info, item)
		}
	})
	t.Run("Telemetry for Marketplace plugins is returned", func(t *testing.T) {
		telemetryService.trackPluginConfig(telemetryService.srv.Config(), marketplaceServer.URL)

		var batches []batch
		collectBatches(&batches)

		for _, b := range batches {
			if b.Event == TRACK_CONFIG_PLUGIN {
				assert.Contains(t, b.Properties, "enable_testplugin")
				assert.Contains(t, b.Properties, "version_testplugin")

				// Confirm known plugins are not present
				assert.NotContains(t, b.Properties, "enable_jira")
				assert.NotContains(t, b.Properties, "version_jira")
			}
		}
	})

	t.Run("Telemetry for known plugins is returned, if request to Marketplace fails", func(t *testing.T) {
		telemetryService.trackPluginConfig(telemetryService.srv.Config(), "http://some.random.invalid.url")

		var batches []batch
		collectBatches(&batches)

		for _, b := range batches {
			if b.Event == TRACK_CONFIG_PLUGIN {
				assert.NotContains(t, b.Properties, "enable_testplugin")
				assert.NotContains(t, b.Properties, "version_testplugin")

				// Confirm known plugins are present
				assert.Contains(t, b.Properties, "enable_jira")
				assert.Contains(t, b.Properties, "version_jira")
			}
		}
	})

	t.Run("SendDailyTelemetryNoRudderKey", func(t *testing.T) {
		telemetryService.sendDailyTelemetry(false)

		select {
		case <-data:
			require.Fail(t, "Should not send telemetry when the rudder key is not set")
		case <-time.After(time.Second * 1):
			// Did not receive telemetry
		}
	})

	t.Run("SendDailyTelemetryDisabled", func(t *testing.T) {
		*cfg.LogSettings.EnableDiagnostics = false
		defer func() {
			*cfg.LogSettings.EnableDiagnostics = true
		}()

		telemetryService.sendDailyTelemetry(true)

		select {
		case <-data:
			require.Fail(t, "Should not send telemetry when they are disabled")
		case <-time.After(time.Second * 1):
			// Did not receive telemetry
		}
	})
}
