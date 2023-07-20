# Golang URL Shortener 🐳

**URL shortener implemented in Golang, aiming to provide concise aliases for long URLs**

## Stack 👨‍💻
- **Golang**
- **Fiber (web framework)**
- **Redis (database)**

## Getting Started 🚀

_Create a new file named .env in the ./app directory
(application root) and define the environment variables
used by the application. You can use the provided example
in .env.example file to set up your variables._

> git clone https://github.com/rtsoy/url-shortener

> cd url-shortener

> docker-compose up -d

## Endpoints 🛣️

### 1. Create Shortened URL 📎
### **POST /api/v1**


| Name   | Required | Type   | Description                                                                                      |
|--------|----------|--------|--------------------------------------------------------------------------------------------------|
| url    | required | string | The original long URL to be shortened.                                                           |
| short  | optional | string | Custom short code for the shortened URL. If not provided, a random UUID will be generated.       |
| expiry | optional | int    | Expiry time in hours for the shortened URL. If not provided, the URL will expire after 24 hours. |

**Request**

```azure
{
    "url": "https://github.com/rtsoy/url-shortener",
    "short": "test132456",
    "expiry": 1
}
```

**Response**

```azure
{
    "url": "https://github.com/rtsoy/url-shortener",
    "short": "localhost:3000/test132456",
    "expiry": 1,
    "rate_limit": 9,
    "rate_limit_reset": 30
}
```

### 2. Redirect to the Original URL 🔄
### GET /:url

**When a user visits http://localhost:3000/test132456, 
they will be redirected to https://github.com/rtsoy/url-shortener.** 

## Contributing 🤝
**If you would like to contribute to this project, 
feel free to submit a pull request.**





