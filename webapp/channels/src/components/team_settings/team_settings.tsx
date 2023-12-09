// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import type {Team} from '@mattermost/types/teams';

import AccessTab from './team_access_tab';
import InfoTab from './team_info_tab';

type Props = {
    activeTab: string;
    hasChanges: boolean;
    hasChangeTabError: boolean;
    closeModal: () => void;
    setHasChanges: (hasChanges: boolean) => void;
    setHasChangeTabError: (hasChangesError: boolean) => void;
    team?: Team;
};

const TeamSettings = ({
    activeTab = '',
    closeModal,
    team,
    hasChanges,
    hasChangeTabError,
    setHasChanges,
    setHasChangeTabError,
}: Props): JSX.Element | null => {
    if (!team) {
        return null;
    }

    let result;
    switch (activeTab) {
    case 'info':
        result = (
            <InfoTab
                team={team}
                hasChanges={hasChanges}
                setHasChanges={setHasChanges}
                hasChangeTabError={hasChangeTabError}
                setHasChangeTabError={setHasChangeTabError}
            />
        );
        break;
    case 'access':
        result = (
            <AccessTab
                team={team}
                hasChanges={hasChanges}
                setHasChanges={setHasChanges}
                hasChangeTabError={hasChangeTabError}
                setHasChangeTabError={setHasChangeTabError}
            />
        );
        break;
    default:
        result = (
            <div/>
        );
        break;
    }

    return result;
};

export default TeamSettings;
