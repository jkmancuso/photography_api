import math

def binary_search(low, high):

    if not exists(low):
        return low
    
    #base case
    if low >= high or low == high-1:
        return min(low,high)

    mid=math.floor((high-low)/2)+low
    
    if exists(mid+1):
        print(f"Searching {mid+1}-{high}")
        return binary_search(mid+1,high)
    else:
        print(f"Searching {low}-{mid}")
        return binary_search(low,mid)


# simulate DB check
def exists(num):
    if num <=5:
        return True
    else:
        return False

print(binary_search(1,100))

'''
1…..100
1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22

check 1…yes (if no, then you are done, sync all)
next number to check: 50
check 50: no
next number to check: 25
check 25: no
next number to check 12
check 12: yes
next number to check: 19
check 19: yes
next number to check: 23
check 23: no
next number to check: 21 
check 21: yes
next number to check: 22
check 22: yes
if low bound is 1 less  then upper bound you have your inflection point '''