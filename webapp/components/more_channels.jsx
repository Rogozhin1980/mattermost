// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import SearchableChannelList from './searchable_channel_list.jsx';

import ChannelStore from 'stores/channel_store.jsx';
import UserStore from 'stores/user_store.jsx';
import TeamStore from 'stores/team_store.jsx';

import Constants from 'utils/constants.jsx';
import * as AsyncClient from 'utils/async_client.jsx';
import {joinChannel, searchMoreChannels} from 'actions/channel_actions.jsx';

import React from 'react';
import PureRenderMixin from 'react-addons-pure-render-mixin';
import {Modal} from 'react-bootstrap';
import {FormattedMessage} from 'react-intl';
import {browserHistory} from 'react-router/es6';

const CHANNELS_CHUNK_SIZE = 50;
const CHANNELS_PER_PAGE = 50;
const SEARCH_TIMEOUT_MILLISECONDS = 100;

export default class MoreChannels extends React.Component {
    constructor(props) {
        super(props);

        this.onChange = this.onChange.bind(this);
        this.handleJoin = this.handleJoin.bind(this);
        this.handleHide = this.handleHide.bind(this);
        this.handleExit = this.handleExit.bind(this);
        this.nextPage = this.nextPage.bind(this);
        this.search = this.search.bind(this);

        this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this);

        this.searchTimeoutId = 0;

        this.state = {
            show: true,
            search: false,
            channels: null,
            serverError: null
        };
    }

    componentDidMount() {
        ChannelStore.addChangeListener(this.onChange);
        AsyncClient.getMoreChannelsPage(0, CHANNELS_CHUNK_SIZE * 2);
    }

    componentWillUnmount() {
        ChannelStore.removeChangeListener(this.onChange);
    }

    handleHide() {
        this.setState({show: false});
    }

    handleExit() {
        if (this.props.onModalDismissed) {
            this.props.onModalDismissed();
        }
    }

    onChange(force) {
        if (this.state.search && !force) {
            return;
        }

        this.setState({
            channels: ChannelStore.getMoreChannelsList(),
            serverError: null
        });
    }

    nextPage(page) {
        AsyncClient.getMoreChannelsPage((page + 1) * CHANNELS_PER_PAGE, CHANNELS_PER_PAGE);
    }

    handleJoin(channel, done) {
        joinChannel(
            channel,
            () => {
                browserHistory.push(TeamStore.getCurrentTeamRelativeUrl() + '/channels/' + channel.name);
                if (done) {
                    done();
                }

                this.handleHide();
            },
            (err) => {
                this.setState({serverError: err.message});
                if (done) {
                    done();
                }
            }
        );
    }

    search(term) {
        if (term === '') {
            this.onChange(true);
            this.setState({search: false});
            return;
        }

        clearTimeout(this.searchTimeoutId);

        this.searchTimeoutId = setTimeout(
            () => {
                searchMoreChannels(
                    term,
                    (channels) => {
                        this.setState({search: true, channels});
                    }
                );
            },
            SEARCH_TIMEOUT_MILLISECONDS
        );
    }

    render() {
        let serverError;
        if (this.state.serverError) {
            serverError = <div className='form-group has-error'><label className='control-label'>{this.state.serverError}</label></div>;
        }

        let createNewChannelButton = (
            <button
                type='button'
                className='btn btn-primary channel-create-btn'
                onClick={this.props.handleNewChannel}
            >
                <FormattedMessage
                    id='more_channels.create'
                    defaultMessage='Create New Channel'
                />
            </button>
        );

        let createChannelHelpText = (
            <p className='secondary-message'>
                <FormattedMessage
                    id='more_channels.createClick'
                    defaultMessage="Click 'Create New Channel' to make a new one"
                />
            </p>
        );

        const isAdmin = TeamStore.isTeamAdminForCurrentTeam() || UserStore.isSystemAdminForCurrentUser();
        const isSystemAdmin = UserStore.isSystemAdminForCurrentUser();

        if (global.window.mm_license.IsLicensed === 'true') {
            if (global.window.mm_config.RestrictPublicChannelManagement === Constants.PERMISSIONS_SYSTEM_ADMIN && !isSystemAdmin) {
                createNewChannelButton = null;
                createChannelHelpText = null;
            } else if (global.window.mm_config.RestrictPublicChannelManagement === Constants.PERMISSIONS_TEAM_ADMIN && !isAdmin) {
                createNewChannelButton = null;
                createChannelHelpText = null;
            }
        }

        return (
            <Modal
                dialogClassName='more-modal more-modal--action'
                show={this.state.show}
                onHide={this.handleHide}
                onExited={this.handleExit}
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        <FormattedMessage
                            id='more_channels.title'
                            defaultMessage='More Channels'
                        />
                    </Modal.Title>
                    {createNewChannelButton}
                </Modal.Header>
                <Modal.Body>
                    <SearchableChannelList
                        channels={this.state.channels}
                        channelsPerPage={CHANNELS_PER_PAGE}
                        nextPage={this.nextPage}
                        search={this.search}
                        handleJoin={this.handleJoin}
                        noResultsText={createChannelHelpText}
                    />
                    {serverError}
                </Modal.Body>
            </Modal>
        );
    }
}

MoreChannels.propTypes = {
    onModalDismissed: React.PropTypes.func,
    handleNewChannel: React.PropTypes.func
};