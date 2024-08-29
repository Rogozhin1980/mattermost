// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {ComponentProps} from 'react';
import React from 'react';
import {DragDropContext, Draggable, Droppable} from 'react-beautiful-dnd';
import styled from 'styled-components';

import type {ChannelBookmark} from '@mattermost/types/channel_bookmarks';
import type {IDMappedObjects} from '@mattermost/types/utilities';

import BookmarkItem from './bookmark_item';
import BookmarksMenu from './channel_bookmarks_menu';
import {useChannelBookmarkPermission, useChannelBookmarks, MAX_BOOKMARKS_PER_CHANNEL, useCanUploadFiles} from './utils';

import './channel_bookmarks.scss';

type Props = {
    channelId: string;
};

function ChannelBookmarks({
    channelId,
}: Props) {
    const {order, bookmarks, reorder} = useChannelBookmarks(channelId);
    const canUploadFiles = useCanUploadFiles();
    const canAdd = useChannelBookmarkPermission(channelId, 'add');
    const hasBookmarks = Boolean(order?.length);
    const limitReached = order.length >= MAX_BOOKMARKS_PER_CHANNEL;

    if (!hasBookmarks && !canAdd) {
        return null;
    }

    const handleOnDragEnd: ComponentProps<typeof DragDropContext>['onDragEnd'] = async ({source, destination, draggableId}) => {
        if (destination) {
            await reorder(draggableId, source.index, destination.index);
        }
    };

    return (
        <DragDropContext
            onDragEnd={handleOnDragEnd}
        >
            <Droppable
                droppableId='channel-bookmarks'
                direction='horizontal'
            >
                {(drop, snap) => {
                    return (
                        <Container
                            ref={drop.innerRef}
                            data-testid='channel-bookmarks-container'
                            {...drop.droppableProps}
                        >
                            {order.map(makeItemRenderer(bookmarks, snap.isDraggingOver))}
                            {drop.placeholder}
                            <BookmarksMenu
                                channelId={channelId}
                                hasBookmarks={hasBookmarks}
                                limitReached={limitReached}
                                canUploadFiles={canUploadFiles}
                            />
                        </Container>
                    );
                }}
            </Droppable>
        </DragDropContext>
    );
}

const makeItemRenderer = (bookmarks: IDMappedObjects<ChannelBookmark>, disableInteractions: boolean) => (id: string, index: number) => {
    return (
        <Draggable
            key={id}
            draggableId={id}
            index={index}
        >
            {(drag, snap) => {
                return (
                    <BookmarkItem
                        key={id}
                        drag={drag}
                        isDragging={snap.isDragging}
                        disableInteractions={disableInteractions}
                        bookmark={bookmarks[id]}
                    />
                );
            }}
        </Draggable>
    );
};

export default ChannelBookmarks;

const Container = styled.div`
    display: flex;
    padding: 0px 6px;
    min-height: 40px;
    align-items: center;
    border-bottom: 1px solid rgba(var(--center-channel-color-rgb), 0.12);
    overflow-x: auto;
`;
