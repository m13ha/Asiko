# BanListApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**addToBanList**](BanListApi.md#addtobanlist) | **POST** /ban-list | Add email to ban list |
| [**getBanList**](BanListApi.md#getbanlist) | **GET** /ban-list | Get user\&#39;s ban list |
| [**removeFromBanList**](BanListApi.md#removefrombanlist) | **DELETE** /ban-list | Remove email from ban list |



## addToBanList

> EntitiesBanListEntry addToBanList(banRequest)

Add email to ban list

Add an email to the user\&#39;s personal ban list.

### Example

```ts
import {
  Configuration,
  BanListApi,
} from '';
import type { AddToBanListRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new BanListApi(config);

  const body = {
    // RequestsBanRequest | Email to ban
    banRequest: ...,
  } satisfies AddToBanListRequest;

  try {
    const data = await api.addToBanList(body);
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
| **banRequest** | [RequestsBanRequest](RequestsBanRequest.md) | Email to ban | |

### Return type

[**EntitiesBanListEntry**](EntitiesBanListEntry.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **400** | Invalid request |  -  |
| **401** | Unauthorized |  -  |
| **409** | Email already on ban list |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getBanList

> Array&lt;EntitiesBanListEntry&gt; getBanList()

Get user\&#39;s ban list

Get a list of all emails on the user\&#39;s personal ban list.

### Example

```ts
import {
  Configuration,
  BanListApi,
} from '';
import type { GetBanListRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new BanListApi(config);

  try {
    const data = await api.getBanList();
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

[**Array&lt;EntitiesBanListEntry&gt;**](EntitiesBanListEntry.md)

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

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## removeFromBanList

> ResponsesSimpleMessage removeFromBanList(banRequest)

Remove email from ban list

Remove an email from the user\&#39;s personal ban list.

### Example

```ts
import {
  Configuration,
  BanListApi,
} from '';
import type { RemoveFromBanListRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new BanListApi(config);

  const body = {
    // RequestsBanRequest | Email to unban
    banRequest: ...,
  } satisfies RemoveFromBanListRequest;

  try {
    const data = await api.removeFromBanList(body);
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
| **banRequest** | [RequestsBanRequest](RequestsBanRequest.md) | Email to unban | |

### Return type

[**ResponsesSimpleMessage**](ResponsesSimpleMessage.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Invalid request |  -  |
| **401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

