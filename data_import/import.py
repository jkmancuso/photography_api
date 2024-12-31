import logging
import os
import datetime
import math
import uuid
import mysql.connector
from dotenv import load_dotenv


class DB:
    def __init__(self, host: str, user: str, password: str, dbname: str):
        self.host=host
        self.user=user
        self.password=password
        self.dbname=dbname

        if not self.host or not self.user or not self.password or not self.dbname:
            logging.error("Missing params")
        
        self.connect()

    def connect(self):
        try:
            self.conn = mysql.connector.connect(
                user=self.user,
                password=self.password,
                host=self.host,
                database=self.dbname)

        except mysql.connector.Error as err:
            logging.error(f"ERROR: {err}")
            logging.error(f"""Conn: 
                            user: {self.user} 
                            password: {self.password}
                            host: {self.host}
                            DB: {self.dbname}"""
                            )
            exit(1)

        self.cursor=self.conn.cursor(dictionary=True)
    
    def close(self):
        self.conn.close()
    
    def get_jobs(self) -> dict:
        jobs_dict={}
        query="select * from jobs"
        self.cursor.execute(query)
        
        #create a dict so that you can loop over the customer orders and
        #lookup by their job id
        for row in self.cursor:
            jobs_dict[row['job_id']] = Job(row['job_name'],row['job_year'],str(uuid.uuid4()))
        
        return jobs_dict
    
    def get_instruments(self) -> dict:
        instruments_dict={}
        query="select * from instruments"
        self.cursor.execute(query)
        
        for row in self.cursor:
            instruments_dict[row['instrument_id']] = Job(row['instrument'],row['section'])
        
        return instruments_dict
    
    def get_groups(self) -> dict:
        groups_dict={}
        query="select * from groups"
        self.cursor.execute(query)
        
        for row in self.cursor:
            groups_dict[row['group_id']] = Job(row['groupname'])
        
        return groups_dict
    
    
    def get_customer_orders(self, jobs_dict: dict[int],
                            instruments_dict: dict[int], 
                            groups_dict: dict[int]) -> dict:
        orders_dict={}
        query="select * from customers where doe>= %s"
        self.cursor.execute(query, datetime.date(2024,1,1))

        for row in self.cursor:
            section={
                "name": instruments_dict[row['instrument_id']]['section'],
                "instrument": instruments_dict[row['instrument_id']]['instrument'],
                "instrument_quantity":row['instrument_quantity'],
                "instrument_position":row['instrument_position']
            }
            orders_dict[row['customer_id']]=NewOrder(str(uuid.uuid4()),
                    jobs_dict[row['job_id']]['job_name'],row['record_num'],
                    row['fname'],row['lname'],row['address'],row['city'],
                    row['state'],row['zip'],row['phone'],row['group_quantity'],
                    groups_dict[row['group_id']]['groupname'],row['group_picturenum'], 
                    row['checknum'],row['amount'],section)
        
        return orders_dict
    
    #Goal: Determine the orders delta (a list of customer ids) we need to sync 
    # from old MySQL customers DB -> new Dynamo orders DB
     
    
    #Step 1: get the most recent update on the old customer DB, check the auto increment id
    def get_highest_customer_id(self) ->int:
        pass

    #Step 2: get the first entry on the old customer DB, (in Jan 01 2024)
    def get_lowest_customer_id(self) ->int:
        pass

    #Step 3: lets make this a binary search problem! Rather than download the entire dynamo 
    # order records, in MySQL we have the highest customer id and we have the lowest possible 
    # customer id (in Jan 01 2024). Thats our range. 
    # Goal: return the highest customer id we have associated with a dynamodb order
    def get_highest_customer_order_in_dynamo(self, low: int, high: int) -> int:
        
        #in this case, dynamo is empty. Send everything (ie start from lower bound)
        if not self.customer_id_exists_in_dynamo(low):
            return low
        
        #base case
        if low >= high or low == high-1:
            return min(low,high)

        mid=math.floor((high-low)/2)+low
        
        if self.customer_id_exists_in_dynamo(mid+1):
            return self.get_highest_customer_order_in_dynamo(mid+1,high)
        else:
            return self.get_highest_customer_order_in_dynamo(low,mid)
        
    def customer_id_exists_in_dynamo(self) ->bool:
        pass

class Job():
    def __init__(self,job_name: str,job_year: str,uuid: str):
        self.job_name=job_name
        self.job_year=job_year
        self.new_uuid=uuid

class NewOrder():
    def __init__(self,NewId: int,JobName: str,JobId: int,RecordNum: int,Fname: str,
                 Lname: str,Address: str,City: str,State: str,Zip: str,
                 Phone: str,GroupQuantity: int,Group: int,GroupPictureNum: str,
                 CheckNum: int,Amount: int,Section: dict):
        self.new_id=NewId
        self.job_name=JobName
        self.job_id=JobId
        self.record_num=RecordNum
        self.fname=Fname
        self.lname=Lname
        self.address=Address
        self.city=City
        self.state=State
        self.zip=Zip
        self.phone=Phone
        self.group_quantity=GroupQuantity
        self.group=Group
        self.group_picture_num=GroupPictureNum
        self.check_num=CheckNum
        self.amount=Amount
        self.section=Section


if __name__ == '__main__':

    load_dotenv()
    
    db = DB(
        os.getenv("MYSQL_HOST"),
        os.getenv("MYSQL_USER"),
        os.getenv("MYSQL_PASSWORD"),
        os.getenv("MYSQL_DB"))

    jobs = db.GetJobs()
    
    db.close()