# Browserpass Communication Protocol

This document describes the protocol used for communication between the browser extension,
and the native host application.

## Response Types

### OK

Consists solely of an `ok` status, an integer app version, and a `response` field. The response
may be of any type.

The app version is an integer, calculated by `(MAJOR * 1000000) + (MINOR * 1000) + PATCH`.

```
{
    "status": "ok",
    "version": <int>,
    "response": <any type>
}
```

### Error

Consists solely of an `error` status, a nonzero integer error code, and an optional `params`
object that provides any parameters that should accompany the error.

```
{
    "status": "error",
    "code": <int>,
    "params": {
       "<paramN>": <valueN>
    }
}
```

When the error message is relaying a message from a native system component, then that message
should be supplied as a `message` parameter.

## List of Error Codes

| Code | Description | Parameters |
| ---- | ----------- | ---------- |
|      |             |            |

## Settings

The `settings` object is a key / value map of individual settings. It's provided by the
browser to the native app as part of every request.

Settings are saved in browser local storage. Each top-level setting is saved separately,
JSON-encoded and saved by its key.

Settings may also be supplied via a `.browserpass.json` file in the root of a password store,
and via parameters in individual `*.gpg` files.

Settings are applied using the following priority, highest first:

1.  Configured by the user in specific `*.gpg` files (e.g. autosubmit: true)
2.  Configured by the user via the extension options
3.  Configured by the user in each store’s `.browserpass.json` file
4.  Defaults shipped with the browser extension

### Global Settings

| Setting      | Description                                          | Default |
| ------------ | ---------------------------------------------------- | ------- |
| gpgPath      | Optional path to gpg binary                          | `null`  |
| defaultStore | Store-specific settings for default store            | `{}`    |
| stores       | List of password stores with store-specific settings | `{}`    |

### Store-specific Settings

| Setting | Description                          | Default |
| ------- | ------------------------------------ | ------- |
| name    | Store name (same as the store key)   | <key>   |
| path    | Path to the password store directory | `""`    |

## Actions

### Configure

Returns a response containing the per-store config. Used to check that the host app
is alive, determine the version at startup, and provide per-store defaults.

#### Request

```
{
    "settings": <settings object>,
    "action": "configure"
}
```

#### Response

```
{

    "status": "ok",
    "version": <int>,
    "response": {
        "defaultPath": "/path/to/default/store",
        "defaultSettings": "<raw contents of $defaultPath/.browserpass.json>",
        “storeSettings”: {
            “storeName”: "<raw contents of storePath/.browserpass.json>"
        }
    }
}
```

### List

Get a list of all `*.gpg` files for each of a provided array of directory paths. The `storeN`
is the name of a password store, the key in `"settings.stores"` object.

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
    "response": {
        "files": {
            "storeN": ["<storeNPath/file1.gpg>", "<...>"],
            "storeN+1": ["<storeN+1Path/file1.gpg>", "<...>"]
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
    "store": "<storeName>",
    "file": "relative/path/to/file.gpg"
}
```

#### Response

```
{
    "status": "ok",
    "version": <int>,
    "response": {
        "data": "<decrypted file contents>"
    }
}
```
