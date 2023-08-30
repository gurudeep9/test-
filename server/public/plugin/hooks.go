// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package plugin

import (
	"io"
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
)

// These assignments are part of the wire protocol used to trigger hook events in plugins.
//
// Feel free to add more, but do not change existing assignments. Follow the naming convention of
// <HookName>ID as the autogenerated glue code depends on that.
const (
	OnActivateID                              = 0
	OnDeactivateID                            = 1
	ServeHTTPID                               = 2
	OnConfigurationChangeID                   = 3
	ExecuteCommandID                          = 4
	MessageWillBePostedID                     = 5
	MessageWillBeUpdatedID                    = 6
	MessageHasBeenPostedID                    = 7
	MessageHasBeenUpdatedID                   = 8
	UserHasJoinedChannelID                    = 9
	UserHasLeftChannelID                      = 10
	UserHasJoinedTeamID                       = 11
	UserHasLeftTeamID                         = 12
	ChannelHasBeenCreatedID                   = 13
	FileWillBeUploadedID                      = 14
	UserWillLogInID                           = 15
	UserHasLoggedInID                         = 16
	UserHasBeenCreatedID                      = 17
	ReactionHasBeenAddedID                    = 18
	ReactionHasBeenRemovedID                  = 19
	OnPluginClusterEventID                    = 20
	OnWebSocketConnectID                      = 21
	OnWebSocketDisconnectID                   = 22
	WebSocketMessageHasBeenPostedID           = 23
	RunDataRetentionID                        = 24
	OnInstallID                               = 25
	OnSendDailyTelemetryID                    = 26
	OnCloudLimitsUpdatedID                    = 27
	deprecatedUserHasPermissionToCollectionID = 28
	deprecatedGetAllUserIdsForCollectionID    = 29
	deprecatedGetAllCollectionIDsForUserID    = 30
	deprecatedGetTopicRedirectID              = 31
	deprecatedGetCollectionMetadataByIdsID    = 32
	deprecatedGetTopicMetadataByIdsID         = 33
	ConfigurationWillBeSavedID                = 34
	NotificationWillBePushedID                = 35
	UserHasBeenDeactivatedID                  = 36
  MessageHasBeenDeletedID                   = 37
	TotalHooksID                              = iota
)

const (
	// DismissPostError dismisses a pending post when the error is returned from MessageWillBePosted.
	DismissPostError = "plugin.message_will_be_posted.dismiss_post"
)

// Hooks describes the methods a plugin may implement to automatically receive the corresponding
// event.
//
// A plugin only need implement the hooks it cares about. The MattermostPlugin provides some
// default implementations for convenience but may be overridden.
type Hooks interface {
	// OnActivate is invoked when the plugin is activated. If an error is returned, the plugin
	// will be terminated. The plugin will not receive hooks until after OnActivate returns
	// without error. OnConfigurationChange will be called once before OnActivate.
	//
	// Minimum server version: 5.2
	OnActivate() error

	// Implemented returns a list of hooks that are implemented by the plugin.
	// Plugins do not need to provide an implementation. Any given will be ignored.
	//
	// Minimum server version: 5.2
	Implemented() ([]string, error)

	// OnDeactivate is invoked when the plugin is deactivated. This is the plugin's last chance to
	// use the API, and the plugin will be terminated shortly after this invocation. The plugin
	// will stop receiving hooks just prior to this method being called.
	//
	// Minimum server version: 5.2
	OnDeactivate() error

	// OnConfigurationChange is invoked when configuration changes may have been made. Any
	// returned error is logged, but does not stop the plugin. You must be prepared to handle
	// a configuration failure gracefully. It is called once before OnActivate.
	//
	// Minimum server version: 5.2
	OnConfigurationChange() error

	// ServeHTTP allows the plugin to implement the http.Handler interface. Requests destined for
	// the /plugins/{id} path will be routed to the plugin.
	//
	// The Mattermost-User-Id header will be present if (and only if) the request is by an
	// authenticated user.
	//
	// Minimum server version: 5.2
	ServeHTTP(c *Context, w http.ResponseWriter, r *http.Request)

	// ExecuteCommand executes a command that has been previously registered via the RegisterCommand
	// API.
	//
	// Minimum server version: 5.2
	ExecuteCommand(c *Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError)

	// UserHasBeenCreated is invoked after a user was created.
	//
	// Minimum server version: 5.10
	UserHasBeenCreated(c *Context, user *model.User)

	// UserWillLogIn before the login of the user is returned. Returning a non empty string will reject the login event.
	// If you don't need to reject the login event, see UserHasLoggedIn
	//
	// Minimum server version: 5.2
	UserWillLogIn(c *Context, user *model.User) string

	// UserHasLoggedIn is invoked after a user has logged in.
	//
	// Minimum server version: 5.2
	UserHasLoggedIn(c *Context, user *model.User)

	// MessageWillBePosted is invoked when a message is posted by a user before it is committed
	// to the database. If you also want to act on edited posts, see MessageWillBeUpdated.
	//
	// To reject a post, return an non-empty string describing why the post was rejected.
	// To modify the post, return the replacement, non-nil *model.Post and an empty string.
	// To allow the post without modification, return a nil *model.Post and an empty string.
	// To dismiss the post, return a nil *model.Post and the const DismissPostError string.
	//
	// If you don't need to modify or reject posts, use MessageHasBeenPosted instead.
	//
	// Note that this method will be called for posts created by plugins, including the plugin that
	// created the post.
	//
	// Minimum server version: 5.2
	MessageWillBePosted(c *Context, post *model.Post) (*model.Post, string)

	// MessageWillBeUpdated is invoked when a message is updated by a user before it is committed
	// to the database. If you also want to act on new posts, see MessageWillBePosted.
	// Return values should be the modified post or nil if rejected and an explanation for the user.
	// On rejection, the post will be kept in its previous state.
	//
	// If you don't need to modify or rejected updated posts, use MessageHasBeenUpdated instead.
	//
	// Note that this method will be called for posts updated by plugins, including the plugin that
	// updated the post.
	//
	// Minimum server version: 5.2
	MessageWillBeUpdated(c *Context, newPost, oldPost *model.Post) (*model.Post, string)

	// MessageHasBeenPosted is invoked after the message has been committed to the database.
	// If you need to modify or reject the post, see MessageWillBePosted
	// Note that this method will be called for posts created by plugins, including the plugin that
	// created the post.
	//
	// Minimum server version: 5.2
	MessageHasBeenPosted(c *Context, post *model.Post)

	// MessageHasBeenUpdated is invoked after a message is updated and has been updated in the database.
	// If you need to modify or reject the post, see MessageWillBeUpdated
	// Note that this method will be called for posts created by plugins, including the plugin that
	// created the post.
	//
	// Minimum server version: 5.2
	MessageHasBeenUpdated(c *Context, newPost, oldPost *model.Post)

	// MessageHasBeenDeleted is invoked after the message has been deleted from the database.
	// Note that this method will be called for posts deleted by plugins, including the plugin that
	// deleted the post.
	//
	// Minimum server version: 9.1
	MessageHasBeenDeleted(c *Context, post *model.Post)

	// ChannelHasBeenCreated is invoked after the channel has been committed to the database.
	//
	// Minimum server version: 5.2
	ChannelHasBeenCreated(c *Context, channel *model.Channel)

	// UserHasJoinedChannel is invoked after the membership has been committed to the database.
	// If actor is not nil, the user was invited to the channel by the actor.
	//
	// Minimum server version: 5.2
	UserHasJoinedChannel(c *Context, channelMember *model.ChannelMember, actor *model.User)

	// UserHasLeftChannel is invoked after the membership has been removed from the database.
	// If actor is not nil, the user was removed from the channel by the actor.
	//
	// Minimum server version: 5.2
	UserHasLeftChannel(c *Context, channelMember *model.ChannelMember, actor *model.User)

	// UserHasJoinedTeam is invoked after the membership has been committed to the database.
	// If actor is not nil, the user was added to the team by the actor.
	//
	// Minimum server version: 5.2
	UserHasJoinedTeam(c *Context, teamMember *model.TeamMember, actor *model.User)

	// UserHasLeftTeam is invoked after the membership has been removed from the database.
	// If actor is not nil, the user was removed from the team by the actor.
	//
	// Minimum server version: 5.2
	UserHasLeftTeam(c *Context, teamMember *model.TeamMember, actor *model.User)

	// FileWillBeUploaded is invoked when a file is uploaded, but before it is committed to backing store.
	// Read from file to retrieve the body of the uploaded file.
	//
	// To reject a file upload, return an non-empty string describing why the file was rejected.
	// To modify the file, write to the output and/or return a non-nil *model.FileInfo, as well as an empty string.
	// To allow the file without modification, do not write to the output and return a nil *model.FileInfo and an empty string.
	//
	// Note that this method will be called for files uploaded by plugins, including the plugin that uploaded the post.
	// FileInfo.Size will be automatically set properly if you modify the file.
	//
	// Minimum server version: 5.2
	FileWillBeUploaded(c *Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string)

	// ReactionHasBeenAdded is invoked after the reaction has been committed to the database.
	//
	// Note that this method will be called for reactions added by plugins, including the plugin that
	// added the reaction.
	//
	// Minimum server version: 5.30
	ReactionHasBeenAdded(c *Context, reaction *model.Reaction)

	// ReactionHasBeenRemoved is invoked after the removal of the reaction has been committed to the database.
	//
	// Note that this method will be called for reactions removed by plugins, including the plugin that
	// removed the reaction.
	//
	// Minimum server version: 5.30
	ReactionHasBeenRemoved(c *Context, reaction *model.Reaction)

	// OnPluginClusterEvent is invoked when an intra-cluster plugin event is received.
	//
	// This is used to allow communication between multiple instances of the same plugin
	// that are running on separate nodes of the same High-Availability cluster.
	// This hook receives events sent by a call to PublishPluginClusterEvent.
	//
	// Minimum server version: 5.36
	OnPluginClusterEvent(c *Context, ev model.PluginClusterEvent)

	// OnWebSocketConnect is invoked when a new websocket connection is opened.
	//
	// This is used to track which users have connections opened with the Mattermost
	// websocket.
	//
	// Minimum server version: 6.0
	OnWebSocketConnect(webConnID, userID string)

	// OnWebSocketDisconnect is invoked when a websocket connection is closed.
	//
	// This is used to track which users have connections opened with the Mattermost
	// websocket.
	//
	// Minimum server version: 6.0
	OnWebSocketDisconnect(webConnID, userID string)

	// WebSocketMessageHasBeenPosted is invoked when a websocket message is received.
	//
	// Minimum server version: 6.0
	WebSocketMessageHasBeenPosted(webConnID, userID string, req *model.WebSocketRequest)

	// RunDataRetention is invoked during a DataRetentionJob.
	//
	// Minimum server version: 6.4
	RunDataRetention(nowTime, batchSize int64) (int64, error)

	// OnInstall is invoked after the installation of a plugin as part of the onboarding.
	// It's called on every installation, not only once.
	//
	// In the future, other plugin installation methods will trigger this hook, e.g. an installation via the Marketplace.
	//
	// Minimum server version: 6.5
	OnInstall(c *Context, event model.OnInstallEvent) error

	// OnSendDailyTelemetry is invoked when the server send the daily telemetry data.
	//
	// Minimum server version: 6.5
	OnSendDailyTelemetry()

	// OnCloudLimitsUpdated is invoked product limits change, for example when plan tiers change
	//
	// Minimum server version: 7.0
	OnCloudLimitsUpdated(limits *model.ProductLimits)

	// ConfigurationWillBeSaved is invoked before saving the configuration to the
	// backing store.
	// An error can be returned to reject the operation. Additionally, a new
	// config object can be returned to be stored in place of the provided one.
	// Minimum server version: 8.0
	ConfigurationWillBeSaved(newCfg *model.Config) (*model.Config, error)

	// NotificationWillBePushed is invoked before a push notification is sent to the push
	// notification server.
	//
	// To reject a notification, return an non-empty string describing why the notification was rejected.
	// To modify the notification, return the replacement, non-nil *model.PushNotification and an empty string.
	// To allow the notification without modification, return a nil *model.PushNotification and an empty string.
	//
	// Note that this method will be called for push notifications created by plugins, including the plugin that
	// created the notification.
	//
	// Minimum server version: 9.0
	NotificationWillBePushed(pushNotification *model.PushNotification, userID string) (*model.PushNotification, string)

	// UserHasBeenDeactivated is invoked when a user is deactivated.
	//
	// Minimum server version: 9.1
	UserHasBeenDeactivated(c *Context, user *model.User)
}
