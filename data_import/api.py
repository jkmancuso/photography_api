import requests
import logging

class API:

    logging.basicConfig(level=logging.INFO)

    def __init__(self, url: str, x_session_id: str):
        self.url=url.rstrip('/') #just incase theres dangling /
        self.x_session_id=x_session_id
        self.headers={
            "x-session-id": self.x_session_id,
            "Content-Type":"application/json"
        }
    
    def get_jobs(self)->list[dict]:
        url=f"{self.url}/jobs"
        logging.info(f"HTTP GET {url} with headers {self.headers}")

        r = requests.get(url,headers=self.headers)

        if r.status_code!=200:
            logging.error(r.json())
            exit(1)
        
        return r.json()

    def get_orders_for_job(self,job_name)->list[dict]:
        orders_dict={}
        url=f"{self.url}/jobs/{self.jobs[job_name]}/orders"

        logging.info(f"HTTP GET {url} with headers {self.headers}")
        r = requests.get(url,headers=self.headers)

        if r.status_code!=200:
            logging.error(r.json())
            exit(1)
        
        return r.json()