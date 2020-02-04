# chihuahua-bot

Slack bot

## Libraries

[Slacker](https://github.com/shomali11/slacker) - Slack Bot Framework

[ansi](https://github.com/mgutz/ansi) - Pretty terminal colors

## Development

- install `go 1.13`
- install `docker`
  - If using windows 10 home, install [docker-toolbox](https://docs.docker.com/toolbox/toolbox_install_windows/) instead
  - If you use toolbox you'll have to use their "quickstart-terminal" unfortunately
  - Should probably just do all of this in windows subsystem for linux (wsl)
- clone this repo
- create `.env` file with the following contents
  - token obtained from https://api.slack.com/legacy/custom-integrations/legacy-tokens once you have permissions to develop on the bot

```
token=slackbot_api_token
```

- in the repo, run `docker-compose up -d`
- to shut down bot, run `docker-compose down`
