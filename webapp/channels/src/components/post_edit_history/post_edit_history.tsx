// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {memo, useEffect, useRef} from 'react';
import {useIntl} from 'react-intl';
import Scrollbars from 'react-custom-scrollbars';

import SearchResultsHeader from 'components/search_results_header';

import LoadingScreen from 'components/loading_screen';

import EditedPostItem from './edited_post_item';

import {PropsFromRedux} from '.';
import AlertIcon from 'components/common/svg_images_components/alert_svg';

import './post_edit_history.scss';

const renderView = (props: Record<string, unknown>): JSX.Element => (
    <div
        {...props}
        className='scrollbar--view'
    />
);

const renderThumbHorizontal = (props: Record<string, unknown>): JSX.Element => (
    <div
        {...props}
        className='scrollbar--horizontal'
    />
);

const renderThumbVertical = (props: Record<string, unknown>): JSX.Element => (
    <div
        {...props}
        className='scrollbar--vertical'
    />
);

const PostEditHistory = ({
    channelDisplayName,
    originalPost,
    postEditHistory,
    errors,
}: PropsFromRedux) => {
    const scrollbars = useRef<Scrollbars | null>(null);
    const {formatMessage} = useIntl();
    const retrieveErrorHeading = formatMessage({
        id: 'post_info.edit.history.retrieveError',
        defaultMessage: 'Unable to load edit history',
    });
    const retrieveErrorSubheading = formatMessage({
        id: 'post_info.edit.history.retrieveErrorVerbose',
        defaultMessage: 'There was an error loading the history for this message. Check your network connection or try again later.',
    });

    useEffect(() => {
        scrollbars.current?.scrollToTop();
    }, [originalPost, postEditHistory]);

    const title = formatMessage({
        id: 'search_header.title_edit.history',
        defaultMessage: 'Edit History',
    });

    const errorContainer: JSX.Element = (
        <div className='edit-post-history__error_container'>
            <div className='edit-post-history__error_item'>
                <AlertIcon
                    width={127}
                    height={127}
                />
                <p className='edit-post-history__error_heading'>
                    {retrieveErrorHeading}
                </p>
                <p className='edit-post-history__error_subheading'>
                    {retrieveErrorSubheading}
                </p>
            </div>
        </div>
    );

    if (postEditHistory.length === 0 && !errors) {
        return (
            <div
                id='rhsContainer'
                className='sidebar-right__body'
            >
                <LoadingScreen
                    style={{
                        display: 'grid',
                        placeContent: 'center',
                        flex: '1',
                    }}
                />
            </div>
        );
    }

    const currentItem = (
        <EditedPostItem
            post={originalPost}
            key={originalPost.id}
            isCurrent={true}
        />
    );

    const postEditItems = [currentItem, ...postEditHistory.map((postEdited) => (
        <EditedPostItem
            key={postEdited.id}
            post={postEdited}
        />
    ))];

    return (
        <div
            id='rhsContainer'
            className='sidebar-right__body'
        >
            <Scrollbars
                ref={scrollbars}
                autoHide={true}
                autoHideTimeout={500}
                autoHideDuration={500}
                renderThumbHorizontal={renderThumbHorizontal}
                renderThumbVertical={renderThumbVertical}
                renderView={renderView}
                style={{display: 'flex', flexDirection: 'column'}}
                
            >
                <SearchResultsHeader>
                    {title}
                    <div className='sidebar--right__title__channel'>{channelDisplayName}</div>
                </SearchResultsHeader>
                {errors ? errorContainer : postEditItems}
            </Scrollbars>
        </div>
    );
};

export default memo(PostEditHistory);
