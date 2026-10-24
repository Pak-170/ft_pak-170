# Registro
curl -i -c cookies.txt -X POST localhost:8080/api/register \
  -d "username=ana&password=12345678&alias=anita"

# Sesión tras el registro
curl -i -b cookies.txt localhost:8080/api/me

# Logout y comprobación (204 y luego 401)
curl -i -b cookies.txt -X POST localhost:8080/api/logout
curl -i -b cookies.txt localhost:8080/api/me

# Login
curl -i -c cookies.txt -X POST localhost:8080/api/login \
  -d "username=ana&password=12345678"

# Sesión, logout final y comprobación (200, 204, 401)
curl -i -b cookies.txt localhost:8080/api/me
curl -i -b cookies.txt -X POST localhost:8080/api/logout
curl -i -b cookies.txt localhost:8080/api/me