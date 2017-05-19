import PropTypes from 'prop-types';

// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import React from 'react';

import Constants from 'utils/constants.jsx';
import PureRenderMixin from 'react-addons-pure-render-mixin';

import {getDateForUnixTicks, isMobile, updateWindowDimensions, openInNewTab} from 'utils/utils.jsx';

import {Link} from 'react-router/es6';
import TeamStore from 'stores/team_store.jsx';

export default class PostTime extends React.Component {
    constructor(props) {
        super(props);

        this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this);
        this.state = {
            currentTeamDisplayName: TeamStore.getCurrent().name,
            width: '',
            height: ''
        };
    }

    componentDidMount() {
        this.intervalId = setInterval(() => {
            this.forceUpdate();
        }, Constants.TIME_SINCE_UPDATE_INTERVAL);
        window.addEventListener('resize', () => {
            updateWindowDimensions(this);
        });
    }

    componentWillUnmount() {
        clearInterval(this.intervalId);
        window.removeEventListener('resize', () => {
            updateWindowDimensions(this);
        });
    }

    renderTimeTag() {
        const date = getDateForUnixTicks(this.props.eventTime);

        return (
            <time
                className='post__time'
                dateTime={date.toISOString()}
                title={date}
            >
                {date.toLocaleString('en', {hour: '2-digit', minute: '2-digit', hour12: !this.props.useMilitaryTime})}
            </time>
        );
    }

    render() {
        return isMobile() ?
            this.renderTimeTag() :
            (
                <Link
                    onClick={openInNewTab.bind(this, `/${this.state.currentTeamDisplayName}/pl/${this.props.postId}`)}
                    className='post__permalink'
                >
                    {this.renderTimeTag()}
                </Link>
            );
    }
}

PostTime.defaultProps = {
    eventTime: 0,
    sameUser: false
};

PostTime.propTypes = {
    eventTime: PropTypes.number.isRequired,
    sameUser: PropTypes.bool,
    compactDisplay: PropTypes.bool,
    useMilitaryTime: PropTypes.bool.isRequired,
    postId: PropTypes.string
};
