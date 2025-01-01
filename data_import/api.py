import requests
import logging
from job import Job
from order import Order

class API:

    logging.basicConfig(level=logging.INFO)

    def __init__(self, url: str, x_session_id: str):
        self.url=url.rstrip('/') #just incase theres dangling /
        self.x_session_id=x_session_id
        self.headers={
            "x-session-id": self.x_session_id,
            "Content-Type":"application/json"
        }
        self.jobs_lookup={}
    
    def get_jobs(self)->list[dict]:
        url=f"{self.url}/jobs"
        logging.info(f"HTTP GET {url} with headers {self.headers}")

        r = requests.get(url,headers=self.headers)

        json=r.json()
        
        if r.status_code!=200:
            logging.error(json)
            exit(1)
        
        for job in json:
            self.jobs_lookup[job['job_name']]=job['id']
        
        return json

    def get_orders_for_job(self,job_name)->list[dict]:
        url=f"{self.url}/jobs/{self.jobs_lookup[job_name]}/orders"

        logging.info(f"HTTP GET {url} with headers {self.headers}")
        r = requests.get(url,headers=self.headers)

        if r.status_code!=200:
            logging.error(r.json())
            exit(1)
        
        logging.info(r.json())
        
        return r.json()
    
    def post_jobs(self,jobs: list[Job]):
        url=f"{self.url}/jobs"

        for job in jobs:
            
            job_data={
                'job_name': job.job_name,
                'job_year': job.job_year
            }

            logging.info(f"HTTP POST {url}")
            logging.info(f"with headers {self.headers}")
            logging.info(f"with payload {job_data}")
            
            r = requests.post(url,json=job_data,headers=self.headers)        

            if r.status_code!=200:
                logging.error(f"Request:{r.request.body}")
                logging.error(f"{r.status_code}: {r.json()}")
                exit(1)
            
            logging.info("SUCCESS")
    
    def post_orders(self,orders: list[Order]):
        url=f"{self.url}/orders"

        for order in orders:
            
            order_data={
                "job_name": order.job_name,
                "job_id": order.job_id,
                "record_num": order.record_num,
                "fname": order.fname,
                "lname": order.lname,
                "address":order.address,
                "city":order.city,
                "state":order.state,
                "zip":order.zip,
                "phone":order.phone,
                "group_quantity":order.group_quantity,
                "group":order.group,
                "group_picture_num":order.group_picture_num,
                "check_num":order.check_num,
                "amount":order.amount,
                "section":order.section
            }

            logging.info(f"HTTP POST {url}")
            logging.info(f"with headers {self.headers}")
            logging.info(f"with payload {order_data}")
            
            r = requests.post(url,json=order_data,headers=self.headers)        

            if r.status_code!=200:
                logging.error(f"Request:{r.request.body}")
                logging.error(f"{r.status_code}: {r.json()}")
                exit(1)
            
            logging.info("SUCCESS")

        