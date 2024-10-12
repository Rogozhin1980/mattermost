// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useCallback} from 'react';
import {useIntl} from 'react-intl';

import SectionNotice from 'components/section_notice';

export default function NotificationPermissionUnsupportedSectionNotice() {
    const intl = useIntl();

    const handleClick = useCallback(async () => {
        // TODO: Change to permalink
        window.open('https://docs.mattermost.com/install/software-hardware-requirements.html#pc-web', '_blank', 'noopener,noreferrer');
    }, []);

    return (
        <div className='extraContentBeforeSettingList'>
            <SectionNotice
                type='danger'
                title={intl.formatMessage({
                    id: 'user.settings.notifications.desktopAndMobile.notificationSection.permissionUnsupported.title',
                    defaultMessage: 'Web browser notifications unsupported',
                })}
                text={intl.formatMessage({
                    id: 'user.settings.notifications.desktopAndMobile.notificationSection.permissionUnsupported.message',
                    defaultMessage: 'You\'re missing important message and call notifications from Mattermost. To start receiving notifications, please update to a supported browser.',
                })}
                tertiaryButton={{
                    text: intl.formatMessage({
                        id: 'user.settings.notifications.desktopAndMobile.notificationSection.permissionUnsupported.button',
                        defaultMessage: 'Update your browser',
                    }),
                    onClick: handleClick,
                }}
            />
        </div>
    );
}

