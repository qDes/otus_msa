# HW1
simple app with one http handler <br>
```
GET /health/
RESPONSE: {"status": "OK"}
```

## build
```
docker build -t health .
```

## run 
```
docker run -p 8000:8000 health
```

check result on http://127.0.0.1:8000/health