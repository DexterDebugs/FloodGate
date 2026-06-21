#This is the backend that pretends to be real. 
# In production, a "users-service" would talk to a database, hash passwords, send emails. 
# Yours does one thing: receive a request, return mock JSON.
#request theeskoni JSON format return chesthadhi. backend idhi

from fastapi import FastAPI, Query  # type: ignore[import]
import time

app = FastAPI()

@app.get("/orders/{order_id}")
def get_order(order_id: int, slow: int = Query(0)):
    if (slow > 0):
        time.sleep(slow/1000);
    return {
        "id": order_id,
        "name": f"Order {order_id}",
        "service": "orders"
    }

