# demo-app

```bash
# Get todos (as user)
curl -u user:userpass http://localhost:8080/api/todos

# Add todo (as admin)
curl -u admin:adminpass -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"text":"Buy milk"}'

# Delete todo with id 1 (as admin)
curl -u admin:adminpass -X DELETE http://localhost:8080/api/todos/1

# Unauthorized add todo (as user, should return 403)
curl -u user:userpass -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"text":"Should fail"}'
```
