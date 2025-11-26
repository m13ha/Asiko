```mermaid
graph TD
    subgraph Frontend (Web App)
        FE_REG[Register Page] --> FE_LOGIN
        FE_LOGIN[Login Page] --> FE_AUTH_CONTEXT(Auth Context: localStorage)
        FE_AUTH_CONTEXT --> FE_API_CLIENT(API Client with Auth Middleware)
    end

    subgraph Backend (Go API)
        BE_REG(POST /api/auth/register)
        BE_VERIFY(POST /api/auth/verify-registration)
        BE_LOGIN(POST /api/login)
        BE_REFRESH(POST /api/auth/refresh)
        BE_AUTH_MIDDLEWARE(Auth Middleware: JWT Validation)
    end

    subgraph Database
        DB_PENDING(Pending Users Table)
        DB_USERS(Users Table)
    end

    %% Frontend to Backend actions
    FE_REG -- 1. User provides name, email, password --> BE_REG
    FE_LOGIN -- 1. User provides email, password --> BE_LOGIN
    FE_API_CLIENT -- 1. API Request (with Access Token) --> BE_AUTH_MIDDLEWARE

    %% Backend internal actions
    BE_REG -- 2. Store Hashed Password & Email Verification --> DB_PENDING
    DB_PENDING -- 3. Email Verification Code --> FE_VERIFY[Verify Page]

    %% User verifies email
    FE_VERIFY -- 4. User provides email, code --> BE_VERIFY
    BE_VERIFY -- 5. Creates User, Generates Tokens --> FE_AUTH_CONTEXT

    %% Login flow details
    BE_LOGIN -- 2. Authenticate User (check hashed password) --> DB_USERS
    DB_USERS -- 3. User Data --> BE_LOGIN
    BE_LOGIN -- 4. Generate Access/Refresh Tokens --> FE_AUTH_CONTEXT

    %% Token Refresh flow
    BE_AUTH_MIDDLEWARE -- 2. If Token Expired/Near Expiry: Request Refresh --> BE_REFRESH
    BE_REFRESH -- 3. Validate Refresh Token, Generate New Tokens --> FE_AUTH_CONTEXT

    %% Auth Middleware to Protected Resources
    BE_AUTH_MIDDLEWARE -- If Authorized --> PROTECTED_RESOURCE[Protected API Endpoint]

    %% Guest booking
    FE_API_CLIENT -- Request Device Token --> BE_DEVICE_TOKEN(POST /api/auth/device-token)
    BE_DEVICE_TOKEN -- Returns Device Token --> FE_API_CLIENT
    FE_API_CLIENT -- Use Device Token for Guest Booking --> BE_BOOKING_GUEST[Guest Booking API]

    style FE_REG fill:#fff,stroke:#333,stroke-width:2px;
    style FE_LOGIN fill:#fff,stroke:#333,stroke-width:2px;
    style FE_AUTH_CONTEXT fill:#f9f,stroke:#333,stroke-width:2px;
    style FE_API_CLIENT fill:#fff,stroke:#333,stroke-width:2px;
    style FE_VERIFY fill:#fff,stroke:#333,stroke-width:2px;

    style BE_REG fill:#add8e6,stroke:#333,stroke-width:2px;
    style BE_VERIFY fill:#add8e6,stroke:#333,stroke-width:2px;
    style BE_LOGIN fill:#add8e6,stroke:#333,stroke-width:2px;
    style BE_REFRESH fill:#add8e6,stroke:#333,stroke-width:2px;
    style BE_AUTH_MIDDLEWARE fill:#add8e6,stroke:#333,stroke-width:2px;
    style BE_DEVICE_TOKEN fill:#add8e6,stroke:#333,stroke-width:2px;
    style BE_BOOKING_GUEST fill:#add8e6,stroke:#333,stroke-width:2px;

    style DB_PENDING fill:#e0b2b2,stroke:#333,stroke-width:2px;
    style DB_USERS fill:#e0b2b2,stroke:#333,stroke-width:2px;

    style PROTECTED_RESOURCE fill:#c7e6ad,stroke:#333,stroke-width:2px;
```