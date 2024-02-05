// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import moment from 'moment-timezone';

import {getCurrentTimezoneFull} from 'mattermost-redux/selectors/entities/timezone';
import {getCurrentUser} from 'mattermost-redux/selectors/entities/users';
import type {NewActionFuncAsync} from 'mattermost-redux/types/actions';

import {updateMe} from './users';

export function autoUpdateTimezone(deviceTimezone: string): NewActionFuncAsync {
    return async (dispatch, getState) => {
        const currentUser = getCurrentUser(getState());
        const currentTimezone = getCurrentTimezoneFull(getState());
        const newTimezoneExists = currentTimezone.automaticTimezone !== deviceTimezone;

        if (currentTimezone.useAutomaticTimezone && newTimezoneExists) {
            const timezone = {
                useAutomaticTimezone: 'true',
                automaticTimezone: deviceTimezone,
                manualTimezone: currentTimezone.manualTimezone,
            };

            const updatedUser = {
                ...currentUser,
                timezone,
            };

            dispatch(updateMe(updatedUser));
        }

        moment.tz.setDefault(deviceTimezone);

        return {data: true};
    };
}
