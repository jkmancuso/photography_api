class Order:
    def __init__(self,Id: int,JobName: str,JobId: int,JobYear: int, 
                 RecordNum: int,Fname: str,Lname: str,Address: str,City: str,State: str,
                 Zip: str, Phone: str,GroupQuantity: int,Group: int,GroupPictureNum: str,
                 CheckNum: int,Amount: int,Section: dict):
        self.id=Id
        self.job_name=JobName
        self.job_id=JobId
        self.job_year=JobYear
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