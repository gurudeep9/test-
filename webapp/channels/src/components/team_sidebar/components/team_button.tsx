// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import classNames from 'classnames';
import React, {useCallback, useMemo} from 'react';
import {Draggable} from 'react-beautiful-dnd';
import {defineMessages, useIntl} from 'react-intl';
import {Link} from 'react-router-dom';

import {mark, trackEvent} from 'actions/telemetry_actions.jsx';

import CopyUrlContextMenu from 'components/copy_url_context_menu';
import TeamIcon from 'components/widgets/team_icon/team_icon';
import WithTooltip from 'components/with_tooltip';
import {ShortcutKeys} from 'components/with_tooltip/shortcut';

import {isDesktopApp} from 'utils/user_agent';

const messages = defineMessages({
    nameUndefined: {
        id: 'team.button.name_undefined',
        defaultMessage: 'This team does not have a name',
    },
});

interface Props {
    btnClass?: string;
    url: string;
    displayName?: string;
    content?: React.ReactNode;
    tip: string | JSX.Element;
    order?: number;
    showOrder?: boolean;
    active?: boolean;
    unread?: boolean;
    mentions?: number;
    placement?: 'left' | 'right' | 'top' | 'bottom';
    teamIconUrl?: string | null;
    switchTeam: (url: string) => void;
    isDraggable?: boolean;
    teamIndex?: number;
    teamId?: string;
    isInProduct?: boolean;
    hasUrgent?: boolean;
}

export default function TeamButton({
    btnClass,
    url,
    displayName,
    order,
    unread,
    mentions,
    teamIconUrl,
    isDraggable = false,
    switchTeam,
    teamIndex,
    teamId,
    tip,
    ...otherProps
}: Props) {
    const {formatMessage} = useIntl();

    const handleSwitch = useCallback((e: React.MouseEvent) => {
        mark('TeamLink#click');
        e.preventDefault();
        switchTeam(url);

        setTimeout(() => {
            trackEvent('ui', 'ui_team_sidebar_switch_team');
        }, 0);
    }, [switchTeam, url]);

    let teamClass: string = otherProps.active ? 'active' : '';
    const isNotCreateTeamButton: boolean = !url.endsWith('create_team') && !url.endsWith('select_team');

    let badge: JSX.Element | undefined;

    let ariaLabel = formatMessage({
        id: 'team.button.ariaLabel',
        defaultMessage: '{teamName} team',
    },
    {
        teamName: displayName,
    });

    if (!teamClass) {
        if (unread && !otherProps.isInProduct) {
            teamClass = 'unread';

            badge = (
                <span className={'unread-badge'}/>
            );
        } else if (isNotCreateTeamButton) {
            teamClass = '';
        } else {
            teamClass = 'special';
        }
        ariaLabel = formatMessage({
            id: 'team.button.unread.ariaLabel',
            defaultMessage: '{teamName} team unread',
        },
        {
            teamName: displayName,
        });

        if (mentions) {
            ariaLabel = formatMessage({
                id: 'team.button.mentions.ariaLabel',
                defaultMessage: '{teamName} team, {mentionCount} mentions',
            },
            {
                teamName: displayName,
                mentionCount: mentions,
            });

            badge = (
                <span className={classNames('badge badge-max-number pull-right small', {urgent: otherProps.hasUrgent})}>{mentions > 99 ? '99+' : mentions}</span>
            );
        }
    }

    ariaLabel = ariaLabel.toLowerCase();

    const content = (
        <TeamIcon
            className={teamClass}
            withHover={true}
            content={otherProps.content || displayName || ''}
            url={teamIconUrl}
        />
    );

    let orderIndicator: JSX.Element | undefined;
    if (typeof order !== 'undefined' && order < 10) {
        if (otherProps.showOrder) {
            orderIndicator = (
                <div className='order-indicator'>
                    {order}
                </div>
            );
        }
    }

    const btn = (
        <WithTeamTooltip
            order={order}
            tip={tip}
            url={url}
        >
            <div className={'team-btn ' + btnClass}>
                {!otherProps.isInProduct && badge}
                {content}
            </div>
        </WithTeamTooltip>
    );

    let teamButton = (
        <Link
            id={`${url.slice(1)}TeamButton`}
            aria-label={ariaLabel}
            to={url}
            onClick={handleSwitch}
        >
            {btn}
        </Link>
    );

    if (isDesktopApp()) {
        // if this is not a "special" team button, give it a context menu
        if (isNotCreateTeamButton) {
            teamButton = (
                <CopyUrlContextMenu
                    link={url}
                    menuId={url}
                >
                    {teamButton}
                </CopyUrlContextMenu>
            );
        }
    }

    return isDraggable ? (
        <Draggable
            draggableId={teamId!}
            index={teamIndex!}
        >
            {(provided, snapshot) => {
                return (
                    <div
                        className='draggable-team-container'
                        ref={provided.innerRef}
                        {...provided.draggableProps}
                        {...provided.dragHandleProps}
                        tabIndex={-1}
                    >
                        <div

                            className={classNames([`team-container ${teamClass}`, {isDragging: snapshot.isDragging}])}
                            tabIndex={-1}
                        >
                            {teamButton}
                            {orderIndicator}
                        </div>
                    </div>
                );
            }}
        </Draggable>
    ) : (
        <div
            className={`team-container ${teamClass}`}
            tabIndex={-1}
        >
            {teamButton}
            {orderIndicator}
        </div>
    );
}

function WithTeamTooltip({
    order,
    tip,
    url,
    children,
}: React.PropsWithChildren<Pick<Props, 'order' | 'tip' | 'url'>>) {
    const intl = useIntl();

    const shortcut = useMemo(() => {
        if (!order || order >= 10) {
            return undefined;
        }

        return {
            default: [ShortcutKeys.ctrl, ShortcutKeys.alt, order.toString()],
            mac: [ShortcutKeys.cmd, ShortcutKeys.option, order.toString()],
        };
    }, [order]);

    return (
        <WithTooltip
            id={`tooltip-${url}`}
            title={tip || intl.formatMessage(messages.nameUndefined)}
            shortcut={shortcut}
            placement='right'
        >
            {children}
        </WithTooltip>
    );
}
