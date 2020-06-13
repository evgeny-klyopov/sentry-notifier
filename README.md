# sentry-notifier

### Install
```bash
mkdir ~/sentry-notifier
cd ~/sentry-notifier
wget https://github.com/evgeny-klyopov/sentry-notifier/releases/download/v0.0.3-alpha/sentry-notifier.linux-amd64.tar.gz
tar -xvf sentry-notifier.linux-amd64.tar.gz
```
### Add config
```bash
wget https://raw.githubusercontent.com/evgeny-klyopov/sentry-notifier/master/example/config.json
```
### Set setting for sentry and telegram
```json
{
  "organization": [
    {
      "name": "sentry-org",
      "token": "sentry-token",
      "projects": [{"name": "sentry-project"}]
    }
  ],
  "default": {
    "sentry": {
      "issue_filter": {
        "query": "is:unresolved",
        "stats_period": "24h"
      },
      "wait_time": 300
    },
    "notifications": {
      "telegram": [{
        "use_proxy": true,
        "proxy": "user:password@ip:port",
        "chat_id": 11111,
        "token": "your_bot_token"
      }]
    }
  }
}
```

### Add service
```go
cat /etc/systemd/system/sentry-notifier.service
[Unit]
Description=sentry-notifier service
ConditionPathExists=/home/YOUR_USER_NAME/sentry-notifier/
After=network.target
 
[Service]
Type=simple
User=YOUR_USER_NAME
Group=YOUR_USER_GROUP

Restart=on-failure
RestartSec=10
startLimitIntervalSec=60

WorkingDirectory=/home/YOUR_USER_NAME/sentry-notifier/
ExecStart=/home/sentry/YOUR_USER_NAME-notifier/sentry-notifier --name=sentry-notifier

SyslogIdentifier=sentry-notifier
 
[Install]
WantedBy=multi-user.target
```
### Enable service
```bash
systemctl enable sentry-notifier
```