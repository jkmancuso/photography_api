from job import Job
from order import Order

class Differ:

    @staticmethod
    def get_Jobs_not_in_Dynamo(jobs_already_in_api: list[Job],jobs_in_DB: list[Job]) -> list[Job]:
        jobnames_in_api=[]
        diff=[]

        for job in jobs_already_in_api:
            jobnames_in_api.append(job.job_name)

        for job in jobs_in_DB:
            if job.job_name not in jobnames_in_api:
                diff.append(job)
        
        return diff

    @staticmethod
    def get_JobOrders_not_in_Dynamo(orders_already_in_api: list[Order],orders_in_DB: list[Order]) -> list[Order]:
        rec_nums=[]
        max_rec_num_in_dynamo=0
        diff=[]

        for order in orders_already_in_api:
            rec_nums.append(order.record_num)

        max_rec_num_in_dynamo=max(rec_nums)

        for order in orders_in_DB:
            if order.record_num> max_rec_num_in_dynamo:
                diff.append(order)

        return diff