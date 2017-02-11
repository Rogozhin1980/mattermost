// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import * as AsyncClient from 'utils/async_client.jsx';
import Client from 'client/web_client.jsx';

export function uploadFile(file, name, channelId, clientId, success, error, progress) {
    Client.uploadFile(
        file,
        name,
        channelId,
        clientId,
        (data) => {
            if (success) {
                success(data);
            }
        },
        (err) => {
            AsyncClient.dispatchError(err, 'uploadFile');

            if (error) {
                error(err);
            }
        },
        progress.bind(this, name)
    );
}
