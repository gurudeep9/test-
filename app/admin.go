// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"io"
	"os"
	"strings"
	"time"

	"runtime/debug"

	"net/http"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/services/mailservice"
	"github.com/mattermost/mattermost-server/utils"
)

func (a *App) GetLogs(page, perPage int) ([]string, *model.AppError) {
	var lines []string
	if a.Cluster != nil && *a.Config().ClusterSettings.Enable {
		lines = append(lines, "-----------------------------------------------------------------------------------------------------------")
		lines = append(lines, "-----------------------------------------------------------------------------------------------------------")
		lines = append(lines, a.Cluster.GetMyClusterInfo().Hostname)
		lines = append(lines, "-----------------------------------------------------------------------------------------------------------")
		lines = append(lines, "-----------------------------------------------------------------------------------------------------------")
	}

	melines, err := a.GetLogsSkipSend(page, perPage)
	if err != nil {
		return nil, err
	}

	lines = append(lines, melines...)

	if a.Cluster != nil && *a.Config().ClusterSettings.Enable {
		clines, err := a.Cluster.GetLogs(page, perPage)
		if err != nil {
			return nil, err
		}

		lines = append(lines, clines...)
	}

	return lines, nil
}

func (a *App) GetLogsSkipSend(page, perPage int) ([]string, *model.AppError) {
	var lines []string

	if a.Config().LogSettings.EnableFile {
		file, err := os.Open(utils.GetLogFileLocation(a.Config().LogSettings.FileLocation))
		if err != nil {
			return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
		}

		defer file.Close()

		var newLine = []byte{'\n'}
		var lineCount int
		const searchPos = -1
		lineEndPos, err := file.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
		}
		for {
			pos, err := file.Seek(searchPos, io.SeekCurrent)
			if err != nil {
				return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
			}

			b := make([]byte, 1)
			_, err = file.ReadAt(b, pos)
			if err != nil {
				return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
			}

			if b[0] == newLine[0] || pos == 0 {
				lineCount++
				if lineCount > page*perPage {
					line := make([]byte, lineEndPos-pos)
					_, err := file.ReadAt(line, pos)
					if err != nil {
						return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
					}
					lines = append(lines, string(line))
				}
				if pos == 0 {
					break
				}
				lineEndPos = pos
			}

			if len(lines) == perPage {
				break
			}
		}

		for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
			lines[i], lines[j] = lines[j], lines[i]
		}
	} else {
		lines = append(lines, "")
	}

	return lines, nil
}

func (a *App) GetClusterStatus() []*model.ClusterInfo {
	infos := make([]*model.ClusterInfo, 0)

	if a.Cluster != nil {
		infos = a.Cluster.GetClusterInfos()
	}

	return infos
}

func (a *App) InvalidateAllCaches() *model.AppError {
	debug.FreeOSMemory()
	a.InvalidateAllCachesSkipSend()

	if a.Cluster != nil {

		msg := &model.ClusterMessage{
			Event:            model.CLUSTER_EVENT_INVALIDATE_ALL_CACHES,
			SendType:         model.CLUSTER_SEND_RELIABLE,
			WaitForAllToSend: true,
		}

		a.Cluster.SendClusterMessage(msg)
	}

	return nil
}

func (a *App) InvalidateAllCachesSkipSend() {
	mlog.Info("Purging all caches")
	a.sessionCache.Purge()
	ClearStatusCache()
	a.Srv.Store.Channel().ClearCaches()
	a.Srv.Store.User().ClearCaches()
	a.Srv.Store.Post().ClearCaches()
	a.Srv.Store.FileInfo().ClearCaches()
	a.Srv.Store.Webhook().ClearCaches()
	a.LoadLicense()
}

func (a *App) GetConfig() *model.Config {
	json := a.Config().ToJson()
	cfg := model.ConfigFromJson(strings.NewReader(json))
	cfg.Sanitize()

	return cfg
}

func (a *App) GetEnvironmentConfig() map[string]interface{} {
	return a.EnvironmentConfig()
}

func (a *App) SaveConfig(cfg *model.Config, sendConfigChangeClusterMessage bool, userId ...string) *model.AppError {
	oldCfg := a.Config()
	cfg.SetDefaults()
	a.Desanitize(cfg)

	if err := cfg.IsValid(); err != nil {
		return err
	}

	if err := utils.ValidateLdapFilter(cfg, a.Ldap); err != nil {
		return err
	}

	if *a.Config().ClusterSettings.Enable && *a.Config().ClusterSettings.ReadOnlyConfig {
		return model.NewAppError("saveConfig", "ent.cluster.save_config.error", nil, "", http.StatusForbidden)
	}

	a.DisableConfigWatch()

	if *oldCfg.SupportSettings.CustomServiceTermsText != *cfg.SupportSettings.CustomServiceTermsText {
		if len(userId) != 1 {
			return model.NewAppError("saveConfig", "ent.cluster.save_config.update_custom_service_terms_no_user.error", nil, "", http.StatusBadRequest)
		}
		if serviceTerms, err := a.CreateServiceTerms(*cfg.SupportSettings.CustomServiceTermsText, userId[0]); err != nil {
			return err
		} else {
			cfg.SupportSettings.CustomServiceTermsId = model.NewString(serviceTerms.Id)
		}
	}

	// no need to temperorily store the text as it will be the as received form API
	cfg.SupportSettings.CustomServiceTermsText = model.NewString("")

	// service terms ID however changes on saving new text. So we need to update it
	// but not save it on config.json.
	serviceTermsId := cfg.SupportSettings.CustomServiceTermsId
	cfg.SupportSettings.CustomServiceTermsId = model.NewString("")

	a.UpdateConfig(func(update *model.Config) {
		*update = *cfg
	})
	a.PersistConfig()
	a.ReloadConfig()
	a.EnableConfigWatch()

	cfg.SupportSettings.CustomServiceTermsId = serviceTermsId

	if a.Metrics != nil {
		if *a.Config().MetricsSettings.Enable {
			a.Metrics.StartServer()
		} else {
			a.Metrics.StopServer()
		}
	}

	if a.Cluster != nil {
		err := a.Cluster.ConfigChanged(cfg, oldCfg, sendConfigChangeClusterMessage)
		if err != nil {
			return err
		}
	}

	// start/restart email batching job if necessary
	a.InitEmailBatching()

	return nil
}

func (a *App) RecycleDatabaseConnection() {
	oldStore := a.Srv.Store

	mlog.Warn("Attempting to recycle the database connection.")
	a.Srv.Store = a.newStore()
	a.Jobs.Store = a.Srv.Store

	if a.Srv.Store != oldStore {
		time.Sleep(20 * time.Second)
		oldStore.Close()
	}

	mlog.Warn("Finished recycling the database connection.")
}

func (a *App) TestEmail(userId string, cfg *model.Config) *model.AppError {
	if len(cfg.EmailSettings.SMTPServer) == 0 {
		return model.NewAppError("testEmail", "api.admin.test_email.missing_server", nil, utils.T("api.context.invalid_param.app_error", map[string]interface{}{"Name": "SMTPServer"}), http.StatusBadRequest)
	}

	// if the user hasn't changed their email settings, fill in the actual SMTP password so that
	// the user can verify an existing SMTP connection
	if cfg.EmailSettings.SMTPPassword == model.FAKE_SETTING {
		if cfg.EmailSettings.SMTPServer == a.Config().EmailSettings.SMTPServer &&
			cfg.EmailSettings.SMTPPort == a.Config().EmailSettings.SMTPPort &&
			cfg.EmailSettings.SMTPUsername == a.Config().EmailSettings.SMTPUsername {
			cfg.EmailSettings.SMTPPassword = a.Config().EmailSettings.SMTPPassword
		} else {
			return model.NewAppError("testEmail", "api.admin.test_email.reenter_password", nil, "", http.StatusBadRequest)
		}
	}
	user, err := a.GetUser(userId)
	if err != nil {
		return err
	}

	T := utils.GetUserTranslations(user.Locale)
	license := a.License()
	if err := mailservice.SendMailUsingConfig(user.Email, T("api.admin.test_email.subject"), T("api.admin.test_email.body"), cfg, license != nil && *license.Features.Compliance); err != nil {
		return model.NewAppError("testEmail", "app.admin.test_email.failure", map[string]interface{}{"Error": err.Error()}, "", http.StatusInternalServerError)
	}

	return nil
}
