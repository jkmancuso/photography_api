import logging
import os
from dotenv import load_dotenv

from api import API
from order import Order
from job import Job
from db import DB
from converter import Converter




if __name__ == '__main__':

    load_dotenv()
    
    db = DB(
        os.getenv("MYSQL_HOST"),
        os.getenv("MYSQL_USER"),
        os.getenv("MYSQL_PASSWORD"),
        os.getenv("MYSQL_DB"))

    jobs = db.get_jobs()
    for job in jobs:
        orders = db.get_customer_orders_for_job(job)
    
    
    db.close()