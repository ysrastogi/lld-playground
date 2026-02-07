# System Requirements

## Overview:
**Design a Rate Limiting system to:**
- Protect backend services from abuse, accidental overload, and malicious traffic
- Enforce fair usage across users, features, and clients
- Provide predictable system behavior under bursty and DDoS scenarios
# Functional Requirements

## User-Level Rate Limiting
**Fixed Quota Limits**
- Limit total number of requests per user
- Example:
    - `10,000 requests / day / user`
- Once exceeded:
    - Reject requests with HTTP `429 Too Many Requests`

**Time Based Window Limit**
Support multiple window types:
- Fixed Window (e.g., 100 req / minute)
- Sliding Window
- Rolling Window

**Burst Handling**
- Allow short-term burst beyond steady rate
- Define:
    - Sustained Rate
    - Burst Capacity
- Example:
    - 100 req/min sustained
    - Burst up to 200 requests within 10 seconds

## Feature-Level Rate Limiting
**Per-Feature Quotas**
- Limit requests for specific APIs or features
- Example:
    - `/generate-report`: 10 req / hour
    - `/login`: 5 req / minute

**Feature + User Combined Limits**
- Enforce stricter limits for expensive endpoints
- Key format example:
    `rate_key = user_id + feature_name`

**Independent Feature Bursting**
- Bursts allowed per feature, not globally
- Prevents one expensive feature from starving others

**Global & System-Level Limits**
- Maximum requests per second across entire system
- Used as a circuit breaker
- Example:
    - If global RPS > threshold → aggressive throttling

## DDoS & Abuse Handling
**IP-Based Throttling**
- Apply stricter limits to:
    - Anonymous users
    - Suspicious IPs
- Progressive penalties:
    - Soft throttle → Hard block → Temporary ban

**Adaptive Throttling**
- Dynamically reduce limits under high load
- Priority order:
    - Anonymous traffic
    - Free-tier users
    - Paid users
- Fail-Safe Behavior
- If rate limiter datastore is unavailable:
    - Fail-Open for critical APIs
    - Fail-Closed for non-critical APIs

# Non-Functional Requirements
## Performance
- Rate limiter check must complete in < 5 ms
- Should handle 100k+ RPS
## Consistency
- Distributed rate limiting must avoid double-counting
- Acceptable error margin must be defined (e.g., ±1%)
## Scalability
- Horizontally scalable
- Stateless API nodes
- Centralized or sharded state store (e.g., Redis)
## Configurability
- Limits configurable without redeploy
- Support dynamic updates per user/feature