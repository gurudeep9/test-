// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import {intlShape, injectIntl, FormattedHTMLMessage, defineMessages} from 'react-intl';

import LoadingScreen from '../loading_screen.jsx';

import ChannelStore from '../../stores/channel_store.jsx';

import * as Client from '../../utils/client.jsx';
import Constants from '../../utils/constants.jsx';

const messages = defineMessages({
    none: {
        id: 'user.settings.hooks_out.none',
        defaultMessage: 'None'
    },
    select: {
        id: 'user.settings.hooks_out.select',
        defaultMessage: '--- Select a channel ---'
    },
    channel: {
        id: 'user.settings.hooks_out.channel',
        defaultMessage: 'Channel: '
    },
    trigger: {
        id: 'user.settings.hooks_out.trigger',
        defaultMessage: 'Trigger Words: '
    },
    regen: {
        id: 'user.settings.hooks_out.regen',
        defaultMessage: 'Regen Token'
    },
    existing: {
        id: 'user.settings.hooks_out.existing',
        defaultMessage: 'Existing outgoing webhooks'
    },
    addDescription: {
        id: 'user.settings.hooks_out.addDescription',
        defaultMessage: 'Create webhooks to send new message events to an external integration. Please see <a href="http://mattermost.org/webhooks">http://mattermost.org/webhooks</a>  to learn more.'
    },
    addTitle: {
        id: 'user.settings.hooks_out.addTitle',
        defaultMessage: 'Add a new outgoing webhook'
    },
    add: {
        id: 'user.settings.hooks_out.add',
        defaultMessage: 'Add'
    },
    only: {
        id: 'user.settings.hooks_out.only',
        defaultMessage: 'Only public channels can be used'
    },
    optional: {
        id: 'user.settings.hooks_out.optional',
        defaultMessage: 'Optional if channel selected'
    },
    comma: {
        id: 'user.settings.hooks_out.comma',
        defaultMessage: 'Comma separated words to trigger on'
    },
    callback: {
        id: 'user.settings.hooks_out.callback',
        defaultMessage: 'Callback URLs:'
    },
    callbackHolder: {
        id: 'user.settings.hooks_out.callbackHolder',
        defaultMessage: 'Each URL must start with http:// or https://'
    },
    callbackDesc: {
        id: 'user.settings.hooks_out.callbackDesc',
        defaultMessage: 'New line separated URLs that will receive the HTTP POST event'
    }
});

class ManageOutgoingHooks extends React.Component {
    constructor() {
        super();

        this.getHooks = this.getHooks.bind(this);
        this.addNewHook = this.addNewHook.bind(this);
        this.updateChannelId = this.updateChannelId.bind(this);
        this.updateTriggerWords = this.updateTriggerWords.bind(this);
        this.updateCallbackURLs = this.updateCallbackURLs.bind(this);

        this.state = {hooks: [], channelId: '', triggerWords: '', callbackURLs: '', getHooksComplete: false};
    }
    componentDidMount() {
        this.getHooks();
    }
    addNewHook(e) {
        e.preventDefault();

        if ((this.state.channelId === '' && this.state.triggerWords === '') ||
                this.state.callbackURLs === '') {
            return;
        }

        const hook = {};
        hook.channel_id = this.state.channelId;
        if (this.state.triggerWords.length !== 0) {
            hook.trigger_words = this.state.triggerWords.trim().split(',');
        }
        hook.callback_urls = this.state.callbackURLs.split('\n');

        Client.addOutgoingHook(
            hook,
            (data) => {
                let hooks = Object.assign([], this.state.hooks);
                if (!hooks) {
                    hooks = [];
                }
                hooks.push(data);
                this.setState({hooks, addError: null, channelId: '', triggerWords: '', callbackURLs: ''});
            },
            (err) => {
                this.setState({addError: err.message});
            }
        );
    }
    removeHook(id) {
        const data = {};
        data.id = id;

        Client.deleteOutgoingHook(
            data,
            () => {
                const hooks = this.state.hooks;
                let index = -1;
                for (let i = 0; i < hooks.length; i++) {
                    if (hooks[i].id === id) {
                        index = i;
                        break;
                    }
                }

                if (index !== -1) {
                    hooks.splice(index, 1);
                }

                this.setState({hooks});
            },
            (err) => {
                this.setState({editError: err.message});
            }
        );
    }
    regenToken(id) {
        const regenData = {};
        regenData.id = id;

        Client.regenOutgoingHookToken(
            regenData,
            (data) => {
                const hooks = Object.assign([], this.state.hooks);
                for (let i = 0; i < hooks.length; i++) {
                    if (hooks[i].id === id) {
                        hooks[i] = data;
                        break;
                    }
                }

                this.setState({hooks, editError: null});
            },
            (err) => {
                this.setState({editError: err.message});
            }
        );
    }
    getHooks() {
        Client.listOutgoingHooks(
            (data) => {
                if (data) {
                    this.setState({hooks: data, getHooksComplete: true, editError: null});
                }
            },
            (err) => {
                this.setState({editError: err.message});
            }
        );
    }
    updateChannelId(e) {
        this.setState({channelId: e.target.value});
    }
    updateTriggerWords(e) {
        this.setState({triggerWords: e.target.value});
    }
    updateCallbackURLs(e) {
        this.setState({callbackURLs: e.target.value});
    }
    render() {
        const {formatMessage} = this.props.intl;

        let addError;
        if (this.state.addError) {
            addError = <label className='has-error'>{this.state.addError}</label>;
        }
        let editError;
        if (this.state.editError) {
            editError = <label className='has-error'>{this.state.editError}</label>;
        }

        const channels = ChannelStore.getAll();
        const options = [];
        options.push(
            <option
                key='select-channel'
                value=''
            >
                {formatMessage(messages.select)}
            </option>
        );

        channels.forEach((channel) => {
            if (channel.type === Constants.OPEN_CHANNEL) {
                options.push(
                    <option
                        key={'outgoing-hook' + channel.id}
                        value={channel.id}
                    >
                        {channel.display_name}
                    </option>
                );
            }
        });

        const hooks = [];
        this.state.hooks.forEach((hook) => {
            const c = ChannelStore.get(hook.channel_id);

            if (!c && hook.channel_id && hook.channel_id.length !== 0) {
                return;
            }

            let channelDiv;
            if (c) {
                channelDiv = (
                    <div className='padding-top'>
                        <strong>{formatMessage(messages.channel)}</strong>{c.display_name}
                    </div>
                );
            }

            let triggerDiv;
            if (hook.trigger_words && hook.trigger_words.length !== 0) {
                triggerDiv = (
                    <div className='padding-top'>
                        <strong>{formatMessage(messages.trigger)}</strong>{hook.trigger_words.join(', ')}
                    </div>
                );
            }

            hooks.push(
                <div
                    key={hook.id}
                    className='webhook__item'
                >
                    <div className='padding-top x2 webhook__url'>
                        <strong>{'URLs: '}</strong><span className='word-break--all'>{hook.callback_urls.join(', ')}</span>
                    </div>
                    {channelDiv}
                    {triggerDiv}
                    <div className='padding-top'>
                        <strong>{'Token: '}</strong>{hook.token}
                    </div>
                    <div className='padding-top'>
                        <a
                            className='text-danger'
                            href='#'
                            onClick={this.regenToken.bind(this, hook.id)}
                        >
                            {formatMessage(messages.regen)}
                        </a>
                        <a
                            className='webhook__remove'
                            href='#'
                            onClick={this.removeHook.bind(this, hook.id)}
                        >
                            <span aria-hidden='true'>{'×'}</span>
                        </a>
                    </div>
                    <div className='padding-top x2 divider-light'></div>
                </div>
            );
        });

        let displayHooks;
        if (!this.state.getHooksComplete) {
            displayHooks = <LoadingScreen/>;
        } else if (hooks.length > 0) {
            displayHooks = hooks;
        } else {
            displayHooks = <div className='padding-top x2'>{formatMessage(messages.none)}</div>;
        }

        const existingHooks = (
            <div className='webhooks__container'>
                <label className='control-label padding-top x2'>{formatMessage(messages.existing)}</label>
                <div className='padding-top divider-light'></div>
                <div className='webhooks__list'>
                    {displayHooks}
                </div>
            </div>
        );

        const disableButton = (this.state.channelId === '' && this.state.triggerWords === '') || this.state.callbackURLs === '';

        return (
            <div key='addOutgoingHook'>
                <FormattedHTMLMessage id='user.settings.hooks_out.addDescription' />
                <div><label className='control-label padding-top x2'>{formatMessage(messages.addTitle)}</label></div>
                <div className='padding-top divider-light'></div>
                <div className='padding-top'>
                    <div>
                        <label className='control-label'>{formatMessage(messages.channel).replace(': ', '')}</label>
                        <div className='padding-top'>
                            <select
                                ref='channelName'
                                className='form-control'
                                value={this.state.channelId}
                                onChange={this.updateChannelId}
                            >
                                {options}
                            </select>
                        </div>
                        <div className='padding-top'>{formatMessage(messages.only)}</div>
                    </div>
                    <div className='padding-top x2'>
                        <label className='control-label'>{formatMessage(messages.trigger)}</label>
                        <div className='padding-top'>
                            <input
                                ref='triggerWords'
                                className='form-control'
                                value={this.state.triggerWords}
                                onChange={this.updateTriggerWords}
                                placeholder={formatMessage(messages.optional)}
                            />
                        </div>
                        <div className='padding-top'>{formatMessage(messages.comma)}</div>
                    </div>
                    <div className='padding-top x2'>
                        <label className='control-label'>{formatMessage(messages.callback)}</label>
                        <div className='padding-top'>
                        <textarea
                            ref='callbackURLs'
                            className='form-control no-resize'
                            value={this.state.callbackURLs}
                            resize={false}
                            rows={3}
                            onChange={this.updateCallbackURLs}
                            placeholder={formatMessage(messages.callbackHolder)}
                        />
                        </div>
                        <div className='padding-top'>{formatMessage(messages.callbackDesc)}</div>
                        {addError}
                    </div>
                    <div className='padding-top padding-bottom'>
                        <a
                            className={'btn btn-sm btn-primary'}
                            href='#'
                            disabled={disableButton}
                            onClick={this.addNewHook}
                        >
                            {formatMessage(messages.add)}
                        </a>
                    </div>
                </div>
                {existingHooks}
                {editError}
            </div>
        );
    }
}

ManageOutgoingHooks.propTypes = {
    intl: intlShape.isRequired
};

export default injectIntl(ManageOutgoingHooks);