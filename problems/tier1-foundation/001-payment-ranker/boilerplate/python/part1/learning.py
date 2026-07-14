# Data class (given — do not modify).
class PaymentMethod:
    def __init__(self, name, cashbackRate, transactionFee, usageCount, easyRefundEligible=False):
        self.name = name
        self.cashbackRate = cashbackRate
        self.transactionFee = transactionFee
        self.usageCount = usageCount
        self.easyRefundEligible = easyRefundEligible

from abc import ABC, abstractmethod

class RankingStrategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

class RewardsMaximizer(RankingStrategy):
    def compare(self, a, b):
        # TODO: implement this
        return False

class LowFeeSeeker(RankingStrategy):
    def compare(self, a, b):
        # TODO: implement this
        return False

class TrustBasedRanker(RankingStrategy):
    def compare(self, a, b):
        # TODO: implement this
        return False

def rank_by_rewards(methods):
    # TODO: implement this
    return methods

def rank_by_low_fee(methods):
    # TODO: implement this
    return methods

def rank_by_trust(methods):
    # TODO: implement this
    return methods
