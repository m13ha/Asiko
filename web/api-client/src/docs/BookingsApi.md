# BookingsApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**bookGuestAppointment**](BookingsApi.md#bookguestappointment) | **POST** /appointments/book | Book an appointment (Guest) |
| [**bookRegisteredUserAppointment**](BookingsApi.md#bookregistereduserappointment) | **POST** /appointments/book/registered | Book an appointment (Registered User) |
| [**cancelBookingByCode**](BookingsApi.md#cancelbookingbycode) | **DELETE** /bookings/{booking_code} | Cancel a booking |
| [**getAvailableSlots**](BookingsApi.md#getavailableslots) | **GET** /appointments/slots/{app_code} | Get available slots for an appointment |
| [**getAvailableSlotsByDay**](BookingsApi.md#getavailableslotsbyday) | **GET** /appointments/slots/{app_code}/by-day | Get available slots for a specific day |
| [**getBookingByCode**](BookingsApi.md#getbookingbycode) | **GET** /bookings/{booking_code} | Get booking by code |
| [**getUserRegisteredBookings**](BookingsApi.md#getuserregisteredbookings) | **GET** /appointments/registered | Get user\&#39;s registered bookings |
| [**rejectBookingByCode**](BookingsApi.md#rejectbookingbycode) | **POST** /bookings/{booking_code}/reject | Reject a booking |
| [**updateBookingByCode**](BookingsApi.md#updatebookingbycode) | **PUT** /bookings/{booking_code} | Update/Reschedule a booking |



## bookGuestAppointment

> EntitiesBooking bookGuestAppointment(booking)

Book an appointment (Guest)

Creates a booking for an appointment as a guest user. Name and email/phone are required.

### Example

```ts
import {
  Configuration,
  BookingsApi,
} from '';
import type { BookGuestAppointmentRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BookingsApi();

  const body = {
    // RequestsBookingRequest | Booking Details
    booking: ...,
  } satisfies BookGuestAppointmentRequest;

  try {
    const data = await api.bookGuestAppointment(body);
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
| **booking** | [RequestsBookingRequest](RequestsBookingRequest.md) | Booking Details | |

### Return type

[**EntitiesBooking**](EntitiesBooking.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **400** | Invalid request payload or validation error |  -  |
| **409** | Slot unavailable or capacity exceeded |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## bookRegisteredUserAppointment

> EntitiesBooking bookRegisteredUserAppointment(booking)

Book an appointment (Registered User)

Creates a booking for an appointment as a registered user.

### Example

```ts
import {
  Configuration,
  BookingsApi,
} from '';
import type { BookRegisteredUserAppointmentRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new BookingsApi(config);

  const body = {
    // RequestsBookingRequest | Booking Details
    booking: ...,
  } satisfies BookRegisteredUserAppointmentRequest;

  try {
    const data = await api.bookRegisteredUserAppointment(body);
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
| **booking** | [RequestsBookingRequest](RequestsBookingRequest.md) | Booking Details | |

### Return type

[**EntitiesBooking**](EntitiesBooking.md)

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
| **401** | Unauthorized |  -  |
| **409** | Slot unavailable or capacity exceeded |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## cancelBookingByCode

> EntitiesBooking cancelBookingByCode(bookingCode)

Cancel a booking

Cancels a booking by its unique booking_code. This is a soft delete.

### Example

```ts
import {
  Configuration,
  BookingsApi,
} from '';
import type { CancelBookingByCodeRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BookingsApi();

  const body = {
    // string | Unique Booking Code
    bookingCode: bookingCode_example,
  } satisfies CancelBookingByCodeRequest;

  try {
    const data = await api.cancelBookingByCode(body);
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
| **bookingCode** | `string` | Unique Booking Code | [Defaults to `undefined`] |

### Return type

[**EntitiesBooking**](EntitiesBooking.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Error while cancelling booking |  -  |
| **404** | Booking not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getAvailableSlots

> GetUserRegisteredBookings200Response getAvailableSlots(appCode)

Get available slots for an appointment

Retrieves a paginated list of all available booking slots for a given appointment.

### Example

```ts
import {
  Configuration,
  BookingsApi,
} from '';
import type { GetAvailableSlotsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BookingsApi();

  const body = {
    // string | Appointment identifier (app_code)
    appCode: appCode_example,
  } satisfies GetAvailableSlotsRequest;

  try {
    const data = await api.getAvailableSlots(body);
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

### Return type

[**GetUserRegisteredBookings200Response**](GetUserRegisteredBookings200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Missing appointment code parameter |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getAvailableSlotsByDay

> GetUserRegisteredBookings200Response getAvailableSlotsByDay(appCode, date)

Get available slots for a specific day

Retrieves a paginated list of available slots for an appointment on a specific day.

### Example

```ts
import {
  Configuration,
  BookingsApi,
} from '';
import type { GetAvailableSlotsByDayRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BookingsApi();

  const body = {
    // string | Appointment identifier (app_code)
    appCode: appCode_example,
    // string | Date in YYYY-MM-DD format
    date: date_example,
  } satisfies GetAvailableSlotsByDayRequest;

  try {
    const data = await api.getAvailableSlotsByDay(body);
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
| **date** | `string` | Date in YYYY-MM-DD format | [Defaults to `undefined`] |

### Return type

[**GetUserRegisteredBookings200Response**](GetUserRegisteredBookings200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Missing or invalid parameters |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getBookingByCode

> EntitiesBooking getBookingByCode(bookingCode)

Get booking by code

Retrieves booking details by its unique booking_code.

### Example

```ts
import {
  Configuration,
  BookingsApi,
} from '';
import type { GetBookingByCodeRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BookingsApi();

  const body = {
    // string | Unique Booking Code
    bookingCode: bookingCode_example,
  } satisfies GetBookingByCodeRequest;

  try {
    const data = await api.getBookingByCode(body);
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
| **bookingCode** | `string` | Unique Booking Code | [Defaults to `undefined`] |

### Return type

[**EntitiesBooking**](EntitiesBooking.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Missing booking_code parameter |  -  |
| **404** | Booking not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getUserRegisteredBookings

> GetUserRegisteredBookings200Response getUserRegisteredBookings()

Get user\&#39;s registered bookings

Retrieves a paginated list of all bookings made by the currently authenticated user.

### Example

```ts
import {
  Configuration,
  BookingsApi,
} from '';
import type { GetUserRegisteredBookingsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new BookingsApi(config);

  try {
    const data = await api.getUserRegisteredBookings();
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
| **401** | Unauthorized |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## rejectBookingByCode

> EntitiesBooking rejectBookingByCode(bookingCode)

Reject a booking

Rejects a booking by its unique booking_code. This is a soft delete.

### Example

```ts
import {
  Configuration,
  BookingsApi,
} from '';
import type { RejectBookingByCodeRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new BookingsApi(config);

  const body = {
    // string | Unique Booking Code
    bookingCode: bookingCode_example,
  } satisfies RejectBookingByCodeRequest;

  try {
    const data = await api.rejectBookingByCode(body);
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
| **bookingCode** | `string` | Unique Booking Code | [Defaults to `undefined`] |

### Return type

[**EntitiesBooking**](EntitiesBooking.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Error while rejecting booking |  -  |
| **401** | Unauthorized |  -  |
| **404** | Booking not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateBookingByCode

> EntitiesBooking updateBookingByCode(bookingCode, booking)

Update/Reschedule a booking

Updates a booking by its unique booking_code. Can be used to reschedule.

### Example

```ts
import {
  Configuration,
  BookingsApi,
} from '';
import type { UpdateBookingByCodeRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BookingsApi();

  const body = {
    // string | Unique Booking Code
    bookingCode: bookingCode_example,
    // RequestsBookingRequest | New Booking Details
    booking: ...,
  } satisfies UpdateBookingByCodeRequest;

  try {
    const data = await api.updateBookingByCode(body);
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
| **bookingCode** | `string` | Unique Booking Code | [Defaults to `undefined`] |
| **booking** | [RequestsBookingRequest](RequestsBookingRequest.md) | New Booking Details | |

### Return type

[**EntitiesBooking**](EntitiesBooking.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Invalid request, validation error, or slot not available |  -  |
| **404** | Booking not found |  -  |
| **409** | Requested slot not available or capacity exceeded |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

