# Data class (given).
class PaymentMethod:
    def __init__(self, name, cashbackRate, transactionFee, usageCount, easyRefundEligible=False):
        self.name = name
        self.cashbackRate = cashbackRate
        self.transactionFee = transactionFee
        self.usageCount = usageCount
        self.easyRefundEligible = easyRefundEligible

# HINT: introduce an abstraction so new ranking rules don't change existing code.

# HINT: pick the field that defines 'better' for this ranking and compare the two.
def rank_by_rewards(methods):
    # TODO: write your solution
    return methods

# HINT: pick the field that defines 'better' for this ranking and compare the two.
def rank_by_low_fee(methods):
    # TODO: write your solution
    return methods

# HINT: pick the field that defines 'better' for this ranking and compare the two.
def rank_by_trust(methods):
    # TODO: write your solution
    return methods
