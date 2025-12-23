# NotificationsApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getNotifications**](NotificationsApi.md#getnotifications) | **GET** /notifications | Get user notifications |
| [**getUnreadNotificationsCount**](NotificationsApi.md#getunreadnotificationscount) | **GET** /notifications/unread-count | Get unread notifications count |
| [**markAllNotificationsAsRead**](NotificationsApi.md#markallnotificationsasread) | **PUT** /notifications/read-all | Mark all notifications as read |



## getNotifications

> GetNotifications200Response getNotifications(page, size)

Get user notifications

Retrieves a paginated list of notifications for the currently authenticated user.

### Example

```ts
import {
  Configuration,
  NotificationsApi,
} from '';
import type { GetNotificationsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new NotificationsApi(config);

  const body = {
    // number | Page number (default: 1) (optional)
    page: 56,
    // number | Page size (default: 10) (optional)
    size: 56,
  } satisfies GetNotificationsRequest;

  try {
    const data = await api.getNotifications(body);
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
| **page** | `number` | Page number (default: 1) | [Optional] [Defaults to `undefined`] |
| **size** | `number` | Page size (default: 10) | [Optional] [Defaults to `undefined`] |

### Return type

[**GetNotifications200Response**](GetNotifications200Response.md)

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


## getUnreadNotificationsCount

> { [key: string]: number; } getUnreadNotificationsCount()

Get unread notifications count

Retrieves the number of unread notifications for the currently authenticated user.

### Example

```ts
import {
  Configuration,
  NotificationsApi,
} from '';
import type { GetUnreadNotificationsCountRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new NotificationsApi(config);

  try {
    const data = await api.getUnreadNotificationsCount();
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

**{ [key: string]: number; }**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | count |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## markAllNotificationsAsRead

> ResponsesSimpleMessage markAllNotificationsAsRead()

Mark all notifications as read

Marks all notifications for the currently authenticated user as read.

### Example

```ts
import {
  Configuration,
  NotificationsApi,
} from '';
import type { MarkAllNotificationsAsReadRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new NotificationsApi(config);

  try {
    const data = await api.markAllNotificationsAsRead();
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

