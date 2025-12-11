## Endpoints To Test

* Gateway health: `http://localhost:8080/health`

* Auth service: `http://localhost:8081/health`, register, login, refresh, profile

* Verification service health: `http://localhost:8082/health`

* ML service: `http://localhost:8000/health`, `http://localhost:8000/api/v1/predict`

* Note: Gateway does not expose ML predict route (returns 404)

## Windows (PowerShell) cURL

* Health:

  * `C:\Windows\System32\curl.exe -s -o NUL -w "%{http_code}\n" http://localhost:8080/health`

  * `C:\Windows\System32\curl.exe -s -o NUL -w "%{http_code}\n" http://localhost:8081/health`

  * `C:\Windows\System32\curl.exe -s -o NUL -w "%{http_code}\n" http://localhost:8082/health`

  * `C:\Windows\System32\curl.exe -s -o NUL -w "%{http_code}\n" http://localhost:8000/health`

* Register:

  * `C:\Windows\System32\curl.exe -s -H "Content-Type: application/json" --data-binary @e:\CODES\JPMCProj\fraud-detection-system\tmp\reg3.json http://localhost:8081/api/v1/auth/register`

* Login:

  * `C:\Windows\System32\curl.exe -s -H "Content-Type: application/json" --data-binary @e:\CODES\JPMCProj\fraud-detection-system\tmp\login3.json http://localhost:8081/api/v1/auth/login`

* Profile:

  * `C:\Windows\System32\curl.exe -s -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8081/api/v1/profile`

* Refresh token:

  * `C:\Windows\System32\curl.exe -s -H "Content-Type: application/json" -d "{\"refresh_token\":\"<REFRESH_TOKEN>\"}" http://localhost:8081/api/v1/auth/refresh`

* ML predict:

  * `C:\Windows\System32\curl.exe -s -H "Content-Type: application/json" --data-binary @e:\CODES\JPMCProj\fraud-detection-system\tmp\predict.json http://localhost:8000/api/v1/predict`

* Gateway ML route check (expected 404):

  * `C:\Windows\System32\curl.exe -s -o NUL -w "%{http_code}\n" http://localhost:8080/api/v1/ml/predict`

## WSL (Ubuntu) cURL

* Health:

  * `curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/health`

  * `curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8081/health`

  * `curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8082/health`

  * `curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8000/health`

* Register:

  * `curl -s -H 'Content-Type: application/json' --data-binary @/mnt/e/CODES/JPMCProj/fraud-detection-system/tmp/reg3.json http://localhost:8081/api/v1/auth/register`

* Login:

  * `curl -s -H 'Content-Type: application/json' --data-binary @/mnt/e/CODES/JPMCProj/fraud-detection-system/tmp/login3.json http://localhost:8081/api/v1/auth/login`

* Capture tokens (requires `jq`):

  * `TOKEN=$(curl -s -H 'Content-Type: application/json' --data-binary @/mnt/e/CODES/JPMCProj/fraud-detection-system/tmp/login3.json http://localhost:8081/api/v1/auth/login | jq -r '.data.access_token')`

  * `REFRESH=$(curl -s -H 'Content-Type: application/json' --data-binary @/mnt/e/CODES/JPMCProj/fraud-detection-system/tmp/login3.json http://localhost:8081/api/v1/auth/login | jq -r '.data.refresh_token')`

* Profile:

  * `curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8081/api/v1/profile`

* Refresh token:

  * `curl -s -H 'Content-Type: application/json' -d '{"refresh_token":"'$REFRESH'"}' http://localhost:8081/api/v1/auth/refresh`

* ML predict:

  * `curl -s -H 'Content-Type: application/json' --data-binary @/mnt/e/CODES/JPMCProj/fraud-detection-system/tmp/predict.json http://localhost:8000/api/v1/predict`

* Gateway ML route check:

  * `curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/api/v1/ml/predict`

## Notes

* Use file-based bodies to avoid PowerShell quoting issues (`--data-binary @file`).

* Replace `<ACCESS_TOKEN>` and `<REFRESH_TOKEN>` with values returned from login in Windows; in WSL use `jq` to capture tokens.

* pgAdmin GUI is at `http://localhost:5050` if you prefer browsing tables visually.

