// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import React from 'react';

import Constants from 'utils/constants.jsx';
import PureRenderMixin from 'react-addons-pure-render-mixin';

import {getDateForUnixTicks, isMobile} from 'utils/utils.jsx';

import {Link} from 'react-router/es6';
import TeamStore from 'stores/team_store.jsx';

export default class PostTime extends React.Component {
    constructor(props) {
        super(props);

        this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this);
        this.state = {
            currentTeamDisplayName: TeamStore.getCurrent().display_name
        };
    }

    componentDidMount() {
        this.intervalId = setInterval(() => {
            this.forceUpdate();
        }, Constants.TIME_SINCE_UPDATE_INTERVAL);
    }

    componentWillUnmount() {
        clearInterval(this.intervalId);
    }

    renderTimeTag() {
        return (
            <time className='post__time'>
                {getDateForUnixTicks(this.props.eventTime).toLocaleString('en', {hour: '2-digit', minute: '2-digit', hour12: !this.props.useMilitaryTime})}
            </time>
        );
    }

    render() {
        return isMobile() ?
            this.renderTimeTag() :
            (
                <Link
                    to={`/${this.state.currentTeamDisplayName}/pl/${this.props.postId}`}
                    target='_blank'
                >
                    {this.renderTimeTag()}
                </Link>
            );
    }
}

PostTime.defaultProps = {
    eventTime: 0,
    sameUser: false
};

PostTime.propTypes = {
    eventTime: React.PropTypes.number.isRequired,
    sameUser: React.PropTypes.bool,
    compactDisplay: React.PropTypes.bool,
    useMilitaryTime: React.PropTypes.bool.isRequired,
    postId: React.PropTypes.string
};
