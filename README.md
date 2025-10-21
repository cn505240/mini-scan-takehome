# Mini-Scan

Hello!

As you've heard by now, Censys scans the internet at an incredible scale. Processing the results necessitates scaling horizontally across thousands of machines. One key aspect of our architecture is the use of distributed queues to pass data between machines.

---

The `docker-compose.yml` file sets up a toy example of a scanner. It spins up a Google Pub/Sub emulator, creates a topic and subscription, and publishes scan results to the topic. It can be run via `docker compose up`.

Your job is to build the data processing side. It should:

1. Pull scan results from the subscription `scan-sub`.
2. Maintain an up-to-date record of each unique `(ip, port, service)`. This should contain when the service was last scanned and a string containing the service's response.

> **_NOTE_**
> The scanner can publish data in two formats, shown below. In both of the following examples, the service response should be stored as: `"hello world"`.
>
> ```javascript
> {
>   // ...
>   "data_version": 1,
>   "data": {
>     "response_bytes_utf8": "aGVsbG8gd29ybGQ="
>   }
> }
>
> {
>   // ...
>   "data_version": 2,
>   "data": {
>     "response_str": "hello world"
>   }
> }
> ```

Your processing application should be able to be scaled horizontally, but this isn't something you need to actually do. The processing application should use `at-least-once` semantics where ever applicable.

You may write this in any languages you choose, but Go would be preferred.

You may use any data store of your choosing, with `sqlite` being one example. Like our own code, we expect the code structure to make it easy to switch data stores.

Please note that Google Pub/Sub is best effort ordering and we want to keep the latest scan. While the example scanner does not publish scans at a rate where this would be an issue, we expect the application to be able to handle extreme out of orderness. Consider what would happen if the application received a scan that is 24 hours old.

---

Please upload the code to a publicly accessible GitHub, GitLab or other public code repository account. This README file should be updated, briefly documenting your solution. Like our own code, we expect testing instructions: whether it’s an automated test framework, or simple manual steps.

To help set expectations, we believe you should aim to take no more than 4 hours on this task.

We understand that you have other responsibilities, so if you think you’ll need more than 5 business days, just let us know when you expect to send a reply.

Please don't hesitate to ask any follow-up questions for clarification.

---

## Development Setup

### Prerequisites

- [asdf](https://asdf-vm.com/) for tool version management
- Docker and Docker Compose

### Installation

1. Install asdf if not already installed:
   ```bash
   git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.14.0
   echo '. "$HOME/.asdf/asdf.sh"' >> ~/.zshrc
   source ~/.zshrc
   ```

2. Install required tools:
   ```bash
   asdf install
   ```

3. Install Go dependencies:
   ```bash
   go mod download
   ```

4. Start the system:
   ```bash
   make up
   ```

---

## Solution Documentation

### Architecture

**Components:**
- Consumer: Stateless message processor
- Database: PostgreSQL with atomic upserts
- Pub/Sub: Google Pub/Sub emulator

**Data Flow:**
1. Scanner publishes scan results to Pub/Sub topic
2. Consumer pulls messages from `scan-sub` subscription
3. Consumer converts raw data to domain models (handles V1/V2 formats)
4. Consumer applies latest-wins logic with timestamp comparison
5. Consumer upserts to database

### Features

- Horizontal scaling supported with stateless consumers
- At-least-once message processing
- Out-of-order message handling with timestamp-based latest-wins
- Database-level race condition protection
- Support for V1 (base64) and V2 (direct string) data formats
- Repository pattern for data store abstraction

### Concurrency Handling

The system handles concurrent processing of messages for the same `(ip, port, service)` through multiple layers of protection:

**Application-Level Logic:**
- Each consumer checks if incoming scan is newer than existing database record
- Older scans are ignored before database interaction
- Reduces database contention and improves performance

**Database-Level Protection:**
- Atomic upsert operations with conflict resolution
- Timestamp-based WHERE clause ensures only newer data overwrites older data
- Row-level locking prevents race conditions during concurrent writes

### Testing

#### Automated Tests
```bash
make test
```

#### Manual Testing
1. Start the system: `make up`
2. Check message processing: `docker compose logs consumer`
3. Verify database records: `PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d scans -c "SELECT COUNT(*) FROM service_scans;"`
4. View sample data: `PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d scans -c "SELECT * FROM service_scans ORDER BY last_scanned DESC LIMIT 5;"`

#### Out-of-Order Message Testing
1. Start system and let it process messages
2. Stop consumer: `docker compose stop consumer`
3. Wait 5 minutes (messages accumulate in subscription)
4. Restart consumer: `docker compose start consumer`
5. Verify latest scans are preserved

#### Data Format Verification
```sql
SELECT ip, port, service, response, last_scanned 
FROM service_scans 
WHERE response LIKE 'service response:%'
ORDER BY last_scanned DESC;
```

### Error Handling

- Database connection issues: Consumer logs errors and retries
- Malformed messages: Consumer logs parsing errors and nacks messages
- Database write failures: Consumer logs write errors and nacks messages

### Development Commands

```bash
make up          # Start the system
make down        # Stop the system
make test        # Run tests
make lint        # Run linter
make lint-fix    # Fix linting issues
make mocks       # Generate mocks
```
