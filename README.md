# TinyWAF

TinyWAF is a lightweight Web Application Firewall designed for easy self-hosting
either on a dedicated machine or on the same machine as your web server.

TinyWAF was conceived after @nevadascout set up awstats to report traffic stats
from apache log files on a web server on the public internet. The traffic stats
revealed thousands of fishing hits to files or endpoints that didn't exist - bots
looking for things they might try to attack. In particular there were a lot of
requests looking for files related to wordpress scripts. He decided to create a
simple open source WAF that could drop in front of Apache on his server and
shield the server from attacks.


## Why does TinyWAF exist?

* Most web sites/apps don't have a WAF protecting them.
* However, most web sites/apps could benefit from a WAF (even a simple one).
* Hosted cloud-based WAFs are too expensive for small websites/apps.


## Guiding principles

In no particular order:

* TinyWAF is designed to run on the same host machine as the web server, however
it should also be possible to run it on a separate, dedicated machine.
* TinyWAF should be as lightweight and performant as possible, with as few
features as possible.
* TinyWAF should not act as a load balancer or reverse proxy (except the bare
minimum to achieve the requirements of a firewall).
* TinyWAF should be simple to configure, and ship with sensible default
security settings (eg. with a set of rules enabled by detault).
* It should be possible to define custom firewall rules and policies for TinyWAF.
* TinyWAF should be thoroughly unit-tested to prevent regressions and issues.


## Development quickstart

TinyWAF is written in Go.

Download the go runtime, clone the repo to your machine, then open a terminal to
the root of the repo and run `go run ./cmd` to launch TinyWAF.


## Hosting quickstart

Consult [the docs](https://tinywaf.com/docs/) to set up TinyWAF on your server.
