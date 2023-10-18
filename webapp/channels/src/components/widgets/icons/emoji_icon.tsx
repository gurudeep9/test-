// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {memo} from 'react';

type Props = React.HTMLAttributes<HTMLSpanElement>;

const EmojiIcon = (props: Props) => {
    return (
        <span {...props}>
            <svg
                width='16px'
                height='16px'
                viewBox='0 0 16 16'
                role='img'
                aria-label='Emoji icon'
            >
                <path d='M9.872 8.00005C10.184 8.00005 10.448 7.88605 10.664 7.65805C10.892 7.43005 11.006 7.16605 11.006 6.86605C11.006 6.56605 10.892 6.30805 10.664 6.09205C10.448 5.86405 10.184 5.75005 9.872 5.75005C9.572 5.75005 9.308 5.86405 9.08 6.09205C8.864 6.30805 8.756 6.56605 8.756 6.86605C8.756 7.16605 8.864 7.43005 9.08 7.65805C9.308 7.88605 9.572 8.00005 9.872 8.00005ZM5.372 8.00005C5.684 8.00005 5.948 7.88605 6.164 7.65805C6.392 7.43005 6.506 7.16605 6.506 6.86605C6.506 6.56605 6.392 6.30805 6.164 6.09205C5.948 5.86405 5.684 5.75005 5.372 5.75005C5.072 5.75005 4.808 5.86405 4.58 6.09205C4.364 6.30805 4.256 6.56605 4.256 6.86605C4.256 7.16605 4.364 7.43005 4.58 7.65805C4.808 7.88605 5.072 8.00005 5.372 8.00005ZM13.22 7.92805C13.244 8.12005 13.256 8.27005 13.256 8.37805C13.256 9.39805 12.998 10.34 12.482 11.204C11.99 12.056 11.318 12.728 10.466 13.22C9.59 13.736 8.642 13.994 7.622 13.994C6.602 13.994 5.654 13.736 4.778 13.22C3.938 12.728 3.266 12.056 2.762 11.204C2.258 10.34 2.006 9.39805 2.006 8.37805C2.006 7.35805 2.264 6.41005 2.78 5.53405C3.272 4.69405 3.944 4.02205 4.796 3.51805C5.66 3.00205 6.602 2.74405 7.622 2.74405C7.73 2.74405 7.88 2.75605 8.072 2.78005C8.144 2.27605 8.306 1.79005 8.558 1.32205C8.15 1.27405 7.838 1.25005 7.622 1.25005C6.338 1.25005 5.138 1.57405 4.022 2.22205C2.954 2.84605 2.102 3.69805 1.466 4.77805C0.818 5.88205 0.494 7.08205 0.494 8.37805C0.494 9.67405 0.818 10.874 1.466 11.978C2.102 13.058 2.954 13.91 4.022 14.534C5.138 15.182 6.338 15.506 7.622 15.506C8.906 15.506 10.106 15.182 11.222 14.534C12.29 13.898 13.142 13.046 13.778 11.978C14.426 10.862 14.75 9.66205 14.75 8.37805C14.75 8.16205 14.726 7.85005 14.678 7.44205C14.21 7.69405 13.724 7.85605 13.22 7.92805ZM11.744 0.494048H13.256V2.74405H15.506V4.25605H13.256V6.50605H11.744V4.25605H9.494V2.74405H11.744V0.494048ZM3.788 9.87205C4.088 10.652 4.586 11.288 5.282 11.78C5.978 12.26 6.758 12.5 7.622 12.5C8.486 12.5 9.266 12.26 9.962 11.78C10.658 11.288 11.162 10.652 11.474 9.87205H3.788Z'/>
            </svg>
        </span>
    );
};

export default memo(EmojiIcon);
