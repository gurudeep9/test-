// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {expect, Locator} from '@playwright/test';

import {NotificationsSettings} from './notification_settings';

export default class AccountSettingsModal {
    readonly container: Locator;

    readonly notificationsSettingsTab;
    readonly notificationsSettings;

    constructor(container: Locator) {
        this.container = container;

        this.notificationsSettingsTab = this.container.locator('#notificationsButton');
        this.notificationsSettings = new NotificationsSettings(container.locator('#notificationSettings'));
    }

    async toBeVisible() {
        await expect(this.container).toBeVisible();
    }

    async openNotificationsTab() {
        await expect(this.notificationsSettingsTab).toBeVisible();
        await this.notificationsSettingsTab.click();
    
        await this.notificationsSettings.toBeVisible();
    }

}

export {AccountSettingsModal};
