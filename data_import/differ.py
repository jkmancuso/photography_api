from job import Job
from order import Order

class Differ:

    @staticmethod
    def get_Jobs_not_in_Dynamo(jobs_already_in_api: list[Job],jobs_in_DB: list[Job]) -> list[Job]:
        pass

    @staticmethod
    def get_JobOrders_not_in_Dynamo(orders_already_in_api: list[Order],orders_in_DB: list[Order]) -> list[Order]:
        pass
