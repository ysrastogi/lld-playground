# Components

The system is divided into three main layers: Models, Services, and Repositories.

## 1. Domain Models (`src/models`)

The business entities that hold state.

### URL
- **URL**: Represents a shortened URL.
  - Attributes: `ID`, `ShortURL`, `LongURL`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, `IsStale`.

## 2. Services (`src/services`)

Contain the business logic.

### URLService
- **Responsibilities**:
  - Shortening a URL.
  - Redirecting a shortened URL to the original URL.
  - Deleting a shortened URL.
  - Updating a shortened URL.
  - Staling a shortened URL.

## 3. Repositories (`src/repositories`)
- **URLRepository**: Store and retrieve URL details.
