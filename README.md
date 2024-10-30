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

> [!WARNING]
> TinyWAF is not ready for production use!


## Why does TinyWAF exist?

* Most web sites/apps don't have a WAF protecting them.
* However, most web sites/apps could benefit from a WAF (even a simple one).
* Hosted cloud-based WAFs are too expensive for small websites/apps.


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
| Request interception and reverse proxy forwarding | :heavy_check_mark: Done |
| Health check endpoint | :heavy_check_mark: Done |
| Define rules in YAML | :construction: In progress |
| Request rule evaluation | :construction: In progress |
| Response rule evaluation | :construction: In progress |
| Rate limiting | :x: Not started |
| HTTPS/TLS support | :x: Not started |
| Websocket forwarding | :x: Not started |
| Metrics/reporting | :x: Not started |
| Custom error pages | :x: Not started |
| AbuseIPDB integration | :x: Not started |
| CrowdSec integration | :x: Not started |


### TinyWAF default/maintained rulesets

|**Ruleset**|**Status**|
|:---|:---|
| Port of OWASP CRS | :x: Not started |
| No Wordpress (for sites not running Wordpress) | :x: Not started |
| No Drupal (for sites not running Drupal) | :x: Not started |
| No Joomla (for sites not running Joomla) | :x: Not started |
| ...others | :x: Not started |


## How are rules evaluated and requests blocked?

<!-- @todo: move this whole section to docs site -->

TinyWAF works similarly to the OWASP Core RuleSet. Rules are evaluated and an
anomaly score is given to the request based on each evaluated rule. If the anomaly
score is above a defined threshold, then the request is blocked.

Rules are run against inbound requests to prevent SQL injection attacks, etc, but
rules also run against oubound requests to prevent information exposure (eg. leaking
server file paths)

Rules are defined in YAML and stored in the TinyWAF config directory. TinyWAF
ships with a set of default rules maintaned by the TinyWAF team, but users can
also write their own rules and choose which ones to enable.

### Anatomy of a rule

Rules are defined in yaml files. Each ruleset yaml file must start with either
`inbound-` or `-outbound-` followed by a hypenated rule group name.

Inside each group file is a `rules` array with the following YAML structure:

* `id (string)` - A unique ID for the rule within this group (file)
* `priority (int)` - [OPTIONAL] Rules are executed in priority order (0 first)
* `inspect (string|string[])` - Which part of the request/response should this rule apply to
* `fields (string|string[])` - [OPTIONAL] Which fields should this rule apply to
* `operators (string[])` - An array of operators to be run
* `action ('block'|'warn'|'none')` - What action to take if a request/response matches this rule

Here's an example rule that will block any request/response with a non-numeric
Content-Length header:

```
rules:
  - id: content-length-not-numeric
    inspect: HEADERS
    fields: "Content-Length"
    operators:
      - regex: ^\d+$
    severity: 3
    action: block
```

## Development quickstart

TinyWAF is written in Go.

Download the go runtime, clone the repo to your machine, then open a terminal to
the root of the repo and run `go run ./cmd` to launch TinyWAF.


## Hosting quickstart

Consult [the docs](https://tinywaf.com/docs/) to set up TinyWAF on your server.
