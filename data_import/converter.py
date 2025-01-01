from job import Job
from order import Order

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
            self.groups[row['job_id']]=row

    @staticmethod
    def convert_apijob_to_Job(row: dict) -> Job:
        return Job(row['id'],row['job_name'],row['job_year'])

    @staticmethod
    def convert_apiorder_to_Order(row: dict) -> Order:
        
        section={
            "name": row['section']['name'],
            "instrument": row['section']['instrument'],
            "quantity":row['section']['quantity'],
            "position":row['section']['position'],
            "picture_num":row['section']['picture_num'],
            }
        
        return Order(row['id'],row['job_name'],row['job_id'],row['job_year'],
                    row['record_num'],row['fname'],row['lname'],row['address'],
                    row['city'],row['state'],row['zip'],row['phone'],
                    row['group_quantity'],row['group'],row['group_picture_num'], 
                    row['check_num'],row['amount'],section)

    @staticmethod
    def convert_DBjob_to_Job(row: dict) -> Job:
        return Job(row['job_id'],row['job_name'],row['job_year'])

    def convert_DBorder_to_Order(self, row: dict) -> Order:
        job_name=self.jobs[row['job_id']]['job_name']
        job_year=self.jobs[row['job_id']]['job_year']
        group=self.groups[row['group_id']]['groupname']

        section={
            "name": self.instruments[row['instrument_id']]['section'],
            "instrument": self.instruments[row['instrument_id']]['instrument'],
            "quantity": row['instrument_quantity'],
            "position": row['instrument_position'],
            "picture_num": row['instrument_picturenum'],
            }

        return Order(row['id'],job_name,row['job_id'],job_year,
                    row['record_num'],row['fname'],row['lname'],row['address'],
                    row['city'],row['state'],row['zip'],row['phone'],
                    row['group_quantity'],group,row['group_picture_num'], 
                    row['check_num'],row['amount'],section)

