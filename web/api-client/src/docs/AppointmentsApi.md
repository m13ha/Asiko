# AppointmentsApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createAppointment**](AppointmentsApi.md#createappointment) | **POST** /appointments | Create a new appointment |
| [**getMyAppointments**](AppointmentsApi.md#getmyappointments) | **GET** /appointments/my | Get appointments created by the user |
| [**getUsersRegisteredForAppointment**](AppointmentsApi.md#getusersregisteredforappointment) | **GET** /appointments/users/{app_code} | Get all bookings for an appointment |



## createAppointment

> EntitiesAppointment createAppointment(appointment)

Create a new appointment

Create a new appointment. Type can be single, group, or party.

### Example

```ts
import {
  Configuration,
  AppointmentsApi,
} from '';
import type { CreateAppointmentRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new AppointmentsApi(config);

  const body = {
    // RequestsAppointmentRequest | Appointment Details
    appointment: ...,
  } satisfies CreateAppointmentRequest;

  try {
    const data = await api.createAppointment(body);
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
| **appointment** | [RequestsAppointmentRequest](RequestsAppointmentRequest.md) | Appointment Details | |

### Return type

[**EntitiesAppointment**](EntitiesAppointment.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **400** | Invalid request payload or validation error |  -  |
| **401** | Authentication required |  -  |
| **500** | Failed to create appointment |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getMyAppointments

> GetMyAppointments200Response getMyAppointments(status, page, size)

Get appointments created by the user

Retrieves a paginated list of appointments created by the currently authenticated user.

### Example

```ts
import {
  Configuration,
  AppointmentsApi,
} from '';
import type { GetMyAppointmentsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new AppointmentsApi(config);

  const body = {
    // Array<string> | Filter by appointment status (pending, ongoing, completed, canceled, expired) (optional)
    status: ...,
    // number | Page number (default: 1) (optional)
    page: 56,
    // number | Page size (default: 10) (optional)
    size: 56,
  } satisfies GetMyAppointmentsRequest;

  try {
    const data = await api.getMyAppointments(body);
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
| **status** | `Array<string>` | Filter by appointment status (pending, ongoing, completed, canceled, expired) | [Optional] |
| **page** | `number` | Page number (default: 1) | [Optional] [Defaults to `undefined`] |
| **size** | `number` | Page size (default: 10) | [Optional] [Defaults to `undefined`] |

### Return type

[**GetMyAppointments200Response**](GetMyAppointments200Response.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Authentication required |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getUsersRegisteredForAppointment

> GetUserRegisteredBookings200Response getUsersRegisteredForAppointment(appCode, page, size)

Get all bookings for an appointment

Retrieves a paginated list of all users/bookings for a specific appointment.

### Example

```ts
import {
  Configuration,
  AppointmentsApi,
} from '';
import type { GetUsersRegisteredForAppointmentRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new AppointmentsApi(config);

  const body = {
    // string | Appointment identifier (app_code)
    appCode: appCode_example,
    // number | Page number (default: 1) (optional)
    page: 56,
    // number | Page size (default: 10) (optional)
    size: 56,
  } satisfies GetUsersRegisteredForAppointmentRequest;

  try {
    const data = await api.getUsersRegisteredForAppointment(body);
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
| **appCode** | `string` | Appointment identifier (app_code) | [Defaults to `undefined`] |
| **page** | `number` | Page number (default: 1) | [Optional] [Defaults to `undefined`] |
| **size** | `number` | Page size (default: 10) | [Optional] [Defaults to `undefined`] |

### Return type

[**GetUserRegisteredBookings200Response**](GetUserRegisteredBookings200Response.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Missing appointment code parameter |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

