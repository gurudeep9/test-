// Copyright (c) 2015 Spinpunch, Inc. All Rights Reserved.
// See License.txt for license information.

var UserStore = require('../stores/user_store.jsx');
var Client = require('../utils/client.jsx');
var AsyncClient = require('../utils/async_client.jsx');

function getStateFromStoresForSessions() {
    return {
        sessions: UserStore.getSessions(),
        server_error: null,
        client_error: null
    };
}

module.exports = React.createClass({
    submitRevoke: function(altId) {
        var self = this;
        Client.revokeSession(altId,
            function(data) {
                AsyncClient.getSessions();
            }.bind(this),
            function(err) {
                state = getStateFromStoresForSessions();
                state.server_error = err;
                this.setState(state);
            }.bind(this)
        );
    },
    componentDidMount: function() {
        UserStore.addSessionsChangeListener(this._onChange);
        AsyncClient.getSessions();
    },
    componentWillUnmount: function() {
        UserStore.removeSessionsChangeListener(this._onChange);
    },
    _onChange: function() {
        this.setState(getStateFromStoresForSessions());
    },
    getInitialState: function() {
        return getStateFromStoresForSessions();
    },
    render: function() {
        var activityList = [];
        var server_error = this.state.server_error ? this.state.server_error : null;

        for (var i = 0; i < this.state.sessions.length; i++) {
            var currentSession = this.state.sessions[i];
            var lastAccessTime = new Date(currentSession.last_activity_at);
            var firstAccessTime = new Date(currentSession.create_at);
            var devicePicture = "";

            if (currentSession.props.platform === "Windows") {
                devicePicture = "windows-picture";
            }
            else if (currentSession.props.platform === "Macintosh" || currentSession.props.platform === "iPhone") {
                devicePicture = "apple-picture";
            }
            
            activityList[i] = (
                <div>
                    <div className="single-device">
                        <div>
                            <div className={devicePicture} />
                            <div className="device-platform-name">{currentSession.props.platform}</div>
                        </div>
                        <div className="activity-info">
                            <div>{"Last activity: " + lastAccessTime.toDateString() + ", " + lastAccessTime.toLocaleTimeString()}</div>
                            <div>{"First time active: " + firstAccessTime.toDateString() + ", " + lastAccessTime.toLocaleTimeString()}</div>
                            <div>{"OS: " + currentSession.props.os}</div>
                            <div>{"Browser: " + currentSession.props.browser}</div>
                            <div>{"Session ID: " + currentSession.alt_id}</div>
                        </div>
                        <div><button onClick={this.submitRevoke.bind(this, currentSession.alt_id)} className="pull-right btn btn-primary">Revoke</button></div>
                        <br/>
                        {i < this.state.sessions.length - 1 ?
                        <div className="divider-light"/>
                        :
                        null
                        }
                    </div>
                </div>
            );
        }

        return (
            <div>
                <div className="modal fade" ref="modal" id="activity_log" tabIndex="-1" role="dialog" aria-hidden="true">
                    <div className="modal-dialog">
                        <div className="modal-content">
                            <div className="modal-header">
                                <button type="button" className="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                                <h4 className="modal-title" id="myModalLabel">Active Devices</h4>
                            </div>
                            <div ref="modalBody" className="modal-body">
                                <form role="form">
                                { activityList }
                                </form>
                                { server_error }
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
});
