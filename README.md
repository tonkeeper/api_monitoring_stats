# TON Status

Service for monitoring the infrastructure around the TON blockchain.

## Installation

```shell
git clone https://github.com/tonkeeper/api_monitoring_stats;
cd api_monitoring_stats;
docker compose up -d --build
```

### Optional

If you want to avoid errors, please add the following tokens before running `docker compose up`
```shell
echo  TONCENTER_API_TOKEN=<token> > .env
echo  CHAINSTACK_TOKEN=<token> > .env
```
