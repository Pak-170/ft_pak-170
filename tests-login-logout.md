# 1. Register
curl -i -X POST localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"juana","password":"12345678","alias":"juanita"}'

# 2. Login (-c saves cookies at cookies.txt)
curl -i -c cookies.txt -X POST localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"juana","password":"12345678"}'

# 3. Protected path (-b sends the cookie)
curl -i -b cookies.txt localhost:8080/api/me

# 4. Logout
curl -i -b cookies.txt -X POST localhost:8080/api/logout

# 5. Again /api/me: must return 401
curl -i -b cookies.txt localhost:8080/api/me