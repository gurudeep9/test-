// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {ReactNode} from 'react';
import classNames from 'classnames';

type Props = {
    withBackButton?: boolean;
    children: ReactNode;
};

const AdminHeader = (props: Props) => {
    return (
        <div
            className={
                classNames(
                    'admin-console__header',
                    {'with-back': props.withBackButton},
                )
            }
        >
            {props.children}
        </div>
    );
};

export default AdminHeader;
