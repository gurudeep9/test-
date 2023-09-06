// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {Overlay} from 'react-bootstrap';

import type {Group} from '@mattermost/types/groups';
import type {UserProfile} from '@mattermost/types/users';

import {Client4} from 'mattermost-redux/client';
import {displayUsername} from 'mattermost-redux/utils/user_utils';

import AtMentionGroup from 'components/at_mention/at_mention_group';
import ProfilePopover from 'components/profile_popover';

import Constants from 'utils/constants';
import {isKeyPressed} from 'utils/keyboard';
import {popOverOverlayPosition} from 'utils/position_utils';
import {getUserOrGroupFromMentionName} from 'utils/post_utils';
import {getViewportSize} from 'utils/utils';

type Props = {
    currentUserId: string;
    mentionName: string;
    teammateNameDisplay: string;
    usersByUsername: Record<string, UserProfile>;
    groupsByName: Record<string, Group>;
    children?: React.ReactNode;
    channelId?: string;
    hasMention?: boolean;
    disableHighlight?: boolean;
    disableGroupHighlight?: boolean;
}

type State = {
    show: boolean;
    target?: HTMLAnchorElement;
    placement?: string;
}

export default class AtMention extends React.PureComponent<Props, State> {
    buttonRef: React.RefObject<HTMLAnchorElement>;

    static defaultProps: Partial<Props> = {
        hasMention: false,
        disableHighlight: false,
        disableGroupHighlight: false,
    };

    constructor(props: Props) {
        super(props);

        this.state = {
            show: false,
        };

        this.buttonRef = React.createRef();
    }

    showOverlay = (target?: HTMLAnchorElement) => {
        const targetBounds = this.buttonRef.current?.getBoundingClientRect();

        if (targetBounds) {
            const placement = popOverOverlayPosition(targetBounds, getViewportSize().h, getViewportSize().h - 240);
            this.setState({target, show: !this.state.show, placement});
        }
    };

    handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        this.showOverlay(e.target as HTMLAnchorElement);
    };

    handleKeyDown = (e: React.KeyboardEvent<HTMLAnchorElement>) => {
        if (isKeyPressed(e, Constants.KeyCodes.ENTER) || isKeyPressed(e, Constants.KeyCodes.SPACE)) {
            e.preventDefault();

            // Prevent propagation so that the message textbox isn't focused
            e.stopPropagation();
            this.showOverlay(e.target as HTMLAnchorElement);
        }
    };

    hideOverlay = () => {
        this.setState({show: false});
    };

    render() {
        const user = getUserOrGroupFromMentionName(this.props.usersByUsername, this.props.mentionName) as UserProfile | '';

        if (!this.props.disableGroupHighlight && !user) {
            const group = getUserOrGroupFromMentionName(this.props.groupsByName, this.props.mentionName) as Group | '';

            if (group && group.allow_reference) {
                const suffix = this.props.mentionName.substring(group.name.length);

                return (
                    <>
                        <span>
                            <AtMentionGroup group={group}/>
                        </span>
                        {suffix}
                    </>
                );
            }
        }

        if (!user) {
            return <React.Fragment>{this.props.children}</React.Fragment>;
        }

        const suffix = this.props.mentionName.substring(user.username.length);
        const displayName = displayUsername(user, this.props.teammateNameDisplay);

        const highlightMention = !this.props.disableHighlight && user.id === this.props.currentUserId;

        return (
            <>
                <span
                    className={highlightMention ? 'mention--highlight' : undefined}
                >
                    <Overlay
                        placement={this.state.placement}
                        show={this.state.show}
                        target={this.state.target}
                        rootClose={true}
                        onHide={this.hideOverlay}
                    >
                        <ProfilePopover
                            className='user-profile-popover'
                            userId={user.id}
                            src={Client4.getProfilePictureUrl(user.id, user.last_picture_update)}
                            hasMention={this.props.hasMention}
                            hide={this.hideOverlay}
                            channelId={this.props.channelId}
                        />
                    </Overlay>
                    <a
                        className='mention-link'
                        onClick={this.handleClick}
                        ref={this.buttonRef}
                        aria-haspopup='dialog'
                        role='button'
                        tabIndex={0}
                        onKeyDown={this.handleKeyDown}
                    >
                        {'@' + displayName}
                    </a>
                </span>
                {suffix}
            </>
        );
    }
}
