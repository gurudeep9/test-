// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';
import {useIntl} from 'react-intl';

import {RefreshIcon} from '@mattermost/compass-icons/components';
import type {Team} from '@mattermost/types/teams';

import SelectTextInput, {type SelectTextInputOption} from 'components/common/select_text_input/select_text_input';
import Input from 'components/widgets/inputs/input/input';
import BaseSettingItem from 'components/widgets/modals/components/base_setting_item';
import CheckboxSettingItem from 'components/widgets/modals/components/checkbox_setting_item';
import ModalSection from 'components/widgets/modals/components/modal_section';
import SaveChangesPanel, {type SaveChangesPanelState} from 'components/widgets/modals/components/save_changes_panel';

import OpenInvite from './open_invite';

import type {PropsFromRedux, OwnProps} from '.';

import './team_access_tab.scss';

type Props = PropsFromRedux & OwnProps;

const AccessTab = (props: Props) => {
    const generateAllowedDomainOptions = (allowedDomains?: string) => {
        if (!allowedDomains || allowedDomains.length === 0) {
            return [];
        }
        const domainList = allowedDomains.includes(',') ? allowedDomains.split(',') : [allowedDomains];
        return domainList.map((domain) => domain.trim());
    };

    const [inviteId, setInviteId] = useState<Team['invite_id']>(props.team?.invite_id ?? '');
    const [allowedDomains, setAllowedDomains] = useState<string[]>(generateAllowedDomainOptions(props.team?.allowed_domains));
    const [showAllowedDomains, setShowAllowedDomains] = useState<boolean>(allowedDomains?.length > 0);
    const [allowOpenInvite, setAllowOpenInvite] = useState<boolean>(props.team?.allow_open_invite ?? false);
    const [saveChangesPanelState, setSaveChangesPanelState] = useState<SaveChangesPanelState>('saving');
    const {formatMessage} = useIntl();

    const handleAllowedDomainsSubmit = async (): Promise<Error | null> => {
        if (allowedDomains.length === 0) {
            return null;
        }
        const {error} = await props.actions.patchTeam({
            id: props.team?.id,
            allowed_domains: allowedDomains.length === 1 ? allowedDomains[0] : allowedDomains.join(', '),
        });
        if (error) {
            return error;
        }
        return null;
    };

    const handleOpenInviteSubmit = async (): Promise<Error | null> => {
        const data = {
            id: props.team?.id,
            allow_open_invite: allowOpenInvite,
        };

        const {error} = await props.actions.patchTeam(data);
        if (error) {
            return error;
        }
        return null;
    };

    const updateAllowedDomains = (domain: string) => {
        props.setHasChanges(true);
        setAllowedDomains((prev) => [...prev, domain]);
    };
    const updateOpenInvite = (value: boolean) => {
        props.setHasChanges(true);
        setAllowOpenInvite(value);
    };
    const handleOnChangeDomains = (allowedDomainsOptions?: SelectTextInputOption[] | null) => {
        props.setHasChanges(true);
        setAllowedDomains(allowedDomainsOptions?.map((domain) => domain.value) || []);
    };
    const handleRegenerateInviteId = async () => {
        const {data, error} = await props.actions.regenerateTeamInviteId(props.team?.id || '');

        if (data?.invite_id) {
            setInviteId(data.invite_id);
            return;
        }

        if (error) {
            // todo sinan: handle with client error
            // setServerError(error.message);
        }
    };

    const handleEnableAllowedDomains = (enabled: boolean) => {
        setShowAllowedDomains(enabled);
        if (!enabled) {
            setAllowedDomains([]);
        }
    };

    const handleCancel = () => {
        setAllowedDomains(generateAllowedDomainOptions(props.team?.allowed_domains));
        setAllowOpenInvite(props.team?.allow_open_invite ?? false);
        handleClose();
    };

    const handleClose = () => {
        setSaveChangesPanelState('saving');
        props.setHasChanges(false);
        props.setHasChangeTabError(false);
    };

    const handleSaveChanges = async () => {
        const allowedDomainError = await handleAllowedDomainsSubmit();
        const openInviteError = await handleOpenInviteSubmit();
        if (allowedDomainError || openInviteError) {
            setSaveChangesPanelState('error');
            return;
        }
        setSaveChangesPanelState('saved');
        props.setHasChangeTabError(false);
    };

    let inviteSection;
    if (props.canInviteTeamMembers) {
        const inviteSectionInput = (
            <div id='teamInviteContainer' >
                <Input
                    id='teamInviteId'
                    className='form-control'
                    type='text'
                    value={inviteId}
                    maxLength={32}
                />
                <button
                    id='regenerateButton'
                    className='btn btn-tertiary'
                    onClick={handleRegenerateInviteId}
                >
                    <RefreshIcon/>
                    {formatMessage({id: 'general_tab.regenerate', defaultMessage: 'Regenerate'})}
                </button>
            </div>
        );

        inviteSection = (
            <BaseSettingItem
                className='access-invite-section'
                title={{id: 'general_tab.codeTitle', defaultMessage: 'Invite Code'}}
                description={{id: 'general_tab.codeLongDesc', defaultMessage: 'The Invite Code is part of the unique team invitation link which is sent to members you’re inviting to this team. Regenerating the code creates a new invitation link and invalidates the previous link.'}}
                content={inviteSectionInput}
                descriptionAboveContent={true}
            />
        );
    }

    const allowedDomainsSectionInput = (
        <div
            id='allowedDomainsSetting'
            className='form-group'
        >
            <CheckboxSettingItem
                inputFieldData={{title: {id: 'general_tab.allowedDomains', defaultMessage: 'Allow only users with a specific email domain to join this team'}, name: 'name'}}
                inputFieldValue={showAllowedDomains}
                handleChange={handleEnableAllowedDomains}
            />
            {showAllowedDomains &&
                <SelectTextInput
                    id='allowedDomains'
                    placeholder={formatMessage({id: 'general_tab.AllowedDomainsExample', defaultMessage: 'corp.mattermost.com, mattermost.com'})}
                    aria-label={formatMessage({id: 'general_tab.allowedDomains.ariaLabel', defaultMessage: 'Allowed Domains'})}
                    value={allowedDomains}
                    onChange={handleOnChangeDomains}
                    handleNewSelection={updateAllowedDomains}
                    isClearable={false}
                    description={formatMessage({id: 'general_tab.AllowedDomainsTip', defaultMessage: 'Seperate multiple domains with a space or comma.'})}
                />
            }
        </div>
    );

    // todo sinan: convert it to <CheckboxSettingItem like in open invite
    const allowedDomainsSection = (
        <BaseSettingItem
            className='access-allowed-domains-section'
            title={{id: 'general_tab.allowedDomainsTitle', defaultMessage: 'Users with a specific email domain'}}
            description={{id: 'general_tab.allowedDomainsInfo', defaultMessage: 'When enabled, users can only join the team if their email matches a specific domain (e.g. "mattermost.org")'}}
            content={allowedDomainsSectionInput}
            descriptionAboveContent={true}
        />
    );

    // todo sinan: check title font size is same as figma
    return (
        <ModalSection
            content={
                <div className='user-settings'>
                    {props.team?.group_constrained ? undefined : allowedDomainsSection}
                    <div className='divider-light'/>
                    <OpenInvite
                        isGroupConstrained={props.team?.group_constrained}
                        allowOpenInvite={allowOpenInvite}
                        setAllowOpenInvite={updateOpenInvite}
                    />
                    <div className='divider-light'/>
                    {props.team?.group_constrained ? undefined : inviteSection}
                    {props.hasChanges ?
                        <SaveChangesPanel
                            handleCancel={handleCancel}
                            handleSubmit={handleSaveChanges}
                            handleClose={handleClose}
                            tabChangeError={props.hasChangeTabError}
                            state={saveChangesPanelState}
                        /> : undefined}
                </div>
            }
        />
    );
};
export default AccessTab;
