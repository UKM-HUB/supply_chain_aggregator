POST /api/v1/auth/login
POST /api/v1/auth/register
POST /api/v1/auth/refresh
POST /api/v1/auth/logout

Request Login 

{
  "token": "jwt-token",
  "refresh_token": "refresh-token",
  "user": {
    "id": 1,
    "name": "Budi",
    "role": "ADMIN"
  }
}

