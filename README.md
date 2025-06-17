# rss-aggregator

saves feed from rss feeders

### docker command for setting up db

```bash
docker run --name rss-db -p 5432:5432 -d -e POSTGRES_PASSWORD=password -e POSTGRES_USER=root -e POSTGRES_DB=rss_db postgres
```
