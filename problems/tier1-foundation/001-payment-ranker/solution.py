"""Payment ranker — Strategy pattern reference solution (Python)."""

from abc import ABC, abstractmethod


class PaymentMethod:
    def __init__(self, name, cashbackRate, transactionFee, usageCount,
                 easyRefundEligible=False):
        self.name = name
        self.cashbackRate = cashbackRate
        self.transactionFee = transactionFee
        self.usageCount = usageCount
        self.easyRefundEligible = easyRefundEligible


class RankingStrategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` should rank strictly before `b`."""


class RewardsMaximizer(RankingStrategy):
    def compare(self, a, b):
        return a.cashbackRate > b.cashbackRate


class LowFeeSeeker(RankingStrategy):
    def compare(self, a, b):
        return a.transactionFee < b.transactionFee


class TrustBasedRanker(RankingStrategy):
    def compare(self, a, b):
        return a.usageCount > b.usageCount


class CompositeStrategy(RankingStrategy):
    def __init__(self, criteria):
        self.criteria = criteria

    def compare(self, a, b):
        for s in self.criteria:
            if s.compare(a, b):
                return True
            if s.compare(b, a):
                return False
        return False


class EasyRefundStrategy(RankingStrategy):
    def __init__(self, prefer):
        self.prefer = prefer

    def compare(self, a, b):
        if not self.prefer:
            return False
        return a.easyRefundEligible and not b.easyRefundEligible


class PaymentRanker:
    def __init__(self, strategy):
        self.strategy = strategy

    def set_strategy(self, strategy):
        self.strategy = strategy

    def rank(self, methods):
        # Translate strict-less comparator into a stable sort via insertion.
        items = list(methods)
        result = []
        for m in items:
            inserted = False
            for i, existing in enumerate(result):
                if self.strategy.compare(m, existing):
                    result.insert(i, m)
                    inserted = True
                    break
            if not inserted:
                result.append(m)
        return result


def rank_by_rewards(methods):
    return PaymentRanker(RewardsMaximizer()).rank(methods)


def rank_by_low_fee(methods):
    return PaymentRanker(LowFeeSeeker()).rank(methods)


def rank_by_trust(methods):
    return PaymentRanker(TrustBasedRanker()).rank(methods)


def rank_composite(methods, criteria):
    return PaymentRanker(CompositeStrategy(criteria)).rank(methods)


def rank_with_refund_filter(methods, preferEasyRefund):
    composite = CompositeStrategy([
        EasyRefundStrategy(preferEasyRefund),
        RewardsMaximizer(),
    ])
    return PaymentRanker(composite).rank(methods)
