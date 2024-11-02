<h1 align="center">🚧 TinyWAF 🚧</h1>

<p align="center">
TinyWAF is a lightweight Web Application Firewall designed for self-hosting.
It can run on the same machine as your web server, or on a dedicated machine.
</p>

> [!WARNING]
> TinyWAF is not ready for production use!


## Why does TinyWAF exist?

TinyWAF was conceived after @nevadascout set up awstats to report traffic stats
from apache log files on a web server on the public internet. The traffic stats
revealed thousands of fishing hits to files or endpoints that didn't exist - bots
looking for things they might try to attack. In particular there were a lot of
requests looking for files related to wordpress scripts. He decided to create a
simple open source WAF that could drop in front of Apache on his server and
shield the server from attacks.

* Most web sites/apps don't have a WAF protecting them.
* However, most web sites/apps could benefit from a WAF (even a simple one).
* Hosted cloud-based WAFs are too expensive for small websites/apps.
* Self-hosting is cool.


## Guiding principles for development

In no particular order:

* TinyWAF is designed to run on the same host machine as the web server, however
it should also be possible to run it on a separate, dedicated machine.
* TinyWAF should be as lightweight and performant as possible, with as few
features as possible.
* TinyWAF should not act as a load balancer or reverse proxy (except the bare
minimum to achieve the requirements of a firewall).
* TinyWAF should be invisible to the user and webserver unless a request/response
is blocked.
* TinyWAF should be simple to configure, and ship with sensible default
security settings (eg. with a set of rules enabled by detault).
* It should be possible to define custom firewall rules and policies for TinyWAF.
* TinyWAF should be thoroughly unit-tested to prevent regressions and issues.


## Development status

### TinyWAF binary

|**Feature**|**Status**|
|:---|:---|
| :bangbang: Request interception and reverse proxy forwarding | :heavy_check_mark: Done |
| Health check endpoint | :heavy_check_mark: Done |
| :bangbang: Define rules in YAML | :heavy_check_mark: Done |
| :bangbang: Request rule evaluation | :heavy_check_mark: Done |
| :bangbang: HTTPS/TLS support | :construction: In progress |
| Rate limiting | :construction: In progress |
| :large_orange_diamond: Performance optimisation | :x: Not started |
| :large_orange_diamond: Metrics/reporting | :x: Not started | <!-- aggregate telemetry for marketing site + reporting for TinyWAF Pro -->
| Websocket forwarding | :x: Not started | <!-- https://github.com/koding/websocketproxy -->
| Response rule evaluation | :x: Not started |
| Custom error pages | :x: Not started |
| AbuseIPDB integration | :x: Not started |
| CrowdSec integration | :x: Not started |

Key
* :bangbang: - Required for production use
* :large_orange_diamond: - Next priority

### TinyWAF default/maintained rulesets

|**Ruleset**|**Status**|
|:---|:---|
| Port of OWASP CRS | :x: Not started |
| Ban AI (block bots scraping data for AI/LLM training) | :x: Not started |
| No Wordpress (for sites not running Wordpress) | :x: Not started |
| No Drupal (for sites not running Drupal) | :x: Not started |
| No Joomla (for sites not running Joomla) | :x: Not started |
| ...others | :x: Not started |


## How are rules evaluated and requests blocked?

<!-- @todo: move this whole section to docs site -->

Rules are defined in YAML and stored in the TinyWAF config directory. TinyWAF
ships with a set of default rules maintaned by the TinyWAF team, but users can
also write their own rules and choose which ones to enable.

Rules are run against requests to prevent SQL injection attacks, etc, but
rules also run against responses to prevent information exposure (eg. leaking
server file paths)

If a request or response matches a defined rule, an action is taken depening on
the rule config. The request may be ignored, warned, ratelimited or blocked.


### Anatomy of a rule

Rules are defined in yaml files. Each ruleset yaml file should start with either
`request-` or `response-` followed by a hypenated rule group name. To disable a
rule file, add `disabled-` at the start of the filename.

Inside each group file is a `rules` array with the following YAML structure:

* `id (string)` - A unique ID for the rule within this group (file)
* `inspect (string|string[])` - Which part of the request/response should this rule apply to
* `whenMethods (string|string[])` - [OPTIONAL] Which request methods should this rule apply to. If not set, applies to all methods
* `fields (string[])` - [OPTIONAL] Which header fields does this request apply to
* `operators` - Which operators to run (contains, exactly, regex + inverse)
* `action ('block'|'ratelimit'|'warn'|'ignore')` - What action to take if a request/response matches this rule

Here's an example rule that will block any request/response with a non-numeric
Content-Length header:

```
rules:
  - id: content-length-not-numeric
    inspect: headers
    fields: "Content-Length"
    operators:
      notregex: ^\d+$
    action: block
```

Here's an example rule that will log a warning about all GET requests to URLs
containing `/signup` or `/login`:

```
rules:
  - id: block-get-signup-login
    whenMethods: get
    inspect: url
    operators:
      contains: "/signup|/login"
    action: warn
```

## Development quickstart

TinyWAF is written in Go.

Download the go runtime, clone the repo to your machine, then open a terminal to
the root of the repo and run `go run ./cmd` to launch TinyWAF.


## TinyWAF installation

Consult [the docs](https://tinywaf.com/docs/) to set up TinyWAF on your server.
