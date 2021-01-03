# heig-MAC/project


![.github/workflows/format.yml](https://github.com/heig-MAC/project/workflows/.github/workflows/format.yml/badge.svg)
![.github/workflows/tests.yml](https://github.com/heig-MAC/project/workflows/.github/workflows/tests.yml/badge.svg)

An Telegram bot for everything related to climbing. This is a semester project done at HEIG-VD.

## Team

| Name                                   |                                  |
|----------------------------------------|----------------------------------|
| Matthieu Burguburu                     | matthieu.burguburu@heig-vd.ch    |
| Alexandre Piveteau                     | alexandre.piveteau@heig-vd.ch    |
| Guy-Laurent Subri                      | guy-laurent.subri@heig-vd.ch     |

## Setting the project up

To run the project locally, you'll need:

- [Docker compose](https://docs.docker.com/compose/); and
- a Telegram [bot API token](https://core.telegram.org/bots/api).

You'll have to create a file named `.env` in `./docker/topologies/dev` :

```sh
> cat ./docker/topologies/dev/.env

TELEGRAM_BOT_DEBUG=false
TELEGRAM_BOT_TOKEN=123_YOUR_TELEGRAM_API_TOKEN
```

Running the bot then done as follows:

```sh
> ./run-compose.sh
```

The bot will remain active until it receives a SIGTERM.
