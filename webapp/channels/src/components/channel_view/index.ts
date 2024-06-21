// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import type {ConnectedProps} from 'react-redux';
import {withRouter} from 'react-router-dom';

import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';
import {getConfig, getLicense} from 'mattermost-redux/selectors/entities/general';
import {getCurrentRelativeTeamUrl} from 'mattermost-redux/selectors/entities/teams';
import {isFirstAdmin} from 'mattermost-redux/selectors/entities/users';

import {goToLastViewedChannel} from 'actions/views/channel';

import type {GlobalState} from 'types/store';

import ChannelView from './channel_view';

function mapStateToProps(state: GlobalState) {
    const channel = getCurrentChannel(state);

    const config = getConfig(state);

    const viewArchivedChannels = config.ExperimentalViewArchivedChannels === 'true';
    const enableOnboardingFlow = config.EnableOnboardingFlow === 'true';
    const enableWebSocketEventScope = config.FeatureFlagWebSocketEventScope === 'true';

    return {
        channelId: channel ? channel.id : '',
        enableOnboardingFlow,
        channelIsArchived: channel ? channel.delete_at !== 0 : false,
        viewArchivedChannels,
        isCloud: getLicense(state).Cloud === 'true',
        teamUrl: getCurrentRelativeTeamUrl(state),
        isFirstAdmin: isFirstAdmin(state),
        enableWebSocketEventScope,
    };
}

const mapDispatchToProps = ({
    goToLastViewedChannel,
});

const connector = connect(mapStateToProps, mapDispatchToProps);

export type PropsFromRedux = ConnectedProps<typeof connector>;

export default withRouter(connector(ChannelView));
