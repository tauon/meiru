# Meiru

## About
Meiru (mail) is a simple program to collect and dump email addresses.

## Requirements
* PostgreSQL >= 9.1
* For TLS, a valid cert (key.pem) and private key (key.pem) are necessary.
  * ```
openssl req -nodes -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 7200
```

## Usage
```shell
$ make
$ ./meiru --help
```
## Example
```shell
$ make
$ ./meiru -l
$ curl -X POST http://127.0.0.1:10025 -d "email=hi@hi.com"
$ ./meiru
hi@hi.com
```

## Testing
A database named "meiru_test" will be created for testing.
```shell
$ make test
```

## License
GPLv3: http://www.gnu.org/licenses/gpl-3.0.html
