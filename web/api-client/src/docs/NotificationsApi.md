# NotificationsApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getNotifications**](NotificationsApi.md#getnotifications) | **GET** /notifications | Get user notifications |
| [**markAllNotificationsAsRead**](NotificationsApi.md#markallnotificationsasread) | **PUT** /notifications/read-all | Mark all notifications as read |



## getNotifications

> GetNotifications200Response getNotifications()

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

  try {
    const data = await api.getNotifications();
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

