# Logalert

Lightweight service to tailing logs and send notifications when patterns match.

Features:
- processing multiple log files simultaneously
- сonfigurable interval for checking new records
- [regexp](https://github.com/google/re2/wiki/Syntax) filtering, multiple filters per file
- exceptions for regexp filters
- aggregation of identical records (within one check interval)
- log rotation support
- sending notifications by e-mail
- sending notifications to Telegram

## Dependencies

OS: Linux  
Platform: amd64  
Go version: 1.19+

## Installing
```
make install
```
## Configuring
See the example in the `config.yml.examlple` file

## Running
```
logalert -config=/path/to/config.yml
```
