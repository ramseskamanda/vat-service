# VAT Check Service

This service checks German VAT numbers against the VIES database

## Installation

Clone the repository, then run:

```bash
make install
```

## Usage

Simply run;

```bash
make run MODE={'debug' | 'release' | 'test'}
```

And you'll be able to send requests to the service at the address given in your terminal – by default `http://localhost:8080`

Here's an example request to send VAT number `DE123456789` to the `production` service:

```bash
curl -d '{"vatNumber": "DE123456789"}' -H "Content-Type: application/json" -X POST http://localhost:3000/data
```

This request will return that `DE123456789` is not a valid VAT number:

```json
{
  "vatNumber": "DE123456789",
  "valid": false,
  "message": ""
}
```

## Shape of the data

Input has the following shape:

```go
type CheckVatRequest struct {
    VATNumber string `json:"vatNumber"`
}
```

Output has the following shape:

```go
type CheckVatResponse struct {
    //The VAT number that was evaluated
    VATNumber string `json:"vatNumber"`
    //Whether the VAT number is valid
    Valid     bool   `json:"valid"`
    //Any error message, e.g. "VAT number should be at least 9 characters"
    Message   string `json:"message"`
}
```
