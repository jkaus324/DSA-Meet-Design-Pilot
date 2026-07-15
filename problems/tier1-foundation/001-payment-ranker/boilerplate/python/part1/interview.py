# Data class (given).
class PaymentMethod:
    def __init__(self, name, cashbackRate, transactionFee, usageCount, easyRefundEligible=False):
        self.name = name
        self.cashbackRate = cashbackRate
        self.transactionFee = transactionFee
        self.usageCount = usageCount
        self.easyRefundEligible = easyRefundEligible

def rank_by_rewards(methods):
    # TODO: write your solution
    return methods

def rank_by_low_fee(methods):
    # TODO: write your solution
    return methods

def rank_by_trust(methods):
    # TODO: write your solution
    return methods
