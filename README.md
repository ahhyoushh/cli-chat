
# Chat-CLI

A simple mail-like client and server!




![Logo](https://i.imgur.com/Jf4Pyys.png)


## Initialization

#### Change the username and password to your username and password in config.json

```bash
{
    "username": username,
    "password": password
}
```

#### Change the database string in controller.go package

```bash
const connectionString = your_database_string
```
## API Reference

#### SignUp
```http
  POST /api/signup
```
```bash
{
    "username": username,
    "password": password
}
```


#### Login

```http
  POST /api/login
```

```bash
{
    "username": username,
    "password": password
}
```
#### Get all messages
```http
  POST /api/getall
```

```bash
{
    "username": username,
    "password": password
}
```

#### Get all unread messages
```http
  POST /api/getunread
```
```bash
{
    "username": username,
    "password": password
}
```


#### Send message
```http
  POST /api/send
```

```bash
{
	"sender": sender_username,
	"receiver": receiver_username,
	"message": message
}

```




## Authors

- [@Ayush Bhalerao](https://www.github.com/ahhyoushh)
- [@Tanishk Bansode](https://github.com/TanishkBansode)

