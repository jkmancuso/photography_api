import logging
import os
from dotenv import load_dotenv

from api import API
from order import Order
from job import Job
from db import DB
from converter import Converter
from differ import Differ



if __name__ == '__main__':

    load_dotenv()
    
    db = DB(
        os.getenv("MYSQL_HOST"),
        os.getenv("MYSQL_USER"),
        os.getenv("MYSQL_PASSWORD"),
        os.getenv("MYSQL_DB"))
    
    api=API(os.getenv("API_URL"),os.getenv("SESSION_ID"))
    
    db_jobs=db.get_jobs()
    api_jobs=api.get_jobs()

    convert=Converter(db.get_instruments(),db.get_groups(),db_jobs)
    
    ######DO JOBS FIRST######
    jobs_in_DB=[convert.convert_DBjob_to_Job(job) for job in db_jobs]

    jobs_already_in_api=[convert.convert_apijob_to_Job(job) for job in api_jobs]
    
    missing_jobs=Differ.get_Jobs_not_in_Dynamo(jobs_already_in_api,jobs_in_DB)
    
    if len(missing_jobs)>0:
        confirm=input("Add the missing Jobs to the DB? [Y/N]")
        print(confirm)
    
    #########################

    ##NOW ORDERS#############
    
    db.close()