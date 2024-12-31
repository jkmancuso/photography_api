import requests
import logging

class API:
    def __init__(self, url: str, x_session_id: str):
        self.url=url.rstrip('/') #just incase theres dangling /
        self.x_session_id=x_session_id
    
    def get_jobs(self)->list:
        r = requests.get(f"{self.url}/jobs")

        if r.status_code!=200:
            logging.error(r.json())
            exit(1)
        
        return r.json()

    def get_orders_for_job(self,job_name)->list:
        orders_dict={}
        r = requests.get(f"{self.url}/jobs/{self.jobs[job_name]}/orders")

        if r.status_code!=200:
            logging.error(r.json())
            exit(1)
        
        return r.json()