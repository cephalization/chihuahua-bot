# chihuahua-bot

Slack bot

## Libraries
---

[Slack](https://github.com/nlopes/slack) - Slack Bot API client

[ansi](https://github.com/mgutz/ansi) - Pretty terminal colors

## Development
---



- install `go 1.13`
- install `docker`
  - If using windows 10 home, install [docker-toolbox](https://docs.docker.com/toolbox/toolbox_install_windows/) instead
  - If you use toolbox you'll have to use their "quickstart-terminal" unfortunately
  - Should probably just do all of this in windows subsystem for linux (wsl)
- clone this repo
- create `.env` file using `.env.sample` (in the repo) as reference
  - token obtained from https://api.slack.com/legacy/custom-integrations/legacy-tokens once you have permissions to develop on the bot

```
token=slackbot_api_token
```

- in the repo, run `docker-compose up -d`
  - if you leave out the `-d` you can see all of the logging and messages
  - `ctrl-c` will kill the bot like this
- to shut down bot, run `docker-compose down`

If docker isn't working or you don't want to use docker you will have to load the environment variables some other way

`token=slackbot_api_token <other env vars from .env> ./chihuahua-bot`

### Useful development commands

This section is mostly dedicated to docker since it makes development so much easier.

- `docker-compose up` will run the last built version of the bot and print output of all containers
  - `docker-compose up --build` will rebuild all the containers, useful when files change
  - `docker-compose up -d` will run the containers in the background
  - I usually run `docker-compose up --build -d`
- `docker-compose down` will tear down all containers
- `docker image prune` will delete stopped containers to get some disk space back
- `docker logs <container>` will show you the output of a container running in the background
  - `docker logs chibot` will show you the bot's output
  - `docker logs mongodb` will show you the db output
- `docker exec -it <container> /bin/bash` will put you in a terminal inside a container
  - `docker exec -it mongodb /bin/bash` will be useful to query the mongodb manually. 
  - You could also download a mongo gui (mongodb compass) and point it to the default mongo ports it provides while the container is running.
