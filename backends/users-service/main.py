#This is the backend that pretends to be real. 
# In production, a "users-service" would talk to a database, hash passwords, send emails. 
# Yours does one thing: receive a request, return mock JSON.
#request theeskoni JSON format return chesthadhi. backend idhi

from fastapi import FastAPI, Query # type: ignore[import]
import time

app = FastAPI()

@app.get("/users/{user_id}")
def get_user(user_id: int, slow: int = Query(0)):   #Query(0) declares and optional query parameter with default as 0
    if slow > 0:
        time.sleep(slow/1000)
    return {
        "id": user_id,
        "name": f"User {user_id}",
        "service": "users"
    }