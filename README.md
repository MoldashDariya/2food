# 2food

2food is a mobile app that helps restaurants sell their surplus food at a discount instead of throwing it away. Customers browse offers from nearby restaurants, reserve and pay for them in the app, and pick the order up at a set time.

**Team:** a final team project for a Go course, built by Dariya Moldash, Alnura Imirova and Saniya Zhaksylykova and moved here from the team's GitLab. My part (Dariya): the catalog (cities, cuisine types, product categories), reviews with photos, restaurant replies and moderation, the user profile and photo uploads to S3. The rest of the project was built by Alnura and Saniya.

The repository is a monorepo with two parts:

| Folder | What it is | Stack |
| --- | --- | --- |
| [`frontend/`](./frontend) | Mobile app for customers and restaurant staff | Expo / React Native, TypeScript |
| [`backend/`](./backend) | REST API, background workers and database migrations | Go, PostgreSQL, Redis |

## How it works

1. A restaurant publishes a product with an original price, a discounted price, the amount left and a pickup time window.
2. A customer picks an offer and checks out. The backend creates the order (`PENDING`) together with a Stripe payment intent that uses manual capture.
3. The customer authorizes the card in the Stripe payment sheet. Stripe calls the backend webhook, the backend locks the products, checks the stock, takes it and marks the order `PAID`. If somebody else took the last item in the meantime, the order is cancelled and the card hold is released.
4. The money is captured in a background job, and the customer gets push notifications on every status change.
5. The restaurant marks the order `READY`, the customer shows the pickup code, and the staff marks the order `COLLECTED`.

Orders that stay unpaid are cancelled automatically (after one hour by default). A delayed audit job re-checks the payment with Stripe and fixes the order if the database and the payment provider disagree.

Order statuses: `PENDING`, `PAID`, `ACCEPTED` (optional step), `READY`, `COLLECTED`, `CANCELLED`, `UNPAID`.

## Features

**Customers**

- Sign in with a phone number (one-time code by SMS) or with Google
- Browse restaurants and offers, search, filter by category and cuisine
- Cart with stock limits, one restaurant per order
- Card payments (including Google Pay) through Stripe
- Order history, order details with pickup code, payment history
- Reviews with photos and restaurant replies
- Push notifications about order status
- Interface in Russian, English, Kazakh, Ukrainian and Chinese

**Restaurant staff (partners)**

- Manage products, prices, stock and pickup times
- Live list of active orders, mark orders as ready or collected, verify the pickup code
- Manage staff (owner and manager roles), reply to reviews

**Administrators**

- Create and delete restaurants, assign owners
- Manage catalog data (cities, cuisine types, product categories)
- Moderate reviews

## Tech stack

**Backend**

- Go 1.25, [chi](https://github.com/go-chi/chi) router, `log/slog`
- PostgreSQL 16 with [golang-migrate](https://github.com/golang-migrate/migrate) migrations
- Redis 7 for OTP codes, OAuth sessions and the job queue
- [asynq](https://github.com/hibiken/asynq) for background jobs and cron (payment confirmation, refunds, audits, push, cleanup)
- Stripe for payments, Twilio for SMS, SMTP for email codes, Firebase Cloud Messaging for push
- S3-compatible storage (MinIO locally) with presigned uploads
- JWT access and refresh tokens, role-based access (customer, partner, admin)
- Prometheus metrics, Swagger / Scalar API documentation

**Frontend**

- Expo SDK 54, React Native 0.81, Expo Router (file-based routing)
- TypeScript, Zustand for state, Axios for the API layer
- i18next for translations, Stripe React Native, Google Sign-In, Expo Notifications
- Tokens are stored in the device's secure storage

## Repository structure

```
.
├── backend
│   ├── cmd/2food          # API server entry point
│   ├── cmd/migrator       # migration runner
│   ├── config             # env-based configuration
│   ├── internal
│   │   ├── app            # wiring of all components and the router
│   │   ├── auth, user, restaurant, product, order, review, catalog, upload
│   │   │                  # one package per feature: service + HTTP transport
│   │   ├── domain         # entities, statuses and domain errors
│   │   ├── infra          # postgres, redis queue, stripe, twilio, smtp, fcm, s3, jwt
│   │   ├── worker         # background job handlers
│   │   └── pkg            # middleware, validation, translations, helpers
│   ├── migrations         # SQL migrations
│   └── docs               # generated Swagger files
└── frontend
    ├── app                # screens (Expo Router): tabs, admin, checkout, login, ...
    ├── components         # reusable UI components
    ├── services           # API calls
    ├── stores             # Zustand stores (auth, cart, admin, notifications, language)
    ├── hooks, utils, constants, types
    └── i18n               # translations (ru, en, kk, uk, zh)
```

## Getting started

### Backend

Requirements: Go 1.25+, Docker.

1. Create the config and start the infrastructure:

   ```bash
   cd backend
   cp .env.example .env
   ```

   `docker-compose.yml` also reads `MINIO_ROOT_USER`, `MINIO_ROOT_PASSWORD`, `MINIO_PORT` and `MINIO_CONSOLE_PORT`, which are not in `.env.example`. Add them to `.env`, for example:

   ```
   MINIO_ROOT_USER=accesskey
   MINIO_ROOT_PASSWORD=secretkey
   MINIO_PORT=9000
   MINIO_CONSOLE_PORT=9001
   ```

   Then start PostgreSQL, Redis and MinIO:

   ```bash
   docker compose up -d postgres redis rustfs
   ```

2. Apply the migrations:

   ```bash
   go run ./cmd/migrator -up -user postgres -password postgres -dbname 2food
   ```

3. Run the server:

   ```bash
   go run ./cmd/2food
   ```

   The API is available at `http://localhost:8080/api/v1`.

Notes:

- On startup the server initialises the email (SMTP), storage (S3) and push (FCM) providers, so `FCM_CREDENTIALS_B64` has to contain real base64-encoded Firebase service account credentials, otherwise the server will not start. Stripe, Twilio and Google values are only used when the matching endpoints are called.
- Create the `2food` bucket in MinIO (console on `http://localhost:9001`) before testing uploads.
- Run the tests with `go test ./...`.

### API documentation

With the server running, open `http://localhost:8080/docs` for the interactive reference. The Swagger files in `backend/docs` are generated from the annotations in the code, so do not edit them by hand. To regenerate them:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/2food/main.go
```

### Frontend

Requirements: Node.js, [pnpm](https://pnpm.io), Xcode (iOS) or Android Studio (Android).

```bash
cd frontend
pnpm install
cp .env.example .env   # set EXPO_PUBLIC_API_BASE_URL to your backend
pnpm expo start
```

The app uses native modules (Stripe, Google Sign-In, push notifications), so it does not run fully in Expo Go. Use a development build:

```bash
pnpm ios       # or: pnpm android
```

Cloud builds are configured in `eas.json` (`development`, `preview` and `production` profiles). Push notifications need the Firebase config files `GoogleService-Info.plist` (iOS) and `google-services.json` (Android). They are not stored in this repository: download them from your Firebase project (Project settings, Your apps) and put them into the `frontend` folder, `app.json` refers to them and native builds fail without them.

### Environment variables

| Where | File | What is inside |
| --- | --- | --- |
| Backend | `backend/.env.example` | server, database, Redis, JWT, OTP, Twilio, SMTP, Google, Stripe, S3, FCM and worker settings |
| Frontend | `frontend/.env.example` | API URL, Google client IDs, Stripe publishable key |

Never commit real `.env` files. Both folders already ignore them.

## Deployment

The backend has a multi-stage `Dockerfile` that builds two images, the API server and the migrator (`--target backend` and `--target migrator`, with `SERVICE_PATH` set to `./cmd/2food/main.go` or `./cmd/migrator/migrator.go`). `backend/.gitlab-ci.yml` is the CI pipeline from the original GitLab setup (tests, image build, deploy); it does not run on GitHub.
