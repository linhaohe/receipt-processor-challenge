# receipt-processor-challenge

## Clone the project

```
$ git https://github.com/linhaohe/receipt-processor-challenge.git
$ cd receipt-processor-challenge
```
## Server startup

```
$ cd receipt-processor-challenge
$ go run main.go
```

## Http Request example
### Get
```
http://localhost:8080/receipts/{{id}}/points
```
An Id must exist in local storeage to calculate point, 

#### Success 
```
{ "points": 109}
```

#### Fail
```
{"description" : "No receipt found for that ID."}
```
### Post
```
http://localhost:8080/receipts/process

//request body
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.2a5"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
```
#### on Success
```
{"id": "ce265348-6a8f-4ae3-ae9c-1af6e6d9e126"}
```
#### Fail
```
{ "description" : "The receipt is invalid."}
```
