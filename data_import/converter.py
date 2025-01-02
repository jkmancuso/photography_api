from job import Job
from order import Order
import logging

class Converter:
    def __init__(self,instruments: list[dict], groups: list[dict], jobs: list[dict]):
        self.instruments={}
        self.groups={}
        self.jobs={}
        
        for row in instruments:
            self.instruments[row['instrument_id']]=row

        for row in groups:
            self.groups[row['group_id']]=row

        for row in jobs:
            self.jobs[row['job_id']]=row

    @staticmethod
    def convert_apijob_to_Job(row: dict) -> Job:
        return Job(row['id'],row['job_name'],int(row['job_year']))

    @staticmethod
    def convert_apiorder_to_Order(row: dict) -> Order:
        section={}
        
        #might be null
        if row['section'].get("name"):
            section={
                "name": row['section']['name'],
                "instrument": row['section']['instrument'],
                "quantity":int(row['section']['quantity']),
                "position":row['section']['position'],
                "picture_num":row['section']['picture_num']
                }
        
        return Order(row['id'],row['job_name'],row['job_id'],
                    int(row['record_num']),row['fname'],row['lname'],row['address'],
                    row['city'],row['state'],row['zip'],row['phone'],
                    int(row['group_quantity']),row['group'],row['group_picture_num'], 
                    int(row['check_num']),int(row['amount']),section)

    @staticmethod
    def convert_DBjob_to_Job(row: dict) -> Job:
        return Job(row['job_id'],row['job_name'],int(row['job_year']))

    def convert_DBorder_to_Order(self, row: dict) -> Order:

        job_name=self.jobs[row['job_id']]['job_name']

        if not row['group_quantity'] or row['group_quantity']==0:
            group=""
            group_quantity=0
            group_picturenum=""
        else:
            group=self.groups[int(row['group_id'])]['groupname']
            group_quantity=int(row['group_quantity'])
            group_picturenum=row['group_picturenum']
        
        section ={}

        if int(row['instrument_quantity'])>0:

            section={
                "name": self.instruments[row['instrument_id']]['section'],
                "instrument": self.instruments[row['instrument_id']]['instrument'],
                "quantity": int(row['instrument_quantity']),
                "position": row['instrument_position'],
                "picture_num": row['instrument_picturenum'],
                }


        return Order(row['customer_id'],job_name,row['job_id'],
                    int(row['record_num']),row['fname'],row['lname'],row['address'],
                    row['city'],row['state'],row['zip'],row['phone'],
                    group_quantity,group,group_picturenum, 
                    int(row['checknum']),int(row['amount']),section)

