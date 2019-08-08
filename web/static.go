// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package web

import (
	"fmt"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/avct/uasurfer"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/utils"
	"github.com/mattermost/mattermost-server/utils/fileutils"
)

var robotsTxt = []byte("User-agent: *\nDisallow: /\n")

func (w *Web) InitStatic() {
	if *w.ConfigService.Config().ServiceSettings.WebserverMode != "disabled" {
		if err := utils.UpdateAssetsSubpathFromConfig(w.ConfigService.Config()); err != nil {
			mlog.Error("Failed to update assets subpath from config", mlog.Err(err))
		}

		staticDir, _ := fileutils.FindDir(model.CLIENT_DIR)
		mlog.Debug(fmt.Sprintf("Using client directory at %v", staticDir))

		subpath, _ := utils.GetSubpathFromConfig(w.ConfigService.Config())

		mime.AddExtensionType(".wasm", "application/wasm")

		staticHandler := staticFilesHandler(http.StripPrefix(path.Join(subpath, "static"), http.FileServer(http.Dir(staticDir))))
		pluginHandler := staticFilesHandler(http.StripPrefix(path.Join(subpath, "static", "plugins"), http.FileServer(http.Dir(*w.ConfigService.Config().PluginSettings.ClientDirectory))))

		if *w.ConfigService.Config().ServiceSettings.WebserverMode == "gzip" {
			staticHandler = gziphandler.GzipHandler(staticHandler)
			pluginHandler = gziphandler.GzipHandler(pluginHandler)
		}

		w.MainRouter.PathPrefix("/static/plugins/").Handler(pluginHandler)
		w.MainRouter.PathPrefix("/static/").Handler(staticHandler)
		w.MainRouter.Handle("/robots.txt", http.HandlerFunc(robotsHandler))
		w.MainRouter.Handle("/unsupported_browser.js", http.HandlerFunc(unsupportedBrowserScriptHandler))
		w.MainRouter.Handle("/{anything:.*}", w.NewStaticHandler(root)).Methods("GET")

		// When a subpath is defined, it's necessary to handle redirects without a
		// trailing slash. We don't want to use StrictSlash on the w.MainRouter and affect
		// all routes, just /subpath -> /subpath/.
		w.MainRouter.HandleFunc("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path += "/"
			http.Redirect(w, r, r.URL.String(), http.StatusFound)
		}))
	}
}

func root(c *Context, w http.ResponseWriter, r *http.Request) {

	if !CheckClientCompatability(r.UserAgent()) {
		w.Header().Set("Cache-Control", "no-store")
		page := utils.NewHTMLTemplate(c.App.HTMLTemplates(), "unsupported_browser")
		page.Props["Title"] = c.App.T("web.error.unsupported_browser.title")
		page.Props["Message"] = c.App.T("web.error.unsupported_browser.message")
		page.Props["App"] = struct {
			LogoSrc                string
			Title                  string
			SupportedVersionString string
			Label64Bit             string
			Link64Bit              string
			Label32Bit             string
			Link32Bit              string
			InstallGuide           string
			InstallGuideLink       string
		}{
			"/static/images/browser-icons/mac.png",
			"Download the App",
			"Version 4.2.1",
			"Download 64-Bit",
			"http://www.mattermost.com",
			"Download 32-Bit",
			"http://www.mattermost.com",
			"Install Guide",
			"http://www.mattermost.com",
		}
		page.Props["Browsers"] = []struct {
			LogoSrc                string
			Title                  string
			SupportedVersionString string
			Src                    string
			GetLatestString        string
		}{{
			"/static/images/browser-icons/chrome.png",
			"Google Chrome",
			"Version 61+",
			"http://www.google.com/chrome",
			"Get the latest Chrome browser",
		}, {
			"/static/images/browser-icons/firefox.png",
			"Mozilla Firefox",
			"Version 61+",
			"http://www.google.com/chrome",
			"Get the latest Firefox browser",
		}}
		page.Props["SystemBrowser"] = struct {
			LogoSrc                string
			Title                  string
			SupportedVersionString string
			LabelOpen              string
			LinkOpen               string
			LinkMakeDefault        string
		}{
			"/static/images/browser-icons/safari.png",
			"Apple Safari",
			"Version 9+",
			"Open Safari",
			"http://www.google.com/chrome",
			"http://www.google.com/chrome",
		}

		ua := uasurfer.Parse(r.UserAgent())
		page.Props["OSVersion"] = ua.OS.Name.String()
		page.RenderToWriter(w)
		return
	}

	if IsApiCall(c.App, r) {
		Handle404(c.App, w, r)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, max-age=31556926, public")

	staticDir, _ := fileutils.FindDir(model.CLIENT_DIR)
	http.ServeFile(w, r, filepath.Join(staticDir, "root.html"))
}

func staticFilesHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=31556926, public")
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}
	w.Write(robotsTxt)
}

func unsupportedBrowserScriptHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}

	templatesDir, _ := fileutils.FindDir("templates")
	http.ServeFile(w, r, filepath.Join(templatesDir, "unsupported_browser.js"))
}
