// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import Client from 'utils/web_client.jsx';
import LoadingScreen from '../loading_screen.jsx';
import UserItem from './user_item.jsx';
import ResetPasswordModal from './reset_password_modal.jsx';

import {FormattedMessage} from 'react-intl';

import React from 'react';

export default class UserList extends React.Component {
    constructor(props) {
        super(props);

        this.getTeamProfiles = this.getTeamProfiles.bind(this);
        this.getCurrentTeamProfiles = this.getCurrentTeamProfiles.bind(this);
        this.doPasswordReset = this.doPasswordReset.bind(this);
        this.doPasswordResetDismiss = this.doPasswordResetDismiss.bind(this);
        this.doPasswordResetSubmit = this.doPasswordResetSubmit.bind(this);

        this.state = {
            teamId: props.team.id,
            users: null,
            serverError: null,
            showPasswordModal: false,
            user: null
        };
    }

    componentDidMount() {
        this.getCurrentTeamProfiles();
    }

    getCurrentTeamProfiles() {
        this.getTeamProfiles(this.props.team.id);
    }

    getTeamProfiles(teamId) {
        Client.getProfilesForTeam(
            teamId,
            (users) => {
                var memberList = [];
                for (var id in users) {
                    if (users.hasOwnProperty(id)) {
                        memberList.push(users[id]);
                    }
                }

                memberList.sort((a, b) => {
                    if (a.username < b.username) {
                        return -1;
                    }

                    if (a.username > b.username) {
                        return 1;
                    }

                    return 0;
                });

                this.setState({
                    users: memberList
                });
            },
            (err) => {
                this.setState({
                    users: null,
                    serverError: err.message
                });
            }
        );
    }

    doPasswordReset(user) {
        this.setState({
            showPasswordModal: true,
            user
        });
    }

    doPasswordResetDismiss() {
        this.setState({
            showPasswordModal: false,
            user: null
        });
    }

    doPasswordResetSubmit() {
        this.getCurrentTeamProfiles();
        this.setState({
            showPasswordModal: false,
            user: null
        });
    }

    componentWillReceiveProps(newProps) {
        this.getTeamProfiles(newProps.team.id);
    }

    render() {
        var serverError = '';
        if (this.state.serverError) {
            serverError = <div className='form-group has-error'><label className='control-label'>{this.state.serverError}</label></div>;
        }

        if (this.state.users == null) {
            return (
                <div className='wrapper--fixed'>
                    <h3>
                        <FormattedMessage
                            id='admin.userList.title'
                            defaultMessage='Users for {team}'
                            values={{
                                team: this.props.team.name
                            }}
                        />
                    </h3>
                    {serverError}
                    <LoadingScreen/>
                </div>
            );
        }

        var memberList = this.state.users.map((user) => {
            return (
                <UserItem
                    key={'user_' + user.id}
                    user={user}
                    refreshProfiles={this.getCurrentTeamProfiles}
                    doPasswordReset={this.doPasswordReset}
                />);
        });

        return (
            <div className='wrapper--fixed'>
                <h3>
                    <FormattedMessage
                        id='admin.userList.title2'
                        defaultMessage='Users for {team} ({count})'
                        values={{
                            team: this.props.team.name,
                            count: this.state.users.length
                        }}
                    />
                </h3>
                {serverError}
                <form
                    className='form-horizontal'
                    role='form'
                >
                    <div className='more-modal__list member-list-holder'>
                        {memberList}
                    </div>
                </form>
                <ResetPasswordModal
                    user={this.state.user}
                    show={this.state.showPasswordModal}
                    team={this.props.team}
                    onModalSubmit={this.doPasswordResetSubmit}
                    onModalDismissed={this.doPasswordResetDismiss}
                />
            </div>
        );
    }
}

UserList.propTypes = {
    team: React.PropTypes.object
};
