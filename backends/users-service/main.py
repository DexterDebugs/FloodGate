#This is the backend that pretends to be real. 
# In production, a "users-service" would talk to a database, hash passwords, send emails. 
# Yours does one thing: receive a request, return mock JSON.
#request theeskoni JSON format return chesthadhi. backend idhi

from fastapi import FastAPI

app = FastAPI()

@app.get("/users/{user_id}")
def get_user(user_id: int):
    return {
        "id": user_id,
        "name": f"User {user_id}",
        "service": "users"
    }