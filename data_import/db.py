
import logging
import mysql.connector

class DB:
    def __init__(self, host: str, user: str, password: str, dbname: str):
        self.host=host
        self.user=user
        self.password=password
        self.dbname=dbname

        if not self.host or not self.user or not self.password or not self.dbname:
            logging.error("Missing params")
            exit(1)
        
        self.connect()
        
        self.jobs=self.get_jobs()
        self.instruments=self.get_instruments()
        self.groups=self.get_groups()

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
    
    def get_jobs(self) -> list[dict]:
        query="select * from jobs"
        self.cursor.execute(query)
        
        return self.cursor.fetchall()
        
    def get_instruments(self) -> list[dict]:
        query="select * from instruments"
        self.cursor.execute(query)
        
        return self.cursor.fetchall()
        
    def get_groups(self) -> list[dict]:
        query="select * from groups"
        self.cursor.execute(query)

        return self.cursor.fetchall()
        
    def get_customer_orders_by_job_id(self, id) -> list[dict]:
        
        query=f"select * from customers where job_id = %s"
        self.cursor.execute(query, (id,))
        
        return self.cursor.fetchall()