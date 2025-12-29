# Notifications Providers

This package defines the email notification contract and the provider selection flow.

## Overview
- `NotificationService` (see `interfaces.go`) is the interface all providers must implement.
- `NewNotificationServiceFromEnv` (see `factory.go`) selects the provider at runtime.
- `AhaSendService` is the default provider.
- `NoopService` is a safe fallback (logs and no-ops) for local/dev or misconfigurations.

## Switching Providers
Set the environment variable:

```
EMAIL_PROVIDER=ahasend
```

Or use the noop provider:

```
EMAIL_PROVIDER=noop
```

If `EMAIL_PROVIDER` is unset, the system defaults to `ahasend`.

## Adding a New Provider
1. Create a new package under `backend/notifications/<provider>`.
2. Implement the `NotificationService` interface from `interfaces.go`.
3. Add a case to `NewNotificationServiceFromEnv` in `factory.go`.
4. Configure any provider-specific environment variables.

## Contract
Each provider must implement:
- `SendBookingConfirmation(booking *entities.Booking) error`
- `SendBookingCancellation(booking *entities.Booking) error`
- `SendBookingRejection(booking *entities.Booking) error`
- `SendVerificationCode(email, code string) error`
- `SendPasswordResetEmail(email, code string) error`

## Notes
- Provider selection happens in `backend/main.go`.
- The AhaSend implementation lives in `ahasend_service.go` and uses the worker publisher in `backend/notifications/ahasend`.
