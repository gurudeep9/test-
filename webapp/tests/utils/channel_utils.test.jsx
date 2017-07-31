// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import * as Utils from 'utils/channel_utils.jsx';
import TeamStore from 'stores/team_store.jsx';
import UserStore from 'stores/user_store.jsx';
import ChannelStore from 'stores/channel_store.jsx';
import Constants from 'utils/constants.jsx';

describe('Channel Utils', () => {
    describe('showDeleteOption', () => {
        test('all users can delete channels on unlicensed instances', () => {
            global.window.mm_license = {IsLicensed: 'false'};
            expect(Utils.showDeleteOptionForCurrentUser(null, true, true, true)).
                toEqual(true);
        });

        test('users cannot delete default channels', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            const channel = {name: Constants.DEFAULT_CHANNEL};
            expect(Utils.showDeleteOptionForCurrentUser(channel, true, true, true)).
                toEqual(false);
        });

        test('system admins can delete private channels, user is system admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPrivateChannelDeletion: Constants.PERMISSIONS_SYSTEM_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.PRIVATE_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, true)).
                toEqual(true);
        });

        test('system admins cannot delete private channels, user is not system admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPrivateChannelDeletion: Constants.PERMISSIONS_SYSTEM_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.PRIVATE_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, false)).
                toEqual(false);
        });

        test('system admins can delete public channels, user is system admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPublicChannelDeletion: Constants.PERMISSIONS_SYSTEM_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.OPEN_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, true)).
                toEqual(true);
        });

        test('system admins cannot delete public channels, user is not system admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPublicChannelDeletion: Constants.PERMISSIONS_SYSTEM_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.OPEN_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, false)).
                toEqual(false);
        });

        test('system admins or team admins can delete private channels, user is system admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPrivateChannelDeletion: Constants.PERMISSIONS_TEAM_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.PRIVATE_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, true)).
                toEqual(true);
        });

        test('system admins or team admins cannot delete private channels, user is not system admin or team admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPrivateChannelDeletion: Constants.PERMISSIONS_TEAM_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.PRIVATE_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, false)).
                toEqual(false);
        });

        test('system admins or team admins can delete public channels, user is system admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPublicChannelDeletion: Constants.PERMISSIONS_TEAM_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.OPEN_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, true)).
                toEqual(true);
        });

        test('system admins or team admins cannot delete public channels, user is not system admin or team admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPublicChannelDeletion: Constants.PERMISSIONS_TEAM_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.OPEN_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, false)).
                toEqual(false);
        });

        test('system admins or team admins can delete private channels, user is team admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPrivateChannelDeletion: Constants.PERMISSIONS_TEAM_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.PRIVATE_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, true, false)).
                toEqual(true);
        });

        test('system admins or team admins can delete public channels, user is team admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPublicChannelDeletion: Constants.PERMISSIONS_TEAM_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.OPEN_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, true, false)).
                toEqual(true);
        });

        test('channel, team, and system admins can delete public channels, user is channel admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPublicChannelDeletion: Constants.PERMISSIONS_CHANNEL_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.OPEN_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, true, false, false)).
                toEqual(true);
        });

        test('channel, team, and system admins can delete private channels, user is channel admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPrivateChannelDeletion: Constants.PERMISSIONS_CHANNEL_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.PRIVATE_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, true, false, false)).
                toEqual(true);
        });

        test('channel, team, and system admins can delete public channels, user is team admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPublicChannelDeletion: Constants.PERMISSIONS_CHANNEL_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.OPEN_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, true, false)).
                toEqual(true);
        });

        test('channel, team, and system admins can delete private channels, user is channel admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPrivateChannelDeletion: Constants.PERMISSIONS_CHANNEL_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.PRIVATE_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, true, false)).
                toEqual(true);
        });

        test('channel, team, and system admins can delete public channels, user is system admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPublicChannelDeletion: Constants.PERMISSIONS_CHANNEL_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.OPEN_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, true)).
                toEqual(true);
        });

        test('channel, team, and system admins can delete private channels, user is system admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPrivateChannelDeletion: Constants.PERMISSIONS_CHANNEL_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.PRIVATE_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, true)).
                toEqual(true);
        });

        test('channel, team, and system admins can delete public channels, user is not admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPublicChannelDeletion: Constants.PERMISSIONS_CHANNEL_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.OPEN_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, false)).
                toEqual(false);
        });

        test('channel, team, and system admins can delete private channels, user is channel admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPrivateChannelDeletion: Constants.PERMISSIONS_CHANNEL_ADMIN};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.PRIVATE_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, false)).
                toEqual(false);
        });

        test('any member can delete public channels, user is not admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPublicChannelDeletion: Constants.PERMISSIONS_ALL};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.OPEN_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, false)).
                toEqual(true);
        });

        test('any member can delete private channels, user is not admin test', () => {
            global.window.mm_license = {IsLicensed: 'true'};
            global.window.mm_config = {RestrictPrivateChannelDeletion: Constants.PERMISSIONS_ALL};

            const channel = {
                name: 'fakeChannelName',
                type: Constants.PRIVATE_CHANNEL
            };
            expect(Utils.showDeleteOptionForCurrentUser(channel, false, false, false)).
                toEqual(true);
        });
    });
});