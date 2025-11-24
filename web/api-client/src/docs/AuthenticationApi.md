# AuthenticationApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createUser**](AuthenticationApi.md#createuser) | **POST** /users | Create a new user (initiate registration) |
| [**generateDeviceToken**](AuthenticationApi.md#generatedevicetoken) | **POST** /auth/device-token | Generate Device Token |
| [**loginUser**](AuthenticationApi.md#loginuser) | **POST** /login | User Login |
| [**logoutUser**](AuthenticationApi.md#logoutuser) | **POST** /logout | User Logout |
| [**resendVerification**](AuthenticationApi.md#resendverification) | **POST** /auth/resend-verification | Resend verification code |
| [**verifyRegistration**](AuthenticationApi.md#verifyregistration) | **POST** /auth/verify-registration | Verify user registration |



## createUser

> ResponsesSimpleMessage createUser(user)

Create a new user (initiate registration)

Register a new user in the system. This will trigger an email verification.

### Example

```ts
import {
  Configuration,
  AuthenticationApi,
} from '';
import type { CreateUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AuthenticationApi();

  const body = {
    // RequestsUserRequest | User Registration Details
    user: ...,
  } satisfies CreateUserRequest;

  try {
    const data = await api.createUser(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **user** | [RequestsUserRequest](RequestsUserRequest.md) | User Registration Details | |

### Return type

[**ResponsesSimpleMessage**](ResponsesSimpleMessage.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **202** | Accepted |  -  |
| **400** | Invalid request payload or validation error |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## generateDeviceToken

> { [key: string]: string; } generateDeviceToken(device)

Generate Device Token

Generate a short-lived token for a given device ID to be used in booking requests.

### Example

```ts
import {
  Configuration,
  AuthenticationApi,
} from '';
import type { GenerateDeviceTokenRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AuthenticationApi();

  const body = {
    // RequestsDeviceTokenRequest | Device ID
    device: ...,
  } satisfies GenerateDeviceTokenRequest;

  try {
    const data = await api.generateDeviceToken(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **device** | [RequestsDeviceTokenRequest](RequestsDeviceTokenRequest.md) | Device ID | |

### Return type

**{ [key: string]: string; }**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Invalid request body or validation error |  -  |
| **500** | Could not generate token |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## loginUser

> ResponsesLoginResponse loginUser(login)

User Login

Authenticate a user and receive a JWT token.

### Example

```ts
import {
  Configuration,
  AuthenticationApi,
} from '';
import type { LoginUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AuthenticationApi();

  const body = {
    // RequestsLoginRequest | Login Credentials
    login: ...,
  } satisfies LoginUserRequest;

  try {
    const data = await api.loginUser(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **login** | [RequestsLoginRequest](RequestsLoginRequest.md) | Login Credentials | |

### Return type

[**ResponsesLoginResponse**](ResponsesLoginResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **202** | Registration pending verification |  -  |
| **400** | Invalid request body or validation error |  -  |
| **401** | Invalid email or password |  -  |
| **500** | Could not generate token |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## logoutUser

> ResponsesSimpleMessage logoutUser()

User Logout

Invalidate the user\&#39;s session.

### Example

```ts
import {
  Configuration,
  AuthenticationApi,
} from '';
import type { LogoutUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AuthenticationApi();

  try {
    const data = await api.logoutUser();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ResponsesSimpleMessage**](ResponsesSimpleMessage.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## resendVerification

> ResponsesSimpleMessage resendVerification(resend)

Resend verification code

Resend a verification code for a pending user registration.

### Example

```ts
import {
  Configuration,
  AuthenticationApi,
} from '';
import type { ResendVerificationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AuthenticationApi();

  const body = {
    // RequestsResendVerificationRequest | Email to resend verification code to
    resend: ...,
  } satisfies ResendVerificationRequest;

  try {
    const data = await api.resendVerification(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **resend** | [RequestsResendVerificationRequest](RequestsResendVerificationRequest.md) | Email to resend verification code to | |

### Return type

[**ResponsesSimpleMessage**](ResponsesSimpleMessage.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **202** | Accepted |  -  |
| **400** | Invalid request payload |  -  |
| **404** | Pending registration not found |  -  |
| **409** | Account already verified |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## verifyRegistration

> ResponsesLoginResponse verifyRegistration(verification)

Verify user registration

Verify a user\&#39;s email address with a code to complete registration.

### Example

```ts
import {
  Configuration,
  AuthenticationApi,
} from '';
import type { VerifyRegistrationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AuthenticationApi();

  const body = {
    // RequestsVerificationRequest | Email and Verification Code
    verification: ...,
  } satisfies VerifyRegistrationRequest;

  try {
    const data = await api.verifyRegistration(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **verification** | [RequestsVerificationRequest](RequestsVerificationRequest.md) | Email and Verification Code | |

### Return type

[**ResponsesLoginResponse**](ResponsesLoginResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **400** | Invalid request payload or verification error |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

