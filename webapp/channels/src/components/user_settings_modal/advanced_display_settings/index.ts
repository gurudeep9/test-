// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import type {Dispatch} from 'redux';
import {bindActionCreators} from 'redux';

import {savePreferences} from 'mattermost-redux/actions/preferences';
import {updateUserActive, revokeAllSessionsForUser} from 'mattermost-redux/actions/users';
import {Preferences as ReduxPreferences} from 'mattermost-redux/constants';
import {getConfig, isPerformanceDebuggingEnabled} from 'mattermost-redux/selectors/entities/general';
import {
    get,
    getUnreadScrollPositionPreference,
    makeGetCategory,
} from 'mattermost-redux/selectors/entities/preferences';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/users';

import {Preferences} from 'utils/constants';

import type {GlobalState} from 'types/store';

import AdvancedDisplaySettings from './advanced_display_settings';

function makeMapStateToProps() {
    const getAdvancedSettingsCategory = makeGetCategory();

    return (state: GlobalState) => {
        const config = getConfig(state);

        const enablePreviewFeatures = config.EnablePreviewFeatures === 'true';
        const enableUserDeactivation = config.EnableUserDeactivation === 'true';

        return {
            advancedSettingsCategory: getAdvancedSettingsCategory(state, Preferences.CATEGORY_ADVANCED_SETTINGS),
            sendOnCtrlEnter: get(state, Preferences.CATEGORY_ADVANCED_SETTINGS, 'send_on_ctrl_enter', 'false'),
            codeBlockOnCtrlEnter: get(state, Preferences.CATEGORY_ADVANCED_SETTINGS, 'code_block_ctrl_enter', 'true'),
            formatting: get(state, Preferences.CATEGORY_ADVANCED_SETTINGS, 'formatting', 'true'),
            joinLeave: get(state, Preferences.CATEGORY_ADVANCED_SETTINGS, 'join_leave', 'true'),
            disableClientPlugins: get(state, ReduxPreferences.CATEGORY_PERFORMANCE_DEBUGGING, ReduxPreferences.NAME_DISABLE_CLIENT_PLUGINS, 'false'),
            disableTelemetry: get(state, ReduxPreferences.CATEGORY_PERFORMANCE_DEBUGGING, ReduxPreferences.NAME_DISABLE_TELEMETRY, 'false'),
            disableTypingMessages: get(state, ReduxPreferences.CATEGORY_PERFORMANCE_DEBUGGING, ReduxPreferences.NAME_DISABLE_TYPING_MESSAGES, 'false'),
            performanceDebuggingEnabled: isPerformanceDebuggingEnabled(state),
            currentUser: getCurrentUser(state),
            unreadScrollPosition: getUnreadScrollPositionPreference(state),
            enablePreviewFeatures,
            enableUserDeactivation,
        };
    };
}

function mapDispatchToProps(dispatch: Dispatch) {
    return {
        actions: bindActionCreators({
            savePreferences,
            updateUserActive,
            revokeAllSessionsForUser,
        }, dispatch),
    };
}

export default connect(makeMapStateToProps, mapDispatchToProps)(AdvancedDisplaySettings);
