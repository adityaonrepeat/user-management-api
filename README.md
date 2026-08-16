# User Management API

A RESTful API in Go for managing users with a name and date of birth. Age is not stored — it is calculated dynamically from `dob` on every read.

Built with Fiber, PostgreSQL, sqlc, Uber Zap, and go-playground/validator.

## Status

Under development. Setup and run instructions will be added as the implementation lands.

## Tech stack

| | |
|---|---|
| HTTP | [Fiber](https://github.com/gofiber/fiber) v2 |
| Database | PostgreSQL |
| Data access | [sqlc](https://sqlc.dev) (generated, no ORM) |
| Logging | [Uber Zap](https://github.com/uber-go/zap) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |

## Documentation

- [`system-design.md`](system-design.md) — architecture, request flow, and design trade-offs
- [`codebase.md`](codebase.md) — package-by-package guide to the code

## Setup

To be written once the implementation is runnable.

## API

To be written once the endpoints are implemented.
