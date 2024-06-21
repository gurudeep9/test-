// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {KatexOptions} from 'katex';
import React, {useEffect, useState} from 'react';
import {FormattedMessage} from 'react-intl';

type Katex = typeof import('katex');

type Props = {
    content: string;
    enableInlineLatex: boolean;
};

const LatexInline = ({
    content,
    enableInlineLatex,
}: Props) => {
    const [katex, setKatex] = useState<Katex | undefined>(undefined);

    useEffect(() => {
        import('katex').then((katex) => {
            setKatex(katex.default);
        });
    }, []);

    if (!enableInlineLatex || katex === undefined) {
        return (
            <span
                className='post-body--code inline-tex'
            >
                {'$' + content + '$'}
            </span>
        );
    }

    try {
        const katexOptions: KatexOptions = {
            throwOnError: false,
            displayMode: false,
            maxSize: 200,
            maxExpand: 100,
            fleqn: true,
        };

        const html = katex.renderToString(content, katexOptions);

        return (
            <span
                className='post-body--code inline-tex'
                dangerouslySetInnerHTML={{__html: html}}
            />
        );
    } catch (e) {
        return (
            <span
                className='post-body--code inline-tex'
            >
                <FormattedMessage
                    id='katex.error'
                    defaultMessage="Couldn't compile your Latex code. Please review the syntax and try again."
                />
            </span>
        );
    }
};

export default React.memo(LatexInline);
