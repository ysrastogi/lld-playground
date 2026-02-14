# Design Patterns

Patterns used in the URL Shortener system.

---

## 1. Layered Architecture

The codebase is organized into distinct horizontal layers with strict dependency direction.

```
Services  →  Repositories  →  Database
    ↓             ↓
  Models        Models
```

| Layer | Package | Responsibility |
|---|---|---|
| **Models** | `src/models` | Domain entities (`URL`, `Metrics`) |
| **Repositories** | `src/repositories` | Data access & persistence |
| **Services** | `services` | Business logic & orchestration |
| **Database** | `src/database` | Connection setup & migrations |
| **Config** | `conf.go` | External configuration loading |

> [!NOTE]
> Each layer only depends on the layer directly below it, keeping coupling low.

---

## 2. Repository Pattern

Data access is abstracted behind the `URLRepository` interface, decoupling business logic from the persistence mechanism (GORM/SQLite).

```go
// Interface — defined in src/repositories/url.go
type URLRepository interface {
    Save(url *models.URL, data map[string]interface{})
    FindByShortURL(shortCode string) (*models.URL, error)
    FindByLongURL(url string) (*models.URL, error)
    Delete(shortCode string)
    // ...
}

// Concrete implementation
type URLRepositoryImpl struct {
    db *gorm.DB
}
```

**Benefits:** The service layer programs against `URLRepository`; the backing store can be swapped (e.g., Postgres, DynamoDB) without touching service code.

---

## 3. Strategy Pattern

Short-code generation is pluggable via the `ShortCodeStrategy` interface.

```go
// services/short_code_generator.go
type ShortCodeStrategy interface {
    generate(shortCode string) string
}

type RandomShortCode struct { ... }   // current implementation
```

**Benefits:** New algorithms (Base62, MD5-hash, counter-based) can be added by implementing the interface — no changes to `URLService`.

---

## 4. Dependency Injection (Constructor Injection)

All dependencies for `URLServiceImpl` are injected through the `NewURLService` constructor, not created internally.

```go
func NewURLService(
    urlRepository repositories.URLRepository,
    cache Cache,
    shortCodeStrategy ShortCodeStrategy,
    rateLimiter RateLimiter,
    analyticsPublisher AnalyticsPublisher,
) URLService { ... }
```

**Benefits:** Makes every dependency explicit, enables easy mocking in tests, and keeps the service free of infrastructure concerns.

---

## 5. Interface Segregation

External infrastructure is accessed through narrow, purpose-specific interfaces defined in the service layer:

| Interface | Methods | Purpose |
|---|---|---|
| `Cache` | `Get`, `Set`, `Delete` | Caching layer (e.g., Redis) |
| `RateLimiter` | `IsAllowed` | Throttle user requests |
| `AnalyticsPublisher` | `Publish` | Emit access metrics |
| `ShortCodeStrategy` | `generate` | Short-code generation |

Each consumer depends only on the slice of behavior it actually needs.

---

## 6. Cache-Aside (Lazy Loading)

`URLServiceImpl.Redirect` implements the cache-aside pattern:

1. **Read from cache** — if hit, return immediately.
2. **On miss** — query the repository, then **populate the cache** before returning.

```go
func (s *URLServiceImpl) Redirect(shortCode string) (string, error) {
    longURL, err := s.cache.Get(shortCode)   // 1. check cache
    if err == nil && longURL != "" {
        return longURL, nil
    }
    url, _ := s.urlRepository.FindByShortURL(shortCode) // 2. fallback to DB
    _ = s.cache.Set(shortCode, url.Url)                 // 3. warm cache
    return url.Url, nil
}
```

Cache is also **proactively set** on `CreateShortURL` and **invalidated** on `Deactivate`.

---

## 7. Soft Delete

`Delete` and `DeleteById` do **not** remove rows; they mark them as stale via a flag update.

```go
func (r *URLRepositoryImpl) Delete(shortCode string) {
    r.db.Model(&url).Update("is_stale", true)
}
```

All read queries filter on `is_stale = false`, making deleted records invisible without losing the data.

---

## 8. Typed Enum (Value Object)

URL lifecycle states are modeled as a custom `Status` type with named constants, avoiding raw strings.

```go
type Status string
const (
    ACTIVE   Status = "ACTIVE"
    INACTIVE Status = "INACTIVE"
    STALE    Status = "STALE"
)
```

---

## 9. Factory Function

Each major component exposes a `New*` constructor that returns the **interface** type, not the concrete struct.

```go
func NewURLRepository(db *gorm.DB) URLRepository       { ... }
func NewURLService(...)             URLService          { ... }
```

Callers never see the implementation, reinforcing the Dependency Inversion Principle.

---

## 10. Externalized Configuration

`conf.go` loads settings from `config.yaml` using Viper and maps them onto typed structs.

```go
type Config struct {
    Redis    Redis    `mapstructure:"redis"`
    Database Database `mapstructure:"database"`
}
```

**Benefits:** Environment-specific values (DB path, Redis URL) live outside code, making deployments flexible without recompilation.

---

## 11. Scheduled Cleanup (Background Job)

`ExpiryScheduler` periodically queries for expired active URLs and soft-deletes them, keeping the main request path free of cleanup logic.

```go
func (e *ExpiryScheduler) schedule() {
    filters := map[string]interface{}{
        "expires_at": time.Now().Add(90 * 24 * time.Hour),
        "status":     "ACTIVE",
    }
    ids, _ := e.UrlRepository.ListURLs(filters)
    for _, id := range ids {
        e.UrlRepository.DeleteById(id.ID)
    }
}
```

---

## Summary

| # | Pattern | Where |
|---|---|---|
| 1 | Layered Architecture | Project structure |
| 2 | Repository | `URLRepository` / `URLRepositoryImpl` |
| 3 | Strategy | `ShortCodeStrategy` / `RandomShortCode` |
| 4 | Dependency Injection | `NewURLService` constructor |
| 5 | Interface Segregation | `Cache`, `RateLimiter`, `AnalyticsPublisher` |
| 6 | Cache-Aside | `Redirect` / `CreateShortURL` |
| 7 | Soft Delete | `Delete`, `DeleteById` |
| 8 | Typed Enum | `Status` constants |
| 9 | Factory Function | `New*` constructors |
| 10 | Externalized Config | `conf.go` + `config.yaml` |
| 11 | Scheduled Cleanup | `ExpiryScheduler` |