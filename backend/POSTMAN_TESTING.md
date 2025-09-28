# Postman Testing Guide

This guide explains how to use the comprehensive Postman collection to test all Appointment Master API endpoints automatically.

## Files Included

- `postman_collection_comprehensive.json` - Complete test suite with all API endpoints
- `postman_environment.json` - Environment variables for testing
- `POSTMAN_TESTING.md` - This guide

## Setup Instructions

### 1. Import Collection and Environment

1. Open Postman
2. Click **Import** button
3. Import both files:
   - `postman_collection_comprehensive.json`
   - `postman_environment.json`
4. Select the "Appointment Master Environment" from the environment dropdown

### 2. Start the API Server

```bash
# Using Docker (recommended)
make docker-dev

# Or locally
make dev

# Verify server is running
curl http://localhost:8888/health
```

### 3. Run the Tests

#### Option A: Run Entire Collection
1. Click on "Appointment Master API - Comprehensive Test Suite" collection
2. Click **Run** button
3. Select all requests
4. Click **Run Appointment Master API**

#### Option B: Run Individual Folders
Run folders in this order for best results:
1. **Authentication** - Sets up user and tokens
2. **Appointments** - Creates test appointments
3. **Bookings** - Tests booking functionality
4. **Ban List Management** - Tests ban list features
5. **Analytics** - Tests analytics endpoints
6. **Health Check** - Verifies API health

## Test Flow Overview

### 1. Authentication Flow
```
Register User → Login User → Generate Device Token
```
- Creates a test user account
- Logs in and stores JWT token
- Generates device token for anti-scalping tests

### 2. Appointment Creation
```
Create Single → Create Group → Create Party → Get My Appointments
```
- Creates three different appointment types
- Stores appointment codes for booking tests
- Verifies appointment creation

### 3. Booking Flow
```
Get Available Slots → Book as Guest → Book as Registered User → 
Book Party with Device Token → Get/Update/Cancel Bookings
```
- Tests all booking scenarios
- Includes anti-scalping protection
- Tests booking management operations

### 4. Advanced Features
```
Ban List Management → Analytics → Booking Rejection
```
- Tests ban list CRUD operations
- Verifies analytics data
- Tests appointment owner controls

## Environment Variables

The collection uses these automatically managed variables:

### Static Variables (Pre-configured)
- `base_url` - API base URL (http://localhost:8888)
- `test_user_name` - Test user name
- `test_user_email` - Test user email
- `test_user_password` - Test user password
- `test_user_phone` - Test user phone
- `device_id` - Device identifier for anti-scalping

### Dynamic Variables (Set by tests)
- `auth_token` - JWT authentication token
- `user_id` - Logged in user ID
- `device_token` - Anti-scalping device token
- `single_app_code` - Single appointment code
- `group_app_code` - Group appointment code
- `party_app_code` - Party appointment code
- `guest_booking_code` - Guest booking code
- `user_booking_code` - Registered user booking code
- `party_booking_code` - Party booking code
- `available_slot_*` - Available slot details

## Test Assertions

Each request includes comprehensive test assertions:

### Status Code Checks
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});
```

### Response Data Validation
```javascript
pm.test("Response has required fields", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('id');
    pm.expect(jsonData).to.have.property('name');
});
```

### Environment Variable Updates
```javascript
pm.test("Token stored", function () {
    var jsonData = pm.response.json();
    pm.environment.set('auth_token', jsonData.token);
});
```

## Customization

### Change Test Data
Edit environment variables to use different test data:
```json
{
    "key": "test_user_email",
    "value": "your-test-email@example.com"
}
```

### Add Custom Tests
Add new test assertions to any request:
```javascript
pm.test("Custom validation", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.custom_field).to.equal("expected_value");
});
```

### Environment-Specific URLs
Create multiple environments for different deployment stages:
- Development: `http://localhost:8888`
- Staging: `https://staging-api.example.com`
- Production: `https://api.example.com`

## Troubleshooting

### Common Issues

1. **Server Not Running**
   ```
   Error: connect ECONNREFUSED 127.0.0.1:8888
   ```
   Solution: Start the API server with `make docker-dev` or `make dev`

2. **Authentication Failures**
   ```
   Status code is 401 (Unauthorized)
   ```
   Solution: Run the Authentication folder first to set up tokens

3. **Database Connection Issues**
   ```
   Status code is 500 (Internal Server Error)
   ```
   Solution: Ensure PostgreSQL is running via Docker Compose

4. **Environment Variables Not Set**
   ```
   Variables like {{auth_token}} not resolved
   ```
   Solution: Ensure environment is selected and Authentication tests have run

### Reset Test Data
To start fresh:
1. Clear all environment variables (except static ones)
2. Restart the API server
3. Run the collection from the beginning

## Advanced Usage

### Newman CLI Testing
Run tests from command line using Newman:

```bash
# Install Newman
npm install -g newman

# Run collection
newman run postman_collection_comprehensive.json \
  -e postman_environment.json \
  --reporters cli,json \
  --reporter-json-export results.json
```

### CI/CD Integration
Add to your CI/CD pipeline:

```yaml
# GitHub Actions example
- name: Run API Tests
  run: |
    newman run postman_collection_comprehensive.json \
      -e postman_environment.json \
      --bail \
      --color off
```

### Load Testing
Use Newman with iterations for load testing:

```bash
newman run postman_collection_comprehensive.json \
  -e postman_environment.json \
  -n 10 \
  --delay-request 1000
```

## Test Coverage

This collection covers:

✅ **Authentication & Authorization**
- User registration and login
- JWT token management
- Device token generation

✅ **Appointment Management**
- All appointment types (Single, Group, Party)
- Anti-scalping levels (None, Standard, Strict)
- Appointment CRUD operations

✅ **Booking System**
- Guest and registered user bookings
- Slot availability checking
- Booking management (get, update, cancel, reject)
- Anti-scalping protection

✅ **Advanced Features**
- Ban list management
- Analytics and reporting
- Health checks

✅ **Error Handling**
- Invalid requests
- Authentication failures
- Resource not found scenarios

## Support

For issues with the Postman collection:
1. Check the troubleshooting section above
2. Verify API server is running and healthy
3. Ensure all environment variables are properly set
4. Review test assertions for specific failure details