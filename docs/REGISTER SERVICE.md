Endpoint
POST /api/v1/register
POST /api/v1/register/verify

Request

{
  "name": "UMKM Maju",
  "phone": "628123456789",
  "password": "123456"
}

Flow
Register
   ↓
Send OTP
   ↓
Verify OTP
   ↓
Create User