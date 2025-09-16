# integtrate-pokeapi
# Project Name

Backend service for managing items and orders, including syncing and analytics.

---

## Setup Instructions

1. **Copy environment file**

```bash
cp .env.example .env
```

2. **Build and run with Docker Compose**

```bash
docker-compose up --build
```

---

## API Endpoints

| Endpoint | Method | Description        |
| -------- | ------ | ------------------ |
| `/items` | GET    | Retrieve all items |
| `/sync`  | POST   | Sync data          |

---

## SQL Tasks

### 1. Number of orders and total amount per status in the last 30 days

```sql
SELECT
    status,
    COUNT(*) AS order_count,
    SUM(amount) AS total_amount
FROM orders
WHERE created_at >= NOW() - INTERVAL 30 DAY
GROUP BY status;
```

### 2. Top 5 customers by total spend

```sql
SELECT
    customer_id,
    SUM(amount) AS total_spent
FROM orders
WHERE status = 'PAID'
GROUP BY customer_id
ORDER BY total_spent DESC
LIMIT 5;
```

---

## Notes

* Ensure `.env` is properly configured before starting the service.
* Docker must be installed to run `docker-compose`.

