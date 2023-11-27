// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState, useRef} from 'react';
import {Modal, type ModalBody} from 'react-bootstrap';
import {useIntl} from 'react-intl';

import TeamSettings from 'components/team_settings';

const SettingsSidebar = React.lazy(() => import('components/settings_sidebar'));

type Props = {
    onExited: () => void;
}

const TeamSettingsModal = (props: Props) => {
    const [activeTab, setActiveTab] = useState('info');
    const [show, setShow] = useState<boolean>(true);
    const modalBodyRef = useRef<ModalBody>(null);
    const {formatMessage} = useIntl();

    const updateTab = (tab: string) => setActiveTab(tab);

    const handleHide = () => setShow(false);

    const handleHidden = () => {
        setActiveTab('info');
        props.onExited();
    };

    const tabs = [
        {name: 'info', uiName: formatMessage({id: 'team_settings_modal.infoTab', defaultMessage: 'Info'}), icon: 'icon icon-information-outline', iconTitle: formatMessage({id: 'generic_icons.info', defaultMessage: 'Info Icon'})},
        {name: 'access', uiName: formatMessage({id: 'team_settings_modal.accessTab', defaultMessage: 'Access'}), icon: 'icon icon-account-multiple-outline', iconTitle: formatMessage({id: 'generic_icons.member', defaultMessage: 'Member Icon'})},
    ];

    return (
        <Modal
            dialogClassName='a11y__modal settings-modal settings-modal--action'
            show={show}
            onHide={handleHide}
            onExited={handleHidden}
            role='dialog'
            aria-labelledby='teamSettingsModalLabel'
            id='teamSettingsModal'
        >
            <Modal.Header
                id='teamSettingsModalLabel'
                closeButton={true}
            >
                <Modal.Title componentClass='h1'>
                    {formatMessage({id: 'team_settings_modal.title', defaultMessage: 'Team Settings'})}
                </Modal.Title>
            </Modal.Header>
            <Modal.Body ref={modalBodyRef}>
                <div className='settings-table'>
                    <div className='settings-links'>
                        <React.Suspense fallback={null}>
                            <SettingsSidebar
                                tabs={tabs}
                                activeTab={activeTab}
                                updateTab={updateTab}
                            />
                        </React.Suspense>
                    </div>
                    <div className='settings-content'>
                        <TeamSettings
                            activeTab={activeTab}
                            closeModal={handleHide}
                        />
                    </div>
                </div>
            </Modal.Body>
        </Modal>
    );
};

export default TeamSettingsModal;
