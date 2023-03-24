// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {FormattedMessage} from 'react-intl';
import {useDispatch, useSelector} from 'react-redux';

import IconMessage from 'components/purchase_modal/icon_message';
import FullScreenModal from 'components/widgets/modals/full_screen_modal';

import {closeModal} from 'actions/views/modals';
import {isModalOpen} from 'selectors/views/modals';
import {GlobalState} from 'types/store';
import {Modal} from 'react-bootstrap';

import './result_modal.scss';

import useOpenSalesLink from 'components/common/hooks/useOpenSalesLink';
import {InquiryType} from 'selectors/cloud';

type Props = {
    type?: string;
    onHide?: () => void;
    icon: JSX.Element;
    title: JSX.Element;
    subtitle: JSX.Element;
    primaryButtonText: JSX.Element;
    primaryButtonHandler: () => void;
    identifier: string;
    contactSupportButtonVisible?: boolean;
    resultType: string;
    ignoreExit: boolean;
};

export default function ResultModal({type, icon, title, subtitle, primaryButtonText, primaryButtonHandler, identifier, contactSupportButtonVisible, resultType, ignoreExit, onHide}: Props) {
    const dispatch = useDispatch();

    const openContactUs = useOpenSalesLink(undefined, InquiryType.Technical);

    const isResultModalOpen = useSelector((state: GlobalState) =>
        isModalOpen(state, identifier),
    );

    const handleHide = () => {
        dispatch(closeModal(identifier));
        if (typeof onHide === 'function') {
            onHide();
        }
    };


    const modalType = `delete-workspace-result_modal__${resultType}`;
    if (type === 'small') {
        return (
            <Modal
                className='ResultModal__small'
                show={isResultModalOpen}
                onHide={handleHide}
            >
                <Modal.Header closeButton={true}/>
                <div className={modalType}>
                    <IconMessage
                        formattedTitle={title}
                        formattedSubtitle={subtitle}
                        error={false}
                        icon={icon}
                        formattedButtonText={primaryButtonText}
                        buttonHandler={primaryButtonHandler}
                        className={'success'}
                        formattedTertiaryButonText={
                            contactSupportButtonVisible ?
                                <FormattedMessage
                                    id={'admin.billing.deleteWorkspace.resultModal.ContactSupport'}
                                    defaultMessage={'Contact Support'}
                                /> :
                                undefined
                        }
                        tertiaryButtonHandler={contactSupportButtonVisible ? openContactUs : undefined}
                    />
                </div>
            </Modal>
        );
    }

    return (
        <FullScreenModal
            show={isResultModalOpen}
            onClose={onHide}
            ignoreExit={ignoreExit}
        >
            <div className={modalType}>
                <IconMessage
                    formattedTitle={title}
                    formattedSubtitle={subtitle}
                    error={false}
                    icon={icon}
                    formattedButtonText={primaryButtonText}
                    buttonHandler={primaryButtonHandler}
                    className={'success'}
                    formattedTertiaryButonText={
                        contactSupportButtonVisible ?
                            <FormattedMessage
                                id={'admin.billing.deleteWorkspace.resultModal.ContactSupport'}
                                defaultMessage={'Contact Support'}
                            /> :
                            undefined
                    }
                    tertiaryButtonHandler={contactSupportButtonVisible ? openContactUs : undefined}
                />
            </div>
        </FullScreenModal>
    );
}
