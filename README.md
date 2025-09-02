# Mattermost Translate

Translate any message in Mattermost with one click using Google Cloud Translation API v2.

-   Adds a “🌐 Translate message” action to each post menu
-   Sends the translation back to you as an ephemeral message
-   Replies in-thread when the original message is threaded
-   Server-side configuration for API key and default target language

## Requirements

-   Mattermost Server 6.2.1+ (tested on newer versions as well)
-   A Google Cloud project with the Cloud Translation API enabled
-   A Google API key with access to Cloud Translation API v2

## Installation

### Install from a release (recommended)

1. Download the plugin bundle tarball from the Releases page (jp.shipmate.mattermost-translate-<version>.tar.gz).
2. In Mattermost System Console, go to:
    - System Console → Plugin Management → Upload Plugin → Upload the tarball
3. Enable the plugin.

### Install from source

Prerequisites: Go 1.20+, Node.js 16+, npm 8+, make

-   Build and package:
    -   make dist
    -   The bundle is created at dist/jp.shipmate.mattermost-translate-<version>.tar.gz
-   Upload via System Console as above, or deploy to a local server using environment variables:
    -   Using personal access token:
        -   export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
        -   export MM_ADMIN_TOKEN=<your_token>
        -   make deploy
    -   Using username/password:
        -   export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
        -   export MM_ADMIN_USERNAME=admin
        -   export MM_ADMIN_PASSWORD=<password>
        -   make deploy

## Configuration

After installation, configure the plugin in System Console → Plugins → Mattermost Translate.

-   GoogleAPIKey (required)
    -   Your Google API key for Cloud Translation API v2.
    -   Create/Manage at https://console.cloud.google.com/apis/credentials and enable the Translation API.
-   DefaultTargetLang (optional)
    -   ISO language code to translate into by default (e.g. en, ja, es). Default: ja

Click Save when done.

## Usage

-   Hover a message, open the “…” post menu, and click “🌐 Translate message”.
-   The translation will be sent back to you as an ephemeral message.
-   If the original message is part of a thread, the reply appears in that thread.

Notes

-   If you want to override the target language per request, this can be added in a future enhancement. Currently the server uses DefaultTargetLang when no target is specified by the client.

## Troubleshooting

-   401 Unauthorized when calling the API
    -   Ensure you are logged in and accessing Mattermost from the same origin (no cross-domain proxy).
    -   Clear cache/hard refresh to load the updated webapp.
-   Translation returns the original text
    -   Verify GoogleAPIKey is valid and the Translation API is enabled in your project.
    -   Check DefaultTargetLang is set correctly (e.g. en to translate Japanese → English).
-   The menu shows but no response
    -   Check System Console → Logs for plugin errors.
    -   Re-upload the latest bundle and refresh the client.

## Development

-   Lint & type-check:
    -   make check-style
-   Build:
    -   make dist
-   Watch & auto-deploy to a local server (requires MM_ADMIN_TOKEN and MM_SERVICESETTINGS_SITEURL):
    -   export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
    -   export MM_ADMIN_TOKEN=<your_token>
    -   make watch

Project structure

-   server/: Go backend (Mattermost plugin hooks and REST API)
-   webapp/: Frontend bundle that registers the post menu action

## License

This project is open source. See LICENSE for details.
