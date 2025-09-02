// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import manifest from '@/manifest';
import type {PluginRegistry} from '@/types/mattermost-webapp';

// Helper to read a cookie value by name
function getCookie(name: string): string | undefined {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) {
        return parts.pop()?.split(';')?.shift();
    }
    return undefined;
}

function getBasePath(): string {
    if (typeof window !== 'undefined') {
        const w = window as unknown as { basename?: string };
        if (typeof w.basename === 'string') {
            return w.basename || '';
        }
    }
    return '';
}

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry) {
    // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        if (registry.registerPostDropdownMenuAction) {
            registry.registerPostDropdownMenuAction(
                'Translate message',
                async (postId: string) => {
                    try {
                        const csrfToken = getCookie('MMCSRF');
                        const basePath = getBasePath();
                        await fetch(`${basePath}/plugins/${manifest.id}/api/v1/translate`, {
                            method: 'POST',
                            credentials: 'same-origin',
                            headers: {
                                'Content-Type': 'application/json',
                                'X-Requested-With': 'XMLHttpRequest',
                                ...(csrfToken ? {'X-CSRF-Token': csrfToken} : {}),
                            },
                            body: JSON.stringify({post_id: postId}),
                        });
                    } catch (e) {
                        // no-op
                    }
                },
            );
        }
    }
}

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void;
        basename?: string;
    }
}

window.registerPlugin(manifest.id, new Plugin());
