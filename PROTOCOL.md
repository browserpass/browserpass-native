# Browserpass Communication Protocol

This document describes the protocol used for communication between the browser extension,
and the native host application.

## Response Types

### OK

Consists solely of an `ok` status, an integer app version, and a `data` field which
may be of any type.

The app version is an integer, calculated by `(MAJOR * 1000000) + (MINOR * 1000) + PATCH`.

```
{
    "status": "ok",
    "version": <int>,
    "data": <any type>
}
```

### Error

Consists solely of an `error` status, a nonzero integer error code, and an optional `params`
object that provides any parameters that should accompany the error.

```
{
    "status": "error",
    "code": <int>,
    "version": <int>,
    "params": {
       "<paramN>": <valueN>
    }
}
```

When the error message is relaying a message from a native system component, then that message
should be supplied as a `message` parameter.

## List of Error Codes

| Code | Description                                                             | Parameters                                                       |
| ---- | ----------------------------------------------------------------------- | ---------------------------------------------------------------- |
| 10   | Unable to parse browser request length                                  | message, error                                                   |
| 11   | Unable to parse browser request                                         | message, error                                                   |
| 12   | Invalid request action                                                  | message, action                                                  |
| 13   | Inaccessible user-configured password store                             | message, action, error, storeId, storePath, storeName            |
| 14   | Inaccessible default password store                                     | message, action, error, storePath                                |
| 15   | Unable to determine the location of the default password store          | message, action, error                                           |
| 16   | Unable to read the default settings of a user-configured password store | message, action, error, storeId, storePath, storeName            |
| 17   | Unable to read the default settings of the default password store       | message, action, error, storePath                                |
| 18   | Unable to list files in a password store                                | message, action, error, storeId, storePath, storeName            |
| 19   | Unable to determine a relative path for a file in a password store      | message, action, error, storeId, storePath, storeName, file      |
| 20   | Invalid password store ID                                               | message, action, storeId                                         |
| 21   | Invalid gpg path                                                        | message, action, error, gpgPath                                  |
| 22   | Unable to detect the location of the gpg binary                         | message, action, error                                           |
| 23   | Invalid password file extension                                         | message, action, file                                            |
| 24   | Unable to decrypt the password file                                     | message, action, error, storeId, storePath, storeName, file      |
| 25   | Unable to list directories in a password store                          | message, action, error, storeId, storePath, storeName            |
| 26   | Unable to determine a relative path for a directory in a password store | message, action, error, storeId, storePath, storeName, directory |
| 27   | The entry contents is missing                                           | message, action                                                  |
| 28   | Unable to determine the recepients for the gpg encryption               | message, action, error, storeId, storePath, storeName, file      |
| 29   | Unable to encrypt the password file                                     | message, action, error, storeId, storePath, storeName, file      |
| 30   | Unable to delete the password file                                      | message, action, error, storeId, storePath, storeName, file      |
| 31   | Unable to determine if directory is empty and can be deleted            | message, action, error, storeId, storePath, storeName, directory |
| 32   | Unable to delete the empty directory                                    | message, action, error, storeId, storePath, storeName, directory |

## Settings

The `settings` object is a key / value map of individual settings. It's provided by the
browser to the native app as part of every request.

Settings are saved in browser local storage. Each top-level setting is saved separately,
JSON-encoded and saved by its key.

Settings may also be supplied via a `.browserpass.json` file in the root of a password store,
and via parameters in individual `*.gpg` files.

Settings are applied using the following priority, highest first:

1.  Configured by the user in specific `*.gpg` files (e.g. autosubmit: true)
1.  Configured by the user in `.browserpass.json` file in specific password stores
1.  Configured by the user via the extension options
1.  Defaults shipped with the browser extension

### Global Settings

| Setting | Description                                          | Default |
| ------- | ---------------------------------------------------- | ------- |
| gpgPath | Optional path to gpg binary                          | `null`  |
| stores  | List of password stores with store-specific settings | `{}`    |

### Store-specific Settings

| Setting | Description                             | Default |
| ------- | --------------------------------------- | ------- |
| id      | Unique store id (same as the store key) | `<key>` |
| name    | Store name                              | `""`    |
| path    | Path to the password store directory    | `""`    |

## Actions

### Configure

Returns a response containing the per-store config. Used to check that the host app
is alive, determine the version at startup, and provide per-store defaults.

#### Request

```
{
    "settings": <settings object>,
    "defaultStoreSettings": <store-specific settings for default store>,
    "action": "configure"
}
```

#### Response

```
{

    "status": "ok",
    "version": <int>,
    "data": {
        "defaultStore": {
            "path": "/path/to/default/store",
            "settings": "<raw contents of $defaultPath/.browserpass.json>",
        },
        “storeSettings”: {
            “storeId”: "<raw contents of storePath/.browserpass.json>"
        }
    }
}
```

### List

Get a list of all `*.gpg` files for each of a provided array of directory paths. The `storeN`
is the ID of a password store, the key in `"settings.stores"` object.

#### Request

```
{
    "settings": <settings object>,
    "action": "list"
}
```

#### Response

```
{
    "status": "ok",
    "version": <int>,
    "data": {
        "files": {
            "storeN": ["<storeNPath/file1.gpg>", "<...>"],
            "storeN+1": ["<storeN+1Path/file1.gpg>", "<...>"]
        }
    }
}
```

### Tree

Get a list of all nested directories for each of a provided array of directory paths. The `storeN`
is the ID of a password store, the key in `"settings.stores"` object.

#### Request

```
{
    "settings": <settings object>,
    "action": "tree"
}
```

#### Response

```
{
    "status": "ok",
    "version": <int>,
    "data": {
        "directories": {
            "storeN": ["<storeNPath/directory1>", "<...>"],
            "storeN+1": ["<storeN+1Path/directory1>", "<...>"]
        }
    }
}
```

### Fetch

Get the decrypted contents of a specific file.

#### Request

```
{
    "settings": <settings object>,
    "action": "fetch",
    "storeId": "<storeId>",
    "file": "relative/path/to/file.gpg"
}
```

#### Response

```
{
    "status": "ok",
    "version": <int>,
    "data": {
        "contents": "<decrypted file contents>"
    }
}
```

### Save

Encrypt the given contents and save to a specific file.

#### Request

```
{
    "settings": <settings object>,
    "action": "save",
    "storeId": "<storeId>",
    "file": "relative/path/to/file.gpg",
    "contents": "<contents to encrypt and save>"
}
```

#### Response

```
{
    "status": "ok",
    "version": <int>
}
```

### Delete

Delete a specific file and empty parent directories caused by the deletion, if any.

#### Request

```
{
    "settings": <settings object>,
    "action": "delete",
    "storeId": "<storeId>",
    "file": "relative/path/to/file.gpg"
}
```

#### Response

```
{
    "status": "ok",
    "version": <int>
}
```

### Echo

Send the `echoResponse` in the request as a response.

#### Request

```
{
    "action": "echo",
    "echoResponse": <anything>
}
```

#### Response

```
<echoResponse>
```
