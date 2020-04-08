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

2.  Welcome message (optional), configurable via the `welcome_message` key in the config.yml.

3.  Community crossposts (optional), configurable via the `crosspost_tags` key in the config.yml:

    ```
    #tag
    ```

    Cthulhu will forward the message in all the groups in which `#tag` is in `crosspost_tags`.

4.  Admin/mod permissions with commands, configurable per group in the config.yml.

## Contributing

TBD

## License

GNU GPL v3
