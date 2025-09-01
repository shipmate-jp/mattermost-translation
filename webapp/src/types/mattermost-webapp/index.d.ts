// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

export interface PluginRegistry {
    registerPostTypeComponent(typeName: string, component: React.ElementType);

    // Register a new action in the post dropdown ("...") menu
    // The exact signature in Mattermost may accept additional parameters. Keep types permissive to avoid build issues.
    registerPostDropdownMenuAction?: (
        text: string,
        action: (postId: string) => void,
        filter?: (...args: any[]) => boolean
    ) => void;

    // Add more if needed from https://developers.mattermost.com/extend/plugins/webapp/reference
}
