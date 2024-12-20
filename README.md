# receipt-processor-challenge
fetch take home assessment

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

#### Sucess 
```
{status:200, "points": 109}
```

#### Fail
```
{status:404, "description" : "No receipt found for that ID."}
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
#### on sucess
```
{status:201, "id": "ce265348-6a8f-4ae3-ae9c-1af6e6d9e126"}
```
#### Fail
```
{status:400, "description" : "The receipt is invalid."}
```
