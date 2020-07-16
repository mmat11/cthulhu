# Cthulhu

![CI](https://github.com/mmat11/cthulhu/workflows/CI/badge.svg)

**Cthulhu** is a Telegram bot fully configurable via a config file.

## Usage

Available commands/functions:

1.  Ban/Unban:

    ```
    /ban
    /unban
    ```

    for these two commands to work, a message with the ban target has to be quoted.

2.  Welcome message, configurable via `welcome_message` key in config.yml.

3.  Community crossposts, configurable via `crosspost_tags` key in config.yml:

    ```
    #tag message
    ```

    Cthulhu will forward _message_ in all the groups in which `#tag` is part of `crosspost_tags`.

4.  Admin/mod permissions with commands, configurable per group.

5.  Broadcast:

    ```
    /broadcast message
    ```

    _message_ will be sent in all groups except the current one.

6.  Prometheus metrics exporter: https://cthulhu.live/d/vjZP1InMk/cthulhu

## Contributing

TBD

## License

GNU GPL v3
