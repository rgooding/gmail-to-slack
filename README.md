# Gmail to Slack

This simple tool gets unread messages with specified labels from a Gmail account
and posts them to Slack channels. Messages are marked as read after posting to 
Slack and can optionally also be archived (removed from the Inbox).

By default the application looks for config.yaml in the current working directory
but this can be overridden by setting an alternative path in the CONFIG_FILE 
environment variable.

## Installation and Configuration

### Building from source
To build just download the source and run `go build`. Currently it has only been
tested with go 1.16 on Linux.

### Configuration

There are some prerequisites which need to be set up before configuring the tool.

#### Slack webhook URL
This is a standard incoming-webhook URL for Slack. The channel associated with
the URL does not matter as it will be overridden in the webhook calls.

#### Gmail API consent screen and OAuth2 client secret
This needs to be set up in the Google Cloud Console. The only OAuth 2.0 scope
required by the tool is `https://www.googleapis.com/auth/gmail.modify`.

The basic steps are: 
1. Create a project
1. Enable the Gmail API
1. Configure the Consent Screen
1. Generate an OAuth 2.0 Client ID.
1. Download the json secret file from the Credentials page.

There is a more detailed step-by-step of how to configure the OAuth 2.0 
settings in this tutorial for `gphotos-sync`:
https://www.linuxuprising.com/2019/06/how-to-backup-google-photos-to-your.html

Once the prerequisites are in place, copy `config.example.yaml` to
`config.yaml` and edit it to your requirements.
