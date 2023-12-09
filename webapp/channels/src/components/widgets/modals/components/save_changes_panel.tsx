// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import classNames from 'classnames';
import React from 'react';
import {FormattedMessage} from 'react-intl';

import {AlertCircleOutlineIcon} from '@mattermost/compass-icons/components';

import './save_changes_panel.scss';

export type SaveChangesPanelState = 'saving' | 'saved' | 'error';

type Props = {
    handleSubmit: () => void;
    handleCancel: () => void;
    handleClose: () => void;
    tabChangeError?: boolean;
    state: SaveChangesPanelState;
}
function SaveChangesPanel({handleSubmit, handleCancel, handleClose, tabChangeError = false, state = 'saving'}: Props) {
    const panelClassName = classNames('mm-save-changes-panel', {error: tabChangeError || state === 'error'}, {saved: state === 'saved'});
    const messageClassName = classNames('mm-save-changes-panel__message', {error: tabChangeError || state === 'error'}, {saved: state === 'saved'});
    const cancelButtonClassName = classNames('mm-save-changes-panel__cancel-btn', {error: tabChangeError || state === 'error'}, {saved: state === 'saved'});
    const saveButtonClassName = classNames('mm-save-changes-panel__save-btn', {error: tabChangeError || state === 'error'}, {saved: state === 'saved'});

    const generateMessage = () => {
        if (tabChangeError || state === 'saving') {
            return (
                <FormattedMessage
                    id='saveChangesPanel.message'
                    defaultMessage='You have unsaved changes'
                />
            );
        }

        if (state === 'error') {
            return (
                <FormattedMessage
                    id='saveChangesPanel.error'
                    defaultMessage='There was an error saving your settings'
                />
            );
        }

        return (
            <FormattedMessage
                id='saveChangesPanel.saved'
                defaultMessage='Settings saved'
            />
        );
    };

    const generateControlButtons = () => {
        if (state === 'saved') {
            return (
                <div className='mm-save-changes-panel__btn-ctr'>
                    <button
                        id='panelCloseButton'
                        type='button'
                        className='btn btn-icon btn-sm'
                        onClick={handleClose}
                    >
                        <i
                            className='icon icon-close'
                        />
                    </button>
                </div>
            );
        }

        return (
            <div className='mm-save-changes-panel__btn-ctr'>
                <button
                    className={cancelButtonClassName}
                    onClick={handleCancel}
                >
                    <FormattedMessage
                        id='saveChangesPanel.cancel'
                        defaultMessage='Undo'
                    />
                </button>
                <button
                    className={saveButtonClassName}
                    onClick={handleSubmit}
                >
                    {state === 'error' ?
                        <FormattedMessage
                            id='saveChangesPanel.tryAgain'
                            defaultMessage='Try again'
                        /> :
                        <FormattedMessage
                            id='saveChangesPanel.save'
                            defaultMessage='Save'
                        />
                    }
                </button>
            </div>
        );
    };

    return (
        <div className={panelClassName}>
            <p className={messageClassName}>
                <AlertCircleOutlineIcon
                    size={18}
                    color={'currentcolor'}
                />
                {generateMessage()}
            </p>
            {generateControlButtons()}
        </div>
    );
}
export default SaveChangesPanel;
