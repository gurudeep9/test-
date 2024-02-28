// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {FormattedMessage, useIntl} from 'react-intl';

import IconButton from '@mattermost/compass-components/components/icon-button'; // eslint-disable-line no-restricted-imports

import OverlayTrigger from 'components/overlay_trigger';
import Tooltip from 'components/tooltip';
import UserSettingsModalLegacy from 'components/user_settings/modal';
import UserSettingsModal from 'components/user_settings_modal';

import Constants, {ModalIdentifiers} from 'utils/constants';

import type {ModalData} from 'types/actions';

type Props = {
    isUserSettingsModalRevamp: boolean;
    actions: {
        openModal: <P>(modalData: ModalData<P>) => void;
    };
};

const SettingsButton = (props: Props): JSX.Element | null => {
    const {formatMessage} = useIntl();

    const tooltip = (
        <Tooltip id='productSettings'>
            <FormattedMessage
                id='global_header.productSettings'
                defaultMessage='Settings'
            />
        </Tooltip>
    );

    const handleClick = () => {
        if (props.isUserSettingsModalRevamp) {
            props.actions.openModal({modalId: ModalIdentifiers.USER_SETTINGS, dialogType: UserSettingsModal});
        } else {
            props.actions.openModal({modalId: ModalIdentifiers.USER_SETTINGS, dialogType: UserSettingsModalLegacy, dialogProps: {isContentProductSettings: true}});
        }
    }

    return (
        <OverlayTrigger
            trigger={['hover', 'focus']}
            delayShow={Constants.OVERLAY_TIME_DELAY}
            placement='bottom'
            overlay={tooltip}
        >
            <IconButton
                size={'sm'}
                icon={'settings-outline'}
                onClick={handleClick}
                inverted={true}
                compact={true}
                aria-haspopup='dialog'
                aria-label={formatMessage({id: 'global_header.productSettings', defaultMessage: 'Settings'})}
            />
        </OverlayTrigger>
    );
};

export default SettingsButton;
