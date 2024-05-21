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

```go
paginator := kvcore.Paginator{
	PageSize: 100,
	OnPagedSuccess: func(i interface{}) {
		cts := i.([]kvcore.Contact)
		// ...
	},
	OnPagedFailure: func(err error) {},
}
filter := ContactFilter{
	Hashtags: []string{ "cool", "awesome" },
	HashtagsAndOr: "AND",
}
api.ListContacts(filter, paginator)
```