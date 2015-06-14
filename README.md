Add Product

```bash
curl -X POST http://localhost:31415/v0/products -d '{
  "data" : [
    {
      "type" : "products",
      "attributes": { 
        "title" : "Product A",
        "description" : "This is a product",
        "active" : true,
        "price"  : 1000 
      }
    }
    ]
}'
```

```json
{
    "data": {
        "attributes": {
            "active": true,
            "description": "This is a product",
            "id": "1",
            "price": 1000,
            "title": "Product A"
        },
        "id": "1",
        "relationships": {},
        "type": "products"
    }
}
```

All Products

```bash
curl http://localhost:31415/v0/products
```

```json
{
    "data": [
        {
            "attributes": {
                "active": true,
                "description": "This is a product",
                "id": "2",
                "price": 1000,
                "title": "Product A"
            },
            "id": "2",
            "relationships": {},
            "type": "products"
        },
        {
            "attributes": {
                "active": true,
                "description": "This is a product",
                "id": "1",
                "price": 10,
                "title": ""
            },
            "id": "1",
            "relationships": {},
            "type": "products"
        }
    ]
}
```