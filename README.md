# Distributed Job Scraper with Actors Model

A small project using the Go Actor Engine. Scrapes job sites for postings related to keywords and location. API has only one endpoint for **POST /find** which initiates the scraping job and the results are sent to the provided email address.

```
make run
```

![jobscr](https://github.com/kerosiinikone/distributed-job-aggregation-service/assets/100020686/cde67059-5860-460a-a315-d6f8cfb32c89)

### TODO

- Job site specific Visitor

Consider these:

- Reuse the same Manager, or spawn new per request ???
- Poison the Manager after job done
