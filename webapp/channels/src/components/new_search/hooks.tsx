// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState, useRef, useEffect, useMemo} from 'react';
import {FormattedMessage} from 'react-intl';
import {useDispatch} from 'react-redux';

import type {Channel} from '@mattermost/types/channels';
import type {ServerError} from '@mattermost/types/errors';

import {autocompleteChannelsForSearch} from 'actions/channel_actions';
import {autocompleteUsersInTeam} from 'actions/user_actions';

import type {ProviderResult} from 'components/suggestion/provider';
import type Provider from 'components/suggestion/provider';
import SearchChannelProvider from 'components/suggestion/search_channel_provider';
import SearchDateProvider from 'components/suggestion/search_date_provider';
import SearchUserProvider from 'components/suggestion/search_user_provider';

import {SearchFileExtensionProvider} from './extension_suggestions_provider';

const useSearchSuggestions = (searchType: string, searchTerms: string, caretPosition: number, getCaretPosition: () => number, setSelectedOption: (idx: number) => void): [ProviderResult<unknown>|null, React.ReactNode] => {
    const dispatch = useDispatch();

    const [providerResults, setProviderResults] = useState<ProviderResult<unknown>|null>(null);
    const [suggestionsHeader, setSuggestionsHeader] = useState<React.ReactNode>(<span/>);

    const suggestionProviders = useRef<Provider[]>([
        new SearchDateProvider(),
        new SearchChannelProvider((term: string, success?: (channels: Channel[]) => void, error?: (err: ServerError) => void) => dispatch(autocompleteChannelsForSearch(term, success, error))),
        new SearchUserProvider((username: string) => dispatch(autocompleteUsersInTeam(username))),
        new SearchFileExtensionProvider(),
    ]);

    const headers = useMemo<React.ReactNode[]>(() => [
        <span/>,
        <FormattedMessage
            id='search_bar.channels'
            defaultMessage='Channels'
        />,
        <FormattedMessage
            id='search_bar.users'
            defaultMessage='Users'
        />,
        <FormattedMessage
            id='search_bar.file_types'
            defaultMessage='File types'
        />
    ], []);

    useEffect(() => {
        setProviderResults(null);
        if (searchType !== '' && searchType !== 'messages' && searchType !== 'files') {
            return;
        }

        let partialSearchTerms = searchTerms.slice(0, caretPosition);
        if (searchTerms.length > caretPosition && searchTerms[caretPosition] !== ' ') {
            return;
        }

        if (caretPosition > 0 && searchTerms[caretPosition - 1] === ' ') {
            return;
        }

        suggestionProviders.current.forEach((provider, idx) => {
            provider.handlePretextChanged(partialSearchTerms, (res: ProviderResult<unknown>) => {
                if (idx ===  3 && searchType !== 'files') {
                    return;
                }
                if (caretPosition !== getCaretPosition()) {
                    return;
                }
                res.items = res.items.slice(0, 10);
                res.terms = res.terms.slice(0, 10);
                setProviderResults(res);
                setSelectedOption(0);
                setSuggestionsHeader(headers[idx]);
            })
        });
    }, [searchTerms, searchType, caretPosition]);

    return [providerResults, suggestionsHeader];
};

export default useSearchSuggestions;
