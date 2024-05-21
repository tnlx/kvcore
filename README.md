# KvCore

KvCore Go API client.

## Usage

Create the API client using the token generated from the portal.

```go
api := kvcore.API{
	Token: "",
}
```

Use the client to perform data retrieval, such as:

### Fetch all contacts by hashtags

Define pagination property, including the size, and operations on each iteration

```go
paginator := kvcore.Paginator{
	PageSize: 100,
	OnPagedSuccess: func(i interface{}) {
		cts := i.([]kvcore.Contact)
		// ...
	},
	OnPagedFailure: func(err error) {},
}
api.ContactsByTags(
    []string{"hashtag_1", "hashtag_2"},
    paginator)
```